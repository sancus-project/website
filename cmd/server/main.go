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

var config = NewConfig()

func init() {
	getopt.FlagLong(&config.Development, "dev", 'd', "Don't hashify static files")
	getopt.FlagLong(&config.Port, "port", 'p', "HTTP port to listen")
	getopt.FlagLong(&config.PIDFile, "pid", 'f', "Path to PID file")
	getopt.FlagLong(&config.GracefulTimeout, "graceful", 't', "Maximum duration to wait for in-flight requests")

	getopt.Parse()

	// TODO: validate config
}

func main() {

	// setup
	log.SetPrefix(fmt.Sprintf("pid:%d ", os.Getpid()))

	listenAddr := fmt.Sprintf(":%v", config.Port)

	s := &http.Server{
		Addr:         listenAddr,
		Handler:      Router(!config.Development),
		ReadTimeout:  config.ReadTimeout,
		WriteTimeout: config.WriteTimeout,
		IdleTimeout:  config.IdleTimeout,
	}

	if config.GracefulTimeout > 0 {
		// Graceful restart mode
		upg, err := tableflip.New(tableflip.Options{
			PIDFile: config.PIDFile,
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
		time.AfterFunc(config.GracefulTimeout, func() {
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
