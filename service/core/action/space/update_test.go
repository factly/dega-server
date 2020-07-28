package space

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"

	"github.com/factly/dega-server/service/core/action/policy"
	"github.com/factly/dega-server/util"
	"github.com/factly/dega-server/util/test"
	"github.com/go-chi/chi"
)

func TestUpdate(t *testing.T) {

	cookie := os.Getenv("COOKIE_ID")
	user_id := os.Getenv("USER_ID")
	org, _ := strconv.Atoi(os.Getenv("ORG_ID"))
	space_str := os.Getenv("SPACE_ID")

	send := &space{
		Name:           "testing updated",
		Slug:           "testing-updated",
		SiteTitle:      "",
		TagLine:        "",
		Description:    "",
		SiteAddress:    "",
		LogoID:         nil,
		LogoMobileID:   nil,
		OrganisationID: org,
	}

	reqBody, _ := json.Marshal(send)
	req := bytes.NewReader(reqBody)

	r := chi.NewRouter()
	r.With(util.CheckUser, util.CheckSpace, util.GenerateOrgnaization, policy.Authorizer).Group(func(r chi.Router) {
		r.Put("/core/spaces/{space_id}", update)
	})

	ts := httptest.NewServer(r)
	defer ts.Close()

	t.Run("Updated successfully", func(t *testing.T) {

		headers := map[string]string{
			"space":  space_str,
			"user":   user_id,
			"cookie": cookie,
		}

		url := fmt.Sprint("/core/spaces/" + space_str)

		_, y, status := test.Request(t, ts, "PUT", url, req, headers)

		if status == http.StatusCreated {
			tempid := y["id"].(float64)
			orgid := y["organisation_id"].(float64)
			os.Setenv("SPACE_ID", strconv.FormatFloat(tempid, 'f', -1, 64))
			os.Setenv("ORG_ID", strconv.FormatFloat(orgid, 'f', -1, 64))
		} else {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusCreated)
		}

	})

}
