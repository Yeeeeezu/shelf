# shelf

static file server. one binary, request logs, nothing else.

```
shelf ./dist -p 3000
```
```
  serving ./dist on http://localhost:3000

  15:04:01  200  1ms     /
  15:04:01  200  0ms     /assets/index.js
  15:04:01  304  0ms     /favicon.ico
```

## usage

```
shelf [dir] [flags]

  -p <port>   port (default: 8080)
  -q          quiet mode, no logs
```

## examples

```
shelf                        # serve current directory
shelf ./dist -p 3000         # serve build output
shelf /var/www -p 80 -q      # production-ish, quiet
```

## build

```
go build -o shelf .
```

Go 1.22+. no dependencies.
