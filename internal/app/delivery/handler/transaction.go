package handler

import (
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
	res := helpers.New[[]dto.Transaction]().ContentType("application/json")
	ctx := r.Context()
	data, errUC := h.uc.GetAll(ctx)
	if errUC != nil {
		res.Code(errUC.Code).Error(errUC.Error).Send(w)
		return
	}

	res.Code(http.StatusOK).Data(data).Send(w)
}
