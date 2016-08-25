all: clean install test vet lint check-gofmt build

build:
	$(GOPATH)/bin/singularity-build

check-gofmt:
	scripts/check_gofmt.sh

clean:
	mkdir -p public/
	rm -f -r public/*

# Long TTL (in seconds) to set on an object in S3. This is suitable for items
# that we expect to only have to invalidate very rarely like images. Although
# we set it for all assets, those that are expected to change more frequently
# like script or stylesheet files are versioned by a path that can be set at
# build time.
LONG_TTL := 86400

# Short TTL (in seconds) to set on an object in S3. This is suitable for items
# that are expected to change more frequently like any HTML file.
SHORT_TTL := 3600

deploy: build check-target-dir
# Note that AWS_ACCESS_KEY_ID will only be set for builds on the master
# branch because it's stored in `.travis.yml` as an encrypted variable.
# Encrypted variables are not made available to non-master branches because
# of the risk of being leaked through a script in a rogue pull request.
ifdef AWS_ACCESS_KEY_ID
	aws --version

	# Force text/html for HTML because we're not using an extension.
	#
	# Note that we don't delete because it could result in a race condition in
	# that files that are uploaded with special directives below could be
	# removed even while the S3 bucket is actively in-use.
	aws s3 sync $(TARGET_DIR) s3://$(S3_BUCKET)/ --acl public-read --cache-control max-age=$(SHORT_TTL) --content-type text/html --exclude 'assets*' $(AWS_CLI_FLAGS)

	# Then move on to assets and allow S3 to detect content type.
	aws s3 sync $(TARGET_DIR)/assets/ s3://$(S3_BUCKET)/assets/ --acl public-read --cache-control max-age=$(LONG_TTL) --delete --follow-symlinks $(AWS_CLI_FLAGS)

	# Upload Atom feed files with their proper content type.
	find $(TARGET_DIR) -name '*.atom' | sed "s|^\$(TARGET_DIR)/||" | xargs -I{} -n1 aws s3 cp $(TARGET_DIR)/{} s3://$(S3_BUCKET)/{} --acl public-read --cache-control max-age=$(SHORT_TTL) --content-type application/xml

	# This one is a bit tricker to explain, but what we're doing here is
	# uploading directory indexes as files at their directory name. So for
	# example, 'articles/index.html` gets uploaded as `articles`.
	#
	# We do all this work because CloudFront/S3 has trouble with index files.
	# An S3 static site can have index.html set to indexes, but CloudFront only
	# has the notion of a "root object" which is an index at the top level.
	#
	# We do this during deploy instead of during build for two reasons:
	#
	# 1. Some directories need to have an index *and* other files. We must name
	#    index files with `index.html` locally though because a file and
	#    directory cannot share a name.
	# 2. The `index.html` files are useful for emulating a live server locally:
	#    Golang's http.FileServer will respect them as indexes.
	find $(TARGET_DIR) -name index.html | egrep -v '$(TARGET_DIR)/index.html' | sed "s|^$(TARGET_DIR)/||" | xargs -I{} -n1 dirname {} | xargs -I{} -n1 aws s3 cp $(TARGET_DIR)/{}/index.html s3://$(S3_BUCKET)/{} --acl public-read --cache-control max-age=$(SHORT_TTL) --content-type text/html
else
	# No AWS access key. Skipping deploy.
endif

install:
	go install $(shell go list ./... | grep -v '/vendor/')

lint:
	$(GOPATH)/bin/golint -set_exit_status

save:
	godep save $(shell go list ./... | grep -v '/vendor/')

serve:
	$(GOPATH)/bin/singularity-serve

test:
	go test $(shell go list ./... | grep -v '/vendor/')

vet:
	go vet $(shell go list ./... | grep -v '/vendor/')

watch:
	fswatch -o articles/ assets/ layouts/ pages/ | xargs -n1 -I{} make build

# This is designed to be compromise between being explicit and readability. We
# can allow the find to discover everything in vendor/, but then the fswatch
# invocation becomes a huge unreadable wall of text that gets dumped into the
# shell. Instead, find all our own *.go files and then just tack the vendor/
# directory on separately (fswatch will watch it recursively).
GO_FILES := $(shell find . -type f -name "*.go" ! -path "./vendor/*")

# We recompile our Go source when a file changes, but we also rebuild the site
# because a change in source may have affected the build formula.
watch-go:
	fswatch -o $(GO_FILES) vendor/ | xargs -n1 -I{} make install build

#
# Helpers
#

check-target-dir:
ifndef TARGET_DIR
	$(error TARGET_DIR is required)
endif
