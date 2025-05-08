package common

func Swap[T any](items []T, a, b uint) {
	tempA := items[a]
	items[a] = items[b]
	items[b] = tempA
}
