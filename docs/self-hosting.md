# Self-hosting

For convenience, we've provided a service that the application uses by default.

Alternatively, you can host your own [NATS Streaming Server](https://docs.nats.io/nats-streaming-concepts/intro) and configure `convey` to use that server.

Some of these methods require using the `convey` `--unsecure` flag. Only use this flag for experimental or development purposes only. We only include these methods as an illustration as traffic to the server is not encrypted. Typically, you should [encrypt connections with TLS](https://docs.nats.io/developing-with-nats/security/tls).

## Host Local Docker container (no TLS)

Start NATS Streaming service as a Docker container:

```sh
docker run -p 4222:4222 nats-streaming:linux
```

Configure `convey` to use this server:

```
convey configure --nats-url nats://localhost:4222 --nats-cluster test-cluster
```

Note: You will need to use the `--unsecure` flag as TLS will not be enabled through this local container.

## Host on Azure Container Instances (no TLS)

Create a resource group and create the container:
```sh
az group create -n nats -l westus
az container create --image nats-streaming:linux --location westus --command-line "/nats-streaming-server -cid test-cluster -mc 0 -ma 30m -mi 10m -D" -g nats -n nats-container --ports 4222 --ip-address Public
```

For the meaning of the `--command-line` arguments, see [NATS Streaming Server - Command Line Arguments](https://docs.nats.io/nats-streaming-server/configuring/cmdline).

Get the IP address of the container and use it below.

Configure `convey` to use this server:

```
convey configure --nats-url nats://IP_ADDRESS:4222 --nats-cluster test-cluster
```

Note: You will need to use the `--unsecure` flag as TLS will not be enabled through this local container.

## Host Local Docker container with self-signed cert (TLS)

Extracted from https://itnext.io/secure-pub-sub-with-nats-fcda983d0612

Download cfssl and cfssljson from https://github.com/cloudflare/cfssl/releases

```sh
echo '
{
  "CN": "Convey CA",
  "key": {
    "algo": "rsa",
    "size": 2048
  },
  "names": [
    {
      "O": "Convey",
      "L": "Portland",
      "ST": "Oregon",
      "C": "US"
    }
  ]
}' > ca.json
```

echo '{
    "signing": {
        "default": {
            "expiry": "43800h"
        },
        "profiles": {   
            "server": {
                "expiry": "43800h",
                "usages": [
                    "signing",
                    "digital signing",
                    "key encipherment",
                    "server auth"
                ]
            }
        }
    }
}' > config.json

echo '{
    "CN": "Server",
    "hosts": [
        "127.0.0.1",
        "messaging.techwhale.io"
    ]
}' > server.json

./cfssl gencert -initca ca.json | ./cfssljson -bare ca
./cfssl gencert -ca=ca.pem -ca-key=ca-key.pem -config=config.json -profile=server server.json | ./cfssljson -bare server

 docker run -p 4222:4222 -v $(pwd)/certs:/certs nats-streaming:linux -mc 0 -tls_client_cacert /certs/ca.pem --encrypt --encryption_key mykey --tlscert /certs/server.pem --tlskey /certs/server-key.pem --tls


Configure `convey` to use this server:

```
convey configure --nats-url nats://localhost:4222 --nats-cluster test-cluster
```

TODO-DEREK Verify it works without --unsecure

If you want to host on a VM instead, it should be fairly straightforward to modify the above and configure `convey` to point to the correct IP or hostname.

## Host on VM with certificate signed by CA (TLS)

TODO-DEREK Complete this using Lets Encrypt.