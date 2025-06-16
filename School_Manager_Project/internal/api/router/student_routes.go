package router

import (
	"net/http"
	"restapi/internal/api/handlers"
)

func registerStudentRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /students/", handlers.GetStudentsHandler)
	mux.HandleFunc("POST /students/", handlers.PostStudentHandler)
	mux.HandleFunc("DELETE /students/", handlers.DeleteStudentsHandler)
	mux.HandleFunc("PATCH /students/", handlers.PatchStudentsHandler)

	mux.HandleFunc("GET /students/{id}", handlers.GetOneStudentHandler)
	mux.HandleFunc("PUT /students/{id}", handlers.UpdateStudentHandler)
	mux.HandleFunc("PATCH /students/{id}", handlers.PatchOneStudentHandler)
	mux.HandleFunc("DELETE /students/{id}", handlers.DeleteOneStudentHandler)
}
