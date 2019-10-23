package main

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Trainer struct {
	Name string
	Age  int
	City string
}

type Order struct {
	Userid   int
	Good     int
	Price    int
	Counter  int
	Totalsum int
	Status   int // 0 - created, 1 - comfimed and aviable, 2 - not aviable at this moment, 11 - comfied by user, 12 - canceled by user
	Todate   string
	Totime   string
	Isdone   bool
}

func UpdateOrder(guid string, s Order) int {
	fmt.Println("Start UpdateOrder function with guid =", guid)
	fmt.Print("update for guid=", guid)
	fmt.Println("user =", s.Userid)
	fmt.Println("good =", s.Good)
	fmt.Println("price =", s.Price)
	fmt.Println("counter =", s.Counter)
	fmt.Println("totalsum =", s.Totalsum)
	fmt.Println("status =", s.Status)
	fmt.Println("todate =", s.Todate)
	fmt.Println("totime =", s.Totime)
	fmt.Println("isdone =", s.Isdone)
	return 0
}

func CreateOrder(s Order) string {
	fmt.Println("Start CreateOrder function")
	fmt.Println("user =", s.Userid)
	fmt.Println("good =", s.Good)
	fmt.Println("price =", s.Price)
	fmt.Println("counter =", s.Counter)
	fmt.Println("totalsum =", s.Totalsum)
	fmt.Println("status =", s.Status)
	fmt.Println("todate =", s.Todate)
	fmt.Println("totime =", s.Totime)
	fmt.Println("isdone =", s.Isdone)

	// Insert a single document
	insertResult, err := collection.InsertOne(context.TODO(), ash)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Inserted a single document: ", insertResult.InsertedID)

	return ""
}

func main() {

	// Set client options
	clientOptions := options.client().ApplyURI("mongodb://localhost:27017")

	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!")

	// Get a handle for your collection
	//collection := client.Database("test").Collection("trainers")
	collection := client.Database("test").Collection("orders")

	t1 := Order{Userid: 1, Good: 1,
		Price: 120, Counter: 1,
		Totalsum: 120,
		Status:   0, // 0 - created, 1 - comfimed and aviable, 2 - not aviable at this moment, 11 - comfied by user, 12 - canceled by user
		Todate:   "26-10-2019",
		Totime:   "с 18 до 19 часов",
		Isdone:   false}

	err = CreateOrder(t1)
	if err != nil {
		log.Fatal("Not created order", err)
	}
	// Close the connection once no longer needed
	err = client.Disconnect(context.TODO())

	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("Connection to MongoDB closed.")
	}

	return
	// Some dummy data to add to the Database
	ash := Trainer{"Ash", 10, "Pallet Town"}
	misty := Trainer{"Misty", 10, "Cerulean City"}
	brock := Trainer{"Brock", 15, "Pewter City"}

	// Insert a single document
	insertResult, err := collection.InsertOne(context.TODO(), ash)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Inserted a single document: ", insertResult.InsertedID)

	// Insert multiple documents
	trainers := []interface{}{misty, brock}

	insertManyResult, err := collection.InsertMany(context.TODO(), trainers)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Inserted multiple documents: ", insertManyResult.InsertedIDs)

	// Update a document
	filter := bson.D{{"name", "Ash"}}

	update := bson.D{
		{"$inc", bson.D{
			{"age", 1},
		}},
	}

	updateResult, err := collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Matched %v documents and updated %v documents.\n", updateResult.MatchedCount, updateResult.ModifiedCount)

	// Find a single document
	var result Trainer

	err = collection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Found a single document: %+v\n", result)

	findOptions := options.Find()
	findOptions.SetLimit(2)

	var results []*Trainer

	// Finding multiple documents returns a cursor
	cur, err := collection.Find(context.TODO(), bson.D{{}}, findOptions)
	if err != nil {
		log.Fatal(err)
	}

	// Iterate through the cursor
	for cur.Next(context.TODO()) {
		var elem Trainer
		err := cur.Decode(&elem)
		if err != nil {
			log.Fatal(err)
		}

		results = append(results, &elem)
	}

	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}

	// Close the cursor once finished
	cur.Close(context.TODO())

	fmt.Printf("Found multiple documents (array of pointers): %+v\n", results)

	// Delete all the documents in the collection
	deleteResult, err := collection.DeleteMany(context.TODO(), bson.D{{}})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Deleted %v documents in the trainers collection\n", deleteResult.DeletedCount)

	// Close the connection once no longer needed
	err = client.Disconnect(context.TODO())

	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("Connection to MongoDB closed.")
	}

}
