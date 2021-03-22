package gateway

import (
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"

	"github.com/smallnest/rpcx/protocol"
)

const (
	XVersion           = "X-RPCX-Version"
	XMessageType       = "X-RPCX-MesssageType"
	XHeartbeat         = "X-RPCX-Heartbeat"
	XOneway            = "X-RPCX-Oneway"
	XMessageStatusType = "X-RPCX-MessageStatusType"
	XSerializeType     = "X-RPCX-SerializeType"
	XMessageID         = "X-RPCX-MessageID"
	XServicePath       = "X-RPCX-ServicePath"
	XServiceMethod     = "X-RPCX-ServiceMethod"
	XMeta              = "X-RPCX-Meta"
	XErrorMessage      = "X-RPCX-ErrorMessage"
)

func HttpRequest2RpcxRequest(r *http.Request) (*protocol.Message, error) {
	req := protocol.NewMessage()
	req.SetMessageType(protocol.Request)

	h := r.Header
	seq := h.Get(XMessageID)
	if seq != "" {
		id, err := strconv.ParseUint(seq, 10, 64)
		if err != nil {
			return nil, err
		}
		req.SetSeq(id)
	}

	heartbeat := h.Get(XHeartbeat)
	if heartbeat != "" {
		req.SetHeartbeat(true)
	}

	oneway := h.Get(XOneway)
	if oneway != "" {
		req.SetOneway(true)
	}

	if h.Get("Content-Encoding") == "gzip" {
		req.SetCompressType(protocol.Gzip)
	}

	st := h.Get(XSerializeType)
	if st != "" {
		rst, err := strconv.Atoi(st)
		if err != nil {
			return nil, err
		}
		req.SetSerializeType(protocol.SerializeType(rst))
	} else {
		return nil, errors.New("empty serialized type")
	}

	meta := h.Get(XMeta)
	if meta != "" {
		metadata, err := url.ParseQuery(meta)
		if err != nil {
			return nil, err
		}
		mm := make(map[string]string)
		for k, v := range metadata {
			if len(v) > 0 {
				mm[k] = v[0]
			}
		}
		req.Metadata = mm
	}

	sp := h.Get(XServicePath)
	if sp != "" {
		req.ServicePath = sp
	} else {
		return nil, errors.New("empty servicepath")
	}

	sm := h.Get(XServiceMethod)
	if sm != "" {
		req.ServiceMethod = sm
	} else {
		return nil, errors.New("empty servicemethod")
	}

	payload, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	req.Payload = payload

	return req, nil
}
