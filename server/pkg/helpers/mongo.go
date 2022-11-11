package helpers

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Repository struct {
	client *mongo.Client
	dbName string
}

func NewStorage(connectionString, dbName string) (*Repository, error) {
	clientOptions := options.Client().ApplyURI(connectionString)

	ctx := context.TODO()
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}

	// Check connection.
	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}

	repo := &Repository{
		client: client,
		dbName: dbName,
	}

	return repo, nil
}

func (r Repository) InsertImage(img string) error {
	coll := r.client.Database(r.dbName).Collection("images")

	doc := bson.M{"image": img}

	_, err := coll.InsertOne(context.TODO(), doc)
	if err != nil {
		return err
	}

	return nil
}
