package main

import (
	"fmt"
	"time"
)

type CreateReviewRequest struct  {
	Title string `json:"title"`
	Director string `json:"director"`
	ReleaseDate string `json:"realeaseDate"`
	Rating string `json:"rating"`
	ReviewNotes string `json:"reviewNotes"`
}

type Review struct{
	ID int `json:"id"`
	Title string `json:"title"`
	Director string `json:"director"`
	ReleaseDate string `json:"realeaseDate"`
	Rating string `json:"rating"`
	ReviewNotes string `json:"reviewNotes"`
	CreatedAt string `json:"createdAt"`
}

func NewReview(title string, director string, releaseDate string, rating string, reviewNotes string) *Review {
	dateTime, err := time.Parse(time.RFC822, releaseDate)
	createdAtTime := time.Now().Format(time.RFC822)

	dateString := dateTime.String()[0:19]
	createdAtString := createdAtTime[0:19]

	if err != nil { 
		fmt.Println("********************** Failed: Ping DB: %w",err)
	}
	return &Review {
		Title: title,
		Director: director,
		ReleaseDate: dateString,
		Rating: rating,
		ReviewNotes: reviewNotes,
		CreatedAt: createdAtString,
	}
}