package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

/*
=== Взаимодействие с ОС ===

Необходимо реализовать свой собственный UNIX-шелл-утилиту с поддержкой ряда простейших команд:


- cd <args> - смена директории (в качестве аргумента могут быть то-то и то)
- pwd - показать путь до текущего каталога
- echo <args> - вывод аргумента в STDOUT
- kill <args> - "убить" процесс, переданный в качесте аргумента (пример: такой-то пример)
- ps - выводит общую информацию по запущенным процессам в формате *такой-то формат*

Так же требуется поддерживать функционал fork/exec-команд

Дополнительно необходимо поддерживать конвейер на пайпах (linux pipes, пример cmd1 | cmd2 | .... | cmdN).


*Шелл — это обычная консольная программа, которая будучи запущенной, в интерактивном сеансе выводит некое приглашение
в STDOUT и ожидает ввода пользователя через STDIN. Дождавшись ввода, обрабатывает команду согласно своей логике
и при необходимости выводит результат на экран. Интерактивный сеанс поддерживается до тех пор, пока не будет введена команда выхода (например \quit).

*/

func main() {
	home, errr := os.UserHomeDir()
	if errr != nil {
		log.Fatal(errr)
	}
	if os.Chdir(home) != nil {
		log.Fatal(errr)
	}

	currentDir := home

	var input string
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print(currentDir, "> $ ")

		b, err := reader.ReadByte()
		if b == 0 && err != nil {
			log.Println(err)
			os.Exit(0)
		}
		reader.UnreadByte()
		input, err = reader.ReadString('\n')
		if err != nil {
			log.Println(err)
			continue
		}
		if input != "" {
			input = strings.TrimSuffix(input, "\n")
			input = strings.TrimSpace(input)
			args := strings.Split(input, " ")
			if len(args) >= 1 {
				input = args[0]
				args = args[1:]
			}
			switch input {
			case "ps":
				ps(args)
			case "kill":
				if err = kill(args); err != nil {
					fmt.Println(err)
				}
			case "echo":
				fmt.Println("Not implemented yet")
			case "pwd":
				pwd()
			case "cd":
				if len(args) > 1 {
					fmt.Println("Too many arguments")
				} else if len(args) == 1 {
					if err := cd(args[0], currentDir); err != nil {
						fmt.Println(err)
					}
				}
				currentDir, err = os.Getwd()
				if err != nil {
					log.Fatal(err)
				}
			case `\quit`:
				os.Exit(0)
			}
			input = ""
			err = nil
		}
	}
}

func pwd() {
	result, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(result)

}

func cd(path string, currentDir string) error {
	dir := currentDir
	if path == ".." {
		if currentDir == "/" {
			return errors.New("cannot go up")
		}
		currentDir = strings.TrimSuffix(currentDir, "/")
		splitted := strings.Split(currentDir, "/")
		dir = strings.Join(splitted[:len(splitted)-1], "/")
		fmt.Println(dir)
	} else if path == "." {
		return nil
	} else {
		firstTry := currentDir + "/" + path
		secondTry := path
		if _, err := os.Stat(firstTry); err == nil {
			dir = firstTry
		} else if _, err := os.Stat(secondTry); err == nil {
			dir = secondTry
		} else {
			return errors.New("no such file or directory")
		}
	}
	if err := os.Chdir(dir); err != nil {
		return err
	}
	return nil
}

func ps(args []string) {
	cmd := exec.Command("ps", args...)
	psOut, err := cmd.Output()
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Println(string(psOut))
}

func kill(args []string) error {
	for _, arg := range args {
		pid, err := strconv.Atoi(arg)
		if err != nil {
			return err
		}
		process, err := os.FindProcess(pid)
		if err != nil {
			return err
		}
		if err = process.Kill(); err != nil {
			return err
		}
	}
	return nil
}
