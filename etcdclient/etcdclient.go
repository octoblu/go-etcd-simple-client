package etcdclient

import (
	"github.com/coreos/etcd/Godeps/_workspace/src/golang.org/x/net/context"
	"github.com/coreos/etcd/client"
)

// EtcdClient interface lets your Get/Set from Etcd
type EtcdClient interface {
	Del(key string) error
	Get(key string) (string, error)
	Set(key, value string) error
	Ls(directory string) ([]string, error)
	LsRecursive(directory string) ([]string, error)
}

// SimpleEtcdClient implements EtcdClient
type SimpleEtcdClient struct {
	etcd client.Client
}

// Dial constructs a new EtcdClient
func Dial(etcdURI string) (EtcdClient, error) {
	etcd, err := client.New(client.Config{
		Endpoints: []string{etcdURI},
	})
	if err != nil {
		return nil, err
	}
	return &SimpleEtcdClient{etcd}, nil
}

// Del deletes a key from Etcd
func (etcdClient *SimpleEtcdClient) Del(key string) error {
	api := client.NewKeysAPI(etcdClient.etcd)
	_, err := api.Delete(context.Background(), key, nil)
	return err
}

// Get gets a value in Etcd
func (etcdClient *SimpleEtcdClient) Get(key string) (string, error) {
	api := client.NewKeysAPI(etcdClient.etcd)
	response, err := api.Get(context.Background(), key, nil)
	if err != nil {
		if client.IsKeyNotFound(err) {
			return "", nil
		}
		return "", err
	}
	return response.Node.Value, nil
}

// Set sets a value in Etcd
func (etcdClient *SimpleEtcdClient) Set(key, value string) error {
	api := client.NewKeysAPI(etcdClient.etcd)
	_, err := api.Set(context.Background(), key, value, nil)
	return err
}

// Ls returns all the keys available in the directory
func (etcdClient *SimpleEtcdClient) Ls(directory string) ([]string, error) {
	api := client.NewKeysAPI(etcdClient.etcd)
	options := &client.GetOptions{Sort: true, Recursive: false}
	response, err := api.Get(context.Background(), directory, options)

	if err != nil {
		if client.IsKeyNotFound(err) {
			return make([]string, 0), nil
		}
		return make([]string, 0), err
	}

	return nodesToStringSlice(response.Node.Nodes), nil
}

// LsRecursive returns all the keys available in the directory, recursively
func (etcdClient *SimpleEtcdClient) LsRecursive(directory string) ([]string, error) {
	api := client.NewKeysAPI(etcdClient.etcd)
	options := &client.GetOptions{Sort: true, Recursive: true}
	response, err := api.Get(context.Background(), directory, options)

	if err != nil {
		if client.IsKeyNotFound(err) {
			return make([]string, 0), nil
		}
		return make([]string, 0), err
	}

	return nodesToStringSlice(response.Node.Nodes), nil
}

func nodesToStringSlice(nodes client.Nodes) []string {
	var keys []string

	for _, node := range nodes {
		keys = append(keys, node.Key)

		for _, key := range nodesToStringSlice(node.Nodes) {
			keys = append(keys, key)
		}
	}

	return keys
}
