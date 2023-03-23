package main

import "log"

const Version = "0.0.1"

func main() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	root := root()
	root.AddCommand(
		server(),
		tunnel(),
	)
	root.Execute()
}
