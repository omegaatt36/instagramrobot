package main

import (
	"github.com/feelthecode/instagramrobot/src/instagram"
)

func main() {
	code := "CSft2G5pFgr"

	ig := instagram.API{}

	response, err := ig.GetPostWithCode(code)
	if err != nil {
		return
	}

	_ = response
}
