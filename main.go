package main

import (
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", HandleDownload)
	http.HandleFunc("/__lanshare_upload", HandleUpload) // TODO: differentiate using HTTP methods instead of a special URL

	err := http.ListenAndServe("0.0.0.0:8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}
