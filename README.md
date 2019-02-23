# convey

A command-line tool that makes sharing pipes between machines easy.

## Download

get.convey.sh/linux

get.convey.sh/windows

This will download the latest release for your platform.

## Usage

In Terminal 1:
```bash
echo "Hello world" | convey
21f50fba373e11e9990a72000872a940
```

In Terminal 2:
```bash
convey 21f50fba373e11e9990a72000872a940
Hello world
```

## Configuration

By default, configuration is loaded from `$HOME/.convey.yaml`.

This is an example of `.convey.yaml`:
```yaml
NatsURL: nats://localhost:4223
NatsClusterID: test-cluster
```

Use the --config flag on the command line to change the config file if needed.

## Development
```bash
go get -u github.com/derekbekoe/convey
cd $GOPATH/src/github.com/derekbekoe/convey
go run main.go
go build -o bin/convey
```

```bash
go build -ldflags "-X github.com/derekbekoe/convey/cmd.VersionGitCommit=$(git rev-list -1 HEAD) -X github.com/derekbekoe/convey/cmd.VersionGitTag=VERSION" -o bin/convey
```

Cross-compile
```bash
go get github.com/mitchellh/gox
gox -ldflags "-X github.com/derekbekoe/convey/cmd.VersionGitCommit=$(git rev-list -1 HEAD) -X github.com/derekbekoe/convey/cmd.VersionGitTag=0.0.1" -os="linux darwin" -arch="amd64" -output="bin/{{.Dir}}_{{.OS}}_{{.Arch}}"
```
See https://golang.org/doc/install/source#environment

## Starting NATS Streaming Server

```bash
docker run -p 4223:4223 -p 8223:8223 nats-streaming -p 4223 -m 8223
```

## FAQ

**How do I try it out?**

Start the local container, download convey, specify the configuration then run.

If you'd like to share between multiple devices, host the server in a location where your devices can access.
