package main

import (
	"os"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestHandleAction(t *testing.T) {
	os.Remove("urls.json.gz")
	shortened, redirect, err := shortenURL("www.google.com")
	if len(shortened) != 1 {
		t.Errorf("Got %s for some reason", shortened)
	}
	if redirect == true {
		t.Error("Incorrectly redirecting")
	}
	if err != nil {
		t.Error(err)
	}

	s := newShortenedURL()
	if s == shortened {
		t.Error("New shortened URL should be different")
	}

	shortened, redirect, err = shortenURL(shortened)
	if shortened != "http://www.google.com" {
		t.Errorf("Got %s for some reason", shortened)
	}
	if redirect == false {
		t.Error("Incorrectly NOT redirecting")
	}
	if err != nil {
		t.Error(err)
	}
	shortened, redirect, err = shortenURL("asldkfjaslkdf")
	if err == nil {
		t.Error("Should throw error!")
	}

}

func TestUtils(t *testing.T) {
	if len(RandString(10)) != 10 {
		t.Error("RandString is weird")
	}
	if RandString(3) == RandString(3) {
		t.Error("RandString should be different")
	}

	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.Use(gin.Logger())
	r.HTMLRender = loadTemplates("index.html")
	r.GET("/*action", handleAction)
}
