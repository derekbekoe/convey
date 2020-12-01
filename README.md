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
<p>A command-line tool that makes it easy to pipe between machines.</p>
<p>Learn more at <a href="https://blog.derekbekoe.com/convey"><em>Convey: Pipe between machines</em></a></p>
</div>

```bash
echo "Hello world" | convey
vibrant_allen
```
```bash
convey vibrant_allen
Hello world
```

# Features

- Pipe between hosts with an idiomatic interface using the standard `|` symbol.
- Easily pipe files between hosts.
- Does not require any open ports between your clients.
- Short channel names allow for easy typing such as `vibrant_allen`.
- Supports colors through [ANSI escape codes](https://en.wikipedia.org/wiki/ANSI_escape_code#Colors).
- Supports Linux, macOS and Windows.
- No dependencies to install.
- Powered by [NATS](https://nats.io/), a CNCF project.
- Pre-configured to use our hosted service so you can get started right away.
  - Data in encrypted in transit with TLS and encrypted in the in-memory store.
  - Data is deleted after 10 minutes of inactivity or 24 hours.
- Self-hosting is available if you'd prefer.

# Getting Started

### 1. Install

#### Linux
```bash
mkdir -p ~/bin
wget -qO ~/bin/convey https://get.convey.sh/linux
chmod +x ~/bin/convey
~/bin/convey -h
```

#### macOS
```bash
mkdir -p ~/bin
curl -sLo ~/bin/convey https://get.convey.sh/macos
chmod +x ~/bin/convey
~/bin/convey -h
```

#### Windows  
```powershell
Invoke-WebRequest https://get.convey.sh/windows -OutFile convey.exe
.\convey.exe -h
```

### 2. Configure Keyfile

Configure a keyfile. This can be a local filepath, accessible URL or file download link.

```bash
convey configure --keyfile FILE
```

The keyfile should be a secret file that can be easily accessed on the machines you want to use `convey` with.  
Your keyfiles don't leave your machine. We create a fingerprint from this file and use that fingerprint only.  
Some examples are:
- a text file (e.g. `~/.ssh/id_rsa.pub`)
- an image file
- a file with randomly generated bytes - [how to](https://unix.stackexchange.com/questions/33629/how-can-i-populate-a-file-with-random-data)
- raw URL to gist - [GitHub gist](https://gist.github.com)

### 3. First Use

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
      --long-names            Use standard uuid format for channel names 
      --overwrite             Overwrite current configuration
      --fingerprint string    (advanced) If you know the fingerprint you want to use (SHAKE-256 hex), you can set it directly instead of using --keyfile
      --nats-cacert string    (advanced) Local path to CA certificate used by NATS server
      --nats-cluster string   (advanced) NATS cluster id
      --nats-url string       (advanced) NATS server url
  -h, --help                  help for configure
```

By default, configuration is loaded from `$HOME/.convey.yaml`.

[Further configuration docs](docs/configuration.md)

# Development

```bash
go get -u github.com/derekbekoe/convey
cd $GOPATH/src/github.com/derekbekoe/convey
go run main.go
```

[Further development docs](docs/development.md)

# Self-hosting

For convenience, we've provided a hosted service that `convey` uses by default.  
This hosted service uses TLS to ensure communications are encrypted.

Alternatively, you can host your own [NATS Streaming Server](https://docs.nats.io/nats-streaming-concepts/intro) and configure `convey` to use that server.

[Further self-hosting docs](docs/self-hosting.md)

# License

Convey source code is available under the [MIT License](LICENSE).
