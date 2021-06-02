# Incident Management API
 An incident is an event that could lead to loss of, or disruption to, an organization's operations, services, or functions. Incident management is a term describing the activities of an organization to identify, analyze, and correct hazards to prevent future re-occurrence.

- REST API
- Developed in Go
- Postgres database (You can use any)
- Automatic migrations and seeding of data

## How to run
- Install Golang
- Clone this repository
- Create a "config.json" file in "config/" directory

### Config file format
```JSON
{
   "port":8080,
   "lp":"logs/logs.log",
   "db":{
      "dt":"postgres",
      "cs":"host=example.ap-south-1.rds.amazonaws.com port=5432 dbname=<db-name> user=<user-name> password=<password>"
   }
}
```
"port" is port number in which you want to run your API server. "lp" is path where you want to store your log files. Database configuration is defined in "db" key, "dt" is database name (Current only Postgres supported), "cs" is connection string in given format.


After setting up config.json file and databasce credentials just run server using,
#### go run main.go