package server

import (
	"encoding/json"
	"errors"
	"github.com/gofiber/fiber/v2"
	"strings"
)

type MessageStruct struct{}

type Structure struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

var Message MessageStruct

var GuardsUrl = []string{
	"file_list",
	"file",
	"validation",
	"import",
	"modify_user",
	"validation",
	"upload_append",
	"upload",
	"search",
}

func (MessageStruct) Success(data any, ctx *fiber.Ctx) (string, error) {
	relData, err := json.Marshal(Structure{Status: true, Message: "success", Data: data})
	if err == nil {
		ctx.Send(relData)
	}
	return string(relData), err
}

func (MessageStruct) Error(errorMessage string, ctx *fiber.Ctx) (error, error) {
	relData, err := json.Marshal(Structure{Status: false, Message: errorMessage, Data: nil})
	if err == nil {
		ctx.Send(relData)
	}
	return errors.New(string(relData)), err
}

func UrlGuardsInclude(path string) bool {
	result := false
	api := strings.Split(path, "/")
	if Include(GuardsUrl, api[2]) {
		result = true
	}
	return result
}

func Include[T comparable](slice []T, target T) bool {
	for i := 0; i < len(slice); i++ {
		sItem := slice[i]
		if sItem == target {
			return true
		}
	}
	return false
}

func FasterError(ctx *fiber.Ctx, msg string) error {
	Message.Error(msg, ctx)
	return errors.New(msg)
}
