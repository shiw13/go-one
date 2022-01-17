package rescue

import "github.com/shiw13/go-one/pkg/logger"

func Recover(cleanups ...func()) {
	for _, cleanup := range cleanups {
		cleanup()
	}

	if e := recover(); e != nil {
		logger.DPanicf("%v", e)
	}
}
