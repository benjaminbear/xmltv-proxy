# XMLTV Proxy Server

![GitHub Workflow Status](https://img.shields.io/github/workflow/status/benjaminbear/xmltv-proxy/goreleaser)

![GitHub release (latest by date)](https://img.shields.io/github/v/release/benjaminbear/xmltv-proxy)
![Go version](https://img.shields.io/github/go-mod/go-version/benjaminbear/xmltv-proxy?filename=go.mod)
![License](https://img.shields.io/github/license/benjaminbear/xmltv-proxy)

This docker service grabs **epgdata** from a source (tv spielfilm) and deploys **xmltv** format via http server.

## Usage

### Docker image

Pull repository

```bash
docker pull bbaerthlein/xmltv-proxy
```


Run container:

```bash
docker run -p 8080:8080 bbaerthlein/xmltv-proxy
```

### Binary

Download and extract the binary for your os and architecture from the [release page](https://github.com/benjaminbear/xmltv-proxy/releases/).

Start the binary: (example for linux)

```bash
XMLTV_PORT=1234 ./xmltv-proxy_linux_x64
```

## Environment variables

```
XMLTV_DAYS: (optional, default=7) Count of days to download from your epg source and serve: 1-14
XMLTV_TIMEZONE: (optional, default=system) Set your timezone for correct airtime: e.g. "Europe/Berlin"
XMLTV_DAILY_DOWNLOAD: (optional, default=false) Force downloading latest epg every day
XMLTV_PORT: (optional, default=8080) Use specific port for webservice
XMLTV_INSECURE_TLS: (optional, default=false) Allow insecure http/s requests
```

## Mount directories (docker)

If you want persist downloaded epgdata files add:

```
-v /your/local/path:/epg/epgdata_files
```

## Docker hub

https://hub.docker.com/r/bbaerthlein/xmltv-proxy