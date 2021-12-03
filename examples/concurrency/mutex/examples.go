package mu

import "fmt"

const (
	mutexLocked      = 1 << iota // mutex is locked // 1
	mutexWoken                   // 2
	mutexStarving                // 4
	mutexWaiterShift = iota      // 3
)

func printMutexMath() {
	fmt.Printf("%b %d \n", mutexLocked, mutexLocked)
	fmt.Printf("%b %d \n", mutexWoken, mutexWoken)
	fmt.Printf("%b %d \n", mutexStarving, mutexStarving)

	fmt.Printf("%b %d \n \n", mutexLocked | mutexStarving, mutexLocked | mutexStarving)

	// old & mutexStarving
	fmt.Printf("%b %d \n", mutexLocked&mutexStarving, mutexLocked&mutexStarving)
	fmt.Printf("%b %d \n", mutexWoken&mutexStarving, mutexWoken&mutexStarving)
	fmt.Printf("%b %d \n", mutexStarving&mutexStarving, mutexStarving&mutexStarving)
	fmt.Printf("%b %d \n \n", mutexWaiterShift&mutexStarving, mutexWaiterShift&mutexStarving)

	// old >> mutexWaiterShift
	fmt.Printf("%b %d \n", mutexLocked&mutexWoken, mutexLocked&mutexWoken)
	fmt.Printf("%b %d \n", mutexWoken&mutexWoken, mutexWoken&mutexWoken)
	fmt.Printf("%b %d \n \n", mutexStarving&mutexWoken, mutexStarving&mutexWoken)

	// old & mutexWoken
	fmt.Printf("%b %d \n", mutexLocked&mutexWoken, mutexLocked&mutexWoken)
	fmt.Printf("%b %d \n", mutexWoken&mutexWoken, mutexWoken&mutexWoken)
	fmt.Printf("%b %d \n \n", mutexStarving&mutexWoken, mutexStarving&mutexWoken)

	// old | mutexWoken
	fmt.Printf("%b %d \n", mutexLocked|mutexWoken, mutexLocked|mutexWoken)
	fmt.Printf("%b %d \n", mutexWoken|mutexWoken, mutexWoken|mutexWoken)
	fmt.Printf("%b %d \n \n", mutexStarving|mutexWoken, mutexStarving|mutexWoken)

	// old&(mutexLocked|mutexStarving)
	fmt.Printf("%b %d \n", mutexLocked&(mutexLocked|mutexStarving), mutexLocked&(mutexLocked|mutexStarving))
	fmt.Printf("%b %d \n", mutexWoken&(mutexLocked|mutexStarving), mutexWoken&(mutexLocked|mutexStarving))
	fmt.Printf("%b %d \n", mutexStarving&(mutexLocked|mutexStarving), mutexWoken&(mutexLocked|mutexStarving))
}
