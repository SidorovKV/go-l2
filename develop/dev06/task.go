package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

/*
=== Утилита cut ===

Принимает STDIN, разбивает по разделителю (TAB) на колонки, выводит запрошенные

Поддержать флаги:
-f - "fields" - выбрать поля (колонки)
-d - "delimiter" - использовать другой разделитель
-s - "separated" - только строки с разделителем

Программа должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/

func main() {
	args := os.Args[1:]
	files := make([]*os.File, 0)
	input := make([]string, 0)
	for _, v := range args {
		if v[0] != '-' || strings.Contains(v, ".") || strings.Contains(v, "/") {
			file, err := os.Open(v)
			if err != nil {
				log.Fatal(err)
			}
			files = append(files, file)
		}
	}
	if len(files) == 0 {
		reader := bufio.NewReader(os.Stdin)
		for txt, err := reader.ReadString('\n'); err == nil; txt, err = reader.ReadString('\n') {
			txt = strings.TrimSuffix(txt, "\n")
			input = append(input, txt)
		}
	}

	for _, file := range files {
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			input = append(input, scanner.Text())
		}
		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}
		file.Close()
	}

	output, err := cutUtility(input, args)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(output)
}

func cutUtility(input []string, args []string) (string, error) {
	result := &strings.Builder{}
	delimiter := "\t"
	flags := &strings.Builder{}
	fields := make([]uint, 0)

	if len(args) != 0 {
		flags.Grow(len(args))
	}
	for _, v := range args {
		if v[0] == '-' && !strings.Contains(v, ".") && !strings.Contains(v, "/") {
			if len(v) > 1 {
				if v[1] == 'f' && v[2] == '=' {
					f := v[3:]
					flds := strings.Split(f, ",")
					for _, column := range flds {
						if strings.Contains(column, "-") {
							s := strings.Split(column, "-")
							if len(s) == 2 {
								c1, err1 := strconv.Atoi(s[0])
								c2, err2 := strconv.Atoi(s[1])
								if err1 == nil && err2 == nil {
									for c := c1; c <= c2; c++ {
										fields = append(fields, uint(c))
									}
								}
							} else {
								return "", errors.New("invalid field specification")
							}
						} else {
							c, err := strconv.Atoi(column)
							if err == nil && c >= 0 {
								fields = append(fields, uint(c))
							}
						}
					}
				} else if v[1] == 'd' && v[2] == '=' {
					delimiter = v[3:]
				} else {
					flags.WriteString(v[1:])
				}
			}
		}
	}

	if len(fields) == 0 {
		return "", errors.New("you must specify a list of fields. For example -f=1,5 or -f=1-5")
	}

	isSeparated := strings.Contains(flags.String(), "s")
	for _, v := range input {
		if isSeparated {
			if !strings.Contains(v, delimiter) {
				continue
			}
		}
		columns := strings.Split(v, delimiter)
		for _, field := range fields {
			result.WriteString(columns[field] + delimiter)
		}
		result.WriteString("\n")
	}
	return result.String(), nil
}
