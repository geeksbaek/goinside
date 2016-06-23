package goinside

import (
	"fmt"
	"log"
	"testing"
)

func TestGetAllGall(t *testing.T) {
	galls, err := GetAllGall()
	if err != nil {
		log.Fatal(err)
	}
	for _, v := range galls {
		fmt.Println(v)
	}
}
