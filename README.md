# Diskcollector
## Introduction
A small Go application with corresponding library that I have written for a still nonexisting project. The applications main purpose is to support metric collection about PVCs when there is no vendor support from the storageclass. The idea is then to supply this application as a sidecar container to the running pod and expose metrics about the pvc.

## Installation
```sh
~ » git clone git@github.com:gelm0/diskcollector.git
~ » cd diskcollector/
diskcollector » go mod tidy
diskcollector » make
```

## Usage
```sh
diskstat
```
Will run the application on port :8080 collecting metrics about the mountpath "/" by default exposing metrics in a [prometheus exposition format](https://prometheus.io/docs/instrumenting/exposition_formats/#text-based-format).
```sh
~ » curl localhost:8080/metrics
# HELP disk_available_bytes Total available disk space left in bytes
# TYPE disk_available_bytes gauge
disk_available_bytes{mount="/"} 7.8872551424e+10
# HELP disk_free_bytes Total free space disk in bytes
# TYPE disk_free_bytes gauge
disk_free_bytes{mount="/"} 7.3487065088e+10
# HELP disk_size_bytes Total disk space in bytes
# TYPE disk_size_bytes gauge
disk_size_bytes{mount="/"} 1.05089261568e+11
# HELP disk_used_bytes Total usage of the disk in bytes
# TYPE disk_used_bytes gauge
disk_used_bytes{mount="/"} 2.6216710144e+10

```

Following **options** can be supplied to diskstat
- **-path** - Path of the mountpath which disk usage you want to monitor
- **-certFile** - Path to certificate for serving TLS. The certfile should include the whole certificate chain. This means certificate, intermediates and CA. If this option is supplied keyFile option must be also supplied.
- **-keyFile** - Path to private key for the certificate.


Following **Environment variables** can be supplied to diskstat
- **METRICS_PORT** - Application port. Defaults to 8080.
- **METRICS_ADDRESS** - Application address. Defaults to localhost.
