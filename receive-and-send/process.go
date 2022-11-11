package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/tencent-connect/botgo/dto"
	"github.com/tencent-connect/botgo/dto/message"
	"github.com/tencent-connect/botgo/openapi"
)

// Processor is a struct to process message
type Processor struct {
	api openapi.OpenAPI
}

const (
	CmdWordDragon     = "成语接龙"
	CmdStopWordDragon = "停止接龙"
	CmdExplainWord    = "查看释义"
)
const (
	StopTip           = "欢迎下次使用！"
	ToStartTip        = "接龙还没有开始哦！输入[成语接龙]开始接龙游戏！"
	NotWordTip        = "输入的不是成语哦,再试试吧！"
	NotMatchDragonTip = "输入的成语没有接到上一个成语哦,再试试吧！"
	NormalTip         = "输入[成语接龙]开始游戏！游戏中可回复[查看释义]查看成语含义！"
)

var lastWord string
var play = false

// ProcessMessage is a function to process message
func (p Processor) ProcessMessage(input string, data *dto.WSATMessageData, words map[string]string) error {
	ctx := context.Background()
	cmd := message.ParseCommand(input)

	beginWord := getBeginWord(words)

	switch cmd.Cmd {
	case CmdWordDragon:
		play = true
		p.sendReplyByString(ctx, data.ID, data.ChannelID, beginWord)
		lastWord = beginWord
	case CmdStopWordDragon:
		if play {
			play = false
			p.sendReplyByString(ctx, data.ID, data.ChannelID, StopTip)
		} else {
			p.sendReplyByString(ctx, data.ID, data.ChannelID, ToStartTip)
		}
	default:
		if play {
			if isWordLegal(cmd.Cmd, words) && isWordDragon(cmd.Cmd, lastWord) {
				nextWord := getWord(cmd.Cmd, words)
				p.sendReplyByString(ctx, data.ID, data.ChannelID, nextWord)
				lastWord = nextWord
			} else if cmd.Cmd == CmdExplainWord {
				value := getWordMeaning(lastWord, words)
				p.sendReplyByString(ctx, data.ID, data.ChannelID, value)
			} else if !isWordLegal(cmd.Cmd, words) {
				p.sendReplyByString(ctx, data.ID, data.ChannelID, NotWordTip)
			} else if !isWordDragon(cmd.Cmd, lastWord) {
				p.sendReplyByString(ctx, data.ID, data.ChannelID, NotMatchDragonTip)
			}
		} else {
			p.sendReplyByString(ctx, data.ID, data.ChannelID, NormalTip)
		}
	}
	return nil
}

// ProcessInlineSearch is a function to process inline search
func (p Processor) ProcessInlineSearch(interaction *dto.WSInteractionData) error {
	if interaction.Data.Type != dto.InteractionDataTypeChatSearch {
		return fmt.Errorf("interaction data type not chat search")
	}
	search := &dto.SearchInputResolved{}
	if err := json.Unmarshal(interaction.Data.Resolved, search); err != nil {
		log.Println(err)
		return err
	}
	if search.Keyword != "test" {
		return fmt.Errorf("resolved search key not allowed")
	}
	searchRsp := &dto.SearchRsp{
		Layouts: []dto.SearchLayout{
			{
				LayoutType: 0,
				ActionType: 0,
				Title:      "内联搜索",
				Records: []dto.SearchRecord{
					{
						Cover: "https://pub.idqqimg.com/pc/misc/files/20211208/311cfc87ce394c62b7c9f0508658cf25.png",
						Title: "内联搜索标题",
						Tips:  "内联搜索 tips",
						URL:   "https://www.qq.com",
					},
				},
			},
		},
	}
	body, _ := json.Marshal(searchRsp)
	if err := p.api.PutInteraction(context.Background(), interaction.ID, string(body)); err != nil {
		log.Println("api call putInteractionInlineSearch  error: ", err)
		return err
	}
	return nil
}

func (p Processor) dmHandler(data *dto.WSATMessageData) {
	dm, err := p.api.CreateDirectMessage(
		context.Background(), &dto.DirectMessageToCreate{
			SourceGuildID: data.GuildID,
			RecipientID:   data.Author.ID,
		},
	)
	if err != nil {
		log.Println(err)
		return
	}

	toCreate := &dto.MessageToCreate{
		Content: "默认私信回复",
	}
	_, err = p.api.PostDirectMessage(
		context.Background(), dm, toCreate,
	)
	if err != nil {
		log.Println(err)
		return
	}
}

func genReplyContent(data *dto.WSATMessageData) string {
	var tpl = `你好：%s
在子频道 %s 收到消息。
收到的消息发送时时间为：%s
当前本地时间为：%s

消息来自：%s
`

	msgTime, _ := data.Timestamp.Time()
	return fmt.Sprintf(
		tpl,
		message.MentionUser(data.Author.ID),
		message.MentionChannel(data.ChannelID),
		msgTime, time.Now().Format(time.RFC3339),
		getIP(),
	)
}

func genReplyArk(data *dto.WSATMessageData) *dto.Ark {
	return &dto.Ark{
		TemplateID: 23,
		KV: []*dto.ArkKV{
			{
				Key:   "#DESC#",
				Value: "这是 ark 的描述信息",
			},
			{
				Key:   "#PROMPT#",
				Value: "这是 ark 的摘要信息",
			},
			{
				Key: "#LIST#",
				Obj: []*dto.ArkObj{
					{
						ObjKV: []*dto.ArkObjKV{
							{
								Key:   "desc",
								Value: "这里展示的是 23 号模板",
							},
						},
					},
					{
						ObjKV: []*dto.ArkObjKV{
							{
								Key:   "desc",
								Value: "这是 ark 的列表项名称",
							},
							{
								Key:   "link",
								Value: "https://www.qq.com",
							},
						},
					},
				},
			},
		},
	}
}
