package main

import (
	"fmt"
	"log"
	"time"
)

const files = 5

func main() {
	names := make(chan string)
	ciphertext := make(chan string)
	cleartext := make(chan string)
	done := make(chan bool)

	go download(names, ciphertext)
	go decrypt(ciphertext, cleartext)
	go upload(cleartext, done)

	for i := 0; i < files; i++ {
		file := fmt.Sprintf("file%03d.txt", i)
		log.Printf("main: Processing %s", file)
		names <- file
	}

	close(names)
	<-done
}

func download(names <-chan string, ciphertext chan<- string) {
	defer close(ciphertext)
	for s := range names {
		log.Printf("download: Received %s", s)
		time.Sleep(3 * time.Second)
		ciphertext <- "lengthy ciphertextryped text from " + s
		log.Printf("download: Completed %s", s)
	}
}

func decrypt(ciphertext <-chan string, cleartext chan<- string) {
	defer close(cleartext)
	for s := range ciphertext {
		log.Printf("decrypt: Received %s", s)
		time.Sleep(1 * time.Second)
		cleartext <- "lenghty clear text of " + s
		log.Printf("decrypt: Completed %s", s)

	}
}

func upload(cleartext <-chan string, done chan<- bool) {
	defer close(done)
	for s := range cleartext {
		log.Printf("upload: Received %s", s)
		time.Sleep(3 * time.Second)
		log.Printf("upload: Completed %s", s)
	}
	done <- true
}
