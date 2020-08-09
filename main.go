package main

import (
	//"fmt"

	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"strings"

	"math/rand"
	"time"

	"github.com/SevereCloud/vksdk/api"
	longpoll "github.com/SevereCloud/vksdk/longpoll-user"
	wrapper "github.com/SevereCloud/vksdk/longpoll-user/v3"
	"github.com/buger/jsonparser"
	scribble "github.com/nanobox-io/golang-scribble"

	. "vkbot/utils"
)

type LastPost struct {
	WallID int64
	ID     int
}

func main() {
	db, err := scribble.New("./BDs", nil)
	if err != nil {
		log.Fatal(err)
	}

	config := ConfigParse()

	wallCount := 0

	jsonparser.ArrayEach(config, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		wallCount++
	}, "walls")

	accessToken, err := jsonparser.GetString(config, "access_token")
	if err != nil {
		log.Fatal(err)
	}
	appeal, err := jsonparser.GetString(config, "appeal")
	if err != nil {
		log.Fatal(err)
	}
	helpText, err := jsonparser.GetString(config, "help_text")
	if err != nil {
		log.Fatal(err)
	}

	//fmt.Println(config)
	vk := api.NewVK(accessToken)
	//fmt.Println(settings[1])

	users, err := vk.UsersGet(api.Params{})
	if err != nil {
		log.Fatal(err)
	}
	BotID := users[0].ID

	log.Printf("Инициализация бота с ID %v\n", BotID)

	go pingOnline(vk)
	go getWallPosts(vk, db, config, appeal, BotID, helpText, wallCount)

	lp, err := longpoll.NewLongpoll(vk, 3)
	if err != nil {
		main()
	}

	w := wrapper.NewWrapper(lp)

	// event with code 4
	w.OnNewMessage(func(m wrapper.NewMessage) {

		mess, err := vk.MessagesGetByIDExtended(api.Params{
			"message_ids": m.MessageID,
		})
		if err != nil {
			log.Fatal(err)
		}

		//fmt.Println(m.MessageID)

		if mess.Items[0].FromID != BotID {
			if m.PeerID-2000000000 < 0 {

				photoURL := ""

				for _, a := range mess.Items[0].Attachments {
					//fmt.Println(a.Type)
					if a.Type == "photo" {
						//fmt.Println(a.Photo.Sizes[len(a.Photo.Sizes)-1].URL)
						photoURL = a.Photo.Sizes[len(a.Photo.Sizes)-1].URL
						break
					}
				}
				//fmt.Println(photoURL)

				log.Printf("Получено сообщение от пользователя \"%s\" с текстом: %s\n", mess.Profiles[0].FirstName+" "+mess.Profiles[0].LastName, m.Text)

				otv, attach, _ := getOtvet(m.Text, mess.Items[0].FromID, vk, false, m.PeerID-2000000000, BotID, helpText, photoURL, 0, false, db)

				data, err := ioutil.ReadFile("blacklist.txt")
				if err != nil {
					panic(err)
				}
				s := string(data)
				ss := strings.Split(s, "\n")
				for i := range ss {
					if ss[i] != "" || ss[i] != " " {
						if strings.Contains(strings.ToLower(m.Text), ss[i]) {
							onAttemptToBlock := "Маму свою заблокируй! Слышал, ты, пузырик пакостный?!"
							onAttemptToBlock, _ = jsonparser.GetString(config, "onAttemptToBlock")
							otv, attach = onAttemptToBlock, ""
						}
					}
				}

				_, err = vk.MessagesSend(api.Params{
					"peer_id":    m.PeerID,
					"message":    otv,
					"attachment": attach,
					"random_id":  0,
				})

				log.Printf("Отправлен ответ пользователю \"%s\" с текстом: \"%s\" и вложением: %s\n", mess.Profiles[0].FirstName+" "+mess.Profiles[0].LastName, otv, attach)

				if err != nil {
					log.Fatal(err)
				}
				//fmt.Printf("4 wrapper.NewMessage: %v\n", m)
			} else {
				//fmt.Printf("4 wrapper.NewMessage: %v\n", m)
				if strings.HasPrefix(strings.ToLower(m.Text), strings.ToLower(appeal)) || strings.HasPrefix(strings.ToLower(m.Text), "бот, ") {
					messageArr := strings.Split(m.Text, " ")
					messageArr[0] = ""

					message := ""

					for i, v := range messageArr {

						if v != "" {
							if i == 1 {
								message = strings.Join([]string{message, v}, "")
							} else {
								message = strings.Join([]string{message, v}, " ")
							}

						}
					}

					fmt.Println(message)

					photoURL := ""

					for _, a := range mess.Items[0].Attachments {
						if a.Type == "photo" {
							photoURL = a.Photo.Sizes[len(a.Photo.Sizes)-1].URL
							break
						}
					}

					chatInfo, _ := vk.MessagesGetChat(api.Params{
						"chat_id": m.PeerID - 2000000000,
					})

					log.Printf("Получено сообщение от пользователя \"%s\" из беседы \"%s\" с текстом: %s\n", mess.Profiles[0].FirstName+" "+mess.Profiles[0].LastName, chatInfo.Title, m.Text)

					//message = strings.Replace(strings.ToLower(message), appeal, "", 1)
					//message = strings.Replace(strings.ToLower(message), "бот, ", "", 1)

					otv, attach, _ := getOtvet(message, mess.Items[0].FromID, vk, true, m.PeerID-2000000000, BotID, helpText, photoURL, m.MessageID, false, db)

					data, err := ioutil.ReadFile("blacklist.txt")
					if err != nil {
						panic(err)
					}
					s := string(data)
					ss := strings.Split(s, "\n")
					for i := range ss {
						if ss[i] != "" || ss[i] != " " {
							if strings.Contains(strings.ToLower(m.Text), ss[i]) {
								onAttemptToBlock := "Маму свою заблокируй! Слышал, ты, пузырик пакостный?!"
								onAttemptToBlock, _ = jsonparser.GetString(config, "onAttemptToBlock")
								otv, attach = onAttemptToBlock, ""
							}
						}
					}

					_, err = vk.MessagesSend(api.Params{
						"peer_id":    m.PeerID,
						"message":    otv,
						"attachment": attach,
						"random_id":  0,
						"reply_to":   m.MessageID,
					})

					log.Printf("Отправлено сообщение для пользователя \"%s\" из беседы \"%s\" с текстом: \"%s\" и вложением: \"%s\"\n", mess.Profiles[0].FirstName+" "+mess.Profiles[0].LastName, chatInfo.Title, otv, attach)

					if err != nil {
						log.Fatal(err)
					}
				}
			}
		}
	})

	if err := lp.Run(); err != nil {
		main()
	}

	lp.Shutdown()
	//lp.Client.CloseIdleConnections()

	/*
		var vopros string
		for {
			fmt.Print("Вопрос: ")
			myscanner := bufio.NewScanner(os.Stdin)
			myscanner.Scan()
			vopros = myscanner.Text()
			if vopros != "exit" {
				otv, attach, _ := getOtvet(vopros)
				fmt.Println("Ответ: ", otv)
				fmt.Println("Вложение: ", attach)
				vopros = ""
				fmt.Println(" ")

			} else {
				break
			}
		}
	*/
}

func getWallPosts(vk *api.VK, db *scribble.Driver, config []byte, appeal string, BotID int, helpText string, wallCount int) {
	for {
		log.Println("Проверка стен")
		go checkWalls(vk, db, config, appeal, BotID, helpText)

		duration := 5000 / 24 / wallCount / 60

		time.Sleep(time.Minute * time.Duration(duration))
	}
}

func pingOnline(vk *api.VK) {
	for {
		vk.AccountSetOnline(api.Params{
			"voip": 0,
		})
		time.Sleep(time.Minute * 5)
	}
}

func checkWalls(vk *api.VK, db *scribble.Driver, config []byte, appeal string, BotID int, helpText string) {
	jsonparser.ArrayEach(config, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		ownid, err := jsonparser.GetInt(value)
		if err != nil {
			log.Fatal(err)
		}
		posts, err := vk.WallGetExtended(api.Params{
			"owner_id": ownid,
			"filter":   "others",
			"count":    1,
		})
		if err != nil {
			log.Println(err)
		} else {
			last_post := LastPost{}
			_ = db.Read("last_post_ids", strconv.Itoa(offset), &last_post)
			if posts.Items[0].ID > last_post.ID {
				if strings.HasPrefix(strings.ToLower(posts.Items[0].Text), strings.ToLower(appeal)) || strings.HasPrefix(strings.ToLower(posts.Items[0].Text), "бот, ") {
					messageArr := strings.Split(posts.Items[0].Text, " ")
					messageArr[0] = ""

					message := ""

					for i, v := range messageArr {

						if v != "" {
							if i == 1 {
								message = strings.Join([]string{message, v}, "")
							} else {
								message = strings.Join([]string{message, v}, " ")
							}

						}
					}

					log.Printf("Получено сообщение со стены %v с текстом: %s\n", ownid, message)

					photoURL := ""

					for _, a := range posts.Items[0].Attachments {
						if a.Type == "photo" {
							photoURL = a.Photo.Sizes[len(a.Photo.Sizes)-1].URL
							break
						}
					}

					time.Sleep(time.Second)

					rand.Seed(time.Now().UnixNano())

					otv, attach, _ := getOtvet(message, posts.Items[0].FromID, vk, false, posts.Items[0].FromID, BotID, helpText, photoURL, 0, true, db)

					data, err := ioutil.ReadFile("blacklist.txt")
					if err != nil {
						panic(err)
					}
					s := string(data)
					ss := strings.Split(s, "\n")
					for i := range ss {
						if ss[i] != "" || ss[i] != " " {
							if strings.Contains(strings.ToLower(message), ss[i]) {
								onAttemptToBlock := "Маму свою заблокируй! Слышал, ты, пузырик пакостный?!"
								onAttemptToBlock, _ = jsonparser.GetString(config, "onAttemptToBlock")
								otv, attach = onAttemptToBlock, ""
							}
						}
					}

					vk.WallCreateComment(api.Params{
						"owner_id":    ownid,
						"post_id":     posts.Items[0].ID,
						"message":     otv,
						"attachments": attach,
						"guid":        rand.Intn(12343242345235),
					})

					//log.Println(resp)

					/*
						_, err := vk.MessagesSend(api.Params{
							"peer_id":    m.PeerID,
							"message":    otv,
							"attachment": attach,
							"random_id":  0,
							"reply_to":   m.MessageID,
						})
					*/
					if err != nil {
						log.Println(err)
					}

				}
				db.Write("last_post_ids", strconv.Itoa(offset), LastPost{ownid, posts.Items[0].ID})
			}
		}
		time.Sleep(time.Second)
	}, "walls")
}
