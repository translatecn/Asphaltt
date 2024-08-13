func isAnagram(s string, t string) bool {
    if len(s) != len(t) { // not the same length, not anagram
        return false
    }
    return compare(countChar(s), countChar(t))
}

// countChar counts char of `s` in lower case
func countChar(s string) []int {
    res := [26]int{}
    for i := 0; i < len(s); i++ {
        res[s[i]-'a']++
    }
    return res[:]
}

// compare compares whether the countings are same
func compare(c0, c1 []int) bool {
    for k := range c0 {
        if c0[k] != c1[k] {
            return false
        }
    }
    return true
}
