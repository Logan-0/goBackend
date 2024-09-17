package main

type Review struct{
	ID int
	Title string
	Director string
	ReleaseDate int
	Rating float64
	Review string
}

func NewReview(title string, director string, releaseDate int, rating float64, review string) *Review {
	return &Review {
		Director: director,
		ReleaseDate: releaseDate,
		Rating: rating,
		Review: review,
	}
}