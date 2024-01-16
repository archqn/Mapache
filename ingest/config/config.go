package config

import "os"

var Version = "1.0.0"
var Env = os.Getenv("ENV")
var Port = os.Getenv("PORT")

var AuthSigningKey = os.Getenv("AUTH_SIGNING_KEY")

var DatabaseHost = os.Getenv("DB_HOST")
var DatabasePort = os.Getenv("DB_PORT")
var DatabaseName = os.Getenv("DB_NAME")
var DatabaseUser = os.Getenv("DB_USER")
var DatabasePassword = os.Getenv("DB_PASSWORD")

var Banner = `
███╗   ███╗ █████╗ ██████╗  █████╗  ██████╗██╗  ██╗███████╗    ██╗███╗   ██╗ ██████╗ ███████╗███████╗████████╗
████╗ ████║██╔══██╗██╔══██╗██╔══██╗██╔════╝██║  ██║██╔════╝    ██║████╗  ██║██╔════╝ ██╔════╝██╔════╝╚══██╔══╝
██╔████╔██║███████║██████╔╝███████║██║     ███████║█████╗      ██║██╔██╗ ██║██║  ███╗█████╗  ███████╗   ██║   
██║╚██╔╝██║██╔══██║██╔═══╝ ██╔══██║██║     ██╔══██║██╔══╝      ██║██║╚██╗██║██║   ██║██╔══╝  ╚════██║   ██║   
██║ ╚═╝ ██║██║  ██║██║     ██║  ██║╚██████╗██║  ██║███████╗    ██║██║ ╚████║╚██████╔╝███████╗███████║   ██║   
╚═╝     ╚═╝╚═╝  ╚═╝╚═╝     ╚═╝  ╚═╝ ╚═════╝╚═╝  ╚═╝╚══════╝    ╚═╝╚═╝  ╚═══╝ ╚═════╝ ╚══════╝╚══════╝   ╚═╝
`
