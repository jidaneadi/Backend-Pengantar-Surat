package utils

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
)

type ChaCha20 struct {
	Key     [32]byte
	Nonce   [12]byte
	Counter uint32
}

func quarterRound(a, b, c, d *uint32) {
	*a += *b
	*d ^= *a
	*d = (*d << 16) | (*d >> (32 - 16))
	*c += *d
	*b ^= *c
	*b = (*b << 12) | (*b >> (32 - 12))
	*a += *b
	*d ^= *a
	*d = (*d << 8) | (*d >> (32 - 8))
	*c += *d
	*b ^= *c
	*b = (*b << 7) | (*b >> (32 - 7))
}

func (c *ChaCha20) chacha20Block() [64]byte {
	var state [16]uint32
	var keystream [64]byte

	state[0] = 0x61707865
	state[1] = 0x3320646e
	state[2] = 0x79622d32
	state[3] = 0x6b206574
	for i := 0; i < 8; i++ {
		state[4+i] = uint32(c.Key[i*4]) | uint32(c.Key[i*4+1])<<8 | uint32(c.Key[i*4+2])<<16 | uint32(c.Key[i*4+3])<<24
	}
	state[12] = c.Counter
	for i := 0; i < 3; i++ {
		state[13+i] = uint32(c.Nonce[i*4]) | uint32(c.Nonce[i*4+1])<<8 | uint32(c.Nonce[i*4+2])<<16 | uint32(c.Nonce[i*4+3])<<24
	}

	workingState := state

	for i := 0; i < 10; i++ {
		quarterRound(&workingState[0], &workingState[4], &workingState[8], &workingState[12])
		quarterRound(&workingState[1], &workingState[5], &workingState[9], &workingState[13])
		quarterRound(&workingState[2], &workingState[6], &workingState[10], &workingState[14])
		quarterRound(&workingState[3], &workingState[7], &workingState[11], &workingState[15])
		quarterRound(&workingState[0], &workingState[5], &workingState[10], &workingState[15])
		quarterRound(&workingState[1], &workingState[6], &workingState[11], &workingState[12])
		quarterRound(&workingState[2], &workingState[7], &workingState[8], &workingState[13])
		quarterRound(&workingState[3], &workingState[4], &workingState[9], &workingState[14])
	}

	for i := 0; i < 16; i++ {
		workingState[i] += state[i]
		keystream[i*4] = byte(workingState[i])
		keystream[i*4+1] = byte(workingState[i] >> 8)
		keystream[i*4+2] = byte(workingState[i] >> 16)
		keystream[i*4+3] = byte(workingState[i] >> 24)
	}

	return keystream
}

func (c *ChaCha20) XORKeyStream(dst, src []byte) {
	var block [64]byte
	for len(src) > 0 {
		block = c.chacha20Block()
		c.Counter++
		n := len(src)
		if n > 64 {
			n = 64
		}
		for i := 0; i < n; i++ {
			dst[i] = src[i] ^ block[i]
		}
		src = src[n:]
		dst = dst[n:]
	}
}

func EncryptDecryptFile2(chacha *ChaCha20, inputFilePath string, outputFilePath string) error {
	input, err := os.Open(inputFilePath)
	if err != nil {
		return err
	}
	defer input.Close()

	data, err := io.ReadAll(input)
	if err != nil {
		return err
	}

	output := make([]byte, len(data))
	chacha.XORKeyStream(output, data)

	err = os.WriteFile(outputFilePath, output, 0644)
	if err != nil {
		return err
	}

	fmt.Println("msg", "Sukses melakukan enkripsi/dekripsi")
	return nil
}

// <<<<<<<<<<=========================Enkrip filenya============================>>>>>>>>>>>>>>>
func EncryptDecryptFile(chacha *ChaCha20, inputFile multipart.File, outputFilePath string) error {
	data, err := io.ReadAll(inputFile)
	if err != nil {
		return err
	}

	output := make([]byte, len(data))
	chacha.XORKeyStream(output, data)

	err = os.WriteFile(outputFilePath, output, 0644)
	if err != nil {
		return err
	}

	fmt.Println("msg", "Sukses melakukan enkripsi")
	return nil
}

func DecryptFile(chacha *ChaCha20, inputFilePath, outputFilePath string) error {
	// Baca data terenkripsi dari file
	encryptedData, err := os.ReadFile(inputFilePath)
	if err != nil {
		return err
	}

	// Dekripsi data
	decryptedData := make([]byte, len(encryptedData))
	chacha.Counter = 1 // Reset counter untuk dekripsi
	chacha.XORKeyStream(decryptedData, encryptedData)

	// Simpan hasil dekripsi ke file output
	err = os.WriteFile(outputFilePath, decryptedData, 0644)
	if err != nil {
		return err
	}

	return nil
}

// func EncryptDecryptFile(c *fiber.Ctx, chacha *ChaCha20, inputFile *multipart.FileHeader, outputFile string) error {
// 	src, err := inputFile.Open()
// 	if err != nil {
// 		return err
// 	}
// 	defer src.Close()

// 	dst, err := os.Create(outputFile)
// 	if err != nil {
// 		return err
// 	}
// 	defer dst.Close()

// 	data, err := io.ReadAll(src)
// 	if err != nil {
// 		return err
// 	}

// 	output := make([]byte, len(data))
// 	chacha.XORKeyStream(output, data)

// 	_, err = dst.Write(output)
// 	if err != nil {
// 		return err
// 	}

// 	fmt.Println("msg", "Sukses melakukan enkripsi")
// 	return nil
// }
