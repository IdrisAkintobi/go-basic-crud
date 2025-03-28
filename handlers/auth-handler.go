package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/IdrisAkintobi/go-basic-crud/handlers/dto"
	"github.com/IdrisAkintobi/go-basic-crud/middlewares"
	"github.com/IdrisAkintobi/go-basic-crud/services"
	"github.com/IdrisAkintobi/go-basic-crud/utils"
	"github.com/jackc/pgx/v5"
)

type AuthHandler struct {
	as *services.AuthService
}

func NewAuthHandler(pg *pgx.Conn) *AuthHandler {
	return &AuthHandler{as: services.NewAuthService(pg)}
}

func (ah *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var loginData dto.AuthLoginReqDTO

	err := json.NewDecoder(r.Body).Decode(&loginData)
	if err != nil {
		utils.SendErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	reqFingerprint := r.Context().Value(utils.FPCtxKey).(middlewares.UserFingerprint)

	loginData.IPAddress = reqFingerprint.IPAddress
	loginData.UserAgent = reqFingerprint.UserAgent

	token, err := ah.as.SignIn(&loginData)
	if err != nil {
		utils.SendErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	resp := dto.AuthLoginResDTO{Token: token}

	utils.SendSuccessResponse(w, resp, http.StatusCreated)
}
