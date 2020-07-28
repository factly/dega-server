package space

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/factly/dega-server/service/core/action/policy"
	"github.com/factly/dega-server/util"
	"github.com/factly/dega-server/util/test"
	"github.com/go-chi/chi"
)

func TestDelete(t *testing.T) {

	space := os.Getenv("SPACE_ID")

	user := os.Getenv("USER_ID")
	cookie := os.Getenv("COOKIE_ID")

	headers := map[string]string{
		"space":  space,
		"user":   user,
		"cookie": cookie,
	}

	r := chi.NewRouter()
	r.With(util.CheckUser, util.CheckSpace, util.GenerateOrgnaization, policy.Authorizer).Group(func(r chi.Router) {
		r.Delete("/core/spaces/{space_id}", delete)
	})

	ts := httptest.NewServer(r)
	defer ts.Close()

	url := fmt.Sprint("/core/spaces/" + space)

	t.Run("Delete Successful", func(t *testing.T) {
		_, _, status := test.Request(t, ts, "DELETE", url, nil, headers)

		if status != http.StatusCreated {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusCreated)
		}
	})
}
