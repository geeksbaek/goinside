package goinside_test

import (
	"net/url"

	"github.com/geeksbaek/goinside"
)

func ExampleSession_ThumbsUp() {
	s.ThumbsUp(article)   // 추천
	s.ThumbsDown(article) // 비추천
}

// 추천이나 비추천은 내부적으로 해당 Article 구조체에 gallInfoDetail 구조체 값이
// 설정되어 있는지 확인한다. 추천과 비추천을 위해선 반드시 필요한 값들이므로
// 해당 값이 설정되어 있지 않다면 스스로 http request 과정을 거쳐 해당 값들을
// fetch하게 된다. 그러나 아래 같이 ThumbsUp 함수가 병렬적으로 실행되는 경우,
// 한 번만 fetch하면 될 gallInfoDetail를 불필요하게 여러번 fetch하는 일이
// 일어날 수 있다. 따라서 PrefetchDetail 라는 메소드를 호출하여 병렬 실행이 
// 일어나기 전에 한 번만 호출하여 낭비를 막을 수 있다.
func ExampleSession_PrefetchDetail() {
	s := goinside.Guest("닉네임", "비밀번호")
	proxys := []*url.URL{} // 프록시의 슬라이스가 있다고 가정

    s.PrefetchDetail(article)
    for _, proxy := range proxys {
        dummy := goinside.Guest("", "")
        dummy.SetTransport(proxy)
        go dummy.ThumbsUp(article)
    }
}
