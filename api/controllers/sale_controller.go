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

func (server *Server) CreateSale(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	sale := models.Sale{}
	err = json.Unmarshal(body, &sale)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	sale.Prepare(server.DB)
	err = sale.Validate(server.DB)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	uid, err := auth.ExtractTokenID(r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}
	if uid != sale.AuthorID {
		responses.ERROR(w, http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
		return
	}
	saleCreated, err := sale.SaveSale(server.DB)
	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		responses.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}
	w.Header().Set("Lacation", fmt.Sprintf("%s%s/%d", r.Host, r.URL.Path, saleCreated.ID))
	responses.JSON(w, http.StatusCreated, saleCreated)
}

func (server *Server) GetSales(w http.ResponseWriter, r *http.Request) {

	sale := models.Sale{}

	sales, err := sale.FindAllSales(server.DB)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	responses.JSON(w, http.StatusOK, sales)
}

func (server *Server) GetSale(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	pid, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	sale := models.Sale{}

	saleReceived, err := sale.FindSaleByID(server.DB, pid)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	responses.JSON(w, http.StatusOK, saleReceived)
}

func (server *Server) UpdateSale(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	// Check if the sale id is valid
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

	// Check if the sale exist
	sale := models.Sale{}
	err = server.DB.Debug().Model(models.Sale{}).Where("id = ?", pid).Take(&sale).Error
	if err != nil {
		responses.ERROR(w, http.StatusNotFound, errors.New("Sale not found"))
		return
	}

	// If a user attempt to update a sale not belonging to him
	if uid != sale.AuthorID {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}
	// Read the data Sold
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	// Start processing the request data
	saleUpdate := models.Sale{}
	err = json.Unmarshal(body, &saleUpdate)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	//Also check if the request user id is equal to the one gotten from token
	if uid != saleUpdate.AuthorID {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	saleUpdate.Prepare(server.DB)
	err = saleUpdate.Validate(server.DB)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	saleUpdate.ID = sale.ID //this is important to tell the model the sale id to update, the other update field are set above

	saleUpdated, err := saleUpdate.UpdateASale(server.DB)

	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		responses.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}
	responses.JSON(w, http.StatusOK, saleUpdated)
}

func (server *Server) DeleteSale(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	// Is a valid sale id given to us?
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

	// Check if the sale exist
	sale := models.Sale{}
	err = server.DB.Debug().Model(models.Sale{}).Where("id = ?", pid).Take(&sale).Error
	if err != nil {
		responses.ERROR(w, http.StatusNotFound, errors.New("Unauthorized"))
		return
	}

	// Is the authenticated user, the owner of this sale?
	if uid != sale.AuthorID {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}
	_, err = sale.DeleteASale(server.DB, pid, uid)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	w.Header().Set("Entity", fmt.Sprintf("%d", pid))
	responses.JSON(w, http.StatusNoContent, "")
}
