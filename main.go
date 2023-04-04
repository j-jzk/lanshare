package main

import (
	"fmt"
	"net/http"
)

func main() {
	//http.ListenAndServe("0.0.0.0:8080", http.FileServer(http.Dir(".")))
	http.HandleFunc("/", HandleDownload)
	http.HandleFunc("/__lanshare_upload", HandleUpload) // TODO: differentiate using HTTP methods instead of a special URL
	http.ListenAndServe("0.0.0.0:8080", nil)
}
