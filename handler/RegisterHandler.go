package handler

import (
	"context"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
)

type RegisterInput struct {
	CitizenId string `json:"cid"`
	Name      string `json:"name"`
	Birthdate string `json:"birthdate"`
	Mobile    string `json:"mobile"`
}

type ReturnRegister struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    *DataResult `json:"data,omitempty"`
}

type DataResult struct {
	RegisterCode string `json:"register_code,omitempty"`
	Status       string `json:"status,omitempty"`
}

type GetStatusResult struct {
	Status string `db:"status"`
}

type GetVerifyInput struct {
	CitizenId    string `json:"cid"`
	RegisterCode string `json:"register_code"`
}

func RegisterHandler(registerFunc RegisterFunc, setRegisterRedisFunc SetRegisterRedisFunc) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var regInput RegisterInput

		err := c.BodyParser(&regInput)
		if err != nil {
			return fiber.NewError(999, "cannot Parser to body")
		}
		key := fmt.Sprintf("REGISTER:%v", regInput.CitizenId)

		Birthdate, checkAge := checkAge(regInput.Birthdate)
		if checkAge == false {
			data := ReturnRegister{
				Code:    200,
				Message: "User is underage cannot register",
			}
			return c.JSON(data)
		}

		resgfisResult, err := registerFunc(regInput.CitizenId, regInput.Name, Birthdate, regInput.Mobile)
		if err != nil {
			if strings.Contains(err.Error(), "tbl_register_cid_key") {
				data := ReturnRegister{
					Code:    200,
					Message: "User already registerd",
				}
				return c.JSON(data)
			} else {
				data := ReturnRegister{
					Code:    200,
					Message: "error has occurred. please contact your system administrator",
				}
				return c.JSON(data)
			}
		}

		err = setRegisterRedisFunc(c.Context(), key, resgfisResult)
		if err != nil {
			data := ReturnRegister{
				Code:    200,
				Message: "error has occurred. please contact your system administrator",
			}
			return c.JSON(data)
		}
		data := ReturnRegister{
			Code:    200,
			Message: "success",
			Data: &DataResult{
				RegisterCode: resgfisResult,
			},
		}
		return c.JSON(data)
	}
}

type RegisterFunc func(citizenId string, name string, birthdate string, mobile string) (string, error)

type SetRegisterRedisFunc func(ctx context.Context, key string, value string) error

func NewRegisterFunc(db *sqlx.DB) RegisterFunc {
	return func(citizenId string, name string, birthdate string, mobile string) (string, error) {
		randNumStr := getRegisterCode()
		query := "insert into tbl_register (cid, name, birthdate, mobile, status, register_code) values ($1, $2, $3, $4, $5, $6)"

		tx, err := db.Begin()
		if err != nil {
			return "", err
		}
		res, err := tx.Exec(
			query,
			citizenId,
			name,
			birthdate,
			mobile,
			"Pending",
			randNumStr,
		)
		if err != nil {
			return "", err
		}
		err = tx.Commit()
		if err != nil {
			return "Commit Failure", err
		}
		_ = res
		return randNumStr, nil
	}
}

func getRegisterCode() string {
	// random number for register_code
	rand.Seed(time.Now().UnixNano())
	randNum := rand.Intn(9000) + 1000
	randNumStr := strconv.Itoa(randNum)
	return randNumStr
}

func NewRegisterRedisFunc(redisClient *redis.Client) SetRegisterRedisFunc {
	return func(ctx context.Context, key string, value string) error {
		return redisClient.Set(ctx, key, value, 0).Err()
	}
}

func checkAge(birthdate string) (string, bool) {
	birthDate, err := time.Parse("02-01-2006", birthdate)
	if err != nil {
		return "date format is not right", false
	}
	dt := birthDate.Format("02012006")
	today := time.Now()
	ages := today.Sub(birthDate).Hours() / 24 / 365

	if ages > 15 {
		return dt, true
	}
	return dt, false
}

func GetStatusHandler(getStatusFunc GetStatusFunc) fiber.Handler {
	return func(c *fiber.Ctx) error {
		citizenId := c.Params("cid")
		statusRes, err := getStatusFunc(citizenId)
		if err != nil {
			if strings.Contains(err.Error(), "sql: no rows in result set") {
				data := ReturnRegister{
					Code:    200,
					Message: "Invilid Citizen ID",
				}
				return c.JSON(data)
			} else {
				data := ReturnRegister{
					Code:    200,
					Message: "error has occurred. please contact your system administrator",
				}
				return c.JSON(data)
			}
		}
		data := ReturnRegister{
			Code:    200,
			Message: "success",
			Data: &DataResult{
				Status: statusRes,
			},
		}
		return c.JSON(data)
	}
}

type GetStatusFunc func(citizenId string) (string, error)

func NewGetStatusFunc(db *sqlx.DB) GetStatusFunc {
	return func(citizenId string) (string, error) {
		query := "select status  from  tbl_register where  cid=$1"
		res := GetStatusResult{}
		err := db.Get(&res, query, citizenId)
		if err != nil {
			return "", err
		}
		status := res.Status
		return status, nil
	}
}

func SetVerifyHandler(getRedisfunc SetVerifyFunc, delRedis DelVerifyFunc, updateStatus UpdateVerifyFunc) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var veriInput GetVerifyInput

		err := c.BodyParser(&veriInput)
		if err != nil {
			return fiber.NewError(999, "cannot Parser to body")
		}
		key := fmt.Sprintf("REGISTER:%v", veriInput.CitizenId)
		result, err := getRedisfunc(c.Context(), key)
		if err != nil {
			data := ReturnRegister{
				Code:    200,
				Message: "Invalid Citizen ID",
			}
			return c.JSON(data)
		}
		if result != veriInput.RegisterCode {
			data := ReturnRegister{
				Code:    200,
				Message: "Invalid Register Code",
			}
			return c.JSON(data)
		}
		err = delRedis(c.Context(), key)
		if err != nil {
			data := ReturnRegister{
				Code:    200,
				Message: "error has occurred. please contact your system administrator",
			}
			return c.JSON(data)
		}
		err = updateStatus(veriInput.CitizenId)
		if err != nil {
			data := ReturnRegister{
				Code:    200,
				Message: "error has occurred. please contact your system administrator",
			}
			return c.JSON(data)
		}
		data := ReturnRegister{
			Code:    200,
			Message: "success",
		}
		return c.JSON(data)
	}
}

type SetVerifyFunc func(ctx context.Context, key string) (string, error)

type DelVerifyFunc func(ctx context.Context, key string) error

type UpdateVerifyFunc func(citizenId string) error

func NewSetVerifyFunc(redisClient *redis.Client) SetVerifyFunc {
	return func(ctx context.Context, key string) (string, error) {
		return redisClient.Get(ctx, key).Result()
	}
}

func NewDelVerifyFunc(redisClient *redis.Client) DelVerifyFunc {
	return func(ctx context.Context, key string) error {
		return redisClient.Del(ctx, key).Err()
	}
}

func NewUpdateVerifyFunc(db *sqlx.DB) UpdateVerifyFunc {
	return func(citizenId string) error {
		query := "update tbl_register set status = 'ACTIVE' , updated_date  = now()  where cid=$1"
		tx, err := db.Begin()
		if err != nil {
			return err
		}
		res, err := tx.Exec(query, citizenId)
		if err != nil {
			return err
		}
		_ = res
		err = tx.Commit()
		if err != nil {
			return err
		}
		return nil
	}
}
