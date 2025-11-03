package main

import (
    "fmt"
	//"time"
    //"slices"
    //"crypto/sha256"

)

const s string = "constant"

func main() {
    /* fmt.Println(s)

    const n = 2
    const d = 3e8 / n
    fmt.Println(d)
    fmt.Println(int64(d))

	var i int =0
    for ; i <10 ; i++ {
		fmt.Print(i," ")
	}

	for n := range 6 {
        if n%2 == 0 {continue}
		
        fmt.Println(n)
    }

	if num := 9; num < 0 {
        fmt.Println(num, "is negative")
    } else if num < 10 {
        fmt.Println(num, "has 1 digit")
    } else {
        fmt.Println(num, "has multiple digits")
    } 

	t := time.Now()
	fmt.Println(t)
	fmt.Println(t.Weekday())
    
    var nums = make([]int, 3)
    jj:=make([]int, 3)
    jj[0] = 1
    nums[0] = 1
    nums[1] = 2
    nums = append(nums,jj...)
    fmt.Println(nums)
    */
    maped:=make(map[int]int)
    maped[1] = 18
    maped[2] = 2
    key,value:= maped[1]
    fmt.Println(key,value)



    /*
    var m = map[string]int{"a":1,"b":2}
    fmt.Println(m["c"])

    h := sha256.Sum256([]byte("xv"))
    fmt.Println(h)
    */
}
 
