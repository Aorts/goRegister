package register_handler

import (
	"database/sql"
	"goEx/api"
	"net/http"

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
			if sql.ErrNoRows == err {
				return c.Status(http.StatusNotFound).JSON(api.Err(404, "error has occurred. please contact your system administrator"))
			}
			return c.Status(http.StatusInternalServerError).JSON(api.Err(500, "error has occurred. please contact your system administrator"))
		}
		return c.Status(http.StatusOK).JSON(api.StatusSuccess(200, "success", statusRes))
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
