package commands

import (
	"github.com/SevereCloud/vksdk/api"
	"github.com/buger/jsonparser"
	scribble "github.com/nanobox-io/golang-scribble"
	"github.com/o1egl/govatar"
	"os"
	"strconv"
	"strings"
	. "vkbot/utils"
)

func Gen(alist []string, userID, ChatID int, db *scribble.Driver, vk *api.VK) (otvet string, attachments string, similar float64) {
	config := ConfigParse()
	if len(alist) > 1 {
		switch strings.ToLower(alist[1]) {
		case "ава", "ava":

			StoreStatistic("ген ава", db)

			if len(alist) > 3 {

				sex := alist[2]
				alist[0], alist[1], alist[2] = "", "", ""

				nickname := ""

				for i, v := range alist {
					if v != "" {
						if i == 1 {
							nickname = strings.Join([]string{nickname, v}, "")
						} else {
							nickname = strings.Join([]string{nickname, v}, " ")
						}

					}
				}

				userdata, err := vk.UsersGet(api.Params{
					"user_ids": userID,
				})

				switch sex {
				case "муж", "male":
					err = govatar.GenerateFileForUsername(govatar.MALE, nickname, strconv.Itoa(userdata[0].ID)+".png")
				case "жен", "female":
					err = govatar.GenerateFileForUsername(govatar.FEMALE, nickname, strconv.Itoa(userdata[0].ID)+".png")
				default:
					answer, _ := jsonparser.GetString(config, "message_list", "genAva", "incorrectField")
					answer = GetRegularData(answer, userID, vk)
					return answer, "", 0
				}

				if err != nil {
					answer, _ := jsonparser.GetString(config, "message_list", "genAva", "error")
					answer = GetRegularData(answer, userID, vk)
					return answer, "", 0
				}

				ava, err := os.Open(strconv.Itoa(userdata[0].ID) + ".png")
				if err != nil {
					answer, _ := jsonparser.GetString(config, "message_list", "genAva", "error")
					answer = GetRegularData(answer, userID, vk)
					return answer, "", 0
				}

				photosPhoto, err := vk.UploadMessagesPhoto(2000000000+ChatID, ava)

				_ = os.Remove(strconv.Itoa(userdata[0].ID) + ".png")

				answer, _ := jsonparser.GetString(config, "message_list", "genAva", "commandDone")
				answer = strings.Replace(answer, "%value%", nickname, 1)
				answer = GetRegularData(answer, userID, vk)

				return answer, "photo" + strconv.Itoa(photosPhoto[0].OwnerID) + "_" + strconv.Itoa(photosPhoto[0].ID), 0
			} else if len(alist) == 2 {
				userdata, err := vk.UsersGet(api.Params{
					"user_ids":  userID,
					"fields":    "sex",
					"name_case": "gen",
				})

				if userdata[0].Sex == 1 {
					err = govatar.GenerateFileForUsername(govatar.FEMALE, strconv.Itoa(userdata[0].ID), strconv.Itoa(userdata[0].ID)+".png")
				} else {
					err = govatar.GenerateFileForUsername(govatar.MALE, strconv.Itoa(userdata[0].ID), strconv.Itoa(userdata[0].ID)+".png")
				}

				if err != nil {
					answer, _ := jsonparser.GetString(config, "message_list", "genAva", "error")
					answer = GetRegularData(answer, userID, vk)
					return answer, "", 0
				}

				ava, err := os.Open(strconv.Itoa(userdata[0].ID) + ".png")
				if err != nil {
					answer, _ := jsonparser.GetString(config, "message_list", "genAva", "error")
					answer = GetRegularData(answer, userID, vk)
					return answer, "", 0
				}

				photosPhoto, err := vk.UploadMessagesPhoto(2000000000+ChatID, ava)
				if err != nil {
					return "Ошиб очка: " + err.Error(), "", 0
				}

				_ = os.Remove(strconv.Itoa(userdata[0].ID) + ".png")

				answer, _ := jsonparser.GetString(config, "message_list", "genAva", "commandDone")
				answer = strings.Replace(answer, "%value%", "[id"+strconv.Itoa(userdata[0].ID)+"|"+userdata[0].FirstName+" "+userdata[0].LastName+"]", 1)
				answer = GetRegularData(answer, userID, vk)

				return answer, "photo" + strconv.Itoa(photosPhoto[0].OwnerID) + "_" + strconv.Itoa(photosPhoto[0].ID), 0
			}
			answer, _ := jsonparser.GetString(config, "message_list", "genAva", "incorrectField")
			answer = GetRegularData(answer, userID, vk)
			return answer, "", 0
		}
		answer, _ := jsonparser.GetString(config, "message_list", "genAva", "undefined")
		answer = GetRegularData(answer, userID, vk)
		return answer, "", 0
	}
	answer, _ := jsonparser.GetString(config, "message_list", "genAva", "undefined")
	answer = GetRegularData(answer, userID, vk)
	return answer, "", 0
}
