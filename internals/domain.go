package booking

type booking struct{
	id string,
	movie_id string,
	seat_id string,
	user_id string,
	status string
}

type booking_store interface{
	book(b booking)error
	list_bookings(movie_id string)[]booking
	list_all_detail(id string)[]booking
}