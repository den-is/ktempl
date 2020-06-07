package render

import (
	"fmt"
	"os"
	"strings"

	"github.com/den-is/ktempl/pkg/logging"
)

// Converts slice of strings of form "key=value" into map[key]value
func StringSliceToStringMap(s []string) map[string]string {

	result := make(map[string]string)

	if len(s) > 0 {
		for _, v := range s {
			kvSlice := strings.SplitN(v, "=", 2)
			if len(kvSlice) <= 1 {
				// TODO: Validate. do that validation using separate function, during early stages of the app initialization
				logging.LogWithFields(
					logging.Fields{
						"component": "render",
					}, "error", "Bad set value provided:", v)
				os.Exit(1)
			}
			result[kvSlice[0]] = kvSlice[1]
		}
	}

	return result

}

func StringifyStringsMap(m map[string]string) string {

	if m == nil {
		return ""
	}

	var s []string
	for k, v := range m {
		s = append(s, fmt.Sprintf("%s=%s", k, v))
	}

	return strings.Join(s, ",")

}
