package main

import (
	"flag"
	"fmt"
	"os"
)

const BlockSize = 8 // 64 бита (8 байт)

// Простая (фиктивная) таблица перестановок для упрощения примера
var IP = [BlockSize]byte{2, 6, 3, 1, 4, 8, 5, 7} // Начальная перестановка
var FP = [BlockSize]byte{4, 1, 3, 5, 7, 2, 8, 6} // Конечная перестановка

// Генерация раундовых ключей (здесь просто сдвигаем ключ, но в DES сложнее)
func generateSubKeys(key []byte) [][]byte {
	subKeys := make([][]byte, 16)
	for i := 0; i < 16; i++ {
		subKeys[i] = make([]byte, BlockSize)
		for j := 0; j < BlockSize; j++ {
			subKeys[i][j] = key[(j+i)%len(key)] // Циклический сдвиг
		}
	}
	return subKeys
}

// Функция перестановки (универсальная)
func permute(input []byte, table [BlockSize]byte) []byte {
	output := make([]byte, len(input))
	for i, pos := range table {
		output[i] = input[pos-1]
	}
	return output
}

// Основная DES-функция (очень упрощённая версия)
func desEncryptBlock(block, key []byte) []byte {
	block = permute(block, IP)      // Начальная перестановка
	subKeys := generateSubKeys(key) // Генерация ключей

	// 16 раундов (упрощённая логика)
	for i := 0; i < 16; i++ {
		for j := 0; j < BlockSize; j++ {
			block[j] ^= subKeys[i][j] // XOR с подраундовым ключом (вместо F-функции)
		}
	}

	block = permute(block, FP) // Итоговая перестановка
	return block
}

func processFile(filePath string, key []byte, decrypt bool) {
	file, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Println("Ошибка чтения файла:", err)
		return
	}

	if len(file)%BlockSize != 0 {
		fmt.Println("Файл должен быть кратен 8 байтам!")
		return
	}

	result := make([]byte, len(file))
	for i := 0; i < len(file); i += BlockSize {
		block := file[i : i+BlockSize]
		if decrypt {
			result[i] = desEncryptBlock(block, key)[0] // Упрощённый процесс (расшифрование аналогично)
		} else {
			result[i] = desEncryptBlock(block, key)[0]
		}
	}

	os.WriteFile(filePath, result, 0666)
	fmt.Println("Файл обработан:", filePath)
}

func main() {
	file := flag.String("file", "", "Файл для обработки")
	decrypt := flag.Bool("decrypt", false, "Расшифровать файл")
	flag.Parse()

	if *file == "" {
		fmt.Println("Укажите файл через -file")
		return
	}

	key := []byte("12345678") // 8 байтовый ключ
	processFile(*file, key, *decrypt)
}
