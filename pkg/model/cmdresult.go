package model

type CmdResult struct {
	Output     string
	ExitCode   uint32
	ErrMsg     error
	Log        string
	NextPrompt bool
}
