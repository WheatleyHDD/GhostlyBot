package commands

import (
	//"fmt"

	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/SevereCloud/vksdk/api"

	"fmt"
	"image"
	"image/color"
	_ "image/jpeg"
	_ "image/png"
	"math"
	"os"

	"github.com/adam-lavrik/go-imath/ix"
	"github.com/buger/jsonparser"
	"github.com/fogleman/gg"

	. "vkbot/utils"
)

func clamp(value, maxv, minv int) int {
	return ix.Max(ix.Min(value, maxv), minv)
}

func FindGif(q string, vk *api.VK, height int) (otvet string, attachments string, similar float64) {
	config := ConfigParse()
	if len(q) <= 1 {
		q = "gif"
	}
	//fmt.Println(q)
	dcs, err := vk.DocsSearch(api.Params{
		"q": q,
	})
	if err != nil {
		return "Ошиб очка: " + err.Error(), "", 0
	}
	rand.Seed(time.Now().UnixNano())
	roffset := rand.Intn(clamp(dcs.Count-10, dcs.Count, 1))
	docs, err := vk.DocsSearch(api.Params{
		"q":      q,
		"count":  1000,
		"offset": roffset,
	})
	if err != nil {
		return "Ошиб очка: " + err.Error(), "", 0
	}
	gif := ""
	cnt := 0
	for _, doc := range docs.Items {
		if cnt == height {
			break
		}
		if doc.Ext == "gif" {
			gif = strings.Join([]string{gif, "doc", strconv.Itoa(doc.OwnerID), "_", strconv.Itoa(doc.ID), ","}, "")
			cnt = cnt + 1
		}
	}
	//fmt.Println(gif)
	if gif != "" {
		answer, _ := jsonparser.GetString(config, "message_list", "find", "gif", "found")
		return answer, gif, 0
	}
	answer, _ := jsonparser.GetString(config, "message_list", "find", "gif", "notFound")
	return answer, "", 0
}

func FindPhoto(q string, vk *api.VK, height int) (otvet string, attachments string, similar float64) {
	config := ConfigParse()
	photos, err := vk.PhotosSearch(api.Params{
		"q":     q,
		"count": 1000,
		"sort":  0,
	})
	if err != nil {
		return "Ошиб очка: " + err.Error(), "", 0
	}
	if photos.Count >= 1 {
		if len(photos.Items) > 0 {
			roffset := 0

			if photos.Count > height {
				rand.Seed(time.Now().UnixNano())
				roffset = rand.Intn(clamp(photos.Count-10, len(photos.Items)-12, 1))
			}

			//log.Print(photos)
			attch := ""
			//cnt := 0
			if photos.Count < height {
				for i := 0; i < photos.Count; i++ {

					attch = strings.Join([]string{attch, "photo", strconv.Itoa(photos.Items[i].OwnerID), "_", strconv.Itoa(photos.Items[i].ID), ","}, "")
				}
			} else {
				for i := roffset; i < roffset+10; i++ {

					attch = strings.Join([]string{attch, "photo", strconv.Itoa(photos.Items[i].OwnerID), "_", strconv.Itoa(photos.Items[i].ID), ","}, "")
				}
			}
			//fmt.Println(attch)
			answer, _ := jsonparser.GetString(config, "message_list", "find", "photo", "found")
			return answer, attch, 0
		}
		return "Братуха, все круто, но сейчас вк на меня накинули лимит, из-за которого я не могу отправить тебе фото. Попробуй еще раз через некоторое время.", "", 0
	}
	answer, _ := jsonparser.GetString(config, "message_list", "find", "photo", "notFound")
	return answer, "", 0
}

func FindVideo(q string, vk *api.VK, height int) (otvet string, attachments string, similar float64) {
	config := ConfigParse()
	if q == "" {
		answer, _ := jsonparser.GetString(config, "message_list", "find", "video", "queryMissing")
		return answer, "", 0
	}

	videos, err := vk.VideoSearch(api.Params{
		"q":     q,
		"count": 200,
		"adult": 1,
	})
	if err != nil {
		return "Ошиб очка: " + err.Error(), "", 0
	}
	if videos.Count >= 1 {
		roffset := 0

		if videos.Count > 10 {
			rand.Seed(time.Now().UnixNano())
			roffset = rand.Intn(clamp(videos.Count-height, len(videos.Items)-6, 1))
		}

		//log.Print(photos)
		attch := ""
		//cnt := 0
		if videos.Count < height {
			for i := 0; i < videos.Count; i++ {
				attch = strings.Join([]string{attch, "video", strconv.Itoa(videos.Items[i].OwnerID), "_", strconv.Itoa(videos.Items[i].ID), ","}, "")
			}
		} else {
			for i := roffset; i < roffset+height; i++ {
				attch = strings.Join([]string{attch, "video", strconv.Itoa(videos.Items[i].OwnerID), "_", strconv.Itoa(videos.Items[i].ID), ","}, "")
			}
		}
		//fmt.Println(attch)
		answer, _ := jsonparser.GetString(config, "message_list", "find", "video", "found")
		return answer, attch, 0
	}
	answer, _ := jsonparser.GetString(config, "message_list", "find", "video", "notFound")
	return answer, "", 0
}

func CreateLobsterMeme(label, path, saveFile string, ChatID int, vk *api.VK) string {

	width, height := getImageDimension(path)
	dc := gg.NewContext(width, height)
	dc.SetRGB(1, 1, 1)
	dc.Clear()
	im, err := gg.LoadJPG(path)
	if err != nil {
		return "Ошиб очка: " + err.Error()
	}

	dc.DrawImage(im, 0, 0)

	grad := gg.NewLinearGradient(250, float64(height), 250, math.Round(float64(height*35/40)))
	grad.AddColorStop(0, color.RGBA{0, 0, 0, 75})
	grad.AddColorStop(1, color.RGBA{0, 0, 0, 0})

	dc.SetFillStyle(grad)
	dc.DrawRectangle(0, math.Round(float64(height*35/40)), float64(width), float64(height*5/40))
	dc.Fill()

	fontsize := float64(ix.Min(height*36/400, width*36/400))

	dc.SetRGB(0, 0, 0)
	if err := dc.LoadFontFace("lobster.ttf", fontsize); err != nil {
		return "Ошиб очка: " + err.Error()
	}
	dc.DrawStringAnchored(label, float64(width/2)-1, float64(height)-fontsize+1, 0.5, 0.9)

	dc.SetRGB(1, 1, 1)
	if err := dc.LoadFontFace("lobster.ttf", fontsize); err != nil {
		return "Ошиб очка: " + err.Error()
	}
	dc.DrawStringAnchored(label, float64(width/2), float64(height)-fontsize, 0.5, 0.9)
	dc.SavePNG(saveFile)
	_ = os.Remove(path)

	meme, err := os.Open(saveFile)
	if err != nil {
		return "Ошиб очка: " + err.Error()
	}

	photosPhoto, err := vk.UploadMessagesPhoto(2000000000+ChatID, meme)

	_ = os.Remove(saveFile)

	return "photo" + strconv.Itoa(photosPhoto[0].OwnerID) + "_" + strconv.Itoa(photosPhoto[0].ID)
}

func getImageDimension(imagePath string) (int, int) {
	file, err := os.Open(imagePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
	}

	image, _, err := image.DecodeConfig(file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %v\n", imagePath, err)
	}
	return image.Width, image.Height
}
