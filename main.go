package main

import (
	"fmt"
	"log"
	"time"
)

func main() {
	done := make(chan bool)
	defer close(done)
	pipeline := upload(done, decrypt(done, download(done, list(done))))
	for v := range pipeline {
		fmt.Println("main: Completed", v)
	}
}

func list(done <-chan bool) <-chan string {
	names := make(chan string)
	go func() {
		defer close(names)
		for _, f := range []string{"file1", "file2", "file3", "file4", "file5"} {
			log.Printf("list: Found %s", f)
			select {
			case <-done:
				return
			case names <- f:
			}
		}
	}()
	return names
}

func download(done <-chan bool, names <-chan string) <-chan string {
	ciphertext := make(chan string)
	go func() {
		defer close(ciphertext)
		for s := range names {
			log.Printf("download: Received %s", s)
			time.Sleep(3 * time.Second)
			select {
			case <-done:
				return
			case ciphertext <- "lengthy ciphertext from " + s:
			}
			log.Printf("download: Completed %s", s)
		}
	}()
	return ciphertext
}

func decrypt(done <-chan bool, ciphertext <-chan string) <-chan string {
	plaintext := make(chan string)
	go func() {
		defer close(plaintext)
		for s := range ciphertext {
			log.Printf("decrypt: Received %s", s)
			time.Sleep(1 * time.Second)
			select {
			case <-done:
				return
			case plaintext <- "lenghty cleartext of " + s:
			}
			log.Printf("decrypt: Completed %s", s)
		}
	}()
	return plaintext
}

func upload(done <-chan bool, cleartext <-chan string) <-chan string {
	result := make(chan string)
	go func() {
		defer close(result)
		for s := range cleartext {
			log.Printf("upload: Received %s", s)
			time.Sleep(3 * time.Second)
			log.Printf("upload: Completed %s", s)
			select {
			case <-done:
				return
			case result <- s:
			}
		}
	}()
	return result
}
