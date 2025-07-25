package convert

import (
	"time"

	"github.com/BBitQNull/SSHoneyNet/core/procsystem"
	proc_pb "github.com/BBitQNull/SSHoneyNet/pb/procsystem"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func ConvertPcbFromEndpoint(src *procsystem.PCB) *proc_pb.Pcb {
	return &proc_pb.Pcb{
		Pid:       src.PID,
		Ppid:      src.PPID,
		Command:   src.Command,
		State:     src.State,
		Tgid:      src.TGID,
		User:      src.User,
		Createdat: timestamppb.New(src.CreatedAt),
		Exitedat:  timestamppb.New(*src.ExitedAt),
	}
}

func ConvertPcbFromPb(src *proc_pb.Pcb) *procsystem.PCB {
	return &procsystem.PCB{
		PID:     src.Pid,
		PPID:    src.Ppid,
		Command: src.Command,
		TGID:    src.Tgid,
		User:    src.User,
		State:   src.State,
		CreatedAt: func(ts *timestamppb.Timestamp) time.Time {
			if ts != nil {
				return ts.AsTime()
			}
			return time.Time{}
		}(src.Createdat),
		ExitedAt: ProtoTimestampToPtr(src.Exitedat),
	}
}

func ProtoTimestampToPtr(ts *timestamppb.Timestamp) *time.Time {
	if ts == nil {
		return nil
	}
	t := ts.AsTime()
	return &t
}

func ConvertPcbListFromEndpoint(src []*procsystem.PCB) []*proc_pb.Pcb {
	result := make([]*proc_pb.Pcb, len(src))
	for i, item := range src {
		result[i] = ConvertPcbFromEndpoint(item)
	}
	return result
}

func ConvertPcbLitsFromPb(src []*proc_pb.Pcb) []*procsystem.PCB {
	result := make([]*procsystem.PCB, len(src))
	for i, item := range src {
		result[i] = ConvertPcbFromPb(item)
	}
	return result
}
