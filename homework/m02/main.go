package main

import (
	"log"
	"m02/quark"
	"m02/service"
	"net/http"
)

func main() {
	router := quark.New()
	router.Use(quark.AccessLog) // middleware
	s := service.NewService()
	c := service.NewController(s)
	c.Register(router)
	log.Fatal(http.ListenAndServe(":8080", router))
}
