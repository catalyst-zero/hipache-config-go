package hipache

import (
	"reflect"
	"testing"
)

const (
	TEST_DOMAIN = "hipache.test.example.tld"
	HOST        = "127.0.0.1:6379"
)

func TestRedisHipacheConfigClient(t *testing.T) {
	client, err := DialHipacheConfig(HOST)
	if err != nil {
		t.Skipf("Failed to connect to local redis server, skipping.")
	}

	// Create
	client.BindingCreate(TEST_DOMAIN)
	client.BindingAddHost(TEST_DOMAIN, "192.168.2.100")
	client.BindingAddHost(TEST_DOMAIN, "192.168.2.101")
	client.BindingAddHost(TEST_DOMAIN, "192.168.2.102")

	dump, _ := client.BindingGet(TEST_DOMAIN)
	expected := Binding{DomainName: TEST_DOMAIN, Hosts: []string{"192.168.2.100", "192.168.2.101", "192.168.2.102"}}

	if !reflect.DeepEqual(dump, expected) {
		t.Fatalf("Expected %v to contain three hosts.", dump)
	}

	// Test delete
	client.BindingDelete(TEST_DOMAIN)
	dump, err = client.BindingGet(TEST_DOMAIN)

	if err == nil {
		t.Fatalf("Expected error for deleted binding!")
	}
	if dump.DomainName != TEST_DOMAIN {
		t.Fatalf("Received wrong domain name: %s", dump.DomainName)
	}
	if len(dump.Hosts) != 0 {
		t.Fatalf("Expected no hosts in response, but has %d hosts", len(dump.Hosts))
	}

}

func TestRedisHipacheConfigClient_BindingCreate__returns_error_if_already_exists(t *testing.T) {
	client, err := DialHipacheConfig(HOST)
	if err != nil {
		t.Skipf("Failed to connect to local redis server, skipping.")
	}

	client.BindingCreate(TEST_DOMAIN)

	err = client.BindingCreate(TEST_DOMAIN)
	switch err.(type) {
	case *BindingAlreadyExistsError:
		return
	default:
		t.Fatalf("Unexpected error for second BindingCreate, expected BindingAlreadyExistsError: %v", err)
	}
}

func TestLocalHipacheConfig__UseCase1(t *testing.T) {
	config := NewLocalHipacheConfig()

	config.BindingCreate(TEST_DOMAIN)
	err := config.BindingCreate(TEST_DOMAIN)

	switch err.(type) {
	case *BindingAlreadyExistsError:
		goto create_ok
	default:
		t.Fatalf("Second BindingCreate() did not raise expected BindingAlreadyExistsError()")
	}

create_ok:

	binding, err := config.BindingGet(TEST_DOMAIN)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if binding.DomainName != TEST_DOMAIN {
		t.Fatalf("BindingGet() returned wrong domain name: %s", binding.DomainName)
	}
	if len(binding.Hosts) != 0 {
		t.Fatalf("BindingGet() returned hostnames, expected none: %v", binding.Hosts)
	}

	// Add host 1
	if err := config.BindingAddHost(TEST_DOMAIN, "10.17.0.100"); err != nil {
		t.Fatalf("BindingAddHost() failed")
	}

	binding, err = config.BindingGet(TEST_DOMAIN)
	if len(binding.Hosts) != 1 && binding.Hosts[0] != "10.17.0.100" {
		t.Fatalf("BindingGet() returned unexpected host: %v", binding.Hosts)
	}

	// Add Host 2
	if err := config.BindingAddHost(TEST_DOMAIN, "10.17.0.101"); err != nil {
		t.Fatalf("BindingAddHost() failed")
	}

	binding, err = config.BindingGet(TEST_DOMAIN)
	if len(binding.Hosts) != 1 && binding.Hosts[0] != "10.17.0.100" && binding.Hosts[1] != "10.17.0.101" {
		t.Fatalf("BindingGet() returned unexpected hosts: %v", binding.Hosts)
	}

}

func TestLocalHipacheConfig_BindingGet__non_existant(t *testing.T) {
	config := NewLocalHipacheConfig()

	_, err := config.BindingGet(TEST_DOMAIN)
	switch err.(type) {
	case *BindingNotFoundError:
		return
	default:
		t.Fatalf("Expected BindingGet() to return BindingNotFoundError")
	}
}
