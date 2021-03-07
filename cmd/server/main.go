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
	var upg *tableflip.Upgrader
	var l net.Listener
	var err error

	// include pid on the logs
	log.SetPrefix(fmt.Sprintf("pid:%d ", os.Getpid()))

	if config.GracefulTimeout > 0 {
		// Graceful restart mode
		upg, err = tableflip.New(tableflip.Options{
			PIDFile: config.PIDFile,
		})
		if err != nil {
			log.Fatal(err)
		}
		defer upg.Stop()
	}

	// prepare server
	s := &http.Server{
		Addr:         fmt.Sprintf(":%v", config.Port),
		Handler:      Router(!config.Development),
		ReadTimeout:  config.ReadTimeout,
		WriteTimeout: config.WriteTimeout,
		IdleTimeout:  config.IdleTimeout,
	}

	// listen service port
	if upg != nil {
		l, err = upg.Listen("tcp", s.Addr)
	} else {
		l, err = net.Listen("tcp", s.Addr)
	}

	if err != nil {
		log.Fatal(err)
	}

	log.Print(fmt.Sprintf("Listening %s", s.Addr))

	if upg != nil {
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

		// start servicing
		go func() {
			err := s.Serve(l)
			if err != nil && err != http.ErrServerClosed {
				log.Fatal(err)
			}
		}()

		// notify being ready for service
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
		err = s.Serve(l)
		if err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}
}
