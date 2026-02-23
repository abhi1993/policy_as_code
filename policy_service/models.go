package main

import (
	"context"
	"errors"
	"fmt"
	"os"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type Policy struct {
	ID     string `bson:"_id"`
	Name   string `json:"name" bson:"name"`
	Author string `json:"author" bson:"author"`
	Date   string `json:"published_date" bson:"date"`
	Code   string `json:"code" bson:"code"`
}

type PolicyRepository struct {
	collection *mongo.Collection
	indexModel *mongo.IndexModel
}

func (repo PolicyRepository) listPolicies(token string) ([]Policy, string, error) {
	var policies []Policy
	var cursor *mongo.Cursor
	var err error

	fmt.Printf("Got here into list policies with token: %s", token)
	opts := options.Find().SetLimit(10).SetSort(bson.D{{Key: "_id", Value: 1}})
	if token == "" {
		cursor, err = repo.collection.Find(context.Background(), bson.D{}, opts)
	} else {
		filter := bson.D{{Key: "_id", Value: bson.D{{Key: "$gt", Value: token}}}}
		cursor, err = repo.collection.Find(context.Background(), filter, opts)
	}
	fmt.Println("Got here into listpolicies 3")
	if err != nil {
		return nil, "", err
	}
	fmt.Println("Got here into listpolicies 3.5")
	err = cursor.All(context.Background(), &policies)
	fmt.Printf("Got here into listpolicies 3.6 %s \n", err)
	if err != nil {
		return nil, "", err
	}

	fmt.Println("printing all policies to return 1:", policies)
	var nextToken string
	if len(policies) > 0 {
		fmt.Println("printing all policies to return 2:", policies)
		nextToken = policies[len(policies)-1].ID
	}
	fmt.Println("Got here into listpolicies 4")
	return policies, nextToken, nil
}

func (repo PolicyRepository) getPolicyById(id string) (Policy, error) {
	ret_policy := Policy{}
	returned_object := repo.collection.FindOne(context.Background(), bson.M{"_id": id})
	if returned_object.Err() != nil {
		fmt.Printf("Mongodb call resulted in error: %s", returned_object.Err())
		return Policy{}, errors.New(returned_object.Err().Error())
	} else {
		fmt.Printf("Object returned is %s", returned_object)
		returned_object.Decode(&ret_policy)
		return ret_policy, nil

	}

}

func (repo PolicyRepository) getPolicyByName(name string) (Policy, error) {
	ret_policy := Policy{}
	returned_object := repo.collection.FindOne(context.Background(), bson.M{"name": name})
	if returned_object.Err() != nil {
		fmt.Printf("Mongodb call resulted in error: %s", returned_object.Err())
		return Policy{}, errors.New(returned_object.Err().Error())
	} else {
		fmt.Printf("Object returned is %s", returned_object)
		returned_object.Decode(&ret_policy)
		return ret_policy, nil
	}
}

func (repo PolicyRepository) addPolicy(pol Policy) error {
	id := pol.ID
	fmt.Printf("adding policy %s to database", id)
	x, err := repo.collection.InsertOne(context.Background(), pol)

	if err != nil {
		fmt.Println("error inserting policy:", err)
		return err
	} else {
		fmt.Println("inserted policy id ", x.InsertedID)
		return nil
	}
}

func (repo PolicyRepository) updatePolicy(pol Policy) error {
	name := pol.Name

	// Step 1: find by name to get the ID
	existing := Policy{}
	err := repo.collection.FindOne(context.Background(), bson.M{"name": pol.Name}).Decode(&existing)
	if err != nil {
		return err
	}

	// Step 2: update by ID, replacing all fields
	pol.ID = existing.ID // carry over the original ID
	filter := bson.D{{Key: "_id", Value: existing.ID}}
	update := bson.D{{Key: "$set", Value: pol}}
	fmt.Printf("for update found name: %s and generated filter %s and update %s\n", name, filter, update)
	resp, err := repo.collection.UpdateOne(context.Background(), filter, update)

	fmt.Println("made update call")
	if err != nil {
		fmt.Fprintln(os.Stderr, "error updating policy:", err)
		return err
	} else if resp.MatchedCount == 0 {
		fmt.Fprintln(os.Stderr, "error updating policy: no documents matched")
		return errors.New("error updating policy: no documents matched")
	} else {
		fmt.Printf("policy successfully acked:%s matched count:%s modified count:%s \n", resp.Acknowledged, resp.MatchedCount, resp.MatchedCount)
		return err
	}

}

func (repo PolicyRepository) deletePolicy(name string) error {
	filter := bson.D{{Key: "name", Value: name}}

	resp, err := repo.collection.DeleteOne(context.Background(), filter)

	if err != nil {
		fmt.Println("error deleting policy:", err)
		return err
	} else if resp.DeletedCount == 0 {
		fmt.Fprintln(os.Stderr, "error deleting policy: no documents matched")
		return errors.New("error deleting policy: no documents matched")
	} else {
		fmt.Printf("policy successfully deleted %s\n", resp)
		return nil
	}

}
