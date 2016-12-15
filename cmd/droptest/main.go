package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"vallon.me/dropfs"
)

var APIKey string

func init() {
	flag.StringVar(&APIKey, "token", "", "required to authenticate your account with the server")

	flag.Parse()

	if APIKey == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}
}

func main() {
	dfs := dropfs.NewFS(APIKey, "/public")
	log.Fatal(http.ListenAndServe(":8000", http.FileServer(dfs)))
}
