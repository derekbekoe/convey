# NATS Streaming Server

You can host your own [NATS Streaming Server](https://docs.nats.io/nats-streaming-concepts/intro) and configure `convey` to use that server.

#### Deploy to a local Docker container

```bash
docker run -p 4222:4222 nats-streaming:linux
convey configure --nats-url nats://localhost:4222 --nats-cluster test-cluster
```

You will need to use the `--unsecure` flag as TLS will not be enabled through this local container.


## Deploy to Azure Container Instances

We only include this as an illustration to keep the command simple as traffic to the container is not secure:

```
az group create -n nats -l westus
az container create --image nats-streaming:linux --location westus --command-line "/nats-streaming-server -cid convey-demo-cluster -mc 0 -ma 30m -mi 10m -D" -g nats -n nats-container --ports 4222 --ip-address Public
```

Connect from the cient:
```
convey configure --nats-url nats://IP_ADDRESS:4222 --nats-cluster convey-demo-cluster
convey --unsecure
```


Extracted from https://itnext.io/secure-pub-sub-with-nats-fcda983d0612

Download cfssl and cfssljson from https://github.com/cloudflare/cfssl/releases

echo '{
    "CN": "Techwhale CA",
    "key": {
        "algo": "rsa",
        "size": 2048
    },
    "names": [
        {
               "C": "FR",
               "L": "Antibes",
               "O": "Techwhale",
               "ST": "PACA"
        }
    ]
}' > ca.json

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

 TODO-DEREK Support custom ca.pem to support connecting to self-signed tls connection
 TODO-DEREK Figure out if clients can do things like list all channels etc.
