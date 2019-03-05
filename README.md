# convey

A command-line tool that makes sharing pipes between machines easy.

## Install

**Install on Linux:**
```bash
wget -qO convey https://get.convey.sh/linux
chmod +x ~/bin/convey
```

**Install on Mac OS:**
```bash
wget -qO ~/bin/convey https://get.convey.sh/macos
chmod +x ~/bin/convey
```

**Install on Windows:**  
> Download from https://get.convey.sh/windows

This will download the latest release for your platform.  
Builds are available for the `amd64` architecture.  
To directly from GitHub, visit https://github.com/derekbekoe/convey/releases/latest.  

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

Set configuration with the `convey configure` command.
```bash
convey configure --nats-url nats://localhost:4222 --nats-cluster test-cluster
```

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

## Host your own NATS Streaming Server

**Deploy to a local Docker container**

```bash
docker run -p 4222:4222 nats-streaming:linux
convey configure --nats-url nats://localhost:4222 --nats-cluster test-cluster
```

**Deploy to an Azure Container Instance**

note: We only include this as an illustration to keep the command simple as traffic is not encrypted.
```bash
az container create --image nats-streaming:linux --ports 4222 --ip-address Public -g RG -n nats1
convey configure --nats-url nats://<IPADDRESS>:4222 --nats-cluster test-cluster
```

**TLS**

openssl req -x509 -nodes -days 3650 -newkey rsa:2048 -keyout certs/key.pem -out certs/cert.pem -subj "/C=US/ST=Texas/L=Austin/O=AwesomeThings/CN=localhost"

```
cd /Users/derekb/go/src/github.com/derekbekoe/convey
docker run -p 4222:4222 -v $(pwd)/certs:/certs nats-streaming:linux -tls_client_cert /certs/cert.pem -tls_client_key /certs/key.pem  --tlscert /certs/cert.pem --tlskey /certs/key.pem --tlsverify=false --tls=false -secure=false
```

## Platform Builds
```bash
go get github.com/mitchellh/gox
gox -ldflags "-X github.com/derekbekoe/convey/cmd.VersionGitCommit=$(git rev-list -1 HEAD) -X github.com/derekbekoe/convey/cmd.VersionGitTag=VERSION" -os="linux darwin" -arch="amd64" -output="bin/{{.Dir}}_{{.OS}}_{{.Arch}}"
```
See https://golang.org/doc/install/source#environment

## FAQ

**How do I try it out?**

Start the local container, download convey, specify the configuration then run.

If you'd like to share between multiple devices, host the server in a location where your devices can access.
