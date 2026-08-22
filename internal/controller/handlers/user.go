package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/belyaevedu/remote-code-service/internal/domain"
	"github.com/belyaevedu/remote-code-service/internal/port"
)

type UserHandlers struct {
	userSvc port.UserService
}

func NewUserHandlers(userSvc port.UserService) *UserHandlers {
	return &UserHandlers{userSvc: userSvc}
}

// body for /register and /login
type authRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type registerResponse struct {
	Message string `json:"message"`
}

type loginResponse struct {
	Token string `json:"token"`
}

// POST /register
func (h *UserHandlers) Register(w http.ResponseWriter, r *http.Request) {
	req, ok := decodeAuth(w, r)
	if !ok {
		return
	}

	if err := h.userSvc.Register(req.Username, req.Password); err != nil {
		if errors.Is(err, domain.ErrInvalidCredentials) {
			WriteJSON(w, http.StatusBadRequest, ErrorResponse{Error: err.Error()})
			return
		}
		if errors.Is(err, domain.ErrUserAlreadyExists) {
			WriteJSON(w, http.StatusConflict, ErrorResponse{Error: err.Error()})
			return
		}
		log.Printf("error raised in userhandlers' register: %v\n", err)
		WriteJSON(w, http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	WriteJSON(w, http.StatusCreated, registerResponse{Message: "user registered"})
}

// POST /login
func (h *UserHandlers) Login(w http.ResponseWriter, r *http.Request) {
	req, ok := decodeAuth(w, r)
	if !ok {
		return
	}

	token, err := h.userSvc.Login(req.Username, req.Password)
	if err != nil {
		if errors.Is(err, domain.ErrInvalidCredentials) {
			WriteJSON(w, http.StatusUnauthorized, ErrorResponse{Error: err.Error()})
			return
		}
		log.Printf("error raised in userhandlers' login: %v\n", err)
		WriteJSON(w, http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	WriteJSON(w, http.StatusOK, loginResponse{Token: token})
}

func decodeAuth(w http.ResponseWriter, r *http.Request) (*authRequest, bool) {
	var req authRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteJSON(w, http.StatusBadRequest, ErrorResponse{Error: "invalid json body"})
		return nil, false
	}
	if req.Username == "" || req.Password == "" {
		WriteJSON(w, http.StatusBadRequest, ErrorResponse{Error: "username and password are required"})
		return nil, false
	}
	return &req, true
}
