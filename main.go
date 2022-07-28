package main

import (
	"flag"
	"log"
	"os"
	"os/exec"
	"strings"

	socks5 "github.com/szampardi/ss5/socks5"
)

var conf struct {
	Address  string
	Port     string
	User     string
	Password string
	Daemon   bool
}

func init() {
	flag.StringVar(&conf.Address, "a", "", "listen address (default ::)")
	flag.StringVar(&conf.Port, "p", "1080", "listen port (default 1080)")
	flag.StringVar(&conf.User, "u", "", "proxy user")
	flag.StringVar(&conf.Password, "x", "", "proxy password")
	flag.BoolVar(&conf.Daemon, "d", false, "run as daemon")

	flag.Parse()

	if conf.Address == "" {
		conf.Address, _ = os.LookupEnv("PROXY_ADDR")
	}
	if conf.Port == "1080" || conf.Port == "" {
		tp, set := os.LookupEnv("PROXY_PORT")
		if set {
			conf.Port = tp
		}
	}

	if conf.User == "" {
		conf.User, _ = os.LookupEnv("PROXY_USER")
	}
	if conf.Password == "" {
		conf.Password, _ = os.LookupEnv("PROXY_PASS")
	}

}

func main() {

	socks5conf := &socks5.Config{
		Logger: log.New(os.Stdout, "\t", log.LstdFlags),
	}

	if conf.User != "" && conf.Password != "" {
		c := socks5.StaticCredentials{conf.User: conf.Password}
		a := socks5.UserPassAuthenticator{Credentials: c}
		socks5conf.AuthMethods = []socks5.Authenticator{a}
	} else {
		log.Print("Running without user/password protection!")
	}

	if conf.Daemon {
		params := os.Args[1:]
		i := 0
		for ; i < len(params); i++ {
			if strings.Contains("-d", params[i]) {
				params[i] = "-d=false"
				break
			}
		}
		cmd := exec.Command(os.Args[0], params...)
		cmd.Start()
		log.Println("Started ", os.Args[0], "daemon with PID =", cmd.Process.Pid)
		os.Exit(0)
	} else {
		srv, e := socks5.New(socks5conf)
		if e != nil {
			log.Fatal(e)
		}

		log.Printf("Proxy starting on %s\n", conf.Address+":"+conf.Port)
		if e := srv.ListenAndServe("tcp", conf.Address+":"+conf.Port); e != nil {
			log.Fatal(e)
		}
	}
}
