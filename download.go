package main

import (
	"errors"
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
	"github.com/dustin/go-humanize"
)

type DownloadHandler struct {
	// shouldn't this be AllowUploads
	allowUploads bool
}

func (dh *DownloadHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// open the requested file
	fsPath := path.Join(".", r.URL.Path)
	f, err := os.Open(fsPath)
	if handleFileError(err, w, fsPath) {
		return
	}
	defer f.Close()

	// check if the file is a dir
	stat, err := f.Stat()
	if handleFileError(err, w, fsPath) {
		return
	}

	if stat.IsDir() {
		log.Printf("GET %s - directory listing\n", fsPath)
		writeDirectoryListing(w, fsPath, dh.allowUploads)
	} else {
		log.Printf("GET %s - file\n", fsPath)
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
	{{if .AllowUploads}}
		<form method="POST" enctype="multipart/form-data" action="/__lanshare_upload">
			<input type="file" name="file" />
			<input type="hidden" name="dir" value="{{.Cwd}}" />
			<button type="submit">Upload</button>
		</form>
	{{end}}
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
		LANShare {{.ServerVersion}} by <a href="https://j-jzk.cz">j-jzk</a>.
	</footer>
</body></html>
`

type listingTemplateParams struct {
	Cwd     string
	AllowUploads bool
	Entries []templateListEntry
	ServerVersion string
}
type templateListEntry struct {
	Name     string
	SizeInfo string
}

func writeDirectoryListing(w http.ResponseWriter, path string, allowUploads bool) {
	if !strings.HasSuffix(path, "/") {
		path += "/"
	}

	dirEntries, err := os.ReadDir(path)
	if handleFileError(err, w, path) {
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
			sizeInfo = humanize.Bytes(uint64(stat.Size()))
			name = stat.Name()
		}

		//templEntries = append(templEntries, templateListEntry{Name: stat.Name(), SizeInfo: sizeInfo})
		templEntries[i+1] = templateListEntry{Name: name, SizeInfo: sizeInfo}
	}

	templ := template.Must(template.New("").Parse(dirListingTemplate))
	templ.Execute(w, listingTemplateParams{Cwd: path, Entries: templEntries, AllowUploads: allowUploads, ServerVersion: VERSION})
}

// UTIL
func handleFileError(err error, w http.ResponseWriter, path string) bool {
	if err == nil {
		return false
	} else {
		if errors.Is(err, fs.ErrNotExist) {
			w.WriteHeader(404)
			fmt.Fprintf(w, "404 not found\n")
			log.Printf("GET %s - 404\n", path)
		} else {
			w.WriteHeader(500)
			fmt.Fprintf(w, "Unexpected error: %s\n", err)
			log.Printf("GET %s - ERROR: %s\n", path, err)
		}

		return true
	}
}
