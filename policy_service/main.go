package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func main() {

	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		mongoURI = "mongodb://127.0.0.1:27017"
	}
	client, err := mongo.Connect(options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(context.TODO())

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal("could not ping mongo:", err)
	}
	fmt.Println("connected to mongodb!")

	collection := client.Database("policy_as_code").Collection("policies")

	indexModel := mongo.IndexModel{
		Keys:    bson.D{{Key: "name", Value: 1}},
		Options: options.Index().SetUnique(true),
	}
	_, err = collection.Indexes().CreateOne(context.Background(), indexModel)
	if err != nil {
		log.Fatal("failed to create index:", err)
	}

	repo := &PolicyRepository{collection: collection, indexModel: &indexModel}
	handler := &PolicyHandler{repo: repo}

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	r.Mount("/policies", PolicyRoutes(handler))

	http.ListenAndServe(":3000", r)

}

func PolicyRoutes(policyHandler *PolicyHandler) chi.Router {

	r := chi.NewRouter()
	fmt.Println("Got here")
	r.Get("/", policyHandler.ListPolicies)
	r.Post("/", policyHandler.CreatePolicy)
	r.Get("/{name}", policyHandler.GetPolicies)
	r.Put("/{name}", policyHandler.UpdatePolicy)
	r.Delete("/{name}", policyHandler.DeletePolicy)

	return r
}
