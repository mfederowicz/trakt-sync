package str

import (
	"testing"
)

func TestHaveIdFiled(t *testing.T) {
	var myInt int64 = 12345678

	in := &IDs{Trakt: &myInt}

	// test that Trakt id is exist
	if got, want := in.HaveID("Trakt"), true; got != want {
		t.Errorf("in.HaveId('Trakt') is %v, want %v", in.HaveID("Trakt"), want)
	}

	// test that Imdb id is not exist
	if got, want := in.HaveID("Imdb"), false; got != want {
		t.Errorf("in.HaveId('Imdb') is %v, want %v", in.HaveID("Imdb"), want)
	}
}
