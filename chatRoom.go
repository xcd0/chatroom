package chat

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/google/uuid"
	"golang.org/x/exp/slices"
)

type ChatRoom struct {
	RoomID   string   `json:"room_id"`  // チャットルーム固有のUUID
	RoomName string   `json:"name"`     // チャットルームにつける名前
	Users    []string `json:"users"`    // チャットルームに参加しているユーザー UserのUserID
	Messages []string `json:"messages"` // チャットメッセージリスト ChatMessageのPostID
	Dir      string   `json:"dir"`      // チャットルームの情報を保持するディレクトリ名
}

func (r *ChatRoom) AddChatMessage(userName string, content string) *ChatMessage {
	uid := chat.UserList[userName].UserID
	if exist := slices.Contains(r.Users, uid); !exist {
		r.Users = append(r.Users, uid)
	}
	return newChatMessage(r.RoomID, uid, content)
}

func mkdir(dir string) error {
	if f, err := os.Stat(dir); !os.IsNotExist(err) && f.IsDir() {
		// ディレクトリがあった
	} else {
		if err := os.Mkdir(dir, 0755); err != nil {
			log.Println(err) // ディレクトリが作成できなかった
			return err
		}
	}
	return nil
}

func NewChatRoom(roomName string) error {
	if _, exist := chat.RoomList[roomName]; exist {
		err := errors.New(fmt.Sprintf("同じ名前のチャットルームがすでにあります。違う名前にして下さい。: name : %s", roomName))
		log.Println(err)
		return err
	}

	chat.RoomList[roomName] = &ChatRoom{
		RoomID:   uuid.NewString(),
		RoomName: roomName,
		Users:    []string{},
		Messages: []string{},
	}

	if err := mkdir(chat.Config.Dir); err != nil {
		log.Println(err) // ディレクトリが作成できなかった
		return err
	}

	chat.RoomList[roomName].Dir = filepath.Join(chat.Config.Dir, fmt.Sprintf("%s_%s", chat.RoomList[roomName].RoomName, chat.RoomList[roomName].RoomID))
	if err := mkdir(chat.RoomList[roomName].Dir); err != nil {
		log.Println(err) // ディレクトリが作成できなかった
		return err
	}

	return nil
}
