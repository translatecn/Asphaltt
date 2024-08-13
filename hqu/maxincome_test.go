package hqu

import "testing"

func TestMaxIncome(t *testing.T) {
	testMaxincome([]int{7, -1, -1, -1}, 7, t)
	testMaxincome([]int{7}, 7, t)
	testMaxincome([]int{1, 2, 3, 4, 5}, 15, t)
	testMaxincome([]int{1, -1, 2, -3, -4, 5, 6}, 11, t)
	testMaxincome([]int{6, -1, 5, 4, -2}, 14, t)
	testMaxincome([]int{2, -3, -4, 3, 4, -2, -1}, 7, t)
}

func testMaxincome(nums []int, expected int, t *testing.T) {
	res := MaxIncome(nums)
	if res != expected {
		t.Logf("nums: %v, result:%d, expected:%d", nums, res, expected)
		t.Fail()
	}
}
