package fs_transport

import (
	"context"
	"errors"
	"log"

	"github.com/BBitQNull/SSHoneyNet/core/filesystem"
	fs_endpoint "github.com/BBitQNull/SSHoneyNet/modules/filesystem/endpoint"
	fs_Pb "github.com/BBitQNull/SSHoneyNet/pb/filesystem"
	"github.com/BBitQNull/SSHoneyNet/pkg/utils/convert"
	grpctransport "github.com/go-kit/kit/transport/grpc"
)

type grpcServer struct {
	fs_Pb.UnimplementedFileManageServer
	createFile        grpctransport.Handler
	findMetaData      grpctransport.Handler
	createDynamicFile grpctransport.Handler
	mkdir             grpctransport.Handler
	remove            grpctransport.Handler
	writeFile         grpctransport.Handler
	readFile          grpctransport.Handler
	listChildren      grpctransport.Handler
}

func decodeFSRequest(ctx context.Context, grpcReq interface{}) (interface{}, error) {
	req, ok := grpcReq.(*fs_Pb.FileRequest)
	if !ok {
		log.Printf("failed to assert *fs_Pb.FileRequest")
		return nil, errors.New("failed to assert *fs_Pb.FileRequest")
	}
	return fs_endpoint.FSRequest{
		Path:          req.Path,
		Content:       req.Content,
		Flag:          req.Flag,
		GeneratorType: req.GeneratorType,
	}, nil
}

func encodeFSResponse(ctx context.Context, response interface{}) (interface{}, error) {
	resp, ok := response.(fs_endpoint.FSResponse)
	if !ok {
		log.Printf("failed to assert fs_endpoint.FSResponse")
		return nil, errors.New("failed to assert fs_endpoint.FSResponse")
	}
	return &fs_Pb.FileResponse{
		Result:   resp.Result,
		Metadata: convert.ConvertMetadataToPb(resp.Metadata),
		Children: convert.ConvertChildrenToPb(resp.Children),
	}, nil
}

func NewFSServer(svc filesystem.FSService) fs_Pb.FileManageServer {
	return &grpcServer{
		createFile: grpctransport.NewServer(
			fs_endpoint.MakeCreateFileEndpoint(svc),
			decodeFSRequest,
			encodeFSResponse,
		),
		findMetaData: grpctransport.NewServer(
			fs_endpoint.MakeFindMetaDataEndpoint(svc),
			decodeFSRequest,
			encodeFSResponse,
		),
		createDynamicFile: grpctransport.NewServer(
			fs_endpoint.MakeCreateDynamicFileEndpoint(svc),
			decodeFSRequest,
			encodeFSResponse,
		),
		mkdir: grpctransport.NewServer(
			fs_endpoint.MakeMkdirEndpoint(svc),
			decodeFSRequest,
			encodeFSResponse,
		),
		writeFile: grpctransport.NewServer(
			fs_endpoint.MakeWriteFileEndpoint(svc),
			decodeFSRequest,
			encodeFSResponse,
		),
		remove: grpctransport.NewServer(
			fs_endpoint.MakeRemoveEndpoint(svc),
			decodeFSRequest,
			encodeFSResponse,
		),
		readFile: grpctransport.NewServer(
			fs_endpoint.MakeReadFileEndpoint(svc),
			decodeFSRequest,
			encodeFSResponse,
		),
		listChildren: grpctransport.NewServer(
			fs_endpoint.MakeListChildrenEndpoint(svc),
			decodeFSRequest,
			encodeFSResponse,
		),
	}
}

func (s *grpcServer) CreateFile(ctx context.Context, req *fs_Pb.FileRequest) (*fs_Pb.FileResponse, error) {
	_, rep, err := s.createFile.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*fs_Pb.FileResponse), nil
}
func (s *grpcServer) FindMetaData(ctx context.Context, req *fs_Pb.FileRequest) (*fs_Pb.FileResponse, error) {
	_, rep, err := s.findMetaData.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*fs_Pb.FileResponse), nil
}
func (s *grpcServer) CreateDynamicFile(ctx context.Context, req *fs_Pb.FileRequest) (*fs_Pb.FileResponse, error) {
	_, rep, err := s.createDynamicFile.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*fs_Pb.FileResponse), nil
}
func (s *grpcServer) Mkdir(ctx context.Context, req *fs_Pb.FileRequest) (*fs_Pb.FileResponse, error) {
	_, rep, err := s.mkdir.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*fs_Pb.FileResponse), nil
}
func (s *grpcServer) Remove(ctx context.Context, req *fs_Pb.FileRequest) (*fs_Pb.FileResponse, error) {
	_, rep, err := s.remove.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*fs_Pb.FileResponse), nil
}
func (s *grpcServer) WriteFile(ctx context.Context, req *fs_Pb.FileRequest) (*fs_Pb.FileResponse, error) {
	_, rep, err := s.writeFile.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*fs_Pb.FileResponse), nil
}
func (s *grpcServer) ReadFile(ctx context.Context, req *fs_Pb.FileRequest) (*fs_Pb.FileResponse, error) {
	_, rep, err := s.readFile.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*fs_Pb.FileResponse), nil
}
func (s *grpcServer) ListChildren(ctx context.Context, req *fs_Pb.FileRequest) (*fs_Pb.FileResponse, error) {
	_, rep, err := s.listChildren.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*fs_Pb.FileResponse), nil
}
