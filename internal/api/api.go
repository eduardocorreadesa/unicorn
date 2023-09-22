package api

import (
	"context"
	"encoding/json"
	"net/http"
	"path/filepath"
	"unicorn/docs"
	"unicorn/internal/configuration"
	"unicorn/internal/controller"

	"github.com/gorilla/mux"

	httpSwagger "github.com/swaggo/http-swagger"
)

type API struct {
	Router            *mux.Router
	Config            configuration.Config
	Ctx               context.Context
	UnicornController controller.UnicornController
}

func NewAPI(ctx context.Context, config configuration.Config,
	unicornController controller.UnicornController) *API {
	return &API{
		Config:            config,
		Ctx:               ctx,
		Router:            mux.NewRouter(),
		UnicornController: unicornController,
	}
}

//	@title			UNICORN API
//	@version		1.0
//	@description	This is a API for UNICORN application
func (a *API) SetupServer() {

	docs.SwaggerInfo.Host = a.Config.Server.Host
	docs.SwaggerInfo.Version = "1.0"

	r := a.Router

	r.HandleFunc(filepath.Join(a.Config.Server.Prefix, "/unicorn/{id}"), a.getUnicorn).Methods(http.MethodGet)
	r.HandleFunc(filepath.Join(a.Config.Server.Prefix, "/unicorn"), a.createUnicorn).Methods(http.MethodGet)

	//Swagger
	r.PathPrefix("/swagger").Handler(httpSwagger.WrapHandler)

	http.ListenAndServe(a.Config.Server.Port, r)
}

// getUnicorn
//
//	@Summary	Search Unicorn Process.
//	@ID			getUnicorn
//	@Tags		Unicorn
//	@Produce	json
//	@Param		id	path		string				true	"id"
//	@Success	200	{object}	[]domain.Unicorn	"ok"
//	@Failure	400	{object}	nil					"Bad request"
//	@Failure	500	{object}	nil					"Error in process"
//	@Router		/api/unicorn/{id} [get]
func (a *API) getUnicorn(w http.ResponseWriter, r *http.Request) {
	resp, statusCode, result, err := a.UnicornController.GetUnicorn(w, r)
	if err != nil {
		respondWithJSON(resp, statusCode, err.Error())
		return
	}

	respondWithJSON(resp, statusCode, result)
}

// createUnicorn
//
//	@Summary	Create Unicorn Process.
//	@ID			createUnicorn
//	@Tags		Unicorn
//	@Produce	json
//	@Param		amount	query		string					false	"/api/create-unicorn?amount=..."
//	@Success	200		{object}	domain.UnicornProcess	"ok"
//	@Failure	400		{object}	nil						"Bad request"
//	@Failure	500		{object}	nil						"Error in process"
//	@Router		/api/unicorn [get]
func (a *API) createUnicorn(w http.ResponseWriter, r *http.Request) {
	resp, statusCode, result, err := a.UnicornController.CreateUnicorn(w, r)
	if err != nil {
		respondWithJSON(resp, statusCode, err.Error())
		return
	}

	respondWithJSON(resp, statusCode, result)
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		respondWithJSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
