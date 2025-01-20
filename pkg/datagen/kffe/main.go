package kffe

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/go-resty/resty/v2"
	uuid2 "github.com/google/uuid"

	"github.com/smart-fm/kf-api/domain/repository"
	"github.com/smart-fm/kf-api/endpoints/common"
	"github.com/smart-fm/kf-api/infrastructure/mysql/dao"
)

// const dev = `http://localhost:8081`
// const code = `/s/EkCLyM/BJfmus/ak8BXI.html` // code.
// const cid = `TM-J9pWlL8GfI`

type GenRequest struct {
	Host   string `json:"url"`
	QRCode string `json:"code"`
	CardId string `json:"card_id"`
}

func Gen(request GenRequest) {

	r := resty.New()
	// 1. 生成30个用户.
	var uuids []string
	for i := 0; i < 30; i++ {
		req := r.NewRequest().SetBody(
			map[string]string{
				"code": request.QRCode,
			},
		).SetHeader("Content-Type", "application/json")

		rsp, err := req.Post(request.Host + "/api/kf-fe/qrcode/scan")
		if err != nil {
			log.Fatal(err)
			return
		}
		var scanResp ScanResponse
		err = json.Unmarshal(rsp.Body(), &scanResp)
		if err != nil {
			log.Fatal(err)
			return
		}
		if scanResp.Code != 200 {
			log.Fatal(scanResp.Message)
			return
		}

		uuid := scanResp.Data.UserInfo.Uuid
		uuids = append(uuids, uuid)
	}
	fmt.Println("--->", len(uuids))
	// 2. 生成消息.
	for _, id := range uuids {
		generateMessageOneUser(id, request.CardId)
	}
}

func generateMessageOneUser(uuid string, cardId string) {
	r := repository.KFMessageRepository{}
	for i := 0; i < 100; i++ {
		msgType := common.MessageType(
			gofakeit.RandomString(
				[]string{
					string(common.MessageTypeText),
					string(common.MessageTypeImage),
					string(common.MessageTypeVideo),
				},
			),
		)
		msg := dao.KFMessage{
			MsgId:   uuid2.NewString(),
			MsgType: msgType,
			GuestId: uuid,
			CardId:  cardId,
			Content: randContent(msgType),
			IsKf:    gofakeit.RandomInt([]int{1, 2}),
		}
		err := r.SaveOne(context.Background(), &msg)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func randContent(mt common.MessageType) string {
	var content string
	switch mt {
	case common.MessageTypeText:
		content = fmt.Sprintf("这是一段长文本,这是一段长文本,这是一段长文本,这是一段长文本,这是一段长文本")
	case common.MessageTypeImage:
		content = "/text.png"
	case common.MessageTypeVideo:
		content = "/QQ20250120-134814-HD.mp4"
	}
	return content
}

type ScanResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		UserInfo struct {
			Uuid     string `json:"uuid"`
			Avatar   string `json:"avatar"`
			NickName string `json:"nickName"`
		} `json:"userInfo"`
		IsNewUser bool `json:"isNewUser"`
	} `json:"data"`
}
