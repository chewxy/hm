# hm [![Build Status](https://travis-ci.org/chewxy/hm.svg?branch=master)](https://travis-ci.org/chewxy/hm)

Package hm is a simple Hindley-Milner type inference system in Go. It provides the necessary data structures and functions for creating such a system. 

# Installation #

This package is go-gettable: `go get -u github.com/chewxy/hm`

There are very few dependencies that this package uses. Therefore there isn't a need for vendoring tools.

Here is a listing of the dependencies of `hm`:

|Package|Used For|Vitality|Notes|Licence|
|-------|--------|--------|-----|-------|
|[errors](https://github.com/pkg/errors)|Error wrapping|Gorgonia won't die without it.|Stable API for the past 6 months|[errors licence](https://github.com/pkg/errors/blob/master/LICENSE) (MIT/BSD-like)|
|[testify/assert](https://github.com/stretchr/testify)|Testing|Can do without but will be a massive pain in the ass to test||[testify licence](https://github.com/stretchr/testify/blob/master/LICENSE) (MIT/BSD-like)|

# Usage

TODO: Write this


# Contributing

This library is developed using Github. Therefore the workflow is very github-centric. 

# Licence

Package `hm` is licenced under the MIT licence.