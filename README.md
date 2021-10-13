# `monoctl`

[![pipeline status](https://gitlab.figo.systems/platform/monoskope/monoctl/badges/main/pipeline.svg)](https://gitlab.figo.systems/platform/monoskope/monoctl/-/commits/main)
[![coverage report](https://gitlab.figo.systems/platform/monoskope/monoctl/badges/main/coverage.svg)](https://gitlab.figo.systems/platform/monoskope/monoctl/-/commits/main)

![monoctl logo](logo/monoctl.png)

`monoctl` implements the cli for Monoskope.

## Documentation

### Getting started

* Build executable with: `make go-build-monoctl`. This will build an executable with your architecture. Rename the one for your system to `monoctl`
* Set up configuration: `./monoctl config init -u api.monoskope.your.domain:443`

The creation of the configuration will also trigger an authentication with the Gitlab OpenID provider used by Monoskope. If you have a valid session in your browser this will be re-used.

`monoctl` also uses the keyring integration with the [`zalando/go-keyring`](https://github.com/zalando/go-keyring) library. When starting `monoctl` this may result in a dialog box appearing, that requests your password.

### General

* Architecture and more in [GDrive](https://drive.google.com/drive/folders/1QEewDHF0LwSLr6aUVoHvMWrFgaJfJLty)
* Docs on the almighty [Makefile](docs/Makefile.md)

## Building

### Private Monoskope Repository

The build process needs to access private Git repositories that
contain modules used by `monoctl`. To enable the go toolchain to
access the code, the `GOPRIVATE` environment variable needs to be
set. This happens automatically in the `go-mod` target of the
Makefile but you need to set the variable before running the go
module commands separately during development.

```shell
export GOPRIVATE="gitlab.figo.systems/platform"
```

Alternatively you can set this permanently using the following command:

```shell
go env -w GOPRIVATE="gitlab.figo.systems/*" 
```
