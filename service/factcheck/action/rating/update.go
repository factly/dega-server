package rating

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/factly/dega-server/config"
	"github.com/factly/dega-server/errors"
	"github.com/factly/dega-server/service/factcheck/model"
	"github.com/factly/dega-server/util"
	"github.com/factly/dega-server/util/slug"
	"github.com/factly/x/renderx"
	"github.com/go-chi/chi"
)

// update - Update rating by id
// @Summary Update a rating by id
// @Description Update rating by ID
// @Tags Rating
// @ID update-rating-by-id
// @Produce json
// @Consume json
// @Param X-User header string true "User ID"
// @Param X-Space header string true "Space ID"
// @Param rating_id path string true "Rating ID"
// @Param Rating body rating false "Rating"
// @Success 200 {object} model.Rating
// @Router /factcheck/ratings/{rating_id} [put]
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

	ratingID := chi.URLParam(r, "id")
	id, err := strconv.Atoi(ratingID)

	if err != nil {
		errors.Render(w, errors.Parser(errors.InvalidID()), 404)
		return
	}

	rating := &rating{}
	json.NewDecoder(r.Body).Decode(&rating)

	result := &model.Rating{}
	result.ID = uint(id)

	// check record exists or not
	err = config.DB.Where(&model.Rating{
		SpaceID: uint(sID),
	}).First(&result).Error

	if err != nil {
		errors.Render(w, errors.Parser(errors.DBError()), 404)
		return
	}

	var ratingSlug string

	if result.Slug == rating.Slug {
		ratingSlug = result.Slug
	} else if rating.Slug != "" && slug.Check(rating.Slug) {
		ratingSlug = slug.Approve(rating.Slug, sID, config.DB.NewScope(&model.Rating{}).TableName())
	} else {
		ratingSlug = slug.Approve(slug.Make(rating.Name), sID, config.DB.NewScope(&model.Rating{}).TableName())
	}

	config.DB.Model(&result).Updates(model.Rating{
		Name:        rating.Name,
		Slug:        ratingSlug,
		MediumID:    rating.MediumID,
		Description: rating.Description,
		Base: config.Base{
			UpdatedByID: &uID,
		},
	}).Preload("Medium").First(&result)

	renderx.JSON(w, http.StatusOK, result)
}
