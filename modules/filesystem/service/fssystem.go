package fs_service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/BBitQNull/SSHoneyNet/core/filesystem"
)

// 基础文件类
type BaseFile struct {
	Name     string
	Parent   *Directory
	Metadata filesystem.FileInfo
}

// 常规文件
type RegularFile struct {
	BaseFile
	Content []byte
}

// 目录
type Directory struct {
	BaseFile
	Children  map[string]filesystem.FileNode
	Generator string `json:"generator,omitempty"`
}

// 动态文件
type DynamicFile struct {
	BaseFile
	Generator func() ([]byte, error)
}

// 链接文件
type Symlink struct {
	BaseFile
	TargetPath string // 指向的路径，比如 "/bin/busybox"
}

// 文件系统
type FileSystem struct {
	mu   sync.RWMutex
	Root filesystem.FileNode
}

// 文件系统服务
type fsService struct {
	fs *FileSystem
}

func NewFSService(fs *FileSystem) filesystem.FSService {
	return &fsService{fs: fs}
}

func buildFromJSON(jf *filesystem.JSONFile, parent *Directory) filesystem.FileNode {
	modTime := time.Unix(int64(jf.MTime), 0)
	meta := filesystem.FileInfo{
		Name:       jf.Name,
		Path:       "", // 可以在 build 完后补充 path
		Size:       jf.Size,
		Mode:       filesystem.ModeFromString(jf.Mode), // 你需要实现这个辅助函数
		OwnerUID:   0,
		OwnerGID:   0,
		ModTime:    modTime,
		CreateTime: modTime,
		AccessTime: modTime,
		NLink:      1,
	}

	switch jf.Mode {
	case "dir":
		dir := &Directory{
			BaseFile: BaseFile{
				Name:     jf.Name,
				Parent:   parent,
				Metadata: meta,
			},
			Children:  make(map[string]filesystem.FileNode),
			Generator: "",
		}
		for _, child := range jf.Children {
			childNode := buildFromJSON(child, dir)
			dir.Children[child.Name] = childNode
		}
		return dir

	case "file":
		return &RegularFile{
			BaseFile: BaseFile{
				Name:     jf.Name,
				Parent:   parent,
				Metadata: meta,
			},
			Content: []byte{}, // 真实内容你可能没有
		}

	case "link":
		return &Symlink{
			BaseFile: BaseFile{
				Name:     jf.Name,
				Parent:   parent,
				Metadata: meta,
			},
			TargetPath: jf.Target,
		}
	case "dynamic":
		dir := &Directory{
			BaseFile: BaseFile{
				Name:     jf.Name,
				Parent:   parent,
				Metadata: meta,
			},
			Children:  make(map[string]filesystem.FileNode),
			Generator: jf.Generator, // 记录下来，暂不执行
		}
		return dir

	default:
		panic(fmt.Sprintf("unknown mode: %s", jf.Mode))
	}
}

func (fs *FileSystem) populatePaths() {
	var walk func(node filesystem.FileNode, path string)
	walk = func(node filesystem.FileNode, path string) {
		meta := node.Stat()
		meta.Path = path
		node.SetMeta(meta)

		if dir, ok := node.(*Directory); ok {
			for name, child := range dir.Children {
				childPath := filepath.Join(path, name)
				walk(child, childPath)
			}
		}
	}
	walk(fs.Root, "/")
}

func NewFileSystem(jsonPath string) (*FileSystem, error) {
	data, err := os.ReadFile(jsonPath)
	if err != nil {
		return nil, err
	}

	var root filesystem.JSONFile
	if err := json.Unmarshal(data, &root); err != nil {
		return nil, err
	}

	rootNode := buildFromJSON(&root, nil)
	rootDir, ok := rootNode.(*Directory)
	if !ok {
		return nil, fmt.Errorf("root is not a directory")
	}

	fs := &FileSystem{
		Root: rootDir,
	}
	// 可选：为所有节点填充 Path 字段（递归）
	fs.populatePaths()

	return fs, nil
}

var generatorRegistry = map[string]func() ([]byte, error){
	"proc_pid_status": func() ([]byte, error) {
		return []byte("Name:\tmyproc\nState:\tR (running)\n"), nil
	},
	"proc_pid_cmdline": func() ([]byte, error) {
		return []byte("/bin/bash"), nil
	},
}

// 常规文件实现filenode接口
func (f *RegularFile) GetName() string {
	return f.Name
}
func (f *RegularFile) IsDir() bool {
	return false
}
func (f *RegularFile) GetPath() string {
	if f.Parent == nil {
		return "/"
	}
	return filepath.Join(f.Parent.GetPath(), f.Name)
}
func (f *RegularFile) Read() ([]byte, error) {
	return f.Content, nil
}
func (f *RegularFile) Write(data []byte, flag string) error {
	f.Metadata.ModTime = time.Now()
	switch flag {
	case "overwrite":
		f.Content = make([]byte, len(data))
		copy(f.Content, data)
		f.Metadata.Size = int64(len(data))
	case "append":
		f.Content = append(f.Content, data...)
		f.Metadata.Size = int64(len(f.Content))
	default:
		return fmt.Errorf("unsupported write flag: %s", flag)
	}
	return nil
}
func (f *RegularFile) Stat() filesystem.FileInfo {
	return f.Metadata
}
func (f *RegularFile) ListChildren() ([]filesystem.FileNode, error) {
	return nil, fmt.Errorf("cannot list children inside regular file '%s'", f.Name)
}
func (f *RegularFile) Find(path string) (filesystem.FileNode, error) {
	return nil, fmt.Errorf("cannot find path '%s' inside regular file '%s'", path, f.Name)
}
func (f *RegularFile) SetMeta(meta filesystem.FileInfo) {
	f.Metadata = meta
}

// 目录实现filenode接口
func (d *Directory) GetName() string {
	return d.Name
}
func (d *Directory) IsDir() bool {
	return true
}
func (d *Directory) GetPath() string {
	if d.Parent == nil {
		return "/"
	}
	return filepath.Join(d.Parent.GetPath(), d.Name)
}
func (d *Directory) Read() ([]byte, error) {
	return nil, fmt.Errorf("cannot read from a directory: %s", d.Name)
}
func (d *Directory) Write(data []byte, flag string) error {
	return fmt.Errorf("cannot write to a directory: %s", d.Name)
}
func (d *Directory) Stat() filesystem.FileInfo {
	return d.Metadata
}
func (d *Directory) ListChildren() ([]filesystem.FileNode, error) {
	if d.Children == nil {
		return []filesystem.FileNode{}, nil
	}
	result := make([]filesystem.FileNode, 0, len(d.Children))
	for _, item := range d.Children {
		if item == nil {
			continue
		}
		result = append(result, item)
	}
	return result, nil
}
func (d *Directory) Find(path string) (filesystem.FileNode, error) {
	cleanedPath := filepath.Clean(path)
	if cleanedPath == "." || cleanedPath == "/" {
		return d, nil
	}
	parts := strings.Split(cleanedPath, "/")
	current := filesystem.FileNode(d)
	for _, part := range parts {
		switch part {
		case ".":
			continue
		case "..":
			// 回到父目录（当前类型为目录）
			dir, ok := current.(*Directory)
			if !ok {
				return nil, fmt.Errorf("cannot go up from non-directory: %s", current.GetName())
			}
			if dir.Parent == nil {
				return nil, fmt.Errorf("already at root directory")
			}
			current = dir.Parent
		default:
			dir, ok := current.(*Directory)
			if !ok {
				return nil, fmt.Errorf("not a directory: %s", current.GetName())
			}
			child, exists := dir.Children[part]
			if !exists {
				return nil, fmt.Errorf("path not found: %s", part)
			}
			current = child
		}
	}
	return current, nil
}
func (d *Directory) SetMeta(meta filesystem.FileInfo) {
	d.Metadata = meta
}

// 动态文件实现filenode接口
func (df *DynamicFile) GetName() string {
	return df.Name
}
func (df *DynamicFile) IsDir() bool {
	return false
}
func (df *DynamicFile) GetPath() string {
	if df.Parent == nil {
		return "/"
	}
	return filepath.Join(df.Parent.GetPath(), df.Name)
}
func (df *DynamicFile) Read() ([]byte, error) {
	if df.Generator == nil {
		return nil, fmt.Errorf("no generator function defined for dynamic file: %s", df.Name)
	}
	return df.Generator()
}
func (df *DynamicFile) Write(data []byte, flag string) error {
	return fmt.Errorf("cannot write to dynamic file: %s", df.Name)
}
func (df *DynamicFile) Stat() filesystem.FileInfo {
	return df.Metadata
}
func (df *DynamicFile) ListChildren() ([]filesystem.FileNode, error) {
	return nil, fmt.Errorf("cannot list children inside DynamicFile file '%s'", df.Name)
}
func (df *DynamicFile) Find(path string) (filesystem.FileNode, error) {
	return nil, fmt.Errorf("cannot find path '%s' inside DynamicFile file '%s'", path, df.Name)
}
func (df *DynamicFile) SetMeta(meta filesystem.FileInfo) {
	df.Metadata = meta
}

// 链接文件实现filenode接口
func (s *Symlink) GetName() string {
	return s.Name
}
func (s *Symlink) IsDir() bool {
	return false
}
func (s *Symlink) GetPath() string {
	if s.Parent == nil {
		return "/"
	}
	return filepath.Join(s.Parent.GetPath(), s.Name)
}
func (s *Symlink) Read() ([]byte, error) {
	return []byte(s.TargetPath), nil
}
func (s *Symlink) Write(data []byte, flag string) error {
	return fmt.Errorf("cannot write to a symlink '%s'", s.Name)
}
func (s *Symlink) Stat() filesystem.FileInfo {
	return s.Metadata
}
func (s *Symlink) ListChildren() ([]filesystem.FileNode, error) {
	return nil, fmt.Errorf("cannot list children inside symlink '%s'", s.Name)
}
func (s *Symlink) Find(path string) (filesystem.FileNode, error) {
	return nil, fmt.Errorf("cannot find path '%s' inside symlink '%s'", path, s.Name)
}
func (s *Symlink) SetMeta(meta filesystem.FileInfo) {
	s.Metadata = meta
}

// FileSystem实现方法
func (fs *FileSystem) Find(path string) (filesystem.FileNode, error) {
	fs.mu.RLock()
	defer fs.mu.RUnlock()
	return fs.Root.Find(path)
}
func (fs *FileSystem) CreateFile(path string, content []byte) error {
	fs.mu.Lock()
	defer fs.mu.Unlock()
	cleanedPath := filepath.Clean(path)
	if cleanedPath == "." || cleanedPath == "/" {
		return fmt.Errorf("invalid file path: %s", path)
	}
	dirPath := filepath.Dir(cleanedPath)
	fileName := filepath.Base(cleanedPath)
	current := fs.Root
	dirParts := strings.Split(dirPath, "/")
	for _, part := range dirParts {
		if part == "" {
			continue
		}

		dir, ok := current.(*Directory)
		if !ok {
			return fmt.Errorf("not a directory in path: %s", part)
		}

		child, exists := dir.Children[part]
		if !exists {
			// 如果不存在，则新建目录
			newDir := &Directory{
				BaseFile: BaseFile{
					Name:   part,
					Parent: dir,
					Metadata: filesystem.FileInfo{
						Name:       part,
						Path:       filepath.Join(dir.GetPath(), part),
						Mode:       filesystem.ModeDir,
						OwnerUID:   0,
						OwnerGID:   0,
						ModTime:    time.Now(),
						CreateTime: time.Now(),
						AccessTime: time.Now(),
						NLink:      3,
					},
				},
				Children: make(map[string]filesystem.FileNode),
			}
			dir.Children[part] = newDir
			current = newDir
		} else {
			// 向下进入
			current = child
		}
	}

	// 确保最终目录是 *Directory
	parentDir, ok := current.(*Directory)
	if !ok {
		return fmt.Errorf("invalid parent directory for file: %s", dirPath)
	}

	// 检查文件是否已存在
	if _, exists := parentDir.Children[fileName]; exists {
		return fmt.Errorf("file already exists: %s", fileName)
	}

	// 创建 RegularFile
	newFile := &RegularFile{
		BaseFile: BaseFile{
			Name:   fileName,
			Parent: parentDir,
			Metadata: filesystem.FileInfo{
				Name:       fileName,
				Path:       filepath.Join(parentDir.GetPath(), fileName),
				Mode:       filesystem.ModeFile,
				Size:       int64(len(content)),
				OwnerUID:   0,
				OwnerGID:   0,
				ModTime:    time.Now(),
				CreateTime: time.Now(),
				AccessTime: time.Now(),
				NLink:      1,
			},
		},
		Content: make([]byte, len(content)),
	}
	copy(newFile.Content, content)

	// 加入到父目录中
	parentDir.Children[fileName] = newFile

	return nil
}
func (fs *FileSystem) CreateDynamicFile(path string, generator func() ([]byte, error)) error {
	fs.mu.Lock()
	defer fs.mu.Unlock()
	cleanedPath := filepath.Clean(path)
	if cleanedPath == "." || cleanedPath == "/" {
		return fmt.Errorf("invalid file path: %s", path)
	}
	dirPath := filepath.Dir(cleanedPath)
	fileName := filepath.Base(cleanedPath)
	current := fs.Root
	dirParts := strings.Split(dirPath, "/")
	for _, part := range dirParts {
		if part == "" {
			continue
		}

		dir, ok := current.(*Directory)
		if !ok {
			return fmt.Errorf("not a directory in path: %s", part)
		}

		child, exists := dir.Children[part]
		if !exists {
			// 如果不存在，则新建目录
			newDir := &Directory{
				BaseFile: BaseFile{
					Name:   part,
					Parent: dir,
					Metadata: filesystem.FileInfo{
						Name:       part,
						Path:       filepath.Join(dir.GetPath(), part),
						Mode:       filesystem.ModeDir,
						OwnerUID:   0,
						OwnerGID:   0,
						ModTime:    time.Now(),
						CreateTime: time.Now(),
						AccessTime: time.Now(),
						NLink:      3,
					},
				},
				Children: make(map[string]filesystem.FileNode),
			}
			dir.Children[part] = newDir
			current = newDir
		} else {
			// 向下进入
			current = child
		}
	}

	// 确保最终目录是 *Directory
	parentDir, ok := current.(*Directory)
	if !ok {
		return fmt.Errorf("invalid parent directory for file: %s", dirPath)
	}

	// 检查文件是否已存在
	if _, exists := parentDir.Children[fileName]; exists {
		return fmt.Errorf("file already exists: %s", fileName)
	}

	// 创建 RegularFile
	newFile := &DynamicFile{
		BaseFile: BaseFile{
			Name:   fileName,
			Parent: parentDir,
			Metadata: filesystem.FileInfo{
				Name:       fileName,
				Path:       filepath.Join(parentDir.GetPath(), fileName),
				Mode:       filesystem.ModeDynamic,
				OwnerUID:   0,
				OwnerGID:   0,
				ModTime:    time.Now(),
				CreateTime: time.Now(),
				AccessTime: time.Now(),
				NLink:      1,
			},
		},
		Generator: generator,
	}

	// 加入到父目录中
	parentDir.Children[fileName] = newFile

	return nil
}

func (fs *FileSystem) Mkdir(path string) error {
	fs.mu.Lock()
	defer fs.mu.Unlock()
	cleanedPath := filepath.Clean(path)
	if cleanedPath == "." || cleanedPath == "/" {
		return fmt.Errorf("invalid file path: %s", path)
	}
	dirParts := strings.Split(strings.Trim(cleanedPath, "/"), "/")
	current := fs.Root
	for _, part := range dirParts {
		if part == "" {
			continue
		}

		dir, ok := current.(*Directory)
		if !ok {
			return fmt.Errorf("not a directory in path: %s", part)
		}

		child, exists := dir.Children[part]
		if !exists {
			// 如果不存在，则新建目录
			newDir := &Directory{
				BaseFile: BaseFile{
					Name:   part,
					Parent: dir,
					Metadata: filesystem.FileInfo{
						Name:       part,
						Path:       filepath.Join(dir.GetPath(), part),
						Mode:       filesystem.ModeDir,
						OwnerUID:   0,
						OwnerGID:   0,
						ModTime:    time.Now(),
						CreateTime: time.Now(),
						AccessTime: time.Now(),
						NLink:      3,
					},
				},
				Children: make(map[string]filesystem.FileNode),
			}
			dir.Children[part] = newDir
			current = newDir
		} else {
			if childDir, ok := child.(*Directory); ok {
				current = childDir
			} else {
				return fmt.Errorf("path segment '%s' exists and is not a directory", part)
			}
		}
	}
	return nil
}
func (fs *FileSystem) traverseToDir(path string) (*Directory, error) {
	cleaned := filepath.Clean(path)
	if cleaned == "/" || cleaned == "." {
		return fs.Root.(*Directory), nil
	}
	current := fs.Root.(*Directory)
	for _, part := range strings.Split(cleaned, "/") {
		if part == "" {
			continue
		}
		child, ok := current.Children[part]
		if !ok {
			return nil, fmt.Errorf("directory not found: %s", part)
		}
		dir, ok := child.(*Directory)
		if !ok {
			return nil, fmt.Errorf("not a directory: %s", part)
		}
		current = dir
	}
	return current, nil
}

func (fs *FileSystem) Remove(path string) error {
	fs.mu.Lock()
	defer fs.mu.Unlock()

	cleanedPath := filepath.Clean(path)
	if cleanedPath == "." || cleanedPath == "/" {
		return fmt.Errorf("cannot remove root or current directory: %s", path)
	}

	dirPath := filepath.Dir(cleanedPath)
	baseName := filepath.Base(cleanedPath)

	// 找到父目录
	parentDir, err := fs.traverseToDir(dirPath)
	if err != nil {
		return fmt.Errorf("parent directory not found: %w", err)
	}

	target, exists := parentDir.Children[baseName]
	if !exists {
		return fmt.Errorf("file or directory not found: %s", path)
	}

	// 如果是目录，且不为空，则报错（非递归删除）
	if dir, ok := target.(*Directory); ok {
		if len(dir.Children) > 0 {
			return fmt.Errorf("directory not empty: %s", path)
		}
	}

	// 删除该文件/目录
	delete(parentDir.Children, baseName)
	return nil
}

func (fs *FileSystem) WriteFile(path string, data []byte, flag string) error {
	fs.mu.Lock()
	defer fs.mu.Unlock()
	node, err := fs.Find(path)
	if err != nil {
		return err
	}
	return node.Write(data, flag)
}
func (fs *FileSystem) ReadFile(path string) ([]byte, error) {
	fs.mu.RLock()
	defer fs.mu.RUnlock()
	node, err := fs.Find(path)
	if err != nil {
		return nil, err
	}
	return node.Read()
}

// FSServive
func (fs *fsService) Find(ctx context.Context, path string) (filesystem.FileNode, error) {
	return fs.fs.Find(path)
}
func (fs *fsService) CreateFile(ctx context.Context, path string, content []byte) error {
	return fs.fs.CreateFile(path, content)
}
func (fs *fsService) CreateDynamicFile(ctx context.Context, path string, generatorType string) error {
	generator, ok := generatorRegistry[generatorType]
	if !ok {
		return fmt.Errorf("unknown generator type: %s", generatorType)
	}
	return fs.fs.CreateDynamicFile(path, generator)
}
func (fs *fsService) Mkdir(ctx context.Context, path string) error {
	return fs.fs.Mkdir(path)
}
func (fs *fsService) Remove(ctx context.Context, path string) error {
	return fs.fs.Remove(path)
}
func (fs *fsService) WriteFile(ctx context.Context, path string, data []byte, flag string) error {
	return fs.fs.WriteFile(path, data, flag)
}
func (fs *fsService) ReadFile(ctx context.Context, path string) ([]byte, error) {
	return fs.fs.ReadFile(path)
}
func (fs *fsService) FindMetaData(ctx context.Context, path string) (filesystem.FileInfo, error) {
	resp, err := fs.fs.Find(path)
	if err != nil {
		return filesystem.FileInfo{}, fmt.Errorf("file is not found: %s", path)
	}
	return resp.Stat(), nil
}
func (fs *fsService) ListChildren(ctx context.Context, path string) ([]filesystem.FileNode, error) {
	node, err := fs.fs.Find(path)
	if err != nil {
		log.Printf("ListChildren: path not found: %s", path)
		return nil, err
	}
	dir, ok := node.(*Directory)
	if !ok {
		log.Printf("ListChildren: not a directory: %s", path)
		return nil, fmt.Errorf("not a directory: %s", path)
	}
	if !node.IsDir() {
		return nil, fmt.Errorf("not a directory")
	}
	return dir.ListChildren()
}
