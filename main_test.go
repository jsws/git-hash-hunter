package main

import (
	"reflect"
	"testing"
)

func TestCalcHash(t *testing.T) {
	commit := "tree 0de2407152f79bbc7fd18bde361e4572dfdef8f8\n" +
		"parent 543a377b8ace7d6f8cf4169d556cb4b70abe35d2\n" +
		"author jsws <test@users.noreply.github.com> 1574977452 +0000\n" +
		"committer jsws <test@users.noreply.github.com> 1574977452 +0000\n\n" +
		"First commit\n"

	hash := calcHash(commit)
	expected := "0ddaac8671212b14ee0e80871a0795d0964f0718"
	if hash != expected {
		t.Errorf("calcHash was incorrect, got: %s, want: %s.", hash, expected)
	}
}

func TestSplitCommit(t *testing.T) {
	commit := "tree 0de2407152f79bbc7fd18bde361e4572dfdef8f8\n" +
		"parent 543a377b8ace7d6f8cf4169d556cb4b70abe35d2\n" +
		"author jsws <test@users.noreply.github.com> 1574977452 +0000\n" +
		"committer jsws <test@users.noreply.github.com> 1574977452 +0000\n\n" +
		"First commit\n"

	header, message := splitCommit(commit)

	expectedMessage := "First commit"
	if message != expectedMessage {
		t.Errorf("splitCommit was incorrect, got: %s, want: %s.", message, expectedMessage)
	}

	expectedHeader := "tree 0de2407152f79bbc7fd18bde361e4572dfdef8f8\n" +
		"parent 543a377b8ace7d6f8cf4169d556cb4b70abe35d2\n" +
		"author jsws <test@users.noreply.github.com> 1574977452 +0000\n" +
		"committer jsws <test@users.noreply.github.com> 1574977452 +0000"
	if header != expectedHeader {
		t.Errorf("splitCommit was incorrect, got: %s, want: %s.", header, expectedHeader)
	}
}

func TestGetPermutation(t *testing.T) {

	// Test
	permutation := getPermutation(2)
	expectedPermutation := []byte{32}

	if !reflect.DeepEqual(permutation, expectedPermutation) {
		t.Errorf("splitCommit was incorrect, got: %s, want: %s.", permutation, expectedPermutation)
	}
}
