package main

import (
	"log"
	"time"
)

type job struct {
	name    string
	process bool
	data    string
}

func main() {
	done := make(chan any)
	defer close(done)
	start := time.Now()
	pipeline := remove(done, upload(done, decrypt(done, download(done, filter(done, list(done))))))
	for v := range pipeline {
		log.Printf("main: %s processed with result %s", v.name, v.data)
	}
	log.Printf("main: took %v", time.Since(start))
}

func list(done <-chan any) <-chan job {
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
		log.Print("list: stage complete")
	}()
	return output
}

func filter(done <-chan any, input <-chan job) <-chan job {
	output := make(chan job)
	go func() {
		defer close(output)
		for j := range input {
			j.process = j.name != "file3"
			select {
			case <-done:
				return
			case output <- j:
			}
		}
		log.Print("filter: stage complete")
	}()
	return output
}
func download(done <-chan any, input <-chan job) <-chan job {
	output := make(chan job)
	go func() {
		defer close(output)
		for j := range input {
			if j.process {
				log.Printf("ftp: GET %s", j.name)
				time.Sleep(3 * time.Second)
				j.data = "lengthy ciphertext"
			}
			select {
			case <-done:
				return
			case output <- j:
			}
		}
		log.Print("download: stage complete")
	}()
	return output
}

func decrypt(done <-chan any, input <-chan job) <-chan job {
	output := make(chan job)
	go func() {
		defer close(output)
		for j := range input {
			if j.process {
				log.Printf("decrypt: %s", j.name)
				time.Sleep(1 * time.Second)
				j.data = "lengthy plaintext"
			}
			select {
			case <-done:
				return
			case output <- j:
			}
		}
		log.Print("decrypt: stage complete")
	}()
	return output
}

func upload(done <-chan any, input <-chan job) <-chan job {
	output := make(chan job)
	go func() {
		defer close(output)
		for j := range input {
			if j.process {
				log.Printf("s3: PUT %s", j.name)
				time.Sleep(3 * time.Second)
			}
			select {
			case <-done:
				return
			case output <- j:
			}
		}
		log.Print("upload: stage complete")
	}()
	return output
}

func remove(done <-chan any, input <-chan job) <-chan job {
	output := make(chan job)
	go func() {
		defer close(output)
		for j := range input {
			log.Printf("ftp: DEL %s", j.name)
			j.data = "OK"
			select {
			case <-done:
				return
			case output <- j:
			}
		}
		log.Print("remove: stage complete")
	}()
	return output
}
