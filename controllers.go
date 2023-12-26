package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type user struct {
	Name string `json:"name"`
	City string `json:"city"`
	Age  int    `json:"age"`
}

var userCol = db().Database("gotest").Collection("user")

func createUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var person user
	err := json.NewDecoder(r.Body).Decode(&person)
	if err != nil {
		fmt.Print(err)
	}
	insertResult, err := userCol.InsertOne(context.TODO(), person)
	if err != nil {
		log.Fatal(err)
	}
	json.NewEncoder(w).Encode(insertResult.InsertedID)
}

func getUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var body user
	err := json.NewDecoder(r.Body).Decode(&body)

	if err != nil {
		fmt.Print(err)
	}

	var results []primitive.M

	data, er := userCol.Find(context.TODO(), bson.D{{"city", body.City}})

	if er != nil {
		fmt.Println(er)
	}

	for data.Next(context.TODO()) {
		var result primitive.M
		e := data.Decode(&result)

		if e != nil {
			fmt.Println(e)
			http.Error(w, e.Error(), http.StatusInternalServerError)
			return
		}
		results = append(results, result)

	}

	json.NewEncoder(w).Encode(results)

}

func updateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	type reqBody struct {
		Name string `json:"name"`
		City string `json:"city"`
	}
	var body reqBody
	err := json.NewDecoder(r.Body).Decode(&body)

	if err != nil {
		fmt.Println(err)
	}

	filter := bson.D{{"name", body.Name}}
	after := options.After

	dataOpt := options.FindOneAndUpdateOptions{
		ReturnDocument: &after,
	}

	update := bson.D{{"$set", bson.D{{"city", body.City}}}}

	updatedData := userCol.FindOneAndUpdate(context.TODO(), filter, update, &dataOpt)

	var result primitive.M

	er := updatedData.Decode(&result)

	if er != nil {
		http.Error(w, er.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(result)
}

func deleteUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)["id"]

	_id, err := primitive.ObjectIDFromHex(params)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	opts := options.Delete().SetCollation(&options.Collation{})

	filter := bson.D{{"_id", _id}}

	res, er := userCol.DeleteOne(context.TODO(), filter, opts)

	if er != nil {
		http.Error(w, er.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(res.DeletedCount)

}
