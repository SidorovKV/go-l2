package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
)

/*
=== Утилита grep ===

Реализовать утилиту фильтрации (man grep)

Поддержать флаги:
-A - "after" печатать +N строк после совпадения
-B - "before" печатать +N строк до совпадения
-C - "context" (A+B) печатать ±N строк вокруг совпадения
-c - "count" (количество строк)
-i - "ignore-case" (игнорировать регистр)
-v - "invert" (вместо совпадения, исключать)
-F - "fixed", точное совпадение со строкой, не паттерн
-n - "line num", печатать номер строки

Программа должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/

type flags struct {
	A int
	B int
	C int
	c bool
	i bool
	v bool
	F bool
	n bool
}

func main() {
	f := parseFlags()
	nonFlagArgs := flag.Args()
	output, err := grep(f, nonFlagArgs)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(output)
}

func parseFlags() *flags {
	f := &flags{}
	flag.IntVar(&f.A, "A", 0, "Print NUM lines of trailing context after matching lines.  Places  a  line  containing  a  group  separator  (--)  between contiguous groups of matches.")
	flag.IntVar(&f.B, "B", 0, "Print  NUM  lines  of  leading  context  before  matching  lines.  Places a line containing a group separator (--) between contiguous groups of matches.")
	flag.IntVar(&f.C, "C", 0, "Print NUM lines of output context.  Places a line containing a group separator (--) between contiguous groups of  matches.")
	flag.BoolVar(&f.c, "c", false, "Suppress normal output; instead print a count of matching lines for each input file.")
	flag.BoolVar(&f.i, "i", false, "Ignore case distinctions in patterns and input data, so that characters that differ only in case match each other.")
	flag.BoolVar(&f.v, "v", false, "Invert the sense of matching, to select non-matching lines.")
	flag.BoolVar(&f.F, "F", false, "Interpret PATTERN as fixed strings, not regular expressions.")
	flag.BoolVar(&f.n, "n", false, "Prefix each line of output with the 1-based line number within its input file.")
	flag.Parse()
	return f
}

func grep(flags *flags, nonFlagArgs []string) (string, error) {
	var input map[string][]string
	var err error
	if len(nonFlagArgs) < 1 {
		return "", errors.New("no pattern provided")
	}
	if len(nonFlagArgs) == 1 {
		var stdin []string
		if stdin, err = readInput(); err != nil {
			return "", err
		}

		input = make(map[string][]string)
		input["stdin"] = stdin
	}
	if len(nonFlagArgs) > 1 {
		if input, err = readFiles(nonFlagArgs[1:]); err != nil {
			return "", err
		}
	}

	var out string
	linesAfter := flags.A + flags.C
	linesBefore := flags.B + flags.C
	if linesAfter < 0 {
		linesAfter = 0
	}
	if linesBefore < 0 {
		linesBefore = 0
	}

	if flags.c {
		out = countMatches(nonFlagArgs[0], nonFlagArgs[1:], input, flags.i, flags.F, flags.v)
	} else {
		out = findMatches(nonFlagArgs[0], nonFlagArgs[1:], input, uint(linesBefore), uint(linesAfter), flags.i, flags.F, flags.v, flags.n)
	}
	return out, nil
}

func countMatches(
	patterns string,
	fileNames []string,
	input map[string][]string,
	isIgnoreCase bool,
	isFixed bool,
	isInvert bool,
) string {

	out := &strings.Builder{}
	pttrns := parsePatterns(patterns)
	for _, p := range pttrns {
		var regex *regexp.Regexp
		if !isFixed {
			if isIgnoreCase {
				p = strings.ToLower(p)
			}
			regex = regexp.MustCompile(p)
		}
		for _, fileName := range fileNames {
			lines, ok := input[fileName]
			if !ok {
				continue
			}
			count := 0
			for _, line := range lines {
				if isIgnoreCase {
					p = strings.ToLower(p)
					line = strings.ToLower(line)
				}

				if isInvert {
					if isFixed {
						if !strings.Contains(line, p) {
							count++
						}
					} else {
						found := regex.FindAllString(line, -1)
						if len(found) == 0 {
							count++
						}
					}
				} else {
					if isFixed {
						if strings.Contains(line, p) {
							count++
						}
					} else {
						found := regex.FindAllString(line, -1)
						if len(found) > 0 {
							count++
						}
					}
				}
			}
			if count > 0 {
				out.WriteString(fileName + ":" + fmt.Sprint(count) + "\n")
			}
		}
	}
	return out.String()
}

func findMatches(
	patterns string,
	fileNames []string,
	input map[string][]string,
	linesBefore uint,
	linesAfter uint,
	isIgnoreCase bool,
	isFixed bool,
	isInvert bool,
	isLineNum bool,
) string {
	out := &strings.Builder{}
	pttrns := parsePatterns(patterns)
	for _, p := range pttrns {
		var regex *regexp.Regexp
		if !isFixed {
			if isIgnoreCase {
				p = strings.ToLower(p)
			}
			regex = regexp.MustCompile(p)
		}
		for _, fileName := range fileNames {
			out.WriteString(fileName + ":\n")
			lines, ok := input[fileName]
			if !ok {
				continue
			}
			lenght := uint(len(lines))
			for lineNum, line := range lines {
				var afterIdx, beforeIdx uint
				if uint(lineNum)+linesAfter < lenght-1 {
					afterIdx = linesAfter + uint(lineNum)
				} else {
					afterIdx = lenght - 1
				}
				if uint(lineNum)-linesBefore > 0 {
					beforeIdx = uint(lineNum) - linesBefore
				} else {
					beforeIdx = 0
				}
				if isIgnoreCase {
					p = strings.ToLower(p)
					line = strings.ToLower(line)
				}

				if isInvert {
					if isFixed {
						if !strings.Contains(line, p) {
							toPrint := lines[beforeIdx : afterIdx+1]
							for n, l := range toPrint {
								if isLineNum {
									out.WriteString(fmt.Sprintf("%d:", beforeIdx+uint(n)))
								}
								out.WriteString(l + "\n")
							}
							if len(toPrint) > 1 {
								out.WriteString("--")
							}
						}
					} else {
						found := regex.FindAllString(line, -1)
						if len(found) == 0 {
							toPrint := lines[beforeIdx : afterIdx+1]
							for n, l := range toPrint {
								if isLineNum {
									out.WriteString(fmt.Sprintf("%d:", beforeIdx+uint(n)))
								}
								out.WriteString(l + "\n")
							}
							if len(toPrint) > 1 {
								out.WriteString("--")
							}
						}
					}
				} else {
					if isFixed {
						if strings.Contains(line, p) {
							toPrint := lines[beforeIdx : afterIdx+1]
							for n, l := range toPrint {
								if isLineNum {
									out.WriteString(fmt.Sprintf("%d:", beforeIdx+uint(n)))
								}
								out.WriteString(l + "\n")
							}
							if len(toPrint) > 1 {
								out.WriteString("--")
							}
						}
					} else {
						found := regex.FindAllString(line, -1)
						if len(found) > 0 {
							toPrint := lines[beforeIdx : afterIdx+1]
							for n, l := range toPrint {
								if isLineNum {
									out.WriteString(fmt.Sprintf("%d:", beforeIdx+uint(n)))
								}
								out.WriteString(l + "\n")
							}
							if len(toPrint) > 1 {
								out.WriteString("--\n")
							}
						}
					}
				}
			}
			out.WriteString("---\n")
		}
	}

	return out.String()
}

func readInput() ([]string, error) {
	input := make([]string, 0, 5)
	reader := bufio.NewReader(os.Stdin)
	var err error
	var txt string
	for txt, err = reader.ReadString('\n'); err == nil; txt, err = reader.ReadString('\n') {
		txt = strings.TrimSuffix(txt, "\n")
		input = append(input, txt)
	}
	if err != nil {
		return nil, err
	}
	return input, nil
}

func readFiles(files []string) (map[string][]string, error) {
	out := make(map[string][]string)

	for _, file := range files {
		var input []string
		var f *os.File
		var err error
		if _, err = os.Stat(file); err == nil {
			if f, err = os.Open(file); err != nil {
				return nil, err
			}
		} else {
			var curDir string
			if curDir, err = os.Getwd(); err != nil {
				return nil, err
			}
			path := curDir + "/" + file
			if _, err = os.Stat(path); err != nil {
				return nil, err
			}

			if f, err = os.Open(path); err != nil {
				return nil, err
			}
		}

		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			input = append(input, scanner.Text())
		}
		f.Close()
		out[file] = input
	}
	return out, nil
}

func parsePatterns(s string) []string {
	if strings.Contains(s, "|") {
		return strings.Split(s, "|")
	}
	return []string{s}
}
