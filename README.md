# OpcuaMqttBridge

This project bridges data from multiple OPC UA servers to MQTT brokers.

## Features
- Stores OPC UA servers, MQTT brokers, and point mappings in a SQLite database.
- Connects to multiple servers and brokers concurrently.
- Subscribes to OPC UA nodes and publishes updates to configured MQTT topics.
- Designed for future expansion with web-based configuration.

## Building
To build the bridge, ensure Go is installed and run:

```bash
go mod tidy
go build ./...
```

The build requires network access to download dependencies.

## Running
After populating the database with server, broker, and point information, run:

```bash
go run ./cmd/bridge
```

The service will read configuration from `bridge.db` (or `DB_PATH` environment variable) and start bridging.

