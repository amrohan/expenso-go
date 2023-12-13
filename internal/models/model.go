package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Response struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type Transaction struct {
	Id         primitive.ObjectID `json:"id" bson:"_id"`
	Title      string             `json:"title" bson:"title"`
	Amount     int                `json:"amount" bson:"amount"`
	Date       time.Time          `json:"date" bson:"date"`
	CategoryId string             `json:"categoryId" bson:"categoryId"`
	Type       string             `json:"type" bson:"type"`
	ImageUrl   string             `json:"imageUrl" bson:"imageUrl"`
	AccountId  string             `json:"accountId" bson:"accountId"`
	UserId     string             `json:"userId" bson:"userId"`
	CreatedAt  time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAt  time.Time          `json:"updatedAt" bson:"updatedAt"`
	DeletedAt  time.Time          `json:"deletedAt" bson:"deletedAt"`
	IsDeleted  bool               `json:"isDeleted" bson:"isDeleted"`
	IsActive   bool               `json:"isActive" bson:"isActive"`
}

type Category struct {
	Id        primitive.ObjectID `json:"id" bson:"_id"`
	Title     string             `json:"title" bson:"title"`
	Icon      string             `json:"imageUrl" bson:"imageUrl"`
	UserId    string             `json:"userId" bson:"userId"`
	CreatedAt time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time          `json:"updatedAt" bson:"updatedAt"`
	DeletedAt time.Time          `json:"deletedAt" bson:"deletedAt"`
	IsDeleted bool               `json:"isDeleted" bson:"isDeleted"`
	IsActive  bool               `json:"isActive" bson:"isActive"`
}

type Account struct {
	Id        primitive.ObjectID `json:"id" bson:"_id"`
	Title     string             `json:"title" bson:"title"`
	Icon      string             `json:"imageUrl" bson:"imageUrl"`
	UserId    string             `json:"userId" bson:"userId"`
	CreatedAt time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time          `json:"updatedAt" bson:"updatedAt"`
	DeletedAt time.Time          `json:"deletedAt" bson:"deletedAt"`
	IsDeleted bool               `json:"isDeleted" bson:"isDeleted"`
	IsActive  bool               `json:"isActive" bson:"isActive"`
}

type User struct {
	Id         primitive.ObjectID `json:"id" bson:"_id"`
	Username   string             `json:"username" bson:"username"`
	Name       string             `json:"name" bson:"name"`
	Email      string             `json:"email" bson:"email"`
	Password   string             `json:"password" bson:"password"`
	ImageUrl   string             `json:"imageUrl" bson:"imageUrl"`
	CreatedAt  time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAt  time.Time          `json:"updatedAt" bson:"updatedAt"`
	DeletedAt  time.Time          `json:"deletedAt" bson:"deletedAt"`
	IsDeleted  bool               `json:"isDeleted" bson:"isDeleted"`
	IsActive   bool               `json:"isActive" bson:"isActive"`
	IsVerified bool               `json:"isVerified" bson:"isVerified"`
}
