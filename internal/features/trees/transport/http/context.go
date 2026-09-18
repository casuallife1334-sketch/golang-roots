package http

import (
	"fmt"
	coreerrors "genealogy-tree/internal/core/errors"
	"genealogy-tree/internal/core/security"
	"genealogy-tree/internal/core/transport/http/request"
	"net/http"
)

func getUserID(r *http.Request) (string, error) {
	userID, ok := security.UserIDFromContext(r.Context())
	if !ok {
		return "", coreerrors.ErrUnauthorized
	}
	return userID, nil
}

func getTreeID(r *http.Request) (string, error) {
	treeID := r.PathValue("tree_id")
	if treeID == "" {
		return "", fmt.Errorf("tree_id is empty: %w", coreerrors.ErrInvalidArgument)
	}
	if _, err := request.GetULIDPathValue(r, "tree_id"); err != nil {
		return "", err
	}
	return treeID, nil
}
