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

Build and run with Docker:

```bash
docker build -t foxess-exporter .
docker run -e FOXESS_API_KEY=your-api-key-here -p 8080:8080 foxess-exporter
```
