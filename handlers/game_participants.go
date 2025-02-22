package handlers

import (
	"errors"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/personal-project/pitch-league/models"
	"github.com/personal-project/pitch-league/repository"
)

type GameParticipantsHandler struct {
	BaseHandler[models.GameParticipants]
	gameParticipantsRepository repository.IGameParticipantsRepository
}

func NewGameParticipantsHandler(r repository.IGameParticipantsRepository) GameParticipantsHandler {
	return GameParticipantsHandler{
		BaseHandler: BaseHandler[models.GameParticipants]{
			baseRepository: r,
		},
		gameParticipantsRepository: r,
	}
}

func (h GameParticipantsHandler) CreateGameParticipants(ctx *fiber.Ctx) error {
	var vm models.GameParticipantsCreateVM
	if err := ctx.BodyParser(&vm); err != nil {
		return errorResult(ctx, err)
	}

	gamePart := vm.ToDBModel(models.GameParticipants{})
	if err := h.gameParticipantsRepository.CreateGameParticipants(ctx.Context(), gamePart); err != nil {
		return errorResult(ctx, errors.New("Oyuncu oyuna eklenirken bir hata oluştu"))
	}

	return successResult(ctx, "Oyuncu oyuna başarıyla eklendi!")
}

func (h GameParticipantsHandler) GetAllGameParticipants(ctx *fiber.Ctx) error {
	gameParts, err := h.gameParticipantsRepository.GetAllGameParticipants(ctx.Context())
	if err != nil {
		return errorResult(ctx, errors.New("Oyuncular getirilirken bir hata oluştu"))
	}

	var result []models.GameParticipantsDetailVM
	for _, gamePart := range gameParts {
		vm := models.GameParticipantsDetailVM{}
		result = append(result, vm.FromDBModel(gamePart))
	}

	return successResult(ctx, result)
}

func (h GameParticipantsHandler) GetByGameParticipantsID(ctx *fiber.Ctx) error {
	id, err := strconv.ParseInt(ctx.Params("id"), 10, 64)
	if err != nil {
		return errorResult(ctx, errors.New("Geçersiz kullanıcı id"))
	}

	gamePart, err := h.gameParticipantsRepository.GetByGameParticipantsID(ctx.Context(), id)
	if err != nil {
		return errorResult(ctx, errors.New("Oyuncu bilgileri getirilirken hata oluştu"))
	}

	vm := models.GameParticipantsDetailVM{}
	result := vm.FromDBModel(*gamePart)

	return successResult(ctx, result)
}

func (h GameParticipantsHandler) GetGameParticipantsUsers(ctx *fiber.Ctx) error {
	id, err := strconv.ParseInt(ctx.Params("id"), 10, 64)
	if err != nil {
		return errorResult(ctx, errors.New("Geçersiz oyun id"))
	}

	users, err := h.gameParticipantsRepository.GetGameParticipantsUsers(ctx.Context(), uint(id))
	if err != nil {
		return errorResult(ctx, errors.New("Oyuncular getirilirken hata oluştu"))
	}

	vm := models.GameParticipantsUsersVM{}
	result := vm.FromDBModel(uint(id), users)

	return successResult(ctx, result)
}

func (h GameParticipantsHandler) DeleteByGameParticipantsID(ctx *fiber.Ctx) error {
	id, err := strconv.ParseInt(ctx.Params("id"), 10, 64)
	if err != nil {
		return errorResult(ctx, errors.New("Geçersiz oyuncu id"))
	}

	err = h.gameParticipantsRepository.DeleteByGameParticipantsID(ctx.Context(), id)
	if err != nil {
		return errorResult(ctx, errors.New("Oyuncu oyundan silinirken hata oluştu"))
	}

	return successResult(ctx, "Oyuncu oyundan başarıyla silindi!")
}
