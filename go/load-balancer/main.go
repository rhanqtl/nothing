package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"
)

/*
 * iota is a keyword in Go representing an incrementing
 *+ sequence (from 0, i.e., zero-indexed), it can be
 *+ used in expressions, and all constants afterwards
 *+ will have a value evaluated from the expression
 *+ where iota is replaced by the corresponding value
 *+ for that line.
 *
 * For example,
 *   const (
 *       Zero int = iota  // Zero  -> 0
 *       One              // One   -> 1
 *       _                // We don't need the name
 *       Three            // Three -> 3
 *   )
 *
 *   const (
 *       _  = iota              // just ignore the value
 *       KB = 1 << (10 * iota)  // 1KB = 2^(10 * 1)B
 *       MB                     // 1MB = 2^(10 * 2)B
 *       GB                     // 1GB = 2^(10 * 3)B
 *       TB                     // 1TB = 2^(10 * 4)B
 *   )
 *
 * Note: the value is evaluated *per-line*.
 *   const (
 *       Zero, One  = iota, iota + 1  // Zero -> 0, One   -> 1
 *       Two, Three                   // Two  -> 1, Three -> 2
 *       Four                         // Four -> 2
 *   )
 */
const (
	// Attempts -> 0
	Attempts int = iota
	// Retry    -> 1
	Retry
)

var balancer ServerPool

// GetAttemptsFromContext returns the attempts for request
func getAttemptsFromContext(r *http.Request) int {
	if attempts, ok := r.Context().Value(Attempts).(int); ok {
		return attempts
	}
	return 1
}

// GetAttemptsFromContext returns the attempts for request
func getRetryFromContext(r *http.Request) int {
	if retry, ok := r.Context().Value(Retry).(int); ok {
		return retry
	}
	return 0
}

// loadBalance load balances the incoming request
func loadBalance(w http.ResponseWriter, r *http.Request) {
	attempts := getAttemptsFromContext(r)
	if attempts > 3 {
		log.Printf("%s(%s) Max attempts reached, terminating\n", r.RemoteAddr, r.URL.Path)
		http.Error(w, "Service not available", http.StatusServiceUnavailable)
		return
	}

	peer, err := balancer.schedule()
	if err != nil {
		http.Error(w, "Service not available", http.StatusServiceUnavailable)
		return
	}
	peer.serve(w, r)
}

func healthCheck() {
	tick := time.NewTicker(9 * time.Second)
	for {
		select {
		case <-tick.C:
			log.Println("Starting health check")
			balancer.doHealthCheck()
			log.Println("Health check completed")
		}
	}
}

func main() {
	var serverList string
	var port int
	flag.StringVar(&serverList, "servers", "", "use comma (,) to separate")
	flag.IntVar(&port, "port", 9527, "port to serve")
	flag.Parse()

	if len(serverList) == 0 {
		log.Fatal("Please provide at least 1 server")
	}

    tokens := strings.Split(serverList, ",")
    for _, tok := range tokens {
    	theUrl, err := url.Parse(tok)
    	if err != nil {
    		log.Fatal(err)
		}
		proxy := httputil.NewSingleHostReverseProxy(theUrl)
		proxy.ErrorHandler = func(writer http.ResponseWriter, req *http.Request, e error) {
			log.Printf("Server %s has error: %s\n", theUrl.Host, e)
			retries := getRetryFromContext(req)
			if retries < 3 {
				select {
				case <-time.After(10 * time.Millisecond):
					ctx := context.WithValue(req.Context(), Retry, retries + 1)
					proxy.ServeHTTP(writer, req.WithContext(ctx))
				}
				return
			}

			balancer.setServerState(theUrl, false)

			attempts := getAttemptsFromContext(req)
			log.Printf("%s(%s) Attempting retry %d\n", req.RemoteAddr, req.URL.Path, attempts)
			ctx := context.WithValue(req.Context(), Attempts, attempts + 1)
			loadBalance(writer, req.WithContext(ctx))
		}

		balancer.addServer(&Server{
			theUrl: theUrl,
			alive: true,
			reverseProxy: proxy,
		})
		log.Printf("Configured server: %s\n", theUrl)
	}

	// create http server
	server := http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: http.HandlerFunc(loadBalance),
	}

	// start health checking
	go healthCheck()

	log.Printf("Load Balancer started at :%d\n", port)
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
