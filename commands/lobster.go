package commands

import (
	"github.com/SevereCloud/vksdk/api"
	"github.com/buger/jsonparser"
	"github.com/imroc/req"
	scribble "github.com/nanobox-io/golang-scribble"
	"strconv"
	"strings"
	. "vkbot/utils"
)

func Lobster(text, photo string, ChatID, MID, userID int, onWall bool, db *scribble.Driver, vk *api.VK) (otvet string, attachments string, similar float64) {
	config := ConfigParse()

	if !onWall {

		if text != "" {
			if photo != "" {

				StoreStatistic("лобстер", db)

				vk.MessagesSend(api.Params{
					"peer_id":   ChatID + 2000000000,
					"message":   "Подожди немного. Эта операция не быстрая и может занять от пары секунд до 1 минуты",
					"random_id": 0,
					"reply_to":  MID,
				})

				r, _ := req.Get(photo)
				r.ToFile(strconv.Itoa(userID) + "_source.jpg")
				meme := CreateLobsterMeme(text, strconv.Itoa(userID)+"_source.jpg", strconv.Itoa(userID)+".png", ChatID, vk)
				if strings.HasPrefix(meme, "Ошиб очка") {
					return meme, "", 0
				}
				return "Держи", meme, 0
			}
			answer, _ := jsonparser.GetString(config, "message_list", "lobster", "nonPhoto")
			answer = GetRegularData(answer, userID, vk)
			return answer, "", 0
		}
		answer, _ := jsonparser.GetString(config, "message_list", "lobster", "nonText")
		answer = GetRegularData(answer, userID, vk)
		return answer, "", 0
	}
	answer, _ := jsonparser.GetString(config, "commandOnWallDontWorks")
	answer = GetRegularData(answer, userID, vk)
	return answer, "", 0
}
