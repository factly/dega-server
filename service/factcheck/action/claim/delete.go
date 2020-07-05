package claim

import (
	"net/http"
	"strconv"

	"github.com/factly/dega-server/config"
	"github.com/factly/dega-server/errors"
	"github.com/factly/dega-server/service/factcheck/model"
	"github.com/factly/dega-server/util"
	"github.com/factly/x/renderx"
	"github.com/go-chi/chi"
)

// delete - Delete claim by id
// @Summary Delete a claim
// @Description Delete claim by ID
// @Tags Claim
// @ID delete-claim-by-id
// @Param X-User header string true "User ID"
// @Param X-Space header string true "Space ID"
// @Param claim_id path string true "Claim ID"
// @Success 200
// @Router /factcheck/claims/{claim_id} [delete]
func delete(w http.ResponseWriter, r *http.Request) {

	sID, err := util.GetSpace(r.Context())
	if err != nil {
		errors.Render(w, errors.Parser(errors.InternalServerError()), 500)
		return
	}

	claimID := chi.URLParam(r, "claim_id")
	id, err := strconv.Atoi(claimID)

	if err != nil {
		errors.Render(w, errors.Parser(errors.InvalidID()), 404)
		return
	}

	result := &model.Claim{}

	result.ID = uint(id)

	// check record exists or not
	err = config.DB.Where(&model.Claim{
		SpaceID: uint(sID),
	}).First(&result).Error

	if err != nil {
		errors.Render(w, errors.Parser(errors.DBError()), 404)
		return
	}

	config.DB.Delete(&result)

	renderx.JSON(w, http.StatusOK, nil)
}
