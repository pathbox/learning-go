package holmes

import (
	"sync"
	"testing"
)

// "time"

func TestStdErrLogger(t *testing.T) {
	defer Start().Stop()
	Infoln("Hello, Mike")
	Warnln("This might be painful but...")
	Errorln("You have to go through it until sunshine comes out")
	Infoln("Those were the days hard work forever pays")
}

func TestFileLoggerEveryMinute(t *testing.T) {
	defer Start(LogFilePath("./log"), EveryMinute).Stop()
	for i := 0; i < 100; i++ {
		// after one minute a new log created, uncomment it to have a try!
		// time.Sleep(time.Second)
		Infof("%s", "Jingle bells, jingle bells,")
		Warnf("%s", "Jingle all the way.")
		Errorf("%s", "Oh! what fun it is to ride")
		Infof("%s", "In a one-horse open sleigh.")
	}
}

func TestFileLoggerMultipleGoroutine(t *testing.T) {
	defer Start(LogFilePath("./log"), EveryHour).Stop()
	wg := sync.WaitGroup{}
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			Infof("%s", "Wake up, Neo")
			Warnf("%s", "The Matrix has you...")
			Errorf("%s", "Follow the white rabbit")
			Infof("%s", "Knock knock!")
			wg.Done()
		}()
	}
	wg.Wait()
}

func TestAlsoStdout(t *testing.T) {
	defer Start(EveryHour, AlsoStdout).Stop()
	Infof("%s\n", "If by life you were deceived,")
	Warnf("%s\n", "Don't be dismal, don't be wild!")
	Errorf("%s", "In the day of grief, be mild.")
	Infof("%s", "Merry days will come, believe.")
}

func TestPrintStack(t *testing.T) {
	defer Start(LogFilePath("./log"), EveryHour, PrintStack, AlsoStdout).Stop()
	for i := 0; i < 100; i++ {
		Infof("%s", "If by life you were deceived,")
		Warnf("%s", "Don't be dismal, don't be wild!")
		Errorf("%s", "In the day of grief, be mild.")
		Infof("%s", "Merry days will come, believe.")
	}
}

func BenchmarkFileLoggerSingleGoroutine(b *testing.B) {
	defer Start(LogFilePath("./log"), EveryHour).Stop()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Infof("%s", "Wake up, Neo")
		Warnf("%s", "The Matrix has you...")
		Errorf("%s", "Follow the white rabbit")
		Infof("%s", "Knock knock!")
	}
}

func BenchmarkFileLoggerMultipleGoroutine(b *testing.B) {
	defer Start(LogFilePath("./log"), EveryHour).Stop()
	wg := sync.WaitGroup{}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		wg.Add(1)
		go func() {
			Infof("%s", "Wake up, Neo")
			Warnf("%s", "The Matrix has you...")
			Errorf("%s", "Follow the white rabbit")
			Infof("%s", "Knock knock!")
			wg.Done()
		}()
	}
	wg.Wait()
}
