package main

import (
	"log"
	"time"
)

func main() {
	names := make(chan string)
	ciphertext := make(chan string)
	cleartext := make(chan string)
	done := make(chan bool)

	go download(names, ciphertext)
	go decrypt(ciphertext, cleartext)
	go upload(cleartext, done)

	for _, f := range []string{"file1", "file2", "file3", "file4", "file5"} {
		log.Printf("main: Processing %s", f)
		names <- f
	}

	close(names)
	<-done
}

func download(names <-chan string, ciphertext chan<- string) {
	defer close(ciphertext)
	for s := range names {
		log.Printf("download: Received %s", s)
		time.Sleep(3 * time.Second)
		ciphertext <- "lengthy ciphertext from " + s
		log.Printf("download: Completed %s", s)
	}
}

func decrypt(ciphertext <-chan string, cleartext chan<- string) {
	defer close(cleartext)
	for s := range ciphertext {
		log.Printf("decrypt: Received %s", s)
		time.Sleep(1 * time.Second)
		cleartext <- "lenghty cleartext of " + s
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
