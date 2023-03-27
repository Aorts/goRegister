package register_handler

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func newFiberApp() *fiber.App {
	app := fiber.New()

	return app
}

func TestRegisterPass(t *testing.T) {
	newRegister := func(citizenId string, name string, birthdate string, mobile string) (string, error) {
		return "7777", nil
	}
	setRedisFunc := func(ctx context.Context, key string, value string) error {
		return nil
	}

	handler := RegisterHandler(newRegister, setRedisFunc)
	params := []byte(`
		{
			"cid": "2469903155290",
    		"name": "AOR",
    		"birthdate" : "16-02-2000",
    		"mobile": "0999999999"
		}`)
	body := bytes.NewReader(params)
	app := newFiberApp()
	app.Post("/api/verify", handler)
	req := httptest.NewRequest(http.MethodPost, "/api/verify", body)
	req.Header.Add("Content-Type", "application/json")
	res, err := app.Test(req)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusCreated, res.StatusCode)
}

func TestRegisterCIDKEY(t *testing.T) {
	newRegister := func(citizenId string, name string, birthdate string, mobile string) (string, error) {
		return "", errors.New("tbl_register_cid_key")
	}
	setRedisFunc := func(ctx context.Context, key string, value string) error {
		return nil
	}

	handler := RegisterHandler(newRegister, setRedisFunc)
	params := []byte(`
		{
			"cid": "2469903155290",
    		"name": "AOR",
    		"birthdate" : "16-02-2000",
    		"mobile": "0999999999"
		}`)
	body := bytes.NewReader(params)
	app := newFiberApp()
	app.Post("/api/register", handler)
	req := httptest.NewRequest(http.MethodPost, "/api/register", body)
	req.Header.Add("Content-Type", "application/json")
	res, err := app.Test(req)
	assert.Equal(t, err, nil)
	assert.Equal(t, http.StatusConflict, res.StatusCode)
}

func TestRegisterRedisErr(t *testing.T) {
	newRegister := func(citizenId string, name string, birthdate string, mobile string) (string, error) {
		return "", nil
	}
	setRedisFunc := func(ctx context.Context, key string, value string) error {
		return errors.New("someting")
	}

	handler := RegisterHandler(newRegister, setRedisFunc)
	params := []byte(`
		{
			"cid": "2469903155290",
    		"name": "AOR",
    		"birthdate" : "16-02-2000",
    		"mobile": "0999999999"
		}`)
	body := bytes.NewReader(params)
	app := newFiberApp()
	app.Post("/api/register", handler)
	req := httptest.NewRequest(http.MethodPost, "/api/register", body)
	req.Header.Add("Content-Type", "application/json")
	fmt.Println(req)
	res, err := app.Test(req)
	assert.Equal(t, err, nil)
	assert.Equal(t, http.StatusInternalServerError, res.StatusCode)
}

func TestRegisterUnderAge(t *testing.T) {
	newRegister := func(citizenId string, name string, birthdate string, mobile string) (string, error) {
		return "", nil
	}
	setRedisFunc := func(ctx context.Context, key string, value string) error {
		return errors.New("someting")
	}

	handler := RegisterHandler(newRegister, setRedisFunc)
	params := []byte(`
		{
			"cid": "2469903155290",
    		"name": "AOR",
    		"birthdate" : "16-02-2012",
    		"mobile": "0999999999"
		}`)
	body := bytes.NewReader(params)
	app := newFiberApp()
	app.Post("/api/register", handler)
	req := httptest.NewRequest(http.MethodPost, "/api/register", body)
	req.Header.Add("Content-Type", "application/json")
	fmt.Println(req)
	res, err := app.Test(req)
	assert.Equal(t, err, nil)
	assert.Equal(t, http.StatusForbidden, res.StatusCode)
}
