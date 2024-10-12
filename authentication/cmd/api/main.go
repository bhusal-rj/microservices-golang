package main

import "database/sql"

const webPort = "3003"

type Config struct {
	DB     *sql.DB
	Models data.Models
}

func main() {

}
