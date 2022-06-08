package middleware

import (
	"context"
	"encoding/base64"
	"math/rand"

	"k8s.io/klog/v2"

	"github.com/authzed/spicedb-operator/pkg/libctrl"
	"github.com/authzed/spicedb-operator/pkg/libctrl/handler"
)

var CtxSyncID = libctrl.NewContextKey[string]()

func NewSyncID(size uint8) string {
	buf := make([]byte, size)
	rand.Read(buf) // nolint
	str := base64.StdEncoding.EncodeToString(buf)
	return str[:size]
}

func SyncIDMiddleware(in handler.Handler) handler.Handler {
	return handler.NewHandlerFromFunc(func(ctx context.Context) {
		_, ok := CtxSyncID.Value(ctx)
		if !ok {
			ctx = CtxSyncID.WithValue(ctx, NewSyncID(5))
		}
		in.Handle(ctx)
	}, in.ID())
}

func KlogMiddleware(ref klog.ObjectRef) libctrl.HandlerMiddleware {
	return func(in handler.Handler) handler.Handler {
		return handler.NewHandlerFromFunc(func(ctx context.Context) {
			klog.V(4).InfoS("entering handler", "syncID", CtxSyncID.MustValue(ctx), "object", ref, "handler", in.ID())
			in.Handle(ctx)
		}, in.ID())
	}
}