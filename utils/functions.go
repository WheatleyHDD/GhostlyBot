package utils

import (
	"github.com/SevereCloud/vksdk/api"
	"github.com/buger/jsonparser"
	scribble "github.com/nanobox-io/golang-scribble"
	"log"
	"strings"
)

type statistic struct {
	Command string
	Count   int
}

type stat_que struct {
	Que   string
	Count int
}

func StoreStatistic(command string, db *scribble.Driver) {
	config := ConfigParse()
	statSetting, _ := jsonparser.GetBoolean(config, "statistic")
	if statSetting {
		stat := statistic{command, 0}
		_ = db.Read("statistic", command, &stat)
		stat.Count++
		db.Write("statistic", command, statistic{command, stat.Count})
	}
}

func StoreStatisticFindQue(command string, db *scribble.Driver) {
	config := ConfigParse()
	statSetting, _ := jsonparser.GetBoolean(config, "statistic")
	if statSetting {
		statq := stat_que{command, 0}
		_ = db.Read("stat_que", command, &statq)
		statq.Count = statq.Count + 1
		db.Write("stat_que", command, stat_que{command, statq.Count})
	}
}

func GetRegularData(a string, userID int, vk *api.VK) string {
	regulars := []string{"%username%", "%usersurname%", "%userbirthdate%", "%userphoto%"}
	userdata, err := vk.UsersGet(api.Params{
		"user_ids":  userID,
		"fields":    "photo_id, verified, sex, bdate, city, country, home_town, has_photo, photo_50, photo_100, photo_200_orig, photo_200, photo_400_orig, photo_max, photo_max_orig, online, domain, has_mobile, contacts, site, education, universities, schools, status, last_seen, followers_count, common_count, occupation, nickname, relatives, relation, personal, connections, exports, activities, interests, music, movies, tv, books, games, about, quotes, can_post, can_see_all_posts, can_see_audio, can_write_private_message, can_send_friend_request, is_favorite, is_hidden_from_feed, timezone, screen_name, maiden_name, crop_photo, is_friend, friend_status, career, military, blacklisted, blacklisted_by_me, can_be_invited_group",
		"name_case": "Nom",
	})
	if err != nil {
		log.Fatal(err)
	}

	retText := a
	for _, v := range regulars {
		switch v {
		case "%username%":
			retText = strings.Replace(retText, v, userdata[0].FirstName, -1)
		case "%usersurname%":
			retText = strings.Replace(retText, v, userdata[0].LastName, -1)
		case "%userbirthdate%":
			if userdata[0].Bdate == "" {
				retText = strings.Replace(retText, v, "<ДАТА РОЖДЕНИЯ ЗАКРЫТА>", -1)
			} else {
				retText = strings.Replace(retText, v, userdata[0].Bdate, -1)
			}
		case "%userphoto%":
			retText = strings.Replace(retText, v, strings.Join([]string{"photo", userdata[0].PhotoID}, ""), -1)
		}

	}
	return retText
}
