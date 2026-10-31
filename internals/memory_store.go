package booking

type memorystore struct{
	bookings map[string]bookings
}

func newmemorystore() *memorystore{
	return &memorystore{
		bookings map[string]bookings{}
	}
}

func() book(b booking)error
func() list_bookings(movie_id string)[]booking
func() list_all_detail(id string)[]booking







  