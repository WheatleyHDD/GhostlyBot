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

func remove(slice []int, s int) []int {
	return append(slice[:s], slice[s+1:]...)
}

func Who(onWall, isChat bool, ChatID, userID, BotID int, db *scribble.Driver, vk *api.VK) (otvet string, attachments string, similar float64) {
	config := ConfigParse()
	if !onWall {
		if isChat {

			StoreStatistic("кто", db)

			chat, err := vk.MessagesGetChat(api.Params{
				"chat_id": ChatID,
			})
			if err != nil {
				return "Ошиб очка: " + err.Error(), "", 0
			}
			chUsers := chat.Users
			for i, id := range chat.Users {
				if id == BotID || id < 0 {
					chUsers = remove(chat.Users, i)
				}
			}
			rand.Seed(time.Now().UnixNano())
			h := rand.Intn(len(chUsers) - 1)
			user, err := vk.UsersGet(api.Params{
				"user_ids": chUsers[h],
			})
			if err != nil {
				return "Ошиб очка: " + err.Error(), "", 0
			}
			answer, _ := jsonparser.GetString(config, "message_list", "who", "commandDone")
			answer = strings.Replace(answer, "%value%", "[id"+strconv.Itoa(user[0].ID)+"|"+user[0].FirstName+" "+user[0].LastName+"]", 1)
			answer = GetRegularData(answer, userID, vk)
			return answer, "", 0
		}
		answer, _ := jsonparser.GetString(config, "message_list", "who", "commandInPM")
		answer = GetRegularData(answer, userID, vk)
		return answer, "", 0
	}
	answer, _ := jsonparser.GetString(config, "commandOnWallDontWorks")
	answer = GetRegularData(answer, userID, vk)
	return answer, "", 0
}
