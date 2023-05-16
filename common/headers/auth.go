package headers

import (
	"encoding/base64"
	"fmt"
)

func Authorization(token string) string {
	return fmt.Sprintf(
		"Bearer %s",
		base64.StdEncoding.EncodeToString([]byte(token)),
	)
}
