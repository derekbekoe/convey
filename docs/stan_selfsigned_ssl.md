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