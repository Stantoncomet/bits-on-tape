package main

import (
	"fmt"
)

func main() {
	data := "H"

	fmt.Printf("\n%b\n", strToByte(data))
	// fmt.Printf("\n%v\n", strToHex(data))
	//fmt.Printf("\n%v\n", strToBin(data))

	// binToAudio(strToBin(data))

	listen()

	// save_rec()
}
