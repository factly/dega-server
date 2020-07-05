package claim

import (
	"encoding/json"
	"net/http"

	"github.com/factly/dega-server/config"
	"github.com/factly/dega-server/errors"
	"github.com/factly/dega-server/service/factcheck/model"
	"github.com/factly/dega-server/util"
	"github.com/factly/dega-server/util/slug"
	"github.com/factly/x/renderx"
	"github.com/factly/x/validationx"
)

// create - Create claim
// @Summary Create claim
// @Description Create claim
// @Tags Claim
// @ID add-claim
// @Consume json
// @Produce json
// @Param X-User header string true "User ID"
// @Param X-Space header string true "Space ID"
// @Param Claim body claim true "Claim Object"
// @Success 201 {object} model.Claim
// @Failure 400 {array} string
// @Router /factcheck/claims [post]
func create(w http.ResponseWriter, r *http.Request) {

	sID, err := util.GetSpace(r.Context())
	if err != nil {
		errors.Render(w, errors.Parser(errors.InternalServerError()), 500)
		return
	}

	uID, err := util.GetUser(r.Context())
	if err != nil {
		errors.Render(w, errors.Parser(errors.InternalServerError()), 500)
		return
	}

	claim := &claim{}

	err = json.NewDecoder(r.Body).Decode(&claim)

	if err != nil {
		errors.Render(w, errors.Parser(errors.DecodeError()), 422)
		return
	}

	validationError := validationx.Check(claim)

	if validationError != nil {
		errors.Render(w, validationError, 422)
		return
	}

	var claimSlug string
	if claim.Slug != "" && slug.Check(claim.Slug) {
		claimSlug = claim.Slug
	} else {
		claimSlug = slug.Make(claim.Title)
	}

	result := &model.Claim{
		Title:         claim.Title,
		Slug:          slug.Approve(claimSlug, sID, config.DB.NewScope(&model.Claim{}).TableName()),
		ClaimDate:     claim.ClaimDate,
		CheckedDate:   claim.CheckedDate,
		ClaimSources:  claim.ClaimSources,
		Description:   claim.Description,
		ClaimantID:    claim.ClaimantID,
		RatingID:      claim.RatingID,
		Review:        claim.Review,
		ReviewTagLine: claim.ReviewTagLine,
		ReviewSources: claim.ReviewSources,
		SpaceID:       uint(sID),
	}
	result.CreatedByID = &uID

	err = config.DB.Model(&model.Claim{}).Create(&result).Error

	if err != nil {
		return
	}

	config.DB.Model(&model.Claim{}).Preload("Rating").Preload("Claimant").Preload("Rating.Medium").Preload("Claimant.Medium").Find(&result)

	renderx.JSON(w, http.StatusCreated, result)
}
