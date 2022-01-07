package main

import (
	"flag"
	"log"
	"net/http"
)

func main() {
	port := flag.String("p", "8080", "port to serve on")
	flag.Parse()

	fs := http.FileServer(http.Dir("./site/"))
	http.Handle("/", http.StripPrefix("/", fs))

	log.Printf("Serving on HTTP port: %s\n", *port)
	log.Fatal(http.ListenAndServe(":"+*port, nil))
}
