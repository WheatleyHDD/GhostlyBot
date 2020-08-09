package main

import (
	//"fmt"
	"io/ioutil"
	"log"
	"math/rand"

	"github.com/SevereCloud/vksdk/api"
	_ "github.com/adam-lavrik/go-imath/ix"
	"strings"
	"time"

	"github.com/masatana/go-textdistance"

	"github.com/buger/jsonparser"
	scribble "github.com/nanobox-io/golang-scribble"

	. "vkbot/commands"
	. "vkbot/utils"
)

func similarity(s1, s2 string) (similarity float64) {
	normal1, normal2 := strings.ToLower(s1), strings.ToLower(s2)
	return textdistance.JaroWinklerDistance(normal1, normal2)
}

func getRandAnswers(k string) []string {
	data, err := ioutil.ReadFile("PhasesDB.txt")
	if err != nil {
		log.Fatal(err)
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

func getAnsw(a string, userID int, vk *api.VK) (string, string, float64) {
	var lastSimilarity float64
	//var lastPhase [3]string
	data, err := ioutil.ReadFile("PhasesDB.txt")
	if err != nil {
		return "Ошиб очка: " + err.Error(), "", 0
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
	wreg := GetRegularData(gra[1], userID, vk)
	attwreg := GetRegularData(gra[2], userID, vk)
	return strings.Replace(wreg, "-n", "\n", -1), attwreg, lastSimilarity
}

func getOtvet(a string, userID int, vk *api.VK, isChat bool, ChatID, BotID int, helpText string, attachment string, MID int, onWall bool, db *scribble.Driver) (otvet string, attachments string, similar float64) {
	config := ConfigParse()
	alist := strings.Split(a, " ")
	//fmt.Println(alist)
	switch strings.ToLower(alist[0]) {
	case "найди":
		if len(alist) > 1 {
			ftype := alist[1]
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
				StoreStatisticFindQue(que, db)
			} else {
				StoreStatisticFindQue("empty", db)
			}
			return Find(ftype, que, onWall, db, vk)
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
			return Who(onWall, isChat, ChatID, userID, BotID, db, vk)
		}
	case "команды", "помощь":

		StoreStatistic("команды", db)

		return strings.Replace(helpText, "-n", "\n", -1), "", 0
	case "лобстер", "lobster":
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

		return Lobster(label, attachment, ChatID, MID, userID, onWall, db, vk)

	case "фас":

		StoreStatistic("фас", db)

		return "Гав-гав", "", 0
	case "выбери", "выбери:":
		return Choose(a, userID, db, vk)
	case "рулеточка", "roll":
		return Roll(alist, userID, db, vk)
	case "инфа":
		return Infa(alist, userID, db, vk)
	case "когда":
		return When(alist, userID, db, vk)
	case "о":
		if len(alist) == 2 || strings.ToLower(alist[1]) == "боте" {

			StoreStatistic("о боте", db)

			answer, _ := jsonparser.GetString(config, "about")
			answer = GetRegularData(answer, userID, vk)
			return answer, "", 0

		}
	case "ген", "gen":
		Gen(alist, userID, ChatID, db, vk)
	}
	reta, retattach, _ := getAnsw(a, userID, vk)
	return reta, retattach, 0
}
