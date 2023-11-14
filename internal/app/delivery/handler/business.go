package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/thansetan/62teknologi-backend-test-fathan-arsyadani/internal/app/delivery/dto"
	"github.com/thansetan/62teknologi-backend-test-fathan-arsyadani/internal/app/usecase"
	"github.com/thansetan/62teknologi-backend-test-fathan-arsyadani/internal/helpers"
)

type businessHandler struct {
	uc usecase.BusinessUsecase
}

func NewBusinessHandler(uc usecase.BusinessUsecase) *businessHandler {
	return &businessHandler{uc}
}

type createBusiness helpers.Response[dto.CreateBusinessResponse]

// CreateBusiness godoc
//
//	@Summary		Create business
//	@Description	Add new business to database
//	@Tags			businesses
//	@Param			Body	body	dto.BusinessRequest	true	"data required to create a new business"
//	@Produce		json
//	@Success		200	{object}	createBusiness{data=dto.CreateBusinessResponse}
//	@Failure		400	{object}	getBusiness{error=string}
//	@Failure		500	{object}	getBusiness{error=string}
//	@Router			/businesses [post]
func (h *businessHandler) Create(w http.ResponseWriter, r *http.Request) {
	var reqData dto.BusinessRequest
	resp := helpers.New[*dto.CreateBusinessResponse]()

	err := json.NewDecoder(r.Body).Decode(&reqData)
	if err != nil {
		resp.Code(http.StatusBadRequest).Error(err).Send(w)
		return
	}

	ctx := r.Context()
	data, err := h.uc.Create(ctx, reqData)
	if err != nil {
		var errUC helpers.ResponseError
		if errors.As(err, &errUC) {
			var errValidation helpers.ValidationError
			if errors.As(errUC.Err, &errValidation) {
				resp.Code(errUC.Code()).Errors(errValidation.ErrSlice()).Send(w)
			} else {
				resp.Code(errUC.Code()).Error(errUC).Send(w)
			}
		} else {
			resp.Code(http.StatusInternalServerError).Error(helpers.ErrInternal).Send(w)
		}
		return
	}

	resp.Code(http.StatusCreated).Data(&data).Send(w)
}

// UpdateBusiness godoc
//
//	@Summary		Update business
//	@Description	Update business data by ID or alias
//	@Tags			businesses
//	@Param			business_id_or_alias	path	string				true	"business ID or alias"
//	@Param			Body					body	dto.BusinessRequest	true	"data required to update business data"
//	@Produce		json
//	@Success		204
//	@Failure		400	{object}	getBusiness{error=string}
//	@Failure		404	{object}	getBusiness{error=string}
//	@Failure		500	{object}	getBusiness{error=string}
//	@Router			/businesses/{business_id_or_alias} [put]
func (h *businessHandler) Update(w http.ResponseWriter, r *http.Request) {
	var reqData dto.BusinessRequest
	resp := helpers.New[error]()

	err := json.NewDecoder(r.Body).Decode(&reqData)
	if err != nil {
		resp.Code(http.StatusBadRequest).Error(err).Send(w)
		return
	}

	idOrAlias := mux.Vars(r)["id_or_alias"]
	ctx := r.Context()

	err = h.uc.Update(ctx, idOrAlias, reqData)
	if err != nil {
		var errUC helpers.ResponseError
		if errors.As(err, &errUC) {
			var errValidation helpers.ValidationError
			if errors.As(errUC.Err, &errValidation) {
				resp.Code(errUC.Code()).Errors(errValidation.ErrSlice()).Send(w)
			} else {
				resp.Code(errUC.Code()).Error(errUC).Send(w)
			}
		} else {
			resp.Code(http.StatusInternalServerError).Error(helpers.ErrInternal).Send(w)
		}
		return
	}

	resp.Code(http.StatusNoContent).Send(w)
}

type getBusiness helpers.Response[dto.BusinessResponse]

// GetBusiness godoc
//
//	@Summary		Get business
//	@Description	Get a single business by ID or alias
//	@Tags			businesses
//	@Param			business_id_or_alias	path	string	true	"business ID or alias"
//	@Produce		json
//	@Success		200	{object}	getBusiness{data=dto.BusinessResponse}
//	@Failure		404	{object}	getBusiness{error=string}
//	@Failure		500	{object}	getBusiness{error=string}
//	@Router			/businesses/{business_id_or_alias} [get]
func (h *businessHandler) GetBusiness(w http.ResponseWriter, r *http.Request) {
	resp := helpers.New[*dto.BusinessResponse]()
	idOrAlias := mux.Vars(r)["id_or_alias"]

	ctx := context.WithValue(r.Context(), "host", r.Host)
	business, err := h.uc.GetBusiness(ctx, idOrAlias)
	if err != nil {
		var errUC helpers.ResponseError
		if errors.As(err, &errUC) {
			resp.Code(errUC.Code()).Error(errUC).Send(w)
		} else {
			resp.Code(http.StatusInternalServerError).Error(helpers.ErrInternal).Send(w)
		}
		return
	}

	resp.Code(http.StatusOK).Data(&business).Send(w)
}

// DeleteBusiness godoc
//
//	@Summary		Delete business
//	@Description	Delete a single business by ID or alias
//	@Tags			businesses
//	@Param			business_id_or_alias	path	string	true	"business ID or alias"
//	@Success		204
//	@Failure		404	{object}	getBusiness{error=string}
//	@Failure		500	{object}	getBusiness{error=string}
//	@Router			/businesses/{business_id_or_alias} [delete]
func (h *businessHandler) Delete(w http.ResponseWriter, r *http.Request) {
	resp := helpers.New[error]()
	idOrAlias := mux.Vars(r)["id_or_alias"]
	ctx := r.Context()

	err := h.uc.Delete(ctx, idOrAlias)
	if err != nil {
		var errUC helpers.ResponseError
		if errors.As(err, &errUC) {
			resp.Code(errUC.Code()).Error(errUC).Send(w)
		} else {
			resp.Code(http.StatusInternalServerError).Error(helpers.ErrInternal).Send(w)
		}
		return
	}

	resp.Code(http.StatusNoContent).Send(w)
}

type getBusinessesResponse helpers.Response[[]dto.BusinessResponse]

// SearchBusinesses godoc
//
//	@Summary		Search businesses
//	@Description	Search for businesses by provided parameters
//	@Tags			businesses
//	@Param			location		query		string	true	"business location (ex:jakarta/roma/new york/etc)"
//	@Param			categories		query		string	false	"business categories. can be multiple. comma separated (ex:italian,japanese,etc)"
//	@Param			transactions	query		string	false	"business transaction type. can be multiple. comma separated (ex:delivery,pickup,etc)"
//	@Param			open_now		query		bool	false	"is the business open right now. open_now and open_at can't be used together"
//	@Param			per_page		query		int		false	"number of business(es) to display per page"																			default(5)		minimum(1)
//	@Param			page			query		int		false	"page number you want to see"																							default(1)		minimum(1)
//	@Param			open_at			query		string	false	"will businesses open at given time. format: 24h, between 0000 to 2359. open_now and open_at can't be used together"	minlength(4)	maxlength(4)
//	@Success		200				{object}	getBusinessesResponse{data=[]dto.BusinessResponse,metadata=dto.Metadata}
//	@Failure		400				{object}	getBusinessesResponse{error=string}
//	@Failure		500				{object}	getBusinessesResponse{error=string}
//	@Router			/businesses/search [get]
func (h *businessHandler) GetBusinesses(w http.ResponseWriter, r *http.Request) {
	queryParams := dto.BusinessQueryParams{}
	resp := helpers.New[*[]dto.BusinessResponse]()

	err := helpers.ParseQuery(r.URL.Query(), &queryParams)
	if err != nil {
		resp.Code(http.StatusBadRequest).Error(err).Send(w)
		return
	}

	ctx := context.WithValue(r.Context(), "host", r.Host)
	data, err := h.uc.GetBusinesses(ctx, queryParams)
	if err != nil {
		var errUC helpers.ResponseError
		if errors.As(err, &errUC) {
			var errValidation helpers.ValidationError
			if errors.As(errUC.Err, &errValidation) {
				resp.Code(errUC.Code()).Errors(errValidation.ErrSlice()).Send(w)
			} else {
				resp.Code(errUC.Code()).Error(errUC).Send(w)
			}
		} else {
			resp.Code(http.StatusInternalServerError).Error(helpers.ErrInternal).Send(w)
		}
		return
	}

	resp.Code(http.StatusOK).Metadata(data.Metadata).Data(&data.Businesses).Send(w)
}
