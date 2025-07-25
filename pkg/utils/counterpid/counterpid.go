package counterpid

import "sync"

var (
	shellPidCounter = make(map[int64]int64)
	counterLock     sync.Mutex
)

func GeneratePidFromShell(shellPid int64) int64 {
	counterLock.Lock()
	defer counterLock.Unlock()

	// 如果不存在，初始化
	if _, ok := shellPidCounter[shellPid]; !ok {
		// 从 shellPid + 1 开始
		shellPidCounter[shellPid] = shellPid + 1
	} else {
		// 否则在上次基础上递增
		shellPidCounter[shellPid]++
	}

	return shellPidCounter[shellPid]
}
