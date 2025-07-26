package filesystem

import (
	"context"
	"time"
)

type FileMode int32

const (
	ModeUnknown FileMode = iota
	ModeFile             // 普通文件
	ModeDir              // 目录
	ModeDynamic          // 动态文件（通过 Generator 生成内容）
)

type FileNode interface {
	GetName() string
	IsDir() bool
	GetPath() string
	Read() ([]byte, error)
	Write(data []byte, flag string) error
	Stat() FileInfo
	ListChildren() ([]FileNode, error)
	Find(path string) (FileNode, error)
}

type FileInfo struct {
	Name       string    // "passwd"
	Path       string    // "/etc/passwd"
	Size       int64     // 字节数
	Mode       FileMode  // 权限+类型（Go 自带类型，方便格式化）
	OwnerUID   int32     // 文件拥有者
	OwnerGID   int32     // 所属组
	ModTime    time.Time // 最后修改时间
	AccessTime time.Time // 上次访问时间
	CreateTime time.Time // 创建时间
	NLink      int32     // 硬链接数（通常目录是 2+，文件是 1）
}

type FSService interface {
	Find(ctx context.Context, path string) (FileNode, error)
	CreateFile(ctx context.Context, path string, content []byte) error
	CreateDynamicFile(ctx context.Context, path string, generatorType string) error
	Mkdir(ctx context.Context, path string) error
	Remove(ctx context.Context, path string) error
	WriteFile(ctx context.Context, path string, data []byte, flag string) error
	ReadFile(ctx context.Context, path string) ([]byte, error)
	FindMetaData(ctx context.Context, path string) (FileInfo, error)
}
