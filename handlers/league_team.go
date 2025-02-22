package handlers

import (
	"errors"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/personal-project/pitch-league/models"
	"github.com/personal-project/pitch-league/repository"
)

type LeagueTeamHandler struct {
	BaseHandler[models.LeagueTeam]
	leagueTeamRepository repository.ILeagueTeamRepository
}

func NewLeagueTeamHandler(r repository.ILeagueTeamRepository) LeagueTeamHandler {
	return LeagueTeamHandler{
		BaseHandler: BaseHandler[models.LeagueTeam]{
			baseRepository: r,
		},
		leagueTeamRepository: r,
	}
}

func (h LeagueTeamHandler) CreateLeagueTeam(ctx *fiber.Ctx) error {
	var vm models.LeagueTeamCreateVM
	if err := ctx.BodyParser(&vm); err != nil {
		return errorResult(ctx, err)
	}

	leagueTeam := vm.ToDBModel(models.LeagueTeam{})
	if err := h.leagueTeamRepository.CreateLeagueTeam(ctx.Context(), leagueTeam); err != nil {
		return errorResult(ctx, errors.New("Takım lige eklenirken bir hata oluştu"))
	}

	return successResult(ctx, "Takım lige başarıyla eklendi!")
}

func (h LeagueTeamHandler) GetAllLeagueTeams(ctx *fiber.Ctx) error {
	leagueTeams, err := h.leagueTeamRepository.GetAllLeagueTeam(ctx.Context())
	if err != nil {
		return errorResult(ctx, errors.New("Lig takımları getirilirken bir hata oluştu"))
	}

	var result []models.LeagueTeamDetailVM
	for _, leagueTeam := range leagueTeams {
		vm := models.LeagueTeamDetailVM{}
		result = append(result, vm.FromDBModel(leagueTeam))
	}

	return successResult(ctx, result)
}

func (h LeagueTeamHandler) GetByLeagueTeamID(ctx *fiber.Ctx) error {
	id, err := strconv.ParseInt(ctx.Params("id"), 10, 64)
	if err != nil {
		return errorResult(ctx, errors.New("Geçersiz takım id"))
	}

	leagueTeam, err := h.leagueTeamRepository.GetByLeagueTeamID(ctx.Context(), id)
	if err != nil {
		return errorResult(ctx, errors.New("Takım bilgileri getirilirken hata oluştu"))
	}

	vm := models.LeagueTeamDetailVM{}
	result := vm.FromDBModel(*leagueTeam)

	return successResult(ctx, result)
}

func (h LeagueTeamHandler) GetByLeagueID(ctx *fiber.Ctx) error {
	id, err := strconv.ParseInt(ctx.Params("id"), 10, 64)
	if err != nil {
		return errorResult(ctx, errors.New("Geçersiz lig id"))
	}

	leagueTeams, err := h.leagueTeamRepository.GetByLeagueID(ctx.Context(), id)
	if err != nil {
		return errorResult(ctx, errors.New("Lig takımları getirilirken hata oluştu"))
	}

	var result []models.LeagueTeamDetailVM
	for _, leagueTeam := range leagueTeams {
		vm := models.LeagueTeamDetailVM{}
		result = append(result, vm.FromDBModel(leagueTeam))
	}

	return successResult(ctx, result)
}

func (h LeagueTeamHandler) DeleteByLeagueTeamID(ctx *fiber.Ctx) error {
	id, err := strconv.ParseInt(ctx.Params("id"), 10, 64)
	if err != nil {
		return errorResult(ctx, errors.New("Geçersiz takım id"))
	}

	err = h.leagueTeamRepository.DeleteByLeagueTeamID(ctx.Context(), id)
	if err != nil {
		return errorResult(ctx, errors.New("Takım ligden silinirken hata oluştu"))
	}

	return successResult(ctx, "Takım ligden başarıyla silindi!")
}
