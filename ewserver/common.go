package ewserver

// PathName represents a path of something being accessed
type PathName string

// UserName represents an entity accessing, or acting on something
type UserName string

func (u UserName) String() string {
	return string(u)
}

// Bytes returns the username as a byte slice.
func (u UserName) Bytes() []byte {
	return []byte(u)
}

// APIKey represents an API Key
type APIKey string

// Bytes returns the username as as byte slice.
func (k APIKey) Bytes() []byte {
	return []byte(k)
}
