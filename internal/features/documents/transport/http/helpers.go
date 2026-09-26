package http

import (
	"fmt"
	"genealogy-tree/internal/core/domain"
	"genealogy-tree/internal/core/transport/http/request"
	"genealogy-tree/internal/core/validation"
	"net/http"
)

func documentOwner(r *http.Request) (domain.DocumentOwner, error) {
	ownerType := r.URL.Query().Get("owner_type")
	if err := validation.ValidateDocumentOwnerType(ownerType); err != nil {
		return domain.DocumentOwner{}, fmt.Errorf("owner_type is invalid: %w", err)
	}
	ownerID, err := request.GetULIDQueryValue(r, "owner_id")
	if err != nil {
		return domain.DocumentOwner{}, err
	}
	return domain.DocumentOwner{Type: domain.DocumentOwnerType(ownerType), ID: ownerID}, nil
}
