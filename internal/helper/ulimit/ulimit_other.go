// +build !linux

package ulimit

func GetUlimit() (uint64, error) {
	return 2048, nil
}
