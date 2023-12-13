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

func CreateUser(w http.ResponseWriter, r *http.Request) {
	user := models.User{
		Id:        primitive.NewObjectID(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		helpers.SendResponse(w, http.StatusBadRequest, "Please valid body", nil, err)
		return
	}

	client, err := db.GetMongoClient()
	if err != nil {
		helpers.SendResponse(w, http.StatusInternalServerError, "Error connecting to mongo", nil, err)
		return
	}

	collection := client.Database(db.Database).Collection(string(db.UserCollection))
	if _, err = collection.InsertOne(context.TODO(), user); err != nil {
		helpers.SendResponse(w, http.StatusInternalServerError, "Error creating user", nil, err)
		return
	}
	helpers.SendResponse(w, http.StatusOK, "User created successfully", user, nil)
}

func GetAllUsers(w http.ResponseWriter, r *http.Request) {
	client, err := db.GetMongoClient()
	if err != nil {
		helpers.SendResponse(w, http.StatusInternalServerError, "Error connecting to mongo", nil, err)
		return
	}

	collection := client.Database(db.Database).Collection(string(db.UserCollection))
	cursor, err := collection.Find(context.TODO(), bson.D{
		primitive.E{},
	})
	if err != nil {
		helpers.SendResponse(w, http.StatusInternalServerError, "Error getting users", nil, err)
		return
	}

	var users []models.User
	if err = cursor.All(context.TODO(), &users); err != nil {
		helpers.SendResponse(w, http.StatusInternalServerError, "Error getting users", nil, err)
		return
	}
	helpers.SendResponse(w, http.StatusOK, "Users fetched successfully", users, nil)
}

func GetUserById(w http.ResponseWriter, r *http.Request) {
	id, err := primitive.ObjectIDFromHex(chi.URLParam(r, "id"))
	if err != nil {
		helpers.SendResponse(w, http.StatusBadRequest, "Invalid user id", nil, err)
		return
	}

	client, err := db.GetMongoClient()
	if err != nil {
		helpers.SendResponse(w, http.StatusInternalServerError, "Error connecting to mongo", nil, err)
		return
	}

	collection := client.Database(db.Database).Collection(string(db.UserCollection))
	var user models.User
	if err = collection.FindOne(context.TODO(), bson.M{"_id": id}).Decode(&user); err != nil {
		helpers.SendResponse(w, http.StatusInternalServerError, "Error getting user", nil, err)
		return
	}
	helpers.SendResponse(w, http.StatusOK, "User fetched successfully", user, nil)
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	user := models.User{
		UpdatedAt: time.Now(),
	}

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		helpers.SendResponse(w, http.StatusBadRequest, "Please valid body", nil, err)
		return
	}
	client, err := db.GetMongoClient()
	if err != nil {
		helpers.SendResponse(w, http.StatusInternalServerError, "Error connecting to mongo", nil, err)
		return
	}
	collection := client.Database(db.Database).Collection(string(db.UserCollection))
	_, err = collection.UpdateOne(context.TODO(), bson.M{"_id": user.Id}, bson.M{"$set": user})
	if err != nil {
		helpers.SendResponse(w, http.StatusInternalServerError, "Error updating user", nil, err)
		return
	}
	helpers.SendResponse(w, http.StatusOK, "User updated successfully", user, nil)
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	id, err := primitive.ObjectIDFromHex(chi.URLParam(r, "id"))
	if err != nil {
		helpers.SendResponse(w, http.StatusBadRequest, "Invalid user id", nil, err)
		return
	}

	client, err := db.GetMongoClient()
	if err != nil {
		helpers.SendResponse(w, http.StatusInternalServerError, "Error connecting to mongo", nil, err)
		return
	}

	collection := client.Database(db.Database).Collection(string(db.UserCollection))
	if _, err = collection.DeleteOne(context.TODO(), bson.M{"_id": id}); err != nil {
		helpers.SendResponse(w, http.StatusInternalServerError, "Error deleting user", nil, err)
		return
	}
	helpers.SendResponse(w, http.StatusOK, "User deleted successfully", nil, nil)
}

func GetUserByEmail(w http.ResponseWriter, r *http.Request) {
	email := chi.URLParam(r, "email")

	client, err := db.GetMongoClient()
	if err != nil {
		helpers.SendResponse(w, http.StatusInternalServerError, "Error connecting to mongo", nil, err)
		return
	}

	collection := client.Database(db.Database).Collection(string(db.UserCollection))
	var user models.User
	if err = collection.FindOne(context.TODO(), bson.M{"email": email}).Decode(&user); err != nil {
		helpers.SendResponse(w, http.StatusInternalServerError, "Error getting user", nil, err)
		return
	}
	helpers.SendResponse(w, http.StatusOK, "User fetched successfully", user, nil)
}

func GetUserByUsername(w http.ResponseWriter, r *http.Request) {
	username := chi.URLParam(r, "username")

	client, err := db.GetMongoClient()
	if err != nil {
		helpers.SendResponse(w, http.StatusInternalServerError, "Error connecting to mongo", nil, err)
		return
	}

	collection := client.Database(db.Database).Collection(string(db.UserCollection))
	var user models.User
	if err = collection.FindOne(context.TODO(), bson.M{"username": username}).Decode(&user); err != nil {
		helpers.SendResponse(w, http.StatusInternalServerError, "Error getting user", nil, err)
		return
	}
	helpers.SendResponse(w, http.StatusOK, "User fetched successfully", user, nil)
}

func GetAllDeletedUser(w http.ResponseWriter, r *http.Request) {
	client, err := db.GetMongoClient()
	if err != nil {
		helpers.SendResponse(w, http.StatusInternalServerError, "Error connecting to mongo", nil, err)
		return
	}

	collection := client.Database(db.Database).Collection(string(db.UserCollection))
	cursor, err := collection.Find(context.TODO(), bson.D{
		primitive.E{Key: "isDeleted", Value: true},
	})
	if err != nil {
		helpers.SendResponse(w, http.StatusInternalServerError, "Error getting users", nil, err)
		return
	}

	var users []models.User
	if err = cursor.All(context.TODO(), &users); err != nil {
		helpers.SendResponse(w, http.StatusInternalServerError, "Error getting users", nil, err)
		return
	}
	helpers.SendResponse(w, http.StatusOK, "Users fetched successfully", users, nil)
}

func RestoreUser(w http.ResponseWriter, r *http.Request) {
	id, err := primitive.ObjectIDFromHex(chi.URLParam(r, "id"))
	if err != nil {
		helpers.SendResponse(w, http.StatusBadRequest, "Invalid user id", nil, err)
		return
	}

	client, err := db.GetMongoClient()
	if err != nil {
		helpers.SendResponse(w, http.StatusInternalServerError, "Error connecting to mongo", nil, err)
		return
	}

	collection := client.Database(db.Database).Collection(string(db.UserCollection))
	if _, err = collection.UpdateOne(context.TODO(), bson.M{"_id": id}, bson.M{"$set": bson.M{"isDeleted": false}}); err != nil {
		helpers.SendResponse(w, http.StatusInternalServerError, "Error restoring user", nil, err)
		return
	}
	helpers.SendResponse(w, http.StatusOK, "User restored successfully", nil, nil)
}
