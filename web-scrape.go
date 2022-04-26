package main

import (
	"flag"
	"fmt"
	"os"
	"reefwn/web-scrape-cli/cmd"
	"time"
)

func main () {
	imgCmd := flag.NewFlagSet("image", flag.ExitOnError)
	imgUrl := imgCmd.String("u", "", "url to scrape")
	imgAttr := imgCmd.String("a", "src", "attribute to scrape")
	imgFolder := imgCmd.String("f", "images", "folder to store images")

	start := time.Now()

	switch os.Args[1] {
		case "image": 
			cmd.HandleImageCmd(imgCmd, imgUrl, imgAttr, imgFolder)
			break
		default:
			fmt.Println("unrecognized command")
	}

	fmt.Printf("\nTook: %f secs\n", time.Since(start).Seconds())
}
