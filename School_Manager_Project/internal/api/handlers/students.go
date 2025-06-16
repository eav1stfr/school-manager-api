package handlers

import (
	"encoding/json"
	"net/http"
	"restapi/internal/models"
	"restapi/internal/sqlconnect"
	"restapi/utils"
	"strconv"
)

func GetOneStudentHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	realID, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, utils.InvalidIdError.Error(), utils.InvalidIdError.GetStatusCode())
		return
	}

	student, err := sqlconnect.GetStudentByID(realID)
	if err != nil {
		if appErr, ok := err.(*utils.AppErrors); ok {
			http.Error(w, appErr.Error(), appErr.GetStatusCode())
			return
		}
		http.Error(w, utils.UnknownInternalServerError.Error(), utils.UnknownInternalServerError.GetStatusCode())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(student)
	if err != nil {
		http.Error(w, utils.ErrorEncodingData.Error(), utils.ErrorEncodingData.GetStatusCode())
		return
	}
}

func GetStudentsHandler(w http.ResponseWriter, r *http.Request) {
	var studentList []models.Student
	studentList, err := sqlconnect.GetStudentsDBHandler(studentList, r)
	if err != nil {
		if appErr, ok := err.(*utils.AppErrors); ok {
			http.Error(w, appErr.Error(), appErr.GetStatusCode())
			return
		}
		http.Error(w, utils.UnknownInternalServerError.Error(), utils.UnknownInternalServerError.GetStatusCode())
		return
	}
	response := struct {
		Status string           `json:"status"`
		Count  int              `json:"count"`
		Data   []models.Student `json:"data"`
	}{
		Status: "success",
		Count:  len(studentList),
		Data:   studentList,
	}

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, utils.ErrorEncodingData.Error(), utils.ErrorEncodingData.GetStatusCode())
		return
	}
	w.Header().Set("Content-Type", "application/json")
}

func PostStudentHandler(w http.ResponseWriter, r *http.Request) {
	var newStudents []models.Student
	err := json.NewDecoder(r.Body).Decode(&newStudents)
	if err != nil {
		http.Error(w, "Error decoding data", http.StatusBadRequest)
		return
	}

	err = utils.ValidateStudentPost(newStudents)
	if err != nil {
		if appErr, ok := err.(*utils.AppErrors); ok {
			http.Error(w, appErr.Error(), appErr.GetStatusCode())
			return
		}
	}

	addedStudents, err := sqlconnect.AddStudentsDBHandler(newStudents)
	if err != nil {
		if appErr, ok := err.(*utils.AppErrors); ok {
			http.Error(w, appErr.Error(), appErr.GetStatusCode())
			return
		}
		http.Error(w, utils.UnknownInternalServerError.Error(), utils.UnknownInternalServerError.GetStatusCode())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	response := struct {
		Status string           `json:"status"`
		Count  int              `json:"count"`
		Data   []models.Student `json:"data"`
	}{
		Status: "success",
		Count:  len(addedStudents),
		Data:   addedStudents,
	}
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, utils.ErrorEncodingData.Error(), utils.ErrorEncodingData.GetStatusCode())
		return
	}
}

// UpdateTeacherHandler PUT teachers - update all fields
func UpdateStudentHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, utils.InvalidIdError.Error(), utils.InvalidIdError.GetStatusCode())
		return
	}
	var updatedStudent models.Student
	err = json.NewDecoder(r.Body).Decode(&updatedStudent)
	if err != nil {
		http.Error(w, utils.InvalidRequestBodyError.Error(), utils.InvalidRequestBodyError.GetStatusCode())
		return
	}

	updatedStudentFromDB, err := sqlconnect.UpdateStudentDBHandler(id, updatedStudent)
	if err != nil {
		if appErr, ok := err.(*utils.AppErrors); ok {
			http.Error(w, appErr.Error(), appErr.GetStatusCode())
			return
		}
		http.Error(w, utils.UnknownInternalServerError.Error(), utils.UnknownInternalServerError.GetStatusCode())
		return
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(updatedStudentFromDB)
	if err != nil {
		http.Error(w, utils.ErrorEncodingData.Error(), utils.ErrorEncodingData.GetStatusCode())
	}
}

// PatchStudentsHandler PATCH /teachers/
func PatchStudentsHandler(w http.ResponseWriter, r *http.Request) {
	var updates []map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&updates)
	if err != nil {
		http.Error(w, utils.InvalidRequestBodyError.Error(), utils.InvalidRequestBodyError.GetStatusCode())
		return
	}
	err = sqlconnect.PatchStudentsDBHandler(updates)
	if err != nil {
		if appErr, ok := err.(*utils.AppErrors); ok {
			http.Error(w, appErr.Error(), appErr.GetStatusCode())
			return
		}
		http.Error(w, utils.UnknownInternalServerError.Error(), utils.UnknownInternalServerError.GetStatusCode())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// PatchOneTeacherHandler patch method - only update received fields
func PatchOneStudentHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, utils.InvalidIdError.Error(), utils.InvalidIdError.GetStatusCode())
		return
	}

	var updates map[string]interface{}
	err = json.NewDecoder(r.Body).Decode(&updates)
	if err != nil {
		http.Error(w, utils.InvalidRequestBodyError.Error(), utils.InvalidRequestBodyError.GetStatusCode())
		return
	}

	existingStudent, err := sqlconnect.PatchOneStudentDBHandler(id, updates)
	if err != nil {
		if appErr, ok := err.(*utils.AppErrors); ok {
			http.Error(w, appErr.Error(), appErr.GetStatusCode())
			return
		}
		http.Error(w, utils.UnknownInternalServerError.Error(), utils.UnknownInternalServerError.GetStatusCode())
		return
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(existingStudent)
	if err != nil {
		http.Error(w, utils.ErrorEncodingData.Error(), utils.ErrorEncodingData.GetStatusCode())
	}
}

func DeleteOneStudentHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, utils.InvalidIdError.Error(), utils.InvalidIdError.GetStatusCode())
		return
	}
	err = sqlconnect.DeleteOneStudentDBHandler(id)
	if err != nil {
		if appErr, ok := err.(*utils.AppErrors); ok {
			http.Error(w, appErr.Error(), appErr.GetStatusCode())
			return
		}
		http.Error(w, utils.UnknownInternalServerError.Error(), utils.UnknownInternalServerError.GetStatusCode())
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
		Status: "Student successfully deleted",
		ID:     id,
	}
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, utils.ErrorEncodingData.Error(), utils.ErrorEncodingData.GetStatusCode())
	}
}

func DeleteStudentsHandler(w http.ResponseWriter, r *http.Request) {
	var ids []int

	err := json.NewDecoder(r.Body).Decode(&ids)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	deletedIds, err := sqlconnect.DeleteStudentsDBHandler(ids)
	if err != nil {
		if appErr, ok := err.(*utils.AppErrors); ok {
			http.Error(w, appErr.Error(), appErr.GetStatusCode())
			return
		}
		http.Error(w, utils.UnknownInternalServerError.Error(), utils.UnknownInternalServerError.GetStatusCode())
		return
	}
	w.Header().Set("Content-Type", "application/json")

	response := struct {
		Status     string `json:"status"`
		DeletedIDs []int  `json:"deleted_ids"`
	}{
		Status:     "Student(s) successfully deleted",
		DeletedIDs: deletedIds,
	}
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, utils.ErrorEncodingData.Error(), utils.ErrorEncodingData.GetStatusCode())
	}
}
