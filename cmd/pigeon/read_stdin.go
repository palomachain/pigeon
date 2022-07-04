package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"syscall"

	"golang.org/x/term"
)

func readPasswordStdin() string {
	bytePassword, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		panic(err)
	}
	return string(bytePassword)
}

func readLineFromStdin(isPassword bool) string {
	var text string
	if isPassword {
		text = readPasswordStdin()
	} else {
		var err error
		reader := bufio.NewReader(os.Stdin)
		text, err = reader.ReadString('\n')
		if err != nil {
			panic(err)
		}
	}
	text = strings.TrimSpace(text)
	return text
}

func readNTimes(n int, fnc func(string) bool) func(string) bool {
	count := 0
	return func(input string) bool {
		if count >= n {
			return true
		}
		if fnc(input) {
			return true
		}
		count++
		return false
	}
}

func doubleReadInput(prefix string, isPassword bool, failedAttempts int) string {
	for i := failedAttempts; i > 0; i-- {
		fmt.Print(prefix)
		text1 := readLineFromStdin(isPassword)
		fmt.Println("\nEnter it again")
		fmt.Print(prefix)
		text2 := readLineFromStdin(isPassword)
		if text1 == text2 {
			return text1
		}
		fmt.Println("\nInputs do not match. Please try again.")
	}
	panic("failed reading the input. failed too many times.")
}
