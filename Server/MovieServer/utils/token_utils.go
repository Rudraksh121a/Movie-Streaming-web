package utils

import (
	"context"
	"os"
	"time"

	"github.com/Rudraksh121a/Movie-Streaming-web/database"
	jwt "github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type SignedDetails struct {
	Email     string
	FirstName string
	LastName  string
	Role      string
	UserId    string
	jwt.RegisteredClaims
}

var SECRET_KEY string = os.Getenv("SECRET_KEY")
var REFRESH_TOKEN_SECRET string = os.Getenv("REFRESH_TOKEN_SECRET")
var userCollection *mongo.Collection = database.OpenCollection("Users")

func GenerateAllTokens(email, firstName, lastName, role, userId string) (string, string, error) {
	clams := &SignedDetails{
		Email:     email,
		FirstName: firstName,
		LastName:  lastName,
		Role:      role,
		UserId:    userId,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "Stream",
			IssuedAt:  &jwt.NumericDate{time.Now()},
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, clams)
	signedToken, err := token.SignedString([]byte(SECRET_KEY))
	if err != nil {
		return "", "", nil
	}

	refreshClams := &SignedDetails{
		Email:     email,
		FirstName: firstName,
		LastName:  lastName,
		Role:      role,
		UserId:    userId,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "Stream",
			IssuedAt:  &jwt.NumericDate{time.Now()},
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClams)
	signedRefreshToken, err := refreshToken.SignedString([]byte(REFRESH_TOKEN_SECRET))
	if err != nil {
		return "", "", nil
	}
	return signedToken, signedRefreshToken, nil
}

func UpdateAllToken(userId, token, refreshToken string) (err error) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	updateAt, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	updateData := bson.M{
		"$set": bson.M{
			"token":         token,
			"refresh_token": refreshToken,
			"update_at":     updateAt,
		},
	}
	_, err = userCollection.UpdateOne(ctx, bson.M{"user_id": userId}, updateData)
	if err != nil {
		return err
	}
	return nil
}
