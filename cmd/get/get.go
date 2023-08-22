package main

import (
	"context"
	"crypto/tls"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

var (
	fileName    string
	fullURLFile string
)

func main() {

	fullURLFile := flag.String("url", "https://ftp.mgts.by/test/100Mb", "url to download")
	insecure := flag.Bool("insecure", false, "set true for ignore tls host check")
	dns := flag.String("dns", "", "set dns ip address, default use system dns resolver")
	flag.Parse()

	if *dns != "" {
		dnsResolverIP := fmt.Sprintf("%v:53", *dns) // DNS resolver.
		dnsResolverProto := "udp"                   // Protocol to use for the DNS resolver
		dnsResolverTimeoutMs := 5000                // Timeout (ms) for the DNS resolver (optional)
		fmt.Printf("dns = '%v://%v'\n", dnsResolverProto, dnsResolverIP)

		dialer := &net.Dialer{
			Resolver: &net.Resolver{
				PreferGo: true,
				Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
					d := net.Dialer{
						Timeout: time.Duration(dnsResolverTimeoutMs) * time.Millisecond,
					}
					return d.DialContext(ctx, dnsResolverProto, dnsResolverIP)
				},
			},
		}

		dialContext := func(ctx context.Context, network, addr string) (net.Conn, error) {
			return dialer.DialContext(ctx, network, addr)
		}

		http.DefaultTransport.(*http.Transport).DialContext = dialContext
	}

	// Build fileName from fullPath
	fileURL, err := url.Parse(*fullURLFile)
	if err != nil {
		log.Fatal(err)
	}
	path := fileURL.Path
	segments := strings.Split(path, "/")
	fileName = segments[len(segments)-1]

	// Create blank file
	file, err := os.Create(fileName)
	if err != nil {
		log.Fatal(err)
	}
	client := http.Client{
		CheckRedirect: func(r *http.Request, via []*http.Request) error {
			r.URL.Opaque = r.URL.Path
			return nil
		},
	}

	if *insecure {
		http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	}

	// Put content on file
	resp, err := client.Get(*fullURLFile)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	size, err := io.Copy(file, resp.Body)

	defer file.Close()

	fmt.Printf("Downloaded a file %s with size %d", fileName, size)

}
