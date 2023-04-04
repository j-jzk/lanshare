package main

import (
	"errors"
	"fmt"
	"html/template"
	"io/fs"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
)

func HandleDownload(w http.ResponseWriter, r *http.Request) {
	// open the requested file
	fsPath := path.Join(".", r.URL.Path)
	f, err := os.Open(fsPath)
	if handleFileError(err, w) {
		return
	}
	defer f.Close()

	// check if the file is a dir
	stat, err := f.Stat()
	if handleFileError(err, w) {
		return
	}

	if stat.IsDir() {
		writeDirectoryListing(w, fsPath)
	} else {
		defaultFileServer := http.FileServer(http.Dir("."))
		defaultFileServer.ServeHTTP(w, r)
	}
}

// DIRECTORY LISTING
const dirListingTemplate = `
<!DOCTYPE html><html>
<head>
	<meta charset="UTF-8" />
	<style>
		body {
			font-family: sans-serif;
		}
		footer {
			margin-top: 20px;
			font-style: italic;
		}
	</style>
</head><body>
	<h1><code>{{.Cwd}}</code></h1>
	<form method="POST" enctype="multipart/form-data" action="/__lanshare_upload">
		<input type="file" name="file" />
		<input type="hidden" name="dir" value="{{.Cwd}}" />
		<button type="submit">Upload</button>
	</form>
	<table>
		<thead>
			<tr><th>name</th><th>size</th></tr>
		</thead>
		<tbody>
			{{range .Entries}}
				<tr>
					<td><a href="{{.Name}}">{{.Name}}</a>
					<td>{{.SizeInfo}}</td>
				</tr>
			{{end}}
		</tbody>
	</table>
	<footer>
		LANShare
	</footer>
</body></html>
`

type listingTemplateParams struct {
	Cwd     string
	Entries []templateListEntry
}
type templateListEntry struct {
	Name     string
	SizeInfo string
}

func writeDirectoryListing(w http.ResponseWriter, path string) {
	if !strings.HasSuffix(path, "/") {
		path += "/"
	}

	dirEntries, err := os.ReadDir(path)
	if handleFileError(err, w) {
		return
	}

	templEntries := make([]templateListEntry, len(dirEntries)+1)
	templEntries[0] = templateListEntry{Name: "..", SizeInfo: "DIR"}
	for i, f := range dirEntries {
		stat, _ := f.Info()

		var sizeInfo, name string
		if stat.IsDir() {
			sizeInfo = "DIR"
			name = stat.Name() + "/"
		} else {
			sizeInfo = strconv.FormatInt(stat.Size(), 10)
			name = stat.Name()
		}

		//templEntries = append(templEntries, templateListEntry{Name: stat.Name(), SizeInfo: sizeInfo})
		templEntries[i+1] = templateListEntry{Name: name, SizeInfo: sizeInfo}
	}

	templ := template.Must(template.New("").Parse(dirListingTemplate))
	templ.Execute(w, listingTemplateParams{Cwd: path, Entries: templEntries})
}

// UTIL
func handleFileError(err error, w http.ResponseWriter) bool {
	if err == nil {
		return false
	} else {
		if errors.Is(err, fs.ErrNotExist) {
			w.WriteHeader(404)
			fmt.Fprintf(w, "404 not found\n")
		} else {
			w.WriteHeader(500)
			fmt.Fprintf(w, "Unexpected error: %s\n", err)
		}

		return true
	}
}
