package http

import (
	"errors"
	coreerrors "genealogy-tree/internal/core/errors"
	"genealogy-tree/internal/core/security"
	corehttp "genealogy-tree/internal/core/transport/http/response"
	"genealogy-tree/internal/features/users/repository"
	"net/http"
)

// GetMe godoc
// @Summary Получение текущего пользователя
// @Description Получение профиля пользователя по идентификатору из access JWT
// @Tags users
// @Produce json
// @Success 200 {object} UserResponse "Текущий пользователь"
// @Failure 401 {object} corehttp.ErrorResponse "Unauthorized"
// @Failure 404 {object} corehttp.ErrorResponse "User not found"
// @Failure 500 {object} corehttp.ErrorResponse "internal server error"
// @Router /users/me [get]
func (h *UsersHTTPHandler) GetMe(w http.ResponseWriter, r *http.Request) {
	userID, ok := security.UserIDFromContext(r.Context())
	if !ok {
		corehttp.Error(w, coreerrors.ErrUnauthorized, "authentication is required")
		return
	}
	user, err := h.usersService.GetUser(r.Context(), userID)
	if errors.Is(err, repository.ErrNotFound) {
		corehttp.Error(w, err, "the requested user was not found")
		return
	}
	if err != nil {
		corehttp.Error(w, err, "could not get current user")
		return
	}
	corehttp.JSON(w, http.StatusOK, UserResponse{
		ID: user.ID, Email: user.Email, CreatedAt: user.CreatedAt, UpdatedAt: user.UpdatedAt,
	})
}
