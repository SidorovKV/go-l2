package main

import (
	"fmt"
	"strings"
	"unicode"
)

/*
=== Задача на распаковку ===

Создать Go функцию, осуществляющую примитивную распаковку строки, содержащую повторяющиеся символы / руны, например:
	- "a4bc2d5e" => "aaaabccddddde"
	- "abcd" => "abcd"
	- "45" => "" (некорректная строка)
	- "" => ""
Дополнительное задание: поддержка escape - последовательностей
	- qwe\4\5 => qwe45 (*)
	- qwe\45 => qwe44444 (*)
	- qwe\\5 => qwe\\\\\ (*)

В случае если была передана некорректная строка функция должна возвращать ошибку. Написать unit-тесты.

Функция должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/

func main() {
	fmt.Println(unpackEscape(`a4bc2d5e`))
}

func unpack(s string) (string, error) {
	if s == "" {
		return "", nil
	}

	out := &strings.Builder{}
	runes := []rune(s)
	l := len(runes)
	out.Grow(len(s))

	if unicode.IsDigit(runes[0]) {
		return "", fmt.Errorf("invalid string")
	}

	var lastRune rune

	for i := 0; i < l; i++ {
		c := runes[i]
		if unicode.IsPrint(c) && !unicode.IsDigit(c) {
			out.WriteRune(c)
			lastRune = c
		}

		if unicode.IsDigit(c) {
			n := int(c - '0')
			for i < l-1 && unicode.IsDigit(runes[i+1]) {
				i++
				n = n*10 + int(runes[i]-'0')
			}
			if out.Cap() < out.Len()+n*2 {
				out.Grow(n * 2)
			}
			for j := 1; j < n; j++ {
				out.WriteRune(lastRune)
			}
		}
	}
	return out.String(), nil
}

func unpackEscape(s string) (string, error) {
	if s == "" {
		return "", nil
	}

	out := &strings.Builder{}
	runes := []rune(s)
	l := len(runes)
	out.Grow(len(s))

	if unicode.IsDigit(runes[0]) {
		return "", fmt.Errorf("invalid string")
	}

	var lastRune rune

	for i := 0; i < l; i++ {
		if runes[i] == '\\' {
			i++
			if i >= l {
				break
			}
			c := runes[i]
			if unicode.IsPrint(c) {
				out.WriteRune(c)
				lastRune = c
			}
			continue
		}

		c := runes[i]
		if unicode.IsPrint(c) && !unicode.IsDigit(c) {
			out.WriteRune(c)
			lastRune = c
		}
		if unicode.IsDigit(c) {
			n := int(c - '0')
			for j := i; j < l-1 && unicode.IsDigit(runes[j+1]); {
				j++
				n = n*10 + int(runes[j]-'0')
			}
			if out.Cap() < out.Len()+n*2 {
				out.Grow(n * 2)
			}
			for j := 1; j < n; j++ {
				out.WriteRune(lastRune)
			}
		}
	}

	return out.String(), nil
}
