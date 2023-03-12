package chat

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"strings"
	"time"
)

func genUrl(host string, port int) (string, error) {
	return fmt.Sprintf("%s:%d", host, port), nil
}

func httpKeep() error {

	url, err := genUrl(chat.Config.Host, chat.Config.Port)
	if err != nil {
		log.Println(err)
		return err
	}

	log.Println(url)

	listener, err := net.Listen("tcp", url)
	if err != nil {
		log.Println(err)
		return err
	}
	fmt.Println("Server is running at ", url)

	errChan := make(chan error, 1)
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println(err)
			return err
		}

		go func() {
			defer conn.Close()
			fmt.Printf("Accept %v\n", conn.RemoteAddr()) // Accept後のソケットで何度も応答を返すためにループ
			for {
				// タイムアウトを設定
				conn.SetReadDeadline(time.Now().Add(chat.Config.KeepAlive)) // リクエストを読み込む
				request, err := http.ReadRequest(bufio.NewReader(conn))
				if err != nil {
					// タイムアウトもしくはソケットクローズ時は終了
					// それ以外はエラーにする
					neterr, ok := err.(net.Error) // ダウンキャスト
					if ok && neterr.Timeout() {
						fmt.Println("Timeout")
						break
					} else if err == io.EOF {
						break
					}
					log.Println(err)
					errChan <- err
					return
				}
				// リクエストを表示
				dump, err := httputil.DumpRequest(request, true)
				if err != nil {
					log.Println(err)
					errChan <- err
				}
				fmt.Println(string(dump))
				content := "Hello World\n"
				// レスポンスを書き込む
				// HTTP/1.1かつ、ContentLengthの設定が必要
				response := http.Response{
					StatusCode:    200,
					ProtoMajor:    1,
					ProtoMinor:    1,
					ContentLength: int64(len(content)), Body: io.NopCloser(strings.NewReader(content)),
				}
				response.Write(conn)
			}
			errChan <- nil
		}()
	}
}
