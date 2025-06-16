package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"restapi/internal/models"
	"restapi/internal/sqlconnect"
	"restapi/utils"
	"strconv"
)

func GetOneTeacherHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	realID, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, utils.InvalidIdError.Error(), utils.InvalidIdError.GetStatusCode())
		return
	}

	teacher, err := sqlconnect.GetTeacherByID(realID)
	if err != nil {
		if appErr, ok := err.(*utils.AppErrors); ok {
			http.Error(w, appErr.Error(), appErr.GetStatusCode())
			return
		}
		http.Error(w, utils.UnknownInternalServerError.Error(), utils.UnknownInternalServerError.GetStatusCode())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(teacher)
	if err != nil {
		http.Error(w, utils.ErrorEncodingData.Error(), utils.ErrorEncodingData.GetStatusCode())
	}
}

func GetTeachersHandler(w http.ResponseWriter, r *http.Request) {
	var teacherList []models.Teacher
	teacherList, err := sqlconnect.GetTeachersDBHandler(r)
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
		Data   []models.Teacher `json:"data"`
	}{
		Status: "success",
		Count:  len(teacherList),
		Data:   teacherList,
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, utils.ErrorEncodingData.Error(), utils.ErrorEncodingData.GetStatusCode())
	}
}

func PostTeacherHandler(w http.ResponseWriter, r *http.Request) {
	var newTeachers []models.Teacher
	err := json.NewDecoder(r.Body).Decode(&newTeachers)
	if err != nil {
		http.Error(w, utils.InvalidRequestBodyError.Error(), utils.InvalidRequestBodyError.GetStatusCode())
		return
	}

	err = utils.ValidateTeacherPost(newTeachers)
	if err != nil {
		if appErr, ok := err.(*utils.AppErrors); ok {
			http.Error(w, appErr.Error(), appErr.GetStatusCode())
			return
		}
	}

	addedTeachers, err := sqlconnect.AddTeachersDBHandler(newTeachers)
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
		Data   []models.Teacher `json:"data"`
	}{
		Status: "success",
		Count:  len(addedTeachers),
		Data:   addedTeachers,
	}
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, utils.ErrorEncodingData.Error(), utils.ErrorEncodingData.GetStatusCode())
	}
}

// UpdateTeacherHandler PUT teachers - update all fields
func UpdateTeacherHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, utils.InvalidIdError.Error(), utils.InvalidIdError.GetStatusCode())
		return
	}
	var updatedTeacher models.Teacher
	err = json.NewDecoder(r.Body).Decode(&updatedTeacher)
	if err != nil {
		http.Error(w, utils.InvalidRequestBodyError.Error(), utils.InvalidRequestBodyError.GetStatusCode())
		return
	}

	updatedTeacherFromDB, err := sqlconnect.UpdateTeacherDBHandler(id, updatedTeacher)
	if err != nil {
		if appErr, ok := err.(*utils.AppErrors); ok {
			http.Error(w, appErr.Error(), appErr.GetStatusCode())
			return
		}
		http.Error(w, utils.UnknownInternalServerError.Error(), utils.UnknownInternalServerError.GetStatusCode())
		return
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(updatedTeacherFromDB)
	if err != nil {
		http.Error(w, utils.ErrorEncodingData.Error(), utils.ErrorEncodingData.GetStatusCode())
	}
}

// PatchTeachersHandler PATCH /teachers/
func PatchTeachersHandler(w http.ResponseWriter, r *http.Request) {
	var updates []map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&updates)
	if err != nil {
		http.Error(w, utils.InvalidRequestBodyError.Error(), utils.InvalidRequestBodyError.GetStatusCode())
		return
	}
	err = sqlconnect.PatchTeachersDBHandler(updates)
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
func PatchOneTeacherHandler(w http.ResponseWriter, r *http.Request) {
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

	existingTeacher, err := sqlconnect.PatchOneTeacherDBHandler(id, updates)
	if err != nil {
		if appErr, ok := err.(*utils.AppErrors); ok {
			http.Error(w, appErr.Error(), appErr.GetStatusCode())
			return
		}
		http.Error(w, utils.UnknownInternalServerError.Error(), utils.UnknownInternalServerError.GetStatusCode())
		return
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(existingTeacher)
	if err != nil {
		http.Error(w, utils.ErrorEncodingData.Error(), utils.ErrorEncodingData.GetStatusCode())
	}
}

func DeleteOneTeacherHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, utils.InvalidIdError.Error(), utils.InvalidIdError.GetStatusCode())
		return
	}
	err = sqlconnect.DeleteOneTeacherDBHandler(id)
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
		Status: "Teacher successfully deleted",
		ID:     id,
	}
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, utils.ErrorEncodingData.Error(), utils.ErrorEncodingData.GetStatusCode())
	}
}

func DeleteTeachersHandler(w http.ResponseWriter, r *http.Request) {
	var ids []int

	err := json.NewDecoder(r.Body).Decode(&ids)
	if err != nil {
		http.Error(w, utils.InvalidRequestBodyError.Error(), utils.InvalidRequestBodyError.GetStatusCode())
		return
	}

	deletedIds, err := sqlconnect.DeleteTeachersDBHandler(ids)
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
		Status:     "Teacher(s) successfully deleted",
		DeletedIDs: deletedIds,
	}
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, utils.ErrorEncodingData.Error(), utils.ErrorEncodingData.GetStatusCode())
	}
}

func GetStudentsListForTeacher(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, utils.InvalidIdError.Error(), utils.InvalidIdError.GetStatusCode())
		return
	}
	studentsList, err := sqlconnect.GetStudentsListForTeacherDBHandler(id)
	if err != nil {
		if appErr, ok := err.(*utils.AppErrors); ok {
			http.Error(w, appErr.Error(), appErr.GetStatusCode())
			return
		}
		http.Error(w, utils.UnknownInternalServerError.Error(), utils.UnknownInternalServerError.GetStatusCode())
		return
	}
	response := struct {
		Status string `json:"status"`
		Count  int    `json:"count"`
		Data   []models.Student
	}{
		Status: "success",
		Count:  len(studentsList),
		Data:   studentsList,
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, utils.ErrorEncodingData.Error(), utils.ErrorEncodingData.GetStatusCode())
		return
	}
}

func GetStudentCountForTeacher(w http.ResponseWriter, r *http.Request) {
	// allowed only for admin, manager, exec
	_, err := utils.AuthorizeUser(r.Context().Value("role").(string), "admin", "manager", "exec")
	fmt.Println(r.Context().Value("role"))
	if err != nil {
		if appErr, ok := err.(*utils.AppErrors); ok {
			http.Error(w, appErr.Error(), appErr.GetStatusCode())
			return
		}
		http.Error(w, utils.UnknownInternalServerError.Error(), utils.UnknownInternalServerError.GetStatusCode())
		return
	}
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, utils.InvalidIdError.Error(), utils.InvalidIdError.GetStatusCode())
		return
	}
	studentCount, err := sqlconnect.GetStudentCountForTeacherDBHandler(id)
	if err != nil {
		if appErr, ok := err.(*utils.AppErrors); ok {
			http.Error(w, appErr.Error(), appErr.GetStatusCode())
			return
		}
		http.Error(w, utils.UnknownInternalServerError.Error(), utils.UnknownInternalServerError.GetStatusCode())
		return
	}
	response := struct {
		Status string `json:"status"`
		Count  int    `json:"count"`
	}{
		Status: "success",
		Count:  studentCount,
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, utils.ErrorEncodingData.Error(), utils.ErrorEncodingData.GetStatusCode())
		return
	}
}
