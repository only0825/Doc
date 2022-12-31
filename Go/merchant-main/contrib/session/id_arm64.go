package session

/* 
#include <mach/mach_time.h>
C.mach_absolute_time()
*/


import (
	"time"
)

func Cputicks() (t uint64) {

	sec := time.Now().UnixNano() / 1e6
	return uint64(sec)
}
