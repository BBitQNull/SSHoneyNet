package convert

import (
	"github.com/BBitQNull/SSHoneyNet/core/commandparser"
	pbcom "github.com/BBitQNull/SSHoneyNet/pb/common"
)

func ConvertScript(src *commandparser.Script) *pbcom.Script {
	if src == nil {
		return nil
	}
	return &pbcom.Script{
		Lines: ConvertCmdLines(src.Lines), // slice
	}
}

func ConvertCmdLines(src []*commandparser.CommandLine) []*pbcom.CommandLine {
	result := make([]*pbcom.CommandLine, len(src))
	for i, comment := range src {
		result[i] = &pbcom.CommandLine{
			Pipeline: ConvertPipeline(comment.Pipeline),
			Redir:    ConvertRedir(comment.Redir),
		}
	}
	return result
}

func ConvertPipeline(src []*commandparser.Command) []*pbcom.Command {
	result := make([]*pbcom.Command, len(src))
	for i, comment := range src {
		result[i] = &pbcom.Command{
			Name:  comment.Name,
			Flags: ConvertFlagWithValue(comment.Flags),
			Args:  ConvertArgument(comment.Args),
		}
	}
	return result
}

func ConvertFlagWithValue(src []commandparser.FlagWithValue) []*pbcom.FlagWithValue {
	result := make([]*pbcom.FlagWithValue, len(src))
	for i, comment := range src {
		result[i] = &pbcom.FlagWithValue{
			Name:  comment.Name,
			Value: comment.Value,
		}
	}
	return result
}

func ConvertArgument(src []commandparser.Argument) []*pbcom.Argument {
	result := make([]*pbcom.Argument, len(src))
	for i, comment := range src {
		result[i] = &pbcom.Argument{
			Value: comment.Value,
		}
	}
	return result
}

func ConvertRedir(src *commandparser.Redirection) *pbcom.Redirection {
	if src == nil {
		return nil
	}
	return &pbcom.Redirection{
		File: src.File,
	}
}

func ConvertScriptFormpb(src *pbcom.Script) *commandparser.Script {
	if src == nil || len(src.Lines) == 0 {
		return nil
	}
	return &commandparser.Script{
		Lines: ConvertCmdLinesFrompb(src.Lines), // slice
	}
}

func ConvertCmdLinesFrompb(src []*pbcom.CommandLine) []*commandparser.CommandLine {
	result := make([]*commandparser.CommandLine, len(src))
	for i, comment := range src {
		result[i] = &commandparser.CommandLine{
			Pipeline: ConvertPipelineFrompb(comment.Pipeline),
			Redir:    ConvertRedirFrompb(comment.Redir),
		}
	}
	return result
}

func ConvertPipelineFrompb(src []*pbcom.Command) []*commandparser.Command {
	result := make([]*commandparser.Command, len(src))
	for i, comment := range src {
		result[i] = &commandparser.Command{
			Name:  comment.Name,
			Flags: ConvertFlagWithValueFrompb(comment.Flags),
			Args:  ConvertArgumentFrompb(comment.Args),
		}
	}
	return result
}

func ConvertFlagWithValueFrompb(src []*pbcom.FlagWithValue) []commandparser.FlagWithValue {
	result := make([]commandparser.FlagWithValue, len(src))
	for i, comment := range src {
		result[i] = commandparser.FlagWithValue{
			Name:  comment.Name,
			Value: comment.Value,
		}
	}
	return result
}

func ConvertArgumentFrompb(src []*pbcom.Argument) []commandparser.Argument {
	result := make([]commandparser.Argument, len(src))
	for i, comment := range src {
		result[i] = commandparser.Argument{
			Value: comment.Value,
		}
	}
	return result
}

func ConvertRedirFrompb(src *pbcom.Redirection) *commandparser.Redirection {
	if src == nil {
		return nil
	}
	return &commandparser.Redirection{
		File: src.File,
	}
}
