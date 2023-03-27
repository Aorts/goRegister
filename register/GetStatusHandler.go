package handler

import (
	"goEx/api"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
)

type GetStatusResult struct {
	Status string `db:"status"`
}

func GetStatusHandler(getStatusFunc GetStatusFunc) fiber.Handler {
	return func(c *fiber.Ctx) error {
		citizenId := c.Params("cid")
		statusRes, err := getStatusFunc(citizenId)
		if err != nil {
			if strings.Contains(err.Error(), "sql: no rows in result set") {
				return c.JSON(api.Err(404, "error has occurred. please contact your system administrator"))
			} else {
				data := ReturnResponse{
					Code:    500,
					Message: "error has occurred. please contact your system administrator",
				}
				return c.JSON(data)
			}
		}
		data := ReturnResponse{
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
