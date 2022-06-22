# Makefile

```
Usage:
  make <target>
  add-license      Adds the license to every file
  check-license    Checks thath the license is set on every file

General
  help             Display this help.
  mod              Do go mod tidy, download, verify
  vet              Do go ver
  go               Do go mod / vet / lint /test
  run              run monoctl, use `ARGS="get user"` to pass arguments
  test             run all tests
  test-ci          run all tests in CICD
  coverage         print coverage from coverprofiles
  gomock-get       download gomock
  lint             go lint
  tools            Target to install all required tools into TOOLS_DIR
  clean            Target clean up tools in TOOLS_DIR
  build-clean      clean up binaries
  build-monoctl-linux  build monoctl for linux
  build-monoctl-osx  build monoctl for osx
  build-monoctl-win  build monoctl for windows
  build-monoctl-all  build monoctl for linux, osx and windows
  rebuild-mocks    rebuild go mocks
```