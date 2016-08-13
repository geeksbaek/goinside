package goinside

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

func ExampleLogin() {
	auth := readAuth("auth.json")
	s, err := Login(auth.ID, auth.PW)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(s.MemberSessionDetail)
}

func ExampleMemberArticleWrite() {
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

func ExampleMemberArticleDelete() {
	ExampleMemberArticleWrite()

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

func ExampleMemberCommentWrite() {
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

func ExampleMemberCommentDelete() {
	ExampleMemberCommentWrite()

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

func ExampleMemberAction() {
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

func ExampleMemberReport() {
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

func ExampleLoginFailed() {
	_, err := Login("", "")
	if err != errLoginFailed {
		log.Fatal("로그인에 실패하였습니다.")
	}
}
