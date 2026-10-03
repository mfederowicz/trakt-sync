package str

import (
	"encoding/json"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestUserNotesUnmarshalJSON(t *testing.T) {
	sean := &UserProfile{Username: String("sean"), Private: Bool(false), IDs: &IDs{Slug: String("sean")}}
	tests := []struct {
		name string
		json string
		want UserNotes
	}{
		{
			name: "nested user",
			json: `{"user":{"username":"sean","private":false,"ids":{"slug":"sean"}},"notes":"watch it"}`,
			want: UserNotes{User: sean, Notes: String("watch it")},
		},
		{
			name: "flat profile",
			json: `{"username":"sean","private":false,"ids":{"slug":"sean"},"notes":"watch it"}`,
			want: UserNotes{User: sean, Notes: String("watch it")},
		},
		{
			name: "flat profile without notes",
			json: `{"username":"sean","private":false,"ids":{"slug":"sean"}}`,
			want: UserNotes{User: sean},
		},
		{
			name: "nested user wins over flat keys",
			json: `{"user":{"username":"sean","private":false,"ids":{"slug":"sean"}},"username":"other","notes":"watch it"}`,
			want: UserNotes{User: sean, Notes: String("watch it")},
		},
		{
			name: "notes only",
			json: `{"notes":"watch it"}`,
			want: UserNotes{Notes: String("watch it")},
		},
		{
			name: "null user",
			json: `{"user":null,"notes":"watch it"}`,
			want: UserNotes{Notes: String("watch it")},
		},
		{
			name: "empty object",
			json: `{}`,
			want: UserNotes{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := UserNotes{}
			if err := json.Unmarshal([]byte(tt.json), &got); err != nil {
				t.Fatalf("unmarshal error: %v", err)
			}
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("user notes mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestUserNotesUnmarshalJSONErrors(t *testing.T) {
	tests := []struct {
		name string
		json string
	}{
		{name: "not an object", json: `"sean"`},
		{name: "wrong type in a flat profile", json: `{"username":1,"notes":"watch it"}`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := UserNotes{}
			if err := json.Unmarshal([]byte(tt.json), &got); err == nil {
				t.Errorf("expected an error, got %v", got)
			}
		})
	}
}

func TestRecommendationUserNotesShapes(t *testing.T) {
	data := `{"title":"Andor","favorited_by":[{"username":"sean","notes":"flat"}],"recommended_by":[{"user":{"username":"justin"},"notes":"nested"},null]}`
	got := Recommendation{}
	if err := json.Unmarshal([]byte(data), &got); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}
	want := Recommendation{
		Title:         String("Andor"),
		FavoritedBy:   &[]UserNotes{{User: &UserProfile{Username: String("sean")}, Notes: String("flat")}},
		RecommendedBy: &[]UserNotes{{User: &UserProfile{Username: String("justin")}, Notes: String("nested")}, {}},
	}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("recommendation mismatch (-want +got):\n%s", diff)
	}

	out, err := json.Marshal(got.FavoritedBy)
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}
	if want := `[{"user":{"username":"sean"},"notes":"flat"}]`; string(out) != want {
		t.Errorf("marshalled %s, want %s", out, want)
	}
}
