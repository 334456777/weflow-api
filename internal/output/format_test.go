package output

import "testing"

func TestFormatContent_RevokeMsg(t *testing.T) {
	input := `<?xml version="1.0"?><sysmsg type="revokemsg"><revokemsg><content>"星魂" 撤回了一条消息</content><revoketime>0</revoketime></revokemsg></sysmsg>`
	want := `["星魂" 撤回了一条消息]`
	got := FormatContent(input)
	if got != want {
		t.Errorf("FormatContent(revoke) = %q, want %q", got, want)
	}
}

func TestFormatContent_PlainText(t *testing.T) {
	got := FormatContent("hello world")
	if got != "hello world" {
		t.Errorf("FormatContent(plain) = %q, want %q", got, "hello world")
	}
}

func TestFormatContent_AppShare(t *testing.T) {
	input := `<msg><appmsg><title>测试标题</title><url>https://example.com</url></appmsg><appinfo><appname>测试App</appname></appinfo></msg>`
	got := FormatContent(input)
	want := "【测试App】测试标题 https://example.com"
	if got != want {
		t.Errorf("FormatContent(appshare) = %q, want %q", got, want)
	}
}
