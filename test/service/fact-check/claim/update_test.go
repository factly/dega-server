package claim

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/factly/dega-server/service"
	"github.com/factly/dega-server/test"
	"github.com/factly/dega-server/test/service/fact-check/claimant"
	"github.com/factly/dega-server/test/service/fact-check/rating"
	"github.com/gavv/httpexpect/v2"
	"gopkg.in/h2non/gock.v1"
)

var updatedClaim = map[string]interface{}{
	"title":           "Claim",
	"claim_date":      time.Time{},
	"checked_date":    time.Time{},
	"claim_sources":   "GOI",
	"description":     test.NilJsonb(),
	"claimant_id":     uint(1),
	"rating_id":       uint(1),
	"review":          "Succesfully reviewed",
	"review_tag_line": "tag line",
	"review_sources":  "TOI",
}

func TestClaimUpdate(t *testing.T) {
	mock := test.SetupMockDB()

	testServer := httptest.NewServer(service.RegisterRoutes())
	gock.New(testServer.URL).EnableNetworking().Persist()
	defer gock.DisableNetworking()
	defer testServer.Close()

	// create httpexpect instance
	e := httpexpect.New(t, testServer.URL)

	t.Run("invalid claim id", func(t *testing.T) {
		test.CheckSpaceMock(mock)
		e.PUT(path).
			WithPath("claim_id", "invalid_id").
			WithHeaders(headers).
			Expect().
			Status(http.StatusNotFound)
	})

	t.Run("claim record not found", func(t *testing.T) {
		test.CheckSpaceMock(mock)
		recordNotFoundMock(mock)

		e.PUT(path).
			WithPath("claim_id", "100").
			WithHeaders(headers).
			WithJSON(updatedClaim).
			Expect().
			Status(http.StatusNotFound)
	})
	t.Run("Unable to decode claim data", func(t *testing.T) {

		test.CheckSpaceMock(mock)

		e.PUT(path).
			WithPath("claim_id", 1).
			WithHeaders(headers).
			Expect().
			Status(http.StatusUnprocessableEntity)
		test.ExpectationsMet(t, mock)

	})

	t.Run("Unprocessable claim", func(t *testing.T) {

		test.CheckSpaceMock(mock)

		e.PUT(path).
			WithPath("claim_id", 1).
			WithHeaders(headers).
			WithJSON(invalidData).
			Expect().
			Status(http.StatusUnprocessableEntity)
		test.ExpectationsMet(t, mock)

	})

	t.Run("update claim", func(t *testing.T) {
		updatedClaim["slug"] = "claim"
		test.CheckSpaceMock(mock)

		claimSelectWithSpace(mock)

		claimUpdateMock(mock, updatedClaim, nil)

		result := e.PUT(path).
			WithPath("claim_id", 1).
			WithHeaders(headers).
			WithJSON(updatedClaim).
			Expect().
			Status(http.StatusOK).JSON().Object().ContainsMap(updatedClaim)
		validateAssociations(result)
		test.ExpectationsMet(t, mock)
	})

	t.Run("update claim by id with empty slug", func(t *testing.T) {
		test.CheckSpaceMock(mock)
		updatedClaim["slug"] = "claim"
		claimSelectWithSpace(mock)

		slugCheckMock(mock, Data)

		claimUpdateMock(mock, updatedClaim, nil)

		result := e.PUT(path).
			WithPath("claim_id", 1).
			WithHeaders(headers).
			WithJSON(dataWithoutSlug).
			Expect().
			Status(http.StatusOK).JSON().Object().ContainsMap(updatedClaim)
		validateAssociations(result)
		test.ExpectationsMet(t, mock)

	})

	t.Run("update claim with different slug", func(t *testing.T) {
		test.CheckSpaceMock(mock)
		updatedClaim["slug"] = "claim-test"

		claimSelectWithSpace(mock)

		mock.ExpectQuery(`SELECT slug, space_id FROM "claims"`).
			WithArgs(fmt.Sprint(updatedClaim["slug"], "%"), 1).
			WillReturnRows(sqlmock.NewRows([]string{"slug", "space_id"}))

		claimUpdateMock(mock, updatedClaim, nil)

		result := e.PUT(path).
			WithPath("claim_id", 1).
			WithHeaders(headers).
			WithJSON(updatedClaim).
			Expect().
			Status(http.StatusOK).JSON().Object().ContainsMap(updatedClaim)
		validateAssociations(result)
		test.ExpectationsMet(t, mock)

	})
	t.Run("claimant do not belong to same space", func(t *testing.T) {
		updatedClaim["slug"] = "claim"
		test.CheckSpaceMock(mock)

		claimSelectWithSpace(mock)

		mock.ExpectBegin()
		claimant.EmptyRowMock(mock)
		mock.ExpectRollback()

		e.PUT(path).
			WithPath("claim_id", 1).
			WithHeaders(headers).
			WithJSON(updatedClaim).
			Expect().
			Status(http.StatusInternalServerError)
		test.ExpectationsMet(t, mock)

	})

	t.Run("rating do not belong to same space", func(t *testing.T) {
		updatedClaim["slug"] = "claim"
		test.CheckSpaceMock(mock)

		claimSelectWithSpace(mock)

		mock.ExpectBegin()
		claimant.SelectWithSpace(mock)
		rating.EmptyRowMock(mock)
		mock.ExpectRollback()

		e.PUT(path).
			WithPath("claim_id", 1).
			WithHeaders(headers).
			WithJSON(updatedClaim).
			Expect().
			Status(http.StatusInternalServerError)
		test.ExpectationsMet(t, mock)

	})

}