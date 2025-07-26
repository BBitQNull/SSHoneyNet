package convert

import (
	"github.com/BBitQNull/SSHoneyNet/core/filesystem"
	fs_Pb "github.com/BBitQNull/SSHoneyNet/pb/filesystem"
)

func ConvertChildrenToPb(src []filesystem.FileNodeInfo) []*fs_Pb.FileNodeInfo {
	if src == nil {
		return nil
	}
	result := make([]*fs_Pb.FileNodeInfo, 0, len(src))
	for _, item := range src {
		result = append(result, &fs_Pb.FileNodeInfo{
			Name:    item.Name,
			Path:    item.Path,
			IsDir:   item.IsDir,
			Mode:    item.Mode,
			ModTime: item.ModTime,
			Size:    item.Size,
		})
	}
	return result
}

func ConvertChildrenFromPb(src []*fs_Pb.FileNodeInfo) []filesystem.FileNodeInfo {
	if src == nil {
		return nil
	}
	result := make([]filesystem.FileNodeInfo, 0, len(src))
	for _, item := range src {
		result = append(result, filesystem.FileNodeInfo{
			Name:    item.Name,
			Path:    item.Path,
			IsDir:   item.IsDir,
			Mode:    item.Mode,
			ModTime: item.ModTime,
			Size:    item.Size,
		})
	}
	return result
}
