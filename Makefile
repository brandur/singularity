# The name of the Go executable.
bin = singularity

all: clean compile test vet lint check-gofmt build

build:
	./$(bin)

check-gofmt:
	scripts/check_gofmt.sh

clean:
	mkdir -p public/
	rm -f -r public/*

compile:
	go build -o $(bin)

deploy: build
# Note that AWS_ACCESS_KEY_ID will only be set for builds on the master
# branch because it's stored in `.travis.yml` as an encrypted variable.
# Encrypted variables are not made available to non-master branches because
# of the risk of being leaked through a script in a rogue pull request.
ifdef AWS_ACCESS_KEY_ID
	aws --version

	# Force text/html for HTML because we're not using an extension.
	aws s3 sync ./public/ s3://$(S3_BUCKET)/ --acl public-read --content-type text/html --delete --exclude 'assets*' $(AWS_CLI_FLAGS)

	# Then move on to assets and allow S3 to detect content type.
	aws s3 sync ./public/assets/ s3://$(S3_BUCKET)/assets/ --acl public-read --delete --follow-symlinks $(AWS_CLI_FLAGS)
endif

lint:
	# Use a full path here for the benefit of Travis.
	$(GOPATH)/bin/golint

	# Hack to workaround the fact that Golint doesn't produce a non-zero exit
	# code on failure because Go Core team is always right and everyone else is
	# always wrong:
	#
	#     https://github.com/golang/lint/issues/65
	#
	test -z "$$(golint .)"

save-deps:
	godep save ./...

serve:
	./$(bin) serve

test:
	go test

vet:
	go vet

watch:
	fswatch -o articles/ assets/ layouts/ pages/ | xargs -n1 -I{} make build
