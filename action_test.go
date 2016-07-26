package goinside_test

func ExampleSession_ThumbsUp() {
	s.ThumbsUp(article)   // 추천
	s.ThumbsDown(article) // 비추천
}
