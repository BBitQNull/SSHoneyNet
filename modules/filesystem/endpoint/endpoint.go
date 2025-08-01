package fs_endpoint

import (
	"context"
	"errors"
	"log"

	"github.com/BBitQNull/SSHoneyNet/core/filesystem"
	"github.com/go-kit/kit/endpoint"
)

type FSRequest struct {
	Path          string
	Content       []byte
	Flag          string
	GeneratorType string
}

type FSResponse struct {
	Result   []byte
	Metadata filesystem.FileInfo
	Children []filesystem.FileNodeInfo
}

func MakeCreateFileEndpoint(svc filesystem.FSService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(FSRequest)
		if !ok {
			log.Printf("failed to assert MakeCreateFileEndpoint:")
			return nil, errors.New("failed to assert MakeCreateFileEndpoint")
		}
		err := svc.CreateFile(ctx, req.Path, req.Content)
		if err != nil {
			log.Printf("error: MakeCreateFileEndpoint:")
			return nil, err
		}
		return FSResponse{}, nil
	}
}

func MakeFindMetaDataEndpoint(svc filesystem.FSService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(FSRequest)
		if !ok {
			log.Printf("failed to assert MakeFindMetaDataEndpoint:")
			return nil, errors.New("failed to assert MakeFindMetaDataEndpoint")
		}
		resp, err := svc.FindMetaData(ctx, req.Path)
		if err != nil {
			log.Printf("error: MakeFindMetaDataEndpoint:")
			return nil, err
		}
		return FSResponse{Metadata: resp}, nil
	}
}

func MakeMkdirEndpoint(svc filesystem.FSService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(FSRequest)
		if !ok {
			log.Printf("failed to assert MakeMkdirEndpoint:")
			return nil, errors.New("failed to assert MakeMkdirEndpoint")
		}
		err := svc.Mkdir(ctx, req.Path)
		if err != nil {
			log.Printf("error: MakeMkdirEndpoint:")
			return nil, err
		}
		return FSResponse{}, nil
	}
}
func MakeRemoveEndpoint(svc filesystem.FSService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(FSRequest)
		if !ok {
			log.Printf("failed to assert MakeRemoveEndpoint:")
			return nil, errors.New("failed to assert MakeRemoveEndpoint")
		}
		err := svc.Remove(ctx, req.Path)
		if err != nil {
			log.Printf("error: MakeRemoveEndpoint:")
			return nil, err
		}
		return FSResponse{}, nil
	}
}
func MakeWriteFileEndpoint(svc filesystem.FSService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(FSRequest)
		if !ok {
			log.Printf("failed to assert MakeWriteFileEndpoint:")
			return nil, errors.New("failed to assert MakeWriteFileEndpoint")
		}
		err := svc.WriteFile(ctx, req.Path, req.Content, req.Flag)
		if err != nil {
			log.Printf("error: MakeWriteFileEndpoint:")
			return nil, err
		}
		return FSResponse{}, nil
	}
}
func MakeReadFileEndpoint(svc filesystem.FSService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(FSRequest)
		if !ok {
			log.Printf("failed to assert MakeReadFileEndpoint:")
			return nil, errors.New("failed to assert MakeReadFileEndpoint")
		}
		resp, err := svc.ReadFile(ctx, req.Path)
		if err != nil {
			log.Printf("error: MakeReadFileEndpoint:")
			return nil, err
		}
		return FSResponse{Result: resp}, nil
	}
}

func MakeCreateDynamicFileEndpoint(svc filesystem.FSService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(FSRequest)
		if !ok {
			log.Printf("failed to assert MakeCreateDynamicFileEndpoint:")
			return nil, errors.New("failed to assert MakeCreateDynamicFileEndpoint")
		}
		err := svc.CreateDynamicFile(ctx, req.Path, req.GeneratorType)
		if err != nil {
			log.Printf("error: MakeCreateDynamicFileEndpoint:")
			return nil, err
		}
		return FSResponse{}, nil
	}
}

func MakeListChildrenEndpoint(svc filesystem.FSService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(FSRequest)
		if !ok {
			log.Printf("failed to assert MakeListChildrenEndpoint:")
			return nil, errors.New("failed to assert MakeListChildrenEndpoint")
		}
		resp, err := svc.ListChildren(ctx, req.Path)
		if err != nil {
			log.Printf("error: MakeListChildrenEndpoint:")
			return nil, err
		}
		result := make([]filesystem.FileNodeInfo, 0, len(resp))
		for _, item := range resp {
			info := item.Stat()
			result = append(result, filesystem.FileNodeInfo{
				Name:    item.GetName(),
				IsDir:   item.IsDir(),
				Path:    item.GetPath(),
				Size:    info.Size,
				Mode:    int32(info.Mode),
				ModTime: info.ModTime.Unix(),
			})
		}
		return FSResponse{
			Children: result,
		}, nil
	}
}
