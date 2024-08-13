package hqu

// MaxIncome gets the max income from a number array
func MaxIncome(nums []int) int {
	if len(nums) == 0 {
		panic("MaxIncome unsupports empty number array")
	}
	if len(nums) == 1 {
		return nums[0]
	}

	maxIncome := int(-1 << 63)
	left, right := 0, len(nums)-1
	for i := left; i < len(nums); i++ {
		tmp := sum(nums[left : i+1])
		if tmp >= maxIncome {
			right, maxIncome = i, tmp
		}
	}

	for j := right; j >= 0; j-- {
		tmp := sum(nums[j : right+1])
		if tmp > maxIncome {
			left, maxIncome = j, tmp
		}
	}
	return maxIncome
}

func sum(nums []int) int {
	res := 0
	for i := range nums {
		res += nums[i]
	}
	return res
}
