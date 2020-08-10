package commands

import (
	"os"
	"strconv"
	"strings"
	. "vkbot/utils"

	"github.com/SevereCloud/vksdk/api"
	"github.com/buger/jsonparser"
	"github.com/evalphobia/google-tts-go/googletts"
	"github.com/imroc/req"
	scribble "github.com/nanobox-io/golang-scribble"
)

func Say(alist []string, userID, ChatID int, db *scribble.Driver, vk *api.VK) (otvet string, attachments string, similar float64) {
	config := ConfigParse()

	if len(alist) > 1 {

		StoreStatistic("скажи", db)

		lang := alist[0]
		alist[0] = ""
		que := ""
		for i, v := range alist {
			if v != "" {
				if i == 1 {
					que = strings.Join([]string{que, v}, "")
				} else {
					que = strings.Join([]string{que, v}, " ")
				}

			}
		}

		if lang == "say" {
			lang = "en"
		} else {
			lang = "ru"
		}

		url, err := googletts.GetTTSURL(que, lang)
		if err != nil {
			return "Ошиб очка: " + err.Error(), "", 0
		}

		r, _ := req.Get(url)
		r.ToFile(strconv.Itoa(userID) + "_message.ogg")
		file, err := os.Open(strconv.Itoa(userID) + "_message.ogg")
		if err != nil {
			return "Ошиб очка: " + err.Error(), "", 0
		}
		docsDoc, err := vk.UploadMessagesDoc(ChatID+2000000000, "audio_message", "Говорю "+que, "", file)
		_ = os.Remove(strconv.Itoa(userID) + "_message.ogg")
		return "", "doc" + strconv.Itoa(docsDoc.AudioMessage.OwnerID) + "_" + strconv.Itoa(docsDoc.AudioMessage.ID), 0
	}
	answer, _ := jsonparser.GetString(config, "message_list", "say", "undefined")
	answer = GetRegularData(answer, userID, vk)
	return answer, "", 0
}
