package render

import (
	"os"
	"strings"

	"github.com/den-is/ktempl/pkg/logging"
)

// Converts slice of strings of form "key=value" into map[key]value
func StringSliceToStringMap(s []string) map[string]string {

	result := make(map[string]string)

	if len(s) > 0 {
		for _, v := range s {
			kv_s := strings.SplitN(v, "=", 2)
			if len(kv_s) <= 1 {
				// TODO: Validate. do that validation using separate function, during early stages of the app initialization
				logging.LogWithFields(
					logging.Fields{
						"component": "render",
					}, "error", "Bad set value provided:", v)
				os.Exit(1)
			}
			result[kv_s[0]] = kv_s[1]
		}
	}

	return result

}
