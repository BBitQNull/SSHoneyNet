package fs_client

import (
	"context"
	"fmt"

	"github.com/BBitQNull/SSHoneyNet/core/filesystem"
	fs_Pb "github.com/BBitQNull/SSHoneyNet/pb/filesystem"
	"github.com/BBitQNull/SSHoneyNet/pkg/utils/convert"
	"github.com/go-kit/kit/endpoint"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	"google.golang.org/grpc"
)

type RawFSRequest struct {
	Path           string
	Content        []byte
	Flag           string
	Generator_type string
}

type RawFSResponse struct {
	Result   []byte
	Metadata filesystem.FileInfo
}

type FSManageClient struct {
	CreateFile        endpoint.Endpoint
	FindMetaData      endpoint.Endpoint
	CreateDynamicFile endpoint.Endpoint
	Mkdir             endpoint.Endpoint
	Remove            endpoint.Endpoint
	WriteFile         endpoint.Endpoint
	ReadFile          endpoint.Endpoint
}

func NewFSManageClient(conn *grpc.ClientConn) FSManageClient {
	return FSManageClient{
		CreateFile:        MakeCreateFileEndpoint(conn),
		FindMetaData:      MakeFindMetaDataEndpoint(conn),
		CreateDynamicFile: MakeCreateDynamicFileEndpoint(conn),
		Mkdir:             MakeMkdirEndpoint(conn),
		Remove:            MakeRemoveEndpoint(conn),
		WriteFile:         MakeWriteFileEndpoint(conn),
		ReadFile:          MakeReadFileEndpoint(conn),
	}
}

func encodeFileRequest(_ context.Context, request interface{}) (interface{}, error) {
	req, ok := request.(*RawFSRequest)
	if !ok {
		return nil, fmt.Errorf("encodeFileRequest: expected *RawFSRequest but got %T", request)
	}
	return &fs_Pb.FileRequest{
		Path:          req.Path,
		Content:       req.Content,
		Flag:          req.Flag,
		GeneratorType: req.Generator_type,
	}, nil
}

func decodeFileResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp, ok := response.(*fs_Pb.FileResponse)
	if !ok {
		return nil, fmt.Errorf("decodeFileResponse: expected *ProcResponse but got %T", response)

	}
	return &RawFSResponse{
		Result:   resp.Result,
		Metadata: convert.ConvertMetadataFromPb(resp.Metadata),
	}, nil
}

func MakeCreateFileEndpoint(conn *grpc.ClientConn) endpoint.Endpoint {
	return grpctransport.NewClient(
		conn,
		"pb.FileManage",
		"CreateFile",
		encodeFileRequest,
		decodeFileResponse,
		fs_Pb.FileResponse{},
	).Endpoint()
}
func MakeCreateDynamicFileEndpoint(conn *grpc.ClientConn) endpoint.Endpoint {
	return grpctransport.NewClient(
		conn,
		"pb.FileManage",
		"CreateDynamicFile",
		encodeFileRequest,
		decodeFileResponse,
		fs_Pb.FileResponse{},
	).Endpoint()
}
func MakeMkdirEndpoint(conn *grpc.ClientConn) endpoint.Endpoint {
	return grpctransport.NewClient(
		conn,
		"pb.FileManage",
		"Mkdir",
		encodeFileRequest,
		decodeFileResponse,
		fs_Pb.FileResponse{},
	).Endpoint()
}
func MakeRemoveEndpoint(conn *grpc.ClientConn) endpoint.Endpoint {
	return grpctransport.NewClient(
		conn,
		"pb.FileManage",
		"Remove",
		encodeFileRequest,
		decodeFileResponse,
		fs_Pb.FileResponse{},
	).Endpoint()
}
func MakeWriteFileEndpoint(conn *grpc.ClientConn) endpoint.Endpoint {
	return grpctransport.NewClient(
		conn,
		"pb.FileManage",
		"WriteFile",
		encodeFileRequest,
		decodeFileResponse,
		fs_Pb.FileResponse{},
	).Endpoint()
}
func MakeReadFileEndpoint(conn *grpc.ClientConn) endpoint.Endpoint {
	return grpctransport.NewClient(
		conn,
		"pb.FileManage",
		"ReadFile",
		encodeFileRequest,
		decodeFileResponse,
		fs_Pb.FileResponse{},
	).Endpoint()
}
func MakeFindMetaDataEndpoint(conn *grpc.ClientConn) endpoint.Endpoint {
	return grpctransport.NewClient(
		conn,
		"pb.FileManage",
		"FindMetaData",
		encodeFileRequest,
		decodeFileResponse,
		fs_Pb.FileResponse{},
	).Endpoint()
}
