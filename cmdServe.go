package chat

import (
	"context"
	"flag"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

type cmdServe struct{}

func (s *cmdServe) name() string {
	return "serve"
}

func (s *cmdServe) description() string {
	return "serve"
}

func test() {
	NewUser("testuser-1")
	NewUser("testuser-2")
	NewUser("testuser-3")
	chat.RoomList["random"].AddChatMessage("testuser-1", "testuser-1.")
	chat.RoomList["random"].AddChatMessage("testuser-2", "testuser-2.")
	chat.RoomList["random"].AddChatMessage("testuser-3", "testuser-3.")

	log.Printf("%q", chat.UserList)
	log.Printf("%q", chat.RoomList)
	log.Printf("%q", chat.MessageList)
	log.Printf("%q", chat.Config)
}

func (s *cmdServe) run(ctx context.Context, argv []string, outStream io.Writer, errStream io.Writer) error {
	flagSet := flag.NewFlagSet("serve", flag.ContinueOnError)
	flagSet.SetOutput(errStream)

	if err := flagSet.Parse(argv); err != nil {
		log.Println(err)
		return err
	}

	//log.Println("serve : argv %v", argv[0])

	if err := InitializeChat(); err != nil {
		log.Println(err)
		return err
	}

	if err := LoadDataFromFile(); err != nil {
		log.Println(err)
		// 過去に保存したデータがないor読み込めない
		// 初期状態で起動
		// chatを初期化し、botユーザーとrandomルームを作成する。
		CreateDefaultChat()
	}

	if len(argv) > 0 {
		if num, err := strconv.Atoi(argv[0]); err != nil {
			// 引数がおかしい
			log.Println("serverの後ろにはポート番号が指定できます。")
			return err
		} else {
			chat.Config.Port = num
		}
	}

	// chanの保存用goroutine
	ch := make(chan struct{}) // 通知用チャネル
	go DataKeep(ch)

	ch <- struct{}{} // これでchatが保存される。古いファイルは削除される。

	/*
		ki, err := kibela.New(version)
		if err != nil {
			return err
		}
		if flagSet.NArg() < 1 {
			return xerrors.New("usage: kibelasync push [md files]")
		}
		for _, f := range flagSet.Args() {
			md, err := kibela.LoadMD(f)
			if err != nil {
				return err
			}
			if err := ki.PushMD(ctx, md); err != nil {
				return err
			}
		}
	*/

	if err := httpKeep(); err != nil {
		log.Println(err)
		os.Exit(1)
	}
	return nil
}

func InitializeChat() error {
	log.Println("InitializeChat()")
	chat.Config = ProgramConfig{
		KeepAlive: 5 * time.Second,
		Port:      80,
		Host:      "localhost",
		Dir:       ".",
	}
	if p, err := os.Getwd(); err == nil {
		chat.Config.Dir = filepath.Join(p, "chatroom_data")
	}

	if f, err := os.Stat(chat.Config.Dir); !os.IsNotExist(err) && f.IsDir() {
		// ディレクトリがあった
	} else {
		if err := mkdir(chat.Config.Dir); err != nil {
			return err
		}
	}

	return nil
}

func CreateDefaultChat() {
	log.Println("add NewUser bot")
	// 初期の部屋を作成 botユーザーの作成
	if err := NewUser("bot"); err != nil {
		// これはあり得ないと思いたい
		log.Println(err)
	}
	log.Println("add NewRoom random")
	if err := NewChatRoom("random"); err != nil {
		chat.RoomList["random"].AddChatMessage("bot", "チャットの開始")
	}
}
