package goinside

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

func TestLogin(t *testing.T) {
	auth := readAuth("auth.json")
	s, err := Login(auth.ID, auth.PW)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(s.memberSessionDetail)
}

func TestMemberArticleWrite(t *testing.T) {
	auth := readAuth("auth.json")
	s, err := Login(auth.ID, auth.PW)
	if err != nil {
		log.Fatal(err)
	}

	draft := NewArticleDraft("programming", "test", "test", `C:\Users\geeks\Pictures\1469023529.jpg`)

	err = s.Write(draft)
	if err != nil {
		log.Fatal(err)
	}
}

func TestMemberArticleDelete(t *testing.T) {
	TestMemberArticleWrite(t)

	auth := readAuth("auth.json")
	s, err := Login(auth.ID, auth.PW)
	if err != nil {
		log.Fatal(err)
	}

	list, err := FetchList("http://gall.dcinside.com/board/view/?id=programming", 1)
	if err != nil {
		log.Fatal(err)
	}

	err = s.Delete(list.Articles[0])
	if err != nil {
		log.Fatal(err)
	}
}

func TestMemberCommentWrite(t *testing.T) {
	auth := readAuth("auth.json")
	s, err := Login(auth.ID, auth.PW)
	if err != nil {
		log.Fatal(err)
	}

	list, err := FetchList("http://gall.dcinside.com/board/view/?id=programming", 1)
	if err != nil {
		log.Fatal(err)
	}

	draft := NewCommentDraft(list.Articles[0], "123")
	err = s.Write(draft)
	if err != nil {
		log.Fatal(err)
	}
}

func TestMemberCommentDelete(t *testing.T) {
	TestMemberCommentWrite(t)

	auth := readAuth("auth.json")
	s, err := Login(auth.ID, auth.PW)
	if err != nil {
		log.Fatal(err)
	}

	list, err := FetchList("http://gall.dcinside.com/board/view/?id=programming", 1)
	if err != nil {
		log.Fatal(err)
	}

	article, err := FetchArticle(list.Articles[0].URL)
	if err != nil {
		log.Fatal(err)
	}

	// 아직 지원하지 않음
	err = s.Delete(article.Detail.Comments[0])
	if err != nil {
		log.Fatal(err)
	}
}

func TestMemberAction(t *testing.T) {
	auth := readAuth("auth.json")
	s, err := Login(auth.ID, auth.PW)
	if err != nil {
		log.Fatal(err)
	}

	list, err := FetchList("http://gall.dcinside.com/board/view/?id=programming", 1)
	if err != nil {
		log.Fatal(err)
	}

	err = s.ThumbsUp(list.Articles[0])
	if err != nil {
		log.Fatal(err)
	}

	err = s.ThumbsDown(list.Articles[0])
	if err != nil {
		log.Fatal(err)
	}
}

func TestMemberReport(t *testing.T) {
	auth := readAuth("auth.json")
	s, err := Login(auth.ID, auth.PW)
	if err != nil {
		log.Fatal(err)
	}

	list, err := FetchList("http://gall.dcinside.com/board/view/?id=programming", 1)
	if err != nil {
		log.Fatal(err)
	}

	err = s.Report(list.Articles[0], "test")
	if err != nil {
		log.Fatal(err)
	}
}
