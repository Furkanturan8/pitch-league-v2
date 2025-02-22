package handlers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/personal-project/pitch-league/models"
	"github.com/personal-project/pitch-league/repository"
)

type GameHandler struct {
	BaseHandler[models.Game]
	gameRepository repository.IGameRepository
}

func NewGameHandler(r repository.IGameRepository) GameHandler {
	return GameHandler{
		BaseHandler: BaseHandler[models.Game]{
			baseRepository: r,
		},
		gameRepository: r,
	}
}

func (h GameHandler) CreateGame(ctx *fiber.Ctx) error {
	var vm models.GameCreateVM
	if err := ctx.BodyParser(&vm); err != nil {
		return errorResult(ctx, err)
	}

	game := vm.ToDBModel(models.Game{})
	if err := h.gameRepository.CreateGame(ctx.Context(), game); err != nil {
		return errorResult(ctx, err)
	}

	return successResult(ctx, "Game başarıyla eklendi!")
}

func (h GameHandler) GetAllGames(ctx *fiber.Ctx) error {
	games, err := h.gameRepository.GetAllGame(ctx.Context())
	if err != nil {
		return errorResult(ctx, err)
	}

	var result []models.GameDetailVM
	for _, game := range games {
		vm := models.GameDetailVM{}
		result = append(result, vm.FromDBModel(game))
	}

	return successResult(ctx, result)
}

func (h GameHandler) GetByGameID(ctx *fiber.Ctx) error {
	id, err := strconv.ParseInt(ctx.Params("id"), 10, 64)
	if err != nil {
		return errorResult(ctx, err)
	}

	game, err := h.gameRepository.GetByGameID(ctx.Context(), id)
	if err != nil {
		return errorResult(ctx, err)
	}

	vm := models.GameDetailVM{}
	result := vm.FromDBModel(*game)

	return successResult(ctx, result)
}

func (h GameHandler) DeleteByGameID(ctx *fiber.Ctx) error {
	id, err := strconv.ParseInt(ctx.Params("id"), 10, 64)
	if err != nil {
		return errorResult(ctx, err)
	}

	if err := h.gameRepository.DeleteByGameID(ctx.Context(), id); err != nil {
		return errorResult(ctx, err)
	}

	return successResult(ctx, "Game başarıyla silindi!")
}

func (h GameHandler) UpdateGameByID(ctx *fiber.Ctx) error {
	id, err := strconv.ParseInt(ctx.Params("id"), 10, 64)
	if err != nil {
		return errorResult(ctx, err)
	}

	game, err := h.gameRepository.GetByGameID(ctx.Context(), id)
	if err != nil {
		return errorResult(ctx, err)
	}

	var vm models.GameCreateVM
	if err := ctx.BodyParser(&vm); err != nil {
		return errorResult(ctx, err)
	}

	updatedGame := vm.ToDBModel(*game)
	if err := h.gameRepository.UpdateGame(ctx.Context(), updatedGame); err != nil {
		return errorResult(ctx, err)
	}

	return successResult(ctx, "Game başarıyla güncellendi!")
}
