package controller

import "github.com/dlclark/regexp2"

var (
	Reip   = regexp2.MustCompile(`((2(5[0-5]|[0-4]\d))|[0-1]?\d{1,2})(\.((2(5[0-5]|[0-4]\d))|[0-1]?\d{1,2})){3}`, 0)
	Readdr = regexp2.MustCompile(`(?<=location":").*?(?=")`, 0)
)
