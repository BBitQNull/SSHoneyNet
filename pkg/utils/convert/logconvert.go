package convert

import (
	"github.com/BBitQNull/SSHoneyNet/core/log"
	log_Pb "github.com/BBitQNull/SSHoneyNet/pb/log"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func ConvertLogEntryFormPb(src *log_Pb.LogEntry) log.LogEntry {
	return log.LogEntry{
		Timestamp: ProtoTimestampToTime(src.Timestamp),
		SessionID: src.SessionID,
		UserInput: src.UserInput,
		Output:    src.Output,
		IP:        src.Ip,
	}
}

func ConvertLogEntryToPb(src log.LogEntry) *log_Pb.LogEntry {
	return &log_Pb.LogEntry{
		Timestamp: timestamppb.New(src.Timestamp),
		SessionID: src.SessionID,
		UserInput: src.UserInput,
		Output:    src.Output,
		Ip:        src.IP,
	}
}
