package manage

import (
	"log"
	"testing"
)

func TestPortFormat(t *testing.T) {
	port := "8080:8080,9090:9090"
	err := checkPortFormat(port)
	if err != nil {
		log.Println(err.Error())
	}
}
