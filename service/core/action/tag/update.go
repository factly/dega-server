package tag

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/factly/dega-server/config"
	"github.com/factly/dega-server/errors"
	"github.com/factly/dega-server/service/core/model"
	"github.com/factly/dega-server/util"
	"github.com/factly/dega-server/util/slug"
	"github.com/factly/x/renderx"
	"github.com/go-chi/chi"
)

// update - Update tag by id
// @Summary Update a tag by id
// @Description Update tag by ID
// @Tags Tag
// @ID update-tag-by-id
// @Produce json
// @Consume json
// @Param X-User header string true "User ID"
// @Param tag_id path string true "Tag ID"
// @Param X-Space header string true "Space ID"
// @Param Tag body tag false "Tag"
// @Success 200 {object} model.Tag
// @Router /core/tags/{tag_id} [put]
func update(w http.ResponseWriter, r *http.Request) {
	tagID := chi.URLParam(r, "tag_id")
	id, err := strconv.Atoi(tagID)

	if err != nil {
		errors.Parser(w, errors.InvalidID, 404)
		return
	}

	sID, err := util.GetSpace(r.Context())

	if err != nil {
		errors.Parser(w, errors.InternalServerError, 500)
		return
	}

	uID, err := util.GetUser(r.Context())
	if err != nil {
		errors.Parser(w, errors.InternalServerError, 500)
		return
	}

	tag := &tag{}
	json.NewDecoder(r.Body).Decode(&tag)

	result := &model.Tag{}
	result.ID = uint(id)

	// check record exists or not
	err = config.DB.Where(&model.Tag{
		SpaceID: uint(sID),
	}).First(&result).Error

	if err != nil {
		errors.Parser(w, err.Error(), 404)
		return
	}

	var tagSlug string

	if result.Slug == tag.Slug {
		tagSlug = result.Slug
	} else if tag.Slug != "" && slug.Check(tag.Slug) {
		tagSlug = slug.Approve(tag.Slug, sID, config.DB.NewScope(&model.Tag{}).TableName())
	} else {
		tagSlug = slug.Approve(slug.Make(tag.Name), sID, config.DB.NewScope(&model.Tag{}).TableName())
	}

	config.DB.Model(&result).Updates(model.Tag{
		Name:        tag.Name,
		Slug:        tagSlug,
		Description: tag.Description,
		Base: config.Base{
			UpdatedByID: &uID,
		},
	}).First(&result)

	renderx.JSON(w, http.StatusOK, result)
}
