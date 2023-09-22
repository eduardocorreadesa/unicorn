package controller

import (
	"context"
	"errors"
	"net/http"
	"unicorn/internal/configuration"
	"unicorn/internal/domain"
	"unicorn/internal/service"

	"github.com/go-playground/validator"
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
)

type UnicornController interface {
	GetUnicorn(w http.ResponseWriter, r *http.Request) (http.ResponseWriter, int, []domain.Unicorn, error)
	CreateUnicorn(w http.ResponseWriter, r *http.Request) (http.ResponseWriter, int, domain.UnicornProcess, error)
}

func NewUnicornController(ctx context.Context, config configuration.Config, unicornService service.UnicornService) UnicornController {
	return &unicorn{
		ctx:            ctx,
		config:         config,
		unicornService: unicornService,
	}
}

type unicorn struct {
	ctx            context.Context
	config         configuration.Config
	unicornService service.UnicornService
}

func (u *unicorn) CreateUnicorn(w http.ResponseWriter, r *http.Request) (http.ResponseWriter, int, domain.UnicornProcess, error) {
	var (
		params   domain.URLParamsDefault
		decoder  = schema.NewDecoder()
		validate = validator.New()
		err      error
		id       string
	)

	err = decoder.Decode(&params, r.URL.Query())
	if err != nil {
		return w, http.StatusBadRequest, domain.UnicornProcess{}, err
	}

	err = validate.Struct(&params)
	if err != nil {
		return w, http.StatusBadRequest, domain.UnicornProcess{}, err
	}

	id = u.unicornService.CreateAsync(params.Amount)

	return w, http.StatusOK, domain.UnicornProcess{RequestID: id}, nil
}

func (u *unicorn) GetUnicorn(w http.ResponseWriter, r *http.Request) (http.ResponseWriter, int, []domain.Unicorn, error) {
	var (
		param     = mux.Vars(r)
		requestID string
	)

	requestID = param["id"]

	res, ok := u.unicornService.GetResult(requestID)
	if !ok {
		return w, http.StatusOK, res.Unicorns, errors.New("We are still processing your Unicorn. Please make sure the request_id is correct.")
	}

	return w, http.StatusOK, res.Unicorns, nil
}
