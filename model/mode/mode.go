package mode

import "strings"

type Mode uint8

const (
	Server Mode = iota + 1
	Update
	Crawl
	Convert
	Check
)

func NewFromInt(i int) Mode {
	switch i {
	case 1:
		return Server
	case 2:
		return Update
	case 3:
		return Crawl
	case 4:
		return Convert
	case 5:
		return Check
	default:
		return Server
	}
}

func NewFromString(str string) Mode {
	switch strings.ToLower(str) {
	case "server":
		return Server
	case "update":
		return Update
	case "crawl":
		return Crawl
	case "convert":
		return Convert
	case "check":
		return Check
	default:
		return Server
	}
}

func RunTask[T any](currentMode Mode, taskMode Mode, task func(T, bool), arg T) {
	switch currentMode {
	case Server:
		task(arg, false)
	case Update:
		task(arg, false)
	case Crawl:
		if taskMode == Convert {
			return
		}
		task(arg, false)
	case Convert:
		if taskMode == Crawl {
			return
		}
		task(arg, false)
	case Check:
		task(arg, true)
	default:
	}
}
