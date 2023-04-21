package main

import (
	"fmt"
	"log"
	"time"
)

type job struct {
	name string
	data string
}

func main() {
	done := make(chan bool)
	defer close(done)
	pipeline := upload(done, decrypt(done, download(done, list(done))))
	for v := range pipeline {
		fmt.Println("main: Completed", v)
	}
}

func list(done <-chan bool) <-chan job {
	names := make(chan job)
	go func() {
		defer close(names)
		for _, f := range []string{"file1", "file2", "file3", "file4", "file5"} {
			log.Printf("list: Found %s", f)
			select {
			case <-done:
				return
			case names <- job{name: f}:
			}
		}
	}()
	return names
}

func download(done <-chan bool, names <-chan job) <-chan job {
	ciphertext := make(chan job)
	go func() {
		defer close(ciphertext)
		for j := range names {
			log.Printf("download: Received job %s", j.name)
			time.Sleep(3 * time.Second)
			select {
			case <-done:
				return
			case ciphertext <- job{j.name, "lengthy ciphertext"}:
			}
			log.Printf("download: Completed %s", j.name)
		}
	}()
	return ciphertext
}

func decrypt(done <-chan bool, ciphertext <-chan job) <-chan job {
	plaintext := make(chan job)
	go func() {
		defer close(plaintext)
		for j := range ciphertext {
			log.Printf("decrypt: Received %s", j.name)
			time.Sleep(1 * time.Second)
			select {
			case <-done:
				return
			case plaintext <- job{j.name, "lenghty cleartext"}:
			}
			log.Printf("decrypt: Completed %s", j.name)
		}
	}()
	return plaintext
}

func upload(done <-chan bool, cleartext <-chan job) <-chan job {
	result := make(chan job)
	go func() {
		defer close(result)
		for j := range cleartext {
			log.Printf("upload: Received %s", j.name)
			time.Sleep(3 * time.Second)
			log.Printf("upload: Completed %s", j.name)
			select {
			case <-done:
				return
			case result <- job{j.name, "OK"}:
			}
		}
	}()
	return result
}
