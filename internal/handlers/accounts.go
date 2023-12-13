package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/amrohan/expenso-go/internal/db"
	"github.com/amrohan/expenso-go/internal/helpers"
	"github.com/amrohan/expenso-go/internal/models"
	"github.com/go-chi/chi/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func CreateAccount(w http.ResponseWriter, r *http.Request) {
	account := models.Account{
		Id:        primitive.NewObjectID(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := json.NewDecoder(r.Body).Decode(&account); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	client, err := db.GetMongoClient()
	if err != nil {
		helpers.SendResponse(w, http.StatusInternalServerError, "Couldnt connect to db", nil, err)
	}
	collection := client.Database(db.Database).Collection(string(db.AccountCollection))

	data, err := collection.InsertOne(context.TODO(), account)
	if err != nil {
		helpers.SendResponse(w, http.StatusInternalServerError, "Couldnt create account", nil, err)
	}
	helpers.SendResponse(w, http.StatusOK, "Account created", data, nil)
}

func GetAllAccount(w http.ResponseWriter, r *http.Request) {
	client, err := db.GetMongoClient()
	if err != nil {
		helpers.SendResponse(w, http.StatusInternalServerError, "Error connecting to database", nil, err)
		return
	}
	collection := client.Database(db.Database).Collection(string(db.AccountCollection))
	var accounts []*models.Account
	cur, err := collection.Find(context.TODO(), bson.D{
		primitive.E{},
	})
	if err != nil {
		helpers.SendResponse(w, http.StatusInternalServerError, "Error finding accounts", nil, err)
		return
	}
	for cur.Next(context.TODO()) {
		var account models.Account

		err := cur.Decode(&account)
		if err != nil {
			helpers.SendResponse(w, http.StatusInternalServerError, "Error decoding account", nil, err)
			return
		}
		accounts = append(accounts, &account)
	}
	helpers.SendResponse(w, http.StatusOK, "Accounts found", accounts, nil)
}

func GetAccountById(w http.ResponseWriter, r *http.Request) {
	var account models.Account
	id, err := primitive.ObjectIDFromHex(chi.URLParam(r, "id"))
	if err != nil {
		helpers.SendResponse(w, http.StatusBadRequest, "Invalid id", nil, err)
	}

	client, err := db.GetMongoClient()
	if err != nil {
		helpers.SendResponse(w, http.StatusInternalServerError, "Couldnt connect to db", nil, err)
	}
	collection := client.Database(db.Database).Collection(string(db.TransactionCollection))

	err = collection.FindOne(context.TODO(), bson.M{"_id": id}).Decode(&account)
	if err != nil {
		helpers.SendResponse(w, http.StatusInternalServerError, "Couldnt find transaction", nil, err)
	}
	helpers.SendResponse(w, http.StatusOK, "Transaction found", account, nil)
}

func GetAccountsByUserId(w http.ResponseWriter, r *http.Request) {
	var accounts []*models.Account
	id := chi.URLParam(r, "id")
	client, err := db.GetMongoClient()
	if err != nil {
		helpers.SendResponse(w, http.StatusInternalServerError, "Couldnt connect to db", nil, err)
	}
	collection := client.Database(db.Database).Collection(string(db.AccountCollection))
	cur, err := collection.Find(context.TODO(), bson.M{"user_id": id})
	if err != nil {
		helpers.SendResponse(w, http.StatusInternalServerError, "Couldnt find transaction", nil, err)
		return
	}
	for cur.Next(context.TODO()) {
		var account models.Account
		if err := cur.Decode(&account); err != nil {
			helpers.SendResponse(w, http.StatusInternalServerError, "Error decoding account", nil, err)
		}
		accounts = append(accounts, &account)
	}
	helpers.SendResponse(w, http.StatusOK, "Accounts found", accounts, nil)
}

func UpdateAccount(w http.ResponseWriter, r *http.Request) {
	account := models.Account{
		UpdatedAt: time.Now(),
	}

	if err := json.NewDecoder(r.Body).Decode(&account); err != nil {
		helpers.SendResponse(w, http.StatusBadRequest, "Please send valid json", nil, err)
		return
	}

	client, err := db.GetMongoClient()
	if err != nil {
		helpers.SendResponse(w, http.StatusInternalServerError, "Couldnt connect to db", nil, err)
	}
	collection := client.Database(db.Database).Collection(string(db.AccountCollection))

	data, err := collection.UpdateOne(r.Context(), bson.M{"id": account.Id}, bson.M{"$set": account})
	if err != nil {
		helpers.SendResponse(w, http.StatusInternalServerError, "Error updating account", nil, err)
		return
	}
	helpers.SendResponse(w, http.StatusOK, "Account updated successfully", data, nil)
}

func DeleteAccount(w http.ResponseWriter, r *http.Request) {
	id, err := primitive.ObjectIDFromHex(chi.URLParam(r, "id"))
	if err != nil {
		helpers.SendResponse(w, http.StatusBadRequest, "Invalid id", nil, err)
	}
	client, err := db.GetMongoClient()
	if err != nil {
		helpers.SendResponse(w, http.StatusInternalServerError, "Couldnt connect to db", nil, err)
	}
	collection := client.Database(db.Database).Collection(string(db.AccountCollection))

	data, err := collection.DeleteOne(r.Context(), bson.M{"id": id})
	if err != nil {
		helpers.SendResponse(w, http.StatusInternalServerError, "Error deleting account", nil, err)
		return
	}
	helpers.SendResponse(w, http.StatusOK, "Account deleted successfully", data, nil)
}
