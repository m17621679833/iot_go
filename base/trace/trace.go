package base

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	logFacede "iot_go/base/log"
	"math/rand"
	"net"
	"os"
	"strings"
	"time"
)

var LocalIP = net.ParseIP("127.0.0.1")

const (
	Undefined    = "_undef"
	MysqlFailed  = "_com_mysql_failure"
	MysqlSuccess = "_com_mysql_success"
	RedisFailed  = "_com_redis_failure"
	RedisSuccess = "_com_redis_success"
	HttpFailed   = "_com_http_failure"
	HttpSuccess  = "_com_http_success"
	TcpSuccess   = "_com_tcp_success"
	TcpFailed    = "_com_tcp_fail"
	RequestIn    = "_com_request_in"
	RequestOut   = "_com_request_out"
)

const (
	_bizTag       = "bizTag"
	_traceId      = "traceId"
	_spanId       = "spanId"
	_childSpanId  = "childSpanId"
	_tagBizPrefix = "_com_"
	_tagBizUndef  = "_com_undef"
)

type Logger struct {
}

type Trace struct {
	TraceId     string
	SpanId      string
	Caller      string
	SrcMethod   string
	HintCode    int64
	HintContent string
}

type TraceContext struct {
	Trace
	CSpanId string
}

func CheckTag(tag string) string {
	if strings.HasPrefix(tag, _tagBizPrefix) {
		return tag
	}
	return _tagBizUndef
}

func (l *Logger) Info(trace *TraceContext, bizTag string, m map[string]interface{}) {
	m[_bizTag] = CheckTag(bizTag)
	m[_traceId] = trace.TraceId
	m[_childSpanId] = trace.CSpanId
	m[_spanId] = trace.SpanId
	logFacede.Info(parseParams(m))
}

func (l *Logger) Warning(trace *TraceContext, bizTag string, m map[string]interface{}) {
	m[_bizTag] = CheckTag(bizTag)
	m[_traceId] = trace.TraceId
	m[_childSpanId] = trace.CSpanId
	m[_spanId] = trace.SpanId
	logFacede.Warning(parseParams(m))
}

func (l *Logger) Error(trace *TraceContext, bizTag string, m map[string]interface{}) {
	m[_bizTag] = CheckTag(bizTag)
	m[_traceId] = trace.TraceId
	m[_childSpanId] = trace.CSpanId
	m[_spanId] = trace.SpanId
	logFacede.Error(parseParams(m))
}

func (l *Logger) Trace(trace *TraceContext, bizTag string, m map[string]interface{}) {
	m[_bizTag] = CheckTag(bizTag)
	m[_traceId] = trace.TraceId
	m[_childSpanId] = trace.CSpanId
	m[_spanId] = trace.SpanId
	logFacede.Trace(parseParams(m))
}

func (l *Logger) Debug(trace *TraceContext, bizTag string, m map[string]interface{}) {
	m[_bizTag] = CheckTag(bizTag)
	m[_traceId] = trace.TraceId
	m[_childSpanId] = trace.CSpanId
	m[_spanId] = trace.SpanId
	logFacede.Debug(parseParams(m))
}

func (l *Logger) Close() {
	logFacede.Close()
}

func parseParams(m map[string]interface{}) string {
	var bizTag = Undefined
	if bizTag, have := m[_bizTag]; have {
		if val, ok := bizTag.(string); ok {
			bizTag = val
		}
	}
	for _key, _val := range m {
		if _key == _bizTag {
			continue
		}
		bizTag = bizTag + "||" + fmt.Sprintf("%v=%+v", _key, _val)
	}
	bizTag = strings.Trim(fmt.Sprintf("%q", bizTag), "\"")
	return bizTag
}

func NewSpanId() string {
	timestamp := uint32(time.Now().Unix())
	ipToLong := binary.BigEndian.Uint32(LocalIP.To4())
	b := bytes.Buffer{}
	b.WriteString(fmt.Sprintf("%08x", ipToLong^timestamp))
	b.WriteString(fmt.Sprintf("%08x", rand.Int31()))
	return b.String()
}
func calcTraceId(ip string) (traceId string) {
	now := time.Now()
	timestamp := uint32(now.Unix())
	timeNano := now.UnixNano()
	pid := os.Getpid()

	b := bytes.Buffer{}
	netIP := net.ParseIP(ip)
	if netIP == nil {
		b.WriteString("00000000")
	} else {
		b.WriteString(hex.EncodeToString(netIP.To4()))
	}
	b.WriteString(fmt.Sprintf("%08x", timestamp&0xffffffff))
	b.WriteString(fmt.Sprintf("%04x", timeNano&0xffff))
	b.WriteString(fmt.Sprintf("%04x", pid&0xffff))
	b.WriteString(fmt.Sprintf("%06x", rand.Int31n(1<<24)))
	b.WriteString("b0") // 末两位标记来源,b0为go

	return b.String()
}
