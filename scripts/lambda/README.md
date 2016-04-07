## Travis Rebuild Script

These are two tiny scripts designed for use with AWS Lambda that will trigger a
rebuild on a Travis repository. The idea here is to put this on a schedule so
that if there isn't a lot of activity on a source repository we'll find out if
something goes wrong with the CI process (say the credentials expire for
example).

* `index.js` contains a Lambda script to trigger a rebuild.
* `event.js` is a test "event" for use with Lambda.

`index.js` must be configured with a repository and a Travis token which can be
obtained via CLI:

    gem install travis
    travis login --org
    travis token

The scripts can be tested locally with Node's lambda-local package:

    npm install -g lambda-local
    lambda-local -l index.js -h handler -e event.js
