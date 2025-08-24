# CE Node Exporter

A Prometheus exporter for Compiler Explorer infrastructure that provides system metrics about mount points.

## Metrics

The exporter provides the following metrics:

- `ce_node_mounts_total{type="..."}` - Total number of mount points by filesystem type
  - `type="all"` - Total count of all mounts
  - `type="cefs1"` - Count of squashfs mounts at `/efs/compiler-explorer/`
  - `type="cefs2"` - Count of mounts at `/cefs/XX/*` (where XX is a 2-character hash prefix)
  - Other filesystem types (ext4, tmpfs, squashfs, etc.)

## Building

```bash
./build.sh
```

## Running

```bash
# Default port 9100
./ce-node-exporter

# Custom port
./ce-node-exporter --listen-address=:9200
```

## Systemd Service

Install the service:
```bash
sudo cp ce-node-exporter.service /etc/systemd/system/
sudo systemctl daemon-reload
sudo systemctl enable ce-node-exporter
sudo systemctl start ce-node-exporter
```

## Configuration

The service can be configured by editing `/etc/systemd/system/ce-node-exporter.service` and modifying the `--listen-address` parameter in the `ExecStart` line.

## Endpoints

- `/metrics` - Prometheus metrics endpoint

## Requirements

- Go 1.22 or later (builds with Go 1.24.2 from `/opt/compiler-explorer/golang-1.24.2/`)
- Linux system with `/proc/mounts` available