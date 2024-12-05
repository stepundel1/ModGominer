package main

import (
	"crypto/rand"
	"encoding/binary"
	"encoding/hex"
	"log"
)

func generateRandomNonce() string {

	// newSoloMiner(background)
	nonce := make([]byte, 4)
	_, err := rand.Read(nonce)
	if err != nil {
		log.Printf("Ошибка генерации случайного решения: %v", err)
		return ""
	}
	log.Printf("Сгенерировано новое решение 23234: %x", nonce)
	return hex.EncodeToString(nonce)
}

func generateRandomNonce2() uint32 {
	var nonce uint32
	err := binary.Read(rand.Reader, binary.LittleEndian, &nonce)
	if err != nil {
		log.Printf("Ошибка генерации nonce: %v", err)
	}
	return nonce
}
