package handlers

import (
	"errors"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/personal-project/pitch-league/models"
	"github.com/personal-project/pitch-league/repository"
)

type MatchHandler struct {
	BaseHandler[models.Match]
	matchRepository repository.IMatchRepository
}

func NewMatchHandler(r repository.IMatchRepository) *MatchHandler {
	return &MatchHandler{
		BaseHandler: BaseHandler[models.Match]{
			baseRepository: r,
		},
		matchRepository: r,
	}
}

func (h *MatchHandler) CreateMatch(ctx *fiber.Ctx) error {
	var vm models.MatchCreateVM
	if err := ctx.BodyParser(&vm); err != nil {
		return errorResult(ctx, err)
	}

	match := vm.ToDBModel(models.Match{})

	if err := h.matchRepository.CreateMatch(ctx.Context(), match); err != nil {
		return errorResult(ctx, errors.New("Maç oluşturulurken bir hata oluştu"))
	}

	if err := h.matchRepository.UpdateLeagueStandings(ctx.Context(), match); err != nil {
		return errorResult(ctx, errors.New("Lig sıralaması güncellenirken bir hata oluştu"))
	}

	return successResult(ctx, "Maç bilgileri başarıyla eklendi!")
}

func (h *MatchHandler) GetAllMatches(ctx *fiber.Ctx) error {
	matches, err := h.matchRepository.GetAllMatch(ctx.Context())
	if err != nil {
		return errorResult(ctx, errors.New("Maçlar getirilirken bir hata oluştu"))
	}

	var matchDetailVMs []models.MatchDetailVM
	for _, match := range matches {
		vm := models.MatchDetailVM{}.FromDBModel(match)
		matchDetailVMs = append(matchDetailVMs, vm)
	}

	return successResult(ctx, matchDetailVMs)
}

func (h *MatchHandler) GetByMatchID(ctx *fiber.Ctx) error {
	id, err := strconv.ParseInt(ctx.Params("id"), 10, 64)
	if err != nil {
		return errorResult(ctx, errors.New("Geçersiz match id"))
	}

	match, err := h.matchRepository.GetByMatchID(ctx.Context(), id)
	if err != nil {
		return errorResult(ctx, errors.New("Maç getirilirken hata oluştu"))
	}

	vm := models.MatchDetailVM{}.FromDBModel(*match)
	return successResult(ctx, vm)
}

func (h *MatchHandler) DeleteByMatchID(ctx *fiber.Ctx) error {
	id, err := strconv.ParseInt(ctx.Params("id"), 10, 64)
	if err != nil {
		return errorResult(ctx, errors.New("Geçersiz match id"))
	}

	err = h.matchRepository.DeleteByMatchID(ctx.Context(), id)
	if err != nil {
		return errorResult(ctx, errors.New("Maç silinirken hata oluştu"))
	}

	return successResult(ctx, "Maç bilgileri başarıyla silindi!")
}
