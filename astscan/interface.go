package astscan

import (
	"fmt"
	"go/ast"
	"regexp"
	"unicode"
)

var (
	reChinesePunctuation = regexp.MustCompile("/·|，|。|《|》|‘|’|”|“|；|：|【|】|？|（|）|、/")
)

// types for scanning
const (
	TypeString  = "string"
	TypeComment = "comment"
)

// Item for scanning callback
type Item struct {
	Pkg, Fileline, Type, Value string
	Node                       ast.Node
}

func (i Item) String() string {
	return fmt.Sprintf("{Pkg:%s Fileline:%s Type:%s Value:%s}",
		i.Pkg, i.Fileline, i.Type, i.Value)
}

// Callback is the scanning callbacker
type Callback func(Item)

// Checker is the scanning checker
type Checker func(s string) bool

// CheckChinese check whether s contains Chinese character
func CheckChinese(s string) bool {
	chars := unicode.Scripts["Han"]
	for _, r := range s {
		if unicode.Is(chars, r) || reChinesePunctuation.Match([]byte(string(r))) {
			return true
		}
	}
	return false
}
