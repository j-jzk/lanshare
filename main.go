package main

import (
	"log"
	"flag"
	"fmt"
	"net/http"
)

func main() {
	// handle command line flags
	allowUploads := flag.Bool("u", false, "whether to allow uploads (default false)")
	help := false
	flag.BoolVar(&help, "help", false, "display help")
	flag.BoolVar(&help, "h", false, "display help")

	flag.Parse()

	if help {
		fmt.Println("lanshare -u | [-h|-help]")
		flag.PrintDefaults()
	} else {
		runServer(*allowUploads)
	}
}

func runServer(allowUploads bool) {
	http.Handle("/", &DownloadHandler{allowUploads: allowUploads})
	if allowUploads {
		http.HandleFunc("/__lanshare_upload", HandleUpload) // TODO: differentiate using HTTP methods instead of a special URL
	}

	err := http.ListenAndServe("0.0.0.0:8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}
