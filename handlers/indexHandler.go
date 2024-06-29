package handlers

import (
	"fmt"
	"net/http"
	"simple-web-server/utils"
	"strings"
)

type IndexHandler struct {
	logger utils.Logger
}

// NewCustomFileServer creates a new CustomFileServer
func NewIndexHandler(logger utils.Logger) *IndexHandler {
	return &IndexHandler{logger: logger}
}

func (ih *IndexHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		if r.URL.Path != "/" {
			response := "Not found"
			http.Error(w, response, http.StatusNotFound)
			return
		}
		ih.logger.Info(fmt.Sprintf("%s %s %s", r.Method, r.RemoteAddr, r.URL))
		fmt.Fprint(w, "OK")
	case "POST":
		err := r.ParseForm()
		if err != nil {
			// If there was an error, respond with a 400 Bad Request status and an error message
			response := "Error parsing form"
			http.Error(w, response, http.StatusBadRequest)
			return
		}
		formValues := ""
		for key, values := range r.Form {
			formValues += " " + key + "=" + strings.Join(values, ":")
		}
		ih.logger.Info(fmt.Sprintf("%s %s %s%s", r.Method, r.RemoteAddr, r.URL, formValues))
		fmt.Fprint(w, "OK")
	default:
		strResponse := "Method not allowed"
		ih.logger.Warning(fmt.Sprintf("%s BLOCKED %s %s", r.Method, r.RemoteAddr, r.URL))
		http.Error(w, strResponse, http.StatusMethodNotAllowed)
	}
}
