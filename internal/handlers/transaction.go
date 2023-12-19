package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/amrohan/expenso-go/internal/db"
	"github.com/amrohan/expenso-go/internal/helpers"
	"github.com/amrohan/expenso-go/internal/models"
	"github.com/go-chi/chi/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func CreateTransaction(w http.ResponseWriter, r *http.Request) {
	transaction := models.Transaction{
		Id:        primitive.NewObjectID(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := json.NewDecoder(r.Body).Decode(&transaction); err != nil {
		helpers.SendResponse(w, http.StatusBadRequest, "Couldnt decode request", nil, err)
	}
	transaction.Id = primitive.NewObjectID()

	client, err := db.GetMongoClient()
	if err != nil {
		helpers.SendResponse(w, http.StatusInternalServerError, "Couldnt connect to db", nil, err)
	}

	collection := client.Database(db.Database).Collection(string(db.TransactionCollection))

	data, err := collection.InsertOne(context.TODO(), transaction)
	if err != nil {
		helpers.SendResponse(w, http.StatusInternalServerError, "Couldnt insert transaction", nil, err)
	}
	helpers.SendResponse(w, http.StatusOK, "Transaction created", data, nil)
}

func GetAllTransaction(w http.ResponseWriter, r *http.Request) {
	var transactions []*models.Transaction

	client, err := db.GetMongoClient()
	if err != nil {
		helpers.SendResponse(w, http.StatusInternalServerError, "Couldnt connect to db", nil, err)
	}
	collection := client.Database(db.Database).Collection(string(db.TransactionCollection))

	cur, err := collection.Find(context.TODO(), bson.D{
		primitive.E{},
	})

	if err != nil {
		helpers.SendResponse(w, http.StatusInternalServerError, "Couldnt find transactions", nil, err)
	}

	for cur.Next(context.TODO()) {
		var transaction models.Transaction
		err := cur.Decode(&transaction)
		if err != nil {
			helpers.SendResponse(w, http.StatusInternalServerError, "Couldnt decode transactions", nil, err)
		}
		transactions = append(transactions, &transaction)
	}
	helpers.SendResponse(w, http.StatusOK, "Transactions fetched successfully", transactions, nil)
}

func GetTransactionById(w http.ResponseWriter, r *http.Request) {
	var transaction models.Transaction

	id, err := primitive.ObjectIDFromHex(chi.URLParam(r, "id"))
	if err != nil {
		helpers.SendResponse(w, http.StatusBadRequest, "Invalid id", nil, err)
	}

	client, err := db.GetMongoClient()
	if err != nil {
		helpers.SendResponse(w, http.StatusInternalServerError, "Couldnt connect to db", nil, err)
	}
	collection := client.Database(db.Database).Collection(string(db.TransactionCollection))

	err = collection.FindOne(context.TODO(), bson.M{"_id": id}).Decode(&transaction)
	if err != nil {
		helpers.SendResponse(w, http.StatusInternalServerError, "Couldnt find transaction", nil, err)
	}
	helpers.SendResponse(w, http.StatusOK, "Transaction found", transaction, nil)
}

func GetTransactionByUserId(w http.ResponseWriter, r *http.Request) {
	var transactions []*models.Transaction
	userId := chi.URLParam(r, "id")

	client, err := db.GetMongoClient()
	if err != nil {
		helpers.SendResponse(w, http.StatusInternalServerError, "Couldnt connect to db", nil, err)
	}
	collection := client.Database(db.Database).Collection(string(db.TransactionCollection))

	cur, err := collection.Find(context.TODO(), bson.M{"userId": userId})
	if err != nil {
		helpers.SendResponse(w, http.StatusInternalServerError, "Couldnt find transactions", nil, err)
	}

	for cur.Next(context.TODO()) {
		var transaction models.Transaction
		err := cur.Decode(&transaction)
		if err != nil {
			helpers.SendResponse(w, http.StatusInternalServerError, "Couldnt decode transactions", nil, err)
		}

		transactions = append(transactions, &transaction)
	}

	helpers.SendResponse(w, http.StatusOK, "Transactions found", transactions, nil)
}

func GetTransactionByAccountId(w http.ResponseWriter, r *http.Request) {
	var transactions []*models.Transaction
	accountId := chi.URLParam(r, "id")
	client, err := db.GetMongoClient()
	if err != nil {
		helpers.SendResponse(w, http.StatusInternalServerError, "Couldnt connect to db", nil, err)
	}
	collection := client.Database(db.Database).Collection(string(db.TransactionCollection))

	cur, err := collection.Find(context.TODO(), bson.M{"accountId": accountId})
	if err != nil {
		helpers.SendResponse(w, http.StatusInternalServerError, "Couldnt find transactions", nil, err)
	}

	for cur.Next(context.TODO()) {
		var transaction models.Transaction
		err := cur.Decode(&transaction)
		if err != nil {
			helpers.SendResponse(w, http.StatusInternalServerError, "Couldnt decode transactions", nil, err)
		}

		transactions = append(transactions, &transaction)
	}
	helpers.SendResponse(w, http.StatusOK, "Transactions found", transactions, nil)
}

func GetTransactionByCategoryId(w http.ResponseWriter, r *http.Request) {
	var transactions []*models.Transaction
	categoryId := chi.URLParam(r, "id")

	client, err := db.GetMongoClient()
	if err != nil {
		helpers.SendResponse(w, http.StatusInternalServerError, "Couldnt connect to db", nil, err)
	}
	collection := client.Database(db.Database).Collection(string(db.TransactionCollection))

	cur, err := collection.Find(context.TODO(), bson.M{"categoryId": categoryId})
	if err != nil {
		helpers.SendResponse(w, http.StatusInternalServerError, "Couldnt find transactions", nil, err)
	}

	for cur.Next(context.TODO()) {
		var transaction models.Transaction
		err := cur.Decode(&transaction)
		if err != nil {
			helpers.SendResponse(w, http.StatusInternalServerError, "Couldnt decode transactions", nil, err)
		}

		transactions = append(transactions, &transaction)
	}
	helpers.SendResponse(w, http.StatusOK, "Transactions found", transactions, nil)
}

func GetTransactionByMonthAndYear(w http.ResponseWriter, r *http.Request) {
	var transactions []*models.Transaction

	monthStr := chi.URLParam(r, "month")
	yearStr := chi.URLParam(r, "year")

	month, err := strconv.Atoi(monthStr)
	if err != nil {
		helpers.SendResponse(w, http.StatusBadRequest, "Please send valid month", nil, err)
		return
	}

	year, err := strconv.Atoi(yearStr)
	if err != nil {
		helpers.SendResponse(w, http.StatusBadRequest, "Please send valid year", nil, err)
		return
	}

	startDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	endDate := startDate.AddDate(0, 1, 0)

	client, err := db.GetMongoClient()
	if err != nil {
		helpers.SendResponse(w, http.StatusInternalServerError, "Couldnt connect to db", nil, err)
		return
	}
	collection := client.Database(db.Database).Collection(string(db.TransactionCollection))

	filter := bson.M{
		"date": bson.M{
			"$gte": startDate,
			"$lt":  endDate,
		},
	}

	cur, err := collection.Find(context.TODO(), filter)
	if err != nil {
		helpers.SendResponse(w, http.StatusInternalServerError, "Couldnt find transactions", nil, err)
		return
	}

	for cur.Next(context.TODO()) {
		var transaction models.Transaction
		err := cur.Decode(&transaction)
		if err != nil {
			helpers.SendResponse(w, http.StatusInternalServerError, "Couldnt decode transactions", nil, err)
			return
		}

		transactions = append(transactions, &transaction)
	}
	if transactions == nil {
		helpers.SendResponse(w, http.StatusOK, "No transactions found", transactions, nil)
		return
	}
	helpers.SendResponse(w, http.StatusOK, "Transactions found", transactions, nil)
}

func GetTransactionByMonthAndYearByUserId(w http.ResponseWriter, r *http.Request) {
	var transactions []*models.Transaction

	monthStr := chi.URLParam(r, "month")
	yearStr := chi.URLParam(r, "year")
	userId := chi.URLParam(r, "userId")

	month, err := strconv.Atoi(monthStr)
	if err != nil {
		helpers.SendResponse(w, http.StatusBadRequest, "Please send valid month", nil, err)
		return
	}

	year, err := strconv.Atoi(yearStr)
	if err != nil {
		helpers.SendResponse(w, http.StatusBadRequest, "Please send valid year", nil, err)
		return
	}

	startDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	endDate := startDate.AddDate(0, 1, 0)

	client, err := db.GetMongoClient()
	if err != nil {
		helpers.SendResponse(w, http.StatusInternalServerError, "Couldnt connect to db", nil, err)
		return
	}
	collection := client.Database(db.Database).Collection(string(db.TransactionCollection))
	filter := bson.M{
		"date": bson.M{
			"$gte": startDate,
			"$lt":  endDate,
		},
		"userId": userId,
	}
	cur, err := collection.Find(context.TODO(), filter)
	if err != nil {
		helpers.SendResponse(w, http.StatusInternalServerError, "Couldnt find transactions", nil, err)
		return
	}
	for cur.Next(context.TODO()) {
		var transaction models.Transaction
		err := cur.Decode(&transaction)
		if err != nil {
			helpers.SendResponse(w, http.StatusInternalServerError, "Couldnt decode transactions", nil, err)
			return
		}

		transactions = append(transactions, &transaction)
	}
	if transactions == nil {
		helpers.SendResponse(w, http.StatusOK, "No transactions found", transactions, nil)
		return
	}
	helpers.SendResponse(w, http.StatusOK, "Transactions found", transactions, nil)
}

func UpdateTransaction(w http.ResponseWriter, r *http.Request) {
	var transaction models.Transaction

	if err := json.NewDecoder(r.Body).Decode(&transaction); err != nil {
		helpers.SendResponse(w, http.StatusBadRequest, "Couldnt decode request", nil, err)
	}

	client, err := db.GetMongoClient()
	if err != nil {
		helpers.SendResponse(w, http.StatusInternalServerError, "Couldnt connect to db", nil, err)
	}
	collection := client.Database(db.Database).Collection(string(db.TransactionCollection))

	data, err := collection.UpdateOne(context.TODO(), bson.M{"_id": transaction.Id}, bson.M{"$set": transaction})
	if err != nil {
		helpers.SendResponse(w, http.StatusInternalServerError, "Couldnt update transaction", nil, err)
	}
	helpers.SendResponse(w, http.StatusOK, "Transaction updated", data, nil)
}

func DeleteTransaction(w http.ResponseWriter, r *http.Request) {
	id, err := primitive.ObjectIDFromHex(chi.URLParam(r, "id"))
	if err != nil {
		helpers.SendResponse(w, http.StatusBadRequest, "Invalid id", nil, err)
	}

	client, err := db.GetMongoClient()
	if err != nil {
		helpers.SendResponse(w, http.StatusInternalServerError, "Couldnt connect to db", nil, err)
	}
	collection := client.Database(db.Database).Collection(string(db.TransactionCollection))

	data, err := collection.DeleteOne(context.TODO(), bson.M{"_id": id})
	if err != nil {
		helpers.SendResponse(w, http.StatusInternalServerError, "Couldnt delete transaction", nil, err)
	}
	helpers.SendResponse(w, http.StatusOK, "Transaction deleted", data, nil)
}
