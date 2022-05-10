package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"

	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var db *gorm.DB
var err error

// Type ini adalah repressentasi dari product
type Products struct {
	ID    int    `json:"id"`
	Code  string `json:"code"`
	Name  string `json:"name"`
	Price string `json:"price" sql:"type:decimal(16,2)"`
}

type Result struct {
	Code    int         `json:"code"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}

func main() {
	db, err = gorm.Open("mysql", "root:@/e_commerce?charset=utf8&parseTime=true")
	if err != nil {
		log.Println("Connection Fail", err)
	} else {
		log.Println("Connection Estabilished")
	}

	db.AutoMigrate(&Products{})
	handleRequests()
}

func handleRequests() {
	log.Println("Start development on 127.0.0.1:9999")
	myRouter := mux.NewRouter().StrictSlash(true)

	myRouter.HandleFunc("/", Homepage)
	myRouter.HandleFunc("/api/products", createProduct).Methods("POST")
	myRouter.HandleFunc("/api/products", getProduct).Methods("GET")
	myRouter.HandleFunc("/api/products/", getProduct).Methods("GET")
	myRouter.HandleFunc("/api/products/{id}", findProduct).Methods("GET")
	myRouter.HandleFunc("/api/products/{id}", deleteProduct).Methods("DELETE")

	log.Fatal(http.ListenAndServe("127.0.0.1:9999", myRouter))

}

func Homepage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Welcome")
}

func createProduct(w http.ResponseWriter, r *http.Request) {

	payloads, _ := ioutil.ReadAll(r.Body)
	var product Products
	json.Unmarshal(payloads, &product)

	w.Header().Set("Content-Type", "application/json")

	res := Result{Code: 200, Data: product, Message: "Success Create Product"}

	if err := db.Create(&product).Error; err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		res = Result{Code: 500, Data: nil, Message: err.Error()}
	} else {
		w.WriteHeader(http.StatusOK)
	}

	result, err := json.Marshal(res)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Write(result)
}

func getProduct(w http.ResponseWriter, r *http.Request) {
	products := []Products{}
	db.Find(&products)

	res := Result{Code: 200, Data: products, Message: "Success get products"}
	result, err := json.Marshal(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(result)
}

func findProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	productsID := vars["id"]
	// db.Find(&products)
	var product Products
	db.First(&product, productsID)

	res := Result{Code: 200, Data: product, Message: "Success get products"}
	result, err := json.Marshal(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(result)
}

func updateProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	productsID := vars["id"]

	payloads, _ := ioutil.ReadAll(r.Body)
	var productUpdate Products
	json.Unmarshal(payloads, &productUpdate)

	var product Products
	db.First(&product, productsID)
	db.Model(&product).Update(productUpdate)

	res := Result{Code: 200, Data: product, Message: "Success get products"}
	result, err := json.Marshal(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(result)
}

func deleteProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	productsID := vars["id"]

	var product Products
	db.First(&product, productsID)
	db.Delete(&product)

	res := Result{Code: 200, Data: product, Message: "Success get products"}
	result, err := json.Marshal(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(result)
}
