package register_handler

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetVerifyPASS(t *testing.T) {
	getRedisfunc := func(ctx context.Context, key string) (string, error) {
		return "0000", nil
	}
	delRedis := func(ctx context.Context, key string) error {
		return nil
	}

	updateStatus := func(citizenId string) error {
		return nil
	}

	handler := SetVerifyHandler(getRedisfunc, delRedis, updateStatus)
	params := []byte(`
		{
			"cid": "2469903155290",
    		"register_code": "0000"
		}`)
	body := bytes.NewReader(params)
	app := newFiberApp()
	app.Post("/api/register", handler)
	req := httptest.NewRequest(http.MethodPost, "/api/register", body)
	req.Header.Add("Content-Type", "application/json")
	res, err := app.Test(req)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
}

func TestSetVerifyNotFoundCID(t *testing.T) {
	getRedisfunc := func(ctx context.Context, key string) (string, error) {
		return "", errors.New("not fond")
	}
	delRedis := func(ctx context.Context, key string) error {
		return nil
	}

	updateStatus := func(citizenId string) error {
		return nil
	}

	handler := SetVerifyHandler(getRedisfunc, delRedis, updateStatus)
	params := []byte(`
		{
			"cid": "2469903155290",
    		"register_code": "0000"
		}`)
	body := bytes.NewReader(params)
	app := newFiberApp()
	app.Post("/api/register", handler)
	req := httptest.NewRequest(http.MethodPost, "/api/register", body)
	req.Header.Add("Content-Type", "application/json")
	res, err := app.Test(req)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusNotFound, res.StatusCode)
}

func TestSetVerifynotMatch(t *testing.T) {
	getRedisfunc := func(ctx context.Context, key string) (string, error) {
		return "0000", nil
	}
	delRedis := func(ctx context.Context, key string) error {
		return nil
	}

	updateStatus := func(citizenId string) error {
		return nil
	}

	handler := SetVerifyHandler(getRedisfunc, delRedis, updateStatus)
	params := []byte(`
		{
			"cid": "2469903155290",
    		"register_code": "1111"
		}`)
	body := bytes.NewReader(params)
	app := newFiberApp()
	app.Post("/api/register", handler)
	req := httptest.NewRequest(http.MethodPost, "/api/register", body)
	req.Header.Add("Content-Type", "application/json")
	res, err := app.Test(req)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusNotFound, res.StatusCode)
}

func TestSetVerifyDelFail(t *testing.T) {
	getRedisfunc := func(ctx context.Context, key string) (string, error) {
		return "0000", nil
	}
	delRedis := func(ctx context.Context, key string) error {
		return errors.New("delete fail")
	}

	updateStatus := func(citizenId string) error {
		return nil
	}

	handler := SetVerifyHandler(getRedisfunc, delRedis, updateStatus)
	params := []byte(`
		{
			"cid": "2469903155290",
    		"register_code": "0000"
		}`)
	body := bytes.NewReader(params)
	app := newFiberApp()
	app.Post("/api/register", handler)
	req := httptest.NewRequest(http.MethodPost, "/api/register", body)
	req.Header.Add("Content-Type", "application/json")
	res, err := app.Test(req)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusServiceUnavailable, res.StatusCode)
}

func TestSetVerifyUpdateFail(t *testing.T) {
	getRedisfunc := func(ctx context.Context, key string) (string, error) {
		return "0000", nil
	}
	delRedis := func(ctx context.Context, key string) error {
		return nil
	}

	updateStatus := func(citizenId string) error {
		return errors.New("UpdateFail")
	}

	handler := SetVerifyHandler(getRedisfunc, delRedis, updateStatus)
	params := []byte(`
		{
			"cid": "2469903155290",
    		"register_code": "0000"
		}`)
	body := bytes.NewReader(params)
	app := newFiberApp()
	app.Post("/api/register", handler)
	req := httptest.NewRequest(http.MethodPost, "/api/register", body)
	req.Header.Add("Content-Type", "application/json")
	res, err := app.Test(req)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusServiceUnavailable, res.StatusCode)
}
