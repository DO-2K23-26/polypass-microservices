package sql

type Config struct {
	//path to migrations
	Migrations string `mapstructure:"migrations"`
	//username for database
	Username string `mapstructure:"username"`
	//password for database
	Password string `mapstructure:"password"`
	//host for database
	Host string `mapstructure:"host"`
	//port for database
	Port int `mapstructure:"port"`
	//database name
	Dbname string `mapstructure:"dbname"`
}
