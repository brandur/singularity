default: lint test

lint:
	jslint *.js

test:
	lambda-local -l index.js -h handler -e event.js
