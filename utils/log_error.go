package utils

// Log error.
// Example func foo() (err, val)
// err, val := foo()
// logError(err)
func LogError(err error) {
	if err != nil {
		panic(err)
	}
}
