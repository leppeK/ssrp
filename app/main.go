package main

import (
	"crypto/tls"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/spf13/pflag"
)

var (
	allowedGroups = make(map[string]bool)
	targetURL     string
	listenAddr    string
	insecure      bool
	groups        []string
)

func main() {
	pflag.StringVarP(&targetURL, "target", "t", "localhost:8080", "Target URL to proxy to")
	pflag.StringVarP(&listenAddr, "listen", "l", ":3000", "Address to listen on")
	pflag.BoolVarP(&insecure, "insecure", "i", false, "Ignore SSL certificate errors")
	pflag.StringSliceVarP(&groups, "group", "g", []string{}, "Allowed groups (can be specified multiple times)")

	pflag.Parse()

	// Process groups into allowedGroups map
	for _, group := range groups {
		allowedGroups[group] = true
	}

	if len(allowedGroups) == 0 {
		log.Fatal("At least one group must be specified using --group or -g")
	}

	if !strings.Contains(targetURL, "://") {
		targetURL = "http://" + targetURL
	}

	targetParsed, err := url.Parse(targetURL)
	if err != nil {
		log.Fatal("Invalid target URL:", err)
	}

	proxy := httputil.NewSingleHostReverseProxy(targetParsed)

	// Optional: Handle insecure TLS connections
	if insecure {
		proxy.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
	}

	// Custom director
	proxy.Director = func(req *http.Request) {
		req.URL.Scheme = targetParsed.Scheme
		req.URL.Host = targetParsed.Host
		req.Host = targetParsed.Host
		req.Header.Add("X-Forwarded-By", "secure-proxy")
	}

	http.HandleFunc("/", handler(proxy))
	log.Printf("Starting proxy server on %s", listenAddr)
	log.Printf("Authorized requests will be forwarded to %v", targetURL)
	log.Fatal(http.ListenAndServe(listenAddr, nil))
}

func handler(proxy *httputil.ReverseProxy) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		groups := strings.Split(r.Header.Get("X-Groups"), ",")
		authorized := false
		for _, group := range groups {
			group = strings.TrimSpace(group)
			if allowedGroups[group] {
				authorized = true
				break
			}
		}

		if !authorized {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		proxy.ServeHTTP(w, r)
	}
}
