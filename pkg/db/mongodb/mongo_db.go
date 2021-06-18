package mongodb

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Greeting struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	GreetedPerson string `bson:"greetedPerson"`
	Greeting string `bson:"greeting"`
}

type MongoDBContext struct {
	client *mongo.Client
	database string
	collection string
}

func NewMongoDBContext(dataSource string) (*MongoDBContext, error) {
	client, err := getClient(dataSource)
	if err != nil {
		return nil, err
	}

	ctx := &MongoDBContext{
		client: client,
		database: "go-greetings",
		collection: "greetings",
	}

	err = ctx.initDB()
	if err != nil {
		return nil, err
	}
	
	return ctx, nil
}

func getClient(dataSource string) (*mongo.Client, error) {
	var err error
	clientOptions := options.Client().ApplyURI(dataSource).SetDirect(true)
	client, err := mongo.NewClient(clientOptions)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	if err = client.Connect(ctx); err != nil {
		return nil, err
	}
	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("could not ping mongodb: %v", err)
	}

	return client, nil
}

func (dbCtx *MongoDBContext) initDB() error {
	return nil
}

func (dbCtx *MongoDBContext) SaveGreeting(greetedPerson string, greeting string) error {
	
	greetings := dbCtx.client.Database(dbCtx.database).Collection(dbCtx.collection)
	filter := bson.M{
		"greetedPerson": greetedPerson,
	}
	upsert := true
	opts := options.UpdateOptions{
		Upsert: &upsert,
	}
	greetingDoc := Greeting{GreetedPerson: greetedPerson, Greeting: greeting}
	update := bson.M{
        "$set": greetingDoc,
    }
	_, err := greetings.UpdateOne(context.Background(), filter, update, &opts)
	
	return err
}

func (dbCtx *MongoDBContext) GetGreeting(greetedPerson string) (*string, error) {
	greetings := dbCtx.client.Database(dbCtx.database).Collection(dbCtx.collection)
	filter := bson.M{
		"greetedPerson": greetedPerson,
	}
	result := greetings.FindOne(context.Background(), filter)
	if result.Err() != nil {
		return nil, result.Err()
	}

	greetingDoc := Greeting{}
	result.Decode(&greetingDoc)

	return &greetingDoc.Greeting, nil
}
