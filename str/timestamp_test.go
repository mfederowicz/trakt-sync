package str

import (
	"encoding/json"
	"testing"
	"time"
)

func TestTimestampJSON(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want time.Time
		out  string
	}{
		{name: "date and time", in: `"2026-10-01T14:07:14Z"`, want: time.Date(2026, time.October, 1, 14, 7, 14, 0, time.UTC), out: `"2026-10-01T14:07:14Z"`},
		{name: "milliseconds are dropped on encode", in: `"2026-10-01T14:07:14.000Z"`, want: time.Date(2026, time.October, 1, 14, 7, 14, 0, time.UTC), out: `"2026-10-01T14:07:14Z"`},
		{name: "date only", in: `"2026-10-01"`, want: time.Date(2026, time.October, 1, 0, 0, 0, 0, time.UTC), out: `"2026-10-01"`},
		{name: "midnight keeps its time", in: `"2026-10-01T00:00:00Z"`, want: time.Date(2026, time.October, 1, 0, 0, 0, 0, time.UTC), out: `"2026-10-01T00:00:00Z"`},
		{name: "midnight with an offset keeps its time and offset", in: `"2026-10-01T00:00:00+02:00"`, want: time.Date(2026, time.September, 30, 22, 0, 0, 0, time.UTC), out: `"2026-10-01T00:00:00+02:00"`},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var ts Timestamp
			if err := json.Unmarshal([]byte(tc.in), &ts); err != nil {
				t.Fatalf("decode %s: %v", tc.in, err)
			}
			if !ts.Equal(tc.want) {
				t.Errorf("decoded time is %v, want %v", ts.Time, tc.want)
			}
			out, err := json.Marshal(ts)
			if err != nil {
				t.Fatalf("encode: %v", err)
			}
			if string(out) != tc.out {
				t.Errorf("encoded time is %s, want %s", out, tc.out)
			}
		})
	}
}

// TestTimestampTimezone checks the encoded form of values moved to another timezone, as the client does with WithTimezone.
func TestTimestampTimezone(t *testing.T) {
	warsaw := time.FixedZone("CEST", 2*60*60)
	newYork := time.FixedZone("EDT", -4*60*60)
	cases := []struct {
		name string
		in   string
		loc  *time.Location
		out  string
	}{
		{name: "local midnight east of UTC keeps its time", in: `"2026-09-30T22:00:00Z"`, loc: warsaw, out: `"2026-10-01T00:00:00+02:00"`},
		{name: "local midnight west of UTC keeps its time", in: `"2026-10-01T04:00:00Z"`, loc: newYork, out: `"2026-10-01T00:00:00-04:00"`},
		{name: "UTC midnight is a local time", in: `"2026-10-01T00:00:00Z"`, loc: warsaw, out: `"2026-10-01T02:00:00+02:00"`},
		{name: "date only east of UTC stays the same day", in: `"2026-10-01"`, loc: warsaw, out: `"2026-10-01"`},
		{name: "date only west of UTC stays the same day", in: `"2026-10-01"`, loc: newYork, out: `"2026-10-01"`},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var ts Timestamp
			if err := json.Unmarshal([]byte(tc.in), &ts); err != nil {
				t.Fatalf("decode %s: %v", tc.in, err)
			}
			ts.Time = ts.Time.In(tc.loc)
			out, err := json.Marshal(ts)
			if err != nil {
				t.Fatalf("encode: %v", err)
			}
			if string(out) != tc.out {
				t.Errorf("encoded time is %s, want %s", out, tc.out)
			}
		})
	}
}

// TestTimestampRoundTrip checks that a value encoded and decoded again is the same moment.
func TestTimestampRoundTrip(t *testing.T) {
	warsaw := time.FixedZone("CEST", 2*60*60)
	for _, want := range []time.Time{
		time.Date(2026, time.October, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2026, time.October, 1, 0, 0, 0, 0, warsaw),
		time.Date(2026, time.October, 1, 14, 7, 14, 0, warsaw),
	} {
		out, err := json.Marshal(Timestamp{Time: want})
		if err != nil {
			t.Fatalf("encode %v: %v", want, err)
		}
		var got Timestamp
		if err := json.Unmarshal(out, &got); err != nil {
			t.Fatalf("decode %s: %v", out, err)
		}
		if !got.Equal(want) {
			t.Errorf("%s decodes to %v, want %v", out, got.Time, want)
		}
	}
}

func TestTimestampInvalid(t *testing.T) {
	var ts Timestamp
	err := json.Unmarshal([]byte(`"yesterday"`), &ts)
	if err == nil {
		t.Fatal("error is nil, want an invalid format error")
	}
	if got, want := err.Error(), "invalid timestamp format: yesterday"; got != want {
		t.Errorf("error is %q, want %q", got, want)
	}
}

func TestTimestampUTC(t *testing.T) {
	warsaw := time.FixedZone("CEST", 2*60*60)
	ts := Timestamp{Time: time.Date(2026, time.October, 1, 16, 7, 14, 0, warsaw)}

	got := ts.UTC()
	if got.Location() != time.UTC {
		t.Errorf("location is %v, want UTC", got.Location())
	}
	if got.Hour() != 14 {
		t.Errorf("hour is %d, want 14", got.Hour())
	}
	if ts.Location() != warsaw {
		t.Errorf("original location changed to %v", ts.Location())
	}
}

func TestToken(t *testing.T) {
	const created = 1790863634
	device := &DeviceToken{
		AccessToken:  String("access"),
		TokenType:    String("bearer"),
		ExpiresIn:    Int64(7200),
		RefreshToken: String("refresh"),
		Scope:        String("public"),
		CreatedAt:    Int64(created),
	}
	token := device.ToToken()
	want := Token{AccessToken: "access", TokenType: "bearer", RefreshToken: "refresh", Scope: "public", ExpiresIn: 7200, CreatedAt: created}
	if *token != want {
		t.Errorf("ToToken() is %+v, want %+v", *token, want)
	}
	if got, wantPoint := token.ExpirationPoint(), time.Unix(created+7200, 0); !got.Equal(wantPoint) {
		t.Errorf("ExpirationPoint() is %v, want %v", got, wantPoint)
	}

	now := time.Now().Unix()
	old := &Token{CreatedAt: now - 7200, ExpiresIn: 3600}
	if !old.Expired() {
		t.Error("Expired() is false for a token that ended an hour ago")
	}
	if got := old.ExpiritySeconds(); got >= 0 {
		t.Errorf("ExpiritySeconds() is %d for an expired token, want a negative value", got)
	}

	fresh := &Token{CreatedAt: now, ExpiresIn: 3600}
	if fresh.Expired() {
		t.Error("Expired() is true for a token valid for an hour")
	}
	if got := fresh.ExpiritySeconds(); got < 3500 || got > 3600 {
		t.Errorf("ExpiritySeconds() is %d, want about 3600", got)
	}
}
