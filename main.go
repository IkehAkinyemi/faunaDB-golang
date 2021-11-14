package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type Post struct {
	Id string `json:"Id"`
	Title string `json:"Title"`
	Content string `json:"Content"`
	Summary string `json:"Summary"`
}

var Blog []Post

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to REST API with FaunaDB")
	fmt.Println("We've reached the Home page endpoint!")
}

func getAllPosts(w http.ResponseWriter, r *http.Request) {
	fmt.Println("You just encountered the get-all-post endpoint!")
	json.NewEncoder(w).Encode(Blog)
}

func getSinglePost(w http.ResponseWriter, r *http.Request) {
	paras := mux.Vars(r)
	id := paras["id"]

	for _, post := range Blog {
		if post.Id == id {
			json.NewEncoder(w).Encode(post)
		}
	}
}

func createNewPost(w http.ResponseWriter, r *http.Request) {
	req, _ := ioutil.ReadAll(r.Body)
	var post Post
	json.Unmarshal(req, &post)
	Blog = append(Blog, post)
	json.NewEncoder(w).Encode(post)
}

func deletePost(w http.ResponseWriter, r *http.Request) {
	paras := mux.Vars(r)
	id := paras["id"]

	for i, post := range Blog {
		if post.Id == id {
			Blog = append(Blog[:i], Blog[i + 1:]...)
		}
	}
}

func updatePost(w http.ResponseWriter, r *http.Request) {
	req, _ := ioutil.ReadAll(r.Body)
	var rPost Post
	json.Unmarshal(req, &rPost)

	paras := mux.Vars(r)
	id := paras["id"]

	for i, post := range Blog {
		if post.Id == id {
			Blog[i] = rPost
		}
	}
}

func handleRequests() {
	r := mux.NewRouter().StrictSlash(true)
	r.HandleFunc("/", homePage)
	r.HandleFunc("/blog", getAllPosts)
	r.HandleFunc("/blog/{id}", getSinglePost).Methods("GET")
	r.HandleFunc("/blog/post", createNewPost)
	r.HandleFunc("/blog/{id}", deletePost).Methods("DELETE")
	r.HandleFunc("/blog/{id}", updatePost).Methods("PUT")
	log.Fatal(http.ListenAndServe(":10000", r))
}

func main() {
	Blog = []Post{
		{Id:"1", Title: "Welcome to Fauna blog posts", Content: "In this article, we'll build a simple REST", Summary: "This article featured Golang"},
		{Id: "2", Title: "Adieu to Fauna blog posts", Content: "In the next article, we'll build a simple micro-service", Summary: "This article featured Golang architectures"},
	}
	handleRequests()
}