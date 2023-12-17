package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/amrohan/expenso-go/internal/db"
	"github.com/amrohan/expenso-go/internal/helpers"
	"github.com/amrohan/expenso-go/internal/models"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type Claims struct {
	Username string `json:"username"`
	Claims   jwt.RegisteredClaims
}

var jwtKey = []byte("n8&DfaT4cCJ424EEuwd")

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func Validate(u models.User) bool {
	return u.Username != "" && u.Email != "" && u.Password != ""
}

func RegisterUser(w http.ResponseWriter, r *http.Request) {
	user := models.User{
		Id:        primitive.NewObjectID(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		helpers.SendResponse(w, http.StatusBadRequest, "Please send valid body", nil, err)
		return
	}
	if !Validate(user) {
		helpers.SendResponse(w, http.StatusBadRequest, "Please send valid body", nil, nil)
		return
	}

	HashPassword, err := HashPassword(user.Password)
	if err != nil {
		helpers.SendResponse(w, http.StatusBadRequest, "Please send valid body", nil, err)
		return
	}
	user.Password = HashPassword

	client, err := db.GetMongoClient()
	if err != nil {
		helpers.SendResponse(w, http.StatusInternalServerError, "Cannot connect to database", nil, err)
		return
	}
	collection := client.Database(db.Database).Collection(string(db.UserCollection))

	// Check if the username or email already exists
	filter := bson.D{
		primitive.E{Key: "$or", Value: []bson.D{
			{{Key: "username", Value: user.Username}},
			{{Key: "email", Value: user.Email}},
		}},
	}
	var existingUser models.User
	if err := collection.FindOne(context.TODO(), filter).Decode(&existingUser); err != nil {
		if err == mongo.ErrNoDocuments {
			// No existing user with the same username or email, proceed with insertion
			if _, err := collection.InsertOne(context.TODO(), user); err != nil {
				helpers.SendResponse(w, http.StatusInternalServerError, "Unable to create user", nil, err)
				return
			}
			helpers.SendResponse(w, http.StatusOK, "User registered successfully", nil, nil)
		} else {
			helpers.SendResponse(w, http.StatusInternalServerError, "Unable to check for existing user", nil, err)
			return
		}
	} else {
		// Existing user found
		if existingUser.Username == user.Username {
			helpers.SendResponse(w, http.StatusBadRequest, "Username already exists", nil, nil)
			return
		} else if existingUser.Email == user.Email {
			helpers.SendResponse(w, http.StatusBadRequest, "Email already exists", nil, nil)
			return
		}
	}
}

func LoginUser(w http.ResponseWriter, r *http.Request) {
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		helpers.SendResponse(w, http.StatusBadRequest, "Please send valid body", nil, err)
		return
	}

	client, err := db.GetMongoClient()
	if err != nil {
		helpers.SendResponse(w, http.StatusInternalServerError, "Cannot connect to database", nil, err)
		return
	}
	collection := client.Database(db.Database).Collection(string(db.UserCollection))

	// Check if the username or email exists
	filter := bson.D{
		primitive.E{Key: "$or", Value: []bson.D{
			{{Key: "username", Value: user.Username}},
			{{Key: "email", Value: user.Email}},
		}},
	}
	var existingUser models.User
	if err := collection.FindOne(context.TODO(), filter).Decode(&existingUser); err != nil {
		if err == mongo.ErrNoDocuments {
			helpers.SendResponse(w, http.StatusNotFound, "User not found", nil, nil)
			return
		}

		helpers.SendResponse(w, http.StatusInternalServerError, "Unable to check for existing user", nil, err)

		return
	}
	if !CheckPasswordHash(user.Password, existingUser.Password) {
		helpers.SendResponse(w, http.StatusUnauthorized, "Incorrect password", nil, nil)
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username":   existingUser.Username,
		"isVerified": existingUser.IsVerified,
		"exp":        time.Now().AddDate(0, 0, 7).Unix(),
		"iat":        time.Now().Unix(),
		"nbf":        time.Now().Unix(),
		"sub":        user.Id,
		"aud":        "expenso-go",
		"iss":        "expenso-go",
	})
	tokenString, err := token.SignedString([]byte(jwtKey))
	if err != nil {
		helpers.SendResponse(w, http.StatusInternalServerError, "Error signing token", nil, err)
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:  "token",
		Value: tokenString,
	})
	helpers.SendResponse(w, http.StatusOK, "Login successful", map[string]string{"token": tokenString}, nil)

}

func LogoutUser(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:   "token",
		Value:  "",
		MaxAge: -1,
	})
	helpers.SendResponse(w, http.StatusOK, "Logout successful", nil, nil)
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("token")
		if err != nil {
			if err == http.ErrNoCookie {
				helpers.SendResponse(w, http.StatusUnauthorized, "Unauthorized", nil, nil)
				return
			}
			helpers.SendResponse(w, http.StatusBadRequest, "Bad Request", nil, err)
			return
		}
		tokenStr := cookie.Value

		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			// Validate the signing method
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			// Return the secret key
			return []byte("n8&DfaT4cCJ424EEuwd"), nil
		})

		if err != nil || !token.Valid {
			helpers.SendResponse(w, http.StatusUnauthorized, "Token has expired buddy", nil, err)
			return
		}

		_, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			helpers.SendResponse(w, http.StatusUnauthorized, "Unauthorized", nil, nil)
			return
		}
		next.ServeHTTP(w, r)
	})
}
