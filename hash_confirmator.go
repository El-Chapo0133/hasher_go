package hash_confirmator

import (
	"fmt"
	"sync"
	"io/ioutil"
	"strings"
	"strconv"
	"math"
	Graybane "../graybane"
)

const (
	NBTESTS = 100000
	CHARS = 10
	THREADS = 8
	NBTESTPERCORE = NBTESTS / THREADS
)
var (
	wg sync.WaitGroup
	words []string
)

type Confirmator struct {}
type Hash struct {
	hash []int
	length int
	input string
}
type List struct {
	hashs []Hash
}
type ListException struct {
	hashs [][]Hash
}

type Confirm interface {
	Test()
}

func (Confirmator) BrithdayAttackCount(hash []int) int {
	possibilities := (len(hash) - 4) * 94 * 2
	return int(math.Sqrt(float64(possibilities)))
}

func (Confirmator) Test() (bool, error) {
	list := new(List)
	listException := new(ListException)

	data, err := ioutil.ReadFile("C:\\Users\\levequel\\P_\\go_folder\\docs\\temp.txt")
	if err != nil {
		fmt.Println("Error reading the file")
		fmt.Println(err)
		return false, err
	}
	words = strings.Split(string(data), "\n")

	fmt.Println("Working on", strconv.Itoa(THREADS), "threads")
	initTests(list)
	output := confirm(list, &listException.hashs)
	fmt.Println(output)
	fmt.Println("Inputs tried:", strconv.Itoa(THREADS * NBTESTPERCORE))
	fmt.Println("Errors:", strconv.Itoa(len(listException.hashs)))
	return true, nil
}

func initTests(list *List) {
	for index := 1; index <= THREADS; index++ {
		coreInit(index, &list.hashs)
	}
	return
}

func coreInit(indexThread int, list *[]Hash) {
	defer fmt.Println("Init", strconv.Itoa(int(float64(THREADS * indexThread) / (THREADS * THREADS) * 100)),"%")
	hasher := new(Graybane.Hash)
	for index := 1; index <= NBTESTPERCORE; index++ {
		input := getString(index * indexThread)
		output := getHash(hasher, input)
		len := len(output)
		*list = append(*list, Hash{output, len, input})
	}
	return
}

func confirm(list *List, listExeption *[][]Hash) bool {
	var result bool = true
	parts := split(&list.hashs)
	for index := 1; index <= THREADS; index++ {
		fmt.Println("Confirm Thread n°", strconv.Itoa(index), "started")
		wg.Add(1)
		go coreConfirmList(parts[index - 1], &list.hashs, &result, listExeption, index)
	}
	wg.Wait()
	return result
}

func coreConfirmList(list []Hash, fullList *[]Hash, result *bool, listExeption *[][]Hash, indexThread int) {
	defer wg.Done(); fmt.Println("Confirm Thread n°", strconv.Itoa(indexThread)," finished")
	var LENGTH int = len(list)
	var FULLLISTLENGTH = len(*fullList)
	for index := 0; index < LENGTH; index++ {
		for index_t := 0; index_t < FULLLISTLENGTH; index_t++ {
			if Graybane.CompareTwoHash(list[index].hash, (*fullList)[index_t].hash) && list[index].input != (*fullList)[index_t].input {
				*result = false
				*listExeption = append(*listExeption, []Hash{list[index], (*fullList)[index_t]})
			} /*else {
				fmt.Println(Graybane.ConvertToString(list[index].hash), Graybane.ConvertToString((*fullList)[index_t].hash))
			}*/
		}
	}
	return
}

func split(list *[]Hash) [THREADS][]Hash {
	defer fmt.Println("Splitting made in", strconv.Itoa(THREADS), "parts")
	var LENGTH int = len(*list)
	var toReturn [THREADS][]Hash
	for index := 0; index < LENGTH; index++ {
		toReturn[index % THREADS] = append(toReturn[index % THREADS], (*list)[index])
	}
	return toReturn
}


func has(hash []int, list []Hash) bool {
	var LENGTH int = len(list)
	for i := 0; i < LENGTH; i++ {
		if equal(hash, list[i].hash) {
			return true
		}
	}
	return false
}

func equal(a, b []int) bool {
        if len(a) != len(b) {
                return false
        }
        for i, v := range a {
                if v != b[i] {
                        return false
                }
        }
        return true
}

func getHash(hasher *Graybane.Hash, input string) []int {
	return hasher.Get(input, 1)
}

func getString(index int) string {
	return words[index]
}
