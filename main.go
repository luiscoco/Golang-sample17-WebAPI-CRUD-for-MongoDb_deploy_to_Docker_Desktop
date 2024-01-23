package main

import (
    "context"
    "encoding/json"
    "net/http"
    "github.com/gorilla/mux"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
    "go.mongodb.org/mongo-driver/bson/primitive"
    "log"
)

// Book struct to map with MongoDB documents
type Book struct {
    ID     primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
    Title  string             `json:"title,omitempty" bson:"title,omitempty"`
    Author string             `json:"author,omitempty" bson:"author,omitempty"`
    ISBN   string             `json:"isbn,omitempty" bson:"isbn,omitempty"`
}

var collection *mongo.Collection

func main() {
    connectDB()
    r := mux.NewRouter()

    r.HandleFunc("/books", createBook).Methods("POST")
    r.HandleFunc("/books", getBooks).Methods("GET")
    r.HandleFunc("/books/{id}", getBook).Methods("GET")
    r.HandleFunc("/books/{id}", updateBook).Methods("PUT")
    r.HandleFunc("/books/{id}", deleteBook).Methods("DELETE")

    log.Fatal(http.ListenAndServe(":8080", r))
}

func connectDB() {
    // Change the URI to "mongodb://mongodb:27017" to connect to the MongoDB container
    clientOptions := options.Client().ApplyURI("mongodb://mongodb:27017")
    client, err := mongo.Connect(context.TODO(), clientOptions)
    if err != nil {
        log.Fatal(err)
    }

    collection = client.Database("bookstore").Collection("books")
}

func createBook(w http.ResponseWriter, r *http.Request) {
    var book Book
    json.NewDecoder(r.Body).Decode(&book)
    
    result, err := collection.InsertOne(context.TODO(), book)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    json.NewEncoder(w).Encode(result)
}

func getBooks(w http.ResponseWriter, r *http.Request) {
    var books []Book
    cursor, err := collection.Find(context.TODO(), bson.M{})
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    defer cursor.Close(context.TODO())

    for cursor.Next(context.TODO()) {
        var book Book
        cursor.Decode(&book)
        books = append(books, book)
    }

    json.NewEncoder(w).Encode(books)
}

func getBook(w http.ResponseWriter, r *http.Request) {
    id, err := primitive.ObjectIDFromHex(mux.Vars(r)["id"])
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    
    var book Book
    err = collection.FindOne(context.TODO(), bson.M{"_id": id}).Decode(&book)
    if err != nil {
        http.Error(w, err.Error(), http.StatusNotFound)
        return
    }

    json.NewEncoder(w).Encode(book)
}

func updateBook(w http.ResponseWriter, r *http.Request) {
    id, err := primitive.ObjectIDFromHex(mux.Vars(r)["id"])
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    
    var book Book
    json.NewDecoder(r.Body).Decode(&book)

    _, err = collection.UpdateOne(
        context.TODO(),
        bson.M{"_id": id},
        bson.D{
            {"$set", bson.D{
                {"title", book.Title},
                {"author", book.Author},
                {"isbn", book.ISBN},
            }},
        },
    )
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    json.NewEncoder(w).Encode(book)
}

func deleteBook(w http.ResponseWriter, r *http.Request) {
    id, err := primitive.ObjectIDFromHex(mux.Vars(r)["id"])
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    
    _, err = collection.DeleteOne(context.TODO(), bson.M{"_id": id})
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusNoContent)
}
