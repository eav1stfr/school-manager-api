package router

import (
	"net/http"
	"restapi/internal/api/handlers"
)

func Router() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/", handlers.HandleRoot)

	registerTeacherRoutes(mux)

	registerStudentRoutes(mux)

	registerExecs(mux)

	return mux
}
