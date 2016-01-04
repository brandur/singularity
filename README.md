# Singularity

[![Travis status](https://travis-ci.org/brandur/singularity.svg?branch=master)](https://travis-ci.org/brandur/singularity)

Build with:

    make

The results will appear in `public/`.

## Deployment

[documentation on encrypted variables][travis-encrypted]

    gem install travis
    travis encrypt AWS_ACCESS_KEY=...
    travis encrypt AWS_SECRET_KEY=...

Add the results in the `env` section of `.travis.yml`.

[travis-encrypted]: https://docs.travis-ci.com/user/environment-variables/#Encrypted-Variables
