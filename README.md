# Singularity

[![Travis status](https://travis-ci.org/brandur/singularity.svg?branch=master)](https://travis-ci.org/brandur/singularity)

A demonstration of a very simple static site generator that deploys to S3
through Travis CI.

## Build

Build with:

    make

The results will appear in `public/`.

## Deployment

Deploy locally by first making sure that you have awscli installed:

    pip install awscli

Then set appropriate AWS keys and go for it:

    export AWS_ACCESS_KEY_ID=...
    export AWS_SECRET_ACCESS_KEY=...
    export S3_BUCKET=singularity.brandur.org
    make deploy

### Travis

Deploy from Travis by configuring it to use [encrypted
versions][travis-encrypted] of your AWS keys:

    gem install travis
    travis encrypt AWS_ACCESS_KEY_ID=...
    travis encrypt AWS_SECRET_ACCESS_KEY=...

Add the results in the `env` section of `.travis.yml` and `git push origin
master`.

[travis-encrypted]: https://docs.travis-ci.com/user/environment-variables/#Encrypted-Variables
