package main

type Update struct {
	UpdateId int     `json:"update_id"`
	Message  Message `json:"message"`
}

type Message struct {
	Chat Chat   `json:"chat"`
	Text string `json:"text"`
}

type Chat struct {
	Id int `json:"id"`
}

type RestResponse struct {
	Result []Update `json:"result"`
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
