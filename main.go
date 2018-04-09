package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/thedoctorde/vk-api"
	"github.com/thedoctorde/vk-repost-bot/tg"
	"github.com/thedoctorde/vk-repost-bot/vk"
	"log"
	"os"
	"strconv"
	"strings"
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

func (h *History) Has(groupId, postId int64) bool {
	h.lock.Lock()
	defer h.lock.Unlock()

	if _, ok := h.storage[groupId]; ok {
		if _, ok2 := h.storage[groupId][postId]; ok2 {
			return true
		}
	}
	return false
}

func main() {
	//var postedRecordIds = map[string]struct{}{}
	err := godotenv.Load("settings.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	tokenTg := os.Getenv("TOKEN_TG")
	tgChannelName := os.Getenv("CHANNEL_NAME")
	tokenVk := os.Getenv("TOKEN_VK")
	groupsStr := os.Getenv("GROUPS")
	updateStr := os.Getenv("UPDATE")
	update, _ := strconv.Atoi(updateStr)
	updatePeriod := int64(update)
	vkManager, err := vk.NewManager(tokenVk)
	if err != nil {
		fmt.Println(err)
	}
	var groups []int64
	a := strings.Split(groupsStr, ",")
	for i := range a {
		n, err := strconv.Atoi(a[i])
		if err != nil {
			continue
		}
		num := int64(n)
		groups = append(groups, num)
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
			_, wallPosts, _, _, errReq := vkManager.Client.GetWall(dst, 10, 0, "", false)
			if errReq != nil {
				fmt.Println(errReq)
				continue
			}
			for _, wallPost := range wallPosts {
				if (unixNow - wallPost.Date) > 1200 {
					continue
				}
				if history.Has(wallPost.OwnerID, wallPost.ID) {
					continue
				}
				history.Add(wallPost.OwnerID, wallPost.ID, wallPost.Text, wallPost.Date)
				err = b.SendMessage(tgChannelName, wallPost.Text)
				if err != nil {
					runes := []rune(wallPost.Text)
					l := len(runes)
					var runes1 []rune
					var runes2 []rune
					for i := range runes {
						if i < l/2 {
							runes1 = append(runes1, runes[i])
						} else {
							runes2 = append(runes2, runes[i])
						}
					}
					err = b.SendMessage(tgChannelName, string(runes1))
					err = b.SendMessage(tgChannelName, string(runes2))
					fmt.Println(err)
				}

			}
		}
		time.Sleep(time.Duration(updatePeriod) * time.Second)
	}

	wg := sync.WaitGroup{}
	wg.Add(1)
	wg.Wait()

}
