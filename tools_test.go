package goinside

import "testing"

var testData = []struct {
	argument string
	result   string
}{
	{"http://gall.dcinside.com/board/lists/?id=programming", "http://m.dcinside.com/list.php?id=programming"},
	{"http://gall.dcinside.com/board/view/?id=programming&no=618317&page=1", "http://m.dcinside.com/view.php?id=programming&no=618317"},
	{"http://gall.dcinside.com/board/lists/?id=programming&page=2&exception_mode=best", "http://m.dcinside.com/list.php?id=programming&page=2"},
	{"http://gall.dcinside.com/mgallery/board/lists/?id=seohyunjin", "http://m.dcinside.com/list.php?id=seohyunjin"},
	{"http://gall.dcinside.com/mgallery/board/view/?id=seohyunjin&no=18949&page=1", "http://m.dcinside.com/view.php?id=seohyunjin&no=18949"},
	{"http://gall.dcinside.com/mgallery/board/lists/?id=seohyunjin&page=1&exception_mode=recommend", "http://m.dcinside.com/list.php?id=seohyunjin&page=1"},
	{"http://m.dcinside.com/list.php?id=seohyunjin", "http://m.dcinside.com/list.php?id=seohyunjin"},
	{"http://m.dcinside.com/view.php?id=seohyunjin&no=18954&page=1", "http://m.dcinside.com/view.php?id=seohyunjin&no=18954"},
	{"http://m.dcinside.com/list.php?id=seohyunjin&page=1&exception_mode=recommend", "http://m.dcinside.com/list.php?id=seohyunjin&page=1"},
	{"http://m.dcinside.com/list.php?id=programming&page=2&exception_mode=best", "http://m.dcinside.com/list.php?id=programming&page=2"},
	{"http://m.dcinside.com/view.php?id=programming&no=618286&page=2", "http://m.dcinside.com/view.php?id=programming&no=618286"},
	{"http://m.dcinside.com/", "http://m.dcinside.com/"},
	{"http://m.dcinside.com/category_gall_total.html", "http://m.dcinside.com/category_gall_total.html"},
	{"http://www.dcinside.com/", "http://www.dcinside.com/"},
	{"http://m.dcinside.com/category_gall_total.html", "http://m.dcinside.com/category_gall_total.html"},
	{"http://m.dcinside.com/comment_more_new.php", "http://m.dcinside.com/comment_more_new.php"},
	{"https://dcid.dcinside.com/join/mobile_app_login.php", "https://dcid.dcinside.com/join/mobile_app_login.php"},
	{"http://upload.dcinside.com/_app_write_api.php", "http://upload.dcinside.com/_app_write_api.php"},
	{"http://m.dcinside.com/api/gall_del.php", "http://m.dcinside.com/api/gall_del.php"},
	{"http://m.dcinside.com/api/comment_ok.php", "http://m.dcinside.com/api/comment_ok.php"},
	{"http://m.dcinside.com/api/comment_del.php", "http://m.dcinside.com/api/comment_del.php"},
	{"http://m.dcinside.com/api/_recommend_up.php", "http://m.dcinside.com/api/_recommend_up.php"},
	{"http://m.dcinside.com/api/_recommend_down.php", "http://m.dcinside.com/api/_recommend_down.php"},
	{"http://m.dcinside.com/api/report_upload.php", "http://m.dcinside.com/api/report_upload.php"},
}

func TestToMobileURL(t *testing.T) {
	for _, d := range testData {
		if r := ToMobileURL(d.argument); r != d.result {
			t.Errorf("%v 의 결과값이 %v (이)가 아닙니다. r=%v",
				d.argument, d.result, r)
		}
	}
}
