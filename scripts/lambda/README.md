## Travis Rebuild Script

These are two tiny scripts designed for use with AWS Lambda that will trigger a
rebuild on a Travis repository. The idea here is to put this on a schedule so
that if there isn't a lot of activity on a source repository we'll find out if
something goes wrong with the CI process (say the credentials expire for
example).

* `index.js` contains a Lambda script to trigger a rebuild.
* `event.js` is a test "event" for use with Lambda.

You'll need a Travis token to run `index.js`. One can be obtained via their
CLI:

    gem install travis
    travis login --org
    travis token

The scripts can also be tested locally with Node's lambda-local package:

    npm install -g jslint
    npm install -g lambda-local
    REPOSITORY=brandur/singularity TRAVIS_TOKEN= make
