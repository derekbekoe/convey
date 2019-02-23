# convey

A command-line tool that makes sharing pipes between machines easy.

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

Cross-compile
```bash
env GOOS=linux GOARCH=amd64 go build
```
See https://golang.org/doc/install/source#environment

## Startings NATS Streaming Server

```bash
docker run -p 4223:4223 -p 8223:8223 nats-streaming -p 4223 -m 8223
```
