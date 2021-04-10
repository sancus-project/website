package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/cloudflare/tableflip"
	"github.com/pborman/getopt/v2"

	"github.com/sancus-project/website/web"
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
	// include pid on the logs
	log.SetPrefix(fmt.Sprintf("pid:%d ", os.Getpid()))

	// Graceful restart mode
	upg, err := tableflip.New(tableflip.Options{
		PIDFile: config.PIDFile,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer upg.Stop()

	// prepare server
	r := web.Router{
		HashifyAssets: !config.Development,
	}
	if err := r.Compile(); err != nil {
		log.Fatal(err)
	}

	s := &http.Server{
		Addr:         fmt.Sprintf(":%v", config.Port),
		Handler:      r.Handler(),
		ReadTimeout:  config.ReadTimeout,
		WriteTimeout: config.WriteTimeout,
		IdleTimeout:  config.IdleTimeout,
	}

	// listen service port
	l, err := upg.Listen("tcp", s.Addr)
	if err != nil {
		log.Fatal(err)
	}

	log.Print(fmt.Sprintf("Listening %s", l.Addr()))

	// watch signals
	go func() {
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, syscall.SIGHUP, syscall.SIGUSR2, syscall.SIGINT, syscall.SIGTERM)

		for signum := range sig {

			switch signum {
			case syscall.SIGHUP:
				// attempt to reload config on SIGHUP
				if err := r.Reload(); err != nil {
					log.Println("Reload failed:", err)
				}
			case syscall.SIGUSR2:
				// attempt to upgrade on SIGUSR2
				if err := upg.Upgrade(); err != nil {
					log.Println("Upgrade failed:", err)
				}
			case syscall.SIGINT, syscall.SIGTERM:
				// terminate on SIGINT or SIGTERM
				log.Println("Terminating...")
				upg.Stop()
			}
		}
	}()

	// start servicing
	go func() {
		err := s.Serve(l)
		if err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	// notify being ready for service
	if err = upg.Ready(); err != nil {
		log.Fatal(err)
	}
	<-upg.Exit()

	if config.GracefulTimeout > 0 {
		// graceful shutdown timeout
		time.AfterFunc(config.GracefulTimeout, func() {
			log.Println("Graceful shutdown timed out")
			os.Exit(1)
		})
	}

	// Wait for connections to drain.
	s.Shutdown(context.Background())
}
