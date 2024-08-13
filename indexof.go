package main

import (
	"context"
	"flag"
	"fmt"
	"indexof/config"
	"indexof/files/indexof"
	"indexof/handle"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

var (
	configFile string
	err        error
)

func init() {
	flag.StringVar(&configFile, "f", "config.json", "User-defined configuration files")
}

func main() {
	flag.Parse()
	err = config.LoadConfig(configFile)
	if err != nil {
		log.Fatalf("load config err=%v", err)
		return
	}
	if err := indexof.LoadIndexOf(); err != nil {
		log.Fatalln(err)
		return
	}
	server := http.Server{
		Addr:    fmt.Sprintf(":%d", config.Global.Port),
		Handler: handle.New(),
	}
	go func() {
		if err := server.ListenAndServe(); err != nil {
			log.Fatalf("service start failed %v", err)
		}
	}()
	log.Println(fmt.Sprintf("server running on port http://0.0.0.0:%d", config.Global.Port))
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	s := <-c
	log.Println(fmt.Sprintf("Receiving the signal %s", s))
	ctx, cancelFunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFunc()
	if err := server.Shutdown(ctx); err != nil {
		log.Println("server shutdown fail")
	}
	log.Println("server exit")
}
