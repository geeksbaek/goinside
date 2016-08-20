# goinside [![Build Status](https://travis-ci.org/geeksbaek/goinside.svg?branch=master)](https://travis-ci.org/geeksbaek/goinside) [![GoDoc](https://godoc.org/github.com/geeksbaek/goinside?status.svg)](https://godoc.org/github.com/geeksbaek/goinside)

이 라이브러리는 디시인사이드 비공식 API 입니다.
API에 대한 설명은 [godoc](https://godoc.org/github.com/geeksbaek/goinside)에서 보실 수 있습니다. 

## Install
```
go get -u github.com/geeksbaek/goinside/...
```

## See also

- [goinside-image-crawler](https://github.com/geeksbaek/goinside-image-crawler)
- [goinside-gallog-cleaner](https://github.com/geeksbaek/goinside-gallog-cleaner)

## Update

### 2016-08-20

갤러리 정보를 가져오는 부분을 디시인사이드 API를 통해 가져오도록 수정하였습니다. 이에 따라 List, Article, Comment 구조체가 모두 변경되었습니다. 이제 디시인사이드 API가 제공되지 않는 gallog 패키지를 제외하면 goinside의 모든 기능은 디시인사이드 API를 통해 구현됩니다.

여기서 말하는 디시인사이드 API는, 디시인사이드 공식 App 및 서드파티 App에서 사용되는 HTTP 기반의 API를 말합니다. 기존의 goinside가 디시인사이드 Web을 파싱하여 비공식 API를 제공하는 형태였다면, 이제는 디시인사이드 HTTP API의 Go 버전 Wrapper를 구현하는 형태가 되었습니다.

이제 FetchList 함수는 ListItem 구조체의 슬라이스을 반환합니다. ListItem은 IP 멤버 변수를 가집니다. 디시인사이드 API를 사용하게 되면서 개별 Article을 각각 Fetch하지 않아도 Article Author의 IP를 알 수 있게 되었습니다.

이제 댓글에서 DCcon과 보이스 리플을 구분합니다. Comment 구조체는 Type이라는 이름의 CommentType 타입이 댓글의 타입을 표현합니다. 또한 Comment의 HTMLContent 메소드는 해당 타입에 대응되는 완성된 HTML 코드를 반환합니다. (DCcon일 경우 img, 보이스 리플일 경우 audio element)

### 2016-08-09

하위 패키지 `github.com/geeksbaek/goinside/gallog ` 를 추가하고 갤로그 관련 API를 추가하였습니다.

### 2016-08-07

API의 구조가 크게 변경되었습니다. godoc을 확인해주세요.

## 주의

현재 개발 중이며 언제든지 API 구조가 변경될 수 있습니다.

Jongyeol Baek <geeksbaek@gmail.com>
