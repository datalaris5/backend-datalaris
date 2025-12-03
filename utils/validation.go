package utils

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgconn"
)

func BindJSON[T any](c *gin.Context) (T, error) {
	var input T
	if err := c.ShouldBindJSON(&input); err != nil {
		return input, err
	}
	return input, nil
}

func ParseDBError(err error) string {
	if err == nil {
		return ""
	}

	// pgx
	if pgErr, ok := err.(*pgconn.PgError); ok {
		return fmt.Sprintf("Postgres error %s: %s (detail: %s, hint: %s)",
			pgErr.Code, pgErr.Message, pgErr.Detail, pgErr.Hint)
	}

	return err.Error()
}
