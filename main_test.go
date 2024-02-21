package main

import (
	"slices"
	"testing"
)

func Test_deDuplicate(t *testing.T) {
    intList := []int{ 1, 2, 2, 3, 4, 4, 5 }
    intWant := []int{ 1, 2, 3, 4, 5 }
    intDeDup := deDuplicate(intList)
    
    if !slices.Equal(intDeDup, intWant) {
        t.Fatalf("got %v, expected %v", intDeDup, intWant)
    }

    strList := []string{ "a", "b", "b", "c", "d", "d", "e" }
    strWant := []string{ "a", "b", "c", "d", "e" }
    strDeDup := deDuplicate(strList)
    
    if !slices.Equal(strDeDup, strWant) {
        t.Fatalf("got %v, expected %v", strDeDup, strWant)
    }

    strList = []string{ "a", "a", "a", "a", "a", "a", "a", "a", "a", "a" }
    strWant = []string{ "a" }
    strDeDup = deDuplicate(strList)
    
    if !slices.Equal(strDeDup, strWant) {
        t.Fatalf("got %v, expected %v", strDeDup, strWant)
    }
}
