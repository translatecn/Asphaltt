package hqu

// MergeSort sorts a number array with merge-sort algorithm
func MergeSort(nums []int) {
	mergeSort(nums, make([]int, 0, len(nums)), 0, len(nums)-1)
}

func mergeSort(nums, buff []int, lo, hi int) {
	if hi-lo < 2 {
		return
	}

	mid := (lo + hi + 1) / 2 // the middle index
	mergeSort(nums, buff, lo, mid-1)
	mergeSort(nums, buff, mid, hi)

	if nums[mid-1] < nums[mid] { // it's sorted
		return
	}

	buff = buff[:0] // reuse a buffer
	i, j := lo, mid
	for i < mid && j <= hi {
		if nums[i] <= nums[j] {
			buff = append(buff, nums[i])
			i++
		} else {
			buff = append(buff, nums[j])
			j++
		}
	}
	copy(nums[lo:hi+1], buff) // restore to the original array
}
