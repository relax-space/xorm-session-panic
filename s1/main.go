package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo"

	"github.com/chneau/limiter"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"github.com/pangpanglabs/goutils/echomiddleware"
)

// panic when get db from context
func main() {
	db, err := xorm.NewEngine("mysql", "root:1234@tcp(127.0.0.1:3306)/fruit")
	if err != nil {
		fmt.Println(err)
		return
	}
	db.Sync(new(Fruit))
	e := echo.New()
	e.Use(echomiddleware.ContextDB("xorm-session-panic", db, echomiddleware.KafkaConfig{
		Brokers: []string{
			"127.0.0.1:9092",
		},
		Topic: "behaviorlog",
	}))

	e.GET("/ping", func(c echo.Context) error {
		limit := limiter.New(10)
		for index := 0; index < 1000; index++ {
			limit.Execute(func() {
				if _, err := query(c.Request().Context()); err != nil {
					fmt.Println(err)
				}
			})
		}
		limit.Wait()
		fmt.Println("done")
		return c.String(http.StatusOK, "pong")
	})
	e.Start(":8080")

}

func query(ctx context.Context) (fruits []*Fruit, err error) {
	err = DB(ctx).Find(&fruits)
	return
}

func DB(ctx context.Context) *xorm.Session {
	v := ctx.Value(echomiddleware.ContextDBName)
	if v == nil {
		panic("ctx is not exist")
	}
	if db, ok := v.(*xorm.Session); ok {
		return db
	}
	if db, ok := v.(*xorm.Engine); ok {
		return db.NewSession()
	}
	panic("DB is not exist")
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
