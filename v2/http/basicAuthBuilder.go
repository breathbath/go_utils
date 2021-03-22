package http

import "encoding/base64"

func BuildBasicAuthString(login, pass string) string {
	auth := login + ":" + pass
	return base64.StdEncoding.EncodeToString([]byte(auth))
}
