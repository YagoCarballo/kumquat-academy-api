package endpoints

import (
	"net/http"
)

func SayHello(name string) (int, map[string]interface{}) {
	return http.StatusOK, map[string]interface{}{
		"message": "Hello, " + name + "!!",
	}
}
