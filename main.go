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
	host := flag.String("host", "0.0.0.0", "the host to listen on")
	port := flag.Int("p", 8080, "the port to listen on")

	flag.Parse()

	// validate host parameter
	if net.ParseIP(*host) == nil {
		log.Fatalln("Invalid host ip address")
	}

	if help {
		fmt.Println("lanshare [-u] [-host <host>] [-p <port>] | [-h|-help]")
		flag.PrintDefaults()
	} else {
		runServer(*allowUploads, *host, *port)
	}
}

func runServer(allowUploads bool, host string, port int) {
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
	printAddresses(host, port)

	err := http.ListenAndServe(fmt.Sprintf("%s:%d", host, port), nil)
	if err != nil {
		log.Fatal(err)
	}
}

func printAddresses(host string, port int) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return
	}

	// there is a check if using a custom host
	// this replaces array addrs with a single IP network with mask /32 (IPv4) and /128 (IPv6)
	if host != "0.0.0.0" {
		ip := net.ParseIP(host)

		if ip == nil {
			log.Fatalf("Invalid host ip address: %s\n", host)
		}

		addrs = []net.Addr{&net.IPNet{IP: ip, Mask: net.CIDRMask(32, 32)}}
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
