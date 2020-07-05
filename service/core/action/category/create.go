package category

import (
	"encoding/json"
	"net/http"

	"github.com/factly/dega-server/config"
	"github.com/factly/dega-server/errors"
	"github.com/factly/dega-server/service/core/model"
	"github.com/factly/dega-server/util"
	"github.com/factly/dega-server/util/slug"
	"github.com/factly/x/renderx"
	"github.com/factly/x/validationx"
)

// create - Create category
// @Summary Create category
// @Description Create category
// @Tags Category
// @ID add-category
// @Consume json
// @Produce json
// @Param X-User header string true "User ID"
// @Param X-Space header string true "Space ID"
// @Param Category body category true "Category Object"
// @Success 201 {object} model.Category
// @Failure 400 {array} string
// @Router /core/categories [post]
func create(w http.ResponseWriter, r *http.Request) {

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

	category := &category{}

	json.NewDecoder(r.Body).Decode(&category)

	validationError := validationx.Check(category)

	if validationError != nil {
		errors.Validator(w, validationError)
		return
	}

	var categorySlug string
	if category.Slug != "" && slug.Check(category.Slug) {
		categorySlug = category.Slug
	} else {
		categorySlug = slug.Make(category.Name)
	}

	result := &model.Category{
		Name:        category.Name,
		Description: category.Description,
		Slug:        slug.Approve(categorySlug, sID, config.DB.NewScope(&model.Category{}).TableName()),
		ParentID:    category.ParentID,
		MediumID:    category.MediumID,
		SpaceID:     uint(sID),
	}

	result.CreatedByID = &uID

	err = config.DB.Model(&model.Category{}).Create(&result).Error

	if err != nil {
		errors.Parser(w, err.Error(), 404)
		return
	}

	config.DB.Model(&model.Category{}).Preload("Medium").First(&result)

	renderx.JSON(w, http.StatusCreated, result)
}
