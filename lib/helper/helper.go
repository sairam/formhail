package helper

import (
	"encoding/json"
	"io"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"time"
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

// PostJSONtoURL uses a zero buffer mechanism to make a http.Post call as json
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

// GetFirstValue gets the first value from the array of options by key
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

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

// RandString generates random alphabets
// source: http://stackoverflow.com/questions/22892120/how-to-generate-a-random-string-of-a-fixed-length-in-golang
// RandStringBytesMaskImprSrc
func RandString(n int) string {
	var src = rand.NewSource(time.Now().UnixNano())
	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return string(b)
}
