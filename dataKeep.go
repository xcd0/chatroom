package chat

import (
	"encoding/gob"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
	"time"
)

// DataKeepはプログラムの起動と同時にgoroutineとして起動して
// チャネルに通知があったときにデータを保存します。
// 必ずgo DataKeep(ch)のように実行すること。
func DataKeep(ch <-chan struct{}) {

	for {
		<-ch // 送信されてきたらchatを保存する
		filename := filepath.Join(chat.Config.Dir, fmt.Sprintf("chat_%s.gob", time.Now().Format("20060102150405")))
		if err := encodeChatWithGob(filename); err != nil {
			log.Println(err) // 終了はしない
		}
		// 古いファイルを削除
		if err := RemoveOldFiles(5); err != nil {
			fmt.Println(err)
		}
	}

}

func GenFileTimestampList() []string {
	// ファイル名から保存日時を取得してソートする
	var timestamps []string
	files, err := ioutil.ReadDir(chat.Config.Dir)
	if err != nil {
		log.Println(err)
		panic(err)
	}
	for _, f := range files {
		if !f.IsDir() {
			name := f.Name()
			e := filepath.Ext(name)
			if name[:5] == "chat_" && e == ".gob" {
				//log.Println("file : ", name)
				timestamps = append(timestamps, name[5:len(name)-4])
			}
		}
	}
	sort.Strings(timestamps)
	return timestamps
}

func LoadDataFromFile() error {
	//log.Printf("dir: %q", chat.Config.Dir)
	timestamps := GenFileTimestampList()
	//log.Println("dir : ", files)

	//log.Println("timestamps : ", timestamps)

	// 最新のファイルを読み込む
	var latestFilename string
	if len(timestamps) > 0 {
		latestTimestamp := timestamps[len(timestamps)-1]
		latestFilename = filepath.Join(chat.Config.Dir, fmt.Sprintf("chat_%s.gob", latestTimestamp))
	} else {
		err := errors.New(fmt.Sprintf("no data files found in : %s", latestFilename))
		log.Println(err)
		panic(err)
		return err
	}
	log.Println("load file : %s", latestFilename)
	return decodeChatWithGob(latestFilename)
}

func RemoveOldFiles(numToKeep int) error {
	timestamps := GenFileTimestampList()

	// 最新のファイルを除く古いファイルを削除
	numFiles := len(timestamps)
	if numFiles <= numToKeep {
		return nil
	}
	for _, ts := range timestamps[:numFiles-numToKeep] {
		filename := filepath.Join(chat.Config.Dir, fmt.Sprintf("chat_%s.gob", ts))
		err := os.Remove(filename)
		log.Println("old data file : ", filename, " removed.")
		if err != nil {
			log.Println(err)
			return err
		}
	}
	return nil
}

func encodeChatWithGob(path string) error {
	f, err := os.Create(path)
	defer f.Close()
	if err != nil {
		log.Println(err)
		return err
	}
	return gob.NewEncoder(f).Encode(chat)
}

func decodeChatWithGob(path string) error {
	f, err := os.Open(path)
	defer f.Close()
	if err != nil {
		log.Println(err)
		return err
	}
	return gob.NewDecoder(f).Decode(chat)
}
