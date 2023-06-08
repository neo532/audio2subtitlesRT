package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	nls "github.com/aliyun/alibabacloud-nls-go-sdk"
)

const (
	AKID  = "LTAI5tKW2uTBohLaGv3f9QmK"
	AKKEY = "0Qm8OnbM5kaRV7pYvwsEy4MVcOs6Lt"
	//online key
	APPKEY = "BIUkvGjFD43DmIv1"
	TOKEN  = "cdfff4d6b6084fcc88c4afacad108381"
)

func onTaskFailed(text string, param interface{}) {
	logger, ok := param.(*nls.NlsLogger)
	if !ok {
		log.Default().Fatal("invalid logger")
		return
	}

	logger.Debugln("TaskFailed:", text)
}

func onStarted(text string, param interface{}) {
	logger, ok := param.(*nls.NlsLogger)
	if !ok {
		log.Default().Fatal("invalid logger")
		return
	}

	logger.Debugln("onStarted:", text)
}

func onSentenceBegin(text string, param interface{}) {
	logger, ok := param.(*nls.NlsLogger)
	if !ok {
		log.Default().Fatal("invalid logger")
		return
	}

	logger.Debugln("onSentenceBegin:", text)
}

func onSentenceEnd(text string, param interface{}) {
	logger, ok := param.(*nls.NlsLogger)
	if !ok {
		log.Default().Fatal("invalid logger")
		return
	}

	logger.Debugln("onSentenceEnd:", text)
}

func onResultChanged(text string, param interface{}) {
	logger, ok := param.(*nls.NlsLogger)
	if !ok {
		log.Default().Fatal("invalid logger")
		return
	}

	var result nls.CommonRequest
	json.Unmarshal([]byte(text), &result)
	fmt.Println(fmt.Sprintf("%+v", result.Payload["result"]))

	logger.Debugln("onResultChanged:", text)
}

func onCompleted(text string, param interface{}) {
	logger, ok := param.(*nls.NlsLogger)
	if !ok {
		log.Default().Fatal("invalid logger")
		return
	}

	logger.Debugln("onCompleted:", text)
}

func onClose(param interface{}) {
	logger, ok := param.(*nls.NlsLogger)
	if !ok {
		log.Default().Fatal("invalid logger")
		return
	}

	logger.Debugln("onClosed:")
}

func waitReady(ch chan bool, logger *nls.NlsLogger) error {
	select {
	case done := <-ch:
		{
			if !done {
				logger.Debugln("Wait failed")
				return errors.New("wait failed")
			}
			logger.Debugln("Wait done")
		}
	case <-time.After(20 * time.Second):
		{
			logger.Debugln("Wait timeout")
			return errors.New("wait timeout")
		}
	}
	return nil
}

func main() {

	param := nls.DefaultSpeechTranscriptionParam()
	param.EnableWords = true
	config, _ := nls.NewConnectionConfigWithAKInfoDefault(nls.DEFAULT_URL, APPKEY, AKID, AKKEY)

	logger := nls.NewNlsLogger(os.Stderr, "test", log.LstdFlags|log.Lmicroseconds)
	logger.SetLogSil(false)
	logger.SetDebug(false)
	st, err := nls.NewSpeechTranscription(config, logger,
		onTaskFailed, onStarted,
		onSentenceBegin, onSentenceEnd, onResultChanged,
		onCompleted, onClose, logger)
	if err != nil {
		logger.Fatalln(err)
		return
	}

	logger.Debugln("ST start")
	ready, err := st.Start(param, nil)
	if err != nil {
		st.Shutdown()
		return
	}

	err = waitReady(ready, logger)
	if err != nil {
		st.Shutdown()
		return
	}

	reader := bufio.NewReader(os.Stdin)
	buf := make([]byte, 4*1024)
	for {
		n, err := reader.Read(buf)
		if err != nil {
			break
		}
		st.SendAudioData(buf[:n])
	}

	ready, err = st.Stop()
	if err != nil {
		st.Shutdown()
		return
	}

	err = waitReady(ready, logger)
	if err != nil {
		st.Shutdown()
		return
	}

	logger.Debugln("Sr done")
	st.Shutdown()
	return
}
