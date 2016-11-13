package sample

import (
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"

	mock "./mock_sample"
)

func TestSample(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// sample.writerとio.Writerを共に実装
	w := mock.NewMockwriter(ctrl)

	gomock.InOrder(
		w.EXPECT().Write([]byte("hoge")).Return(4, nil),
		w.EXPECT().Write([]byte("fuga")).Return(4, nil),
	)

	// io.Writerとして渡す
	fmt.Fprintf(w, "hoge")
	fmt.Fprintf(w, "fuga")
}
