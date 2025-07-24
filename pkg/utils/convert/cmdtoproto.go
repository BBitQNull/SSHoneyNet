package convert

import (
	"github.com/BBitQNull/SSHoneyNet/core/commandparser"
	pb "github.com/BBitQNull/SSHoneyNet/pb/cmdparser"
	pbdis "github.com/BBitQNull/SSHoneyNet/pb/dispatcher"
)

func ConvertScript(src *commandparser.Script) *pb.Script {
	return &pb.Script{
		Lines: ConvertCmdLines(src.Lines), // slice
	}
}

func ConvertCmdLines(src []*commandparser.CommandLine) []*pb.CommandLine {
	result := make([]*pb.CommandLine, len(src))
	for i, comment := range src {
		result[i] = &pb.CommandLine{
			Pipeline: ConvertPipeline(comment.Pipeline),
			Redir:    ConvertRedir(comment.Redir),
		}
	}
	return result
}

func ConvertPipeline(src []*commandparser.Command) []*pb.Command {
	result := make([]*pb.Command, len(src))
	for i, comment := range src {
		result[i] = &pb.Command{
			Flags: ConvertFlagWithValue(comment.Flags),
			Args:  ConvertArgument(comment.Args),
		}
	}
	return result
}

func ConvertFlagWithValue(src []commandparser.FlagWithValue) []*pb.FlagWithValue {
	result := make([]*pb.FlagWithValue, len(src))
	for i, comment := range src {
		result[i] = &pb.FlagWithValue{
			Name:  comment.Name,
			Value: comment.Value,
		}
	}
	return result
}

func ConvertArgument(src []commandparser.Argument) []*pb.Argument {
	result := make([]*pb.Argument, len(src))
	for i, comment := range src {
		result[i] = &pb.Argument{
			Value: comment.Value,
		}
	}
	return result
}

func ConvertRedir(src *commandparser.Redirection) *pb.Redirection {
	return &pb.Redirection{
		File: src.File,
	}
}

func ConvertScriptFormpb(src *pbdis.Script) *commandparser.Script {
	return &commandparser.Script{
		Lines: ConvertCmdLinesFrompb(src.Lines), // slice
	}
}

func ConvertCmdLinesFrompb(src []*pbdis.CommandLine) []*commandparser.CommandLine {
	result := make([]*commandparser.CommandLine, len(src))
	for i, comment := range src {
		result[i] = &commandparser.CommandLine{
			Pipeline: ConvertPipelineFrompb(comment.Pipeline),
			Redir:    ConvertRedirFrompb(comment.Redir),
		}
	}
	return result
}

func ConvertPipelineFrompb(src []*pbdis.Command) []*commandparser.Command {
	result := make([]*commandparser.Command, len(src))
	for i, comment := range src {
		result[i] = &commandparser.Command{
			Flags: ConvertFlagWithValueFrompb(comment.Flags),
			Args:  ConvertArgumentFrompb(comment.Args),
		}
	}
	return result
}

func ConvertFlagWithValueFrompb(src []*pbdis.FlagWithValue) []commandparser.FlagWithValue {
	result := make([]commandparser.FlagWithValue, len(src))
	for i, comment := range src {
		result[i] = commandparser.FlagWithValue{
			Name:  comment.Name,
			Value: comment.Value,
		}
	}
	return result
}

func ConvertArgumentFrompb(src []*pbdis.Argument) []commandparser.Argument {
	result := make([]commandparser.Argument, len(src))
	for i, comment := range src {
		result[i] = commandparser.Argument{
			Value: comment.Value,
		}
	}
	return result
}

func ConvertRedirFrompb(src *pbdis.Redirection) *commandparser.Redirection {
	return &commandparser.Redirection{
		File: src.File,
	}
}
