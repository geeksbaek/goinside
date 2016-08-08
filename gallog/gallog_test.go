package gallog

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
)

type tempAuth struct {
	ID, PW string
}

func readAuth(path string) (auth *tempAuth) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}

	auth = &tempAuth{}
	err = json.Unmarshal(data, auth)
	if err != nil {
		log.Fatal(err)
	}
	return
}

// func TestGallogLogin(t *testing.T) {
// 	auth := readAuth("../auth.json")
// 	s, err := Login(auth.ID, auth.PW)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	fmt.Println(s.cookies)
// }

func ExampleFetch() {
	auth := readAuth("../auth.json")
	s, err := Login(auth.ID, auth.PW)
	if err != nil {
		log.Fatal(err)
	}

	data := s.FetchAll()
	fmt.Println("article num :", len(data.as))
	fmt.Println("comment num :", len(data.cs))
}
