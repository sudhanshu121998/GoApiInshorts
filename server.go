package main

import (
    "context"
    "encoding/json"
    "fmt"

    "net/http"
    "time"
    "github.com/gorilla/mux"
    "go.mongodb.org/mongo-driver/bson"
     "go.mongodb.org/mongo-driver/mongo/options"
    "go.mongodb.org/mongo-driver/bson/primitive"
    "go.mongodb.org/mongo-driver/mongo"
)

var client *mongo.Client

type Article struct {
    ID       primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
    Title string                `json:"Title,omitempty" bson:"Title,omitempty"`
    SubTitle  string             `json:"SubTitle,omitempty" bson:"SubTitle,omitempty"`
    Content  string             `json:"Content,omitempty" bson:"Content,omitempty"`
}

func CreateArticleEndpoint(response http.ResponseWriter, request *http.Request) {
    response.Header().Set("content-type", "application/json")
    var article Article
    _ = json.NewDecoder(request.Body).Decode(&article)
    collection := client.Database("Inshort").Collection("article")
    ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
    result, _ := collection.InsertOne(ctx, article)
    json.NewEncoder(response).Encode(result)
     json.NewEncoder(response).Encode(article)
}
func GetArticleEndpoint(response http.ResponseWriter, request *http.Request) {
    response.Header().Set("content-type", "application/json")
    params := mux.Vars(request)
    id, _ := primitive.ObjectIDFromHex(params["id"])
    var article Article
    collection := client.Database("Inshort").Collection("article")
    ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
    err := collection.FindOne(ctx, Article{ID: id}).Decode(&article)
    if err != nil {
        response.WriteHeader(http.StatusInternalServerError)
        response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
        return
    }
    json.NewEncoder(response).Encode(article)
}
func GetAllArticleEndpoint(response http.ResponseWriter, request *http.Request) {
    response.Header().Set("content-type", "application/json")
    var sarticle []Article
    collection := client.Database("Inshort").Collection("article")
    ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
    cursor, err := collection.Find(ctx, bson.M{})
    if err != nil {
        response.WriteHeader(http.StatusInternalServerError)
        response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
        return
    }
    defer cursor.Close(ctx)
    for cursor.Next(ctx) {
        var article Article
        cursor.Decode(&article)
        sarticle = append(sarticle, article)
    }
    if err := cursor.Err(); err != nil {
        response.WriteHeader(http.StatusInternalServerError)
        response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
        return
    }
    json.NewEncoder(response).Encode(sarticle)
}
func main() {
    fmt.Println("Starting the application...")
    ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
    clientOptions := options.Client().ApplyURI("mongodb+srv://sudhanshu:Sudhanshu123@inshortsapi.y8p5y.mongodb.net/Inshort?retryWrites=true&w=majority")
    client, _ = mongo.Connect(ctx, clientOptions)
    router := mux.NewRouter()
    router.HandleFunc("/articles", CreateArticleEndpoint).Methods("POST")
    router.HandleFunc("/articles", GetAllArticleEndpoint).Methods("GET")
    router.HandleFunc("/articles/{id}", GetArticleEndpoint).Methods("GET")
    http.ListenAndServe(":8080", router)
}