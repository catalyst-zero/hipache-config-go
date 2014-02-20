package server

func NewLocalHipacheConfig() HipacheConfig {
	return &localHipacheConfig{
		make(map[string]Binding),
	}
}

type localHipacheConfig struct {
	bindings map[string]Binding
}

func (cfg *localHipacheConfig) BindingCreate(domainName string) error {
	_, ok := cfg.bindings[domainName]
	if ok {
		return &BindingAlreadyExistsError{domainName}
	}

	cfg.bindings[domainName] = Binding{DomainName: domainName}
	return nil
}

func (cfg *localHipacheConfig) BindingDelete(domainName string) error {
	delete(cfg.bindings, domainName)
	return nil
}

func (cfg *localHipacheConfig) BindingGet(domainName string) (Binding, error) {
	value, ok := cfg.bindings[domainName]
	if !ok {
		return value, &BindingNotFoundError{domainName}
	}
	return value, nil
}

func (cfg *localHipacheConfig) BindingAddHost(domainName, backendHostAddress string) error {
	value, ok := cfg.bindings[domainName]
	if !ok {
		return &BindingNotFoundError{domainName}
	}

	value.Hosts = append(value.Hosts, backendHostAddress)
	cfg.bindings[domainName] = value
	return nil
}

func (cfg *localHipacheConfig) BindingRemoveHost(domainName, backendHostAddress string) error {
	value, ok := cfg.bindings[domainName]
	if !ok {
		return &BindingNotFoundError{domainName}
	}

	for i, host := range value.Hosts {
		if host == backendHostAddress {
			copy(value.Hosts[i:], value.Hosts[i+1:])
			cfg.bindings[domainName] = value
			return nil
		}
	}
	return nil
}
