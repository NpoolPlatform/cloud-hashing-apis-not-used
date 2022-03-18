package referral

import (
	"fmt"
)

func CacheKey(appID, userID, usedFor string) string {
	return fmt.Sprintf("%v:%v:%v", appID, userID, usedFor)
}
