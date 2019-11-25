**Multi-platform Builds**
```bash
go get github.com/mitchellh/gox
gox -ldflags "-X github.com/derekbekoe/convey/cmd.VersionGitCommit=$(git rev-list -1 HEAD) -X github.com/derekbekoe/convey/cmd.VersionGitTag=VERSION" -os="linux darwin" -arch="amd64" -output="bin/{{.Dir}}_{{.OS}}_{{.Arch}}"
```
See https://golang.org/doc/install/source#environment

**Go Module Verification**
```bash
go mod tidy
# verification
go build
go test
# list all
go list -m all
```