package handlers

import (
	"errors"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/personal-project/pitch-league/models"
	"github.com/personal-project/pitch-league/repository"
)

type LeagueHandler struct {
	BaseHandler[models.League]
	leagueRepository repository.ILeagueRepository
}

func NewLeagueHandler(r repository.ILeagueRepository) LeagueHandler {
	return LeagueHandler{
		BaseHandler: BaseHandler[models.League]{
			baseRepository: r,
		},
		leagueRepository: r,
	}
}

func (h LeagueHandler) CreateLeague(ctx *fiber.Ctx) error {
	var vm models.LeagueCreateVM
	if err := ctx.BodyParser(&vm); err != nil {
		return errorResult(ctx, err)
	}

	league := vm.ToDBModel(models.League{})
	if err := h.leagueRepository.CreateLeague(ctx.Context(), league); err != nil {
		return errorResult(ctx, errors.New("Lig oluşturulurken bir hata oluştu"))
	}

	return successResult(ctx, "Lig başarıyla eklendi!")
}

func (h LeagueHandler) GetAllLeagues(ctx *fiber.Ctx) error {
	leagues, err := h.leagueRepository.GetAllLeague(ctx.Context())
	if err != nil {
		return errorResult(ctx, errors.New("Ligler getirilirken bir hata oluştu"))
	}

	var result []models.LeagueDetailVM
	for _, league := range leagues {
		vm := models.LeagueDetailVM{}
		result = append(result, vm.FromDBModel(league))
	}

	return successResult(ctx, result)
}

func (h LeagueHandler) GetByLeagueID(ctx *fiber.Ctx) error {
	id, err := strconv.ParseInt(ctx.Params("id"), 10, 64)
	if err != nil {
		return errorResult(ctx, errors.New("Geçersiz lig id"))
	}

	league, err := h.leagueRepository.GetByLeagueID(ctx.Context(), id)
	if err != nil {
		return errorResult(ctx, errors.New("Lig getirilirken hata oluştu"))
	}

	vm := models.LeagueDetailVM{}
	result := vm.FromDBModel(*league)

	return successResult(ctx, result)
}

func (h LeagueHandler) DeleteByLeagueID(ctx *fiber.Ctx) error {
	id, err := strconv.ParseInt(ctx.Params("id"), 10, 64)
	if err != nil {
		return errorResult(ctx, errors.New("Geçersiz lig id"))
	}

	err = h.leagueRepository.DeleteByLeagueID(ctx.Context(), id)
	if err != nil {
		return errorResult(ctx, errors.New("Lig silinirken hata oluştu"))
	}

	return successResult(ctx, "Lig başarıyla silindi!")
}
