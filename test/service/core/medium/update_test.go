package medium

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/factly/dega-server/service"
	"github.com/factly/dega-server/test"
	"github.com/gavv/httpexpect/v2"
	"gopkg.in/h2non/gock.v1"
)

var updatedMedium = map[string]interface{}{
	"name":        "Image",
	"type":        "jpg",
	"title":       "Sample image",
	"description": "desc",
	"caption":     "sample",
	"alt_text":    "sample",
	"file_size":   100,
	"url":         nilJsonb(),
	"dimensions":  "testdims",
}

func TestMediumUpdate(t *testing.T) {
	mock := test.SetupMockDB()

	testServer := httptest.NewServer(service.RegisterRoutes())
	gock.New(testServer.URL).EnableNetworking().Persist()
	defer gock.DisableNetworking()
	defer testServer.Close()

	// create httpexpect instance
	e := httpexpect.New(t, testServer.URL)

	t.Run("invalid medium id", func(t *testing.T) {
		test.CheckSpaceMock(mock)
		e.PUT(path).
			WithPath("medium_id", "invalid_id").
			WithHeaders(headers).
			Expect().
			Status(http.StatusNotFound)
	})

	t.Run("medium record not found", func(t *testing.T) {
		test.CheckSpaceMock(mock)
		recordNotFoundMock(mock)

		e.PUT(path).
			WithPath("medium_id", "100").
			WithHeaders(headers).
			WithJSON(updatedMedium).
			Expect().
			Status(http.StatusNotFound)
	})
	t.Run("Unable to decode medium data", func(t *testing.T) {

		test.CheckSpaceMock(mock)

		e.PUT(path).
			WithPath("medium_id", 1).
			WithHeaders(headers).
			Expect().
			Status(http.StatusUnprocessableEntity)
		test.ExpectationsMet(t, mock)

	})

	t.Run("Unprocessable medium", func(t *testing.T) {

		test.CheckSpaceMock(mock)

		e.PUT(path).
			WithPath("medium_id", 1).
			WithHeaders(headers).
			WithJSON(invalidData).
			Expect().
			Status(http.StatusUnprocessableEntity)
		test.ExpectationsMet(t, mock)

	})

	t.Run("update medium", func(t *testing.T) {
		updatedMedium["slug"] = "image"
		test.CheckSpaceMock(mock)

		SelectWithSpace(mock)

		mediumUpdateMock(mock, updatedMedium, nil)

		SelectWithOutSpace(mock)

		e.PUT(path).
			WithPath("medium_id", 1).
			WithHeaders(headers).
			WithJSON(updatedMedium).
			Expect().
			Status(http.StatusOK).JSON().Object().ContainsMap(updatedMedium)
		test.ExpectationsMet(t, mock)
	})

	t.Run("update medium by id with empty slug", func(t *testing.T) {
		test.CheckSpaceMock(mock)
		updatedMedium["slug"] = "image"
		SelectWithSpace(mock)

		slugCheckMock(mock, Data)

		mediumUpdateMock(mock, updatedMedium, nil)

		SelectWithOutSpace(mock)

		e.PUT(path).
			WithPath("medium_id", 1).
			WithHeaders(headers).
			WithJSON(dataWithoutSlug).
			Expect().
			Status(http.StatusOK).JSON().Object().ContainsMap(updatedMedium)
		test.ExpectationsMet(t, mock)

	})

	t.Run("update medium with different slug", func(t *testing.T) {
		test.CheckSpaceMock(mock)
		updatedMedium["slug"] = "image-test"

		SelectWithSpace(mock)

		mock.ExpectQuery(`SELECT slug, space_id FROM "media"`).
			WithArgs(fmt.Sprint(updatedMedium["slug"], "%"), 1).
			WillReturnRows(sqlmock.NewRows([]string{"slug", "space_id"}))

		t.Log(updatedMedium)
		mediumUpdateMock(mock, updatedMedium, nil)
		mock.ExpectQuery(selectQuery).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows(columns).
				AddRow(1, time.Now(), time.Now(), nil, updatedMedium["name"], updatedMedium["slug"], updatedMedium["type"], updatedMedium["title"], updatedMedium["description"], updatedMedium["caption"], updatedMedium["alt_text"], updatedMedium["file_size"], updatedMedium["url"], updatedMedium["dimensions"], 1))

		e.PUT(path).
			WithPath("medium_id", 1).
			WithHeaders(headers).
			WithJSON(updatedMedium).
			Expect().
			Status(http.StatusOK).JSON().Object().ContainsMap(updatedMedium)
		test.ExpectationsMet(t, mock)

	})

	t.Run("medium not found", func(t *testing.T) {
		test.CheckSpaceMock(mock)
		updatedMedium["slug"] = "toi-test"

		SelectWithSpace(mock)

		mock.ExpectQuery(`SELECT slug, space_id FROM "media"`).
			WithArgs(fmt.Sprint(updatedMedium["slug"], "%"), 1).
			WillReturnRows(sqlmock.NewRows([]string{"slug", "space_id"}))

		mediumUpdateMock(mock, updatedMedium, errors.New("update failed"))

		e.PUT(path).
			WithPath("medium_id", 1).
			WithHeaders(headers).
			WithJSON(updatedMedium).
			Expect().
			Status(http.StatusInternalServerError)
		test.ExpectationsMet(t, mock)

	})

}