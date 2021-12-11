/*
 * Copyright 2021 LiAuto authors.
 * @yangchunxin
 */

package utils

import (
	"encoding/json"
	"fmt"
)

// PrettyPrint xxx
func PrettyPrint(v interface{}) {
	out, _ := json.MarshalIndent(v, "", "  ")
	fmt.Printf("%s\n", out)
}
