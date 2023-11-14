package handler

import (
	"errors"
	"net/http"

	"github.com/thansetan/62teknologi-backend-test-fathan-arsyadani/internal/app/delivery/dto"
	"github.com/thansetan/62teknologi-backend-test-fathan-arsyadani/internal/app/usecase"
	"github.com/thansetan/62teknologi-backend-test-fathan-arsyadani/internal/helpers"
)

type transactionHandler struct {
	uc usecase.TransactionUsecase
}

func NewTransactionHandler(uc usecase.TransactionUsecase) *transactionHandler {
	return &transactionHandler{uc}
}

type getAllTransactionResponse helpers.Response[[]dto.Transaction]

// ListTransactions godoc
//
//	@Summary		List transactions
//	@Description	List all supported transaction types
//	@Tags			transactions
//	@Produce		json
//	@Success		200	{object}	getAllTransactionResponse{data=[]dto.Transaction}
//	@Failure		500	{object}	getAllTransactionResponse{error=string}
//
//	@Router			/transactions [get]
func (h *transactionHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	resp := helpers.New[[]dto.Transaction]()
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
