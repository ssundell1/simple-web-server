# Simple HTTP Web Server

Built for hosting files and logging web requests (both GET and POST forms) when penetration testing web apps.

## Endpoints
```txt
/               For basic logging or requests
/files          For listing files in the hosted directory
/files/{FILE}   For getting desired file
```

## Build
```bash
go build
```

## Run
```txt
Usage of ./simple-web-server:
  -d string
        the directory of static file to host (default "files")
  -l string
        log level: debug, info, warning, error (default "DEBUG")
  -p string
        port to serve on (default "8080")
```

```bash
./simple-web-server
```
Or
```bash
go run main.go
```