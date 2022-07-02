package tg

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPeerIDImpl(t *testing.T) {
	for _, test := range []struct {
		PeerID PeerID
		Want   string
	}{
		{UserID(1), "1"},
		{ChatID(1), "1"},
		{&Chat{ID: ChatID(1)}, "1"},
	} {
		assert.Equal(t, test.Want, test.PeerID.PeerID())
	}
}

func TestUsername_PeerID(t *testing.T) {
	assert.Equal(t, "@username", Username("username").PeerID())
}

func TestChatType_String(t *testing.T) {
	for _, test := range []struct {
		ChatType ChatType
		Want     string
	}{
		{ChatTypePrivate, "private"},
		{ChatTypeGroup, "group"},
		{ChatTypeSupergroup, "supergroup"},
		{ChatTypeChannel, "channel"},
		{ChatTypeSender, "sender"},
		{ChatType(-1), "unknown"},
	} {
		assert.Equal(t, test.Want, test.ChatType.String())
	}
}

func TestChatAction_String(t *testing.T) {
	for _, test := range []struct {
		ChatAction ChatAction
		Want       string
	}{
		{ChatActionTyping, "typing"},
		{ChatActionUploadPhoto, "upload_photo"},
		{ChatActionUploadVideo, "upload_video"},
		{ChatActionRecordVideo, "record_video"},
		{ChatActionRecordVoice, "record_voice"},
		{ChatActionUploadVoice, "upload_voice"},
		{ChatActionUploadDocument, "upload_document"},
		{ChatActionChooseSticker, "choose_sticker"},
		{ChatActionFindLocation, "find_location"},
		{ChatActionRecordVideoNote, "record_video_note"},
		{ChatActionUploadVideoNote, "upload_video_note"},
		{ChatAction(-1), "unknown"},
	} {
		assert.Equal(t, test.Want, test.ChatAction.String())
	}
}

func TestChatType_MarshalJSON(t *testing.T) {
	type sample struct {
		Type ChatType `json:"type"`
	}
	for _, test := range []struct {
		Sample sample
		Want   string
	}{
		{sample{ChatTypePrivate}, `{"type":"private"}`},
		{sample{ChatTypeGroup}, `{"type":"group"}`},
		{sample{ChatTypeSupergroup}, `{"type":"supergroup"}`},
		{sample{ChatTypeChannel}, `{"type":"channel"}`},
		{sample{ChatTypeSender}, `{"type":"sender"}`},
		{sample{ChatType(-1)}, `{"type":"unknown"}`},
	} {
		actual, err := json.Marshal(test.Sample)
		assert.NoError(t, err)

		assert.Equal(t, test.Want, string(actual))
	}
}

func TestChatType_UnmarshalJSON(t *testing.T) {
	type sample struct {
		Type ChatType `json:"type"`
	}
	for _, test := range []struct {
		Input  string
		Sample sample
		Want   ChatType
		Err    bool
	}{
		{`{"type": "private"}`, sample{}, ChatTypePrivate, false},
		{`{"type": "group"}`, sample{}, ChatTypeGroup, false},
		{`{"type": "supergroup"}`, sample{}, ChatTypeSupergroup, false},
		{`{"type": "channel"}`, sample{}, ChatTypeChannel, false},
		{`{"type": "test"}`, sample{}, ChatType(-1), true},
		{`{"type": "sender"}`, sample{}, ChatTypeSender, false},
		{`{"type": {}}`, sample{}, ChatType(-1), true},
	} {
		if test.Err {
			assert.Error(t, json.Unmarshal([]byte(test.Input), &test.Sample))
		} else {
			assert.NoError(t, json.Unmarshal([]byte(test.Input), &test.Sample))
			assert.Equal(t, test.Want, test.Sample.Type)
		}
	}
}

func TestInlineReplyMarkup(t *testing.T) {
	actual := NewInlineKeyboardMarkup(
		NewButtonRow(
			NewInlineKeyboardButtonURL("text", "https://google.com"),
			NewInlineKeyboardButtonCallback("text", "data"),
			NewInlineKeyboardButtonWebApp("text", WebAppInfo{}),
			NewInlineKeyboardButtonLoginURL("text", LoginURL{
				URL: "https://google.com",
			}),
			NewInlineKeyboardButtonSwitchInlineQuery("text", "query"),
			NewInlineKeyboardButtonSwitchInlineQueryCurrentChat("text", "query"),
			NewInlineKeyboardButtonCallbackGame("text"),
			NewInlineKeyboardButtonPay("text"),
		),
	)

	actual.isReplyMarkup()

	assert.EqualValues(t, InlineKeyboardMarkup{
		InlineKeyboard: [][]InlineKeyboardButton{
			{
				{Text: "text", URL: "https://google.com"},
				{Text: "text", CallbackData: "data"},
				{Text: "text", WebApp: &WebAppInfo{}},
				{Text: "text", LoginURL: &LoginURL{URL: "https://google.com"}},
				{Text: "text", SwitchInlineQuery: "query"},
				{Text: "text", SwitchInlineQueryCurrentChat: "query"},
				{Text: "text", CallbackGame: &CallbackGame{}},
				{Text: "text", Pay: true},
			},
		},
	}, actual)

	assert.Equal(t, actual, *actual.Ptr())
}

func TestReplyKeyboardMarkup(t *testing.T) {
	actual := NewReplyKeyboardMarkup(
		NewButtonRow(
			NewKeyboardButton("text"),
			NewKeyboardButtonRequestContact("text"),
			NewKeyboardButtonRequestLocation("text"),
			NewKeyboardButtonRequestPoll("text", KeyboardButtonPollType{}),
			NewKeyboardButtonWebApp("text", WebAppInfo{}),
		),
	).WithResizeKeyboardMarkup().
		WithOneTimeKeyboardMarkup().
		WithInputFieldPlaceholder("text").
		WithSelective()

	actual.isReplyMarkup()

	assert.EqualValues(t, &ReplyKeyboardMarkup{
		Keyboard: [][]KeyboardButton{
			{
				{Text: "text"},
				{Text: "text", RequestContact: true},
				{Text: "text", RequestLocation: true},
				{Text: "text", RequestPoll: &KeyboardButtonPollType{}},
				{Text: "text", WebApp: &WebAppInfo{}},
			},
		},
		ResizeKeyboard:        true,
		OneTimeKeyboard:       true,
		InputFieldPlaceholder: "text",
		Selective:             true,
	}, actual)
}

func TestReplyKeyboardRemove(t *testing.T) {
	actual := NewReplyKeyboardRemove().WithSelective()

	actual.isReplyMarkup()

	assert.EqualValues(t, &ReplyKeyboardRemove{
		RemoveKeyboard: true,
		Selective:      true,
	}, actual)
}

func TestForceReplay(t *testing.T) {
	actual := NewForceReply().WithSelective().WithInputFieldPlaceholder("test")

	actual.isReplyMarkup()

	assert.EqualValues(t, &ForceReply{
		ForceReply:            true,
		Selective:             true,
		InputFieldPlaceholder: "test",
	}, actual)
}

func TestNewButtonLayout(t *testing.T) {
	keyboard := NewButtonLayout(1,
		NewInlineKeyboardButtonCallback("1", "1"),
		NewInlineKeyboardButtonCallback("2", "2"),
		NewInlineKeyboardButtonCallback("3", "3"),
	).Keyboard()

	assert.Equal(t, [][]InlineKeyboardButton{
		{{Text: "1", CallbackData: "1"}},
		{{Text: "2", CallbackData: "2"}},
		{{Text: "3", CallbackData: "3"}},
	}, keyboard)
}
func TestButtonLayout_Add(t *testing.T) {
	for _, test := range []struct {
		Layout *ButtonLayout[KeyboardButton]
		Want   [][]KeyboardButton
	}{
		{
			Layout: NewButtonLayout[KeyboardButton](3).
				Add(NewKeyboardButton("text")),
			Want: [][]KeyboardButton{
				{{Text: "text"}},
			},
		},
		{
			Layout: NewButtonLayout[KeyboardButton](3).
				Add(NewKeyboardButton("text"), NewKeyboardButton("text"), NewKeyboardButton("text")),
			Want: [][]KeyboardButton{
				{{Text: "text"}, {Text: "text"}, {Text: "text"}},
			},
		},
		{
			Layout: NewButtonLayout[KeyboardButton](3).
				Add(NewKeyboardButton("text"), NewKeyboardButton("text"), NewKeyboardButton("text"), NewKeyboardButton("text")),
			Want: [][]KeyboardButton{
				{{Text: "text"}, {Text: "text"}, {Text: "text"}},
				{{Text: "text"}},
			},
		},
		{
			Layout: NewButtonLayout[KeyboardButton](3).
				Add(NewKeyboardButton("text"), NewKeyboardButton("text"), NewKeyboardButton("text"), NewKeyboardButton("text")).
				Add(NewKeyboardButton("text")),
			Want: [][]KeyboardButton{
				{{Text: "text"}, {Text: "text"}, {Text: "text"}},
				{{Text: "text"}},
				{{Text: "text"}},
			},
		},
	} {
		assert.EqualValues(t, test.Want, test.Layout.Keyboard())
	}
}

func TestButtonLayout_Row(t *testing.T) {
	keyboard := NewButtonLayout(1,
		NewKeyboardButton("1"),
		NewKeyboardButton("2"),
		NewKeyboardButton("3"),
	).Row(
		NewKeyboardButton("4"),
		NewKeyboardButton("5"),
		NewKeyboardButton("6"),
		NewKeyboardButton("7"),
	).Keyboard()

	assert.Equal(t, [][]KeyboardButton{
		{{Text: "1"}},
		{{Text: "2"}},
		{{Text: "3"}},
		{{Text: "4"}, {Text: "5"}, {Text: "6"}, {Text: "7"}},
	}, keyboard)
}

func TestButtonLayout_Insert(t *testing.T) {
	for _, test := range []struct {
		Layout *ButtonLayout[KeyboardButton]
		Want   [][]KeyboardButton
	}{
		{
			Layout: NewButtonLayout[KeyboardButton](3).
				Insert(NewKeyboardButton("text")),
			Want: [][]KeyboardButton{
				{{Text: "text"}},
			},
		},
		{
			Layout: NewButtonLayout[KeyboardButton](3).
				Insert(NewKeyboardButton("text"), NewKeyboardButton("text"), NewKeyboardButton("text")),
			Want: [][]KeyboardButton{
				{{Text: "text"}, {Text: "text"}, {Text: "text"}},
			},
		},
		{
			Layout: NewButtonLayout[KeyboardButton](3).
				Insert(NewKeyboardButton("text"), NewKeyboardButton("text"), NewKeyboardButton("text"), NewKeyboardButton("text")),
			Want: [][]KeyboardButton{
				{{Text: "text"}, {Text: "text"}, {Text: "text"}},
				{{Text: "text"}},
			},
		},
		{
			Layout: NewButtonLayout[KeyboardButton](3).
				Insert(NewKeyboardButton("1"), NewKeyboardButton("2")).
				Insert(NewKeyboardButton("3")),
			Want: [][]KeyboardButton{
				{{Text: "1"}, {Text: "2"}, {Text: "3"}},
			},
		},
		{
			Layout: NewButtonLayout[KeyboardButton](3).
				Insert(NewKeyboardButton("1"), NewKeyboardButton("2")).
				Insert(NewKeyboardButton("3")),
			Want: [][]KeyboardButton{
				{{Text: "1"}, {Text: "2"}, {Text: "3"}},
			},
		},
		{
			Layout: NewButtonLayout[KeyboardButton](3).
				Add(NewKeyboardButton("1"), NewKeyboardButton("2"), NewKeyboardButton("3")).
				Insert(NewKeyboardButton("4")).
				Add(NewKeyboardButton("5")),
			Want: [][]KeyboardButton{
				{{Text: "1"}, {Text: "2"}, {Text: "3"}},
				{{Text: "4"}},
				{{Text: "5"}},
			},
		},
	} {
		assert.EqualValues(t, test.Want, test.Layout.Keyboard())
	}
}

func TestNewButtonColumn(t *testing.T) {
	keyboard := NewButtonColumn(
		NewInlineKeyboardButtonCallback("1", "1"),
		NewInlineKeyboardButtonCallback("2", "2"),
		NewInlineKeyboardButtonCallback("3", "3"),
	)

	assert.Equal(t, [][]InlineKeyboardButton{
		{{Text: "1", CallbackData: "1"}},
		{{Text: "2", CallbackData: "2"}},
		{{Text: "3", CallbackData: "3"}},
	}, keyboard)
}

func TestInlineQueryResultMarshalJSON(t *testing.T) {
	for _, test := range []struct {
		Type   string
		Result InlineQueryResult
	}{
		{"audio", InlineQueryResultCachedAudio{}},

		{"document", InlineQueryResultCachedDocument{}},
		{"gif", InlineQueryResultCachedGIF{}},
		{"mpeg4_gif", InlineQueryResultCachedMPEG4GIF{}},
		{"photo", InlineQueryResultCachedPhoto{}},
		{"sticker", InlineQueryResultCachedSticker{}},
		{"video", InlineQueryResultCachedVideo{}},
		{"voice", InlineQueryResultCachedVoice{}},
		{"audio", InlineQueryResultAudio{}},
		{"document", InlineQueryResultDocument{}},
		{"gif", InlineQueryResultGIF{}},
		{"mpeg4_gif", InlineQueryResultMPEG4GIF{}},
		{"photo", InlineQueryResultPhoto{}},
		{"video", InlineQueryResultVideo{}},
		{"voice", InlineQueryResultVoice{}},
		{"article", InlineQueryResultArticle{}},
		{"contact", InlineQueryResultContact{}},
		{"game", InlineQueryResultGame{}},
		{"location", InlineQueryResultLocation{}},
		{"venue", InlineQueryResultVenue{}},
	} {
		t.Run(test.Type, func(t *testing.T) {
			body, err := json.Marshal(test.Result)
			assert.NoError(t, err, "marshal json")

			test.Result.isInlineQueryResult()

			result := struct {
				Type string `json:"type"`
			}{}

			err = json.Unmarshal(body, &result)
			assert.NoError(t, err, "unmarshal json")

			assert.Equal(t, test.Type, result.Type)
		})

	}
}

func TestInputMessageContent(t *testing.T) {
	for _, test := range []InputMessageContent{
		InputTextMessageContent{},
		InputLocationMessageContent{},
		InputVenueMessageContent{},
		InputContactMessageContent{},
		InputInvoiceMessageContent{},
	} {
		assert.Implements(t, (*InputMessageContent)(nil), test)
		test.isInputMessageContent()
	}
}

func TestInputMedia_getMedia(t *testing.T) {
	for _, test := range []InputMedia{
		&InputMediaPhoto{},
		&InputMediaVideo{},
		&InputMediaAudio{},
		&InputMediaAnimation{},
		&InputMediaDocument{},
	} {
		assert.Implements(t, (*InputMedia)(nil), test)

		media, _ := test.getMedia()

		assert.NotNil(t, media)
	}
}

func TestInputMedia_MarshalJSON(t *testing.T) {
	for _, test := range []struct {
		InputMedia InputMedia
		Want       string
	}{
		{
			InputMedia: &InputMediaPhoto{
				Media: FileArg{FileID: "file_id"},
			},
			Want: `{"type":"photo","media":"file_id"}`,
		},
		{
			InputMedia: &InputMediaVideo{
				Media: FileArg{FileID: "file_id"},
			},
			Want: `{"type":"video","media":"file_id"}`,
		},
		{
			InputMedia: &InputMediaAudio{
				Media: FileArg{FileID: "file_id"},
			},
			Want: `{"type":"audio","media":"file_id"}`,
		},
		{
			InputMedia: &InputMediaAnimation{
				Media: FileArg{FileID: "file_id"},
			},
			Want: `{"type":"animation","media":"file_id"}`,
		},
		{
			InputMedia: &InputMediaDocument{
				Media: FileArg{FileID: "file_id"},
			},
			Want: `{"type":"document","media":"file_id"}`,
		},
	} {
		v, err := json.Marshal(test.InputMedia)
		assert.NoError(t, err, "marshal json")
		assert.Equal(t, test.Want, string(v))
	}
}

func TestFileArg_MarshalJSON(t *testing.T) {
	for _, test := range []struct {
		Name    string
		FileArg FileArg
		Want    string
		Err     bool
	}{
		{
			Name:    "FileID",
			FileArg: FileArg{FileID: "file_id"},
			Want:    `"file_id"`,
		},
		{
			Name:    "FileURL",
			FileArg: FileArg{URL: "file_url"},
			Want:    `"file_url"`,
		},
		{
			Name:    "FileAddr",
			FileArg: FileArg{addr: "addr"},
			Want:    `"addr"`,
		},
		{
			Name: "FileUpload",
			FileArg: FileArg{
				Upload: InputFile{},
			},
			Err: true,
		},
	} {
		t.Run(test.Name, func(t *testing.T) {

			v, err := json.Marshal(test.FileArg)

			if test.Err {
				assert.Error(t, err, "marshal json")
			} else {
				assert.NoError(t, err, "marshal json")
				assert.Equal(t, test.Want, string(v))
			}
		})
	}
}

func TestFileArg_getString(t *testing.T) {
	for _, test := range []struct {
		FileArg FileArg
		Want    string
	}{
		{
			FileArg: FileArg{FileID: "file_id"},
			Want:    "file_id",
		},
		{
			FileArg: FileArg{URL: "file_url"},
			Want:    "file_url",
		},
		{
			FileArg: FileArg{addr: "addr"},
			Want:    "addr",
		},
	} {
		assert.Equal(t, test.Want, test.FileArg.getString())
	}
}

func TestBotCommandScope(t *testing.T) {
	for _, test := range []struct {
		Scope BotCommandScope
		Want  string
	}{
		{BotCommandScopeDefault{}, `{"type":"default"}`},
		{BotCommandScopeAllPrivateChats{}, `{"type":"all_private_chats"}`},
		{BotCommandScopeAllGroupChats{}, `{"type":"all_group_chats"}`},
		{BotCommandScopeAllChatAdministrators{}, `{"type":"all_chat_administrators"}`},
		{BotCommandScopeChat{}, `{"type":"chat","chat_id":0}`},
		{BotCommandScopeChatAdministrators{}, `{"type":"chat_administrators","chat_id":0}`},
		{BotCommandScopeChatMember{}, `{"type":"chat_member","chat_id":0,"user_id":0}`},
	} {
		v, err := json.Marshal(test.Scope)
		assert.NoError(t, err, "marshal json")
		assert.Equal(t, test.Want, string(v))
		test.Scope.isBotCommandScope()
	}
}

func TestMenuButton(t *testing.T) {
	for _, test := range []struct {
		Scope MenuButton
		Want  string
	}{
		{MenuButtonDefault{}, `{"type":"default"}`},
		{MenuButtonCommands{}, `{"type":"commands"}`},
		{MenuButtonWebApp{}, `{"type":"web_app","text":"","web_app":{"url":""}}`},
	} {
		v, err := json.Marshal(test.Scope)
		assert.NoError(t, err, "marshal json")
		assert.Equal(t, test.Want, string(v))
		test.Scope.isMenuButton()
	}
}
