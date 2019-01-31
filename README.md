# convey

A command-line tool that makes sharing pipes between machines easy.

## Usage

```bash
echo "Hello world" | convey
```

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
