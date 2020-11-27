package example

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
)

func TestGetFromDB(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish() // 断言 DB.Get() 方法是否被调用

	m := NewMockDB(ctrl)
	m.EXPECT().Get(gomock.Eq("Tom")).Return(100, errors.New("not exist"))

	v := GetFromDB(m, "Tom")
	if v != -1 {
		t.Fatal("expected -1, but got", v)
	}

	// m.EXPECT().Get(gomock.Any()).Return(630, nil)
	// m.EXPECT().Get(gomock.Any()).DoAndReturn(func(key string) (int, error) {
	// 	if key == "Sam" {
	// 		return 630, nil
	// 	}
	// 	return 0, errors.New("not exist")
	// })
}

// mockgen -source db.go -destination=db_mock.go -package=example
// go test . -cover -v
