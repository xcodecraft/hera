package qhydra

import (
	"encoding/json"
	"fmt"
	"log/syslog"
	"os"
	"time"
)

const EVENT_PREFIX string = "pcgameq_"

var qhydras map[string]*Qhydra

func init() {
	qhydras = make(map[string]*Qhydra)
}

type Qhydra struct {
	event    string
	tag      string
	hostname string
	priority syslog.Priority
	writer   *syslog.Writer
}

func GetQhydra(topic string) *Qhydra {
	var qhydra_ *Qhydra
	var exist bool

	if qhydra_, exist = qhydras[topic]; !exist {
		qhydra_ = newQhydra(topic)
		qhydras[topic] = qhydra_
	}

	return qhydra_
}

func newQhydra(event string) *Qhydra {

	qhydra_ := &Qhydra{
		event:    event,
		tag:      EVENT_PREFIX + event,
		priority: syslog.LOG_LOCAL4 | syslog.LOG_INFO,
	}

	if err := qhydra_.init(); err != nil {
		//TODO Log err
		fmt.Println("qhydra_ init errmsg:" + err.Error())
		return nil
	}

	return qhydra_
}

func (this *Qhydra) init() error {
	if writer, err := syslog.New(this.priority, this.tag); err != nil {
		return err
	} else {
		this.writer = writer
	}

	if hostname, err := os.Hostname(); err != nil {
		return err
	} else {
		this.hostname = hostname
	}

	return nil
}

func (this *Qhydra) Trigger(data interface{}, key string) ([]byte, int, error) {

	msg := map[string]interface{}{
		"name": this.event,
		"data": data,
		"host": this.hostname,
		"key":  key,
		"time": time.Now().Format("2006-01-02 15:04:05"),
	}

	msgJson, _ := json.Marshal(msg)
	num, err := this.writer.Write(msgJson)
	return msgJson, num, err
}

func (this *Qhydra) Close() error {
	return this.writer.Close()
}
