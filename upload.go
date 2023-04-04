package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"strings"
)

func HandleUpload(w http.ResponseWriter, r *http.Request) {
	dir := safeDir(r.FormValue("dir"))
	infile, header, err := r.FormFile("file")
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, "Error: %s\n", err)
		return
	}

	// TODO: rename the file if it already exists?
	outfile, err := os.Create(path.Join(".", dir, header.Filename))
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, "Error: %s\n", err)
		return
	}

	_, err = io.Copy(outfile, infile)
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, "Error: %s\n", err)
		return
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
