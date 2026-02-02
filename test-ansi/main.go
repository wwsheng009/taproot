package main

import (
	"fmt"
)

func main() {
	// Direct ANSI escape codes for testing
	
	// Reset
	fmt.Println("\033[0m")
	
	// Test 1: Red foreground, Green background with upper half block
	fmt.Println("\033[38;5;196m\033[48;5;46m▀\033[0m")
	
	// Test 2: RGB foreground, RGB background (truecolor)
	fmt.Println("\033[38;2;255;0;0m\033[48;2;0;255;0m▀\033[0m")
	
	// Test 3: A simple pattern
	for i := 0; i < 10; i++ {
		// Gradient from red to blue
		r := 255 - (i * 25)
		b := i * 25
		fmt.Printf("\033[38;2;%d;0;%dm\033[48;2;0;%d;255m▀", r, b, b)
	}
	fmt.Println("\033[0m")
	
	// Test 4: Full row
	fmt.Println("\n--- Full row test ---")
	for y := 0; y < 5; y++ {
		for x := 0; x < 40; x++ {
			// Create a gradient pattern
			r := (x * 6) % 256
			g := (y * 50) % 256
			b := ((x + y) * 3) % 256
			fmt.Printf("\033[38;2;%d;%d;%dm\033[48;2;%d;%d;%dm▀", r, g, b, b, g, r)
		}
		fmt.Println("\033[0m")
	}
	
	fmt.Println("\n--- Done ---")
}
