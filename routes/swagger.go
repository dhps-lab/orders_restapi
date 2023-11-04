package routes

import (
	"net/http"

	"github.com/go-openapi/runtime/middleware"
	"github.com/gorilla/mux"
)

func Swagger(r *mux.Router) {
	r.Handle("/swagger.yaml", http.FileServer(http.Dir("./")))
	opts := middleware.SwaggerUIOpts{SpecURL: "/swagger.json"}
	sh := middleware.SwaggerUI(opts, nil)
	r.Handle("/docs", sh)

	// documentation for share
	opts1 := middleware.RedocOpts{SpecURL: "/swagger.json", Path: "docs1"}
	sh1 := middleware.Redoc(opts1, nil)
	r.Handle("/docs1", sh1)
}
