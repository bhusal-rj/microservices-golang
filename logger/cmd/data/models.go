package data

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

func New(mongo *mongo.Client) Models {
	client = mongo
	return Models{
		LogEntry: LogEntry(),
	}
}

type Models struct {
	LogEntry LogEntry
}

type LogEntry struct {
	ID        string    `bson:"_id,omitempty"	json:"id,omitempty"`
	Name      string    `bson:"name" json:"name"`
	Data      string    `bson:"data" json:"data"`
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
}

func (l *LogEntry) Insert(entry LogEntry) error {
	collection := client.Database("logs").Collection("logger")
	_, err := collection.InsertOne(context.TODO(), LogEntry{
		Name:      entry.Name,
		Data:      entry.Data,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})

	if err != nil {
		log.Panic("Error inserting the data", err)
		return err
	}
	return nil

}

func (l *LogEntry) All() ([]*LogEntry, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	collection := client.Database("logs").Collection("logs")
	opts := options.Find()
	cursor, err := collection.Find(ctx, bson.D{}, opts)
	if err != nil {
		log.Panic("Error retrivig all log entries", err)
		return nil, err
	}
	var logEntries []*LogEntry
	for cursor.Next(ctx) {
		var item LogEntry
		err := cursor.Decode(&item)
		if err != nil {
			log.Print("Error decoding the slice:", err)
			return nil, err
		} else {
			logEntries = append(logEntries, &item)
		}
	}
	return logEntries, nil
}
