# Makefile
#
# You'll need a few dependencies to run these tasks. See the README.
#

default: lint test

get-deps:
	npm install -g jslint
	npm install -g lambda-local

lint:
	jslint *.js

test:
	lambda-local -l index.js -h handler -e event.js
