package hqu

import "testing"

func TestMergeSort(t *testing.T) {
	testQsort([]int{}, t)
	testQsort([]int{1, 2}, t)
	testQsort([]int{2, 1}, t)
	testQsort([]int{3, 2, 1}, t)
	testQsort([]int{1, 2, 3}, t)
	testQsort([]int{4, 4, 5, 4}, t)
	testQsort([]int{1, 2, 3, 4, 5, 6, 7}, t)
	testQsort([]int{7, 6, 5, 4, 3, 2, 1, 0}, t)
	testQsort([]int{7, 6, 5, 4, 3, 9, 8, 7}, t)
	testQsort([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, t)
}

func testMergeSort(nums []int, t *testing.T) {
	raw := make([]int, len(nums))
	copy(raw, nums)
	MergeSort(nums)
	if !sorted(nums) {
		t.Logf("from:%v -> to:%v", raw, nums)
		t.Fail()
	}
}
