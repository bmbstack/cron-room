package helper

import (
	_ "github.com/joho/godotenv/autoload"
	"github.com/joho/godotenv"
	"path/filepath"
	"os"
	"strconv"
	"fmt"
	"github.com/labstack/gommon/color"
)

type Env struct {
	DebugOn             bool
	DBDialect         string // db_dialect=mysql
	DBHost            string // db_host=127.0.0.1
	DBPort            int64  // db_port=3306
	DBName            string // db_name=forum
	DBUser            string // db_user=wangming
	DBPassword        string // db_password=bmbstack@123
	LogFilePathPrefix string // log_file_path_prefix=/var/log/cron-room-
}

var env *Env

func InitEnv() {
	line := "==============================="
	// parse .env or .env.example
	workingDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	env = &Env{}
	err = godotenv.Load(filepath.Join(workingDir, ".env"))
	if err != nil {
		err = godotenv.Load(filepath.Join(workingDir, ".env.example"))
		if err != nil {
			fmt.Println(fmt.Sprintf(color.Red(".env or .env.example is not exist!")))
			panic(err)
		} else {
			fmt.Println(fmt.Sprintf("%s%s%s",
				color.White(line),
				color.Bold(color.Green(".env.example information")),
				color.White(line)))
		}
	} else {
		fmt.Println(fmt.Sprintf("%s%s%s",
			color.White(line),
			color.Bold(color.Green(".env information")),
			color.White(line)))
	}

	env.DebugOn, err = strconv.ParseBool(os.Getenv("debug_on"))
	if err != nil {
		fmt.Println(fmt.Sprintf(color.Red("env debug is parse error, it must be bool!")))
		panic(err)
	}

	// db
	env.DBDialect = os.Getenv("db_dialect")
	env.DBHost = os.Getenv("db_host")
	env.DBPort, err = strconv.ParseInt(os.Getenv("db_port"), 10, 64)
	if err != nil {
		fmt.Println(fmt.Sprintf(color.Red("env db_port is parse error, it must be int64!")))
		panic(err)
	}
	env.DBName = os.Getenv("db_name")
	env.DBUser = os.Getenv("db_user")
	env.DBPassword = os.Getenv("db_password")

	env.LogFilePathPrefix = os.Getenv("log_file_path_prefix")

	fmt.Println(fmt.Sprintf("env.debug_on=%s", color.Green(env.DebugOn)))
	fmt.Println(fmt.Sprintf("env.db_dialect=%s", color.Green(env.DBDialect)))
	fmt.Println(fmt.Sprintf("env.db_host=%s", color.Green(env.DBHost)))
	fmt.Println(fmt.Sprintf("env.db_port=%s", color.Green(env.DBPort)))
	fmt.Println(fmt.Sprintf("env.db_name=%s", color.Green(env.DBName)))
	fmt.Println(fmt.Sprintf("env.db_user=%s", color.Green(env.DBUser)))
	fmt.Println(fmt.Sprintf("env.db_password=%s", color.Green(env.DBPassword)))
	fmt.Println(fmt.Sprintf("env.log_file_path_prefix=%s", color.Green(env.LogFilePathPrefix)))
}

func GetEnv() *Env {
	return env
}
