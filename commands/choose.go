package commands

import (
	"github.com/SevereCloud/vksdk/api"
	"github.com/buger/jsonparser"
	scribble "github.com/nanobox-io/golang-scribble"
	"math/rand"
	"strings"
	"time"
	. "vkbot/utils"
)

func Choose(a string, userID int, db *scribble.Driver, vk *api.VK) (otvet string, attachments string, similar float64) {
	config := ConfigParse()
	if strings.Contains(a, " или ") {

		StoreStatistic("выбери", db)

		a = strings.Replace(a, "выбери ", "", 1)
		a = strings.Replace(a, "Выбери ", "", 1)
		opinions := strings.Split(a, " или ")
		rand.Seed(time.Now().UnixNano())
		opinion := rand.Intn(len(opinions))

		answer, _ := jsonparser.GetString(config, "message_list", "choose", "commandDone")
		answer = strings.Replace(answer, "%value%", opinions[opinion], 1)
		answer = GetRegularData(answer, userID, vk)
		return answer, "", 0
	}
	answer, _ := jsonparser.GetString(config, "message_list", "choose", "notOpinion")
	answer = GetRegularData(answer, userID, vk)
	return answer, "", 0
}
