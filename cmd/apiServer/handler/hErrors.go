package handler

var (
	LanguageNotSupported = newError(4002, "Language not supported")
	PayloadNotValid      = newError(4003, "Payload not valid")
	TooMuchOutPutErr     = newError(4005, "Too much output")
	TimeOutErr           = newError(4006, "Time out")
	RunCodeErr           = newError(5001, "Run code error")
)

var codeMap = make(map[int]int)

type errorForFront struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

func newError(code int, msg string) errorForFront {
	if _, got := codeMap[code]; got {
		panic("repeat code")
	}
	codeMap[code] = 1
	return errorForFront{code, msg}
}
