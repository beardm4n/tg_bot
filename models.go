package main

type Update struct {
	UpdateId int     `json:"update_id"`
	Message  Message `json:"message"`
}

type Message struct {
	MessageId int    `json:"message_id"`
	Chat      Chat   `json:"chat"`
	Text      string `json:"text"`
	Voice     Voice  `json:"voice"`
}

type Chat struct {
	Id int `json:"id"`
}

type MessageResponse struct {
	Result []Update `json:"result"`
}

type VoiceResponse struct {
	Result File `json:"result"`
}

type BotMessage struct {
	ChatId      int                 `json:"chat_id"`
	Text        string              `json:"text"`
	ReplyMarkup ReplyKeyboardMarkup `json:"reply_markup"`
}

type ReplyKeyboardMarkup struct {
	Keyboard        [][]KeyboardButton `json:"keyboard"`
	OneTimeKeyboard bool               `json:"one_time_keyboard,omitempty"`
}

type KeyboardButton struct {
	Text            string `json:"text"`
	RequestContact  bool   `json:"request_contact,omitempty"`
	RequestLocation bool   `json:"request_location,omitempty"`
}

type Voice struct {
	FileId     string `json:"file_id"`
	FileUniqId string `json:"file_unique_id"`
	Duration   int    `json:"duration"`
}

type File struct {
	FileId     string `json:"file_id"`
	FileUniqId string `json:"file_unique_id"`
	FileSize   int    `json:"file_size"`
	FilePath   string `json:"file_path"`
}
