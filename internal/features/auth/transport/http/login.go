package http

import (
	"genealogy-tree/internal/core/transport/http/request"
	corehttp "genealogy-tree/internal/core/transport/http/response"
	"net/http"
)

// Login godoc
// @Summary Вход пользователя
// @Description Проверка email и пароля с выдачей access JWT
// @Tags auth
// @Accept json
// @Produce json
// @Param request body LoginRequest true "Login тело запроса"
// @Success 200 {object} LoginResponse "Успешный вход"
// @Failure 400 {object} corehttp.ErrorResponse "Bad Request"
// @Failure 401 {object} corehttp.ErrorResponse "Invalid credentials"
// @Failure 500 {object} corehttp.ErrorResponse "internal server error"
// @Router /auth/login [post]
func (h *AuthHTTPHandler) Login(w http.ResponseWriter, r *http.Request) {
	var input LoginRequest
	if err := request.DecodeJSON(r, &input); err != nil {
		corehttp.Error(w, err, "request body contains invalid JSON")
		return
	}
	result, err := h.authService.Login(r.Context(), input.Email, input.Password)
	if err != nil {
		corehttp.Error(w, err, "invalid email or password")
		return
	}
	corehttp.JSON(w, http.StatusOK, LoginResponse{AccessToken: result.AccessToken, TokenType: "Bearer", ExpiresIn: result.ExpiresIn})
}
