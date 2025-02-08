package matrix

import (
	"maunium.net/go/mautrix"
)

func SendMessage(client *mautrix.Client, roomId string, message string) (err error) {
	_, err = client.SendMessageEvent(roomId, "m.room.message", mautrix.Content{
		MsgType: "m.text",
		Body:    message,
	})
	return
}

func EditMessage(client *mautrix.Client, roomId string, message string, eventId string) (err error) {
	_, err = client.SendMessageEvent(roomId, "m.room.message", mautrix.Content{
		MsgType: "m.text",
		Body:    message,
		Format:  "org.matrix.custom.html",
		RelatesTo: &mautrix.RelatesTo{
			RelType: "m.replace",
			EventID: eventId,
		},
	})
	return
}

func ReplyMessage(client *mautrix.Client, roomId string, message string, eventId string) (err error) {
	_, err = client.SendMessageEvent(roomId, "m.room.message", mautrix.Content{
		MsgType: "m.text",
		Body:    message,
		Format:  "org.matrix.custom.html",
		RelatesTo: &mautrix.RelatesTo{
			RelType: "m.in_reply_to",
			EventID: eventId,
		},
	})
	return
}

func UploadFile(client *mautrix.Client, roomId string, file string) (err error) {
	_, err = client.Upload(roomId, file)
	return
}
