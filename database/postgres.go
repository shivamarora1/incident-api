package database

import (
	"database/sql"
	"encoding/json"

	"time"

	"example.com/incident-api/config"
	"example.com/incident-api/models"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

type PostgresDB struct {
	dao *sql.DB
}

func (pDB *PostgresDB) Connect() error {
	db, err := sql.Open("postgres", config.ConfigObj.DBConfig.ConnectionString)
	if err != nil {
		return err
	}
	err = db.Ping()
	if err != nil {
		return err
	}

	pDB.dao = db

	pDB.seed()
	//Seed data here after connection
	config.Logger.Info("Database connected!!!")
	return nil
}

func (pDB *PostgresDB) GetIncidents(incidentId int) (models.Incident, error) {
	//To be implement for getting information about any specific incident
	return models.Incident{}, nil
}

//Get information about all incidents in the DB
//TODO: We can also control the list by providing limit and offset
//TODO: But that functionality has not been implemented yet
func (pDB *PostgresDB) GetAllIncidents(limit, offset int) ([]models.Incident, error) {
	incidents := []models.Incident{}

	//In future we can also consider limit and offset params
	query := `select id,title,category,latitude,longitude,peoples,comments,
						extract(epoch from incidenttime)::int,extract(epoch from updatetime)::int,
						extract(epoch from createdtime)::int
			  from incidents`

	if rows, err := pDB.dao.Query(query); err == nil {
		defer rows.Close()
		for rows.Next() {
			inc := models.Incident{}

			var id, catgory, incidenT, updteTime, cTime sql.NullInt64
			var title, comments, peoplesStr sql.NullString
			var lat, long sql.NullFloat64
			if err := rows.Scan(&id, &title, &catgory, &lat, &long,
				&peoplesStr, &comments, &incidenT, &updteTime, &cTime); err != nil {
				config.Logger.Error("unable to scan row", zap.String("error", err.Error()))
				continue
			}
			inc.Id = int(id.Int64)
			inc.Title = title.String
			inc.Category = int(catgory.Int64)
			inc.Location.Latitude = lat.Float64
			inc.Location.Logitude = long.Float64
			inc.Comments = comments.String
			inc.IncidentDate = time.Unix(incidenT.Int64, 0).Format(time.RFC3339)
			inc.CreateDate = time.Unix(cTime.Int64, 0).Format(time.RFC3339)
			inc.ModifyDate = time.Unix(updteTime.Int64, 0).Format(time.RFC3339)

			peoples := []models.People{}
			if err := json.Unmarshal([]byte(peoplesStr.String), &peoples); err != nil {
				config.Logger.Error("Cannot unmarshall string to object",
					zap.String("string", peoplesStr.String),
					zap.Error(err))
			}
			inc.PeopleAffected = peoples
			incidents = append(incidents, inc)

		}
		return incidents, nil
	} else {
		return nil, err
	}

}

//Save incident in Database
func (pDB *PostgresDB) SaveIncident(inc models.Incident) error {

	query := `insert into incidents (title,category,latitude,longitude,
										peoples,comments,incidenttime,
										updatetime,createdtime)
			values($1,$2,$3,$4,$5,$6,to_timestamp($7),to_timestamp($8),to_timestamp($9));`

	//Converting struct to JSON string
	peoplesStr, err := json.Marshal(inc.PeopleAffected)
	if err != nil {
		config.Logger.Error("Error in marshalling object ",
			zap.Error(err),
			zap.Any("struct", inc.PeopleAffected))
		return err
	}

	if res, err := pDB.dao.Exec(query, inc.Title, inc.Category, inc.Location.Latitude,
		inc.Location.Logitude, peoplesStr, inc.Comments, inc.IncidentDateObj.Unix(),
		inc.ModifyDateObj.Unix(), inc.CreateDateObj.Unix()); err != nil {
		config.Logger.Error("Error in executing query",
			zap.Error(err), zap.String("query", query))
		return err
	} else if lId, err := res.LastInsertId(); err == nil {
		config.Logger.Info("Inserted new incident ",
			zap.Int64("id", lId))
	}
	return nil
}

//Seeding data into data table
func (pDB *PostgresDB) seed() {
	if _, err := pDB.dao.Exec(`CREATE TABLE IF NOT EXISTS incidents(id serial primary key,title text,category int, latitude real, longitude real,
									peoples jsonb,comments text, incidenttime timestamp with time zone default now(),
									updatetime timestamp with time zone default now(), createdtime timestamp with time zone default now());`); err != nil {
		config.Logger.Warn("cannot create incident table", zap.String("error", err.Error()))
	}
	if _, err := pDB.dao.Exec(`insert into incidents (title,category,latitude,
														longitude,peoples,comments,
														incidenttime,updatetime,createdtime)
								values('Customers data leak',1,19.0760,72.8777,
										'[{"name":"Mohit","type":"staff"},{"name":"Priyanka","type":"witness"}]',
										'Data leak breach occured due to careless of Mohit',now()-interval'1 day',now(),now()),
									('Covid Cases',2,19.0760,72.8777,
										'[{"name":"Anita","type":"staff"},{"name":"Sushma","type":"staff"},{"name":"Ajay","type":"witness"}]',
										'Covid symptoms founded in Sushma and Anita',now()-interval'12 hours',now(),now())
								on conflict(id) do nothing;`); err != nil {
		config.Logger.Warn("Unable to insert data into incident table", zap.String("error", err.Error()))
	}
}
