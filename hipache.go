package server

import (
	"fmt"
)

// # Hipache Redis Guide from README
// ## Define backend
// redis-cli rpush frontend:www.dotcloud.com mywebsite
// (integer) 1
// ## Associate two hosts
// redis-cli rpush frontend:www.dotcloud.com http://192.168.0.42:80
// (integer) 2
// redis-cli rpush frontend:www.dotcloud.com http://192.168.0.43:80
// (integer) 3
//
// ## Review configuration
// $ redis-cli lrange frontend:www.dotcloud.com 0 -1
// 1) "mywebsite"
// 2) "http://192.168.0.42:80"
// 3) "http://192.168.0.43:80"
type HipacheConfig interface {
	// BindingCreate() creates a new binding after first checking, it none with the name already exists.
	// If a binding exists, error is of type *BindingAlreadyExistsError
	BindingCreate(domainName string) error

	BindingDelete(domainName string) error

	BindingAddHost(domainName, backendHostAddress string) error

	// BindingRemoveHost() removes the given backendHostAddress from the list of backends for domainName.
	// LREM frontend:{domainName} -1 {backendHostAddress}
	BindingRemoveHost(domainName, backendHostAddress string) error

	// BindingGet() reads all backends and returns them.
	// If no binding with the given name exists, error is of type *BindingNotFoundError, otherwise nil.
	BindingGet(domainName string) (binding Binding, err error)
}

type Binding struct {
	DomainName string
	Hosts      []string
}

type BindingNotFoundError struct {
	DomainName string
}

func (err *BindingNotFoundError) Error() string {
	return fmt.Sprintf("Binding not found: %s", err.DomainName)
}

type BindingAlreadyExistsError struct {
	DomainName string
}

func (err *BindingAlreadyExistsError) Error() string {
	return fmt.Sprintf("Binding already exists: %s", err.DomainName)
}
