package main

import (
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
	pipeline := remove(done, upload(done, decrypt(done, download(done, list(done)))))
	for v := range pipeline {
		log.Printf("main: %s processed with result %s", v.name, v.data)
	}
}

func list(done <-chan bool) <-chan job {
	output := make(chan job)
	go func() {
		defer close(output)
		log.Println("ftp: DIR")
		for _, f := range []string{"file1", "file2", "file3", "file4", "file5"} {
			select {
			case <-done:
				return
			case output <- job{name: f}:
			}
		}
	}()
	return output
}

func download(done <-chan bool, input <-chan job) <-chan job {
	output := make(chan job)
	go func() {
		defer close(output)
		for j := range input {
			log.Printf("ftp: GET %s", j.name)
			time.Sleep(3 * time.Second)
			select {
			case <-done:
				return
			case output <- job{j.name, "lengthy ciphertext"}:
			}
		}
		log.Printf("download complete")
	}()
	return output
}

func decrypt(done <-chan bool, ciphertext <-chan job) <-chan job {
	plaintext := make(chan job)
	go func() {
		defer close(plaintext)
		for j := range ciphertext {
			log.Printf("decrypt: %s", j.name)
			time.Sleep(1 * time.Second)
			select {
			case <-done:
				return
			case plaintext <- job{j.name, "lenghty cleartext"}:
			}
		}
	}()
	return plaintext
}

func upload(done <-chan bool, input <-chan job) <-chan job {
	output := make(chan job)
	go func() {
		defer close(output)
		for j := range input {
			log.Printf("s3: PUT %s", j.name)
			time.Sleep(3 * time.Second)
			select {
			case <-done:
				return
			case output <- job{j.name, ""}:
			}
		}
	}()
	return output
}

func remove(done <-chan bool, input <-chan job) <-chan job {
	output := make(chan job)
	go func() {
		defer close(output)
		for j := range input {
			log.Printf("ftp: RM %s", j.name)
			select {
			case <-done:
				return
			case output <- job{j.name, "OK"}:
			}
		}
	}()
	return output
}
