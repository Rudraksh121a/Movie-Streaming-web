package utils

import (
	"context"
	"errors"
	"os"
	"time"

	"github.com/Rudraksh121a/Movie-Streaming-web/database"
	"github.com/gin-gonic/gin"
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

func GetAccessToken(c *gin.Context) (string, error) {
	authHeader := c.Request.Header.Get("Authorization")
	if authHeader == "" {
		return "", errors.New("Authorization header missing")
	}
	tokenString := authHeader[len("Bearer "):]
	if tokenString == "" {
		return "", errors.New("Bearer token is required")
	}
	return tokenString, nil
}

func ValidateToken(tokenString string) (*SignedDetails, error) {
	clams := &SignedDetails{}
	token, err := jwt.ParseWithClaims(tokenString, clams, func(t *jwt.Token) (interface{}, error) {
		return []byte(SECRET_KEY), nil
	})
	if err != nil {
		return nil, err
	}

	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, err
	}
	if clams.ExpiresAt.Time.Before(time.Now()) {
		return nil, errors.New("Token expired")
	}
	return clams, nil
}

func GetUserIdFromContext(c *gin.Context) (string, error) {
	userId, exists := c.Get("user_Id")
	if !exists {
		return "", errors.New("userid not found in context")
	}
	id, ok := userId.(string)
	if !ok {
		return "", errors.New("userid in context is not a string")
	}
	return id, nil
}

func GetRoleFromContext(c *gin.Context) (string, error) {
	role, exists := c.Get("role")
	if !exists {
		return "", errors.New("role not found in context")
	}
	memberRole, ok := role.(string)
	if !ok {
		return "", errors.New("role in context is not a string")
	}
	return memberRole, nil
}
