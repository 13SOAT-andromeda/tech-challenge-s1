package ports

type Email interface {
	Send(name string, email string, subject string, html string) error
}
