package space

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"

	"github.com/factly/dega-server/config"
	"github.com/factly/dega-server/service/core/model"
	"github.com/factly/dega-server/util"
	"github.com/factly/x/errorx"
	"github.com/factly/x/loggerx"
	"github.com/factly/x/renderx"
)

type organisationUser struct {
	config.Base
	Role string `gorm:"column:role" json:"role"`
}

type orgWithSpace struct {
	config.Base
	Title      string           `gorm:"column:title" json:"title"`
	Slug       string           `gorm:"column:slug;unique_index" json:"slug"`
	Permission organisationUser `json:"permission"`
	Spaces     []model.Space    `json:"spaces"`
}

// list - Get all spaces for a user
// @Summary Show all spaces
// @Description Get all spaces
// @Tags Space
// @ID get-all-spaces
// @Produce  json
// @Param X-User header string true "User ID"
// @Param X-Space header string true "Space ID"
// @Success 200 {array} orgWithSpace
// @Router /core/spaces [get]
func my(w http.ResponseWriter, r *http.Request) {
	uID, err := util.GetUser(r.Context())
	if err != nil {
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.InternalServerError()))
		return
	}

	req, err := http.NewRequest("GET", os.Getenv("KAVACH_URL")+"/organisations/my", nil)
	req.Header.Set("X-User", strconv.Itoa(uID))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.NetworkError()))
		return
	}

	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	allOrg := []orgWithSpace{}
	err = json.Unmarshal(body, &allOrg)

	var allOrgIDs []int

	for _, each := range allOrg {
		allOrgIDs = append(allOrgIDs, int(each.ID))
	}

	var allSpaces = make([]model.Space, 0)

	config.DB.Model(model.Space{}).Where("organisation_id IN (?)", allOrgIDs).Preload("Logo").Preload("LogoMobile").Preload("FavIcon").Preload("MobileIcon").Find(&allSpaces)

	result := make([]orgWithSpace, 0)

	for _, each := range allOrg {
		eachOrgWithAllSpaces := []model.Space{}
		for _, space := range allSpaces {
			if space.OrganisationID == int(each.ID) {
				eachOrgWithAllSpaces = append(eachOrgWithAllSpaces, space)
			}
		}
		each.Spaces = eachOrgWithAllSpaces
		result = append(result, each)
	}

	renderx.JSON(w, http.StatusOK, result)
}
