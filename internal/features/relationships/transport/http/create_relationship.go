package http

import (
	"errors"
	"genealogy-tree/internal/core/domain"
	coreerrors "genealogy-tree/internal/core/errors"
	"genealogy-tree/internal/core/transport/http/request"
	corehttp "genealogy-tree/internal/core/transport/http/response"
	relrepo "genealogy-tree/internal/features/relationships/repository"
	"net/http"
)

// CreateRelationship godoc
// @Summary Создание связи
// @Description Создание новой связи между двумя людьми
// @Tags relationships
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body CreateRelationshipRequest true "CreateRelationship тело запроса"
// @Param tree_id path string true "ULID дерева"
// @Success 201 {object} RelationshipResponse "Успешно созданная связь"
// @Failure 400 {object} corehttp.ErrorResponse "Bad Request"
// @Failure 404 {object} corehttp.ErrorResponse "Person not found"
// @Failure 409 {object} corehttp.ErrorResponse "Conflict"
// @Failure 500 {object} corehttp.ErrorResponse "internal server error"
// @Router /trees/{tree_id}/relationships [post]
func (h *RelationshipsHTTPHandler) CreateRelationship(w http.ResponseWriter, r *http.Request) {
	userID, treeID, err := request.GetTreeContext(r)
	if err != nil {
		corehttp.Error(w, err, "tree access is invalid")
		return
	}
	var requestBody CreateRelationshipRequest
	if err := request.DecodeJSON(r, &requestBody); err != nil {
		corehttp.Error(w, err, "request body contains invalid JSON")
		return
	}
	command := domain.CreateRelationshipCommand{
		Relationship: domain.CreateRelationshipInput{
			Person1ID: requestBody.Person1ID,
			Person2ID: requestBody.Person2ID,
			Type:      requestBody.Type,
			Direction: requestBody.Direction,
			Metadata:  requestBody.Metadata,
		},
		FamilyID: requestBody.FamilyID,
	}
	rel, err := h.relationshipCreator.CreateRelationship(r.Context(), userID, treeID, command)
	if err != nil {
		if errors.Is(err, coreerrors.ErrAmbiguousFamily) {
			corehttp.Error(w, err, "несколько подходящих семей; укажите family_id")
			return
		}
		if errors.Is(err, relrepo.ErrDuplicate) {
			corehttp.Error(w, err, "Такая связь уже существует")
			return
		}
		corehttp.Error(w, err, "could not create relationship")
		return
	}
	corehttp.JSON(w, http.StatusCreated, rel)
}
