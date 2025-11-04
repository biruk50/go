package main

import (
	"fmt"
)

func sum(nums []int ) int {
	res:=0
	for _,num := range nums {
		res+=num
	}

	return res
}
func main() {

	nums:= make([]int,3)
	nums = append(nums,1,2,3)

	fmt.Println(sum(nums))
}