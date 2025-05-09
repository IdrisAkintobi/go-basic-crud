package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/IdrisAkintobi/go-basic-crud/handlers/dto"
	"github.com/IdrisAkintobi/go-basic-crud/middlewares"
	"github.com/IdrisAkintobi/go-basic-crud/services"
	"github.com/IdrisAkintobi/go-basic-crud/utils"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AuthHandler struct {
	as *services.AuthService
}

func NewAuthHandler(pg *pgxpool.Pool) *AuthHandler {
	return &AuthHandler{as: services.NewAuthService(pg)}
}

func (ah *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var loginData dto.AuthLoginReqDTO

	err := json.NewDecoder(r.Body).Decode(&loginData)
	if err != nil {
		utils.SendErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	if loginData.Email == "" || loginData.Password == "" {
		utils.SendErrorResponse(w, "email and password can not be empty", http.StatusBadRequest)
		return
	}

	reqFingerprint := r.Context().Value(utils.FPCtxKey).(*middlewares.UserFingerprint)

	if reqFingerprint.DeviceId == "" {
		utils.SendErrorResponse(w, "device id must be passed in the header", http.StatusBadRequest)
		return
	}

	loginData.DeviceId = reqFingerprint.DeviceId
	loginData.IPAddress = reqFingerprint.IPAddress
	loginData.UserAgent = reqFingerprint.UserAgent

	token, err := ah.as.SignIn(&loginData)
	if err != nil {
		utils.SendErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	resp := dto.AuthLoginResDTO{Token: token}

	utils.SendSuccessResponse(w, resp, http.StatusOK)
}

func (ah *AuthHandler) LogOut(w http.ResponseWriter, r *http.Request) {
	authData := r.Context().Value(utils.AuthUserCtxKey).(*middlewares.AuthData)

	err := ah.as.LogOut(authData.SessionId)
	if err != nil {
		utils.SendErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	utils.SendSuccessResponse(w, nil, http.StatusOK)
}

func (ah *AuthHandler) WhoAmI(w http.ResponseWriter, r *http.Request) {
	authData := r.Context().Value(utils.AuthUserCtxKey).(*middlewares.AuthData)

	userData, err := ah.as.WhoAmI(authData.UserID)
	if err != nil {
		utils.SendErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}
	if userData == nil {
		utils.SendErrorResponse(w, "user not found", http.StatusNotFound)
		return
	}

	utils.SendSuccessResponse(w, userData, http.StatusOK)
}

func (ah *AuthHandler) GetActiveSessions(w http.ResponseWriter, r *http.Request) {
	authData := r.Context().Value(utils.AuthUserCtxKey).(*middlewares.AuthData)

	userSessions, err := ah.as.GetActiveSessions(authData.UserID)
	if err != nil {
		utils.SendErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	utils.SendSuccessResponse(w, userSessions, http.StatusOK)
}
