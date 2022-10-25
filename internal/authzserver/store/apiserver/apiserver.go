// Copyright 2020 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package apiserver

import (
	"google.golang.org/grpc/credentials/insecure"
	"sync"

	pb "github.com/marmotedu/api/proto/apiserver/v1"
	"github.com/marmotedu/iam/internal/authzserver/store"
	"github.com/marmotedu/iam/pkg/log"
	"google.golang.org/grpc"
)

type datastore struct {
	cli pb.CacheClient
}

func (ds *datastore) Secrets() store.SecretStore {
	return newSecrets(ds)
}

func (ds *datastore) Policies() store.PolicyStore {
	return newPolicies(ds)
}

var (
	apiServerFactory store.Factory
	once             sync.Once
)

// GetAPIServerFactoryOrDie return cache instance and panics on any error.
func GetAPIServerFactoryOrDie(address string, clientCA string) store.Factory {
	once.Do(func() {
		var (
			err  error
			conn *grpc.ClientConn
			//creds credentials.TransportCredentials
		)

		//creds, err = credentials.NewClientTLSFromFile(clientCA, "")
		//if err != nil {
		//	log.Panicf("credentials.NewClientTLSFromFile err: %v", err)
		//}

		//conn, err = grpc.Dial(address, grpc.WithBlock(), grpc.WithTransportCredentials(creds))
		conn, err = grpc.Dial(address, grpc.WithBlock(), grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Panicf("Connect to grpc server failed, error: %s", err.Error())
		}

		apiServerFactory = &datastore{pb.NewCacheClient(conn)}
		log.Infof("Connected to grpc server, address: %s", address)
	})

	if apiServerFactory == nil {
		log.Panicf("failed to get apiserver store fatory")
	}

	return apiServerFactory
}
