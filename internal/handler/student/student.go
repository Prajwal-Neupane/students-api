package student

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/Prajwal-Neupane/students-api/internal/storage"
	"github.com/Prajwal-Neupane/students-api/internal/types"
	"github.com/Prajwal-Neupane/students-api/pkg/utils"
	"github.com/go-playground/validator/v10"
)


func New(storage storage.Storage) http.HandlerFunc {

	return  func(w http.ResponseWriter, r *http.Request) {
		slog.Info("Creating a student")

		var student types.Student

		err := json.NewDecoder(r.Body).Decode(&student)

		if errors.Is(err, io.EOF) {
			utils.WriteJson(w, http.StatusBadRequest, utils.GeneralError(err))
			return 
		}

		if err != nil {
			utils.WriteJson(w, http.StatusBadRequest, utils.GeneralError(err))
		}
		
		// _, err := w.Write([]byte("Welcome to student API"))
		// if err != nil {
		// 	panic(err)
		// }
		// w.Write([] byte("Welcome to students api"))



		// Request Validation

		if err:= validator.New().Struct(student); err != nil {

			validateErrors := err.(validator.ValidationErrors)
			utils.WriteJson(w, http.StatusBadRequest, utils.ValidationError(validateErrors))
			return
		}

		lastId, err := storage.CreateStudent(
			student.Name,
			student.Email,
			student.Age,
		)
		slog.Info("user created successfully", slog.String("userId", fmt.Sprint(lastId)))

		if err != nil {
			utils.WriteJson(w, http.StatusInternalServerError, err)
			return 
		}

		


		utils.WriteJson(w, http.StatusCreated, map[string]int64{"id": lastId})
	}
}