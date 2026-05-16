# FoxESS Prometheus Exporter

A simple Prometheus exporter for FoxESS cloud data.

## Configuration

Set the `FOXESS_API_KEY` environment variable with your FoxESS API key.

```bash
export FOXESS_API_KEY=your-api-key-here
```

## Running

```bash
./foxess-exporter
```

The metrics will be available at `http://localhost:8080/metrics`.

## Docker

You can build and run the Docker image using `ko`:

```bash
# Build and run
ko build .
```
