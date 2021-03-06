package search

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/factly/x/meilisearchx"
	"github.com/meilisearch/meilisearch-go"

	"github.com/factly/x/errorx"
	"github.com/factly/x/loggerx"
	"github.com/factly/x/middlewarex"
	"github.com/factly/x/renderx"
	"github.com/factly/x/validationx"
)

// search - Search Entities
// @Summary Global search for all entities
// @Description Global search for all entities
// @Tags Search
// @ID search-entities
// @Produce json
// @Consume json
// @Param X-User header string true "User ID"
// @Param X-Space header string true "Space ID"
// @Param Search body searchQuery false "Search"
// @Success 200
// @Router /core/search [post]
func list(w http.ResponseWriter, r *http.Request) {
	sID, err := middlewarex.GetSpace(r.Context())
	if err != nil {
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.Unauthorized()))
		return
	}

	searchQuery := &searchQuery{}
	err = json.NewDecoder(r.Body).Decode(&searchQuery)
	if err != nil {
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.DecodeError()))
		return
	}

	validationError := validationx.Check(searchQuery)
	if validationError != nil {
		loggerx.Error(errors.New("validation error"))
		errorx.Render(w, validationError)
		return
	}

	// Add spaceId filter
	var filters string = fmt.Sprint("space_id=", sID)
	if len(searchQuery.Filters) > 0 {
		filters = fmt.Sprint(searchQuery.Filters, " AND ", filters)
	}

	result, err := meilisearchx.Client.Search("dega").Search(meilisearch.SearchRequest{
		Query:        searchQuery.Query,
		Limit:        searchQuery.Limit,
		Filters:      filters,
		FacetFilters: searchQuery.FacetFilters,
	})

	if err != nil {
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.InternalServerError()))
		return
	}

	renderx.JSON(w, http.StatusOK, result.Hits)
}
