package main
import (
	"fmt"
	"net/http"
)

func main() {
	fmt.Println("ahoj")

	//http.ListenAndServe("0.0.0.0:8080", http.FileServer(http.Dir(".")))
	http.HandleFunc("/", HandleDownload)
	http.ListenAndServe("0.0.0.0:8080", nil)
}
