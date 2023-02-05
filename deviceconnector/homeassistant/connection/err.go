package connection

import "fmt"

type remoteError struct {
	Code    string
	Message string
}

func (e remoteError) Error() string {
	return fmt.Sprintf("remote error: %v - %v", e.Code, e.Message)
}
