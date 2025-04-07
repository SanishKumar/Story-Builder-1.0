package main

import (
	"context"
	"log"
	"net"
	"sync"

	pb "github.com/SanishKumar/story-builder/proto/story"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

// StoryServer implements the StoryService.
type StoryServer struct {
	pb.UnimplementedStoryServiceServer
	mu       sync.Mutex
	segments []*pb.StorySegment
	clients  map[pb.StoryService_StreamStoryUpdatesServer]struct{}
}

// NewStoryServer initializes our server.
func NewStoryServer() *StoryServer {
	return &StoryServer{
		segments: make([]*pb.StorySegment, 0),
		clients:  make(map[pb.StoryService_StreamStoryUpdatesServer]struct{}),
	}
}

// GenerateSegment calls the free generative AI API to generate text.
func (s *StoryServer) GenerateSegment(ctx context.Context, req *pb.GenerateRequest) (*pb.GenerateResponse, error) {
	// Here, integrate your free generative AI API call.
	// For demo, we'll use a placeholder response.
	generatedText := "Once upon a time, in a land of endless possibilities..."
	return &pb.GenerateResponse{GeneratedText: generatedText}, nil
}

// StreamStoryUpdates streams current story segments to connected clients.
func (s *StoryServer) StreamStoryUpdates(empty *emptypb.Empty, stream pb.StoryService_StreamStoryUpdatesServer) error {
	s.mu.Lock()
	s.clients[stream] = struct{}{}
	// Send initial state
	err := stream.Send(&pb.UpdateResponse{Segments: s.segments})
	s.mu.Unlock()
	if err != nil {
		return err
	}
	// Keep the stream open.
	<-stream.Context().Done()
	s.mu.Lock()
	delete(s.clients, stream)
	s.mu.Unlock()
	return nil
}

// SubmitSegment adds a new story segment and broadcasts it.
func (s *StoryServer) SubmitSegment(ctx context.Context, req *pb.UpdateRequest) (*emptypb.Empty, error) {
	s.mu.Lock()
	// Append the new segment.
	s.segments = append(s.segments, req.GetSegment())
	// Broadcast the update to all clients.
	for client := range s.clients {
		// We send updates in a goroutine so one slow client doesn't block others.
		go func(c pb.StoryService_StreamStoryUpdatesServer) {
			c.Send(&pb.UpdateResponse{Segments: s.segments})
		}(client)
	}
	s.mu.Unlock()
	return &emptypb.Empty{}, nil
}

func main() {
	port := ":50051"
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to listen on %v: %v", port, err)
	}

	grpcServer := grpc.NewServer()
	storyServer := NewStoryServer()
	pb.RegisterStoryServiceServer(grpcServer, storyServer)

	log.Printf("Server listening on %v", port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
