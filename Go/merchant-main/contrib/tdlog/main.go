package tdlog

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/nats-io/nats.go"
)

var (
	timeout time.Duration
	conn    *nats.Conn
)

//字符串特殊字符转译
func addslashes(str string) string {

	tmpRune := []rune{}
	strRune := []rune(str)
	for _, ch := range strRune {
		switch ch {
		case []rune{'\\'}[0], []rune{'"'}[0], []rune{'\''}[0]:
			tmpRune = append(tmpRune, []rune{'\\'}[0])
			tmpRune = append(tmpRune, ch)
		default:
			tmpRune = append(tmpRune, ch)
		}
	}

	return string(tmpRune)
}

func write(fields map[string]string, flags string) error {

	var b strings.Builder

	t := time.Now()

	b.WriteString("INSERT INTO zlog (ts, filename, content, fn, flags, id, project) VALUES(")
	b.WriteByte('"')
	b.WriteString(t.Format("2006-01-02 15:04:05.000"))
	b.WriteByte('"')
	b.WriteByte(',')
	b.WriteByte('"')
	b.WriteString(addslashes(fields["filename"]))
	b.WriteByte('"')
	b.WriteByte(',')
	b.WriteByte('"')
	b.WriteString(addslashes(fields["content"]))
	b.WriteByte('"')
	b.WriteByte(',')
	b.WriteByte('"')
	b.WriteString(fields["fn"])
	b.WriteByte('"')
	b.WriteByte(',')

	b.WriteByte('"')
	b.WriteString(flags)
	b.WriteByte('"')
	b.WriteByte(',')
	b.WriteByte('"')
	b.WriteString(fields["id"])
	b.WriteByte('"')
	b.WriteByte(',')
	b.WriteByte('"')
	b.WriteString(fields["project"])
	b.WriteByte('"')
	b.WriteByte(')')

	err := conn.Publish("zlog", []byte(b.String()))
	return err
}

func Login(fields map[string]string) error {

	var b strings.Builder

	t := time.Now()

	b.WriteString("INSERT INTO login (ts, uid, username, ip) VALUES(")
	b.WriteByte('"')
	b.WriteString(t.Format("2006-01-02 15:04:05.000"))
	b.WriteByte('"')
	b.WriteByte(',')
	b.WriteByte('"')
	b.WriteString(fields["uid"])
	b.WriteByte('"')
	b.WriteByte(',')
	b.WriteByte('"')
	b.WriteString(fields["username"])
	b.WriteByte('"')
	b.WriteByte(',')
	b.WriteByte('"')
	b.WriteString(fields["ip"])
	b.WriteByte('"')
	b.WriteByte(')')

	err := conn.Publish("zlog", []byte(b.String()))
	return err

}

func buildSql(table string, data map[string]string) string {

	ts := time.Now()
	keys := []string{}
	values := []string{}

	for k, v := range data {

		keys = append(keys, k)
		values = append(values, addslashes(v))
	}

	query := fmt.Sprintf("INSERT INTO %s (ts,%s) VALUES(\"%s\",\"%s\")", table, strings.Join(keys, ","), ts.Format("2006-01-02 15:04:05.000"), strings.Join(values, "\",\""))

	return query
}

func WriteLog(table string, data map[string]string) error {

	query := buildSql(table, data)

	err := conn.Publish("zlog", []byte(query))
	return err
}

func InitNatsIO(urls []string, name, password string) *nats.Conn {

	opts := nats.Options{
		Servers:        urls,
		User:           name,
		Password:       password,
		AllowReconnect: true,
		MaxReconnect:   10,
		ReconnectWait:  5 * time.Second,
		Timeout:        1 * time.Second,
		FlusherTimeout: 5 * time.Second,
	}

	nc, err := opts.Connect()
	if err != nil {
		log.Fatalln(err)
	}

	return nc
}

func New(urls []string, name, password string) {

	conn = InitNatsIO(urls, name, password)
}

func Info(fields map[string]string) error {
	return write(fields, "info")
}

func Warn(fields map[string]string) error {
	return write(fields, "warn")
}

func Error(fields map[string]string) error {
	return write(fields, "error")
}
