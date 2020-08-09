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

func When(alist []string, userID int, db *scribble.Driver, vk *api.VK) (otvet string, attachments string, similar float64) {
	config := ConfigParse()

	if len(alist) > 1 {

		StoreStatistic("когда", db)

		h := []string{"дней", "минут", "лет", "секунд", "месяцев", "часов"}

		rand.Seed(time.Now().UnixNano())
		date := strconv.Itoa(rand.Intn(1000)) + " " + h[rand.Intn(len(h)-1)]

		answer, _ := jsonparser.GetString(config, "message_list", "when", "commandDone")
		answer = strings.Replace(answer, "%value%", date, 1)
		answer = GetRegularData(answer, userID, vk)
		return answer, "", 0
	}
	answer, _ := jsonparser.GetString(config, "message_list", "when", "undefined")
	answer = GetRegularData(answer, userID, vk)
	return answer, "", 0
}
