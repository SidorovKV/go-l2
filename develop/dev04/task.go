package main

import (
	"fmt"
	"slices"
)

/*
=== Поиск анаграмм по словарю ===

Напишите функцию поиска всех множеств анаграмм по словарю.
Например:
'пятак', 'пятка' и 'тяпка' - принадлежат одному множеству,
'листок', 'слиток' и 'столик' - другому.

Входные данные для функции: ссылка на массив - каждый элемент которого - слово на русском языке в кодировке utf8.
Выходные данные: Ссылка на мапу множеств анаграмм.
Ключ - первое встретившееся в словаре слово из множества
Значение - ссылка на массив, каждый элемент которого, слово из множества. Массив должен быть отсортирован по возрастанию.
Множества из одного элемента не должны попасть в результат.
Все слова должны быть приведены к нижнему регистру.
В результате каждое слово должно встречаться только один раз.

Программа должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/

func main() {
	d := []string{"столик", "листок", "слиток", "пятак", "пятка", "тяпка"}
	fmt.Println(FindAnagrams(d))
}

func FindAnagrams(dictionary []string) map[string][]string {
	temp := make(map[string][]string, len(dictionary))
	for _, word := range dictionary {
		key := makeKey(word)
		annagrams, ok := temp[key]
		if ok {
			annagrams = append(annagrams, word)
			temp[key] = annagrams
		} else {
			temp[key] = []string{word}
		}
	}
	out := make(map[string][]string, len(temp))

	for _, value := range temp {
		if len(value) > 1 {
			out[value[0]] = value
			slices.Sort(value)
		}
	}
	return out
}

func makeKey(word string) string {
	chars := []rune(word)
	slices.Sort(chars)
	return string(chars)
}
