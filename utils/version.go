package utils

import (
	"runtime"
	"strconv"
	"strings"
)

// GoVer 当前Go环境版本
//
// 例如: "1.18.2"
func GoVer() string {
	return strings.Split(strings.TrimPrefix(runtime.Version(), "go"), " ")[0]
}

// GoVerLt 当前Go环境是否 小于 指定版本
//   - major 目标主版本号 Number类型(如:1)
//   - minor 目标次版本号 Number类型(如:18)
//   - patch 目标修订版本号(可选) Number类型(如:2) 不传默认为0
func GoVerLt[T Number](major, minor T, patch ...T) bool {
	if len(patch) == 0 {
		return checkGoVersion(int(major), int(minor), 0) < 0
	}
	return checkGoVersion(int(major), int(minor), int(patch[0])) < 0
}

// GoVerGt 当前Go环境是否 大于 指定版本
//   - major 目标主版本号 Number类型(如:1)
//   - minor 目标次版本号 Number类型(如:18)
//   - patch 目标修订版本号(可选) Number类型(如:2) 不传默认为0
func GoVerGt[T Number](major, minor T, patch ...T) bool {
	if len(patch) == 0 {
		return checkGoVersion(int(major), int(minor), 0) > 0
	}
	return checkGoVersion(int(major), int(minor), int(patch[0])) > 0
}

// GoVerEq 当前Go环境是否 等于 指定版本
//   - major 目标主版本号 Number类型(如:1)
//   - minor 目标次版本号 Number类型(如:18)
//   - patch 目标修订版本号(可选) Number类型(如:2) 不传默认为0
func GoVerEq[T Number](major, minor T, patch ...T) bool {
	if len(patch) == 0 {
		return checkGoVersion(int(major), int(minor), 0) == 0
	}
	return checkGoVersion(int(major), int(minor), int(patch[0])) == 0
}

// 比较 Golang 版本
//   - major 目标主版本号
//   - minor 目标次版本号
//   - patch 目标修订版本号
//
// 返回 当前版本 -1小于/0等于/1大于 传入版本
func checkGoVersion(major, minor, patch int) int {
	currVerArr := strings.Split(GoVer(), ".")
	currMajor, _ := strconv.Atoi(currVerArr[0])
	currMinor, _ := strconv.Atoi(currVerArr[1])
	currPatch, _ := strconv.Atoi(currVerArr[2])

	var checkMajor int
	if major < currMajor {
		checkMajor = 1
	} else if major > currMajor {
		checkMajor = -1
	}
	if checkMajor != 0 {
		return checkMajor
	}

	var checkMinor int
	if minor < currMinor {
		checkMinor = 1
	} else if minor > currMinor {
		checkMinor = -1
	}
	if checkMinor != 0 {
		return checkMinor
	}

	var checkPatch int
	if patch < currPatch {
		checkPatch = 1
	} else if patch > currPatch {
		checkPatch = -1
	}
	return checkPatch
}
