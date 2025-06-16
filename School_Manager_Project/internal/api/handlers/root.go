package handlers

import (
	"log"
	"net/http"
	"restapi/utils"
)

func HandleRoot(w http.ResponseWriter, r *http.Request) {
	message := "Hello from the root!"

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")

	if _, err := w.Write([]byte(message)); err != nil {
		log.Println("Error on the server:", err)
		http.Error(w, utils.UnknownInternalServerError.Error(), utils.UnknownInternalServerError.GetStatusCode())
		return
	}

	log.Println("Successfully responded to root request")
}
