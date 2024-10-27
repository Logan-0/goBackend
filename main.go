package main

import (
	"fmt"
)

func main() {

	// Initialize Database
	client := InitializeClientAndDB()
	fmt.Println("********************** Success Create Cassandra Session")
	
	// Initialize Keyspace and Database Table
	err := client.CreateReviewTable();
	if err != nil {
		fmt.Println("********************** Success Create Review Table")
	}

	// Initialize Server
	RunNewServer(":8080")
	fmt.Println("********************** Success Create Server")
}