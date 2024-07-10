package main

import (
	"context"
	"fmt"
	"log"

	"github.com/priyangshupal/grpc-raft-consensus/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GRPCTransport struct {
	ListenAddr string
}

func (t *GRPCTransport) Addr() string {
	return t.ListenAddr
}
func (t *GRPCTransport) Dial(s *RaftServer, addr string) error {
	conn, err := grpc.NewClient(
		addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return fmt.Errorf("could not connect to raft replica: %v", err)
	}
	bootstrapServiceClient := pb.NewBootstrapServiceClient(conn)
	response, err := bootstrapServiceClient.AddReplica(
		context.Background(),
		&pb.AddrInfo{Addr: t.Addr()},
	)
	if err != nil {
		return fmt.Errorf("[%s] error while calling bootstapping service: %v", t.Addr(), err)
	}
	if !response.IsAdded {
		return fmt.Errorf("error while adding replica [%s] to [%s]: %v", t.Addr(), addr, err)
	}
	log.Printf("[%s] successfully connected to [%s]\n", s.Transport.Addr(), addr)
	// everything is fine, so add the connection to replicaConnMap
	s.ReplicaConnMapLock.Lock()
	s.ReplicaConnMap[addr] = conn
	s.ReplicaConnMapLock.Unlock()
	return nil
}
