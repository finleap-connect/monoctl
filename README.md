# `monoctl`

[![Build status](https://github.com/finleap-connect/monoctl/actions/workflows/golang.yaml/badge.svg)](https://github.com/finleap-connect/monoctl/actions/workflows/golang.yaml)
[![Coverage Status](https://coveralls.io/repos/github/finleap-connect/monoctl/badge.svg?branch=main)](https://coveralls.io/github/finleap-connect/monoskope?branch=main)
[![Go Report Card](https://goreportcard.com/badge/github.com/finleap-connect/monoctl)](https://goreportcard.com/report/github.com/finleap-connect/monoskope)
[![Go Reference](https://pkg.go.dev/badge/github.com/finleap-connect/monoctl.svg)](https://pkg.go.dev/github.com/finleap-connect/monoctl)
[![GitHub release](https://img.shields.io/github/release/finleap-connect/monoctl.svg)](https://github.com/finleap-connect/monoctl/releases)

![monoctl logo](logo/monoctl.png)

`monoctl` implements the cli for Monoskope.

## Documentation

### Getting started

* Build executable with: `make go-build-monoctl`. This will build an executable with your architecture. Rename the one for your system to `monoctl`
* Set up configuration: `./monoctl config init -u api.monoskope.your.domain:443`

The creation of the configuration will also trigger an authentication with the Gitlab OpenID provider used by Monoskope. If you have a valid session in your browser this will be re-used.

`monoctl` also uses the keyring integration with the [`zalando/go-keyring`](https://github.com/zalando/go-keyring) library. When starting `monoctl` this may result in a dialog box appearing, that requests your password.

### General

* Docs on the almighty [Makefile](docs/Makefile.md)
