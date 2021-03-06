package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"
	"syscall"

	"github.com/judwhite/go-svc/svc"
	"github.com/pangliang/MirServer-Go/gameserver"
	"github.com/pangliang/MirServer-Go/loginserver"
	"github.com/pangliang/MirServer-Go/tools"
)

type program struct {
	loginServer *loginserver.LoginServer
	gameServer  *gameserver.GameServer
}

func main() {
	prg := &program{}
	if err := svc.Run(prg, syscall.SIGINT, syscall.SIGTERM); err != nil {
		log.Fatal(err)
	}
}

func (p *program) Init(env svc.Environment) error {
	if env.IsWindowsService() {
		dir := filepath.Dir(os.Args[0])
		return os.Chdir(dir)
	}
	return nil
}

func (p *program) Start() error {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.SetOutput(os.Stdout)

	tools.MigrateDevDB()

	loginOpt := &loginserver.Option{}
	flagSet := flag.NewFlagSet("loginserver", flag.ExitOnError)
	flagSet.BoolVar(&loginOpt.IsTest, "test.v", false, "")
	flagSet.StringVar(&loginOpt.Address, "login-address", "0.0.0.0:7000", "<addr>:<port> to listen on for TCP clients")
	flagSet.StringVar(&loginOpt.DataSourceName, "dbSource", "./mir2.db", "DataSourceName")
	flagSet.StringVar(&loginOpt.DriverName, "dbDriver", "sqlite3", "database DriverName")
	flagSet.Parse(os.Args[1:])

	p.loginServer = loginserver.New(loginOpt)
	p.loginServer.Main()

	gameOpt := &gameserver.Option{}
	flagSet = flag.NewFlagSet("gameserver", flag.ExitOnError)
	flagSet.BoolVar(&gameOpt.IsTest, "test.v", false, "")
	flagSet.StringVar(&gameOpt.Address, "game-address", "0.0.0.0:7400", "<addr>:<port> to listen on for TCP clients")
	flagSet.StringVar(&gameOpt.DataSourceName, "dbSource", "./mir2.db", "DataSourceName")
	flagSet.StringVar(&gameOpt.DriverName, "dbDriver", "sqlite3", "database DriverName")
	flagSet.Parse(os.Args[1:])

	p.gameServer = gameserver.New(gameOpt)
	p.gameServer.Main()

	return nil
}

func (p *program) Stop() error {
	if p.loginServer != nil {
		p.loginServer.Exit()
	}

	if p.gameServer != nil {
		p.gameServer.Exit()
	}
	return nil
}
