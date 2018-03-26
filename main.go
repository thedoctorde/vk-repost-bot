package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/thedoctorde/vk-api"
	"github.com/thedoctorde/vk-repost-bot/tg"
	"github.com/thedoctorde/vk-repost-bot/vk"
	"log"
	"os"
	"sync"
	"time"
)

type History struct {
	lock    sync.Mutex
	storage map[int64]map[int64]Post
}

type Post struct {
	ID   int64
	Text string
	Date int64
}

func NewHistory() *History {
	return &History{
		storage: make(map[int64]map[int64]Post, 0),
	}
}

func (h *History) Add(groupId, postId int64, text string, date int64) error {
	h.lock.Lock()
	defer h.lock.Unlock()
	if _, ok := h.storage[groupId]; !ok {
		h.storage[groupId] = make(map[int64]Post, 0)
	}
	h.storage[groupId][postId] = Post{
		Text: text,
		ID:   postId,
		Date: date,
	}
	return nil
}

func main() {

	err := godotenv.Load("settings.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	tokenTg := os.Getenv("TOKEN_TG")
	tgChannelName := os.Getenv("CHANNEL_NAME")
	tokenVk := os.Getenv("TOKEN_VK")

	vkManager, err := vk.NewManager(tokenVk)
	if err != nil {
		fmt.Println(err)
	}
	groups := []int64{
		102325800,
		62262331,
		90358638,
		76800969,
		62914251,
		134618886,
		86264049,
		119563192,
		73493509,
		149861205,
		163995482,
	}
	vkManager.FillGroups(groups)

	history := NewHistory()

	b := tg.NewBot(tokenTg)
	for {
		unixNow := time.Now().Unix()
		for key := range vkManager.Groups {
			dst := vkapi.Destination{
				GroupID: key,
			}
			_, items, _, _, errReq := vkManager.Client.GetWall(dst, 4, 0, "", false)
			if errReq != nil {
				fmt.Println(errReq)
			} else {
				for _, item := range items {
					fmt.Println(item)
					if (unixNow - item.Date) < 300 {
						history.Add(item.OwnerID, item.ID, item.Text, item.Date)
						err = b.SendMessage(tgChannelName, item.Text)
						if err != nil {
							fmt.Println(err)
						}
					}
				}
			}
		}
		time.Sleep(5 * time.Minute)
	}

	wg := sync.WaitGroup{}
	wg.Add(1)
	wg.Wait()

}
