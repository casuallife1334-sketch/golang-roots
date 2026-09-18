package http

import (
	coreerrors "genealogy-tree/internal/core/errors"
	"genealogy-tree/internal/core/security"
	"genealogy-tree/internal/core/transport/http/request"
	"net/http"
)

func getTreeContext(r *http.Request) (string, string, error) {
	userID, ok := security.UserIDFromContext(r.Context())
	if !ok {
		return "", "", coreerrors.ErrUnauthorized
	}
	treeID, err := request.GetULIDPathValue(r, "tree_id")
	if err != nil {
		return "", "", err
	}
	return userID, treeID, nil
}
