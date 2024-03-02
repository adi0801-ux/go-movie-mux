package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"go-MongoApi/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/http"
)

const connectionString = "mongodb+srv://adiac:adiacmongo123@go-mongo-mux.q18vfkf.mongodb.net/?retryWrites=true&w=majority&appName=go-mongo-mux"
const dbName = "netflix"
const colName = "watchlist"

// most important
var collection *mongo.Collection

// connect with mongodb
// runs only once in whole project and before main
func init() {
	//client option
	clientOption := options.Client().ApplyURI(connectionString)

	//connect to mongo
	client, err := mongo.Connect(context.TODO(), clientOption)
	if err != nil {
		log.Fatal(err.Error())
	}
	fmt.Println("MongoDb connection success")
	collection = client.Database(dbName).Collection(colName)

	//collection reference is ready
	fmt.Println("Collection instance is ready")
}

// mongo db helper functions
func insertOneMovie(movie model.Netflix) {
	insert, err := collection.InsertOne(context.Background(), movie)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("one movie is inserted:", insert.InsertedID)
}

func updateOneMovie(movieId string) {
	id, _ := primitive.ObjectIDFromHex(movieId)
	filter := bson.M{"_id": id}
	update := bson.M{"$set": bson.M{"watched": true}}

	result, err := collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("modified count", result.ModifiedCount)
}

func deleteOneMovie(movieId string) {
	id, _ := primitive.ObjectIDFromHex(movieId)
	filter := bson.M{"_id": id}
	deleteCount, err := collection.DeleteOne(context.Background(), filter)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("movie got deleted with delete count: ", deleteCount.DeletedCount)
}

func deleteAllMovie() int64 {
	deleteResult, err := collection.DeleteMany(context.Background(), bson.D{{}})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Number of movies deleted: ", deleteResult.DeletedCount)
	return deleteResult.DeletedCount
}

func getAllMovie() []primitive.M {
	cur, err := collection.Find(context.Background(), bson.D{{}})
	if err != nil {
		log.Fatal(err)
	}
	var movies []primitive.M
	for cur.Next(context.Background()) {
		var movie primitive.M
		err := cur.Decode(&movie)
		if err != nil {
			log.Fatal(err)
		}
		movies = append(movies, movie)
	}
	defer cur.Close(context.Background())
	return movies
}

// actual controller - file
func GetAllMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	allMovies := getAllMovie()
	json.NewEncoder(w).Encode(allMovies)
}

func CreateMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	var movie model.Netflix
	_ = json.NewDecoder(r.Body).Decode(&movie)
	insertOneMovie(movie)
	json.NewEncoder(w).Encode(movie)
}

func MarkAsWatched(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Methods", "PUT")
	params := mux.Vars(r)
	updateOneMovie(params["id"])
	json.NewEncoder(w).Encode(params["id"])
}

func DeleteAMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Methods", "DELETE")
	params := mux.Vars(r)
	deleteOneMovie(params["id"])
	json.NewEncoder(w).Encode(params["id"])
}

func DeleteAllMovies(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Methods", "DELETE")
	_ = mux.Vars(r)
	count := deleteAllMovie()
	result := map[string]int64{"Number of movies deleted are: ": count}
	json.NewEncoder(w).Encode(result)
}
