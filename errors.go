package hipache

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
