package main

import (
	"encoding/json"
	"errors"
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

type Data struct {
	URLToString map[string]string `json:"url_to_string"`
	StringToURL map[string]string `json:"string_to_url"`
}

type Filesystem struct {
	data Data
	sync.RWMutex
}

func (s *Filesystem) Init() {
	s.Load()
}

func (s *Filesystem) Load() {
	s.Lock()
	defer s.Unlock()
	if _, err := os.Stat("urls.json"); !os.IsNotExist(err) {
		b, _ := ioutil.ReadFile("urls.json")
		json.Unmarshal(b, &s.data)
	} else {
		s.data.URLToString = make(map[string]string)
		s.data.StringToURL = make(map[string]string)
	}
}

func (s *Filesystem) Save(url string, short string) {
	s.Lock()
	s.data.URLToString[url] = short
	s.data.StringToURL[short] = url
	b, _ := json.MarshalIndent(s.data, "", " ")
	ioutil.WriteFile("urls.json", b, 0644)
	s.Unlock()
}

func (s *Filesystem) GetShortFromURL(url string) (string, error) {
	s.RLock()
	defer s.RUnlock()
	val, ok := s.data.URLToString[url]
	if !ok {
		return "", errors.New("URL does not exist")
	} else {
		return val, nil
	}
}

func (s *Filesystem) GetURLFromShort(short string) (string, error) {
	s.RLock()
	defer s.RUnlock()
	val, ok := s.data.StringToURL[short]
	if !ok {
		return "", errors.New("short does not exist")
	} else {
		return val, nil
	}
}

var fs Filesystem

func init() {
	fs.Init()
}

func newShortenedURL() string {
	for n := 2; n < 10; n++ {
		for i := 0; i < 10; i++ {
			candidate := RandString(n)
			_, err := fs.GetURLFromShort(candidate)
			if err != nil {
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
			url := action
			if !strings.Contains(url, "//") {
				url = strings.Replace(url, "/", "//", 1)
			}

			// Save the URL
			shortened, err := fs.GetShortFromURL(url)
			if err == nil {
				c.String(http.StatusOK, "Already shortened %s as %s/%s", url, Host, shortened)
				return
			}

			// Get a new shortend URL
			shortened = newShortenedURL()
			fs.Save(url, shortened)

			c.String(http.StatusOK, "Generated new URL from %s: %s/%s", url, Host, shortened)
			return
		} else {
			// Redirect the URL if it is shortened
			url, err := fs.GetURLFromShort(action)
			if err == nil {
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
