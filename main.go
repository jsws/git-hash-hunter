package main

import (
	"context"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"sync/atomic"
	"time"
)

var hashesCalculated uint64

func main() {

	if len(os.Args) < 2 {
		fmt.Println("Usage: githh <HASH> \ne.g.\n     githh 0000000")
		os.Exit(1)
	}

	hashToFind := strings.ToLower(os.Args[1])

	// Check given hash is valid.
	re := regexp.MustCompile(`^[0-9a-f]+$`)
	if re.MatchString(hashToFind) != true {
		fmt.Println("Please enter a valid hex string.")
		os.Exit(1)
	}

	// Start timer.
	start := time.Now()

	// Get commit at HEAD.
	commit, err := exec.Command("git", "cat-file", "commit", "HEAD").Output()
	if err != nil {
		fmt.Println("Error getting current commit.")
		os.Exit(1)
	}
	commitStr := string(commit)

	currentHash := calcHash(commitStr)
	fmt.Printf("Original hash: %s\n", currentHash)

	oldHeader, oldMessage := splitCommit(commitStr)

	fmt.Printf("Original commit message: \"%s\"\n", oldMessage)

	padding := make(chan string, 4000)
	results := make(chan string, 1)

	ctx, cancel := context.WithCancel(context.Background())

	// Padding generators, one generate even number permutations the other odd.
	go paddingGenerator(ctx, 0, padding)
	go paddingGenerator(ctx, 1, padding)

	// Workers to calculate hashes.
	go worker(hashToFind, oldMessage, oldHeader, padding, results)
	go worker(hashToFind, oldMessage, oldHeader, padding, results)
	go worker(hashToFind, oldMessage, oldHeader, padding, results)
	go worker(hashToFind, oldMessage, oldHeader, padding, results)

	result := <-results
	cancel()

	fmt.Printf("New hash: %s\n", calcHash(oldHeader+"\n\n"+result))
	fmt.Printf("New Message: \"%s\"\n", result)

	// Get previous commit date.
	previousDate, err := exec.Command("git", "--no-pager", "show", "-s", "--oneline", "--format='%cd'", "HEAD").Output()
	if err != nil {
		fmt.Println("Error getting commit date.")
		os.Exit(1)
	}

	cmd := exec.Command("git", "commit", "--amend", "--cleanup=verbatim", "--no-gpg-sign", "-m"+result)
	// Set environment variable to fake the commit date.
	cmd.Env = append(os.Environ(),
		"GIT_COMMITTER_DATE="+string(previousDate),
	)
	err = cmd.Run()
	if err != nil {
		fmt.Println("Error setting new commit. ")
		os.Exit(1)
	}

	elapsed := time.Since(start)
	fmt.Printf("Took %.3fs @ %.0f H/s\n", elapsed.Seconds(), float64(hashesCalculated)/elapsed.Seconds())

}

// Takes in a commit mesage in the format given from 'git cat-file commit HEAD'
// and returns the header and the commit message.
func splitCommit(commit string) (string, string) {

	// Commit message separated from header by two new lines.
	splitCommit := strings.Split(commit, "\n\n")

	// Header is the first element.
	header := strings.TrimSpace(splitCommit[0])

	// Commit message is everything after first element. Can be more than one
	// element if messagecontains \n\n.
	message := strings.TrimSpace(strings.Join(splitCommit[1:], ""))

	return header, message
}

func paddingGenerator(ctx context.Context, start int, padding chan string) {
	for i := start; true; i = i + 2 {
		i = i + 1
		select {
		case <-ctx.Done():
			return
		default:
			padding <- string(getPermutation(i))
		}
	}
}

// Returns a permutation of the 'chars'.
// Based on the algorithm to convert between radixes, will return a unique byte
// array for a unique 'num' parameter.
func getPermutation(num int) []byte {

	// Characters '\t', '\n' and ' '.
	chars := []byte{9, 10, 32}

	permutation := []byte{}
	for num > 0 {
		digit := int(num % len(chars))
		permutation = append(permutation, chars[digit])
		num = num / len(chars)
	}
	return permutation
}

func worker(hashToFind string, old_message string, old_head string, padding <-chan string, results chan<- string) {
	for pad := range padding {
		result := findHash(hashToFind, old_message, old_head, pad)

		if result != "" {
			results <- result
			break
		}

	}
}

// Calulates hash to see if
func findHash(hashToFind string, oldMessage string, oldHead string, padding string) string {
	atomic.AddUint64(&hashesCalculated, 1)
	newMessage := oldMessage + string(padding) + "\n"
	newHead := oldHead + "\n\n" + newMessage

	newHash := calcHash(newHead)

	if newHash[:len(hashToFind)] == hashToFind {
		return newMessage
	}
	return ""
}

func calcHash(head string) string {

	commit := "commit " + strconv.Itoa(len(head)) + "\000" + string(head)

	h := sha1.New()
	h.Write([]byte(commit))
	bs := h.Sum(nil)

	return hex.EncodeToString(bs)
}
