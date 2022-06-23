package config

import (
	"time"
)

const (
	RegUser   = "reg"
	UnregUser = "unreg"

	/*------------------------------------------------------*/

	BanDuration = time.Hour * 1 // will ban till

	/*------------------------------------------------------*/

	frequency         = time.Minute                   // some interval
	requests          = 100                           // accept requests in interval
	FrequencyRequest  = frequency / requests          // time for 1 request
	CheckAfter        = requests / 2                  // makes n-request and then will check for ddos (half)
	CheckTimeInterval = FrequencyRequest * CheckAfter // n-requests must be not often that this interval

	/*------------------------------------------------------*/

	SessionExpire        = time.Hour * 24 // 1 day (in days)
	SessionExpireByToken = time.Hour * 24 // 1 day (in days)
	SessionMaxAge        = 24 * 60 * 60   // 1 day (in seconds)

	/*------------------------------------------------------*/

	DeleteAfter     = time.Hour * 24 * 31 // 1 month (if sessions/banned_ip not updated more than month in map, it can delete from memory)
	NodeMaxLength   = 10000               // if linked list too big, it will find old records by "DeleteAfter"
	FindOldNodeLoop = NodeMaxLength / 3   // function will look for old keys a given number of times when run
	ClearDuration   = time.Second * 10    // break clearing if timeout
)
