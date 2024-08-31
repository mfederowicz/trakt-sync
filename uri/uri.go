// Package uri used for url operations
package uri

import (
	"fmt"
	"net/url"
	"reflect"
	"sort"
	"strconv"
	"strings"
)

// config slices
var (
	StatusOptions = []string{"running series", "continuing", "in production", "planned", "upcoming", "pilot", "canceled", "ended"}
	EpisodeTypes  = []string{"standard", "series_premiere", "season_premiere", "mid_season_finale", "mid_season_premiere", "season_finale", "series_finale"}
)

// Pagination represents pagination params
type Pagination struct {
	// For paginated result sets, page of results to retrieve.
	Page int `url:"page,omitempty"`
	// For paginated result sets, the number of elements on one page.
	Limit int `url:"limit,omitempty"`
}

// RatingRange represents min/max int parameters
type RatingRange struct {
	Min int `url:"min,omitempty"`
	Max int `url:"max,omitempty"`
}

func (rr RatingRange) String() string {
	if rr.Min <= 0 || rr.Max > 100 {
		return ""
	}

	if rr.Min > rr.Max {
		return ""
	}
	return fmt.Sprintf("%d-%d", rr.Min, rr.Max)
}

// RatingRangeFloat represents min/max float parameters
type RatingRangeFloat struct {
	Min float32 `url:"min,omitempty"`
	Max float32 `url:"max,omitempty"`
}

func (rr RatingRangeFloat) String() string {
	if rr.Min <= 0.0 || rr.Max > 100.0 {
		return ""
	}

	if rr.Min > rr.Max {
		return ""
	}
	return fmt.Sprintf("%.1f-%.1f", rr.Min, rr.Max)

}

// VotesRange represents min/max int votes parameters
type VotesRange struct {
	Min int `url:"min,omitempty"`
	Max int `url:"max,omitempty"`
}

func (r VotesRange) String() string {
	if r.Min <= 0 || r.Max > 100000 {
		return ""
	}

	if r.Min > r.Max {
		return ""
	}

	return fmt.Sprintf("%d-%d", r.Min, r.Max)
}

// ImdbVotesRange represents min/max int imdb votes parameters
type ImdbVotesRange struct {
	Min int `url:"min,omitempty"`
	Max int `url:"max,omitempty"`
}

func (r ImdbVotesRange) String() string {
	if r.Min <= 0 || r.Max > 3000000 {
		return ""
	}

	if r.Min > r.Max {
		return ""
	}

	return fmt.Sprintf("%d-%d", r.Min, r.Max)
}

// TmdbRatingRange represents min/max float tmdb rating parameters
type TmdbRatingRange struct {
	Min float32 `url:"min,omitempty"`
	Max float32 `url:"max,omitempty"`
}

func (r TmdbRatingRange) String() string {
	if r.Min < 0.0 || r.Max > 10.0 {
		return ""
	}

	if r.Min > r.Max || r.Min == r.Max {
		return ""
	}

	return fmt.Sprintf("%.1f-%.1f", r.Min, r.Max)
}

// ListOptions specifies the optional parameters to various List methods that
// support offset pagination.
type ListOptions struct {
	Field          string           `url:"field,omitempty"`
	Runtimes       string           `url:"runtimes,omitempty"`
	Query          string           `url:"query,omitempty"`
	Type           string           `url:"type,omitempty"`
	Extended       string           `url:"extended,omitempty"`
	Years          string           `url:"years,omitempty"`
	Certifications []string         `url:"certifications,omitempty"`
	Genres         []string         `url:"genres,omitempty"`
	EpisodeTypes   []string         `url:"episode_types,omitempty"`
	Countries      []string         `url:"countries,omitempty"`
	Languages      []string         `url:"languages,omitempty"`
	StudioIDs      []int            `url:"studio_ids,omitempty"`
	Status         []string         `url:"status,omitempty"`
	NetworkIDs     []int            `url:"network_ids,omitempty"`
	Ratings        RatingRange      `url:"ratings,omitempty"`
	TmdbVotes      VotesRange       `url:"tmdb_votes,omitempty"`
	ImdbRatings    RatingRange      `url:"imdb_ratings,omitempty"`
	ImdbVotes      ImdbVotesRange   `url:"imdb_votes,omitempty"`
	RtMeters       RatingRange      `url:"rt_meters,omitempty"`
	RtUserMeters   RatingRange      `url:"rt_user_meters,omitempty"`
	Votes          VotesRange       `url:"votes,omitempty"`
	Page           int              `url:"page,omitempty"`
	Limit          int              `url:"limit,omitempty"`
	Metascores     RatingRangeFloat `url:"metascores,omitempty"`
	TmdbRatings    TmdbRatingRange  `url:"tmdb_ratings,omitempty"`
}

// AddQuery adds query parameters to s.
func AddQuery(s string, opts interface{}) (string, error) {

	u, err := url.Parse(s)
	if err != nil {
		return s, err
	}

	// Use reflection to flatten ListOptions fields
	v := reflect.ValueOf(opts)
	qs := url.Values{}

	flatOptsStruct(v, &qs)

	u.RawQuery = EncodeParams(qs)
	return u.String(), nil

}

// CustomTypeHandler defines the function signature for handling custom types
type CustomTypeHandler func(reflect.Value, *url.Values, string)

// customTypeHandlers maps custom types to their corresponding handling functions
var customTypeHandlers = map[reflect.Type]CustomTypeHandler{
	reflect.TypeOf(RatingRange{}):      handleRatingRange,
	reflect.TypeOf(VotesRange{}):       handleVotesRange,
	reflect.TypeOf(TmdbRatingRange{}):  handleTmdbRatingRange,
	reflect.TypeOf(ImdbVotesRange{}):   handleImdbVotesRange,
	reflect.TypeOf(RatingRangeFloat{}): handleMetaCriticRange,
}

// handleFloatRange handles the FloatRange custom type
func handleRatingRange(fieldValue reflect.Value, qs *url.Values, fieldTag string) {
	rr := fieldValue.Interface().(RatingRange)

	if fieldTag != "" && len(rr.String()) > 0 {
		// Remove omitempty tag from the field tag
		fieldTag = strings.Split(fieldTag, ",")[0]
		qs.Add(fieldTag, rr.String())
	}

}

// handleVotesRange handles the VotesRange custom type
func handleVotesRange(fieldValue reflect.Value, qs *url.Values, fieldTag string) {
	rr := fieldValue.Interface().(VotesRange)

	if fieldTag != "" && len(rr.String()) > 0 {
		// Remove omitempty tag from the field tag
		fieldTag = strings.Split(fieldTag, ",")[0]
		qs.Add(fieldTag, rr.String())
	}

}

// handleTmdbRatingRange handles the TmdbRatingRange custom type
func handleTmdbRatingRange(fieldValue reflect.Value, qs *url.Values, fieldTag string) {
	rr := fieldValue.Interface().(TmdbRatingRange)

	if fieldTag != "" && len(rr.String()) > 0 {
		// Remove omitempty tag from the field tag
		fieldTag = strings.Split(fieldTag, ",")[0]
		qs.Add(fieldTag, rr.String())
	}

}

// handleImdbVotesRange handles the ImdbVotesRange custom type
func handleImdbVotesRange(fieldValue reflect.Value, qs *url.Values, fieldTag string) {
	rr := fieldValue.Interface().(ImdbVotesRange)

	if fieldTag != "" && len(rr.String()) > 0 {
		// Remove omitempty tag from the field tag
		fieldTag = strings.Split(fieldTag, ",")[0]
		qs.Add(fieldTag, rr.String())
	}

}

// handleMetaCriticRange handles the RatingRangeFloat custom type
func handleMetaCriticRange(fieldValue reflect.Value, qs *url.Values, fieldTag string) {
	rr := fieldValue.Interface().(RatingRangeFloat)

	if fieldTag != "" && len(rr.String()) > 0 {
		// Remove omitempty tag from the field tag
		fieldTag = strings.Split(fieldTag, ",")[0]
		qs.Add(fieldTag, rr.String())
	}

}

func flatOptsStruct(v reflect.Value, qs *url.Values) error {

	// Check if the value is a pointer
	if v.Kind() == reflect.Ptr {
		// Dereference the pointer to get the underlying value
		v = v.Elem()
	}

	for i := 0; i < v.NumField(); i++ {
		fieldValue := v.Field(i)
		fieldType := v.Type().Field(i)

		// Check if the field value is a custom type
		if handler, ok := customTypeHandlers[fieldValue.Type()]; ok {
			fieldTag := fieldType.Tag.Get("url")
			handler(fieldValue, qs, fieldTag)
			continue
		}

		// Process struct types
		if fieldValue.Kind() == reflect.Struct {
			if err := flatOptsStruct(fieldValue, qs); err != nil {
				return err
			}
		} else {
			fieldTag := fieldType.Tag.Get("url")
			if fieldTag != "" {
				// Remove omitempty tag from the field tag
				fieldTag = strings.Split(fieldTag, ",")[0]

				// Convert field value to string based on its type
				var value string
				switch fieldValue.Kind() {
				case reflect.Slice, reflect.Array:
					if fieldValue.Len() == 0 {
						value = "" // If the slice is empty, return an empty string
					}
					var values []string
					for i := 0; i < fieldValue.Len(); i++ {
						switch fieldValue.Index(i).Kind() {
						case reflect.String:
							values = append(values, fieldValue.Index(i).String())
						case reflect.Int:
							values = append(values, strconv.FormatInt(fieldValue.Index(i).Int(), 10))
						}
					}
					// Join the string slice into a single comma-separated string
					if len(values) > 0 {
						value = strings.Join(values, ",")
					} else {
						value = ""
					}

				case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
					value = fmt.Sprintf("%d", fieldValue.Int())
				case reflect.Float32, reflect.Float64:
					value = fmt.Sprintf("%.2f", fieldValue.Float())
				case reflect.String:
					value = fieldValue.String()
				}

				// Add field to query string only if it's non-empty
				if value != "" && !isEmptyValue(fieldValue) {
					qs.Add(fieldTag, value)
				}
			}
		}
	}
	return nil
}

// isEmptyValue checks if a reflect.Value represents the zero value of its type
func isEmptyValue(v reflect.Value) bool {
	zero := reflect.Zero(v.Type())
	return reflect.DeepEqual(v.Interface(), zero.Interface())
}

// EncodeParams encodes the values for query sorted by key
func EncodeParams(values url.Values) string {

	if len(values) == 0 {
		return ""
	}

	var buf strings.Builder
	keys := make([]string, 0, len(values))
	for k := range values {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		vs := values[k]
		for _, v := range vs {
			if buf.Len() > 0 {
				buf.WriteByte('&')
			}
			buf.WriteString(k)
			buf.WriteByte('=')
			buf.WriteString(v)
		}
	}
	return buf.String()
}

func hasStruct(v reflect.Value) bool {
	for i := 0; i < v.NumField(); i++ {
		if v.Field(i).Kind() == reflect.Struct {
			return true
		}
	}
	return false
}

// SanitizeURL redacts the client_secret parameter from the URL which may be
// exposed to the user.
func SanitizeURL(uri *url.URL) *url.URL {
	if uri == nil {
		return nil
	}
	params := uri.Query()
	if len(params.Get("client_secret")) > 0 {
		params.Set("client_secret", "REDACTED")
		uri.RawQuery = params.Encode()
	}
	return uri
}
