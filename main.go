package main

import "music-subscription/music"

func main()  {
	var a music.App
	a.Initialize("root", "", "bcc_music")
	a.Run(":8080")
}
