package contract

//ServerSettings - common server settings
type ServerSettings struct {
	Port      string
	LogConfig *LogConfig
}

//LogConfig - settings for logger
type LogConfig struct {
	FileName string
	SizeMb   int
}
