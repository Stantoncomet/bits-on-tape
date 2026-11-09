package main

import (
	"fmt"

	"github.com/BlueishLeaf/go-sound/output"
	"github.com/BlueishLeaf/go-sound/sounds"
)

func binToAudio(data []byte) {
	var note_bits []float64
	for _, bit := range data {
		note_bits = append(note_bits, 1000*(float64(bit))+1000)
	}

	s := soundSequence(note_bits, 100)
	output.Play(s)
}

// fs is a short way to write an array of floats.
func fs(fs ...float64) []float64 {
	return fs
}

func soundSequence(freqs []float64, duration float64) sounds.Sound {
	var notes []sounds.Sound
	for _, f := range freqs {
		notes = append(notes, sounds.NewTimedSound(sounds.NewSineWave(f), duration))
	}
	s := sounds.ConcatSounds(notes...)
	return s
}

func strToBin(str string) []byte {
	var out_data []byte
	for _, c := range str {
		char_bin := fmt.Sprintf("%08b", byte(c))
		for _, b := range char_bin {
			out_data = append(out_data, byte(b)-48)
		}
	}
	return out_data
}

func strToHex(str string) []byte {
	out_data := []byte{}
	for _, byte := range []byte(str) {
		out_data = append(out_data, byte/16)
		out_data = append(out_data, byte%16)
	}
	return out_data
}

func strToByte(str string) []byte {
	out_data := []byte{}
	for _, byte := range []byte(str) {
		out_data = append(out_data, byte)
	}
	return out_data
}
