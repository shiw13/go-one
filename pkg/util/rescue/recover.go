package rescue

import "go-one/pkg/logger"

func Recover(cleanups ...func()) {
	for _, cleanup := range cleanups {
		cleanup()
	}

	if e := recover(); e != nil {
		logger.DPanicf("%v", e)
	}
}
