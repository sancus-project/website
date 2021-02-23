package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	reuseport "github.com/kavu/go_reuseport"
	"github.com/pborman/getopt/v2"
	"github.com/rs/seamless"

	"github.com/amery/go-webpack-starter/assets"
)

var (
	waitStarted     = time.Second // wait one second before considering it started
	defaultPidFile  = "/tmp/reuseport.pid"
	defaultPort     = 8080
	defaultGraceful = 60 * time.Second

	devFlag         = getopt.BoolLong("dev", 'd', "Don't hashify static files")
	portListen      = getopt.Uint16Long("port", 'p', uint16(defaultPort), "HTTP port to listen")
	pidFile         = getopt.StringLong("pid", 'f', defaultPidFile, "Seemless restart PID file")
	gracefulTimeout = getopt.DurationLong("graceful", 't', defaultGraceful, "Maximum duration to wait for in-flight requests")
)

func Handler(w http.ResponseWriter, r *http.Request) {
	if d := r.URL.Query().Get("delay"); d != "" {
		if delay, err := time.ParseDuration(d); err == nil {
			time.Sleep(delay)
		}
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Server pid: %d\n", os.Getpid())
}

func main() {

	getopt.Parse()

	listenAddr := fmt.Sprintf(":%v", *portListen)

	s := &http.Server{
		Addr:    listenAddr,
		Handler: http.HandlerFunc(Handler),
	}

	if !*devFlag {
		// service hashified statics on non-dev mode
		s.Handler = assets.Files.Handler(true, s.Handler)
	}

	if !*devFlag && *gracefulTimeout > 0 {
		// Graceful restart mode
		seamless.Init(*pidFile)

		l, err := reuseport.Listen("tcp", listenAddr)
		if err != nil {
			log.Fatal(err)
		}

		seamless.OnShutdown(func() {
			ctx, cancel := context.WithTimeout(context.Background(), *gracefulTimeout)
			defer cancel()

			if err := s.Shutdown(ctx); err != nil {
				log.Print("Graceful shutdown timeout, force closing")
				s.Close()
			}
		})

		go func() {
			time.Sleep(waitStarted)
			seamless.Started()
		}()

		err = s.Serve(l)
		if err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}

		seamless.Wait()
	} else {
		l, err := net.Listen("tcp", s.Addr)
		if err != nil {
			log.Fatal(err)
		}

		log.Print(fmt.Sprintf("Listening %s", s.Addr))
		err = s.Serve(l)
		if err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}
}
