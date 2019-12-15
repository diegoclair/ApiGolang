package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/diegoclair/ApiGolang/api/auth"
	"github.com/diegoclair/ApiGolang/api/models"
	"github.com/diegoclair/ApiGolang/api/responses"
	"github.com/diegoclair/ApiGolang/api/utils/formaterror"
	"github.com/gorilla/mux"
)

func (server *Server) CreateBuy(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	buy := models.Buy{}
	err = json.Unmarshal(body, &buy)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	buy.Prepare(server.DB)
	err = buy.Validate(server.DB)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	uid, err := auth.ExtractTokenID(r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}
	if uid != buy.AuthorID {
		responses.ERROR(w, http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
		return
	}
	buyCreated, err := buy.SaveBuy(server.DB)
	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		responses.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}
	w.Header().Set("Lacation", fmt.Sprintf("%s%s/%d", r.Host, r.URL.Path, buyCreated.ID))
	responses.JSON(w, http.StatusCreated, buyCreated)
}

func (server *Server) GetBuys(w http.ResponseWriter, r *http.Request) {

	buy := models.Buy{}

	buys, err := buy.FindAllBuys(server.DB)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	responses.JSON(w, http.StatusOK, buys)
}

func (server *Server) GetBuy(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	// Check if the buy id is valid
	pid, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	buy := models.Buy{}

	buyReceived, err := buy.FindBuyByID(server.DB, pid)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	responses.JSON(w, http.StatusOK, buyReceived)
}

func (server *Server) UpdateBuy(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	// Check if the buy id is valid
	pid, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	//CHeck if the auth token is valid and  get the user id from it
	uid, err := auth.ExtractTokenID(r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	// Check if the buy exist
	buy := models.Buy{}
	err = server.DB.Debug().Model(models.Buy{}).Where("id = ?", pid).Take(&buy).Error
	if err != nil {
		responses.ERROR(w, http.StatusNotFound, errors.New("Buy not found"))
		return
	}

	// If a user attempt to update a buy not belonging to him
	if uid != buy.AuthorID {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}
	// Read the data buyed
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	// Start processing the request data
	buyUpdate := models.Buy{}
	err = json.Unmarshal(body, &buyUpdate)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	//Also check if the request user id is equal to the one gotten from token
	if uid != buyUpdate.AuthorID {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	buyUpdate.Prepare(server.DB)
	err = buyUpdate.Validate(server.DB)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	buyUpdate.ID = buy.ID //this is important to tell the model the buy id to update, the other update field are set above

	buyUpdated, err := buyUpdate.UpdateABuy(server.DB)

	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		responses.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}
	responses.JSON(w, http.StatusOK, buyUpdated)
}

func (server *Server) DeleteBuy(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	// Is a valid buy id given to us?
	pid, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	// Is this user authenticated?
	uid, err := auth.ExtractTokenID(r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	// Check if the buy exist
	buy := models.Buy{}
	err = server.DB.Debug().Model(models.Buy{}).Where("id = ?", pid).Take(&buy).Error
	if err != nil {
		responses.ERROR(w, http.StatusNotFound, errors.New("Unauthorized"))
		return
	}

	// Is the authenticated user, the owner of this buy?
	if uid != buy.AuthorID {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}
	_, err = buy.DeleteABuy(server.DB, pid, uid)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	w.Header().Set("Entity", fmt.Sprintf("%d", pid))
	responses.JSON(w, http.StatusNoContent, "")
}
