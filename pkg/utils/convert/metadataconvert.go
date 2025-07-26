package convert

import (
	"time"

	"github.com/BBitQNull/SSHoneyNet/core/filesystem"
	fs_Pb "github.com/BBitQNull/SSHoneyNet/pb/filesystem"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func ConvertMetadataToPb(src filesystem.FileInfo) *fs_Pb.Metadata {
	return &fs_Pb.Metadata{
		Name:       src.Name,
		Path:       src.Path,
		Size:       src.Size,
		Filemode:   int32(src.Mode),
		OwnerUID:   src.OwnerUID,
		OwnerGID:   src.OwnerGID,
		ModTime:    timestamppb.New(src.ModTime),
		CreateTime: timestamppb.New(src.CreateTime),
		AccessTime: timestamppb.New(src.AccessTime),
		NLink:      src.NLink,
	}
}

func ProtoTimestampToTime(ts *timestamppb.Timestamp) time.Time {
	if ts == nil {
		return time.Time{}
	}
	t := ts.AsTime()
	return t
}

func ConvertMetadataFromPb(src *fs_Pb.Metadata) filesystem.FileInfo {
	if src == nil {
		return filesystem.FileInfo{}
	}
	return filesystem.FileInfo{
		Name:       src.Name,
		Path:       src.Path,
		Size:       src.Size,
		Mode:       filesystem.FileMode(src.Filemode),
		OwnerUID:   src.OwnerUID,
		OwnerGID:   src.OwnerGID,
		ModTime:    ProtoTimestampToTime(src.ModTime),
		CreateTime: ProtoTimestampToTime(src.CreateTime),
		AccessTime: ProtoTimestampToTime(src.AccessTime),
		NLink:      src.NLink,
	}
}
