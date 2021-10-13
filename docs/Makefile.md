# Makefile

When developing, the `Makefile` comes in handy to help you with various tasks.
There are specific `*.mk` files for things like go, etc. which provides targets for developing with those tools.

The following targets are defined. Please not that there are variables (uppercase) which can be overriden:

| target | Description |
| --------- | ----------- |
| *general* | |
| `clean` | Cleans everything, tools, tmp dir used, whatever |
| `tools` | Install necessary tools to `TOOLS_DIR`, like `ginkgo`, `golangci-lint`, ... |
| `tools-clean` | Removes the tools |
| *go* | |
| `go-mod` | Downloads all require go modules |
| `go-fmt` | Formats all `*.go` files |
| `go-vet` | Vets all go code |
| `go-lint` | Lints all go code |
| `go-run-*` | Runs the app in `cmd/*`, e.g. `go-run-monoctl` to run `monoctl` from sources |
| `go-test` | Runs all go tests |
