package domain

type OrderBookNotFound struct {
	Message string
}

func (err OrderBookNotFound) Error() string {
	return err.Message
}
