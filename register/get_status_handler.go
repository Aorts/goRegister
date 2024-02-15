package register_handler

import (
	"database/sql"
	"errors"
	"goEx/api"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type GetStatusResult struct {
	Status string `db:"status"`
}

func GetStatusHandler(logger *zap.Logger, getStatusFunc GetStatusFunc) fiber.Handler {
	return func(c *fiber.Ctx) error {
		citizenId := c.Params("cid")
		logger.Info(citizenId)
		statusRes, err := getStatusFunc(citizenId)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return c.Status(http.StatusNotFound).JSON(api.Err(404, "error has occurred. please contact your system administrator"))
			}
			return c.Status(http.StatusInternalServerError).JSON(api.Err(500, "error has occurred. please contact your system administrator"))
		}
		return c.Status(http.StatusOK).JSON(api.StatusSuccess(200, "success", statusRes))
	}
}

type GetStatusFunc func(citizenId string) (string, error)

func NewGetStatusFunc() GetStatusFunc {
	return func(citizenId string) (string, error) {
		return citizenId, nil
	}
}
