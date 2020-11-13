package db

type user interface {
	info() (string, string, string)
	setId(string)
}
