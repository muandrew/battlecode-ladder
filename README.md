# Battlecode Legacy

This is a project to make off season Battlecode easier to play. The plan is to support multiple years worth of Battlecode but right now only 2017 is supported.

# Usage

## Setup
* java 8
* golang 1.8
* redis 3.2.4
* unzip
* source [SOURCE\_ME.sh](./SOURCE\_ME.sh) in your `.bash_profile`/`.bash_rc`
* copy [example-bcl-env.sh](./go/src/github.com/muandrew/battlecode-legacy-go/example-bcl-env.sh) to your own `bcl_env.sh`.

**note:** The version numbers are the ones I'm using, and does not mean the exact version is required.

## Running
1. run `debug_run.sh`

## Developing
1. Do yourself a favor and download [GoLand EAP](https://www.jetbrains.com/go/nextversion/) or if you have a license just use [GoLand](https://www.jetbrains.com/go/).
2. With GoLand, open `./go/src/github.com/muandrew/battlecode-legacy-go`.

**note:** I normally develop with Intellij so it just makes more sense to me. You can use whatever you want.

## GraphQL[wip]
* check out ChromeiQL or other out of the box solutions for an easy way to test the GraphQL endpoint.

## Warnings
* Running an initial match is pretty ram heavy, I tried with a 512mb box and the compilation failed. This could also be because I'm using Redis.
* I would not use any coding patterns from this project as a template for your future projects. I'm just derping around.
