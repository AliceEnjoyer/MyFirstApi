package random

import (
	"math/rand"
	"time"
)

func NewRandomAlias(aliasLength int) string {
	var res []byte
	arr := `QWERTYUIOPASDFGHJKLZXCVBNMqwertyuiopasdfghjklzxcvbnm[]{};'",.<>/\?!@#$%^&*()_+-=1234567890~`
	for i := 0; i < aliasLength; i++ {
		res = append(res, arr[rand.New(rand.NewSource(time.Now().UnixNano())).Int31n(int32(len(arr)-1))])
	}
	return string(res)
}
