package dictionary

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/iqbaleff214/kamus-banjar-api/pkg/response"
)

// VoteCounter is an optional dependency injected into the handler to append
// vote counts to word-detail responses. Implemented by community.Service.
type VoteCounter interface {
	GetVoteSummary(word string) (up, down int, err error)
}

// Handler contains method for fiber route handlers.
type Handler interface {
	GetAlphabets(c *fiber.Ctx) error
	GetWordsByAlphabet(c *fiber.Ctx) error
	GetWord(c *fiber.Ctx) error
	Search(c *fiber.Ctx) error
}

type handler struct {
	service Service
	votes   VoteCounter // optional; nil when community features are disabled
}

// NewHandler creates a dictionary handler. Pass a non-nil VoteCounter to include
// vote counts in GET /api/v1/entries/:word responses.
func NewHandler(service Service, votes VoteCounter) Handler {
	return handler{service: service, votes: votes}
}

// GET /api/v1/alphabets
func (h handler) GetAlphabets(c *fiber.Ctx) error {
	alphabets, err := h.service.GetAlphabets()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	c.Set("Cache-Control", "public, max-age=86400")
	return response.OK(c, "All alphabets successfully retrieved.", alphabets)
}

// GET /api/v1/alphabets/:letter
func (h handler) GetWordsByAlphabet(c *fiber.Ctx) error {
	letter := strings.ToLower(c.Params("letter"))

	alphabet, words, err := h.service.GetWordsByAlphabet(letter)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	c.Set("Cache-Control", "public, max-age=86400")
	return response.OK(c, "All words with letter '"+letter+"' successfully retrieved.", map[string]any{
		"letter": alphabet.Letter,
		"total":  alphabet.Total,
		"words":  words,
	})
}

// GET /api/v1/entries/:word
func (h handler) GetWord(c *fiber.Ctx) error {
	word := strings.ToLower(c.Params("word"))

	result, err := h.service.GetWord(word)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	var votes any
	if h.votes != nil {
		up, down, _ := h.votes.GetVoteSummary(word)
		votes = map[string]int{"up": up, "down": down}
	}

	c.Set("Cache-Control", "public, max-age=604800")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"code":    fiber.StatusOK,
		"status":  "success",
		"message": "Definition of word '" + result.Word + "' successfully retrieved.",
		"data":    result,
		"votes":   votes,
	})
}

// GET /api/v1/entries?search=keyword
func (h handler) Search(c *fiber.Ctx) error {
	keyword := c.Query("search")
	data, err := h.service.Search(keyword)

	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	c.Set("Cache-Control", "public, max-age=3600")
	return response.OK(c, "Matching word retrieved successfully.", data)
}
