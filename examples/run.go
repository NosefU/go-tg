package examples

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/nosefu/go-tg"
	"github.com/nosefu/go-tg/tgb"
)

// Run runs bot with given router.
// Exit on error.
func Run(handler tgb.Handler, opts ...tg.ClientOption) {
	ctx := context.Background()

	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt, os.Kill, syscall.SIGTERM)
	defer cancel()

	if err := run(ctx, handler, nil, opts...); err != nil {
		log.Printf("error: %v", err)
		defer os.Exit(1)
	}
}

func run(
	ctx context.Context,
	handler tgb.Handler,
	do func(ctx context.Context, client *tg.Client) error,
	opts ...tg.ClientOption,
) error {
	// define flags
	var (
		flagToken         string
		flagServer        string
		flagTestEnv       bool
		flagWebhookURL    string
		flagWebhookListen string
	)

	// parse flags
	flag.StringVar(&flagToken, "token", "", "Telegram Bot API token")
	flag.StringVar(&flagServer, "server", "https://api.telegram.org", "Telegram Bot API server")
	flag.BoolVar(&flagTestEnv, "test-env", false, "switch bot to test environment")
	flag.StringVar(&flagWebhookURL, "webhook-url", "", "Telegram Bot API webhook URL, if not provdide run in longpoll mode")
	flag.StringVar(&flagWebhookListen, "webhook-listen", ":8000", "Telegram Bot API webhook URL")
	flag.Parse()

	if flagToken == "" {
		return fmt.Errorf("token is required, provide it with -token flag")
	}

	opts = append(opts,
		tg.WithClientServerURL(flagServer),
	)

	if flagTestEnv {
		opts = append(opts, tg.WithClientTestEnv())
	}

	client := tg.New(flagToken, opts...)

	me, err := client.Me(ctx)
	if err != nil {
		return fmt.Errorf("get me: %w", err)
	}

	log.Printf("authorized as %s", me.Username.Link())

	if do != nil {
		return do(ctx, client)
	} else if flagWebhookURL != "" {
		err = tgb.NewWebhook(
			handler,
			client,
			flagWebhookURL,
			tgb.WithDropPendingUpdates(true),
			tgb.WithWebhookLogger(log.Default()),
		).Run(
			ctx,
			flagWebhookListen,
		)
	} else {
		err = tgb.NewPoller(
			handler,
			client,
			tgb.WithPollerLogger(log.Default()),
			tgb.WithPollerAllowedUpdates(
				tg.UpdateTypeMessage,
				tg.UpdateTypeEditedMessage,
				tg.UpdateTypeChannelPost,
				tg.UpdateTypeEditedChannelPost,
				tg.UpdateTypeMessageReaction,
				tg.UpdateTypeMessageReactionCount,
				tg.UpdateTypeInlineQuery,
				tg.UpdateTypeChosenInlineResult,
				tg.UpdateTypeCallbackQuery,
				tg.UpdateTypeShippingQuery,
				tg.UpdateTypePreCheckoutQuery,
				tg.UpdateTypePoll,
				tg.UpdateTypePollAnswer,
				tg.UpdateTypeMyChatMember,
				tg.UpdateTypeChatMember,
				tg.UpdateTypeChatJoinRequest,
				tg.UpdateTypeChatBoost,
				tg.UpdateTypeRemovedChatBoost,
			),
		).Run(ctx)
	}

	return err
}
