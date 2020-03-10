package cryptomood

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	cache "tohan.net/go-practice/src/cache"
	cacheTypes "tohan.net/go-practice/src/cache/types"
)

func ConsumeSentiments(c *cache.Cache, certFile string, server string) {

	creds, err := credentials.NewClientTLSFromFile(certFile, "")
	if err != nil {
		panic(err)
	}
	conn, err := grpc.Dial(server, grpc.WithTransportCredentials(creds), grpc.WithTimeout(5*time.Second), grpc.WithBlock())
	if err != nil {
		panic(fmt.Sprintf("did not connect: %v", err))
	}
	fmt.Println("Connected to cryptomood")

	proxyClient := NewSentimentsClient(conn)

	req := &AggregationCandleFilter{Resolution: "M1", AssetsFilter: &AssetsFilter{Assets: []string{"BTC", "ETH"}, AllAssets: false}}
	sub, err := proxyClient.SubscribeSocialSentiment(context.Background(), req)
	if err != nil {
		panic(err)
	}
	for {
		msg, err := sub.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			panic(err)
		}
		out, err := json.Marshal(msg.Id)
		if err != nil {
			fmt.Println("Sentiment is in wrong format. Cannot process.")
			continue
		}
		c.AddItem(cacheTypes.CacheItem{Key: string(out), Value: msg.Asset})
	}
}
