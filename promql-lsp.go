package main


import (
	"bufio"
	"log"
	"os"
	"golang.org/x/tools/internal/jsonrpc2"
)

func main() {


	conn := jsonrpc2.NewConn(jsonrpc2.NewStream(os.Stdin, os.Stdout))

	
	s := bufio.NewScanner(os.Stdin)
	for s.Scan() {
		log.Println("client says:", s.Text())
	}
}