package http

import (
	"genealogy-tree/internal/core/domain"
	"genealogy-tree/internal/core/transport/http/request"
	corehttp "genealogy-tree/internal/core/transport/http/response"
	"net/http"
)

// Register godoc
// @Summary Регистрация пользователя
// @Description Создание нового аккаунта пользователя
// @Tags auth
// @Accept json
// @Produce json
// @Param request body RegisterRequest true "Register тело запроса"
// @Success 201 {object} RegisterResponse "Успешная регистрация"
// @Failure 400 {object} corehttp.ErrorResponse "Bad Request"
// @Failure 409 {object} corehttp.ErrorResponse "Email already exists"
// @Failure 500 {object} corehttp.ErrorResponse "internal server error"
// @Router /auth/register [post]
func (h *AuthHTTPHandler) Register(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, maxAuthBodySize)
	var input RegisterRequest
	if err := request.DecodeJSON(r, &input); err != nil {
		corehttp.Error(w, err, "request body contains invalid JSON")
		return
	}
	user, err := h.authService.Register(r.Context(), input.Email, input.Password)
	if err != nil {
		corehttp.Error(w, err, "could not register user")
		return
	}
	corehttp.JSON(w, http.StatusCreated, RegisterResponse(domain.User(user)))
}
