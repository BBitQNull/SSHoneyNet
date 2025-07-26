package convert

import (
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
