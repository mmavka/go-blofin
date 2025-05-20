package main

import (
	"flag"
	"log"
	"os"
	"os/exec"
)

func main() {
	example := flag.String("example", "", "Example to run: public, private, ws")
	flag.Parse()

	var cmd *exec.Cmd
	switch *example {
	case "public":
		cmd = exec.Command("go", "run", "public/main.go")
	case "private":
		cmd = exec.Command("go", "run", "private/main.go")
	case "ws":
		cmd = exec.Command("go", "run", "ws/main.go")
	default:
		log.Fatal("Please specify example to run: -example=public|private|ws")
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Dir = "examples"

	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}
}
