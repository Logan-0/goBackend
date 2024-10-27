package main

import (
	"fmt"
	"log"

	"github.com/gocql/gocql"
	_ "github.com/lib/pq"
)

var session *gocql.Session

type Storage interface {
	CreateReview(*Review) error
	UpdateReview(*Review) error
	DeleteReview(int) error
	GetReviewById(int) (*Review, error)
}

type Client struct {
	cassandraSession *gocql.Session
}

func InitializeClientAndDB() *Client {
    cluster := gocql.NewCluster("127.0.0.1")
	if cluster == nil {
		log.Fatal("Failed Create Cassandra Cluster")
	}
    cluster.Consistency = gocql.Quorum
    cluster.Keyspace = "reviewKeyspace"
    session,err := cluster.CreateSession()
	if err != nil {
		log.Fatal("Failed Create Cassandra Session")
	}
    return &Client{cassandraSession: session}
}

func (cassandraClient *Client) CreateReviewTable() error {
	session = cassandraClient.cassandraSession
	createReviewDatabaseAndTableQuery := 
	`CREATE KEYSPACE dev WITH replication = {'class':'SimpleStrategy', 'replication_factor' : 1};
	
	use dev;
	
	CREATE TABLE dev.reviews IF NOT EXISTS (
		id serial,
		title text,
		director text,
		rating serial,
		release_date text,
		review text,
		PRIMARY KEY(id));
		
	CREATE INDEX on dev.reviews(id);`

	err := session.Query(createReviewDatabaseAndTableQuery).Exec()
	if err != nil {
		return err
	}
	fmt.Println("********************** Success Create Review Table")
	return nil
}

func (cassandraClient *Client) DropReviewTable() error {
	session = cassandraClient.cassandraSession
	dropTableQuery := `DROP TABLE reviews`
	err := session.Query(dropTableQuery).Exec()
	if err != nil {
		fmt.Println("********************** Failed Drop Review Table")
		return err
	}
	fmt.Println("********************** Success Drop Review Table")
	return nil
}

func (cassandraClient *Client) CreateReview(*Review) error {
	session = cassandraClient.cassandraSession
	return nil
}

func (cassandraClient *Client) UpdateReview(*Review) error {
	session = cassandraClient.cassandraSession
	return nil
}
func (cassandraClient *Client) DeleteReview(id int) error {
	session = cassandraClient.cassandraSession
	return nil
}
func (cassandraClient *Client) GetReviewById(id int) (*Review, error){
	session = cassandraClient.cassandraSession
	return nil, nil
}