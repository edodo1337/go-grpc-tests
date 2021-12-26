package main

import (
	"grpc-client/pb"
	"log"
	"math/rand"
	"strconv"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

const (
	CLIENTS_COUNT = 5000
	KEY_RANGE     = 150
)

var responseCount int32

type Client struct {
	pbclient *pb.KVStorageServiceClient
	t        *testing.T
}

func NewClient(addr string, t *testing.T) *Client {
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("New client creating %v", err.Error())
	}

	client := pb.NewKVStorageServiceClient(conn)

	return &Client{
		pbclient: &client,
		t:        t,
	}
}

func (c *Client) PutRandom(ctx context.Context) error {
	key := rand.Intn(KEY_RANGE)

	strKey := strconv.Itoa(key)

	item := pb.PutRequest{Key: strKey, Value: strKey}

	_, err := (*c.pbclient).Put(ctx, &item)
	if err != nil {
		log.Printf("Put random err %v", err)

		return err
	}

	return nil
}

func (c *Client) GetRandom(ctx context.Context) error {
	key := rand.Intn(KEY_RANGE)

	strKey := strconv.Itoa(key)

	item := pb.GetRequest{Key: strKey}

	_, err := (*c.pbclient).Get(ctx, &item)
	if err != nil {
		log.Printf("Get random err %v", err)

		return err
	}

	return nil
}

func (c *Client) DeleteRandom(ctx context.Context) error {
	key := rand.Intn(KEY_RANGE)

	strKey := strconv.Itoa(key)

	item := pb.DeleteRequest{Key: strKey}

	_, err := (*c.pbclient).Delete(ctx, &item)
	if err != nil {
		log.Printf("Delete random err %v", err)

		return err
	}

	return nil
}

type OpFunc func(ctx context.Context) error

func (c *Client) RandomOperation(wg *sync.WaitGroup) {
	defer wg.Done()

	OPS_MAP := map[int]OpFunc{
		0: c.PutRandom,
		1: c.GetRandom,
		2: c.DeleteRandom,
	}

	code := rand.Intn(3)

	op, ok := OPS_MAP[code]

	if !ok {
		c.t.Log("Not ok")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()

	if err := op(ctx); err == nil {
		atomic.AddInt32(&responseCount, 1)
	}
}

func TestLoad(t *testing.T) {
	// addr := flag.String("addr", "localhost:9000", "the address to connect to")
	addr := "localhost:9000"

	var clients []*Client
	var wg sync.WaitGroup

	for i := 0; i < CLIENTS_COUNT; i++ {
		c := NewClient(addr, t)
		clients = append(clients, c)
	}

	wg.Add(CLIENTS_COUNT)
	for i := 0; i < CLIENTS_COUNT; i++ {
		go clients[i].RandomOperation(&wg)
		// time.Sleep(1 * time.Millisecond)
	}

	wg.Wait()

	t.Error("?Hello", responseCount)
}
