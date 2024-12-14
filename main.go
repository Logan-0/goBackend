package main

import (
	"fmt"
	"log"
)

func main() {
	// Initialize Database
	client, err := InitializeClientAndDB()
	if err != nil {
		log.Fatal("********************** Failed: Connection to Database ", err.Error())
	}
	fmt.Println("********************** Success: Database Port 5321")

	// Initialize Keyspace and Database Table
	// To be Used Later after testing for devops debugging.
	// err = client.CreateReviewTable()
	// if err != nil {
	// 	log.Fatal("********************** Failed: Create Review Table" + err.Error())
	// }
	// fmt.Println("********************** Success: Create Review Table")

	// Initialize Server
	fmt.Println("********************** Success: Server Running 8080")
	RunNewServer("0.0.0.0:8080", client)
}
