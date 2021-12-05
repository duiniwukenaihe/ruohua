package main

import (
	"github.com/chenhqchn/ruohua/server/router"
	selfCasbin "github.com/chenhqchn/ruohua/server/utils/casbin"
	"github.com/chenhqchn/ruohua/server/utils/config"
	"log"
	"os"
)

func main() {
	// load the config file
	err :=  config.LoadConfig()
	if err != nil {
		log.Fatalf("Config file failed to load：\n%s", err.Error())
	}

	// log module initialization
	config.InitLogger()
	config.L().Debug("Log module initialization is complete")

	//get mysql conn
	config.L().Debug("Start to connecting to DB")
	err = config.InitDB()
	if err != nil {
		config.L().Errorf("Connected to the DB failed：\n%s", err.Error())
		os.Exit(-1)
	}
	config.L().Debug("DB is ready")

	err = selfCasbin.InitEnforcer()
	if err != nil {
		config.L().Errorf("Failed to initialize casbin adapter or enforcer: \n%s", err)
		os.Exit(-1)
	}

	config.L().Debug("Web server nitializing")
	r := router.InitRouter()
	err = r.Run(":8080")
	if err != nil {
		log.Fatalf("starting failed, %s", err.Error())
	}
}
