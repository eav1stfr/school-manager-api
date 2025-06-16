package router

import (
	"net/http"
	"restapi/internal/api/handlers"
)

func registerTeacherRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /teachers/", handlers.GetTeachersHandler)
	mux.HandleFunc("POST /teachers/", handlers.PostTeacherHandler)
	mux.HandleFunc("DELETE /teachers/", handlers.DeleteTeachersHandler)
	mux.HandleFunc("PATCH /teachers/", handlers.PatchTeachersHandler)

	mux.HandleFunc("GET /teachers/{id}", handlers.GetOneTeacherHandler)
	mux.HandleFunc("PUT /teachers/{id}", handlers.UpdateTeacherHandler)
	mux.HandleFunc("PATCH /teachers/{id}", handlers.PatchOneTeacherHandler)
	mux.HandleFunc("DELETE /teachers/{id}", handlers.DeleteOneTeacherHandler)

	mux.HandleFunc("GET /teachers/{id}/students", handlers.GetStudentsListForTeacher)
	mux.HandleFunc("GET /teachers/{id}/studentCount", handlers.GetStudentCountForTeacher)
}
