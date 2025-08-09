package main

import (
    "fmt"
    "time"
)

type CreateReviewRequest struct {
    Title       string `json:"title"`
    Director    string `json:"director"`
    ReleaseDate string `json:"releaseDate"`
    Rating      string `json:"rating"`
    ReviewNotes string `json:"reviewNotes"`
}

type Review struct {
    ID          int    `json:"id"`
    Title       string `json:"title"`
    Director    string `json:"director"`
    ReleaseDate string `json:"releaseDate"`
    Rating      string `json:"rating"`
    ReviewNotes string `json:"reviewNotes"`
    DateCreated string `json:"dateCreated"`
}

func NewReview(title string, director string, releaseDate string, rating string, reviewNotes string) *Review {
    dateTime, err := time.Parse(time.RFC822, releaseDate)
    if err != nil {
        fmt.Printf("Failed to Parse Date: Format 01 Jan 22 00:00 UTC You put: %s\n", releaseDate)
        // Use current time as fallback if parsing fails
        dateTime = time.Now()
    }
    
    dateString := dateTime.Format(time.RFC822)[0:19]
    dateCreated := time.Now().Format(time.RFC822)[0:19]

    return &Review{
        Title:       title,
        Director:    director,
        ReleaseDate: dateString,
        Rating:      rating,
        ReviewNotes: reviewNotes,
        DateCreated: dateCreated,
    }
}