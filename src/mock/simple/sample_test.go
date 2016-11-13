package sample

import (
	"testing"

	"github.com/golang/mock/gomock"

	mock "./mock_sample"
)

func TestSample1(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSample := mock.NewMockSample(ctrl)
	mockSample.EXPECT().Method("hoge").Return(1)

	t.Log("result", mockSample.Method("hoge"))
}

func TestSample2(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSample := mock.NewMockSample(ctrl)

	// 呼び出される順番を指定
	gomock.InOrder(
		mockSample.EXPECT().Method("hoge").Return(1),
		mockSample.EXPECT().Method("fuga").Return(2),
	)
	/* // 上記と同じ
	mockSample.EXPECT().Method("hoge").Return(1).After(
		mockSample.EXPECT().Method("fuga").Return(2),
	)
	*/

	t.Log("result", mockSample.Method("hoge"))
	t.Log("result", mockSample.Method("fuga"))
}

func TestSample3(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSample := mock.NewMockSample(ctrl)

	// 呼び出される回数を指定
	mockSample.EXPECT().Method("hoge").Return(1).Times(2)

	// 何回でも良い場合は，AnyTimes
	// mockSample.EXPECT().Method("hoge").Return(1).AnyTimes()

	t.Log("result", mockSample.Method("hoge"))
	t.Log("result", mockSample.Method("hoge"))
}
