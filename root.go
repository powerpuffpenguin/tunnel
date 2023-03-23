package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/spf13/cobra"
)

func root() *cobra.Command {
	var version bool
	cmd := &cobra.Command{
		Use:   "tunnel",
		Short: "tunnel port mapping between server and client",
		Run: func(cmd *cobra.Command, args []string) {
			if version {
				fmt.Println(Version)
				return
			}
			fmt.Println(`tunnel port mapping between server and client`)
			fmt.Println(runtime.GOOS, runtime.GOARCH, runtime.Version())
			fmt.Println(`version`, Version)
		},
	}
	flags := cmd.Flags()
	flags.BoolVarP(&version, `version`, `v`, false, `print version`)
	return cmd
}

func BasePath() string {
	filename, e := exec.LookPath(os.Args[0])
	if e != nil {
		log.Fatalln(e)
	}

	filename, e = filepath.Abs(filename)
	if e != nil {
		log.Fatalln(e)
	}
	return filepath.Dir(filename)
}
