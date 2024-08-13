package hqu

import "testing"

func TestTopN(t *testing.T) {
	arr := []int{34, 87, 22, 43, 88, 23, 55, 23, 32, 10, 87, 100, 67, 97, 89, 91, 87, 0}
	t.Log(arr)
	t.Log(TopN(arr, 4))
	t.Log(arr)
}

func TestTopK(t *testing.T) {
	arr := []int{34, 87, 22, 43, 88, 23, 55, 23, 32, 10, 87, 100, 67, 97, 89, 91, 87, 0}
	testTopK(arr, 4, t)

	arr = []int{44, 44, 48, 44}
	testTopK(arr, 1, t)
}

func testTopK(nums []int, k int, t *testing.T) {
	t.Log(nums)
	t.Log(TopK(nums, k))
	t.Log(nums)
}
