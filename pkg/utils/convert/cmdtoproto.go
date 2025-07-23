package convert

import (
	"github.com/BBitQNull/SSHoneyNet/core/commandparser"
	pb "github.com/BBitQNull/SSHoneyNet/pb/cmdparser"
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
