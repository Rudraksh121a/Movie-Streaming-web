package controllers

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/Rudraksh121a/Movie-Streaming-web/database"
	"github.com/Rudraksh121a/Movie-Streaming-web/models"
	"github.com/Rudraksh121a/Movie-Streaming-web/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"github.com/tmc/langchaingo/llms/googleai"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var validate = validator.New()
var movieCollection *mongo.Collection = database.OpenCollection("Movies")

func GetMovies() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var movies []models.Movie
		cursor, err := movieCollection.Find(ctx, bson.M{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch movies"})
			return
		}
		defer cursor.Close(ctx)

		if err = cursor.All(ctx, &movies); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to decode movies"})
			return
		}

		c.JSON(http.StatusOK, movies)
	}
}

func GetMovie() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		movieId := c.Param("imdb_id")
		if movieId == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Movie Id is Required"})
			return
		}
		var movie models.Movie

		err := movieCollection.FindOne(ctx, bson.M{"imdb_id": movieId}).Decode(&movie)
		if err != nil {

			c.JSON(http.StatusNotFound, gin.H{"error": "Movie Not Found"})
		}
		c.JSON(http.StatusOK, movie)

	}
}
func AddMovie() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var movie models.Movie
		if err := c.ShouldBindBodyWithJSON(&movie); err != nil {

			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
			return
		}
		if err := validate.Struct(movie); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "validation failed", "detail": err.Error()})
			return
		}

		result, err := movieCollection.InsertOne(ctx, movie)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to add movie"})
			return
		}
		c.JSON(http.StatusCreated, result)
	}
}

func AdminReviewUpdate() gin.HandlerFunc {
	return func(c *gin.Context) {

		role, err := utils.GetRoleFromContext(c)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Role not found in context"})
			return
		}
		if role != "ADMIN" {
			c.JSON(http.StatusForbidden, gin.H{"error": "Only admin can update review"})
			return
		}

		movieId := c.Param("imdb_id")
		if movieId == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Movie Id is Required"})
			return
		}
		var req struct {
			AdminReview string `json:"admin_review"`
		}
		var resp struct {
			RankingName string `json:"ranking_name"`
			AdminReview string `json:"admin_review"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}
		sentiment, rankVal, err := GetReviewRanking(req.AdminReview)
		if err != nil {
			log.Println("Error in GetReviewRanking:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to review ranking", "details": err.Error()})
			return
		}
		filter := bson.D{{Key: "imdb_id", Value: movieId}}

		update := bson.M{
			"$set": bson.M{
				"admin_review": req.AdminReview,
				"ranking": bson.M{
					"ranking_value": rankVal,
					"ranking_name":  sentiment,
				},
			},
		}
		var ctx, cancel = context.WithTimeout(c, 100*time.Second)
		defer cancel()

		var movieCollection *mongo.Collection = database.OpenCollection("Movies")

		result, err := movieCollection.UpdateOne(ctx, filter, update)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating movie"})
			return
		}

		if result.MatchedCount == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "Movie not found"})
			return
		}
		resp.RankingName = sentiment
		resp.AdminReview = req.AdminReview
		c.JSON(http.StatusOK, resp)
	}
}
func GetReviewRanking(adminReview string) (string, int, error) {
	rankings, err := GetRanking()
	if err != nil {
		return "", 0, err
	}
	sentimentDelimited := ""
	for _, ranking := range rankings {
		if ranking.RankingValue != 999 {

			sentimentDelimited += ranking.RankingName + ", "
		}
	}
	sentimentDelimited = strings.Trim(sentimentDelimited, ",")
	err = godotenv.Load()
	if err != nil {
		log.Println("Warning env not found")

	}
	GEMINI_API_KEY := os.Getenv("GEMINI_API_KEY")
	if GEMINI_API_KEY == "" {
		return "", 0, errors.New("GEMINI API key not set in environment variables")
	}
	ctx := context.Background()
	llm, err := googleai.New(ctx, googleai.WithAPIKey(GEMINI_API_KEY), googleai.WithDefaultModel("gemini-2.5-flash-lite"))
	if err != nil {
		return "", 0, err
	}
	basePromptTemplate := os.Getenv("BASE_PROMT_TEMPLATE")
	base_prompt := strings.Replace(basePromptTemplate, "{rankings}", sentimentDelimited, 1)
	response, err := llm.Call(ctx, base_prompt+adminReview)
	if err != nil {
		return "", 0, err
	}
	rankval := 0
	for _, ranking := range rankings {
		if ranking.RankingName == response {
			rankval = ranking.RankingValue
			break
		}
	}
	return response, rankval, nil
}

func GetRanking() ([]models.Ranking, error) {
	var rankings []models.Ranking

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	rankingCollection := database.OpenCollection("Rankings")
	cursor, err := rankingCollection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if err = cursor.All(ctx, &rankings); err != nil {
		return nil, err
	}
	return rankings, nil
}

func GetRecommendedMovies() gin.HandlerFunc {
	return func(c *gin.Context) {
		userid, err := utils.GetUserIdFromContext(c)

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found in context"})
			return
		}
		favourite_genres, err := GetUsersFavouriteGenres(userid)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if err != nil {
			log.Println("waring: .env file not found")
		}
		var recommendedMoviesLimitVal int64 = 5
		recommendedMoviesLimitStr := os.Getenv("RECOMMENDED_MOVIES_LIMIT")
		if recommendedMoviesLimitStr != "" {
			recommendedMoviesLimitVal, _ = strconv.ParseInt(recommendedMoviesLimitStr, 10, 64)
		}
		findOptions := options.Find()

		findOptions.SetSort(bson.D{{Key: "ranking.ranking_value", Value: 1}})

		findOptions.SetLimit(recommendedMoviesLimitVal)

		filter := bson.D{
			{Key: "genre.genre_name", Value: bson.D{
				{Key: "$in", Value: favourite_genres},
			}},
		}

		var ctx, cancel = context.WithTimeout(c, 100*time.Second)
		defer cancel()

		cursor, err := movieCollection.Find(ctx, filter, findOptions)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching recommended movies"})
			return
		}
		defer cursor.Close(ctx)
		var recommendedMovies []models.Movie

		if err := cursor.All(ctx, &recommendedMovies); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, recommendedMovies)
	}
}

func GetUsersFavouriteGenres(userid string) ([]string, error) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	filter := bson.D{{Key: "user_id", Value: userid}}

	projection := bson.M{
		"favourite_genres.genre_name": 1,
		"_id":                         0,
	}

	opts := options.FindOne().SetProjection(projection)
	var result bson.M
	err := userCollection.FindOne(ctx, filter, opts).Decode(&result)
	if err == mongo.ErrNoDocuments {
		return []string{}, nil
	}
	favGenresArray, ok := result["favourite_genres"].(bson.A)

	if !ok {
		return []string{}, errors.New("unable to retrieve favourite genres for user")
	}

	var genreNames []string

	for _, item := range favGenresArray {
		if genreMap, ok := item.(bson.D); ok {
			for _, elem := range genreMap {
				if elem.Key == "genre_name" {
					if name, ok := elem.Value.(string); ok {
						genreNames = append(genreNames, name)
					}
				}
			}
		}
	}

	return genreNames, nil

}
