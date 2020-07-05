package medium

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/factly/dega-server/config"
	"github.com/factly/dega-server/errors"
	"github.com/factly/dega-server/service/core/model"
	"github.com/factly/dega-server/util"
	"github.com/factly/x/renderx"
	"github.com/go-chi/chi"
)

// update - Update medium by id
// @Summary Update a medium by id
// @Description Update medium by ID
// @Tags Medium
// @ID update-medium-by-id
// @Produce json
// @Consume json
// @Param X-User header string true "User ID"
// @Param medium_id path string true "Medium ID"
// @Param X-Space header string true "Space ID"
// @Param Medium body medium false "Medium"
// @Success 200 {object} model.Medium
// @Router /core/media/{medium_id} [put]
func update(w http.ResponseWriter, r *http.Request) {

	sID, err := util.GetSpace(r.Context())
	if err != nil {
		errors.Render(w, errors.Parser(errors.InternalServerError()), 500)
		return
	}

	uID, err := util.GetUser(r.Context())
	if err != nil {
		errors.Parser(w, errors.InternalServerError, 500)
		return
	}

	mediumID := chi.URLParam(r, "medium_id")
	id, err := strconv.Atoi(mediumID)

	if err != nil {
		errors.Render(w, errors.Parser(errors.InvalidID()), 404)
		return
	}

	medium := &medium{}
	json.NewDecoder(r.Body).Decode(&medium)

	result := &model.Medium{}
	result.ID = uint(id)

	// check record exists or not
	err = config.DB.Where(&model.Medium{
		SpaceID: uint(sID),
	}).First(&result).Error

	if err != nil {
		errors.Render(w, errors.Parser(errors.RecordNotFound()), 404)
		return
	}

	config.DB.Model(&result).Updates(model.Medium{
		Name:        medium.Name,
		Slug:        medium.Slug,
		Title:       medium.Title,
		Type:        medium.Title,
		Description: medium.Description,
		AltText:     medium.AltText,
		Caption:     medium.Caption,
		FileSize:    medium.FileSize,
		URL:         medium.URL,
		Dimensions:  medium.Dimensions,
		Base: config.Base{
			UpdatedByID: &uID,
		},
	}).First(&result)

	renderx.JSON(w, http.StatusOK, result)
}
