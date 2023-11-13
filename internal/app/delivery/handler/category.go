package handler

import (
	"net/http"

	"github.com/thansetan/62teknologi-backend-test-fathan-arsyadani/internal/app/delivery/dto"
	"github.com/thansetan/62teknologi-backend-test-fathan-arsyadani/internal/app/usecase"
	"github.com/thansetan/62teknologi-backend-test-fathan-arsyadani/internal/helpers"
)

type categoryHandler struct {
	uc usecase.CategoryUsecase
}

func NewCategoryHandler(uc usecase.CategoryUsecase) *categoryHandler {
	return &categoryHandler{uc}
}

type getAllCategoryResponse helpers.Response[[]dto.Category]

// ListCategories godoc
//
//	@Summary		List categories
//	@Description	List all available categories
//	@Tags			categories
//	@Produce		json
//	@Success		200	{object}	getAllCategoryResponse{data=[]dto.Category}
//	@Failure		500	{object}	getAllCategoryResponse{error=string}
//
//	@Router			/categories [get]
func (h *categoryHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	res := helpers.New[[]dto.Category]().ContentType("application/json")
	ctx := r.Context()
	data, errUC := h.uc.GetAll(ctx)
	if errUC != nil {
		res.Code(errUC.Code).Error(errUC.Error).Send(w)
		return
	}

	res.Code(http.StatusOK).Data(data).Send(w)
}
