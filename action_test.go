package goinside_test

import (
	"net/url"

	"github.com/geeksbaek/goinside"
)

func ExampleSession_ThumbsUp() {
	s.ThumbsUp(article)   // 추천
	s.ThumbsDown(article) // 비추천
}

// 추천, 비추천 요청에는 특별한 값이 필요하다. 이 값은 gallInfoDetail 구조체로
// 정의되어 있는데, 이 값이 없으면 추천, 비추천 요청을 하기 전에 이것을
// 가져오는 작업을 거친다. 그러데 만약 아래와 같이 추천, 비추천 요청을
// 동시적으로 보내는 경우에는 이 값을 여러 번 중복해서 가져오는 낭비가 발생한다.
// 이러한 경우에는 PrefetchDetail 함수를 미리 호출하여 낭비를 피할 수 있다.
func ExampleSession_PrefetchDetail() {
	s, _ := goinside.Guest("닉네임", "비밀번호")
	proxys := []*url.URL{} // 프록시의 슬라이스가 있다고 가정

	s.PrefetchDetail(article)
	for _, proxy := range proxys {
		proxy := proxy
		go func() {
			dummy, _ := goinside.Guest(".", ".")
			dummy.SetTransport(proxy)
			dummy.ThumbsUp(article)
		}()
	}
}
