package models

type Mail struct {
	Dest string
	Body string
}

type Subscriber struct {
	ID                int
	Address           string
	Name              string
	Surname           string
	FavouriteCategory string
}
