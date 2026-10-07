package http

import (
	"genealogy-tree/internal/core/transport/http/request"
	"net/http"
)

func familyContext(r *http.Request) (string, string, string, error) {
	userID, treeID, err := request.GetTreeContext(r)
	if err != nil {
		return "", "", "", err
	}
	familyID, err := request.GetULIDPathValue(r, "family_id")
	return userID, treeID, familyID, err
}
