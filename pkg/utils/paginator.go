package utils

import "fmt"

type PageNavigation struct {
	Max      int64
	Active   string
	PageList []string
}

func Paginate(current int64, total int64) PageNavigation {
	var stringRange []string
	var intRange []int64
	var delta int64 = 2
	var i, l int64

	left := current - delta
	right := current + delta + 1

	for i = 1; i <= total; i++ {
		if i == 1 || i == total || i >= left && i < right {
			intRange = append(intRange, i)
		}

	}

	for idx, i := range intRange {
		if idx != 0 {
			if i-l == 2 {
				stringRange = append(stringRange, fmt.Sprint(l+1))
			} else if i-l != 1 {
				stringRange = append(stringRange, "...")
			}
		}

		stringRange = append(stringRange, fmt.Sprint(i))
		l = i
	}

	return PageNavigation{
		Max:      total,
		Active:   fmt.Sprint(current),
		PageList: stringRange,
	}
}
