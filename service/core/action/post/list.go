package post

import (
	"fmt"
	"net/http"

	"github.com/factly/dega-server/config"
	"github.com/factly/dega-server/service/core/action/author"
	"github.com/factly/dega-server/service/core/model"
	factcheckModel "github.com/factly/dega-server/service/factcheck/model"
	"github.com/factly/dega-server/util"
	"github.com/factly/x/errorx"
	"github.com/factly/x/loggerx"
	"github.com/factly/x/paginationx"
	"github.com/factly/x/renderx"
)

// list response
type paging struct {
	Total int        `json:"total"`
	Nodes []postData `json:"nodes"`
}

// list - Get all posts
// @Summary Show all posts
// @Description Get all posts
// @Tags Post
// @ID get-all-posts
// @Produce  json
// @Param X-User header string true "User ID"
// @Param X-Space header string true "Space ID"
// @Param limit query string false "limit per page"
// @Param page query string false "page number"
// @Param format query string false "format type"
// @Success 200 {array} postData
// @Router /core/posts [get]
func list(w http.ResponseWriter, r *http.Request) {

	sID, err := util.GetSpace(r.Context())
	if err != nil {
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.InternalServerError()))
		return
	}

	tx := config.DB.Preload("Medium").Preload("Tags").Preload("Categories").Preload("Format").Model(&model.Post{}).Where(&model.Post{
		SpaceID: uint(sID),
	})

	format := r.URL.Query().Get("format")

	if format != "" {
		f := model.Format{}
		err = config.DB.Where(&model.Format{
			SpaceID: uint(sID),
			Slug:    format,
		}).First(&f).Error
		if err == nil {
			tx = tx.Where(&model.Post{
				FormatID: f.ID,
			})
		}
	}

	result := paging{}
	result.Nodes = make([]postData, 0)

	posts := make([]model.Post, 0)

	offset, limit := paginationx.Parse(r.URL.Query())

	err = tx.Count(&result.Total).Order("id desc").Offset(offset).Limit(limit).Find(&posts).Error

	if err != nil {
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.DBError()))
		return
	}

	var postIDs []uint
	for _, p := range posts {
		postIDs = append(postIDs, p.ID)
	}

	// fetch all claims related to posts
	postClaims := []factcheckModel.PostClaim{}
	config.DB.Model(&factcheckModel.PostClaim{}).Where("post_id in (?)", postIDs).Preload("Claim").Preload("Claim.Claimant").Preload("Claim.Claimant.Medium").Preload("Claim.Rating").Preload("Claim.Rating.Medium").Find(&postClaims)

	postClaimMap := make(map[uint][]factcheckModel.Claim)
	for _, pc := range postClaims {
		if _, found := postClaimMap[pc.PostID]; !found {
			postClaimMap[pc.PostID] = make([]factcheckModel.Claim, 0)
		}
		postClaimMap[pc.PostID] = append(postClaimMap[pc.PostID], pc.Claim)
	}

	// fetch all authors
	authors, err := author.All(r.Context())

	// fetch all authors related to posts
	postAuthors := []model.PostAuthor{}
	config.DB.Model(&model.PostAuthor{}).Where("post_id in (?)", postIDs).Find(&postAuthors)

	postAuthorMap := make(map[uint][]uint)
	for _, po := range postAuthors {
		if _, found := postAuthorMap[po.PostID]; !found {
			postAuthorMap[po.PostID] = make([]uint, 0)
		}
		postAuthorMap[po.PostID] = append(postAuthorMap[po.PostID], po.AuthorID)
	}

	for _, post := range posts {
		postList := &postData{}

		postList.Authors = make([]model.Author, 0)
		postList.Claims = postClaimMap[post.ID]
		postList.Post = post

		postAuthors, hasEle := postAuthorMap[post.ID]

		if hasEle {
			for _, postAuthor := range postAuthors {
				aID := fmt.Sprint(postAuthor)
				if author, found := authors[aID]; found {
					postList.Authors = append(postList.Authors, author)
				}
			}
		}

		result.Nodes = append(result.Nodes, *postList)
	}

	renderx.JSON(w, http.StatusOK, result)
}
