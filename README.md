# Singularity

[![Travis status](https://travis-ci.org/brandur/singularity.svg?branch=master)](https://travis-ci.org/brandur/singularity)

A demonstration of a very simple static site generator that deploys to S3
through Travis CI.

## Build

Build with:

    make

The results will appear in `public/`.

## Deployment

The repository will deploy to S3 automatically from the Travis build when
changes are pushed to master.

This works by having [encrypted variables][travis-encrypted] configured in
`.travis.yml` for an IAM account with privileges to the production S3 bucket.
These credentials can be reconfigured with:

    gem install travis
    travis encrypt AWS_ACCESS_KEY_ID=...
    travis encrypt AWS_SECRET_ACCESS_KEY=...

### Locally

Deploy locally by first making sure that you have awscli installed:

    pip install awscli

Then set appropriate AWS keys and go for it:

    export AWS_ACCESS_KEY_ID=...
    export AWS_SECRET_ACCESS_KEY=...
    export S3_BUCKET=singularity.brandur.org
    make deploy

[travis-encrypted]: https://docs.travis-ci.com/user/environment-variables/#Encrypted-Variables
