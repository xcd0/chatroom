package chat

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"

	"golang.org/x/xerrors"
)

type ProgramConfig struct {
	KeepAlive time.Duration
	Port      int
	Host      string
	Dir       string // 保持するデータを置くディレクトリパス
}

type Chat struct {
	UserList    map[string]*User        // name -> struct
	RoomList    map[string]*ChatRoom    // roomName -> struct
	MessageList map[string]*ChatMessage // PostID -> struct
	Config      ProgramConfig           // 設定
}

type runner interface {
	name() string
	description() string
	run(context.Context, []string, io.Writer, io.Writer) error
}

var (
	chat = &Chat{ // すべてのデータを保持
		UserList:    map[string]*User{},        // name -> struct
		RoomList:    map[string]*ChatRoom{},    // roomName -> struct
		MessageList: map[string]*ChatMessage{}, // PostID -> struct
		Config:      ProgramConfig{},           // 設定
	}
	cmdName     = ""
	subCommands = []runner{
		&cmdServe{},
		//&cmdPull{},
		//&cmdPush{},
	}
	dispatch          = make(map[string]runner, len(subCommands))
	maxSubcommandName int
)

func init() {
	for _, r := range subCommands {
		n := r.name()
		l := len(n)
		if l > maxSubcommandName {
			maxSubcommandName = l
		}
		dispatch[n] = r
	}
	cmdName = filepath.Base(os.Args[0])
}

func Run(argv []string, outStream, errStream io.Writer) error {
	log.SetOutput(errStream)
	log.SetPrefix(fmt.Sprintf("[%s] ", cmdName))

	if !debug {
		log.SetFlags(log.Ltime | log.Lshortfile)
		log.Println("Debug mode on")
	} else {
		log.SetFlags(0)
	}

	nameAndVer := fmt.Sprintf("%s (v%s rev:%s)", cmdName, version, revision)

	flagSet := flag.NewFlagSet(nameAndVer, flag.ContinueOnError)
	flagSet.SetOutput(errStream)
	flagSet.Usage = func() {
		fmt.Fprintf(flagSet.Output(), "Usage of %s:\n", nameAndVer)
		flagSet.PrintDefaults()
		fmt.Fprintf(flagSet.Output(), "\nCommands:\n")
		formatCommands(flagSet.Output())
	}

	ver := flagSet.Bool("version", false, "display version")
	if err := flagSet.Parse(argv); err != nil {
		return err
	}
	if *ver {
		return printVersion(outStream)
	}

	argv = flagSet.Args() // 上書き
	if len(argv) < 1 {
		return xerrors.New("no subcommand specified")
	}
	r, ok := dispatch[argv[0]]
	if !ok {
		return xerrors.Errorf("unknown subcommand: %s", argv[0])
	}

	if err := r.run(context.Background(), argv[1:], outStream, errStream); err != nil {
		return err
	}
	return nil
}

func printVersion(out io.Writer) error {
	_, err := fmt.Fprintf(out, "%s v%s (rev:%s)\n", cmdName, version, revision)
	return err
}

func formatCommands(out io.Writer) {
	format := fmt.Sprintf("    %%-%ds  %%s\n", maxSubcommandName)
	for _, r := range subCommands {
		fmt.Fprintf(out, format, r.name(), r.description())
	}
}
