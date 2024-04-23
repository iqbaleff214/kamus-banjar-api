package dictionary

import (
	"strings"

	"github.com/gofiber/fiber/v2"
)

// Handler contains method for fiber route handlers
type Handler interface {
	GetAlphabets(c *fiber.Ctx) error
	GetWordsByAlphabet(c *fiber.Ctx) error
	GetWord(c *fiber.Ctx) error
}

// handler as a class
type handler struct {
	service Service
}

// NewHandler is a function to instantiate new handler object
func NewHandler(service Service) Handler {
	return handler{service}
}

// GET /api/v1/alphabets
func (h handler) GetAlphabets(c *fiber.Ctx) error {
	alphabets, err := h.service.GetAlphabets()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(map[string]any{
		"code":    fiber.StatusOK,
		"message": "All alphabets successfully retrieved.",
		"status":  "success",
		"data":    alphabets,
	})
}

// GET /api/v1/alphabets/:letter
func (h handler) GetWordsByAlphabet(c *fiber.Ctx) error {
	letter := strings.ToLower(c.Params("letter"))

	alphabet, words, err := h.service.GetWordsByAlphabet(letter)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(map[string]any{
		"code":    fiber.StatusOK,
		"message": "All words with letter '" + letter + "' successfully retrieved.",
		"status":  "success",
		"data": map[string]any{
			"letter": alphabet.Letter,
			"total":  alphabet.Total,
			"words":  words,
		},
	})
}

// GET /api/v1/entries/:word
func (h handler) GetWord(c *fiber.Ctx) error {
	word := strings.ToLower(c.Params("word"))

	result, err := h.service.GetWord(word)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(map[string]any{
		"code":    fiber.StatusOK,
		"message": "Definition of word '" + result.Word + "' successfully retrieved.",
		"status":  "success",
		"data":    result,
	})
}
