package controllers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/diegoclair/ApiGolang/api/models"
	"github.com/diegoclair/ApiGolang/api/responses"
	"github.com/gorilla/mux"
)

func (server *Server) GetReportsByUserId(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	uid, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	reportsByUser := models.Reports{}
	userReports, err := reportsByUser.FindReportsByUserID(server.DB, uint32(uid))
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	responses.JSON(w, http.StatusOK, userReports)
}

func (server *Server) GetReportsByDate(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	date := vars["date"]
	err := errors.New("You need to pass the date after URL reports/day/ ")

	if date == "" {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	reportsByDate := models.Reports{}
	dateReports, err := reportsByDate.FindReportsByDate(server.DB, date)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	responses.JSON(w, http.StatusOK, dateReports)
}
