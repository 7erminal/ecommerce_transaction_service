package main

import(
	"os"

	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/client/orm/migration"

	_ "github.com/go-sql-driver/mysql"
)

func init(){
	orm.RegisterDataBase("default", "mysql","root:password@tcp(localhost:3306)/AMC")
}

func main(){
	task := "upgrade"
	switch task {
	case "upgrade":
		if err := migration.Upgrade(1739839495); err != nil {
			os.Exit(2)
		}
	case "rollback":
		if err := migration.Rollback("AddColumnToNotificationMessagesTable_20250218_004455"); err != nil {
			os.Exit(2)
		}
	case "reset":
		if err := migration.Reset(); err != nil {
			os.Exit(2)
		}
	case "refresh":
		if err := migration.Refresh(); err != nil {
			os.Exit(2)
		}
	}
}

