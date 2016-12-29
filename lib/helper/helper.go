package helper

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"
)

// StringIn finds a key in a list of strings
func StringIn(list []string, key string) bool {
	for _, k := range list {
		if key == k {
			return true
		}
	}
	return false
}

// PostJSONtoURL ..
func PostJSONtoURL(url string, data interface{}) error {
	pr, pw := io.Pipe()
	go func() {
		// close the writer, so the reader knows there's no more data
		defer pw.Close()

		// write json data to the PipeReader through the PipeWriter
		if err := json.NewEncoder(pw).Encode(data); err != nil {
			log.Print(err)
		}
	}()

	if _, err := http.Post(url, "application/json", pr); err != nil {
		return err
	}
	return nil
}

// GetFirstValue ..
func GetFirstValue(q url.Values, key string) string {
	if len(q[key]) == 0 {
		return ""
	}
	return q[key][0]
}

// Min ..
func Min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
