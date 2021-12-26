package pb

import (
	context "context"
	"grpc-server/pkg/storage"

	"github.com/sirupsen/logrus"
)

type StorageService struct {
	logger  *logrus.Logger
	storage *storage.InMemoryKeyValueStorage
	UnimplementedKVStorageServiceServer
}

func NewStorageService(logger *logrus.Logger) *StorageService {
	s := storage.NewInMemoryKeyValueStorage(300)

	return &StorageService{
		storage: s,
		logger:  logger,
	}
}

func (s *StorageService) Put(ctx context.Context, in *PutRequest) (*Item, error) {
	s.logger.Printf("Put item: %v\n", in)

	s.storage.Put(in.Key, in.Value)

	return &Item{Key: in.Key, Value: in.Value}, nil
}

func (s *StorageService) Get(ctx context.Context, in *GetRequest) (*Item, error) {
	s.logger.Printf("Get item: %v\n", in)

	value, exists := s.storage.Get(in.Key)

	if !exists {
		return &Item{}, nil
	}

	return &Item{Key: in.Key, Value: value}, nil
}

func (s *StorageService) Delete(ctx context.Context, in *DeleteRequest) (*DeleteResponse, error) {
	s.logger.Printf("Delete item: %v\n", in)

	exists := s.storage.Delete(in.Key)

	if !exists {
		return &DeleteResponse{Status: DeletedStatus_KEY_NOT_FOUND}, nil
	}

	return &DeleteResponse{Status: DeletedStatus_DELETED}, nil
}
