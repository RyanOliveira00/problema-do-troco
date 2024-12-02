package main

import (
	"fmt"
	"math"
	"sort"
	"time"
)

const (
	MAX_RECURSION_DEPTH = 1000
	TIMEOUT_DURATION    = 5 * time.Second
)

type timeoutChecker struct {
	startTime time.Time
}

func newTimeoutChecker() *timeoutChecker {
	return &timeoutChecker{startTime: time.Now()}
}

func (t *timeoutChecker) hasTimedOut() bool {
	return time.Since(t.startTime) > TIMEOUT_DURATION
}

// 1. Implementação Recursiva
func coinChangeRecursive(coins []int, amount int) int {
	if amount == 0 {
		return 0
	}

	// Ordena as moedas em ordem decrescente para otimização
	sortedCoins := make([]int, len(coins))
	copy(sortedCoins, coins)
	sort.Sort(sort.Reverse(sort.IntSlice(sortedCoins)))

	checker := newTimeoutChecker()
	result := coinChangeRecursiveHelper(sortedCoins, amount, 0, checker)

	if checker.hasTimedOut() {
		return -3 // Código para timeout
	}
	return result
}

func coinChangeRecursiveHelper(coins []int, amount, depth int, checker *timeoutChecker) int {
	if checker.hasTimedOut() {
		return -3
	}

	if depth > MAX_RECURSION_DEPTH {
		return -2
	}

	if amount == 0 {
		return 0
	}
	if amount < 0 {
		return -1
	}

	minCoins := math.MaxInt32
	for _, coin := range coins {
		if coin > amount {
			continue
		}
		result := coinChangeRecursiveHelper(coins, amount-coin, depth+1, checker)
		if result >= 0 && result < minCoins {
			minCoins = result + 1
		}
	}

	if minCoins == math.MaxInt32 {
		return -1
	}
	return minCoins
}

// 2. Implementação com Memorização
func coinChangeMemoized(coins []int, amount int) int {
	if amount == 0 {
		return 0
	}

	sortedCoins := make([]int, len(coins))
	copy(sortedCoins, coins)
	sort.Sort(sort.Reverse(sort.IntSlice(sortedCoins)))

	memo := make(map[int]int, amount+1)
	checker := newTimeoutChecker()
	result := coinChangeMemoHelper(sortedCoins, amount, memo, checker)

	if checker.hasTimedOut() {
		return -3
	}
	return result
}

func coinChangeMemoHelper(coins []int, amount int, memo map[int]int, checker *timeoutChecker) int {
	if checker.hasTimedOut() {
		return -3
	}

	if amount == 0 {
		return 0
	}
	if amount < 0 {
		return -1
	}

	if val, exists := memo[amount]; exists {
		return val
	}

	minCoins := math.MaxInt32
	for _, coin := range coins {
		if coin > amount {
			continue
		}
		result := coinChangeMemoHelper(coins, amount-coin, memo, checker)
		if result >= 0 && result < minCoins {
			minCoins = result + 1
		}
	}

	if minCoins == math.MaxInt32 {
		minCoins = -1
	}
	memo[amount] = minCoins
	return minCoins
}

// 3. Implementação Iterativa
func coinChangeIterative(coins []int, amount int) int {
	if amount == 0 {
		return 0
	}

	dp := make([]int, amount+1)
	for i := range dp {
		dp[i] = math.MaxInt32
	}
	dp[0] = 0

	sortedCoins := make([]int, len(coins))
	copy(sortedCoins, coins)
	sort.Ints(sortedCoins)

	checker := newTimeoutChecker()
	for i := 1; i <= amount; i++ {
		if checker.hasTimedOut() {
			return -3
		}

		for _, coin := range sortedCoins {
			if coin > i {
				break
			}
			if dp[i-coin] != math.MaxInt32 {
				dp[i] = min(dp[i], dp[i-coin]+1)
			}
		}
	}

	if dp[amount] == math.MaxInt32 {
		return -1
	}
	return dp[amount]
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func TestImplementations(coins []int, amount int) {
	fmt.Printf("\nTestando para amount=%d e coins=%v\n", amount, coins)

	if amount < 0 {
		fmt.Println("Valor inválido: amount deve ser positivo")
		return
	}
	if len(coins) == 0 {
		fmt.Println("Valor inválido: array de moedas está vazio")
		return
	}

	// Teste Recursivo
	start := time.Now()
	result1 := coinChangeRecursive(coins, amount)
	duration1 := time.Since(start)
	printResult("Recursivo", result1, duration1)

	// Teste Memorizado
	start = time.Now()
	result2 := coinChangeMemoized(coins, amount)
	duration2 := time.Since(start)
	printResult("Memorizado", result2, duration2)

	// Teste Iterativo
	start = time.Now()
	result3 := coinChangeIterative(coins, amount)
	duration3 := time.Since(start)
	printResult("Iterativo", result3, duration3)
}

func printResult(method string, result int, duration time.Duration) {
	switch result {
	case -3:
		fmt.Printf("%s: Timeout - execução muito longa (tempo: %v)\n", method, duration)
	case -2:
		fmt.Printf("%s: Profundidade máxima de recursão atingida (tempo: %v)\n", method, duration)
	case -1:
		fmt.Printf("%s: Impossível fazer o troco (tempo: %v)\n", method, duration)
	default:
		fmt.Printf("%s: %d moedas (tempo: %v)\n", method, result, duration)
	}
}

func main() {
	TestImplementations([]int{1, 2, 5}, 11)
	TestImplementations([]int{2, 5, 10, 20, 50}, 100)
	TestImplementations([]int{1, 2, 3, 4}, 999)
	TestImplementations([]int{2, 5, 10}, 3)
}