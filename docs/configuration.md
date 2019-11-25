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
