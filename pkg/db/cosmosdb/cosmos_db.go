package cosmosdb

import (
	"context"
	"fmt"
	"os"

	"github.com/vippsas/go-cosmosdb/cosmosapi"
)

type Logger struct {

}

func (*Logger) Print(args ...interface{}) {
	fmt.Print(args...)
}
func (*Logger) Printf(template string, args ...interface{}) {
	fmt.Printf(template, args...)
}
func (*Logger) Println(args ...interface{}) {
	fmt.Println(args...)
}

type Greeting struct {
	ID string `json:"id"`
	GreetedPerson string `json:"greetedPerson"`
	Greeting      string `json:"greeting"`
}

type CosmosDBContext struct {
	client    *cosmosapi.Client
	database  string
	container string
}

func NewCosmosDBContext(dataSource string) (*CosmosDBContext, error) {
	masterkey := os.Getenv("GOGREETING_COSMOSDB_MASTERKEY")
	if masterkey == "" {
		return nil, fmt.Errorf("no masterkey specified, set GOGREETING_COSMOSDB_MASTERKEY")
	}
	client, err := getClient(dataSource, masterkey)
	if err != nil {
		return nil, err
	}

	ctx := &CosmosDBContext{
		client:    client,
		database:  "go-greetings",
		container: "greetings",
	}

	err = ctx.initDB()
	if err != nil {
		return nil, err
	}

	return ctx, nil
}

func getClient(dataSource string, masterKey string) (*cosmosapi.Client, error) {
	// Create CosmosDB client
	cosmosCfg := cosmosapi.Config{
		MasterKey: masterKey,
	}
	client := cosmosapi.New(dataSource, cosmosCfg, nil, &Logger{})

	return client, nil
}

func (dbCtx *CosmosDBContext) initDB() error {
	_, err := dbCtx.client.GetDatabase(context.Background(), dbCtx.database, nil)
	if err != nil {
		_, err := dbCtx.client.CreateDatabase(context.Background(), dbCtx.database, nil)
		if err != nil {
			return err
		}
	}

	_, err = dbCtx.client.GetCollection(context.Background(), dbCtx.database, dbCtx.container)
	if err != nil {
		colOps := cosmosapi.CreateCollectionOptions{
			Id: dbCtx.container,
			PartitionKey: &cosmosapi.PartitionKey{Paths: []string{"/greetedPerson"}, Kind: "Hash"},
			IndexingPolicy: &cosmosapi.IndexingPolicy{
				IndexingMode: "consistent",
				Automatic: true,
				Included: []cosmosapi.IncludedPath{{Path: "/*"}},
				Excluded: []cosmosapi.ExcludedPath{{Path: "/\"_etag\"/?"}},
			},
		}
		_, err = dbCtx.client.CreateCollection(context.Background(), dbCtx.database, colOps)
		if err != nil {
			return err
		}
	}

	return nil
}

func (dbCtx *CosmosDBContext) SaveGreeting(greetedPerson string, greeting string) error {
	greetingDoc := Greeting{ID: greetedPerson, GreetedPerson: greetedPerson, Greeting: greeting}

	ops := cosmosapi.CreateDocumentOptions{
		PartitionKeyValue: greetedPerson,
		IsUpsert:          true,
	}
	_, _, err := dbCtx.client.CreateDocument(context.Background(), dbCtx.database, dbCtx.container, &greetingDoc, ops)

	return err
}

func (dbCtx *CosmosDBContext) GetGreeting(greetedPerson string) (*string, error) {

	qops := cosmosapi.DefaultQueryDocumentOptions()
	qops.PartitionKeyValue = greetedPerson
	query := cosmosapi.Query{
		Query: "SELECT * FROM greetings g WHERE g.greetedPerson = @greetedPerson OFFSET 0 LIMIT 1",
		Params: []cosmosapi.QueryParam{
			{Name: "@greetedPerson", Value: greetedPerson},
		},
	}

	docs := []Greeting{}

	resp, err := dbCtx.client.QueryDocuments(context.Background(), dbCtx.database, dbCtx.container, query, &docs, qops)
	if err != nil {
		return nil, err
	}

	if resp.Count == 0 {
		return nil, fmt.Errorf("Could not find greeting for %s", greetedPerson)
	}

	return &docs[0].Greeting, nil
}
