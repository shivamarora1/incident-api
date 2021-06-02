package models

import (
	"fmt"
	"time"

	"example.com/incident-api/config"
	"go.uber.org/zap"
)

const (
	iNCIDENT_CATEGORY_SECURITY        = 1
	iNCIDENT_CATEGORY_HEALTH          = 2
	iNCIDENT_CATEGORY_LOSS_PREVENTION = 3
)

type Incident struct {
	Id             int              `json:"id"`
	Category       int              `json:"category"`
	Location       IncidentLocation `json:"location"`
	Title          string           `json:"title"`
	PeopleAffected []People         `json:"people"`
	Comments       string           `json:"comments"`
	IncidentDate   string           `json:"incidentDate"`
	CreateDate     string           `json:"createDate"`
	ModifyDate     string           `json:"modifyDate"`

	IncidentDateObj time.Time `json:"-"`
	CreateDateObj   time.Time `json:"-"`
	ModifyDateObj   time.Time `json:"-"`
}

type IncidentLocation struct {
	Latitude float64 `json:"latitude"`
	Logitude float64 `json:"longitude"`
}

type People struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

//Function is used to validate incident params
func (inc *Incident) Validate() error {

	if inc.Category == 0 || (inc.Category != iNCIDENT_CATEGORY_HEALTH &&
		inc.Category != iNCIDENT_CATEGORY_LOSS_PREVENTION &&
		inc.Category != iNCIDENT_CATEGORY_SECURITY) {
		return fmt.Errorf("Category of incident is not valid")
	} else if inc.Title == "" {
		return fmt.Errorf("Incident title is not valid")
	} else if inc.Location.Latitude == 0 || inc.Location.Logitude == 0 {
		return fmt.Errorf("Logitude or Latitude is not correct")
	} else if inc.IncidentDate == "" {
		return fmt.Errorf("Incident date not specified")
	} else if incDate, err := time.Parse(time.RFC3339, inc.IncidentDate); err != nil {
		config.Logger.Error("Error in parsing date",
			zap.Error(err), zap.String("date", inc.IncidentDate))
		return fmt.Errorf("Incident date is not in YYYY-MM-DDTHH:MM:SSZhh:mm format")
	} else {
		//Validation is completed
		inc.IncidentDateObj = incDate

		if mObj, err := time.Parse(time.RFC3339, inc.ModifyDate); err != nil {
			config.Logger.Info("Error in parsing date",
				zap.Error(err), zap.String("date", inc.ModifyDate))
			//If last update time is not valid or blank consider current time
			inc.ModifyDateObj = time.Now()
		} else {
			inc.ModifyDateObj = mObj
		}

		if cObj, err := time.Parse(time.RFC3339, inc.CreateDate); err != nil {
			config.Logger.Info("Error in parsing date",
				zap.Error(err), zap.String("date", inc.CreateDate))
			//If created time is not valid or blank consider current time
			inc.CreateDateObj = time.Now()
		} else {
			inc.CreateDateObj = cObj
		}
		return nil
	}
}
