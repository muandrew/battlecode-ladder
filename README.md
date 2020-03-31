# Battlecode Legacy

This is a project to make off season Battlecode easier to play. The plan is to support multiple years worth of Battlecode but right now only 2017 is supported.

# Usage

All scripts are assumed to be ran from the base directory

## Setup
* java 8
* golang 1.13.5
* redis 3.2.4
* all tools in [check_dependency.sh](./check_dependency.sh)
* source [SOURCEME.sh](./SOURCEME.sh) in your `.bash_profile`/`.bash_rc`
* copy [example-bcl-env.sh](./go/src/github.com/muandrew/battlecode-legacy-go/example-bcl-env.sh) to your own `bcl_env.sh`.

**note:** The version numbers are the ones I'm using, and does not mean the exact version is required.

## Development
* use `precommit.sh` for any formatting that should be done before comitting
* use `start_bcl.sh` to just build and start the service.

## Deployment
1. run `deploy.sh` this should try to start up any services that are needed

## GraphQL[wip]
* check out ChromeiQL or other out of the box solutions for an easy way to test the GraphQL endpoint.

## Warnings
* Running an initial match is pretty ram heavy, I tried with a 512mb box and the compilation failed. This could also be because I'm using Redis.
* I would not use any coding patterns from this project as a template for your future projects. I'm just derping around.
