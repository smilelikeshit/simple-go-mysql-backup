package main

import (
	"flag"
	"log"
	"net/http"
)

func main() {
	port := flag.String("p", "8000", "port to serve on")
	flag.Parse()

	http.Handle("/", http.FileServer(http.Dir("./dumps")))


	log.Printf("Serving server on HTTP port: %s\n",*port)
	log.Fatal(http.ListenAndServe(":"+*port, nil))
}