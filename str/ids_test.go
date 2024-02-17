package str

import (
	"testing"
)

func TestHaveIdFiled(t *testing.T) {

	var myInt int64 = 12345678

	in := &Ids{Trakt:&myInt}

	
	// test that Trakt id is exist 
	if got, want := in.HaveId("Trakt"), true; got != want {
		t.Errorf("in.HaveId('Trakt') is %v, want %v",in.HaveId("Trakt"), want)
	}

	// test that Imdb id is not exist 
	if got, want := in.HaveId("Imdb"), false; got != want {
		t.Errorf("in.HaveId('Imdb') is %v, want %v",in.HaveId("Imdb"), want)
	}


	
}

