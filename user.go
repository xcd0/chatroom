package chat

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"net/url"
	"strings"
	"time"

	"github.com/google/uuid"
)

type User struct {
	UserID     string    `json:"user_id"`     // ユーザー固有のID
	Name       string    `json:"name"`        //
	Avatar     string    `json:"avatar"`      //
	LastActive time.Time `json:"last_active"` //
	Posts      []string  `json:"posts"`       // 投稿リスト ChatMessageのPostID
}

func NewUser(name string) error {
	if _, exist := chat.UserList[name]; exist {
		err := errors.New("同じ名前のユーザーがいます。違う名前にして下さい。")
		log.Println(err)
		// TODO: ユーザーに名前を変えてもらう処理を別のところで実装
		return err
	}

	chat.UserList[name] = &User{
		UserID:     uuid.NewString(),
		Name:       name,
		LastActive: time.Now(),
	}
	chat.UserList[name].Avatar = genIcon(fmt.Sprintf("%s%s", chat.UserList[name].Name, chat.UserList[name].UserID))
	return nil
}

const (
	gravatarURL = "https://www.gravatar.com/avatar/"
	iconSize    = 256 // アバター画像のサイズ
)

func genIcon(uniqueString string) string {
	hash := md5.Sum([]byte(strings.ToLower(uniqueString)))
	md5Hash := hex.EncodeToString(hash[:])
	return fmt.Sprintf("%s%s?d=identicon&s=%d", gravatarURL, url.PathEscape(md5Hash), iconSize) // Gravatar URLの生成
}
