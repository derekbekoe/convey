<p align="center">
  <img 
    src="https://derekb.blob.core.windows.net/public/convey_1.svg" 
    width="400" border="0" alt="Convey">
</p>
<p align="center">
<a href="https://github.com/derekbekoe/convey/releases"><img src="https://img.shields.io/github/release/derekbekoe/convey.svg?style=flat-square&logo=github&color=%236C63FF" alt="Version"></a>
<a href="https://travis-ci.org/derekbekoe/convey"><img src="https://img.shields.io/travis/derekbekoe/convey/master.svg?style=flat-square&logo=travis" alt="Build Status"></a>
<a href="https://online.visualstudio.com/environments/new?name=👏%20Convey&repo=derekbekoe/convey"><img src="https://img.shields.io/static/v1?style=flat-square&logo=microsoft&label=VS%20Online&message=Create&color=blue" alt="VS Online"></a>
</p>
<div align="center">
<p>A command-line tool that makes it easy to pipe between machines.</p>
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
- Configure it to use short channel names instead of UUIDs for easy typing such as `vibrant_allen`.
- Supports colors through [ANSI escape codes](https://en.wikipedia.org/wiki/ANSI_escape_code#Colors).
- Supports Linux, macOS and Windows.
- No dependencies to install.
- Powered by [NATS](https://nats.io/), a CNCF project.

# Getting Started

## 1. Install

#### Linux
```bash
wget -qO ~/bin/convey https://get.convey.sh/linux
chmod +x ~/bin/convey
~/bin/convey -h
```

#### macOS
```bash
curl -sLo ~/bin/convey https://get.convey.sh/macos
chmod +x ~/bin/convey
~/bin/convey -h
```

#### Windows  
```powershell
Invoke-WebRequest https://get.convey.sh/windows -OutFile convey.exe
.\convey.exe -h
```

## 2. First Use

```bash
convey configure --keyfile URL
echo "Hello world" | convey
<ID>
```
```bash
convey <ID>
Hello world
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

#### Deploy to a local Docker container

```bash
docker run -p 4222:4222 nats-streaming:linux
convey configure --nats-url nats://localhost:4222 --nats-cluster test-cluster
```

You will need to use the `--unsecure` flag as TLS will not be enabled through this local container.

# Development

**Set up**
```bash
go get -u github.com/derekbekoe/convey
cd $GOPATH/src/github.com/derekbekoe/convey
go run main.go
```

[Further development docs](docs/development.md)

# License

Convey source code is available under the [MIT License](LICENSE).
