package v1

import (
	"errors"
	"fmt"
)

// Deprecated: Get is deprecated and will be removed in a future version.
// Use GetNew instead.
func Get() (*Resp, error) {
	fmt.Println("implement me")
	return nil, errors.New("implement me")
}

func GetRAW() []byte {
	panic("implement me")
}
