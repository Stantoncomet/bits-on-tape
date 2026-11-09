package main

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/MarkKremer/microphone"
)

const (
	SAMPLE_RATE = 44100
	// duration bits were encoded in (in milliseconds)
	BIT_LENGTH       = 100
	CAPTURES_PER_BIT = 2
	// captured frequencies every SAMPLE_RATE / CAPTURE_RATE seconds
	// ex: SAMPLE_RATE / 4 is 1/4 second, SAMPLE_RATE * (BIT_LENGTH / 1000 * 2) is two captures per bit
	CAPTURE_RATE = SAMPLE_RATE * BIT_LENGTH / (1000 * CAPTURES_PER_BIT)
	// CAPTURE_RATE = SAMPLE_RATE / 20
)

func binToStr() {
	// var rawbin = "1111011111111111111110111111111111100000011111111111000000000001011111111101000000000000100000000010101001111111111111111111000010000000011111111100000001111111110100110111111111111111110111111111111111111111111110000001100000010000111111111111111111010011011111111111111111111111100101011100000101111111111111111111000101111111111111111111111111111111111101110001011011111111111111101100110111111111111111111100111111011001111111111110011111111111111111111111111111111111111111111111001111111111111111111110111101111111111111111111111111111111111111111111111111111111111111111101011100001011111111111101010111111001011111111111111111100111111111111111111111111101001100110000111111111111111111111110101011101010011111111110011001011100000110100011110001111111111110011100000111110001000100101111111111111111111111111111"

	// for _, c := range rawbin {

	// }

}

func listen() {
	fmt.Println("Starting microphone...")

	microphone.Init()
	defer microphone.Terminate()

	stream, _, _ := microphone.OpenDefaultStream(SAMPLE_RATE, 1)
	defer stream.Close()

	stream.Start()

	do_capture := true

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, os.Kill)
	go func() {
		<-sig
		do_capture = false
	}()
	// start_time := time.Now().UnixMilli()

	// capture frequency from mic every some interval, and add it to slice
	var captured_freqs []int
	for do_capture {
		// make a slice to store samples in
		samples := make([][2]float64, CAPTURE_RATE)

		// fmt.Println("Listening to microphone...")
		// puts stream samples into our samples slice, and pauses execution until samples is full
		stream.Stream(samples)

		// fmt.Printf("%v\n", samples)
		wavelength := estimateWavelength(samples)
		fmt.Printf("wl: %v\n", wavelength)
		frequency := SAMPLE_RATE / wavelength
		captured_freqs = append(captured_freqs, frequency)

		// fmt.Printf("wavelength: %v\n", wavelength)
		fmt.Printf("frequency: %v\n", frequency)
	}
	// fmt.Printf("%v\n", time.Now().UnixMilli()-start_time)

	fmt.Printf("%v\n", captured_freqs)

	var bits []byte
	for _, freq := range captured_freqs {
		if freq > 1500 {
			bits = append(bits, 1)
		} else if freq > 0 {
			bits = append(bits, 0)
		} else {
			bits = append(bits, 255)
		}
	}

	fmt.Printf("%v\n", bits)
}

// does some math to get the average wavelength over a slice of samples

func estimateWavelength(samples [][2]float64) int {
	// fmt.Printf("sample count: %v\n", len(samples))
	// find all the points in which the dB is zero
	var zero_indexes []int
	for i := range len(samples) - 1 {
		// compare this sample with next sample to see if the sign changes (hence a zero between)

		// fmt.Printf("multiplied: %v\n", samples[i][0]*samples[i+1][0])
		if samples[i][0]*samples[i+1][0] < 0 {
			// fmt.Printf("found zero between: %v(%v) %v(%v)\n", samples[i][0], i, samples[i+1][0], i+1)
			zero_indexes = append(zero_indexes, i)
		}
	}
	if len(zero_indexes) < 1 {
		return -1
	}

	// fmt.Printf("all zeros: %v\n", zero_indexes)
	// use the zero locations to get an average distance between 2 neighboring zeros (half wavelength)
	// could also skip every other zero and calculate full wavelength with this distances...
	total_dists := 0
	for i := range len(zero_indexes) - 1 {
		dist := zero_indexes[i+1] - zero_indexes[i]
		// fmt.Printf("dist(%v): %v\n", i, dist)
		total_dists += dist
	}
	average_dist := total_dists / len(zero_indexes)
	if average_dist < 1 {
		average_dist = -1
	}
	wavelength := average_dist * 2

	return wavelength
}
