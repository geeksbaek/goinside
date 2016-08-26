package gallog_test

import (
	"log"
	"time"

	"github.com/geeksbaek/goinside/gallog"
)

func Example() {
	s, err := gallog.Login("ID", "PASSWORD")
	if err != nil {
		log.Fatal(err)
	}
	log.Println(s.Name, "님. 로그인에 성공하였습니다.")

	log.Println("모든 글과 댓글을 불러오는 중입니다. 잠시만 기다려주세요.")
	start := time.Now()
	data := s.FetchAll()

	log.Printf("글 %v개, 댓글 %v개 ", len(data.As), len(data.Cs))
	log.Println("불러오기를 완료하였습니다.")
	log.Println("불러오는 데 걸린 시간 :", time.Since(start))

	log.Println("삭제를 시작합니다. 잠시만 기다려주세요.")
	middle := time.Now()
	s.DeleteAll(data, func(i, n int) {
		log.Printf("%v/%v 완료", i, n)
	})

	log.Println("삭제가 끝났습니다.")
	log.Println("삭제하는 데 걸린 시간 :", time.Since(middle))
}
