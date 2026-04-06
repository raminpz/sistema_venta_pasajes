package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"sistema_venta_pasajes/internal/auth/input"
	"sistema_venta_pasajes/internal/auth/service"
	"sistema_venta_pasajes/internal/auth/util"
	"sistema_venta_pasajes/pkg"
)

// AuthHandler expone los endpoints de autenticación.
type AuthHandler struct {
	svc service.AuthService
}

// NewAuthHandler crea un AuthHandler.
func NewAuthHandler(svc service.AuthService) *AuthHandler {
	return &AuthHandler{svc: svc}
}

// Login autentica al usuario y devuelve el par de tokens JWT.
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var in input.LoginInput
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		pkg.WriteError(w, r, pkg.BadRequest(util.ERR_CODE_BODY_INVALIDO, util.MSG_BODY_REQUERIDO))
		return
	}

	out, err := h.svc.Login(context.Background(), in)
	if err != nil {
		pkg.WriteError(w, r, err)
		return
	}
	pkg.WriteSuccess(w, http.StatusOK, util.MSG_LOGIN_OK, out, nil)
}

// Refresh renueva el access token usando el refresh token.
func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	var in input.RefreshInput
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		pkg.WriteError(w, r, pkg.BadRequest(util.ERR_CODE_BODY_INVALIDO, util.MSG_BODY_REQUERIDO))
		return
	}

	out, err := h.svc.Refresh(context.Background(), in)
	if err != nil {
		pkg.WriteError(w, r, err)
		return
	}
	pkg.WriteSuccess(w, http.StatusOK, util.MSG_REFRESH_OK, out, nil)
}

// Logout revoca el refresh token del usuario.
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	var in input.RefreshInput
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		pkg.WriteError(w, r, pkg.BadRequest(util.ERR_CODE_BODY_INVALIDO, util.MSG_BODY_REQUERIDO))
		return
	}

	if err := h.svc.Logout(context.Background(), in); err != nil {
		pkg.WriteError(w, r, err)
		return
	}
	pkg.WriteSuccess(w, http.StatusOK, util.MSG_LOGOUT_OK, nil, nil)
}
