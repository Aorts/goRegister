package register_handler

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetStatusPass(t *testing.T) {
	getStatusFunc := func(citizenId string) (string, error) { return "ACTIVE", nil }
	handler := GetStatusHandler(getStatusFunc)
	app := newFiberApp()
	app.Get("/api/:cid", handler)
	req := httptest.NewRequest(http.MethodGet, "/api/register", nil)
	res, err := app.Test(req)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
}

func TestGetStatusNotFound(t *testing.T) {
	getStatusFunc := func(citizenId string) (string, error) { return "", errors.New("sql: no rows in result set") }
	handler := GetStatusHandler(getStatusFunc)
	app := newFiberApp()
	app.Get("/api/:cid", handler)
	req := httptest.NewRequest(http.MethodGet, "/api/register", nil)
	res, err := app.Test(req)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusNotFound, res.StatusCode)
}

func TestGetStatusNotConnectDB(t *testing.T) {
	getStatusFunc := func(citizenId string) (string, error) { return "", errors.New("Something") }
	handler := GetStatusHandler(getStatusFunc)
	app := newFiberApp()
	app.Get("/api/:cid", handler)
	req := httptest.NewRequest(http.MethodGet, "/api/register", nil)
	res, err := app.Test(req)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusInternalServerError, res.StatusCode)
}
