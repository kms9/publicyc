package debug

import (
	"encoding/json"
	"fmt"
)

func Print(s interface{}) {
	b, err := json.Marshal(s)
	fmt.Println("Debug: ", string(b), " err:", err)
}
