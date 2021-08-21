package main

import (
	"fmt"

	"github.com/feelthecode/instagramrobot/src/config"
	"github.com/feelthecode/instagramrobot/src/instagram"
	"github.com/feelthecode/instagramrobot/src/utils"
)

func main() {
	utils.RegisterLogger()
	config.Load()

	code := "CSft2G5pFgr"

	ig := instagram.API{}

	response, err := ig.GetPostWithCode(code)
	if err != nil {
		fmt.Print(err.Error())
		return
	}

	fmt.Printf("%+v\n", response)
}
