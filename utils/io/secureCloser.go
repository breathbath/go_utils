package io

import (
	"io"
)

//CloseResourceSecure allows to call closer and handle errors as warnings in logs
func CloseResourceSecure(name string, c io.Closer) {
	if c == nil {
		return
	}

	err := c.Close()
	if err != nil {
		OutputError(err, "", "Failed to close resource '%s': %v", name, c)
	}
}
