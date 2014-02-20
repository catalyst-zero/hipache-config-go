package server

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
)

// DialHipacheConfig() established a new connection to the given host:port
// and returns a new HipacheConfigClient to modify the target redis store.
//
// `connectionString` must be in "host:port" format
func DialHipacheConfig(connectionString string) (HipacheConfig, error) {
	redis, err := redis.Dial("tcp", connectionString)
	return &redisHipacheConfigClient{redis}, err
}

type redisHipacheConfigClient struct {
	// TODO: Should we use a pool for this?
	redis redis.Conn
}

func hipacheFrontendKey(domainName string) string {
	return "frontend:" + domainName
}

func (client *redisHipacheConfigClient) BindingCreate(domainName string) error {
	reply, err := client.redis.Do("EXISTS", hipacheFrontendKey(domainName))
	if err != nil {
		return err
	}
	if reply.(int64) == 1 {
		return &BindingAlreadyExistsError{domainName}
	}

	reply, err = client.redis.Do("RPUSH", hipacheFrontendKey(domainName), domainName)
	if err != nil {
		return err
	}
	if reply.(int64) != 1 {
		return fmt.Errorf("Unexpected list-length after initial rpush: %d", reply)
	}
	return nil
}

func (client *redisHipacheConfigClient) BindingDelete(domainName string) error {
	_, err := client.redis.Do("DEL", hipacheFrontendKey(domainName))
	if err != nil {
		return err
	}
	// we don't check the results, since all non-panic cases are fine
	// if <0 => Panic?
	// if 0 => Key didn't exist
	// if 1 => Ok, key is gone
	// if >1 => panic
	return nil
}

func (client *redisHipacheConfigClient) BindingAddHost(domainName, backendHostAddress string) error {
	_, err := client.redis.Do("RPUSH", hipacheFrontendKey(domainName), backendHostAddress)
	if err != nil {
		return err
	}
	return nil
}

// BindingRemoveHost() removes the given backendHostAddress from the list of backends for domainName.
// LREM frontend:{domainName} -1 {backendHostAddress}
func (client *redisHipacheConfigClient) BindingRemoveHost(domainName, backendHostAddress string) error {
	_, err := client.redis.Do("LREM", hipacheFrontendKey(domainName), "1", backendHostAddress)
	if err != nil {
		return err
	}
	return nil
}

func (client *redisHipacheConfigClient) BindingGet(domainName string) (binding Binding, err error) {
	binding.DomainName = domainName

	hosts, err := redis.Strings(client.redis.Do("LRANGE", hipacheFrontendKey(domainName), 0, -1))
	if err != nil {
		return binding, err
	}
	if len(hosts) == 0 {
		return binding, &BindingNotFoundError{domainName}
	}

	binding.Hosts = hosts[1:]

	return binding, nil
}
