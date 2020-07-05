package factcheck

import (
	"fmt"
	"net/http"

	"github.com/factly/dega-server/config"
	"github.com/factly/dega-server/errors"
	"github.com/factly/dega-server/service/core/action/author"
	coreModel "github.com/factly/dega-server/service/core/model"
	"github.com/factly/dega-server/service/factcheck/model"
	"github.com/factly/dega-server/util"
	"github.com/factly/x/paginationx"
	"github.com/factly/x/renderx"
)

// list response
type paging struct {
	Total int             `json:"total"`
	Nodes []factcheckData `json:"nodes"`
}

// list - Get all factchecks
// @Summary Show all factchecks
// @Description Get all factchecks
// @Tags Factcheck
// @ID get-all-factchecks
// @Produce  json
// @Param X-User header string true "User ID"
// @Param X-Space header string true "Space ID"
// @Param limit query string false "limit per page"
// @Param page query string false "page number"
// @Success 200 {object} paging
// @Router /factcheck/factchecks [get]
func list(w http.ResponseWriter, r *http.Request) {

	sID, err := util.GetSpace(r.Context())
	if err != nil {
		errors.Render(w, errors.Parser(errors.InternalServerError()), 500)
		return
	}

	result := paging{}
	result.Nodes = make([]factcheckData, 0)

	offset, limit := paginationx.Parse(r.URL.Query())

	factchecks := []model.Factcheck{}

	err = config.DB.Model(&model.Factcheck{}).Preload("Medium").Where(&model.Factcheck{
		SpaceID: uint(sID),
	}).Count(&result.Total).Order("id desc").Offset(offset).Limit(limit).Find(&factchecks).Error

	if err != nil {
		return
	}

	// fetch all authors
	authors, err := author.All(r.Context())

	for _, factcheck := range factchecks {
		factcheckList := &factcheckData{}
		factcheckList.Categories = make([]coreModel.Category, 0)
		factcheckList.Tags = make([]coreModel.Tag, 0)
		factcheckList.Claims = make([]model.Claim, 0)
		factcheckList.Authors = make([]coreModel.Author, 0)

		categories := []model.FactcheckCategory{}
		tags := []model.FactcheckTag{}
		claims := []model.FactcheckClaim{}
		factCheckAuthors := []model.FactcheckAuthor{}

		factcheckList.Factcheck = factcheck

		// fetch all categories
		config.DB.Model(&model.FactcheckCategory{}).Where(&model.FactcheckCategory{
			FactcheckID: factcheck.ID,
		}).Preload("Category").Preload("Category.Medium").Find(&categories)

		// fetch all tags
		config.DB.Model(&model.FactcheckTag{}).Where(&model.FactcheckTag{
			FactcheckID: factcheck.ID,
		}).Preload("Tag").Find(&tags)

		// fetch all claims
		config.DB.Model(&model.FactcheckClaim{}).Where(&model.FactcheckClaim{
			FactcheckID: factcheck.ID,
		}).Preload("Claim").Preload("Claim.Claimant").Preload("Claim.Claimant.Medium").Preload("Claim.Rating").Preload("Claim.Rating.Medium").Find(&claims)

		// fetch all post authors
		config.DB.Model(&model.FactcheckAuthor{}).Where(&model.FactcheckAuthor{
			FactcheckID: factcheck.ID,
		}).Find(&factCheckAuthors)

		for _, c := range categories {
			factcheckList.Categories = append(factcheckList.Categories, c.Category)
		}

		for _, t := range tags {
			factcheckList.Tags = append(factcheckList.Tags, t.Tag)
		}

		for _, c := range claims {
			factcheckList.Claims = append(factcheckList.Claims, c.Claim)
		}

		for _, factCheckAuthor := range factCheckAuthors {
			aID := fmt.Sprint(factCheckAuthor.AuthorID)
			if authors[aID].Email != "" {
				factcheckList.Authors = append(factcheckList.Authors, authors[aID])
			}
		}

		result.Nodes = append(result.Nodes, *factcheckList)
	}

	renderx.JSON(w, http.StatusOK, result)
}
