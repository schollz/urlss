package main

import (
	"errors"
	"flag"
	"fmt"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/gin-contrib/multitemplate"
	"github.com/gin-gonic/gin"
	"github.com/goware/urlx"
	"github.com/schollz/jsonstore"
)

var ks *jsonstore.JSONStore

func init() {
	var err error
	ks, err = jsonstore.Open("urls.json.gz")
	if err != nil {
		ks = new(jsonstore.JSONStore)
	}
}

var Port string

func main() {
	gin.SetMode(gin.ReleaseMode)
	flag.StringVar(&Port, "p", "8006", "port (default 8006)")
	flag.Parse()
	r := gin.Default()
	r.Use(gin.Logger())
	r.HTMLRender = loadTemplates("index.html")
	r.GET("/*action", handleAction)
	// Start server
	fmt.Println("Listening on port", Port)
	r.Run(":" + Port) // listen and serve on 0.0.0.0:8080
}

// handleAction performs the shortening or redirecting
func handleAction(c *gin.Context) {
	fmt.Println(c.Request.RequestURI)
	action := c.Request.RequestURI
	action = action[1:len(action)]
	shortened, redirect, err := shortenURL(action)
	if redirect {
		c.Redirect(301, shortened)
	} else {
		errString := ""
		if err != nil {
			errString = err.Error()
		}
		c.HTML(http.StatusOK, "index.html", gin.H{
			"shortened": shortened,
			"error":     errString,
		})
	}
}

func shortenURL(requestURL string) (shortened string, redirect bool, err error) {
	if strings.Contains(requestURL, "http") && !strings.Contains(requestURL, "//") {
		requestURL = strings.Replace(requestURL, "/", "//", 1)
	}
	parsedURL, _ := urlx.Parse(requestURL)
	url, _ := urlx.Normalize(parsedURL)
	if len(url) > 0 && !strings.Contains(url, "favicon") {
		// Check if it is already a URL
		errFound := ks.Get(url, &shortened)
		if errFound != nil {
			// Get a new shortend URL
			shortened = newShortenedURL()
			ks.Set(url, shortened)
			ks.Set(shortened, url)
			go jsonstore.Save(ks, "urls.json.gz")
			log.Printf("Shortened %s to %s", url, shortened)
		}
	} else {
		// Redirect the URL if it is shortened
		err = ks.Get(requestURL, &shortened)
		if err == nil {
			redirect = true
			log.Printf("Redirect %s to %s", requestURL, shortened)
		} else {
			if requestURL == "" {
				err = nil
			} else {
				err = errors.New("Could not find " + requestURL)
			}
		}
	}
	return
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

// newShortenedURL generates a shortened URL,
// stochastically, checking for collisions until
// a selects a free one
func newShortenedURL() string {
	for n := 1; n < 10; n++ {
		for i := 0; i < 10; i++ {
			candidate := RandString(n)
			var foo string
			err := ks.Get(candidate, &foo)
			if err != nil {
				return candidate
			}
		}
	}
	return ""
}

// loadTemplates will use the built-in assets to
// load required templates
func loadTemplates(list ...string) multitemplate.Render {
	r := multitemplate.New()

	for _, x := range list {
		templateString, err := Asset("templates/" + x)
		if err != nil {
			panic(err)
		}

		tmplMessage, err := template.New(x).Parse(string(templateString))
		if err != nil {
			panic(err)
		}

		r.Add(x, tmplMessage)
	}

	return r
}
