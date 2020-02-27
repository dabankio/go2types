package user

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestJSON(t *testing.T) {
	b, _ := json.MarshalIndent(Person{
		Age: 17,
	}, "", "  ")
	fmt.Println(string(b))
}
