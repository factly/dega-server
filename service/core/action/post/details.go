package post

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/factly/dega-server/config"
	"github.com/factly/dega-server/service/core/action/author"
	"github.com/factly/dega-server/service/core/model"
	factcheckModel "github.com/factly/dega-server/service/factcheck/model"
	"github.com/factly/dega-server/util"
	"github.com/factly/x/errorx"
	"github.com/factly/x/loggerx"
	"github.com/factly/x/renderx"
	"github.com/go-chi/chi"
)

// details - Get post by id
// @Summary Show a post by id
// @Description Get post by ID
// @Tags Post
// @ID get-post-by-id
// @Produce  json
// @Param X-User header string true "User ID"
// @Param X-Space header string true "Space ID"
// @Param post_id path string true "Post ID"
// @Success 200 {object} postData
// @Router /core/posts/{post_id} [get]
func details(w http.ResponseWriter, r *http.Request) {

	sID, err := util.GetSpace(r.Context())
	if err != nil {
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.InternalServerError()))
		return
	}

	postID := chi.URLParam(r, "post_id")
	id, err := strconv.Atoi(postID)

	if err != nil {
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.InvalidID()))
		return
	}

	result := &postData{}
	result.Authors = make([]model.Author, 0)
	result.Claims = make([]factcheckModel.Claim, 0)

	postAuthors := []model.PostAuthor{}
	postClaims := []factcheckModel.PostClaim{}
	result.ID = uint(id)

	err = config.DB.Model(&model.Post{}).Preload("Medium").Preload("Tags").Preload("Categories").Preload("Format").Where(&model.Post{
		SpaceID: uint(sID),
	}).First(&result.Post).Error

	if err != nil {
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.RecordNotFound()))
		return
	}

	if result.Format.Slug == "factcheck" {

		config.DB.Model(&factcheckModel.PostClaim{}).Where(&factcheckModel.PostClaim{
			PostID: uint(id),
		}).Preload("Claim").Preload("Claim.Claimant").Preload("Claim.Claimant.Medium").Preload("Claim.Rating").Preload("Claim.Rating.Medium").Find(&postClaims)

		// appending all post claims
		for _, postClaim := range postClaims {
			result.Claims = append(result.Claims, postClaim.Claim)
		}
	}

	// fetch all authors
	config.DB.Model(&model.PostAuthor{}).Where(&model.PostAuthor{
		PostID: uint(id),
	}).Find(&postAuthors)

	// Adding author
	authors, err := author.All(r.Context())
	for _, postAuthor := range postAuthors {
		aID := fmt.Sprint(postAuthor.AuthorID)
		if author, found := authors[aID]; found {
			result.Authors = append(result.Authors, author)
		}
	}

	renderx.JSON(w, http.StatusOK, result)
}
