package utils

import (
	"github.com/go-playground/validator/v10"
	"restapi/internal/models"
)

var validate = validator.New()

func ValidateTeacherPost(newTeachers []models.Teacher) error {
	for _, teacher := range newTeachers {
		err := validate.Struct(teacher)
		if err != nil {
			return MissingFieldsError
		}
	}
	return nil
}

func ValidateStudentPost(newStudents []models.Student) error {
	for _, student := range newStudents {
		err := validate.Struct(student)
		if err != nil {
			return MissingFieldsError
		}
	}
	return nil
}

func ValidateExecPost(newExecs []models.Exec) error {
	for _, student := range newExecs {
		err := validate.Struct(student)
		if err != nil {
			return MissingFieldsError
		}
	}
	return nil
}

func ValidateExecPasswordUpdate(data models.UpdatePasswordRequest) error {
	if data.CurrentPassword == "" || data.NewPassword == "" {
		return MissingFieldsError
	}
	return nil
}
