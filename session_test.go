package goinside

// import (
// 	"encoding/json"
// 	"fmt"
// 	"io/ioutil"
// 	"log"
// 	"os"
// 	"testing"
// )

// func TestLogin(t *testing.T) {
// 	file, e := ioutil.ReadFile("./auth.json")
// 	if e != nil {
// 		fmt.Printf("File error: %v\n", e)
// 		os.Exit(1)
// 	}
// 	var auth struct {
// 		ID string
// 		PW string
// 	}
// 	json.Unmarshal(file, &auth)

// 	s, err := Login(auth.ID, auth.PW)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	articles, comments, err := s.GetGallogData()
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	fmt.Println(articles)
// 	fmt.Println(comments)

// 	s.Logout()
// }
