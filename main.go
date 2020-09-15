package main

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
)

/**********************************************************************
Types
***********************************************************************/
const stackHeightMin = 1
const stackHeightMax = 100
const stackCountMin = 1
const stackCountMax = 100

type pancakeStack struct {
	id        int
	flipCount int
	data      []bool
}

func (stack *pancakeStack) hasIntegrity(str string, min, max int) bool {
	patt := fmt.Sprintf("^[-+]{%d,%d}$", min, max)
	match, _ := regexp.Match(patt, []byte(str))
	if !match {
		log.Printf("ERROR: Line #%d failed to match required pattern of %s.\n discarding %v",
			stack.id, patt, stack.data)
	}
	return match
}

func (stack *pancakeStack) setDataFromString(str string, min, max int) bool {
	o := []bool{}
	plussign := rune(0x2b)
	success := stack.hasIntegrity(str, min, max)
	if success {
		for _, r := range str {
			o = append(o, r == plussign)
		}
	}
	stack.data = o
	return success
}

func (stack *pancakeStack) containsFalse() bool {
	for _, v := range stack.data {
		if !v {
			return true
		}
	}
	return false
}

func (stack *pancakeStack) flip(boundary int) {
	for i := 0; i <= boundary; i++ {
		stack.data[i] = !stack.data[i]
	}
	stack.flipCount++
}

func (stack *pancakeStack) flipAllToHappy(ch chan pancakeStack) {
	for stack.containsFalse() {
		boundary := -1
		started := false
		for i := range stack.data {
			if !stack.data[i] {
				if !started {
					started = true
				}
			}

			if started && stack.data[i] {
				boundary = i - 1
			} else if started && i == len(stack.data)-1 {
				boundary = i
			}

			if boundary > -1 {
				stack.flip(boundary)
				boundary = -1
			}
		}
	}
	ch <- *stack
}

/**********************************************************************
Sample Data
**********************************************************************/
func getData() []string {
	return []string{"5", "-", "-+", "+-", "+++", "--+-"}
}

/**********************************************************************
Helpers
**********************************************************************/
func printResult(set []pancakeStack) {
	for _, stack := range set {
		log.Printf("Case #%d: %d\n", stack.id+1, stack.flipCount)
	}
}

/**
 * check if set has numerical header and that it matches number of cases.
 * return dataset without header if okay. empty otherwise.
 */
func conformingSet(set []string, min, max int) []string {
	head, err := strconv.Atoi(set[0])
	if err != nil {
		log.Printf("ERROR: malformed data. non-numerical header?\n")
		return set[0:0]
	}
	tailset := set[1:]
	size := len(tailset)

	if head != size {
		log.Printf("ERROR: malformed data. header value does not match stack count\n")
		return set[0:0]
	}

	if size < min || size > max {
		log.Printf("ERROR: malformed data. Number of stacks not in range of %d to %d\n", min, max)
		return set[0:0]
	}
	return tailset
}

/**********************************************************************
Main
**********************************************************************/
func main() {
	cases := conformingSet(getData(), stackCountMin, stackCountMax)
	stacks := []pancakeStack{}

	ch := make(chan pancakeStack)
	gocount := 0 // keep track of how many results were expecting on the channel

	for i, c := range cases {
		stacks = append(stacks, pancakeStack{i, 0, nil})
		if stacks[i].setDataFromString(c, stackHeightMin, stackHeightMax) {
			go stacks[i].flipAllToHappy(ch)
			gocount++
		}
	}

	for i := 0; i < gocount; i++ {
		stack := <-ch
		stacks[stack.id] = stack
	}
	printResult(stacks)
}
