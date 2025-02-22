package handlers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/personal-project/pitch-league/models"
	"github.com/personal-project/pitch-league/repository"
)

type TeamHandler struct {
	BaseHandler[models.Team]
	teamRepository repository.ITeamRepository
}

func NewTeamHandler(r repository.ITeamRepository) TeamHandler {
	return TeamHandler{
		BaseHandler: BaseHandler[models.Team]{
			baseRepository: r,
		},
		teamRepository: r,
	}
}

func (h TeamHandler) CreateTeam(ctx *fiber.Ctx) error {
	var vm models.TeamCreateVM
	if err := ctx.BodyParser(&vm); err != nil {
		return errorResult(ctx, err)
	}

	team := vm.ToDBModel(models.Team{})
	if err := h.teamRepository.CreateTeam(ctx.Context(), team); err != nil {
		return errorResult(ctx, err)
	}

	detailVM := models.TeamDetailVM{}
	result := detailVM.FromDBModel(team)

	return successResult(ctx, result)
}

func (h TeamHandler) GetAllTeams(ctx *fiber.Ctx) error {
	teams, err := h.teamRepository.GetAllTeam(ctx.Context())
	if err != nil {
		return errorResult(ctx, err)
	}

	var result []models.TeamDetailVM
	for _, team := range teams {
		vm := models.TeamDetailVM{}
		result = append(result, vm.FromDBModel(team))
	}

	return successResult(ctx, result)
}

func (h TeamHandler) GetByTeamID(ctx *fiber.Ctx) error {
	id, err := strconv.ParseInt(ctx.Params("id"), 10, 64)
	if err != nil {
		return errorResult(ctx, err)
	}

	m, err := h.teamRepository.GetByTeamID(ctx.Context(), id)
	if err != nil {
		return errorResult(ctx, err)
	}

	vm := models.TeamDetailVM{}
	result := vm.FromDBModel(*m)

	return successResult(ctx, result)
}

func (h TeamHandler) DeleteByTeamID(ctx *fiber.Ctx) error {
	id, err := strconv.ParseInt(ctx.Params("id"), 10, 64)
	if err != nil {
		return errorResult(ctx, err)
	}

	if err := h.teamRepository.DeleteByTeamID(ctx.Context(), id); err != nil {
		return errorResult(ctx, err)
	}

	return successResult(ctx, "Team deleted successfully")
}

func (h TeamHandler) UpdateTeamByID(ctx *fiber.Ctx) error {
	id, err := strconv.ParseInt(ctx.Params("id"), 10, 64)
	if err != nil {
		return errorResult(ctx, err)
	}

	m, err := h.teamRepository.GetByTeamID(ctx.Context(), id)
	if err != nil {
		return errorResult(ctx, err)
	}

	var vm models.TeamCreateVM
	if err := ctx.BodyParser(&vm); err != nil {
		return errorResult(ctx, err)
	}

	updatedTeam := vm.ToDBModel(*m)
	if err := h.teamRepository.UpdateTeam(ctx.Context(), updatedTeam); err != nil {
		return errorResult(ctx, err)
	}

	return successResult(ctx, updatedTeam)
}

func (h TeamHandler) JoinTeam(ctx *fiber.Ctx) error {
	teamID, err := strconv.ParseInt(ctx.Params("id"), 10, 64)
	if err != nil {
		return errorResult(ctx, err)
	}

	userID, err := strconv.ParseInt(ctx.Params("userID"), 10, 64)
	if err != nil {
		return errorResult(ctx, err)
	}

	if err := h.teamRepository.AddUserToTeam(ctx.Context(), userID, teamID); err != nil {
		return errorResult(ctx, err)
	}

	return successResult(ctx, "Successfully joined the team")
}
