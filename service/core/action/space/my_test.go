package space

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/factly/dega-server/service/core/action/policy"
	"github.com/factly/dega-server/util"
	"github.com/factly/dega-server/util/test"
	"github.com/go-chi/chi"
)

// Testing my.go
func TestMy(t *testing.T) {

	// space := os.Getenv("SPACE_ID")
	user := os.Getenv("USER_ID")
	cookie := os.Getenv("COOKIE_ID")

	headers := map[string]string{
		"space":  "",
		"user":   user,
		"cookie": cookie,
	}

	r := chi.NewRouter()

	r.With(util.CheckUser, util.CheckSpace, util.GenerateOrgnaization, policy.Authorizer).Group(func(r chi.Router) {
		r.Get("/core/spaces", my)
	})

	ts := httptest.NewServer(r)
	defer ts.Close()

	t.Run("Get Successful", func(t *testing.T) {
		x, y, status := test.Request(t, ts, "GET", "/core/spaces", nil, headers)

		if status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v and %v and %v", status, http.StatusOK, x, y)
		}
	})

}
