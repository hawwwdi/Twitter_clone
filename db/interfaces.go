package db

type user interface {
	Info() (string, string, string)
	SetId(string)
}
