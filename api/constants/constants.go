package constants

const DEFAULT_STATUS = "pending"
const FAILED = "failed"
const SUCCESS = "success"
const TIME_OUT = "timeout"

var StatusMap = map[uint64]string{}

func init() {
	StatusMap = make(map[uint64]string)
	StatusMap[0] = FAILED
	StatusMap[1] = SUCCESS
}
