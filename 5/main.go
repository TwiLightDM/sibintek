package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"

	"github.com/go-audio/audio"
	"github.com/go-audio/wav"
)

func encodeMessage(samples []int, message string) []int {
	messageBytes := append([]byte(message), 0) // null-terminated
	bitIndex := 0

	for i := range samples {
		if bitIndex >= len(messageBytes)*8 {
			break
		}

		byteIndex := bitIndex / 8
		bitInByte := 7 - (bitIndex % 8)
		bit := (messageBytes[byteIndex] >> bitInByte) & 1

		sample := samples[i] &^ 1 // clear LSB
		samples[i] = sample | int(bit)

		bitIndex++
	}

	return samples
}

func decodeMessage(samples []int) string {
	var buffer bytes.Buffer
	var currentByte byte
	bitCount := 0

	for _, sample := range samples {
		lsb := sample & 1
		currentByte = (currentByte << 1) | byte(lsb)
		bitCount++

		if bitCount == 8 {
			if currentByte == 0 {
				break
			}
			buffer.WriteByte(currentByte)
			currentByte = 0
			bitCount = 0
		}
	}

	return buffer.String()
}

func readWavFile(path string) (*audio.IntBuffer, *wav.Decoder, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, nil, err
	}
	defer f.Close()

	decoder := wav.NewDecoder(f)
	if !decoder.IsValidFile() {
		return nil, nil, fmt.Errorf("invalid WAV file")
	}

	buf, err := decoder.FullPCMBuffer()
	if err != nil {
		return nil, nil, err
	}

	return buf, decoder, nil
}

func writeWavFile(path string, buf *audio.IntBuffer) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	enc := wav.NewEncoder(f, buf.Format.SampleRate, buf.SourceBitDepth, buf.Format.NumChannels, 1)
	defer enc.Close()

	return enc.Write(buf)
}

func main() {
	audioPath := flag.String("audio", "", "Path to WAV file")
	decodeFlag := flag.Bool("decode", false, "Set this flag to decode hidden message")
	flag.Parse()

	if *audioPath == "" {
		fmt.Println("Укажите WAV-файл с помощью флага -audio")
		return
	}

	buf, _, err := readWavFile(*audioPath)
	if err != nil {
		panic(err)
	}

	if *decodeFlag {
		msg := decodeMessage(buf.Data)
		fmt.Println("Извлечённое сообщение:", msg)

		err = writeWavFile("decode.wav", buf)
		if err != nil {
			panic(err)
		}
		fmt.Println("Сохранён файл decode.wav")
	} else {
		var message string
		fmt.Print("Введите сообщение для встраивания: ")
		fmt.Scanln(&message)

		buf.Data = encodeMessage(buf.Data, message)
		err := writeWavFile("encode.wav", buf)
		if err != nil {
			panic(err)
		}
		fmt.Println("Сообщение встроено и сохранено в encode.wav")
	}
}
