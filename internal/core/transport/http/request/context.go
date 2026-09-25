package request

import (
	coreerrors "genealogy-tree/internal/core/errors"
	"genealogy-tree/internal/core/security"
	"net/http"
)

func GetUserID(r *http.Request) (string, error) {
	userID, ok := security.UserIDFromContext(r.Context())
	if !ok {
		return "", coreerrors.ErrUnauthorized
	}
	return userID, nil
}

func GetTreeID(r *http.Request) (string, error) {
	return GetULIDPathValue(r, "tree_id")
}

func GetTreeContext(r *http.Request) (string, string, error) {
	userID, err := GetUserID(r)
	if err != nil {
		return "", "", err
	}
	treeID, err := GetTreeID(r)
	if err != nil {
		return "", "", err
	}
	return userID, treeID, nil
}
