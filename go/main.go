package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/julienschmidt/httprouter"
)

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "Welcome!\n")
}

func Hello(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"msg": "hello, %s!", "obj": {}, "num": 123, "arr": [123]}`, ps.ByName("name"))
}

func main() {
	sigs := make(chan os.Signal, 1)
	run("localhost:8080", sigs)
}

func run(addr string, sigs chan os.Signal) {
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	router := httprouter.New()
	router.GET("/", Index)
	router.GET("/hello/:name", Hello)

	srv := &http.Server{Addr: addr, Handler: router}
	go func() {
		err := srv.ListenAndServe()
		if err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	sig := <-sigs
	log.Printf("Received signal %s, shutting down server\n", sig)
}
