package main

import (
	"fmt"
  canal "github.com/siddontang/go-mysql/canal"
	"runtime/debug"
)

type binlogHandler struct {
	canal.DummyEventHandler
	BinlogParser
}

func (h *binlogHandler) OnRow(e *canal.RowsEvent) error {
	defer func() {
		if r := recover(); r != nil {
			fmt.Print(r, " ", string(debug.Stack()))
		}
	}()

	var n = 0
	var k = 1

	if e.Action == canal.UpdateAction {
		n = 1
		k = 2
	}

	switch e.Action {
			case canal.UpdateAction:
				oldUser := User{}
				h.GetBinLogData(&oldUser, e, i-1)
				fmt.Printf("User %d name changed from %s to %s\n", user.Id, oldUser.Name, user.Name)
			case canal.InsertAction:
				fmt.Printf("User %d is created with name %s\n", user.Id, user.Name)
			case canal.DeleteAction:
				fmt.Printf("User %d is deleted with name %s\n", user.Id, user.Name)
			default:
				fmt.Printf("Unknown action")
			}
		}

	}
	return nil
}

func (h *binlogHandler) String() string {
	return "binlogHandler"
}

func binlogListener() {
	c, err := getDefaultCanal()
	if err == nil {
		coords, err := c.GetMasterPos()
		if err == nil {
			c.SetEventHandler(&binlogHandler{})
			c.RunFrom(coords)
		}
	}
}

func getDefaultCanal() (*canal.Canal, error) {
	cfg := canal.NewDefaultConfig()
	cfg.Addr = fmt.Sprintf("%s:%d", "mariadb", 3307)
	cfg.User = "root"
	cfg.Password = "root"
	cfg.Flavor = "mysql"

	cfg.Dump.ExecutionPath = ""

	return canal.NewCanal(cfg)
}