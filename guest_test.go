package goinside

import (
	"log"
	"testing"
)

func TestGuestArticleWrite(t *testing.T) {
	s, _ := Guest("ㅇㅇ", "123")

	draft := NewArticleDraft("programming", "test", "test", `C:\Users\geeks\Pictures\1469023529.jpg`)
	err := s.Write(draft)
	if err != nil {
		log.Fatal(err)
	}
}

func TestGuestArticleDelete(t *testing.T) {
	TestGuestArticleWrite(t)

	s, _ := Guest("ㅇㅇ", "123")

	list, err := FetchList("http://gall.dcinside.com/board/view/?id=programming", 1)
	if err != nil {
		log.Fatal(err)
	}

	err = s.Delete(list.Articles[0])
	if err != nil {
		log.Fatal(err)
	}
}

func TestGuestCommentWrite(t *testing.T) {
	s, _ := Guest("ㅇㅇ", "123")

	list, err := FetchList("http://gall.dcinside.com/board/view/?id=programming", 1)
	if err != nil {
		log.Fatal(err)
	}

	draft := NewCommentDraft(list.Articles[0], "test")
	err = s.Write(draft)
	if err != nil {
		log.Fatal(err)
	}
}

func TestGuestCommentDelete(t *testing.T) {
	TestGuestCommentWrite(t)

	s, _ := Guest("ㅇㅇ", "123")

	list, err := FetchList("http://gall.dcinside.com/board/view/?id=programming", 1)
	if err != nil {
		log.Fatal(err)
	}

	article, err := FetchArticle(list.Articles[0].URL)
	if err != nil {
		log.Fatal(err)
	}

	err = s.Delete(article.Detail.Comments[0])
	if err != nil {
		log.Fatal(err)
	}
}

func TestGuestAction(t *testing.T) {
	s, _ := Guest("ㅇㅇ", "123")

	list, err := FetchList("http://gall.dcinside.com/board/view/?id=programming&no=618139&page=1", 1)
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

func TestGuestReport(t *testing.T) {
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
