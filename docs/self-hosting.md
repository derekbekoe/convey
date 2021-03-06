# Self-hosting

For convenience, we've provided a service that the application uses by default.

Alternatively, you can host your own [NATS Streaming Server](https://docs.nats.io/nats-streaming-concepts/intro) and configure `convey` to use that server.

Some of these methods do not use TLS. Only use the "no TLS" methods for experimental or development purposes only. Typically, you should enable TLS on the NATS Server. See [encrypt connections with TLS](https://docs.nats.io/developing-with-nats/security/tls).

## Host Local Docker container (no TLS)

Start NATS Streaming service as a Docker container:

```sh
docker run -p 4222:4222 nats-streaming:linux
```

Configure `convey` to use this server:

```
convey configure --nats-url nats://localhost:4222 --nats-cluster test-cluster
```

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

## Host Local Docker container with self-signed cert (TLS)

*The below steps to generate the certificates were adapted from this blog post: https://itnext.io/secure-pub-sub-with-nats-fcda983d0612*

In order to secure our NATS server, we will create a self-signed cert and we will sign this certificate with our own self-signed Certification Authority. If you are familiar with this and can create your own certificates, you can skip this part.

Download `cfssl` and `cfssljson` from https://github.com/cloudflare/cfssl/releases.

Prepare some files for `cfssl`:

```sh
echo '
{
  "CN": "Convey self-signed CA",
  "key": {
    "algo": "rsa",
    "size": 2048
  },
  "names": [
    {
      "O": "Convey CA",
      "L": "Portland",
      "ST": "Oregon",
      "C": "US"
    }
  ]
}' > ca.json
```

```sh
echo '
{
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
```

```sh
echo '{
    "CN": "Server",
    "hosts": [
        "127.0.0.1",
        "localhost"
    ]
}' > server.json
```

Generate the certificates:

```sh
cfssl gencert -initca ca.json | cfssljson -bare ca
cfssl gencert -ca=ca.pem -ca-key=ca-key.pem -config=config.json -profile=server server.json | cfssljson -bare server
```

Start the NATS Streaming server:

```sh
docker run -p 4222:4222 -v $(pwd):/certs nats-streaming:linux --cluster_id test-cluster --store MEMORY --max_channels 0 --max_subs 0 --max_msgs 0 --max_bytes 0 --max_age 24h --max_inactivity 10m -tls_client_cacert /certs/ca.pem --encrypt --encryption_key mykey --tlscert /certs/server.pem --tlskey /certs/server-key.pem --tls
```

Where `$(pwd)` is the directory that contains all the created certificates.

Finally, configure `convey` to use this server:

```
convey configure --nats-url nats://localhost:4222 --nats-cluster test-cluster --nats-cacert $(pwd)/ca.pem --keyfile FILE
```

If you want to host on a VM instead, it should be fairly straightforward to modify the above and configure `convey` to point to the correct IP or hostname.

## Host on VM with certificate signed by CA (TLS)

Create a resource group and VM (an Azure VM in this sample):
```sh
az group create -n nats -l westus
az vm create --image UbuntuLTS -g nats -n convey-nats-usw2-1 -l westus2 --size Standard_DS2_v2 --public-ip-address-dns-name convey-nats-usw2-1
az vm open-port -g nats -n convey-nats-usw2-1 --port 80 443 4443 4444
```

SSH into the VM:
```
ssh IP_ADDRESS
```

Use certbot to get your SSL certificate:

https://certbot.eff.org/lets-encrypt/ubuntubionic-other

Install and start NATS Server:

```sh
wget -O nats-server.deb https://github.com/nats-io/nats-server/releases/download/v2.1.2/nats-server-v2.1.2-amd64.deb

nohup nats-server --addr 0.0.0.0 --port 4443 --https_port 4444 --tlscert /etc/letsencrypt/live/convey-nats-usw2-1.westus2.cloudapp.azure.com/fullchain.pem --tlskey /etc/letsencrypt/live/convey-nats-usw2-1.westus2.cloudapp.azure.com/privkey.pem --tls --log /var/log/nats-server &
```

Other releases: https://github.com/nats-io/nats-server/releases

Install and start NATS Streaming Server:

```
wget -O nats-streaming-server.deb https://github.com/nats-io/nats-streaming-server/releases/download/v0.16.2/nats-streaming-server-v0.16.2-amd64.deb
dpkg -i nats-streaming-server.deb 

nohup nats-streaming-server --cluster_id test-cluster --store MEMORY --max_channels 0 --max_subs 0 --max_msgs 0 --max_bytes 0 --max_age 24h --max_inactivity 10m --encrypt --encryption_key mykey --nats_server nats://convey-nats-usw2-1.westus2.cloudapp.azure.com:4443 --log /var/log/nats-streaming-server &
```

Other releases: https://github.com/nats-io/nats-streaming-server/releases

Finally, configure `convey` to use this server:

```
convey configure --nats-url nats://convey-nats-usw2-1.westus2.cloudapp.azure.com:4443 --nats-cluster test-cluster --keyfile FILE
```
