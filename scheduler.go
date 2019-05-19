package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/jasonlvhit/gocron"
	"log"
	"net/http"
	"time"
)

const itcUaURL = "https://itc.ua/"
const newsUpdateTimeMinute = 2

var lastReceivedITCURL string

// Scheduler is a scheduler service struct
type Scheduler struct {
	chatService Chat
	bot         Bot
}

// NewScheduler returns new Bot service
func NewScheduler(b Bot, c Chat) Scheduler {
	return Scheduler{
		chatService: c,
		bot:         b,
	}
}

// StartBotScheduler starts bot Scheduler
func (s *Scheduler) StartBotScheduler() {
	time.Sleep(time.Second*5)

	sc := gocron.NewScheduler()
	sc.Every(newsUpdateTimeMinute).Minutes().Do(s.UpdateNews)

	// starts scheduler
	<-sc.Start()
}

// UpdateNews fetches all updates
func (s *Scheduler) UpdateNews() {
	// ITC.ua block start
	url, err := s.LastITCUAPostURL()
	if err != nil {
		log.Printf("Error while retreiving ITC post %v", err)
		return
	}

	if lastReceivedITCURL == url || url == "" {
		return
	}

	err = s.bot.SendToChats(s.chatService.GetChatIDs(), url)
	if err != nil {
		log.Printf("Error while sending to chats %v", err)
	}
	lastReceivedITCURL = url
	// ITC.ua block end
}

// LastITCUAPostURL parses last post url from ITC.ua
func (s *Scheduler) LastITCUAPostURL() (firstPostURL string, err error) {
	res, err := http.Get(itcUaURL)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return "", fmt.Errorf("got %v", res.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return "", fmt.Errorf("got %v", res.StatusCode)
	}

	// To find all hrefs
	//doc.Find("div.row").Each(func(i int, s *goquery.Selection) {
	//s.Find("div.post").Each(func(i int, s *goquery.Selection) {
	//	str,_:= s.Find("div.col-xs-8").Find("a").Attr("href")
	//	fmt.Println(str)
	//})
	//})

	firstPostURL, found := doc.Find("div.row").
		Find("div.post").
		Find("div.col-xs-8").
		Find("a").
		Attr("href")

	if !found {
		return "", fmt.Errorf("received no url")
	}
	return
}
