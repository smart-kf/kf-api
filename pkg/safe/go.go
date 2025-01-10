package safe

import (
	"fmt"
	"runtime"
)

func Go(fn func()) {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				var b = make([]byte, 1024)
				n := runtime.Stack(b, true)
				fmt.Printf("error panic:%+v \n%s\n", err, string(b[:n]))
			}

			go fn()
		}()
	}()
}
