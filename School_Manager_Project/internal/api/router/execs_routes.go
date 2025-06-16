package router

import (
	"net/http"
	"restapi/internal/api/handlers"
)

func registerExecs(mux *http.ServeMux) {
	mux.HandleFunc("GET /execs/", handlers.GetExecsHandler)
	mux.HandleFunc("POST /execs/", handlers.PostExecsHandler)
	mux.HandleFunc("PATCH /execs/", handlers.PatchExecsHandler)

	mux.HandleFunc("GET /execs/{id}", handlers.GetOneExecHandler)
	mux.HandleFunc("PATCH /execs/{id}", handlers.PatchOneExecHandler)
	mux.HandleFunc("DELETE /execs/{id}", handlers.DeleteOneExecHandler)

	mux.HandleFunc("POST /execs/{id}/updatePassword", handlers.UpdatePasswordHandler)

	mux.HandleFunc("POST /execs/login", handlers.LoginHandler)
	mux.HandleFunc("POST /execs/logout", handlers.LogoutHandler)
	//mux.HandleFunc("POST /execs/forgotPassword", handlers.GetExecsHandler)
	//mux.HandleFunc("POST /execs/resetPassword/reset/{resetCode}", handlers.GetExecsHandler)
}
