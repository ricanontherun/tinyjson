package main

import (
	"fmt"
	"os"
)

func main() {
	file, err := os.OpenFile("./test.log", os.O_RDWR|os.O_APPEND|os.O_CREATE, 0600)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer file.Close()

	// TODO: How to handle ErrShortWrite?
	// Keep writing more bytes (starting at len(b) - n) until we complete.
	written, err := file.Write([]byte("Hello"))
	if err != nil {
		fmt.Println(err.Error())
	}

	fmt.Printf("wrote bytes %d", written)

	stuff := []byte{}
	fmt.Println(file.Read(stuff))
	fmt.Println(stuff)

}
