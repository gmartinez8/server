package serverservice

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"
)

//Router struct to handle multiples methods
//To store a http.HandlerFunc in a key value map
//must match path and method
//map used to assert if pathExist or methodExist
type Router struct {
	defaultRules map[string]map[string]http.HandlerFunc
}

//NewRouter creates a new Router type
func NewRouter() *Router {
	return &Router{
		defaultRules: make(map[string]map[string]http.HandlerFunc),
	}
}

//checkPathRegex check if Path for a DELETE and PUT request
//are valid and have a number and not a word for this example
// url= /users/123 is valid
// url= /users/asd is not valid
func checkPathRegex(path string) (string, bool) {
	regexpPath := strings.Split(path, "/")
	index := len(regexpPath)
	match, err := regexp.MatchString("([0-9]+)", regexpPath[index-1])
	if err != nil {
		fmt.Println(err)
	}
	if match == true {
		path = strings.Replace(path, "/"+regexpPath[index-1], "", 1)
	}
	return path, match
}

//checkGetPathRegex check Path for a GET request
//and clean path if needed
func checkGetPathRegex(path string) (string, bool) {
	regexpPath := strings.Split(path, "/")
	index := len(regexpPath)
	match, err := regexp.MatchString("([0-9]+)", regexpPath[index-1])
	if err != nil {
		fmt.Println(err)
	}
	if match == true {
		//clean the path for GET /{path}/{id}
		//so we can return GET /{path}
		path = strings.Replace(path, "/"+regexpPath[index-1], "", 1)
		return path, match
	}
	//Returns the original path for GET /{path}
	return path, true
}

//FindHandler finds the handler assigned to a route example GET /api
func (rt *Router) FindHandler(path string, method string) (http.HandlerFunc, bool, bool) {
	_, allowedPath := rt.defaultRules[path]
	if method == "DELETE" || method == "PUT" {
		path, allowedPath = checkPathRegex(path)
	}
	if method == "GET" {
		path, allowedPath = checkGetPathRegex(path)
	}
	handlerFunc, allowedMethod := rt.defaultRules[path][method]
	return handlerFunc, allowedMethod, allowedPath
}

//AllowedPath verifies if path is allowed
func (rt *Router) AllowedPath(path string) bool {
	_, allowedPath := rt.defaultRules[path]
	return allowedPath
}

//ServeHTTP for HandlerFunc to check if our Path and Method are allowed
//this is executed on each request
func (rt *Router) ServeHTTP(w http.ResponseWriter, request *http.Request) {
	handlerFunc, allowedMethod, allowedPath := rt.FindHandler(request.URL.Path, request.Method)
	if !allowedPath {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if !allowedMethod {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	handlerFunc(w, request)
}
