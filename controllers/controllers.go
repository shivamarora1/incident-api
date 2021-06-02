package controllers

import (
	"encoding/json"

	"io/ioutil"
	"net/http"

	"example.com/incident-api/config"
	"example.com/incident-api/database"
	"example.com/incident-api/models"
	"example.com/incident-api/utils"
	"go.uber.org/zap"
)

func AllIncidentGetter(w http.ResponseWriter, r *http.Request) {
	var res []byte
	if incidents, err := database.DB.GetAllIncidents(0, 0); err != nil {
		config.Logger.Error("Error in getting all incidents", zap.Error(err))
		res = utils.ResponseFailure(err.Error())
	} else {
		res = utils.ResponseSuccess(incidents)
	}
	w.Write(res)
}

func IncidentSaver(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		config.Logger.Error("error in reading request body", zap.Error(err))
		w.Write(utils.ResponseFailure(err.Error()))
		return
	}

	config.Logger.Info("Recieved body", zap.String("body", string(body)))

	//Parsing it to incident struct
	inc := new(models.Incident)
	if err := json.Unmarshal(body, inc); err != nil {
		config.Logger.Error("Error in unmarshalling body",
			zap.Error(err), zap.String("bodyStr", string(body)))
		w.Write(utils.ResponseFailure(err.Error()))
		return
	}
	if err := inc.Validate(); err != nil {
		config.Logger.Info("invalid incident object", zap.Error(err))
		w.Write(utils.ResponseFailure(err.Error()))
		return
	}
	if err := database.DB.SaveIncident(*inc); err != nil {
		config.Logger.Error("can't save incident", zap.Error(err))
		w.Write(utils.ResponseFailure(err.Error()))
		return
	} else {
		w.Write(utils.ResponseSuccess("Saved new incident"))
	}

}
