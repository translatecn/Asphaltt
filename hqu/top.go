package hqu

// TopN finds the count `n` greatest integers in the array `arr`
func TopN(arr []int, n int) []int {
	if len(arr) == 0 || n < 1 || n >= len(arr) {
		return arr
	}

	top(arr, 0, len(arr)-1, n)
	return arr[:n]
}

// TopK finds the `k`th greatest integers in the array `arr`.
func TopK(arr []int, k int) (int, bool) {
	if len(arr) == 0 || k < 1 || k >= len(arr) {
		return 0, false
	}

	top(arr, 0, len(arr)-1, k)
	return arr[k-1], true
}

func top(nums []int, lo, hi, t int) {
	if lo >= hi {
		return
	}

	pivot := nums[(lo+hi+1)/2] // +1 to pick the middle one
	i, j := lo, hi
	for i < j {
		for nums[i] > pivot && i < j { // do not equal
			i++
		}
		for nums[j] < pivot && i < j { // do not equal
			j--
		}
		if i < j {
			nums[i], nums[j] = nums[j], nums[i]
			i, j = i+1, j-1
		}
	}

	if i == t {
		return
	}

	if i > t {
		top(nums, lo, i, t)
	} else {
		top(nums, i, hi, t)
	}
}
