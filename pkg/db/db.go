package db

import (
	"fmt"

	"github.com/Useurmind/go-greetings/pkg/db/pgsql"
	"github.com/Useurmind/go-greetings/pkg/db/mongodb"
	"github.com/Useurmind/go-greetings/pkg/db/cosmosdb"
)

const DBTypePGSQL = "pgsql"
const DBTypeMongoDB = "mongodb"
const DBTypeCosmosDB = "cosmosdb"

type DBContext interface {
	SaveGreeting(greetedPerson string, greeting string) error
	GetGreeting(greetedPerson string) (*string, error)
}

func NewDBContext(dbType string, dataSource string) (DBContext, error) {
	switch dbType {
	case DBTypePGSQL:
		return pgsql.NewPGSqlDBContext(dataSource)
	case DBTypeMongoDB:
		return mongodb.NewMongoDBContext(dataSource)
	case DBTypeCosmosDB:
		return cosmosdb.NewCosmosDBContext(dataSource)
	}

	return nil, fmt.Errorf("unknown db type %s specified", dbType)
}
