package server

import (
	"net/http"
	"path"
	"strings"

	"github.com/Ow1Dev/Zynra/pkgs/httpsutils"
	"github.com/rs/zerolog/log"
)

func getRoot(w http.ResponseWriter, r *http.Request) {
	cleanPath := path.Clean(r.URL.Path)
	cleanPath = strings.Trim(cleanPath, "/")

	segments := strings.Split(cleanPath, "/")

	log.Debug().Any("segments", segments).Msg("Request path segments")

	if len(segments) != 1 || segments[0] == "" {
		http.Error(w, "Invalid path", http.StatusBadRequest)
		return
	}

	action := segments[0]
	if(action != "echo") {
		http.Error(w, "Only 'echo' is allowed", http.StatusBadRequest)
		return
	}

	httpsutils.Encode(w, http.StatusOK, map[string]string{
		"message": "Hello, World!",
	})
}

func NewRouterServer() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /", getRoot)
	var handler http.Handler = mux
	return handler
}
