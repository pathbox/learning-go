package ninetail

import (
	"io"
	"os"
	"sync"

	"github.com/mattn/go-colorable"
)

type NineTail struct {
	output  io.Writer
	tailers []*Tailer
}

type Config struct {
	Colorize bool
}

func Runner(filenames []string, config Config) (*NineTail, error) {
	var output io.Writer
	if config.Colorize {
		output = colorable.NewColorableStdout() // set colorable stdout
	} else {
		output = colorable.NewNonColorable(os.Stdout)
	}

	tailers, err := NewTailers(filenames)
	if err != nil {
		return nil, err
	}

	return &NineTail{
		output:  output,
		tailers: tailers,
	}, nil
}

func (n *NineTail) Run() {
	var wg sync.WaitGroup

	for _, t := range n.tailers { // support multi files
		wg.Add()
		go func(t *Tailer) {
			t.Do(n.output)
			wg.Done()
		}(t)
	}

	wg.Wait()
}
