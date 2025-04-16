package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"

func getModifiedAlphabet(key string, shift int) string {
	key = strings.ToLower(key)
	uniqueKey := removeDuplicates(key)

	alphabetWithKey := uniqueKey + alphabet

	modifiedAlphabet := removeDuplicates(alphabetWithKey)

	alphabetShifted := modifiedAlphabet[shift:] + modifiedAlphabet[:shift]

	return alphabetShifted
}

func removeDuplicates(input string) string {
	seen := make(map[rune]bool)
	var result []rune
	for _, char := range input {
		if !seen[char] {
			seen[char] = true
			result = append(result, char)
		}
	}
	return string(result)
}

func encrypt(input, modifiedAlphabet string) string {
	return processText(input, modifiedAlphabet, false)
}

func decrypt(input, modifiedAlphabet string) string {
	return processText(input, modifiedAlphabet, true)
}

func processText(input, modifiedAlphabet string, isDecrypt bool) string {
	var result strings.Builder
	normalAlphabet := alphabet

	for _, char := range input {
		if index := strings.IndexRune(normalAlphabet, char); index != -1 {
			var newChar rune
			if isDecrypt {
				newChar = rune(normalAlphabet[index])
			} else {
				newChar = rune(modifiedAlphabet[index])
			}
			result.WriteRune(newChar)
		} else {
			result.WriteRune(char)
		}
	}

	return result.String()
}

func main() {
	filePath := flag.String("file", "", "Путь к текстовому файлу")
	decryptFlag := flag.Bool("decrypt", false, "Если true, то дешифровать, если false - шифровать")
	shift := flag.Int("shift", 3, "Сдвиг по умолчанию 3")
	key := flag.String("key", "rythm", "Ключевое слово для шифрования/дешифрования")

	flag.Parse()

	var inputText string

	if *filePath != "" {
		data, err := os.ReadFile(*filePath)
		if err != nil {
			fmt.Println("Ошибка при чтении файла:", err)
			os.Exit(1)
		}
		inputText = string(data)
	} else {
		fmt.Println("Введите текст для шифрования/дешифрования:")
		fmt.Scanln(&inputText)
	}
	modifiedAlphabet := getModifiedAlphabet(*key, *shift)

	if *decryptFlag {
		result := decrypt(inputText, modifiedAlphabet)
		fmt.Println("Результат дешифрования:", result)
	} else {
		result := encrypt(inputText, modifiedAlphabet)
		fmt.Println("Результат шифрования: \n", result)
	}
}
