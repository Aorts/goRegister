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

func RegisterHandler(registerFunc RegisterFunc, setRegisterRedisFunc SetRegisterRedisFunc) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var regInput RegisterInput

		err := c.BodyParser(&regInput)
		if err != nil {
			return fiber.NewError(999, "cannot Parser to body")
		}
		key := fmt.Sprintf("REGISTER:%s", regInput.CitizenId)

		Birthdate, checkAge := checkAge(regInput.Birthdate)
		if checkAge == false {
			data := ReturnResponse{
				Code:    403,
				Message: "User is underage cannot register",
			}
			return c.JSON(data)
		}

		resgfisResult, err := registerFunc(regInput.CitizenId, regInput.Name, Birthdate, regInput.Mobile)
		if err != nil {
			if strings.Contains(err.Error(), "tbl_register_cid_key") {
				data := ReturnResponse{
					Code:    409,
					Message: "User already registerd",
				}
				return c.JSON(data)
			} else {
				data := ReturnResponse{
					Code:    500,
					Message: "error has occurred. please contact your system administrator",
				}
				return c.JSON(data)
			}
		}

		err = setRegisterRedisFunc(c.Context(), key, resgfisResult)
		if err != nil {
			data := ReturnResponse{
				Code:    500,
				Message: "error has occurred. please contact your system administrator",
			}
			return c.JSON(data)
		}
		data := ReturnResponse{
			Code:    201,
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
			"PENDING",
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
