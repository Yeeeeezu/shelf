package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"time"
)

type logged struct {
	h     http.Handler
	quiet bool
}

type rw struct {
	http.ResponseWriter
	code int
}

func (r *rw) WriteHeader(c int) { r.code = c; r.ResponseWriter.WriteHeader(c) }

func (l *logged) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if l.quiet {
		l.h.ServeHTTP(w, r)
		return
	}
	lw := &rw{ResponseWriter: w, code: 200}
	t := time.Now()
	l.h.ServeHTTP(lw, r)
	color := "\x1b[32m"
	if lw.code >= 400 {
		color = "\x1b[31m"
	} else if lw.code >= 300 {
		color = "\x1b[33m"
	}
	fmt.Printf("  \x1b[2m%s\x1b[0m  %s%-3d\x1b[0m  \x1b[2m%-6v\x1b[0m  %s\n",
		time.Now().Format("15:04:05"), color, lw.code, time.Since(t).Round(time.Millisecond), r.URL.Path)
}

func main() {
	port := flag.String("p", "8080", "port")
	quiet := flag.Bool("q", false, "suppress logs")
	flag.Usage = func() {
		fmt.Print(`
  shelf — static file server

  usage:
    shelf [dir] [-p port] [-q]

  flags:
    -p <port>   port to listen on (default: 8080)
    -q          quiet, no request logs

  examples:
    shelf
    shelf ./dist -p 3000
    shelf /var/www -p 80

`)
	}
	flag.Parse()

	dir := "."
	if flag.NArg() > 0 {
		dir = flag.Arg(0)
	}
	if _, err := os.Stat(dir); err != nil {
		fmt.Fprintf(os.Stderr, "  cannot access %s: %v\n", dir, err)
		os.Exit(1)
	}

	http.Handle("/", &logged{h: http.FileServer(http.Dir(dir)), quiet: *quiet})
	fmt.Printf("  \x1b[2mserving\x1b[0m %s \x1b[2mon\x1b[0m http://localhost:%s\n\n", dir, *port)
	if err := http.ListenAndServe(":"+*port, nil); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
