# Storm Exporter

## Overview

`storm-exporter` is a Prometheus exporter for Apache Storm. It collects metrics
from the Storm UI and exposes them at an HTTP endpoint (`/metrics`) for
Prometheus to scrape.

## Features

- Reads metrics from Apache Storm UI
- Configurable via environment variables
- Comes with a `Dockerfile` to containerize the application
- Includes a `docker-compose.yml` file for testing with preconfigured Prometheus
  and Grafana

## Installation

### Using Docker

```sh
docker build -t storm-exporter .
docker run -p 8080:8080 -e STORM_UI_HOST=localhost:8081 storm-exporter
```

### Using Docker Compose

To run `storm-exporter` along with Prometheus and Grafana, use the provided
`docker-compose.yml`:

```sh
docker-compose up
```

## Configuration

`storm-exporter` supports configuration through environment variables:

| Environment Variable | Default Value    | Description                                                        |
| -------------------- | ---------------- | ------------------------------------------------------------------ |
| `STORM_UI_HOST`      | `localhost:8081` | The host of the Storm UI                                           |
| `EXPORTER_PORT`      | `8080`           | The port on which the metrics endpoint is exposed                  |
| `REFRESH_RATE`       | `5`              | The interval (in seconds) at which the exporter updates it metrics |

## Usage

Once running, Prometheus can scrape metrics from:

```
<http://localhost:8080/metrics>

```

## Monitoring with Prometheus & Grafana

A preconfigured `docker-compose.yml` file is included to spin up
`storm-exporter`, Prometheus, and Grafana. After starting the services, access:

- Prometheus: [http://localhost:9090](http://localhost:9090)
- Grafana: [http://localhost:3000](http://localhost:3000) (Default login:
  `admin/admin`)

## Acknowledgement

This project is largely inspired by the repo
[storm_expoter_prometheus](https://github.com/mr4x2/storm_exporter_prometheus).
Special thanks to the original author for their contributions.

## License

This project is licensed under the MIT License.
