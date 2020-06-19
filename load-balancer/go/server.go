package main

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
)

/**
 * the info of servers the load balancer should know
 */
type Server struct {
	theUrl       *url.URL
	alive        bool
	mux          sync.RWMutex
	reverseProxy *httputil.ReverseProxy
}

func (s *Server) setAlive(alive bool) {
	s.mux.Lock()
	s.alive = alive
	s.mux.Unlock()
}

func (s *Server) isAlive() bool {
	s.mux.RLock()
	defer s.mux.RUnlock()
	return s.alive
}

func (s *Server) serve(rw http.ResponseWriter, req *http.Request) {
	s.reverseProxy.ServeHTTP(rw, req)
}
