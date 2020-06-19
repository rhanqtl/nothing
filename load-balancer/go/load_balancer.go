package main

import (
	"log"
	"net"
	"net/url"
	"sync"
	"time"
)

/**
 * The core responsibility of ServerPool, a.k.a. LoadBalancer
 *+ is selecting a server to handle the request. Therefore,
 *+ the load balancer should have knowledge about where the
 *+ server is (URL).
 *
 * Which server is selected depends on the policy.
 *
 * However, the load balancer must not relay a request to a
 *+ non-available server, so it have to know the status of
 *+ each server (alive).
 *
 * To know about whether or not a server is available, the
 *+ load balancer should check it periodically.
 */

/**
 * servers
 */
type ServerPool struct {
	pool    []*Server
	current uint
	mux     sync.Mutex
}

func (sp *ServerPool) addServer(srv *Server) {
	if srv != nil {
		sp.pool = append(sp.pool, srv)
	}
}

func (sp *ServerPool) removeServer(theURL *url.URL) {
	targetIndex := -1
	for idx, srv := range sp.pool {
		if srv.theUrl.String() == theURL.String() {
			targetIndex = idx
		}
	}
	if targetIndex != -1 {
		sp.pool = append(sp.pool[: targetIndex], sp.pool[targetIndex + 1:]...)
	}
}

func (sp *ServerPool) schedule() (*Server, error) {
	sp.mux.Lock()
	defer sp.mux.Unlock()
	for sp.pool[sp.current].isAlive() {
		target := sp.pool[sp.current]
		sp.current = (sp.current + 1) % uint(len(sp.pool))
		return target, nil
	}
	return nil, NoAvailableServerError("no available server")
}

func (sp *ServerPool) setServerState(theUrl *url.URL, alive bool) {
	for _, srv := range sp.pool {
		if srv.theUrl.String() == theUrl.String() {
			srv.setAlive(alive)
			break
		}
	}
}

func (sp *ServerPool) doHealthCheck() {
	for _, srv := range sp.pool {
		status := "up"
		alive := isServerAlive(srv.theUrl)
		srv.setAlive(alive)
		if !alive {
			status = "down"
		}
		log.Printf("Server %s is %s\n", srv.theUrl, status)
	}
}

func isServerAlive(theUrl *url.URL) bool {
	timeout := 2 * time.Second
	conn, err := net.DialTimeout("tcp", theUrl.Host, timeout)
	if err != nil {
		log.Println("Site unreachable, error: ", err)
		return false
	}
	conn.Close()
	return true
}
