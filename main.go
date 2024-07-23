package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

type Product struct {
	Name        string  `json:"name"`
	ProductType string  `json:"product_type"`
	Picture     string  `json:"picture"`
	Price       float64 `json:"price"`
	Description string  `json:"description"`
}

func main() {
	cnStr := "postgres://postgres:secret@localhost:5432/thalesdb?sslmode=disable"

	db, error := sql.Open("postgres", cnStr)

	defer db.Close()

	if error != nil {
		log.Fatal(error)
	}

	if error = db.Ping(); error != nil {
		log.Fatal(error)
	}

	createProductTable(db)

	// product := Product{"Johnsons® Bedtime Baby Lotion", "Personal Hygiene", "https://www.johnsonsbaby.com.sg/sites/jbaby_sg/files/styles/product_image/public/johnsons-baby-bedtime-baby-lotion-front.jpg", 5.35, "Johnsons Bedtime baby lotion with NaturalCalm®  essence clinically proven to help calm and comfort baby before sleep"}
	// pk := insertProduct(db, product)

	// fmt.Printf("ID = %d\n", pk)

	// fmt.Printf(getAllProducts(db))
	fmt.Printf(getSingleProduct(db))

}

func createProductTable(db *sql.DB) {

	createTablequery := `CREATE TABLE IF NOT EXISTS product (
		id SERIAL PRIMARY KEY NOT NULL, 
		name VARCHAR(255), 
		product_type VARCHAR(100), 
		picture VARCHAR(255), 
		price NUMERIC(10,2), 
		description TEXT,
		created timestamp DEFAULT NOW()
	)`

	_, error := db.Exec(createTablequery)

	if error != nil {
		log.Fatal(error)
	}

}

func insertProduct(db *sql.DB, product Product) int {

	query := `INSERT INTO PRODUCT (name, product_type, picture, price, description) 
	VALUES ($1, $2, $3, $4, $5) RETURNING id`

	var pk int
	err := db.QueryRow(query, product.Name, product.ProductType, product.Picture, product.Price, product.Description).Scan(&pk)

	if err != nil {
		log.Fatal(err)
	}

	return pk
}

func getSingleProduct(db *sql.DB) string {
	//SELECT STATEMENT

	var Name, ProductType, Picture, Description string
	var Price float64

	data := Product{}

	query := `SELECT name, product_type, picture, price, description FROM product`

	dbError := db.QueryRow(query).Scan(&Name, &ProductType, &Picture, &Price, &Description)

	if dbError != nil {
		log.Fatal(dbError)
	}

	// fmt.Printf("Name: %s\n", Name)
	// fmt.Printf("ProductType: %s\n", ProductType)
	// fmt.Printf("Picture: %s\n", Picture)
	// fmt.Printf("Price: %f\n", Price)
	// fmt.Printf("Description: %s\n", Description)

	data = Product{Name, ProductType, Picture, Price, Description}

	jsonString, jsonError := json.Marshal(data)

	if jsonError != nil {
		log.Fatal(jsonError)
	}

	return string(jsonString)

}

func getAllProducts(db *sql.DB) string {

	data := []Product{}
	rows, err := db.Query(`SELECT name, product_type, picture, price, description FROM product`)

	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()
	// to scan DB values
	var Name, ProductType, Picture, Description string
	var Price float64

	for rows.Next() {
		rows.Scan(&Name, &ProductType, &Picture, &Price, &Description)
		if err != nil {
			log.Fatal(err)
		}
		data = append(data, Product{Name, ProductType, Picture, Price, Description})
	}

	//convert struct into json string

	jsonString, jsonError := json.Marshal(data)

	if jsonError != nil {
		log.Fatal(jsonError)
	}

	return string(jsonString)

}
