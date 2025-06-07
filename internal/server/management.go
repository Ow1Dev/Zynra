package server

import (
	"net/http"

	"github.com/Ow1Dev/Zynra/pkgs/httpsutils"
)

func connectHandler(w http.ResponseWriter, r *http.Request) {
	httpsutils.Encode(w, http.StatusOK, map[string]string{
		"message": "Connected successfully",
	})
}

func NewManagementServer() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /connect", connectHandler)
	var handler http.Handler = mux
	return handler
}
