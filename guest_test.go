package goinside

import "log"

func ExampleGuestArticleWrite() {
	s, _ := Guest("ㅇㅇ", "123")

	draft := NewArticleDraft("programming", "test", "test", `C:\Users\geeks\Pictures\1469023529.jpg`)
	err := s.Write(draft)
	if err != nil {
		log.Fatal(err)
	}
}

func ExampleGuestArticleDelete() {
	ExampleGuestArticleWrite()

	s, _ := Guest("ㅇㅇ", "123")

	list, err := FetchList("http://gall.dcinside.com/board/view/?id=programming", 1)
	if err != nil {
		log.Fatal(err)
	}

	for _, v := range list.Articles {
		err = s.Delete(v)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func ExampleGuestCommentWrite() {
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

func ExampleGuestCommentDelete() {
	ExampleGuestCommentWrite()

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

func ExampleGuestAction() {
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

func ExampleGuestReport() {
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
