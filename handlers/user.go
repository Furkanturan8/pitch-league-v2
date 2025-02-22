package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/personal-project/pitch-league/models"
	"github.com/personal-project/pitch-league/repository"
	"strconv"
)

type UserHandler struct {
	BaseHandler[models.User]
}

func NewUserHandler(repo repository.IBaseRepository[models.User]) UserHandler {
	return UserHandler{
		BaseHandler[models.User]{
			baseRepository: repo,
		},
	}
}

func (h UserHandler) CreateUser(ctx *fiber.Ctx) error {
	var createModel models.UserCreate
	if err := ctx.BodyParser(&createModel); err != nil {
		return errorResult(ctx, err)
	}

	user := createModel.ToModel()
	user.Role = models.UserRoleNormal

	createdUser, err := h.baseRepository.Create(ctx.Context(), user)
	if err != nil {
		return errorResult(ctx, err)
	}

	return successResult(ctx, models.ToUserResponse(createdUser))
}

func (h UserHandler) CreateAdmin(ctx *fiber.Ctx) error {
	var createModel models.UserCreate
	if err := ctx.BodyParser(&createModel); err != nil {
		return errorResult(ctx, err)
	}

	user := createModel.ToModel()
	user.Role = models.UserRoleAdmin

	createdUser, err := h.baseRepository.Create(ctx.Context(), user)
	if err != nil {
		return errorResult(ctx, err)
	}

	return successResult(ctx, models.ToUserResponse(createdUser))
}

func (h UserHandler) GetAllUsers(ctx *fiber.Ctx) error {
	users, err := h.baseRepository.GetAll(ctx.Context())
	if err != nil {
		return errorResult(ctx, err)
	}

	userResponses := make([]models.UserResponse, len(users))
	for i, user := range users {
		userResponses[i] = models.ToUserResponse(user)
	}

	return successResult(ctx, userResponses)
}

func (h UserHandler) GetByUserID(ctx *fiber.Ctx) error {
	id, err := strconv.ParseInt(ctx.Params("id"), 10, 64)
	if err != nil {
		return errorResult(ctx, err)
	}

	user, err := h.baseRepository.GetByID(ctx.Context(), id)
	if err != nil {
		return errorResult(ctx, err)
	}

	return successResult(ctx, models.ToUserResponse(user))
}

func (h UserHandler) DeleteByUserID(ctx *fiber.Ctx) error {
	id, err := strconv.ParseInt(ctx.Params("id"), 10, 64)
	if err != nil {
		return errorResult(ctx, err)
	}

	err = h.baseRepository.Delete(ctx.Context(), id)
	if err != nil {
		return errorResult(ctx, err)
	}

	return successResult(ctx, "User deleted successfully")
}

func (h UserHandler) UpdateUserByID(ctx *fiber.Ctx) error {
	id, err := strconv.ParseInt(ctx.Params("id"), 10, 64)
	if err != nil {
		return errorResult(ctx, err)
	}

	user, err := h.baseRepository.GetByID(ctx.Context(), id)
	if err != nil {
		return errorResult(ctx, err)
	}

	var updateModel models.UserUpdate
	if err := ctx.BodyParser(&updateModel); err != nil {
		return errorResult(ctx, err)
	}

	updatedUser := updateModel.ToModel(user)
	err = h.baseRepository.Update(ctx.Context(), updatedUser)
	if err != nil {
		return errorResult(ctx, err)
	}

	return successResult(ctx, models.ToUserResponse(updatedUser))
}
