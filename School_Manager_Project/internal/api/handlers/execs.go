package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"restapi/internal/models"
	"restapi/internal/sqlconnect"
	"restapi/utils"
	"strconv"
	"time"
)

func GetOneExecHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	realID, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, utils.InvalidIdError.Error(), utils.InvalidIdError.GetStatusCode())
		return
	}

	exec, err := sqlconnect.GetExecByID(realID)
	if err != nil {
		if appErr, ok := err.(*utils.AppErrors); ok {
			http.Error(w, appErr.Error(), appErr.GetStatusCode())
			return
		}
		http.Error(w, "unknown internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(exec)
	if err != nil {
		http.Error(w, "Error encoding data", http.StatusInternalServerError)
		return
	}
}

func GetExecsHandler(w http.ResponseWriter, r *http.Request) {
	execList, err := sqlconnect.GetExecsDBHandler(r)
	if err != nil {
		if appErr, ok := err.(*utils.AppErrors); ok {
			http.Error(w, appErr.Error(), appErr.GetStatusCode())
			return
		}
		http.Error(w, "unknown internal server error", http.StatusInternalServerError)
		return
	}
	response := struct {
		Status string        `json:"status"`
		Count  int           `json:"count"`
		Data   []models.Exec `json:"data"`
	}{
		Status: "success",
		Count:  len(execList),
		Data:   execList,
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, "Error processing list of execs", http.StatusInternalServerError)
	}
}

func PostExecsHandler(w http.ResponseWriter, r *http.Request) {
	var newExecs []models.Exec
	err := json.NewDecoder(r.Body).Decode(&newExecs)
	if err != nil {
		http.Error(w, "Error decoding data", http.StatusBadRequest)
		return
	}

	err = utils.ValidateExecPost(newExecs)
	if err != nil {
		if appErr, ok := err.(*utils.AppErrors); ok {
			log.Println("ERROR 1:", err)
			http.Error(w, appErr.Error(), appErr.GetStatusCode())
			return
		}
		http.Error(w, "unknown internal server error", http.StatusInternalServerError)
		return
	}

	addedExecs, err := sqlconnect.AddExecsDBHandler(newExecs)
	if err != nil {
		if appErr, ok := err.(*utils.AppErrors); ok {
			log.Println("ERROR 2:", err)
			http.Error(w, appErr.Error(), appErr.GetStatusCode())
			return
		}
		http.Error(w, "unknown internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	response := struct {
		Status string        `json:"status"`
		Count  int           `json:"count"`
		Data   []models.Exec `json:"data"`
	}{
		Status: "success",
		Count:  len(addedExecs),
		Data:   addedExecs,
	}
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Println("ERROR 3:", err)
		http.Error(w, "Error encoding error", http.StatusInternalServerError)
		return
	}
}

// PatchExecsHandler PATCH /execs/
func PatchExecsHandler(w http.ResponseWriter, r *http.Request) {
	var updates []map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&updates)
	if err != nil {
		http.Error(w, "error decoding data", http.StatusBadRequest)
		return
	}
	err = sqlconnect.PatchExecsDBHandler(updates)
	if err != nil {
		if appErr, ok := err.(*utils.AppErrors); ok {
			http.Error(w, appErr.Error(), appErr.GetStatusCode())
			return
		}
		http.Error(w, "unknown internal server error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// PatchOneTeacherHandler patch method - only update received fields
func PatchOneExecHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, utils.InvalidIdError.Error(), utils.InvalidIdError.GetStatusCode())
		return
	}

	var updates map[string]interface{}
	err = json.NewDecoder(r.Body).Decode(&updates)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	existingExec, err := sqlconnect.PatchOneExecDBHandler(id, updates)
	if err != nil {
		if appErr, ok := err.(*utils.AppErrors); ok {
			http.Error(w, appErr.Error(), appErr.GetStatusCode())
			return
		}
		http.Error(w, "unknown internal server error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(existingExec)
}

func DeleteOneExecHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, utils.InvalidIdError.Error(), utils.InvalidIdError.GetStatusCode())
		return
	}
	err = sqlconnect.DeleteOneExecDBHandler(id)
	if err != nil {
		if appErr, ok := err.(*utils.AppErrors); ok {
			http.Error(w, appErr.Error(), appErr.GetStatusCode())
			return
		}
		http.Error(w, "unknown internal server error", http.StatusInternalServerError)
		return
	}
	// typically for delete requests use the line below (commented)
	// w.WriteHeader(http.StatusNoContent)

	// otherwise, if you want to send a response body, use the code below
	w.Header().Set("Content-Type", "application/json")
	response := struct {
		Status string `json:"status"`
		ID     int    `json:"id"`
	}{
		Status: "Exec successfully deleted",
		ID:     id,
	}
	json.NewEncoder(w).Encode(response)
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var req models.Exec
	// data validation
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, utils.InvalidRequestBodyError.Error(), utils.InvalidRequestBodyError.GetStatusCode())
		return
	}
	defer r.Body.Close()

	if req.Username == "" || req.Password == "" {
		http.Error(w, utils.MissingFieldsError.Error(), utils.MissingFieldsError.GetStatusCode())
		return
	}

	// search for user if user actually exists
	user, err := sqlconnect.LoginDBHandler(req.Username)
	if err != nil {
		if appErr, ok := err.(*utils.AppErrors); ok {
			http.Error(w, appErr.Error(), appErr.GetStatusCode())
			return
		}
		http.Error(w, utils.UnknownInternalServerError.Error(), utils.UnknownInternalServerError.GetStatusCode())
		return
	}

	// is user active
	if user.InactiveStatus {
		http.Error(w, utils.AccountInactiveError.Error(), utils.AccountInactiveError.GetStatusCode())
		return
	}

	// verify password
	_, err = utils.VerifyPassword(user.Password, req.Password)
	if err != nil {
		if appErr, ok := err.(*utils.AppErrors); ok {
			http.Error(w, appErr.Error(), appErr.GetStatusCode())
			return
		}
		http.Error(w, "unknown internal server error", http.StatusInternalServerError)
		return
	}

	// generate token
	token, err := utils.SignToken(user.ID, req.Username, user.Role)
	if err != nil {
		if appErr, ok := err.(*utils.AppErrors); ok {
			http.Error(w, appErr.Error(), appErr.GetStatusCode())
			return
		}
	}

	// send token as a response or as a cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "Bearer",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		Expires:  time.Now().Add(24 * time.Hour),
		SameSite: http.SameSiteStrictMode,
	})

	w.Header().Set("Content-Type", "application/json")
	response := struct {
		Token string `json:"token"`
	}{
		Token: token,
	}
	json.NewEncoder(w).Encode(response)
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "Bearer",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		Expires:  time.Unix(0, 0),
		SameSite: http.SameSiteStrictMode,
	})
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"message": "Logged out successfully"}`))
}

func UpdatePasswordHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	userId, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, utils.InvalidIdError.Error(), utils.InvalidIdError.GetStatusCode())
		return
	}
	var request models.UpdatePasswordRequest
	err = json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, utils.InvalidRequestBodyError.Error(), utils.InvalidRequestBodyError.GetStatusCode())
		return
	}
	r.Body.Close()

	err = utils.ValidateExecPasswordUpdate(request)
	if err != nil {
		if appErr, ok := err.(*utils.AppErrors); ok {
			http.Error(w, appErr.Error(), appErr.GetStatusCode())
			return
		}
		http.Error(w, utils.UnknownInternalServerError.Error(), utils.UnknownInternalServerError.GetStatusCode())
		return
	}

	token, err := sqlconnect.UpdatePasswordInDB(userId, request.NewPassword, request.CurrentPassword)
	if err != nil {
		if appErr, ok := err.(*utils.AppErrors); ok {
			http.Error(w, appErr.Error(), appErr.GetStatusCode())
			return
		}
		http.Error(w, utils.UnknownInternalServerError.Error(), utils.UnknownInternalServerError.GetStatusCode())
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "Bearer",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		Expires:  time.Now().Add(24 * time.Hour),
		SameSite: http.SameSiteStrictMode,
	})

	w.Header().Set("Content-Type", "application/json")
	response := struct {
		Message string `json:"message"`
	}{
		Message: "Password updated successfully",
	}
	json.NewEncoder(w).Encode(response)
}
