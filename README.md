<p align="center">
  <img 
    src="https://derekb.blob.core.windows.net/public/convey_1.svg" 
    width="400" border="0" alt="Convey">
</p>
<p align="center">
<a href="https://github.com/derekbekoe/convey/releases"><img src="https://img.shields.io/github/release/derekbekoe/convey.svg?style=flat-square&logo=github&color=%236C63FF" alt="Version"></a>
<a href="https://travis-ci.org/derekbekoe/convey"><img src="https://img.shields.io/travis/derekbekoe/convey/master.svg?style=flat-square&logo=travis" alt="Build Status"></a>
</p>
<div align="center">
<p>A command-line tool that makes sharing pipes between machines easy.</p>
<p>Learn more at <a href="https://blog.derekbekoe.com/convey"><em>Convey: Pipe between machines</em></a></p>
</div>


```bash
echo "Hello world" | convey
21f50fba373e11e9990a72000872a940
```
```bash
convey 21f50fba373e11e9990a72000872a940
Hello world
```

# Features

- Pipe between hosts with an idomatic interface using the standard `|` symbol.
- Easily pipe files between hosts.
- Does not require any open ports between your clients.
- Supports colors through [ANSI escape codes](https://en.wikipedia.org/wiki/ANSI_escape_code#Colors).
- Supports Linux, macOS and Windows.
- No dependencies to install.
- Powered by [NATS](https://nats.io/), a CNCF project.

# Getting Started

### Linux
```bash
wget -qO convey https://get.convey.sh/linux
chmod +x ~/bin/convey
~/bin/convey -h
```

### macOS
```bash
wget -qO ~/bin/convey https://get.convey.sh/macos
chmod +x ~/bin/convey
~/bin/convey -h
```

### Windows  
```powershell
Invoke-WebRequest https://get.convey.sh/windows -OutFile convey.exe
convey -h
```

### Go
```bash
go get github.com/derekbekoe/convey
```

# Configuration

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

Use the `--config` flag on the command line to change the config file used if needed.

# NATS Streaming Server

You can host your own [NATS Streaming Server](https://nats.io/documentation/streaming/nats-streaming-intro/) and configure `convey` to use that server.

**Deploy to a local Docker container**

```bash
docker run -p 4222:4222 nats-streaming:linux
convey configure --nats-url nats://localhost:4222 --nats-cluster test-cluster
```

**Deploy to an Azure Container Instance**

note: We only include this as an illustration to keep the command simple as traffic is not encrypted.
```bash
az container create --image nats-streaming:linux --ports 4222 --ip-address Public -g <RG> -n nats1
convey configure --nats-url nats://<IPADDRESS>:4222 --nats-cluster test-cluster
```

**TLS**

```bash
openssl req -x509 -nodes -days 3650 -newkey rsa:2048 -keyout certs/key.pem -out certs/cert.pem -subj "/C=US/ST=Texas/L=Austin/O=AwesomeThings/CN=localhost"
```

```bash
cd /Users/derekb/go/src/github.com/derekbekoe/convey
docker run -p 4222:4222 -v $(pwd)/certs:/certs nats-streaming:linux -tls_client_cert /certs/cert.pem -tls_client_key /certs/key.pem  --tlscert /certs/cert.pem --tlskey /certs/key.pem --tlsverify=false --tls=false -secure=false
```

# Development

**Set up**
```bash
go get -u github.com/derekbekoe/convey
cd $GOPATH/src/github.com/derekbekoe/convey
go run main.go
```

**Multi-platform Builds**
```bash
go get github.com/mitchellh/gox
gox -ldflags "-X github.com/derekbekoe/convey/cmd.VersionGitCommit=$(git rev-list -1 HEAD) -X github.com/derekbekoe/convey/cmd.VersionGitTag=VERSION" -os="linux darwin" -arch="amd64" -output="bin/{{.Dir}}_{{.OS}}_{{.Arch}}"
```
See https://golang.org/doc/install/source#environment

# Examples

Click to expand each gif.

<div style="text-align: center">
<img src="https://derekb.blob.core.windows.net/public/blog_convey_vm_1_top.gif" alt="Convey with Top" border="0" width="445">
<img src="https://derekb.blob.core.windows.net/public/blog_convey_vm_1_filecp.gif" alt="Convey for file piping" border="0" width="445">
<img src="https://derekb.blob.core.windows.net/public/blog_convey_vm_1_ms.gif" alt="Convey with millisecond date" border="0" width="445">
<img src="https://derekb.blob.core.windows.net/public/blog_convey_cloudshell_1.gif" alt="Convey with Cloud Shell" border="0" width="445">
</div>

# License

Convey source code is available under the [MIT License](LICENSE).
