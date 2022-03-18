package referral

import (
	"fmt"
)

func cacheKey(appID, userID, usedFor string) string {
	return fmt.Sprintf("%v:%v:%v", appID, userID, usedFor)
}
