package main

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

/*
=== Утилита telnet ===

Реализовать примитивный telnet клиент:
Примеры вызовов:
go-telnet --timeout=10s host port go-telnet mysite.ru 8080 go-telnet --timeout=3s 1.1.1.1 123

Программа должна подключаться к указанному хосту (ip или доменное имя) и порту по протоколу TCP.
После подключения STDIN программы должен записываться в сокет, а данные полученные и сокета должны выводиться в STDOUT
Опционально в программу можно передать таймаут на подключение к серверу (через аргумент --timeout, по умолчанию 10s).

При нажатии Ctrl+D программа должна закрывать сокет и завершаться. Если сокет закрывается со стороны сервера, программа должна также завершаться.
При подключении к несуществующему сервер, программа должна завершаться через timeout.
*/

const EOT byte = 4

func main() {
	args := os.Args[1:]
	result, err := telnet(args)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(result)
}

func parseArgs(args []string) (host string, port string, timeout time.Duration, err error) {
	timeout = 10 * time.Second

	if len(args) < 2 {
		return "", "", 0, errors.New("not enough arguments")
	}

	if len(args) == 3 {
		if strings.Contains(args[0], "--timeout") {
			s := strings.Split(args[0], "=")
			t, err := time.ParseDuration(s[1])
			if err != nil {
				return "", "", 0, err
			}
			if len(args[1]) < len(args[2]) {
				return "", "", 0, errors.New("invalid arguments")
			}
			timeout = t
			host = args[1]
			port = args[2]
		} else if strings.Contains(args[2], "--timeout") {
			s := strings.Split(args[2], "=")
			t, err := time.ParseDuration(s[1])
			if err != nil {
				return "", "", 0, err
			}
			if len(args[0]) < len(args[1]) {
				return "", "", 0, errors.New("invalid arguments")
			}
			timeout = t
			host = args[0]
			port = args[1]
		} else {
			return "", "", 0, errors.New("invalid arguments")
		}
	} else if len(args) == 2 {
		if len(args[0]) < len(args[1]) {
			return "", "", 0, errors.New("invalid arguments")
		}
		host = args[0]
		port = args[1]
	}
	return
}

func telnet(args []string) (string, error) {
	host, port, timeout, err := parseArgs(args)
	if err != nil {
		return "", err
	}

	uri := ""
	if strings.Contains(host, "/") {
		s := strings.SplitN(host, "/", 2)
		uri = s[0] + ":" + port + "/" + s[1]
	} else {
		uri = host + ":" + port
	}

	ctx, cncl := context.WithTimeout(context.Background(), timeout)
	defer cncl()

	go func() {
		reader := bufio.NewReader(os.Stdin)
		for {
			b, _ := reader.ReadByte()
			if b == 0 {
				cncl()
				break
			}
		}
		<-ctx.Done()
		os.Exit(0)
	}()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("https://%s", uri), nil)
	if err != nil {
		return "", err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(data), nil
}
