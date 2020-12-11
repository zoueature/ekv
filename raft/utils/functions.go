/*+-----------------------------+
 *| Author: Zoueature           |
 *+-----------------------------+
 *| Email: zoueature@gmail.com  |
 *+-----------------------------+
 */
package utils

import (
	"math/rand"
	"time"
)

func RandTimer(min, max time.Duration) (*time.Timer, time.Duration) {
	random := rand.New(rand.NewSource(time.Now().UnixNano()))
	randomNum := random.Intn(int(max - min))
	timerNum := time.Duration(randomNum) + min
	//150-300毫秒的随机定时器
	timer := time.NewTimer(timerNum)
	return timer, timerNum
}

func BlockRandTime(min, max time.Duration) time.Duration {
	timer, num := RandTimer(min, max)
	<-timer.C
	return num
}

func BlockStaticTime(num time.Duration) {
	timer := time.NewTimer(num)
	<-timer.C
}
