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

func CreateCategory(w http.ResponseWriter, r *http.Request) {
	category := models.Category{
		Id:        primitive.NewObjectID(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := json.NewDecoder(r.Body).Decode(&category); err != nil {
		helpers.SendResponse(w, http.StatusBadRequest, "Please send valid json", nil, err)
		return
	}

	client, err := db.GetMongoClient()
	if err != nil {
		helpers.SendResponse(w, http.StatusInternalServerError, "Error connecting to database", nil, err)
		return
	}

	collection := client.Database(db.Database).Collection(string(db.CategoryCollection))

	if _, err := collection.InsertOne(r.Context(), category); err != nil {
		helpers.SendResponse(w, http.StatusInternalServerError, "Error creating category", nil, err)
		return
	}
	helpers.SendResponse(w, http.StatusCreated, "Category created successfully", nil, nil)
}

func GetAllCategory(w http.ResponseWriter, r *http.Request) {
	client, err := db.GetMongoClient()
	if err != nil {
		helpers.SendResponse(w, http.StatusInternalServerError, "Error connecting to database", nil, err)
		return
	}

	collection := client.Database(db.Database).Collection(string(db.CategoryCollection))

	var categories []models.Category

	cur, err := collection.Find(context.TODO(), bson.D{
		primitive.E{},
	})

	if err != nil {
		helpers.SendResponse(w, http.StatusInternalServerError, "Error fetching categories", nil, err)
		return
	}
	for cur.Next(r.Context()) {
		var category *models.Category
		if err := cur.Decode(&category); err != nil {
			helpers.SendResponse(w, http.StatusInternalServerError, "Error decoding category", nil, err)
			return
		}
		categories = append(categories, *category)
	}
	helpers.SendResponse(w, http.StatusOK, "Categories fetched successfully", categories, nil)
}

func GetCategoryById(w http.ResponseWriter, r *http.Request) {
	id, err := primitive.ObjectIDFromHex(chi.URLParam(r, "id"))
	if err != nil {
		helpers.SendResponse(w, http.StatusBadRequest, "Invalid id", nil, err)
	}

	client, err := db.GetMongoClient()
	if err != nil {
		helpers.SendResponse(w, http.StatusInternalServerError, "Error connecting to database", nil, err)
		return
	}

	collection := client.Database(db.Database).Collection(string(db.CategoryCollection))

	var category *models.Category

	if err := collection.FindOne(r.Context(), bson.M{"id": id}).Decode(&category); err != nil {
		helpers.SendResponse(w, http.StatusInternalServerError, "Error fetching category", nil, err)
		return
	}
	helpers.SendResponse(w, http.StatusOK, "Category fetched successfully", category, nil)
}

func GetCategoryByUserId(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		helpers.SendResponse(w, http.StatusBadRequest, "Please provide a valid id", nil, nil)
		return
	}

	client, err := db.GetMongoClient()
	if err != nil {
		helpers.SendResponse(w, http.StatusInternalServerError, "Error connecting to database", nil, err)
		return
	}

	collection := client.Database(db.Database).Collection(string(db.CategoryCollection))

	var categories []*models.Category

	cur, err := collection.Find(context.TODO(), bson.M{"user_id": id})

	if err != nil {
		helpers.SendResponse(w, http.StatusInternalServerError, "Error fetching categories", nil, err)
		return
	}
	for cur.Next(r.Context()) {
		var category models.Category
		if err := cur.Decode(&category); err != nil {
			helpers.SendResponse(w, http.StatusInternalServerError, "Error decoding category", nil, err)
			return
		}
		categories = append(categories, &category)
	}
	helpers.SendResponse(w, http.StatusOK, "Categories fetched successfully", categories, nil)
}

func UpdateCategory(w http.ResponseWriter, r *http.Request) {
	category := models.Category{
		UpdatedAt: time.Now(),
	}

	if err := json.NewDecoder(r.Body).Decode(&category); err != nil {
		helpers.SendResponse(w, http.StatusBadRequest, "Please send valid json", nil, err)
		return
	}

	client, err := db.GetMongoClient()
	if err != nil {
		helpers.SendResponse(w, http.StatusInternalServerError, "Error connecting to database", nil, err)
		return
	}
	collection := client.Database(db.Database).Collection(string(db.CategoryCollection))

	data, err := collection.UpdateOne(r.Context(), bson.M{"id": category.Id}, bson.M{"$set": category})
	if err != nil {
		helpers.SendResponse(w, http.StatusInternalServerError, "Error updating category", nil, err)
		return
	}
	helpers.SendResponse(w, http.StatusOK, "Category updated successfully", data, nil)
}

func DeleteCategory(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		helpers.SendResponse(w, http.StatusBadRequest, "Please provide a valid id", nil, nil)
		return
	}

	client, err := db.GetMongoClient()
	if err != nil {
		helpers.SendResponse(w, http.StatusInternalServerError, "Error connecting to database", nil, err)
		return
	}

	collection := client.Database(db.Database).Collection(string(db.CategoryCollection))

	if _, err := collection.DeleteOne(r.Context(), bson.M{"id": id}); err != nil {
		helpers.SendResponse(w, http.StatusInternalServerError, "Error deleting category", nil, err)
		return
	}
	helpers.SendResponse(w, http.StatusOK, "Category deleted successfully", nil, nil)
}
