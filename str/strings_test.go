package str

import (
	"net/http"
	"net/url"
	"testing"
	"time"
)

func TestStringify(t *testing.T) {
	var nilGenre *Genre
	day := Timestamp{Time: time.Date(2026, time.October, 1, 0, 0, 0, 0, time.UTC)}
	cases := []struct {
		name string
		in   any
		want string
	}{
		{name: "string", in: "tron", want: `"tron"`},
		{name: "int", in: 5, want: "5"},
		{name: "nil pointer", in: nilGenre, want: "<nil>"},
		{name: "slice of strings", in: []string{"a", "b"}, want: `["a" "b"]`},
		{name: "struct skips nil fields", in: Genre{Name: String("Action")}, want: `str.Genre{Name:"Action"}`},
		{name: "pointer to struct", in: &Genre{Name: String("Action"), Slug: String("action")}, want: `str.Genre{Name:"Action", Slug:"action"}`},
		{name: "nested struct", in: Movie{Title: String("TRON"), IDs: &IDs{Trakt: Int64(1)}}, want: `str.Movie{Title:"TRON", IDs:str.IDs{Trakt:1}}`},
		{name: "slice of structs", in: []*Genre{{Name: String("Action")}, {Slug: String("drama")}}, want: `[str.Genre{Name:"Action"} str.Genre{Slug:"drama"}]`},
		{name: "timestamp", in: day, want: "str.Timestamp{2026-10-01 00:00:00 +0000 UTC}"},
		{name: "bool pointer field", in: UserProfile{Private: Bool(true), Age: Int(30)}, want: "str.UserProfile{Private:true, Age:30}"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := Stringify(tc.in); got != tc.want {
				t.Errorf("Stringify() is %s, want %s", got, tc.want)
			}
		})
	}
}

func TestStringMethodUsesStringify(t *testing.T) {
	g := Genre{Name: String("Action")}
	if got, want := g.String(), `str.Genre{Name:"Action"}`; got != want {
		t.Errorf("String() is %s, want %s", got, want)
	}
}

func TestContain(t *testing.T) {
	if !ContainString("b", []string{"a", "b"}) {
		t.Error(`ContainString("b") is false, want true`)
	}
	if ContainString("c", []string{"a", "b"}) {
		t.Error(`ContainString("c") is true, want false`)
	}
	if !ContainInt(2, []int{1, 2}) {
		t.Error("ContainInt(2) is false, want true")
	}
	if ContainInt(3, nil) {
		t.Error("ContainInt(3) is true, want false")
	}
}

func TestSliceFlagValues(t *testing.T) {
	s := Slice{}
	for _, v := range []string{"a", "b", "a"} {
		if err := s.Set(v); err != nil {
			t.Fatalf("Set(%q): %v", v, err)
		}
	}
	if got, want := s.String(), "a,b"; got != want {
		t.Errorf("Slice.String() is %q, want %q", got, want)
	}
	if got := (&Slice{}).String(); got != "" {
		t.Errorf("empty Slice.String() is %q, want empty", got)
	}

	si := SliceInt{}
	for _, v := range []int{3, 1, 3} {
		if err := si.Set(v); err != nil {
			t.Fatalf("Set(%d): %v", v, err)
		}
	}
	if got, want := si.String(), "3,1"; got != want {
		t.Errorf("SliceInt.String() is %q, want %q", got, want)
	}
	if got := (&SliceInt{}).String(); got != "" {
		t.Errorf("empty SliceInt.String() is %q, want empty", got)
	}
}

func TestErrorsGetComments(t *testing.T) {
	var empty Errors
	if err := empty.GetComments(); err != nil {
		t.Errorf("GetComments() without comments is %v, want nil", err)
	}

	e := Errors{Comment: &[]string{"must be at least 5 words", "is too long"}}
	err := e.GetComments()
	if err == nil {
		t.Fatal("GetComments() is nil, want an error")
	}
	if got, want := err.Error(), "must be at least 5 words\nis too long"; got != want {
		t.Errorf("GetComments() is %q, want %q", got, want)
	}
}

func TestErrorResponseError(t *testing.T) {
	plain := &ErrorResponse{Message: "boom"}
	if got, want := plain.Error(), `str.ErrorResponse{Message:"boom", ErrorCode:"", Guidance:""}`; got != want {
		t.Errorf("Error() without a response is %s, want %s", got, want)
	}

	u, err := url.Parse("https://api.trakt.tv/oauth/token?client_secret=hidden")
	if err != nil {
		t.Fatalf("parse url: %v", err)
	}
	withResponse := &ErrorResponse{
		Message:  "boom",
		Response: &http.Response{StatusCode: http.StatusTeapot, Request: &http.Request{Method: http.MethodPost, URL: u}},
	}
	if got, want := withResponse.Error(), "POST https://api.trakt.tv/oauth/token?client_secret=REDACTED: 418 boom"; got != want {
		t.Errorf("Error() is %s, want %s", got, want)
	}
	if got, want := u.Query().Get("client_secret"), "hidden"; got != want {
		t.Errorf("request url secret is %q after Error(), want it untouched", got)
	}
}
