package handler

import (
	"bytes"
	"context"
	"testing"
)

func TestRegisterPass(t *testing.T) {
	getRedisFunc := func(ctx context.Context, key string) (int64, error) {
		return 5, nil
	}

	setRedisFunc := func(ctx context.Context, key string, value int64) error {
		return nil
	}

	params := []byte(`
		{
			"cid": "2469903155290",
    		"name": "AOR",
    		"birthdate" : "16-02-2000",
    		"mobile": "0999999999"
		}`)
	body := bytes.NewReader(params)
}
