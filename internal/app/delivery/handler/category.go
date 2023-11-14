package handler

import (
	"errors"
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
	resp := helpers.New[[]dto.Category]()
	ctx := r.Context()
	data, err := h.uc.GetAll(ctx)
	if err != nil {
		var errUC helpers.ResponseError
		if errors.As(err, &errUC) {
			resp.Code(errUC.Code()).Error(errUC).Send(w)
		} else {
			resp.Code(http.StatusInternalServerError).Error(helpers.ErrInternal).Send(w)
		}
		return
	}

	resp.Code(http.StatusOK).Data(data).Send(w)
}
