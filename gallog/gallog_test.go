package gallog

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"testing"
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

func TestGallogLogin(t *testing.T) {
	auth := readAuth("../auth.json")
	s, err := Login(auth.ID, auth.PW)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(s.cookies)
}

func TestFetchGallogArticle(t *testing.T) {
	auth := readAuth("../auth.json")
	s, err := Login(auth.ID, auth.PW)
	if err != nil {
		log.Fatal(err)
	}

	as := s.FetchAllArticle()
	s.DeleteArticle(as)
	fmt.Println("done.")
}

func TestFetchGallogComment(t *testing.T) {
	auth := readAuth("../auth.json")
	s, err := Login(auth.ID, auth.PW)
	if err != nil {
		log.Fatal(err)
	}

	cs := s.FetchAllComment()
	fmt.Println(len(cs))
	s.DeleteComment(cs)
	fmt.Println("done.")
}
