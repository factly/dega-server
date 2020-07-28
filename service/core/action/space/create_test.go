package space

import (
	"bytes"
	"encoding/json"
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

type SpaceRequest struct {
	Org  int    `json:"organisation_id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
}

func TestCreate(t *testing.T) {

	cookie := os.Getenv("COOKIE_ID")
	user_id := os.Getenv("USER_ID")
	org, _ := strconv.Atoi(os.Getenv("ORG_ID"))

	body := &SpaceRequest{
		Org:  org,
		Name: "test space",
		Slug: "test-space",
	}

	reqBody, _ := json.Marshal(body)
	req := bytes.NewReader(reqBody)

	r := chi.NewRouter()

	r.With(util.CheckUser, util.CheckSpace, util.GenerateOrgnaization, policy.Authorizer).Group(func(r chi.Router) {
		r.Post("/core/spaces", create)
	})

	ts := httptest.NewServer(r)
	defer ts.Close()

	// Successful post
	t.Run("Successful Post", func(t *testing.T) {

		headers := map[string]string{
			"space":  "0",
			"user":   user_id,
			"cookie": cookie,
		}
		_, y, status := test.Request(t, ts, "POST", "/core/spaces", req, headers)

		if status == http.StatusCreated {
			tempid := y["id"].(float64)
			orgid := y["organisation_id"].(float64)
			os.Setenv("SPACE_ID", strconv.FormatFloat(tempid, 'f', -1, 64))
			os.Setenv("ORG_ID", strconv.FormatFloat(orgid, 'f', -1, 64))
		} else {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusCreated)
		}

	})

	// Run tests
	t.Run("User Missing", func(t *testing.T) {
		headers := map[string]string{
			"space":  "0",
			"user":   "",
			"cookie": cookie,
		}
		_, _, status := test.Request(t, ts, "POST", "/core/spaces", req, headers)

		if status != http.StatusUnauthorized {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusUnauthorized)
		}
	})

}
