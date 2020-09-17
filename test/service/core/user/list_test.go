package user

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/factly/dega-server/config"
	"github.com/factly/dega-server/service"
	"github.com/factly/dega-server/test"
	"github.com/gavv/httpexpect"
	"gopkg.in/h2non/gock.v1"
)

func TestListUsers(t *testing.T) {
	mock := test.SetupMockDB()

	test.MockServer()
	defer gock.DisableNetworking()

	testServer := httptest.NewServer(service.RegisterRoutes())
	gock.New(testServer.URL).EnableNetworking().Persist()
	defer gock.DisableNetworking()
	defer testServer.Close()

	// create httpexpect instance
	e := httpexpect.New(t, testServer.URL)

	t.Run("get users in space", func(t *testing.T) {
		test.CheckSpaceMock(mock)
		e.GET(path).
			WithHeaders(headers).
			Expect().
			Status(http.StatusOK).
			JSON().
			Object().
			Value("nodes").
			Array().
			Element(0).
			Object().
			Value("email").
			Equal(test.Dummy_AuthorList[0]["email"])
	})

	t.Run("get users when keto is down", func(t *testing.T) {
		test.DisableKetoGock(testServer.URL)
		test.CheckSpaceMock(mock)
		e.GET(path).
			WithHeaders(headers).
			Expect().
			Status(http.StatusServiceUnavailable)
	})

	t.Run("get users when keto is unable to fetch policies", func(t *testing.T) {
		test.DisableKetoGock(testServer.URL)

		gock.New(config.KetoURL + "/engines/acp/ory/regex/roles/(.+)").
			Persist().
			Reply(http.StatusOK).
			JSON(test.Dummy_Role)

		test.CheckSpaceMock(mock)

		e.GET(path).
			WithHeaders(headers).
			Expect().
			Status(http.StatusServiceUnavailable)
	})
}
