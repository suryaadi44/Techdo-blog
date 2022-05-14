package utils

type PageNavigation struct {
	Max      int64
	Active   int64
	PageList []int64
}

func Paginate(current int64, total int64) PageNavigation {
	var stringRange []int64
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
				stringRange = append(stringRange, l+1)
			} else if i-l != 1 {
				stringRange = append(stringRange, -1)
			}
		}

		stringRange = append(stringRange, i)
		l = i
	}

	return PageNavigation{
		Max:      total,
		Active:   current,
		PageList: stringRange,
	}
}
