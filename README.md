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

### 1. Install

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

### 2. First Use

Configure a keyfile. This can be a local filepath or accessible URL.

```bash
convey configure --keyfile FILE
```

```bash
echo "Hello world" | convey
<ID>
```

```bash
convey <ID>
Hello world
```

If you're looking for further ideas on what you can use this application for, see these [examples](docs/examples.md).


# Configuration

Set configuration with the `convey configure` command.

```
Usage:
  convey configure [flags]

Flags:
      --keyfile string        URL or local path to keyfile (at least 64 bytes is required)
      --short-names           Use short channel names (channel conflicts could be more likely for a given keyfile/fingerprint)
      --overwrite             Overwrite current configuration
      --fingerprint string    (advanced) If you know the fingerprint you want to use (SHAKE-256 hex), you can set it directly instead of using --keyfile
      --nats-cluster string   (advanced) NATS cluster id
      --nats-url string       (advanced) NATS server url
      --nats-cacert string    (advanced) Local path to CA certificate used by NATS server
  -h, --help                  help for configure
```

By default, configuration is loaded from `$HOME/.convey.yaml`.

[Further configuration docs](docs/configuration.md)

# Development

**Set up**
```bash
go get -u github.com/derekbekoe/convey
cd $GOPATH/src/github.com/derekbekoe/convey
go run main.go
```

[Further development docs](docs/development.md)

# Self-hosting

For convenience, we've provided a service that the application uses by default.

Alternatively, you can host your own [NATS Streaming Server](https://docs.nats.io/nats-streaming-concepts/intro) and configure `convey` to use that server.

[Further self-hosting docs](docs/self-hosting.md)

# License

Convey source code is available under the [MIT License](LICENSE).
