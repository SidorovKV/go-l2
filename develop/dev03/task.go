package main

import (
	"bufio"
	"cmp"
	"errors"
	"fmt"
	"log"
	"os"
	"regexp"
	"slices"
	"strconv"
	"strings"
)

/*
=== Утилита sort ===

Отсортировать строки (man sort)
Основное

Поддержать ключи

-k — указание колонки для сортировки
-n — сортировать по числовому значению
-r — сортировать в обратном порядке
-u — не выводить повторяющиеся строки

Дополнительное

Поддержать ключи

-M — сортировать по названию месяца
-b — игнорировать хвостовые пробелы
-c — проверять отсортированы ли данные
-h — сортировать по числовому значению с учётом суффиксов

Программа должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/

const EOT byte = '\x04'

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

	result, err := sortUtility(args, input)
	if err != nil {
		log.Fatal(err)
	} else if result == nil {
		os.Exit(0)
	}

	for _, v := range result {
		fmt.Println(v)
	}
}

func sortUtility(args []string, toSort []string) ([]string, error) {
	flags := &strings.Builder{}
	column := 0
	delimiter := " "

	if len(args) != 0 {
		flags.Grow(len(args))
	}
	for _, v := range args {
		if v[0] == '-' && !strings.Contains(v, ".") && !strings.Contains(v, "/") {
			if len(v) > 1 {
				if v[1] == 'k' && v[2] == '=' {
					c, err := strconv.Atoi(v[3:])
					if err != nil {
						return nil, err
					}
					column = c
				} else if v[1] == 'd' && v[2] == '=' {
					delimiter = v[3:]
				} else {
					flags.WriteString(v[1:])
				}
			}

		}
	}

	isReverse := strings.Contains(flags.String(), "r")
	isNumeric := strings.Contains(flags.String(), "n")
	isByColumn := column > 0
	if strings.Contains(flags.String(), "c") {
		line, ok := check(toSort, isNumeric, isByColumn, column, delimiter)
		if !ok {
			return nil, errors.New(fmt.Sprintf("not sorted,  line %d", line))
		} else {
			return nil, nil
		}
	}
	if strings.Contains(flags.String(), "u") {
		toSort = unique(toSort)
	}
	if isNumeric {
		numericSort(toSort, isReverse)
	}
	if isByColumn {
		sortByColumn(toSort, column, delimiter, isNumeric, isReverse)
	}
	if !isNumeric && !isByColumn {
		simpleSort(toSort, isReverse)
	}

	return toSort, nil
}

func check(in []string, isNumeric, isByColumn bool, column int, delimiter string) (uint64, bool) {
	if len(in) < 2 {
		return 0, true
	}

	for i := 1; i < len(in); i++ {
		a, b := in[i-1], in[i]

		if isNumeric && !isByColumn {
			re := regexp.MustCompile(`[-]?\d[\d,]*[\.]?[\d{2}]*`)
			x, y := re.FindAllString(a, -1), re.FindAllString(b, -1)

			if len(x) == len(y) {
				for n, sx := range x {
					xf, _ := strconv.ParseFloat(sx, 64)
					yf, _ := strconv.ParseFloat(y[n], 64)
					if cmp.Compare(xf, yf) > 0 {
						return uint64(i), false
					}
				}
			}
			if len(x) > len(y) {
				return uint64(i), false
			}
		}

		if isByColumn {
			x, y := strings.Split(a, delimiter), strings.Split(b, delimiter)
			if len(x) < column-1 || len(y) < column-1 {
				return uint64(i), false
			}
			if !isNumeric && cmp.Compare(x[column-1], y[column-1]) > 0 {
				return uint64(i), false
			}
			if isNumeric {
				re := regexp.MustCompile(`[-]?\d[\d,]*[\.]?[\d{2}]*`)
				x2, y2 := re.FindAllString(x[column-1], -1), re.FindAllString(y[column-1], -1)

				if len(x2) == len(y2) {
					for n, sx := range x2 {
						xf, _ := strconv.ParseFloat(sx, 64)
						yf, _ := strconv.ParseFloat(y2[n], 64)
						if cmp.Compare(xf, yf) > 0 {
							return uint64(i), false
						}
					}
				}
				if len(x2) > len(y2) {
					return uint64(i), false
				}
			}
		}

		if !isNumeric && !isByColumn {
			if cmp.Compare(a, b) > 0 {
				return uint64(i), false
			}
		}
	}

	return 0, true
}

func simpleSort(in []string, isReverse bool) {
	slices.SortFunc(in, func(a, b string) int {
		if isReverse {
			return -cmp.Compare(a, b)
		}
		return cmp.Compare(a, b)
	})
}

func numericSort(in []string, isReverse bool) {
	slices.SortFunc(in, func(a, b string) int {
		return numericSortFunc(a, b, isReverse)
	})
}

func unique(in []string) []string {
	out := make([]string, 0, len(in))
	m := make(map[string]struct{}, len(in))
	for _, v := range in {
		if _, ok := m[v]; !ok {
			m[v] = struct{}{}
			out = append(out, v)
		}
	}
	return out
}

func sortByColumn(in []string, column int, delimiter string, isNumeric, isReverse bool) {
	slices.SortFunc(in, func(a, b string) int {
		x, y := strings.Split(a, delimiter), strings.Split(b, delimiter)

		if len(x) < column-1 && len(y) < column-1 {
			return 0
		}
		if len(x) < column-1 && len(y) >= column-1 {
			if isReverse {
				return -1
			}
			return 1
		}
		if len(x) >= column-1 && len(y) < column-1 {
			if isReverse {
				return 1
			}
			return -1
		}

		if isNumeric {
			return numericSortFunc(x[column-1], y[column-1], isReverse)
		}

		if isReverse {
			return -cmp.Compare(x[column-1], y[column-1])
		}
		return cmp.Compare(x[column-1], y[column-1])
	})
}

func numericSortFunc(a, b string, isReverse bool) int {
	re := regexp.MustCompile(`[-]?\d[\d,]*[\.]?[\d{2}]*`)
	x, y := re.FindAllString(a, -1), re.FindAllString(b, -1)
	if len(x) == 0 && len(y) == 0 {
		return 0
	}

	if isReverse {
		if len(x) == len(y) {
			for n, sx := range x {
				xf, _ := strconv.ParseFloat(sx, 64)
				yf, _ := strconv.ParseFloat(y[n], 64)
				if cmp.Compare(xf, yf) != 0 {
					return -cmp.Compare(xf, yf)
				}
			}
			return 0
		}
		if len(x) < len(y) {
			return 1
		}
		if len(x) > len(y) {
			return -1
		}
	}
	if len(x) == len(y) {
		for n, sx := range x {
			xf, _ := strconv.ParseFloat(sx, 64)
			yf, _ := strconv.ParseFloat(y[n], 64)
			if cmp.Compare(xf, yf) != 0 {
				return cmp.Compare(xf, yf)
			}
		}
		return 0
	}
	if len(x) < len(y) {
		return -1
	}
	if len(x) > len(y) {
		return 1
	}
	return 0
}
