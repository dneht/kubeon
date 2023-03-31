/*
Copyright (c) 2020, Dash

Licensed under the LGPL, Version 3.0 (the "License");
you may not use this file except in compliance with the License.
*/

package onutil

import (
	"encoding/json"
)

func PrettyJson(jsonData interface{}) ([]byte, error) {
	return json.MarshalIndent(jsonData, "", "  ")
}
