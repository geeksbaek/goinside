// Package goinside 는 Go로 구현한 비공식 디시인사이드 API입니다.
package goinside

/*
goinside는 다음 기능을 지원합니다.

1. 유동닉 또는 고정닉으로 글과 댓글의 작성 및 삭제
2. 추천과 비추천, 신고
3. 특정 갤러리의 특정 페이지에 있는 게시물 및 댓글 가져오기
4. 모든 일반 갤러리, 마이너 갤러리 정보 가져오기
5. 프록시 모드

글이나 댓글을 작성하거나 삭제하려면 우선 세션을 생성해야 합니다.

유동닉 세션은 Guest 함수로 생성하며, 닉네임과 비밀번호를 인자로 받습니다.
빈 문자열을 인자로 넘길 경우 에러를 반환합니다.

        s, err := goinside.Guest("ㅇㅇ", "123")
        err != nil {
                log.Fatal(err)
        }

고정닉 세션은 Login 함수로 생성합니다. 디시인사이드 ID와 비밀번호를 인자로 받습니다.
로그인에 실패할 경우 에러를 반환합니다.

        s, err := goinside.Login("ID", "PASSWORD")
        err != nil {
                log.Fatal(err)
        }

글이나 댓글을 작성하기 위해서는 Draft를 먼저 생성해야 합니다.
Draft를 생성하기 위해 NewArticleDraft, NewCommentDraft 함수가 있습니다.
해당 함수로 생성된 Draft 객체를 Wrtie 메소드로 전달하여 글을 작성합니다.

        draft := NewArticleDraft("programming", "제목", "내용", "이미지 경로")
        if err := s.Write(draft); err != nil {
                log.Fatal(err)
        }

갤러리의 글을 가져오는데는 세션이 필요하지 않습니다.
다음 코드는 programming 갤러리의 개념글 목록 1페이지에 있는 글들을 가져옵니다.

        URL := "http://gall.dcinside.com/board/lists/?id=programming"
        list, err := goinside.FetchBestList(URL, 1)
        if err != nil {
                log.Fatal(err)
        }

가져온 글을 Write 메소드에 넘겨서 바로 재작성 할 수도 있습니다.
그러나 FetchList, FetchBestList 함수로 가져온 Item들은
아직 글의 내용을 알 수 없는 상태입니다.
이 Item이 Write 함수의 인자로 전달될 때는 글의 제목을 그대로 내용으로 쓰도록 되어있습니다.

        if err := s.Write(list.Items[0]); err != nil {
                log.Fatal(err)
        }

해당 세션을 프록시로 전환할 수도 있습니다. 아래 코드의 proxy 변수는 *url.URL
타입이라고 가정합니다.

        s.Connection().SetTransport(proxy)

http 요청에 타임아웃을 설정할 수도 있습니다.

        s.Connection().SetTimeout(time.Second * 3)

*/
