// Package contains example of using tgb.ChatType filter.
package main

import (
	"context"

	"github.com/nosefu/go-tg"
	"github.com/nosefu/go-tg/examples"
	"github.com/nosefu/go-tg/tgb"
)

func main() {
	examples.Run(tgb.NewRouter().
		Message(func(ctx context.Context, msg *tgb.MessageUpdate) error {
			return msg.Answer("this is private chat response").DoVoid(ctx)
		}, tgb.ChatType(tg.ChatTypePrivate)).
		Message(func(ctx context.Context, msg *tgb.MessageUpdate) error {
			return msg.Answer("this is group chat response").DoVoid(ctx)
		}, tgb.ChatType(tg.ChatTypeGroup, tg.ChatTypeSupergroup)),
	)
}
