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

func Infa(alist []string, userID int, db *scribble.Driver, vk *api.VK) (otvet string, attachments string, similar float64) {
	config := ConfigParse()

	if len(alist) > 1 {

		StoreStatistic("инфа", db)

		rand.Seed(time.Now().UnixNano())
		percent := rand.Intn(100)

		answer, _ := jsonparser.GetString(config, "message_list", "infa", "commandDone")
		answer = strings.Replace(answer, "%value%", strconv.Itoa(percent), 1)
		answer = GetRegularData(answer, userID, vk)
		return answer, "", 0
	}
	answer, _ := jsonparser.GetString(config, "message_list", "infa", "undefined")
	answer = GetRegularData(answer, userID, vk)
	return answer, "", 0
}
