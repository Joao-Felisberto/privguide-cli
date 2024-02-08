package main

import "testing"

func TestA(t *testing.T) {
	t.Log("Hello world!")
}

func TestB(t *testing.T) {
	//t.Fatal("DON'T WORRY THIS IS A DRILL!")
	t.Log("FIXED!")
}
