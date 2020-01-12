package gorest

import (
	"context"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"time"
)

type MongoDB struct {
	database *mongo.Database
	log      *log.Entry
}

var g_MongoDB *MongoDB

func GetMongoDB() *MongoDB {
	return g_MongoDB
}

func NewMongoDB() *MongoDB {
	g_MongoDB = &MongoDB{
		database: nil,
		log: log.WithFields(log.Fields{
			"key": "db",
		}),
	}

	if !g_MongoDB.init() {
		g_MongoDB.log.Errorln("Unable to connect to mongodb")
		return nil
	}

	return g_MongoDB
}

func (mongoDB *MongoDB) init() bool {
	var err error
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	mongoDBUrl := Get("mongodb.url").(string)
	if mongoDBUrl == "" {
		mongoDBUrl = "mongodb://localhost:27017"
		mongoDB.log.Infoln("No configuration for mongodb's url found, use the default one: %s", mongoDBUrl)
	}
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoDBUrl))
	if err != nil {
		mongoDB.log.Fatalln("Error to init mongodb. Error: %s", err)
		return false
	}

	ctx, _ = context.WithTimeout(context.Background(), 2*time.Second)
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		mongoDB.log.Fatalln("Error to ping mongodb. Error: %s", err)
		return false
	}

	mongoDB.log.Infoln("Mongodb is init!")
	databaseName := Get("mongodb.database").(string)
	if databaseName == "" {
		databaseName = "gorest"
		mongoDB.log.Infoln("No configuration for mongodb database found, use the default one: %s", databaseName)
	}
	g_MongoDB.database = client.Database(databaseName)
	mongoDB.log.Infoln("Mongodb is connected to gorest database")
	return true
}

func (mongoDB *MongoDB) Close() bool {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err := mongoDB.database.Client().Disconnect(ctx)

	if err != nil {
		mongoDB.log.Fatalln("Error to close mongodb connection. Error: %s", err)
		return false
	}

	mongoDB.log.Infoln("Connection to MongoDB closed.")
	return true
}

func (mongoDB *MongoDB) InsertOne(ctx context.Context, collectionName string, document interface{}) (interface{}, error) {
	collection := mongoDB.database.Collection(collectionName)
	res, err := collection.InsertOne(ctx, document)
	if err != nil {
		mongoDB.log.Errorln("Error to insert %v in %s. Error: %s", document, collectionName, err)
		return nil, err
	}

	return res.InsertedID, nil
}

func (mongoDB *MongoDB) InsertMany(ctx context.Context, collectionName string, documents []interface{}) ([]interface{}, error) {
	collection := mongoDB.database.Collection(collectionName)
	res, err := collection.InsertMany(ctx, documents)
	if err != nil {
		mongoDB.log.Errorln("Error to insert many %v in %s. Error: %s", documents, collectionName, err)
		return nil, err
	}

	return res.InsertedIDs, nil
}

func (mongoDB *MongoDB) UpdateOne(ctx context.Context, collectionName string, update bson.M, filter bson.M) (int64, error) {
	collection := mongoDB.database.Collection(collectionName)
	res, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		mongoDB.log.Errorln("Error to update one for %s collection. Error %s", collectionName, err)
		return 0, err
	}

	return res.ModifiedCount, nil
}

func (mongoDB *MongoDB) UpdateMany(ctx context.Context, collectionName string, update bson.M, filter bson.M) (int64, error) {
	collection := mongoDB.database.Collection(collectionName)
	res, err := collection.UpdateMany(ctx, filter, update)
	if err != nil {
		mongoDB.log.Errorln("Error to update many for %s collection. Error %s", collectionName, err)
		return 0, err
	}

	return res.ModifiedCount, nil
}

func (mongoDB *MongoDB) Find(ctx context.Context, collectionName string, filter bson.M) ([]interface{}, error) {
	collection := mongoDB.database.Collection(collectionName)
	cur, err := collection.Find(ctx, filter)
	if err != nil {
		mongoDB.log.Errorln("Error to find in %s collection for %v filter. Error: %s", collectionName, filter, err)
		return nil, err
	}
	defer func() {
		err := cur.Close(ctx)
		if err != nil {
			mongoDB.log.Errorln("Error to closed cursor. Error: %s", err)
		}
	}()

	var res []interface{}
	for cur.Next(ctx) {
		var result bson.M
		err := cur.Decode(&result)
		if err != nil {
			mongoDB.log.Errorln("Error to decode result. Error: %s", err)
			return nil, err
		}
		res = append(res, result)
	}

	if err := cur.Err(); err != nil {
		mongoDB.log.Errorln("Error in cursor. Error: %s", err)
		return nil, err
	}

	return res, nil
}

func (mongoDB *MongoDB) FindOne(ctx context.Context, collectionName string, filter bson.M) (interface{}, error) {
	collection := mongoDB.database.Collection(collectionName)
	var result interface{}
	err := collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		mongoDB.log.Errorln("Error to find one in %s collection for %s filter. Error: %s", collectionName, filter, err)
		return nil, err
	}

	return result, nil
}

func (mongoDB *MongoDB) DeleteOne(ctx context.Context, collectionName string, filter bson.M) (int64, error) {
	collection := mongoDB.database.Collection(collectionName)
	res, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		mongoDB.log.Errorln("Error to delete one %s in %s collection. Error: %s", filter, collectionName, err)
		return 0, err
	}

	return res.DeletedCount, nil
}

func (mongoDB *MongoDB) DeleteMany(ctx context.Context, collectionName string, filter bson.M) (int64, error) {
	collection := mongoDB.database.Collection(collectionName)
	res, err := collection.DeleteMany(ctx, filter)
	if err != nil {
		mongoDB.log.Errorln("Error to delete many %s in %s collection. Error: %s", filter, collectionName, err)
		return 0, err
	}

	return res.DeletedCount, nil
}
