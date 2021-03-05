package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/cloudflare/tableflip"
	"github.com/pborman/getopt/v2"
)

var (
	waitStarted     = time.Second // wait one second before considering it started
	defaultPidFile  = "/tmp/tableflip.pid"
	defaultPort     = 8080
	defaultGraceful = 60 * time.Second

	devFlag         = getopt.BoolLong("dev", 'd', "Don't hashify static files")
	portListen      = getopt.Uint16Long("port", 'p', uint16(defaultPort), "HTTP port to listen")
	pidFile         = getopt.StringLong("pid", 'f', defaultPidFile, "Path to PID file")
	gracefulTimeout = getopt.DurationLong("graceful", 't', defaultGraceful, "Maximum duration to wait for in-flight requests")
)

func main() {

	// check argments
	getopt.Parse()

	// setup
	log.SetPrefix(fmt.Sprintf("pid:%d ", os.Getpid()))

	listenAddr := fmt.Sprintf(":%v", *portListen)

	s := &http.Server{
		Addr:    listenAddr,
		Handler: Router(!*devFlag),
	}

	if *gracefulTimeout > 0 {
		// Graceful restart mode
		upg, err := tableflip.New(tableflip.Options{
			PIDFile: *pidFile,
		})
		if err != nil {
			log.Fatal(err)
		}
		defer upg.Stop()

		// attempt upgrade on SIGUSR2
		go func() {
			sig := make(chan os.Signal, 1)
			signal.Notify(sig, syscall.SIGUSR2)
			for range sig {
				if err := upg.Upgrade(); err != nil {
					log.Println("Upgrade failed:", err)
				}
			}
		}()

		// listen service port
		l, err := upg.Listen("tcp", listenAddr)
		if err != nil {
			log.Fatal(err)
		}

		log.Print(fmt.Sprintf("Listening %s", s.Addr))

		// starter servicing
		go func() {
			err := s.Serve(l)
			if err != nil && err != http.ErrServerClosed {
				log.Fatal(err)
			}
		}()

		if err = upg.Ready(); err != nil {
			log.Fatal(err)
		}
		<-upg.Exit()

		// graceful shutdown timeout
		time.AfterFunc(*gracefulTimeout, func() {
			log.Println("Graceful shutdown timed out")
			os.Exit(1)
		})

		// Wait for connections to drain.
		s.Shutdown(context.Background())
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
