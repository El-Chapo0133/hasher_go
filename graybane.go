package graybane

import (
	"sync"
	//"fmt"
	//"strconv"
)

var (
	NBCHAR = float64(1024 - len(STARTUP))
	NBCHARINT = int(NBCHAR)
	NBLOOPS int = 128
	STARTUP string = "$21%"
	STARTUPINTS []int = []int {36, 50, 49, 37}
)
const (
	PARTS = 8
)

type Hash struct {}

type Hasher interface {
	Get(input string, ref int) []int
}

var wg_hash sync.WaitGroup

func (Hash) Get(input string, ref int) []int {
	prehash := convert(input)
	hashOnParts := splitOnParts(prehash)
	for index := 0; index < PARTS; index++ {
		wg_hash.Add(1)
		go threadCalc(&hashOnParts[index], ref, index)
	}
	wg_hash.Wait()
	hash := collapseParts(hashOnParts)
	collapse(&hash)
	return hash
}
func threadCalc(part *[]int, ref int, index int) {
	defer wg_hash.Done()
	for index := 0; index < NBLOOPS; index++ {
		applyRef(part, ref)
		applyFibonacci(part)
		applyCollatz(part)
		applyLog(part)
	}
	return
}
func threadConvert(part *[]int, second []int) {
	defer wg_hash.Done()
	var INITLENGTH int = len(second)
	for index := 0; index < INITLENGTH; index++ {
		(*part)[index % NBCHARINT] += second[index % len(second)] * (index + 2)
	}
}
func splitOnParts(input []int) [PARTS][]int {
	var toReturn [PARTS][]int
	for index, value := range input {
		toReturn[index % PARTS] = append(toReturn[index % PARTS], value)
	}
	return toReturn
}
func collapseParts(input [PARTS][]int) []int {
	var toReturn []int
	for index := 0; index < PARTS; index++ {
		toReturn = append(toReturn, input[index]...)
	}
	return toReturn
}
func collapse(input *[]int) {
	*input = append(STARTUPINTS, *input...)
	return
}
func convert(input string) []int {
	var INITLENGTH, length, TEMPLENGTH int = len(input), 0, int(NBCHAR * 1.5)
	var temp string = input
	for length < TEMPLENGTH {
		temp += input
		length += INITLENGTH
	}
	toReturn, second := split(temp)
	for index := 0; index < PARTS; index++ {
		wg_hash.Add(1)
		go threadConvert(&toReturn, second[index])
	}
	wg_hash.Wait()
	return toReturn
}
func split(input string) ([]int, [PARTS][]int) {
	var INITLENGTH int = len(input)
	var first []int
	var second [PARTS][]int
	for index := 0; index < INITLENGTH; index++ {
		if index < NBCHARINT {
			first = append(first, int(input[index]))
		} else {
			second[index % PARTS] = append(second[index % PARTS], int(input[index]))
		}
	}
	return first, second
}
func convertToInt(input string) []int {
	var INITLENGTH int = len(input)
	var toReturn []int
	for index := 0; index < INITLENGTH; index++ {
		toReturn = append(toReturn, int(input[index]))
	}
	return toReturn
}
func applyLog(input *[]int) {
	var INITLENGTH int = len(*input)
	for index := 0; index < INITLENGTH; index++ {
		(*input)[index] = ((*input)[index] % (*input)[(index + 1) % INITLENGTH] * index % 94) + 33
	}
	return
}
func applyRef(input *[]int, ref int) {
	for index, value := range *input {
		(*input)[index] = value + ref
	}
	return
}
func applyFibonacci(input *[]int) {
	var INITLENGTH int = len(*input) - 1
	for index := 1; index < INITLENGTH; index++ {
		(*input)[index + 1] += (*input)[index] + (*input)[index - 1]
	}
	return
}
func applyCollatz(input *[]int) {
	var INITLENGTH int = len(*input)
	for index := 0; index < INITLENGTH; index++ {
		(*input)[index] = 3 * (*input)[index] + 1
	}
	return
}
func ConvertToString(input []int) string {
	var toReturn string = ""
	for _, value := range input {
		toReturn += string(value)
	}
	return toReturn
}
func CompareTwoHash(input1 []int, input2 []int) bool {
	if len(input1) != len(input2) {
		return false
	}
	for index, value := range input1 {
		if value != input2[index] {
			return false
		}
	}
	return true
}