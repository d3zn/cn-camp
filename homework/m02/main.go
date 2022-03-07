package main

import (
	"context"
	"log"
	"m02/quark"
	"m02/service"
	"net/http"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

func main() {
	// Create context that listens for the interrupt signal from the OS.
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	conf, err := service.LoadConfig()
	if err != nil {
		log.Fatalf("load config err: %v", err)
	}
	startTime := time.Now()
	router := quark.New()
	router.Use(quark.AccessLog) // middleware
	s := service.NewService(startTime)
	c := service.NewController(s)
	c.Register(router)

	srv := http.Server{
		Addr:    ":" + strconv.Itoa(conf.Service.Port),
		Handler: router,
	}

	go func() {
		log.Printf("mode: %s", conf.Service.Mode)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen err: %v", err)
		}
	}()
	<-ctx.Done()
	stop()
	log.Println("shutting down gracefully, press Ctrl+C again to force")
	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown: ", err)
	}

	log.Println("Server exiting")
}
