package handlers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/personal-project/pitch-league/models"
	"github.com/personal-project/pitch-league/repository"
)

type FieldHandler struct {
	BaseHandler[models.Field]
	fieldRepository repository.IFieldRepository
}

func NewFieldHandler(r repository.IFieldRepository) FieldHandler {
	return FieldHandler{
		BaseHandler: BaseHandler[models.Field]{
			baseRepository: r,
		},
		fieldRepository: r,
	}
}

func (h FieldHandler) CreateField(ctx *fiber.Ctx) error {
	var vm models.FieldCreateVM
	if err := ctx.BodyParser(&vm); err != nil {
		return errorResult(ctx, err)
	}

	field := vm.ToDBModel(models.Field{})
	if err := h.fieldRepository.CreateField(ctx.Context(), field); err != nil {
		return errorResult(ctx, err)
	}

	return successResult(ctx, "Halısaha başarıyla eklendi!")
}

func (h FieldHandler) GetAllFields(ctx *fiber.Ctx) error {
	fields, err := h.fieldRepository.GetAllField(ctx.Context())
	if err != nil {
		return errorResult(ctx, err)
	}

	var result []models.FieldDetailVM
	for _, field := range fields {
		vm := models.FieldDetailVM{}
		result = append(result, vm.FromDBModel(field))
	}

	return successResult(ctx, result)
}

func (h FieldHandler) GetByFieldID(ctx *fiber.Ctx) error {
	id, err := strconv.ParseInt(ctx.Params("id"), 10, 64)
	if err != nil {
		return errorResult(ctx, err)
	}

	field, err := h.fieldRepository.GetByFieldID(ctx.Context(), id)
	if err != nil {
		return errorResult(ctx, err)
	}

	vm := models.FieldDetailVM{}
	result := vm.FromDBModel(*field)

	return successResult(ctx, result)
}

func (h FieldHandler) DeleteByFieldID(ctx *fiber.Ctx) error {
	id, err := strconv.ParseInt(ctx.Params("id"), 10, 64)
	if err != nil {
		return errorResult(ctx, err)
	}

	if err := h.fieldRepository.DeleteByFieldID(ctx.Context(), id); err != nil {
		return errorResult(ctx, err)
	}

	return successResult(ctx, "Halısaha başarıyla silindi!")
}

func (h FieldHandler) UpdateFieldByID(ctx *fiber.Ctx) error {
	id, err := strconv.ParseInt(ctx.Params("id"), 10, 64)
	if err != nil {
		return errorResult(ctx, err)
	}

	field, err := h.fieldRepository.GetByFieldID(ctx.Context(), id)
	if err != nil {
		return errorResult(ctx, err)
	}

	var vm models.FieldCreateVM
	if err := ctx.BodyParser(&vm); err != nil {
		return errorResult(ctx, err)
	}

	updatedField := vm.ToDBModel(*field)
	if err := h.fieldRepository.UpdateField(ctx.Context(), updatedField); err != nil {
		return errorResult(ctx, err)
	}

	return successResult(ctx, "Halısaha başarıyla güncellendi!")
}
