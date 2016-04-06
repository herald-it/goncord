package models

type Setting struct {
	Database struct {
		Host       string
		DbName     string
		TokenTable string
		UserTable  string
	}
}
