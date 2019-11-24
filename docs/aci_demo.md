## Deploy to Azure Container Instances

We only include this as an illustration to keep the command simple as traffic to the container is not encrypted:

```
az group create -n nats -l westus
az container create --image nats-streaming:linux --location westus --command-line "/nats-streaming-server -cid convey-demo-cluster -mc 0 -ma 30m -mi 10m -D" -g nats -n nats-container --ports 4222 --ip-address Public
```

Connect from the cient:
```
convey configure --nats-url nats://IP_ADDRESS:4222 --nats-cluster convey-demo-cluster
convey --unsecure
```
