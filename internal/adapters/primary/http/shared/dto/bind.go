package dto

import (
	"strconv"

	"github.com/gofiber/fiber/v3"
)

func QueryInt(c fiber.Ctx, key string, def int) int {
	v := c.Query(key)
	if v == "" {
		return def
	}
	i, err := strconv.Atoi(v)
	if err != nil {
		return def
	}
	return i
}
