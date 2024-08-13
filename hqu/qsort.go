package hqu

// QuickSort sorts a number array with quick-sort algorithm
func QuickSort(nums []int) {
	if len(nums) < 2 {
		return
	}

	quickPartition(nums, 0, len(nums)-1)
	// quickSort(nums, 0, len(nums)-1)
}

func quickPartition(nums []int, lo, hi int) {
	if lo >= hi {
		return
	}

	pivot := nums[(lo+hi+1)/2] // +1 to pick the middle one
	i, j := lo, hi
	for i < j {
		for nums[i] < pivot && i < j { // do not equal
			i++
		}
		for nums[j] > pivot && i < j { // do not equal
			j--
		}
		if i < j {
			nums[i], nums[j] = nums[j], nums[i]
			i, j = i+1, j-1
		}
	}
	quickPartition(nums, lo, i-1)
	quickPartition(nums, i, hi) // do not skip the ith number
}

func quickSort(nums []int, lo, hi int) {
	if lo >= hi {
		return
	}

	idx, pivot := lo+1, nums[lo]
	for i := idx; i <= hi; i++ {
		if nums[i] < pivot && i != idx {
			nums[i], nums[idx] = nums[idx], nums[i]
			idx++
		}
	}

	nums[lo], nums[idx-1] = nums[idx-1], nums[lo]
	quickSort(nums, lo, idx-1)
	quickSort(nums, idx, hi)
}
