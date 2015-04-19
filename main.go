package main

import (
	"bufio"
	"os"
	"os/signal"
	"syscall"

	"encoding/json"
	"fmt"

	"log"
	"net"
	"os/exec"
)

type Command struct {
	Name string   `json:"name"`
	Args []string `json:"args"`
}

func main() {
	commands := make(chan Command)
	go ListenTCP(commands, 9999)

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	for {
		select {
		case command := <-commands:
			if command.Name == "volume" {
				exec.Command("amixer", "sset", "'Master'", fmt.Sprintf("%s%%", command.Args[0])).Run()
			}
			fmt.Println(command)

		case <-signals:
			return
		}
	}
}

func ListenTCP(output_stream chan Command, port int) error {
	listener, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", port))
	defer listener.Close()

	if err != nil {
		return err
	}

	for {
		conn, _ := listener.Accept()

		go func() {
			// Receive message
			scanner := bufio.NewScanner(conn)
			for scanner.Scan() {
				// Decode message
				m := Command{}
				if err := json.Unmarshal([]byte(scanner.Text()), &m); err != nil {
					log.Println(err)
					continue
				}
				// Output decode message
				output_stream <- m
			}

			if err := scanner.Err(); err != nil {
				log.Println(err)
			}
			conn.Close()
		}()
	}

	return nil
}
