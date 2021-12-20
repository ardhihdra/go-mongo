package main

import (
	"context"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var (
	MONGO_HOST string
	MONGO_PORT string
	credential = options.Credential{
		AuthMechanism: "SCRAM-SHA-256",
		AuthSource:    "<authenticationDb>",
		Username:      "<username>",
		Password:      "<password>",
	}
)

type Mongodb struct {
	Client *mongo.Client
}

func main() {
	godotenv.Load()
	MONGO_HOST = os.Getenv("MONGODB_HOST")
	MONGO_PORT = os.Getenv("MONGODB_PORT")

	// Connection URI
	// "mongodb://user:pass@sample.host:27017/?maxPoolSize=20&w=majority"
	uri := fmt.Sprintf("mongodb://%s:%s/?maxPoolSize=20&w=majority", MONGO_HOST, MONGO_PORT)

	// Create a new client and connect to the server
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))

	if err != nil {
		panic(err)
	}
	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()
	mongo := Mongodb{client}
	mongo.create()
	mongo.read()
	// Ping the primary
	if err := client.Ping(context.TODO(), readpref.Primary()); err != nil {
		panic(err)
	}

	fmt.Println("Successfully connected and pinged.")
}

func (mc *Mongodb) create() {
	col := mc.Client.Database("tes").Collection("connectionstatus")
	doc := bson.D{{"ip", "172.16.1.121"}, {"well_id", "BTG-P1"}, {"target_mqtt", "mqtt://172.16.1.200"}}

	result, _ := col.InsertOne(context.TODO(), doc)
	fmt.Printf("Inserted document with return: %v\n", result)
}

func (mc *Mongodb) read() {
	col := mc.Client.Database("tes").Collection("connectionstatus")
	filter := bson.D{}
	sort := bson.D{{"well_id", -1}}
	projection := bson.D{{"well_id", "BTG-P1"}, {"target_mqtt", "mqtt://172.16.1.200"}}
	opts := options.FindOne().SetSort(sort).SetProjection(projection)
	var result bson.D
	err := col.FindOne(context.TODO(), filter, opts).Decode(&result)
	if err != nil {
		panic(err)
	}
	fmt.Println(result)
}

func (mc *Mongodb) update() {

}

func (mc *Mongodb) delete() {

}
