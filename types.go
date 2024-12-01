package main

import (
	"fmt"
	"time"
)

type CreateReviewRequest struct  {
	Title string `json:"title"`
	Director string `json:"director"`
	ReleaseDate string `json:"releaseDate"`
	Rating string `json:"rating"`
	ReviewNotes string `json:"reviewNotes"`
}

type Review struct{
	ID int `json:"id"`
	Title string `json:"title"`
	Director string `json:"director"`
	ReleaseDate string `json:"releaseDate"`
	Rating string `json:"rating"`
	ReviewNotes string `json:"reviewNotes"`
	CreatedAt string `json:"createdAt"`
}

func NewReview(title string, director string, releaseDate string, rating string, reviewNotes string) *Review {
	dateTime, err := time.Parse(time.RFC822, releaseDate)
	dateString := dateTime.Format(time.RFC822)[0:19]
	createdAtTime := time.Now().Format(time.RFC822)[0:19]
	if err != nil {
		fmt.Println("Failed to Parse Date: Format 01 Jan 22 00:00 UTC You put: %w", err)
	}

	return &Review {
		Title: title,
		Director: director,
		ReleaseDate: dateString,
		Rating: rating,
		ReviewNotes: reviewNotes,
		CreatedAt: createdAtTime,
	}
}