package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
)

func HandleUpload(w http.ResponseWriter, r *http.Request) {
	// dir := safeDir(r.FormValue("dir"))
	dir := safeDir(r.URL.Path)

	log.Printf("UPLOAD - files:\n")
	r.ParseMultipartForm(32 << 10) // max memory of 32 MB
	for _, header := range r.MultipartForm.File["file"] {
		infile, err := header.Open()
		if handleError(err, w, "") {
			return
		}

		// TODO: rename the file if it already exists?
		localPath := path.Join(".", dir, header.Filename)
		outfile, err := os.Create(localPath)
		if handleError(err, w, localPath) {
			return
		}

		_, err = io.Copy(outfile, infile)
		if handleError(err, w, localPath) {
			return
		}

		log.Printf("\t%s - success\n", localPath)
	}

	w.Header().Add("Location", "/"+dir)
	w.WriteHeader(303)
}

func safeDir(unsafe string) string {
	unsafe = strings.TrimPrefix(unsafe, "/")
	unsafe = strings.ReplaceAll(unsafe, "..", "")
	unsafe = strings.ReplaceAll(unsafe, "//", "/")
	// TODO: notify the user if the dir was changed
	return unsafe
}

func handleError(err error, w http.ResponseWriter, filePath string) bool {
	if err == nil {
		return false
	} else {
		w.WriteHeader(500)
		fmt.Fprintf(w, "Error: %s\n", err)
		log.Printf("UPLOAD %s - ERROR: %s\n", filePath, err)
		return true
	}
}
