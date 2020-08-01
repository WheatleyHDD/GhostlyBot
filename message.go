package main

import (
	//"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"

	"strconv"
	"strings"
	"time"

	"github.com/SevereCloud/vksdk/api"
	"github.com/adam-lavrik/go-imath/ix"
	"github.com/buger/jsonparser"
	"github.com/imroc/req"
	"github.com/masatana/go-textdistance"
	scribble "github.com/nanobox-io/golang-scribble"
	"github.com/o1egl/govatar"
)

func clamp(value, maxv, minv int) int {
	return ix.Max(ix.Min(value, maxv), minv)
}

func getRegularData(a string, userID int, vk *api.VK) string {
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

func similarity(s1, s2 string) (similarity float64) {
	normal1, normal2 := strings.ToLower(s1), strings.ToLower(s2)
	return textdistance.JaroWinklerDistance(normal1, normal2)
}

func getRandAnswers(k string) []string {
	data, err := ioutil.ReadFile("PhasesDB.txt")
	if err != nil {
		panic(err)
	}
	s := string(data)
	ss := strings.Split(s, "\n")
	answers := make([][]string, 0)
	for i := range ss {
		phase := strings.Split(ss[i], "|")
		if strings.ToLower(k) == strings.ToLower(phase[0]) {
			answers = append(answers, phase)
		}
	}
	rand.Seed(time.Now().UnixNano())
	t := 0
	if len(answers)-1 != 0 {
		t = rand.Intn(len(answers) - 1)
	}
	return answers[t]
}

func remove(slice []int, s int) []int {
	return append(slice[:s], slice[s+1:]...)
}

func getAnsw(a string, userID int, vk *api.VK) (string, string, float64) {
	var lastSimilarity float64
	//var lastPhase [3]string
	data, err := ioutil.ReadFile("PhasesDB.txt")
	if err != nil {
		panic(err)
	}
	s := string(data)
	ss := strings.Split(s, "\n")
	lastPhase := strings.Split(ss[0], "|")
	for i := range ss {
		phase := strings.Split(ss[i], "|")

		sim := similarity(a, phase[0])
		if sim > lastSimilarity {
			lastPhase = phase
			lastSimilarity = sim
		}
	}

	gra := getRandAnswers(lastPhase[0])
	wreg := getRegularData(gra[1], userID, vk)
	attwreg := getRegularData(gra[2], userID, vk)
	return strings.Replace(wreg, "-n", "\n", -1), attwreg, lastSimilarity
}

type statistic struct {
	Command string
	Count   int
}

type stat_que struct {
	Que   string
	Count int
}

func getOtvet(a string, userID int, vk *api.VK, isChat bool, ChatID, BotID int, helpText string, attachment string, MID int, onWall bool, db *scribble.Driver) (otvet string, attachments string, similar float64) {
	config := configParse()
	alist := strings.Split(a, " ")
	//fmt.Println(alist)
	switch strings.ToLower(alist[0]) {
	case "найди":
		if len(alist) > 1 {
			switch strings.ToLower(alist[1]) {
			case "гиф", "gif":

				storeStatistic("найди гиф", db)

				alist[0], alist[1] = "", ""
				que := ""
				for i, v := range alist {
					if v != "" {
						if i == 2 {
							que = strings.Join([]string{que, v}, "")
						} else {
							que = strings.Join([]string{que, v}, " ")
						}

					}
				}

				if que != "" {
					storeStatisticFindQue(que, db)
				} else {
					storeStatisticFindQue("empty", db)
				}

				f, s := "", ""

				if onWall {
					f, s, _ = FindGif(que, vk, 2)
				} else {
					f, s, _ = FindGif(que, vk, 5)
				}

				return f, s, 0
			case "фото":

				storeStatistic("найди фото", db)

				alist[0], alist[1] = "", ""
				que := ""
				for i, v := range alist {
					if v != "" {
						if i == 2 {
							que = strings.Join([]string{que, v}, "")
						} else {
							que = strings.Join([]string{que, v}, " ")
						}

					}
				}

				if que != "" {
					storeStatisticFindQue(que, db)
				} else {
					storeStatisticFindQue("empty", db)
				}

				f, s := "", ""
				if onWall {
					f, s, _ = FindPhoto(que, vk, 2)
				} else {
					f, s, _ = FindPhoto(que, vk, 10)
				}

				return f, s, 0
			case "видео":

				storeStatistic("найди видео", db)

				alist[0], alist[1] = "", ""
				que := ""
				for i, v := range alist {
					if v != "" {
						if i == 2 {
							que = strings.Join([]string{que, v}, "")
						} else {
							que = strings.Join([]string{que, v}, " ")
						}

					}
				}

				if que != "" {
					storeStatisticFindQue(que, db)
				} else {
					storeStatisticFindQue("empty", db)
				}

				f, s := "", ""
				if onWall {
					f, s, _ = FindVideo(que, vk, 2)
				} else {
					f, s, _ = FindVideo(que, vk, 4)
				}

				return f, s, 0
			}
			answer, _ := jsonparser.GetString(config, "message_list", "find", "TypeUndefined")
			return answer, "", 0
		}
		answer, _ := jsonparser.GetString(config, "message_list", "find", "TypeUndefined")
		return answer, "", 0
	case "кто", "who":

		if len(alist) > 1 {
			switch strings.ToLower(alist[1]) {
			case "я", "ты", "мы", "он", "они", "вы", "она", "я?", "ты?", "мы?", "он?", "они?", "вы?", "она?":
				reta, retattach, _ := getAnsw(a, userID, vk)
				return reta, retattach, 0
			}
			if !onWall {
				if isChat {

					storeStatistic("кто", db)

					chat, err := vk.MessagesGetChat(api.Params{
						"chat_id": ChatID,
					})
					if err != nil {
						panic(err)
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
						panic(err)
					}
					answer, _ := jsonparser.GetString(config, "message_list", "who", "commandDone")
					answer = strings.Replace(answer, "%value%", "[id"+strconv.Itoa(user[0].ID)+"|"+user[0].FirstName+" "+user[0].LastName+"]", 1)
					answer = getRegularData(answer, userID, vk)
					return answer, "", 0
				}
				answer, _ := jsonparser.GetString(config, "message_list", "who", "commandInPM")
				answer = getRegularData(answer, userID, vk)
				return answer, "", 0
			}
			answer, _ := jsonparser.GetString(config, "commandOnWallDontWorks")
			answer = getRegularData(answer, userID, vk)
			return answer, "", 0
		}
	case "команды", "помощь":

		storeStatistic("команды", db)

		return strings.Replace(helpText, "-n", "\n", -1), "", 0
	case "лобстер", "lobster":
		if !onWall {
			alist[0] = ""
			label := ""
			for i, v := range alist {
				if v != "" {
					if i == 1 {
						label = strings.Join([]string{label, v}, "")
					} else {
						label = strings.Join([]string{label, v}, " ")
					}

				}
			}

			if label != "" {
				if attachment != "" {

					storeStatistic("лобстер", db)

					vk.MessagesSend(api.Params{
						"peer_id":   ChatID + 2000000000,
						"message":   "Подожди немного. Эта операция не быстрая и может занять от пары секунд до 1 минуты",
						"random_id": 0,
						"reply_to":  MID,
					})

					r, _ := req.Get(attachment)
					r.ToFile(strconv.Itoa(userID) + "_source.jpg")
					meme := CreateLobsterMeme(label, strconv.Itoa(userID)+"_source.jpg", strconv.Itoa(userID)+".png", ChatID, vk)
					return "Держи", meme, 0
				}
				answer, _ := jsonparser.GetString(config, "message_list", "lobster", "nonPhoto")
				answer = getRegularData(answer, userID, vk)
				return answer, "", 0
			}
			answer, _ := jsonparser.GetString(config, "message_list", "lobster", "nonText")
			answer = getRegularData(answer, userID, vk)
			return answer, "", 0
		}
		answer, _ := jsonparser.GetString(config, "commandOnWallDontWorks")
		answer = getRegularData(answer, userID, vk)
		return answer, "", 0
	case "фас":

		storeStatistic("фас", db)

		return "Гав-гав", "", 0
	case "выбери", "выбери:":
		if strings.Contains(a, " или ") {

			storeStatistic("выбери", db)

			a = strings.Replace(a, "выбери ", "", 1)
			opinions := strings.Split(a, " или ")
			rand.Seed(time.Now().UnixNano())
			opinion := rand.Intn(len(opinions))

			answer, _ := jsonparser.GetString(config, "message_list", "choose", "commandDone")
			answer = strings.Replace(answer, "%value%", opinions[opinion], 1)
			answer = getRegularData(answer, userID, vk)
			return answer, "", 0
		}
		answer, _ := jsonparser.GetString(config, "message_list", "choose", "notOpinion")
		answer = getRegularData(answer, userID, vk)
		return answer, "", 0
	case "рулеточка", "roll":

		storeStatistic("рулеточка", db)

		if len(alist) == 2 {
			i, err := strconv.Atoi(alist[1])
			if err != nil {
				answer, _ := jsonparser.GetString(config, "message_list", "roll", "notNumber")
				answer = getRegularData(answer, userID, vk)
				return answer, "", 0
			}
			if i > 0 {
				rand.Seed(time.Now().UnixNano())
				val := rand.Intn(i)
				answer, _ := jsonparser.GetString(config, "message_list", "roll", "commandDone")
				answer = strings.Replace(answer, "%value%", strconv.Itoa(val), 1)
				answer = getRegularData(answer, userID, vk)
				return answer, "", 0
			}
			answer, _ := jsonparser.GetString(config, "message_list", "roll", "negativeNumber")
			answer = getRegularData(answer, userID, vk)
			return answer, "", 0
		} else if len(alist) == 3 {
			i1, err := strconv.Atoi(alist[1])
			if err != nil {
				answer, _ := jsonparser.GetString(config, "message_list", "roll", "notNumber")
				answer = getRegularData(answer, userID, vk)
				return answer, "", 0
			}
			if i1 > 0 {
				i2, err := strconv.Atoi(alist[2])
				if err != nil {
					answer, _ := jsonparser.GetString(config, "message_list", "roll", "notNumber")
					answer = getRegularData(answer, userID, vk)
					return answer, "", 0
				}
				if i2 > 0 {
					max, min := i1, i2
					if i1 < i2 {
						max, min = i2, i1
					} else if i1 == i2 {
						answer, _ := jsonparser.GetString(config, "message_list", "roll", "twoSimilarNumber")
						answer = getRegularData(answer, userID, vk)
						return answer, "", 0
					}
					i := max - min
					rand.Seed(time.Now().UnixNano())
					val := rand.Intn(i) + min

					answer, _ := jsonparser.GetString(config, "message_list", "roll", "commandDone")
					answer = strings.Replace(answer, "%value%", strconv.Itoa(val), 1)
					answer = getRegularData(answer, userID, vk)
					return answer, "", 0
				}
				answer, _ := jsonparser.GetString(config, "message_list", "roll", "negativeNumber")
				answer = getRegularData(answer, userID, vk)
				return answer, "", 0
			}
			answer, _ := jsonparser.GetString(config, "message_list", "roll", "negativeNumber")
			answer = getRegularData(answer, userID, vk)
			return answer, "", 0
		}
		answer, _ := jsonparser.GetString(config, "message_list", "roll", "nonNumber")
		answer = getRegularData(answer, userID, vk)
		return answer, "", 0
	case "инфа":
		if len(alist) > 1 {

			storeStatistic("инфа", db)

			rand.Seed(time.Now().UnixNano())
			percent := rand.Intn(100)

			answer, _ := jsonparser.GetString(config, "message_list", "infa", "commandDone")
			answer = strings.Replace(answer, "%value%", strconv.Itoa(percent), 1)
			answer = getRegularData(answer, userID, vk)
			return answer, "", 0
		}
		answer, _ := jsonparser.GetString(config, "message_list", "infa", "undefined")
		answer = getRegularData(answer, userID, vk)
		return answer, "", 0
	case "когда":
		if len(alist) > 1 {

			storeStatistic("когда", db)

			h := []string{"дней", "минут", "лет", "секунд", "месяцев", "часов"}

			rand.Seed(time.Now().UnixNano())
			date := strconv.Itoa(rand.Intn(1000)) + " " + h[rand.Intn(len(h)-1)]

			answer, _ := jsonparser.GetString(config, "message_list", "when", "commandDone")
			answer = strings.Replace(answer, "%value%", date, 1)
			answer = getRegularData(answer, userID, vk)
			return answer, "", 0
		}
		answer, _ := jsonparser.GetString(config, "message_list", "when", "undefined")
		answer = getRegularData(answer, userID, vk)
		return answer, "", 0
	case "о":
		if len(alist) == 2 || strings.ToLower(alist[1]) == "боте" {

			storeStatistic("о боте", db)

			answer, _ := jsonparser.GetString(config, "about")
			answer = getRegularData(answer, userID, vk)
			return answer, "", 0

		}
	case "ген", "gen":
		if len(alist) > 1 {
			switch strings.ToLower(alist[1]) {
			case "ава", "ava":

				storeStatistic("ген ава", db)

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
						answer = getRegularData(answer, userID, vk)
						return answer, "", 0
					}

					if err != nil {
						answer, _ := jsonparser.GetString(config, "message_list", "genAva", "error")
						answer = getRegularData(answer, userID, vk)
						return answer, "", 0
					}

					ava, err := os.Open(strconv.Itoa(userdata[0].ID) + ".png")
					if err != nil {
						answer, _ := jsonparser.GetString(config, "message_list", "genAva", "error")
						answer = getRegularData(answer, userID, vk)
						return answer, "", 0
					}

					photosPhoto, err := vk.UploadMessagesPhoto(2000000000+ChatID, ava)

					_ = os.Remove(strconv.Itoa(userdata[0].ID) + ".png")

					answer, _ := jsonparser.GetString(config, "message_list", "genAva", "commandDone")
					answer = strings.Replace(answer, "%value%", nickname, 1)
					answer = getRegularData(answer, userID, vk)

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
						answer = getRegularData(answer, userID, vk)
						return answer, "", 0
					}

					ava, err := os.Open(strconv.Itoa(userdata[0].ID) + ".png")
					if err != nil {
						answer, _ := jsonparser.GetString(config, "message_list", "genAva", "error")
						answer = getRegularData(answer, userID, vk)
						return answer, "", 0
					}

					photosPhoto, err := vk.UploadMessagesPhoto(2000000000+ChatID, ava)

					_ = os.Remove(strconv.Itoa(userdata[0].ID) + ".png")

					answer, _ := jsonparser.GetString(config, "message_list", "genAva", "commandDone")
					answer = strings.Replace(answer, "%value%", "[id"+strconv.Itoa(userdata[0].ID)+"|"+userdata[0].FirstName+" "+userdata[0].LastName+"]", 1)
					answer = getRegularData(answer, userID, vk)

					return answer, "photo" + strconv.Itoa(photosPhoto[0].OwnerID) + "_" + strconv.Itoa(photosPhoto[0].ID), 0
				}
				answer, _ := jsonparser.GetString(config, "message_list", "genAva", "incorrectField")
				answer = getRegularData(answer, userID, vk)
				return answer, "", 0
			}
			answer, _ := jsonparser.GetString(config, "message_list", "genAva", "undefined")
			answer = getRegularData(answer, userID, vk)
			return answer, "", 0
		}
		answer, _ := jsonparser.GetString(config, "message_list", "genAva", "undefined")
		answer = getRegularData(answer, userID, vk)
		return answer, "", 0
	}
	reta, retattach, _ := getAnsw(a, userID, vk)
	return reta, retattach, 0
}

func storeStatistic(command string, db *scribble.Driver) {
	config := configParse()
	statSetting, _ := jsonparser.GetBoolean(config, "statistic")
	if statSetting {
		stat := statistic{command, 0}
		_ = db.Read("statistic", command, &stat)
		stat.Count++
		db.Write("statistic", command, statistic{command, stat.Count})
	}
}

func storeStatisticFindQue(command string, db *scribble.Driver) {
	config := configParse()
	statSetting, _ := jsonparser.GetBoolean(config, "statistic")
	if statSetting {
		statq := stat_que{command, 0}
		_ = db.Read("stat_que", command, &statq)
		statq.Count = statq.Count + 1
		db.Write("stat_que", command, stat_que{command, statq.Count})
	}
}
