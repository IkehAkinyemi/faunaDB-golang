package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	f "github.com/fauna/faunadb-go/faunadb"
	"github.com/gorilla/mux"
)

type Post struct {
	Id      string `fauna:"Id"`
	Title   string `fauna:"Title"`
	Content string `fauna:"Content"`
	Summary string `fauna:"Summary"`
}

var (
	data = f.ObjKey("data")
	ref  = f.ObjKey("ref")
)
var postId f.RefV

var Blog []Post

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to REST API with FaunaDB")
	fmt.Println("We've reached the Home page endpoint!")
}

func getSinglePost(w http.ResponseWriter, r *http.Request) {
	client := f.NewFaunaClient("fnAEYRolCWACTEED7DnXPbR-4dLKvMSH2KwYPywG")
	// Create a class to store profiles
	_, _ = client.Query(f.CreateClass(f.Obj{"name": "Blog"}))

	req, _ := ioutil.ReadAll(r.Body)
	var post Post
	json.Unmarshal(req, &post)

	// Retrieve profile by its ID
	value, _ := client.Query(f.Get(postId))
	_ = value.At(data).Get(&post)

	json.NewEncoder(w).Encode(value)
}

func createNewPost(w http.ResponseWriter, r *http.Request) {
	client := f.NewFaunaClient("fnAEYRolCWACTEED7DnXPbR-4dLKvMSH2KwYPywG")
	// Create a class to store profiles
	_, _ = client.Query(f.CreateClass(f.Obj{"name": "Blog"}))

	req, _ := ioutil.ReadAll(r.Body)
	var post Post
	json.Unmarshal(req, &post)

	// Save profile at FaunaDB
	newProfile, _ := client.Query(
		f.Create(
			f.Class("Blog"),
			f.Obj{"data": post},
		),
	)

	// Get generated profile ID
	_ = newProfile.At(ref).Get(&postId)

	Blog = append(Blog, post)
	json.NewEncoder(w).Encode(post)
}

func deletePost(w http.ResponseWriter, r *http.Request) {
	client := f.NewFaunaClient("fnAEYRolCWACTEED7DnXPbR-4dLKvMSH2KwYPywG")
	// Create a class to store profiles
	_, _ = client.Query(f.CreateClass(f.Obj{"name": "Blog"}))

	// Delete profile using its ID
	_, _ = client.Query(f.Delete(postId))
}

func updatePost(w http.ResponseWriter, r *http.Request) {
	client := f.NewFaunaClient("fnAEYRolCWACTEED7DnXPbR-4dLKvMSH2KwYPywG")
	// Create a class to store profiles
	_, _ = client.Query(f.CreateClass(f.Obj{"name": "Blog"}))

	// Update existing profile entry
	_, _ = client.Query(
		f.Update(
			postId,
			f.Obj{"data": f.Obj{
				"Title": "Adieu to Fauna blog posts", "Content": "In the next article, we'll build a simple micro-service", "Summary": "This article featured Golang architectures",
			}},
		),
	)

}

func handleRequests() {
	r := mux.NewRouter().StrictSlash(true)
	r.HandleFunc("/", homePage)
	r.HandleFunc("/blog", getSinglePost).Methods("GET")
	r.HandleFunc("/blog/post", createNewPost)
	r.HandleFunc("/blog", deletePost).Methods("DELETE")
	r.HandleFunc("/blog", updatePost).Methods("PUT")
	log.Fatal(http.ListenAndServe(":10000", r))
}

func main() {
	handleRequests()
}
