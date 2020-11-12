//Package server implements a simple library for a http server
//You can create a server, add the routes you want to handle
package server

import (
	"context"
	"log"
	"net/http"
)

//Server struct allows us to define multiples servers if needed
//Each server will have its own Router to avoid conflicts
type Server struct {
	port   string
	router *Router
}

//NewServer creates a new Server
//and asign a NewRouter
func NewServer(port string) *Server {
	return &Server{
		port:   port,
		router: NewRouter(),
	}
}

//Run starts the server and asign *Router  to handle the routes
//Router its a map[string][string]http.HandlerFunc
//map[path][method]http.HandlerFunc
func (sr *Server) Run(ctx context.Context) error {
	log.Printf("HTTP Server is starting to listen on 0.0.0.0%s", sr.port)
	http.Handle("/", sr.router)
	//creating a type Server
	s := &http.Server{Addr: "0.0.0.0" + sr.port}
	ch := make(chan error)

	go func(s *http.Server, ch chan error) {
		ch <- s.ListenAndServe()
	}(s, ch)

	select {
	case err := <-ch:
		log.Printf("server returned and error: %v", err)
		return err
	case <-ctx.Done():
		//if parent context is Done this will be executed
		//log.Println("context is canceled")
		return s.Close()
	}

}

//Handle defines/register the routes i want to handle
//also asign each route a HandlerFunc to handle it
//you can define each HandlerFunc in handlers.go file
func (sr *Server) Handle(path string, method string, handler http.HandlerFunc) {
	//Check if the path already exists
	if !sr.router.AllowedPath(path) {
		//If not path then create a new one
		sr.router.defaultRules[path] = make(map[string]http.HandlerFunc)
	}
	sr.router.defaultRules[path][method] = handler
}
