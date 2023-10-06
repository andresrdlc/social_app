package infrastructure

import (
	"encoding/json"
	"net/http"

	"github.com/unmsmfisi-socialapplication/social_app/internal/register/application"
	"github.com/unmsmfisi-socialapplication/social_app/pkg/utils"
)

type RegisterUserHandler struct {
	useCase *application.RegistrationUseCase
}

func NewRegisterUserHandler(uc *application.RegistrationUseCase) *RegisterUserHandler {
	return &RegisterUserHandler{useCase: uc}
}

func (rh *RegisterUserHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {

	var data struct {
		Phone     string `json:"phone"`
		Email     string `json:"email"`
		User_name string `json:"user_name"`
		Password  string `json:"password"`
	}
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {

		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	_, er := rh.useCase.RegisterUser(data.Phone, data.Email, data.User_name, data.Password)
	if er != nil {
		switch er {
		case application.ErrEmailInUse:
			utils.SendJSONResponse(w, http.StatusBadRequest, "ERROR", "Email already in use")
			return
		case application.ErrFormat:
			utils.SendJSONResponse(w, http.StatusBadRequest, "ERROR", "Not secure password")
			return
		case application.ErrPhone:
			utils.SendJSONResponse(w, http.StatusBadRequest, "ERROR", "Invalid phone format")
			return

		}
	}
	json.NewEncoder(w).Encode(data)
}
