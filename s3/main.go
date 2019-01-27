package main

import (
	"fmt"
	"time"

	"github.com/chneau/limiter"
	_ "github.com/go-sql-driver/mysql"

	"github.com/go-xorm/xorm"
)

//don't panic when get db from session
func main() {
	db, err := xorm.NewEngine("mysql", "root:1234@tcp(127.0.0.1:3306)/fruit")
	if err != nil {
		fmt.Println(err)
		return
	}
	db.Sync(new(Fruit))
	session := db.NewSession()

	limit := limiter.New(10)
	for index := 0; index < 10000; index++ {
		session := *session
		limit.Execute(func() {
			if _, err := query(session); err != nil {
				fmt.Println(err)
			}
		})
	}
	limit.Wait()
	fmt.Println("done")
	time.Sleep(1 * time.Hour)

}

func query(session xorm.Session) (fruits []*Fruit, err error) {
	err = session.Find(&fruits)
	return
}

type Fruit struct {
	Id        int64      `json:"id" xorm:"int64 notnull autoincr pk 'id'"`
	Code      string     `json:"code"`
	Name      string     `json:"name"`
	Color     string     `json:"color"`
	Price     int64      `json:"price"`
	StoreCode string     `json:"storeCode"`
	CreatedAt *time.Time `json:"createdAt" xorm:"created"`
	UpdatedAt *time.Time `json:"updatedAt" xorm:"updated"`
	DeletedAt *time.Time `json:"deletedAt" xorm:"deleted"`
	UniqueId  uint64     `json:"uniqueId"`
}
