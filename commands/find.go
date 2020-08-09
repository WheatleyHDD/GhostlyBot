package commands

import (
	"github.com/SevereCloud/vksdk/api"
	"github.com/buger/jsonparser"
	scribble "github.com/nanobox-io/golang-scribble"
	. "vkbot/utils"
)

func Find(ftype, query string, onWall bool, db *scribble.Driver, vk *api.VK) (otvet string, attachments string, similar float64) {
	config := ConfigParse()
	switch ftype {
	case "гиф", "gif":

		StoreStatistic("найди гиф", db)

		f, s := "", ""

		if onWall {
			f, s, _ = FindGif(query, vk, 2)
		} else {
			f, s, _ = FindGif(query, vk, 5)
		}

		return f, s, 0
	case "фото":

		StoreStatistic("найди фото", db)

		f, s := "", ""

		if onWall {
			f, s, _ = FindPhoto(query, vk, 2)
		} else {
			f, s, _ = FindPhoto(query, vk, 10)
		}

		return f, s, 0
	case "видео":

		StoreStatistic("найди видео", db)

		f, s := "", ""
		if onWall {
			f, s, _ = FindVideo(query, vk, 2)
		} else {
			f, s, _ = FindVideo(query, vk, 4)
		}

		return f, s, 0
	}
	answer, _ := jsonparser.GetString(config, "message_list", "find", "TypeUndefined")
	return answer, "", 0
}
