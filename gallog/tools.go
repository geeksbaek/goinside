package gallog

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type _JSONResponse struct {
	Result bool
	Cause  string
}

// regex
var (
	gallogRe              = regexp.MustCompile(`http://gallog.dcinside.com`)
	gallogArticleURLRe    = regexp.MustCompile(`gid=([^&]+)&cid=([^&]+)&page=.*&pno=([^&]+)&logNo=([^&]+)&mode=([^&']+)`)
	gallogCommentURLRe    = regexp.MustCompile(`gid=([^&]+)&cid=.*&id=&no=([^&]+)&c_no=([^&]+)&logNo=([^&]+)&rpage=.*`)
	gallogGallIDRe        = regexp.MustCompile(`<INPUT TYPE="hidden" NAME="id" value=(?:"|')(.+)(?:"|')>`)
	gallogSecretKeyPairRe = regexp.MustCompile(`<INPUT TYPE="hidden" NAME=".*" id=(?:"|')([^'"]+)(?:"|') value=(?:"|')([^'"]{10,})(?:"|')>`)
	gallogCIDRe           = regexp.MustCompile(`<INPUT TYPE="hidden" NAME="cid" value="([^"]+)">`)
)

// errors
var (
	errParseGallogArticleURL    = errors.New("cannot parse gallog article url")
	errParseGallogCommentURL    = errors.New("cannot parse gallog comment url")
	errUnknownCause             = errors.New("result false with empty cause")
	errParseGallogGallID        = errors.New("cannot find gall id")
	errParseGallogSecreyKeyPair = errors.New("cannot find secret key pair")
	errParseGallogCID           = errors.New("cannot find cid")
)

// formatting
const (
	gallogArticlePageURLFormat   = "http://gallog.dcinside.com/inc/_mainGallog.php?page=%v&gid=%v"
	gallogCommentPageURLFormat   = "http://gallog.dcinside.com/inc/_mainGallog.php?rpage=%v&gid=%v"
	gallogArticleDetailURLFormat = "http://gallog.dcinside.com/inc/_deleteLog.php?gid=%v&cid=%v&page=&pno=%v&logNo=%v&mode=%v"
	gallogCommentDetailURLFormat = "http://gallog.dcinside.com/inc/_deleteLogRep.php?gid=%v&cid=&id=&no=%v&c_no=%v&logNo=%v&rpage="
)

func _Form(m map[string]string) io.Reader {
	data := url.Values{}
	for k, v := range m {
		data.Set(k, v)
	}
	return strings.NewReader(data.Encode())
}

func _ParseGallogArticleURL(URL string) (a *ArticleMicroInfo, err error) {
	matched := gallogArticleURLRe.FindStringSubmatch(URL)
	if len(matched) != 6 {
		err = errParseGallogArticleURL
		return
	}
	a = &ArticleMicroInfo{matched[1], matched[2], matched[3], matched[4], matched[5]}
	return
}

func _ParseGallogCommentURL(URL string) (c *CommentMicroInfo, err error) {
	matched := gallogCommentURLRe.FindStringSubmatch(URL)
	if len(matched) != 5 {
		err = errParseGallogCommentURL
		return
	}
	c = &CommentMicroInfo{matched[1], matched[2], matched[3], matched[4]}
	return
}

func _ParseGallogGallID(body string) (id string, err error) {
	matched := gallogGallIDRe.FindStringSubmatch(body)
	if len(matched) != 2 {
		err = errParseGallogGallID
		return
	}
	id = matched[1]
	return
}

func _ParseGallogSecretKeyPair(body string) (key, val string, err error) {
	matched := gallogSecretKeyPairRe.FindStringSubmatch(body)
	if len(matched) != 3 {
		err = errParseGallogSecreyKeyPair
		return
	}
	key, val = matched[1], matched[2]
	return
}

func _ParseGallogCID(body string) (cid string, err error) {
	matched := gallogCIDRe.FindStringSubmatch(body)
	if len(matched) != 2 {
		err = errParseGallogCID
		return
	}
	cid = matched[1]
	return
}

func _GallogArticlePageURL(gid string, page int) string {
	return fmt.Sprintf(gallogArticlePageURLFormat, page, gid)
}

func _GallogCommentPageURL(gid string, page int) string {
	return fmt.Sprintf(gallogCommentPageURLFormat, page, gid)
}

func _GallogArticleDetailURL(a *ArticleMicroInfo) string {
	return fmt.Sprintf(gallogArticleDetailURLFormat, a.gid, a.cid, a.pno, a.logNo, a.mode)
}

func _GallogCommentDetailURL(a *CommentMicroInfo) string {
	return fmt.Sprintf(gallogCommentDetailURLFormat, a.gid, a.no, a.cno, a.logNo)
}

func _NewGallogDocument(s *Session, URL string) (*goquery.Document, error) {
	resp, err := do("GET", URL, s.cookies, nil, gallogRequestHeader)
	if err != nil {
		log.Fatal(err)
	}
	return goquery.NewDocumentFromResponse(resp)
}

func _CheckResponse(resp *http.Response) error {
	jsonResponse := &_JSONResponse{}
	if err := _ResponseUnmarshal(jsonResponse, resp); err != nil {
		return err
	}
	return _CheckResult(jsonResponse)
}

func _ResponseUnmarshal(data interface{}, resp *http.Response) error {
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	body = []byte(strings.Trim(string(body), "[]"))
	if err := json.Unmarshal(body, data); err != nil {
		return err
	}
	return nil
}

func _CheckResult(jsonResponse *_JSONResponse) error {
	if jsonResponse.Result == false {
		if jsonResponse.Cause != "" {
			return errors.New(jsonResponse.Cause)
		}
		return errUnknownCause
	}
	return nil
}
