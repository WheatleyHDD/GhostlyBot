package commands

import (
	"github.com/SevereCloud/vksdk/api"
	"github.com/buger/jsonparser"
	scribble "github.com/nanobox-io/golang-scribble"
	"math/rand"
	"strconv"
	"strings"
	"time"
	. "vkbot/utils"
)

func Roll(alist []string, userID int, db *scribble.Driver, vk *api.VK) (otvet string, attachments string, similar float64) {
	config := ConfigParse()

	StoreStatistic("рулеточка", db)

	if len(alist) == 2 {
		i, err := strconv.Atoi(alist[1])
		if err != nil {
			answer, _ := jsonparser.GetString(config, "message_list", "roll", "notNumber")
			answer = GetRegularData(answer, userID, vk)
			return answer, "", 0
		}
		if i > 0 {
			rand.Seed(time.Now().UnixNano())
			val := rand.Intn(i)
			answer, _ := jsonparser.GetString(config, "message_list", "roll", "commandDone")
			answer = strings.Replace(answer, "%value%", strconv.Itoa(val), 1)
			answer = GetRegularData(answer, userID, vk)
			return answer, "", 0
		}
		answer, _ := jsonparser.GetString(config, "message_list", "roll", "negativeNumber")
		answer = GetRegularData(answer, userID, vk)
		return answer, "", 0
	} else if len(alist) == 3 {
		i1, err := strconv.Atoi(alist[1])
		if err != nil {
			answer, _ := jsonparser.GetString(config, "message_list", "roll", "notNumber")
			answer = GetRegularData(answer, userID, vk)
			return answer, "", 0
		}
		if i1 > 0 {
			i2, err := strconv.Atoi(alist[2])
			if err != nil {
				answer, _ := jsonparser.GetString(config, "message_list", "roll", "notNumber")
				answer = GetRegularData(answer, userID, vk)
				return answer, "", 0
			}
			if i2 > 0 {
				max, min := i1, i2
				if i1 < i2 {
					max, min = i2, i1
				} else if i1 == i2 {
					answer, _ := jsonparser.GetString(config, "message_list", "roll", "twoSimilarNumber")
					answer = GetRegularData(answer, userID, vk)
					return answer, "", 0
				}
				i := max - min
				rand.Seed(time.Now().UnixNano())
				val := rand.Intn(i) + min

				answer, _ := jsonparser.GetString(config, "message_list", "roll", "commandDone")
				answer = strings.Replace(answer, "%value%", strconv.Itoa(val), 1)
				answer = GetRegularData(answer, userID, vk)
				return answer, "", 0
			}
			answer, _ := jsonparser.GetString(config, "message_list", "roll", "negativeNumber")
			answer = GetRegularData(answer, userID, vk)
			return answer, "", 0
		}
		answer, _ := jsonparser.GetString(config, "message_list", "roll", "negativeNumber")
		answer = GetRegularData(answer, userID, vk)
		return answer, "", 0
	}
	answer, _ := jsonparser.GetString(config, "message_list", "roll", "nonNumber")
	answer = GetRegularData(answer, userID, vk)
	return answer, "", 0
}
