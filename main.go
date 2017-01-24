package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

var data = struct {
	sync.RWMutex
	URLToString map[string]string
	StringToURL map[string]string
}{URLToString: make(map[string]string), StringToURL: make(map[string]string)}

func loadData() {
	if _, err := os.Stat("urltostring.json"); !os.IsNotExist(err) {
		// path/to/whatever does exist
		b, _ := ioutil.ReadFile("urltostring.json")
		data.Lock()
		json.Unmarshal(b, &data.URLToString)
		data.Unlock()
	}
	if _, err := os.Stat("stringtourl.json"); !os.IsNotExist(err) {
		// path/to/whatever does exist
		b, _ := ioutil.ReadFile("stringtourl.json")
		data.Lock()
		json.Unmarshal(b, &data.StringToURL)
		data.Unlock()
	}
}

func saveData() {
	data.RLock()
	defer data.RUnlock()

	b, _ := json.MarshalIndent(data.StringToURL, "", " ")
	ioutil.WriteFile("stringtourl.json", b, 0644)
	b, _ = json.MarshalIndent(data.URLToString, "", " ")
	ioutil.WriteFile("urltostring.json", b, 0644)
}

func init() {
	loadData()
}

func getShortenedURL() string {
	data.RLock()
	defer data.RUnlock()
	for n := 2; n < 10; n++ {
		for i := 0; i < 10; i++ {
			candidate := RandString(n)
			_, ok := data.StringToURL[candidate]
			if !ok {
				return candidate
			}
		}
	}
	return ""
}

var Host, Port string

func main() {
	gin.SetMode(gin.ReleaseMode)
	flag.StringVar(&Host, "h", "", "host (optional)")
	flag.StringVar(&Port, "p", "8006", "port (default 8006)")
	flag.Parse()
	r := gin.Default()
	r.GET("/*action", func(c *gin.Context) {
		action := c.Param("action")
		action = action[1:len(action)]
		if strings.Contains(action, "http") {
			// Save the URL
			data.RLock()
			shortened, ok := data.URLToString[action]
			data.RUnlock()
			if ok {
				c.String(http.StatusOK, "Already shortened %s as %s/%s", action, Host, shortened)
				return
			}

			// Get a new shortend URL
			shortened = getShortenedURL()
			data.Lock()
			data.StringToURL[shortened] = action
			data.URLToString[action] = shortened
			data.Unlock()

			go saveData()
			c.String(http.StatusOK, "Generated new URL from %s: %s/%s", action, Host, shortened)
			return

		} else {
			// Redirect the URL if it is shortened
			data.RLock()
			url, ok := data.StringToURL[action]
			data.RUnlock()
			if ok {
				c.Redirect(301, url)
				return
			} else {
				if action == "" {
					c.String(http.StatusOK, "Usage:\n\n/https://.... to store\n\n/... to redirect")
				} else {
					c.String(http.StatusOK, "Could not find %s", action)
				}
			}
		}
	})
	r.Run(":" + Port) // listen and serve on 0.0.0.0:8080
}

// From http://stackoverflow.com/questions/22892120/how-to-generate-a-random-string-of-a-fixed-length-in-golang
const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

var src = rand.NewSource(time.Now().UnixNano())

func RandString(n int) string {
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
