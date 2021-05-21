package mysql

import "github.com/spf13/pflag"

type DBOptions struct {
	DSN string
}

func NewDBOptions() *DBOptions {
	return &DBOptions{
		DSN: "root:root@tcp(127.0.0.1:3306)/crond?charset=utf8&parseTime=True&loc=Local",
	}
}

func (o *DBOptions) BindFlags(fs *pflag.FlagSet) {
	if o == nil {
		return
	}

	fs.StringVar(&o.DSN, "mysql-dsn", o.DSN, "MySQL DSN")
}
