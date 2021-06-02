package database

import (
	"errors"

	"example.com/incident-api/config"
	"example.com/incident-api/models"
)

const (
	POSTGRES = "postgres"
	///Add other DBs here
)

var DB Database

type Database interface {
	Connect() error
	GetIncidents(incidentId int) (models.Incident, error)
	GetAllIncidents(limit, offset int) ([]models.Incident, error)
	SaveIncident(inc models.Incident) error
}

func InitDatabase() error {
	switch config.ConfigObj.DBConfig.DatabaseType {
	case POSTGRES:
		DB = new(PostgresDB)
		return nil
	default:
		return errors.New("DB of specified config not" +
			" found/Currently only Postgres supported")
	}
}
