package random

import (
	"math/rand"
	"time"
)

func NewRandomAlias(aliasLength int) string {
	res := make([]byte, aliasLength)
	arr := `QWERTYUIOPASDFGHJKLZXCVBNMqwertyuiopasdfghjklzxcvbnm[]{};'",.<>/\?!@#$%^&*()_+-=1234567890~`
	for i := 0; i < aliasLength; i++ {
		res = append(res, arr[rand.New(rand.NewSource(time.Now().UnixNano())).Int31n(int32(len(arr)))])
	}
	return string(res)
}
