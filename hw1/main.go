package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type Book struct {
	Title     string
	Author    string
	ISBN      string
	Publisher string
	Year      int
	Copies    int
}

func (b *Book) addCustomBook() {
	fmt.Print("Enter title: ")
	fmt.Scanln(&b.Title)
	fmt.Print("Enter author: ")
	fmt.Scanln(&b.Author)
	fmt.Print("Enter ISBN: ")
	fmt.Scanln(&b.ISBN)
	fmt.Print("Enter publisher: ")
	fmt.Scanln(&b.Publisher)
	fmt.Print("Enter year: ")
	fmt.Scanln(&b.Year)
	fmt.Print("Enter copies: ")
	fmt.Scanln(&b.Copies)
}

func main() {
	ctx := context.TODO()
	opts := options.Client().ApplyURI("mongodb://localhost:27017")

	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)

	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		log.Fatal(err)
	}

	LibraryDB := client.Database("Library")
	BooksCollection := LibraryDB.Collection("Books")
	defer BooksCollection.Drop(ctx)

	var customBookselecter byte
	var custom Book

	fmt.Println("Do you want to add a custom book? (y/n): ")
	for {
		fmt.Scanf("%c\n", &customBookselecter)

		if customBookselecter == 'y' || customBookselecter == 'n' {
			break
		} else {

			fmt.Println("Please enter (y/n): ")
			time.Sleep(time.Millisecond)

		}
	}
	if customBookselecter == 'y' {
		custom.addCustomBook()
		bookDoc := bson.D{
			{Key: "title", Value: custom.Title},
			{Key: "author", Value: custom.Author},
			{Key: "isbn", Value: custom.ISBN},
			{Key: "publisher", Value: custom.Publisher},
			{Key: "year", Value: custom.Year},
			{Key: "copies", Value: custom.Copies},
		}
		r, err := BooksCollection.InsertOne(ctx, bookDoc)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Inserted book with ID:", r.InsertedID)
	}

	example := bson.D{
		{Key: "title", Value: "War and Peace"},
		{Key: "author", Value: "Leo Tolstoy"},
		{Key: "isbn", Value: "12412"},
		{Key: "publisher", Value: "Moscow"},
		{Key: "year", Value: 1867},
		{Key: "copies", Value: 426000},
	}
	fmt.Println()
	fmt.Println("START OF INSERT ONE")
	r, err := BooksCollection.InsertOne(ctx, example)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Inserted document with ID:", r.InsertedID)
	fmt.Println("END OF INSERT ONE")
	fmt.Println()

	examples := []interface{}{
		bson.D{
			{Key: "title", Value: "1984"},
			{Key: "author", Value: "George Orwell"},
			{Key: "isbn", Value: "978214214"},
			{Key: "publisher", Value: "London"},
			{Key: "year", Value: 1949},
			{Key: "copies", Value: 8},
		},
		bson.D{
			{Key: "title", Value: "Captain's Daughter"},
			{Key: "author", Value: "Alexander Pushkin"},
			{Key: "isbn", Value: "978214231214"},
			{Key: "publisher", Value: "Moscow"},
			{Key: "year", Value: 1836},
			{Key: "copies", Value: 800000},
		},
	}
	fmt.Println()
	fmt.Println("START OF INSERT MANY")
	rs, err := BooksCollection.InsertMany(ctx, examples)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Inserted documents with IDs:", rs.InsertedIDs)
	fmt.Println("END OF INSERT MANY")
	fmt.Println()
	fmt.Println("STARTS OF FIND ONE")
	c := BooksCollection.FindOne(ctx, bson.M{"_id": r.InsertedID})
	var exampleResult Book
	err = c.Decode(&exampleResult)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			fmt.Println("No document found with ID:", r.InsertedID)
		} else {
			log.Fatal(err)
		}
	} else {
		fmt.Printf("Item with ID: %v contains the following:\n", r.InsertedID)
		fmt.Println("Title:", exampleResult.Title)
		fmt.Println("Author:", exampleResult.Author)
		fmt.Println("ISBN:", exampleResult.ISBN)
		fmt.Println("Publisher:", exampleResult.Publisher)
		fmt.Println("Year:", exampleResult.Year)
		fmt.Println("Copies:", exampleResult.Copies)
	}

	objectId, err := primitive.ObjectIDFromHex("6666e80bf3dd29088b524a4e")
	if err != nil {
		log.Fatal(err)
	}
	new_c := BooksCollection.FindOne(ctx, bson.M{"_id": objectId})
	var exampleRes Book
	err = new_c.Decode(&exampleRes)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			fmt.Println("No document found with ID:", objectId)
		} else {
			log.Fatal(err)
		}
	} else {
		fmt.Printf("Item with ID: %v contains the following:\n", objectId)
		fmt.Println("Title:", exampleRes.Title)
		fmt.Println("Author:", exampleRes.Author)
		fmt.Println("ISBN:", exampleRes.ISBN)
		fmt.Println("Publisher:", exampleRes.Publisher)
		fmt.Println("Year:", exampleRes.Year)
		fmt.Println("Copies:", exampleRes.Copies)
	}
	fmt.Println("END OF FIND ONE")

	fmt.Println()
	fmt.Println("START OF UPDATE ONE")
	rUpd, err := BooksCollection.UpdateOne(ctx,
		bson.M{"_id": r.InsertedID},
		bson.D{
			{Key: "$set", Value: bson.M{"copies": 10}},
		},
	)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Modified documents count:", rUpd.ModifiedCount)
	fmt.Println("END OF UPDATE ONE")
	fmt.Println()

	fmt.Println("START OF FIND ONE")
	c = BooksCollection.FindOne(ctx, bson.M{"_id": r.InsertedID})
	var exampleResult2 Book
	err = c.Decode(&exampleResult2)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			fmt.Println("No document found with ID:", r.InsertedID)
		} else {
			log.Fatal(err)
		}
	} else {
		fmt.Printf("Updated item with ID: %v contains the following:\n", r.InsertedID)
		fmt.Println("Title:", exampleResult2.Title)
		fmt.Println("Author:", exampleResult2.Author)
		fmt.Println("ISBN:", exampleResult2.ISBN)
		fmt.Println("Publisher:", exampleResult2.Publisher)
		fmt.Println("Year:", exampleResult2.Year)
		fmt.Println("Copies:", exampleResult2.Copies)
	}
	fmt.Println("END OF FIND ONE")
	fmt.Println()

	fmt.Println("START OF UPDATE MANY")
	rUpd2, err := BooksCollection.UpdateMany(ctx,
		bson.D{{Key: "Copies", Value: bson.D{{Key: "$lt", Value: 5}}}},
		bson.D{
			{Key: "$set", Value: bson.M{"copies": 5}},
		},
	)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Modified documents count:", rUpd2.ModifiedCount)

	exampleAll, err := BooksCollection.Find(ctx, bson.M{})
	if err != nil {
		log.Fatal(err)
	}
	var examplesRes []Book
	if err = exampleAll.All(ctx, &examplesRes); err != nil {
		log.Fatal(err)
	}
	for _, e := range examplesRes {
		fmt.Println("Title:", e.Title)
		fmt.Println("Author:", e.Author)
		fmt.Println("ISBN:", e.ISBN)
		fmt.Println("Publisher:", e.Publisher)
		fmt.Println("Year:", e.Year)
		fmt.Println("Copies:", e.Copies)
		time.Sleep(time.Millisecond)
	}
	fmt.Println("END OF UPDATE MANY")
	fmt.Println()

	fmt.Println("START OF DELETE ONE")
	rDel, err := BooksCollection.DeleteOne(ctx, bson.M{"_id": r.InsertedID})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Count of deleted documents:", rDel.DeletedCount)
	fmt.Println("END OF DELETE ONE")
	fmt.Println()
}
