package main

type Review struct{
	ID int `json:"id"`
	Title string `json:"title"`
	Director string `json:"director"`
	ReleaseDate string `json:"realeaseDate"`
	Rating string `json:"rating"`
	Review string `json:"review"`
}

func NewReview(title string, director string, releaseDate string, rating string, reviewNotes string) *Review {
	return &Review {
		Director: director,
		ReleaseDate: releaseDate,
		Rating: rating,
		Review: reviewNotes,
	}
}