package main

import (
	"flag"
	"fmt"
	"github.com/mdp/qrterminal/v3"
	"log"
	"net"
	"net/http"
	"os"
)

const VERSION = "1.0"

func main() {
	fmt.Printf("LANShare %s by j-jzk. Free software under the BSD license.\n", VERSION)

	// handle command line flags
	allowUploads := flag.Bool("u", false, "whether to allow uploads (default false)")
	help := false
	flag.BoolVar(&help, "help", false, "display help")
	flag.BoolVar(&help, "h", false, "display help")
	port := flag.Int("p", 8080, "the port to listen on")

	flag.Parse()

	if help {
		fmt.Println("lanshare [-u] [-p <port>] | [-h|-help]")
		flag.PrintDefaults()
	} else {
		runServer(*allowUploads, *port)
	}
}

func runServer(allowUploads bool, port int) {
	// we can use Handle("GET /") and HandleFunc("POST /") when we upgrade to a newer version of Go
	downloadHandler := DownloadHandler{allowUploads: allowUploads}
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" && allowUploads {
			HandleUpload(w, r)
		} else {
			downloadHandler.ServeHTTP(w, r)
		}
	})

	fmt.Printf("Listening on port %d.\n", port)
	printAddresses(port)

	err := http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", port), nil)
	if err != nil {
		log.Fatal(err)
	}
}

func printAddresses(port int) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return
	}

	fmt.Println("Open the UI at one of these URLs:")
	var firstViableUrl string
	for _, addr := range addrs {
		ipAddr, ok := addr.(*net.IPNet)
		if ok {
			var ipStr string
			if ipAddr.IP.To4() != nil {
				ipStr = ipAddr.IP.String()
			} else {
				ipStr = fmt.Sprintf("[%s]", ipAddr.IP.String())
			}

			fmt.Printf(" - http://%s:%d", ipStr, port)
			if ipAddr.IP.IsLoopback() {
				fmt.Printf(" (loopback)\n")
			} else {
				// if this is the first non-loopback IP, we will use it for the QR code
				if firstViableUrl == "" {
					firstViableUrl = fmt.Sprintf("http://%s:%d", ipStr, port)
				}
				fmt.Printf("\n")
			}
		}
	}
	fmt.Println("")

	if firstViableUrl != "" {
		qrterminal.GenerateHalfBlock(firstViableUrl, qrterminal.L, os.Stdout)
		fmt.Printf("QR code for: %s\n", firstViableUrl)
		fmt.Println("")
	}
}
