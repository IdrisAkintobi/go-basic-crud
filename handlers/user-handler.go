package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/IdrisAkintobi/go-basic-crud/handlers/dto"
	"github.com/IdrisAkintobi/go-basic-crud/services"
	"github.com/IdrisAkintobi/go-basic-crud/utils"
	"github.com/jackc/pgx/v5"
)

type UserHandler struct {
	us *services.UserService
}

func NewUserHandler(pg *pgx.Conn) *UserHandler {
	return &UserHandler{us: services.NewUserService(pg)}
}

func (uh *UserHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var userData dto.RegisterUserReqDTO

	err := json.NewDecoder(r.Body).Decode(&userData)
	if err != nil {
		utils.SendErrorResponse(w, fmt.Sprintf("error parsing request body: %s", err.Error()), http.StatusBadRequest)
		return
	}

	newUser, err := uh.us.RegisterUser(&userData)
	if err != nil {
		utils.SendErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	utils.SendSuccessResponse(w, newUser, http.StatusCreated)
}
