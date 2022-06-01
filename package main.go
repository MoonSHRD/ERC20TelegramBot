package main

import (
	"log"
	"os"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var numericKeyboard = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("ERC20"),
		tgbotapi.NewKeyboardButton("ERC20Snapshot"),
	),
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("ERC20Votes")),
)

var yesNoKeyboard = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Yes"),
		tgbotapi.NewKeyboardButton("No")),
)

var correctKeyboard = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Name")),

	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Symbol"),
		tgbotapi.NewKeyboardButton("Supply"),
		tgbotapi.NewKeyboardButton("Type")),

	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("It's all correct"),
	),
)

var tgApiKey, err = os.ReadFile(".secret")

var bot, error1 = tgbotapi.NewBotAPI(string(tgApiKey))

type user struct {
	id                int64
	status            int64
	exportTokenName   string
	exportTokenSymbol string
	exportTokenSupply uint64
	exportTokenType   uint64
}

var userDatabase = make(map[int64]user)

func main() {

	bot, err = tgbotapi.NewBotAPI(string(tgApiKey))
	if err != nil {
		log.Panic(err)
	}

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	//greetings & ask for tokenname
	for update := range updates {
		if update.Message != nil {
			if _, ok := userDatabase[update.Message.From.ID]; ok {
				if userDatabase[update.Message.From.ID].status == 0 {
					if updateDb, ok := userDatabase[update.Message.From.ID]; ok {
						updateDb.exportTokenName = update.Message.Text
						updateDb.status = 1
						userDatabase[update.Message.From.ID] = updateDb
					}
					msg := tgbotapi.NewMessage(userDatabase[update.Message.From.ID].id, userDatabase[update.Message.From.ID].exportTokenName+"? That's a cool name! Now tell me the symbol of your token? Usually it's like Bitcoin - BTC, you get the idea")
					bot.Send(msg)

				} else if userDatabase[update.Message.From.ID].status == 1 {
					if updateDb, ok := userDatabase[update.Message.From.ID]; ok {
						updateDb.exportTokenSymbol = update.Message.Text
						updateDb.status = 2
						userDatabase[update.Message.From.ID] = updateDb
					}
					msg := tgbotapi.NewMessage(userDatabase[update.Message.From.ID].id, userDatabase[update.Message.From.ID].exportTokenSymbol+", alright. Now tell me, what's your desired supply of the tokens?")
					bot.Send(msg)

				} else if userDatabase[update.Message.From.ID].status == 2 {
					TokenSupplyString := update.Message.Text
					tokenSupply, err2 := strconv.ParseUint(TokenSupplyString, 10, 64)
					if err2 == nil {
						if updateDb, ok := userDatabase[update.Message.From.ID]; ok {
							updateDb.exportTokenSupply = tokenSupply
							updateDb.status = 3
							userDatabase[update.Message.From.ID] = updateDb
						}
						msg := tgbotapi.NewMessage(userDatabase[update.Message.From.ID].id, TokenSupplyString+" tokens may exist at max, great. Now let's decide about what type of token you want to use - ERC20, ERC20Snapshot or ERC20Votes?")
						msg.ReplyMarkup = numericKeyboard
						bot.Send(msg)
					} else {
						msg := tgbotapi.NewMessage(userDatabase[update.Message.From.ID].id, "Please, enter a number of tokens you want to exist!")
						bot.Send(msg)
					}

				} else if userDatabase[update.Message.From.ID].status == 3 {

					if update.Message.Text == "ERC20Snapshot" || update.Message.Text == "ERC20" || update.Message.Text == "ERC20Votes" {

						var tokenType uint64
						var tokenTypeString string

						if update.Message.Text == "ERC20" {
							tokenType = 0
							tokenTypeString = "ERC20"
						} else if update.Message.Text == "ERC20Snapshot" {
							tokenType = 1
							tokenTypeString = "ERC20Snapshot"
						} else if update.Message.Text == "ERC20Votes" {
							tokenType = 2
							tokenTypeString = "ERC20Votes"
						}

						if updateDb, ok := userDatabase[update.Message.From.ID]; ok {
							updateDb.exportTokenType = tokenType
							updateDb.status = 4
							userDatabase[update.Message.From.ID] = updateDb
						}

						supplyString := strconv.FormatUint(userDatabase[update.Message.From.ID].exportTokenSupply, 10)

						checkMsg := "Okay, let's check it.\n \n" +
							"Token name: " + userDatabase[update.Message.From.ID].exportTokenName + "\n" +
							"Token symbol: " + userDatabase[update.Message.From.ID].exportTokenSymbol + "\n" +
							"Total supply: " + supplyString + "\n" +
							"Token type: " + tokenTypeString + "\n \n" +
							"Is this right?"

						msg := tgbotapi.NewMessage(userDatabase[update.Message.From.ID].id, checkMsg)
						msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
						msg.ReplyMarkup = yesNoKeyboard
						bot.Send(msg)

					} else {
						msg := tgbotapi.NewMessage(update.Message.Chat.ID, "That's not the type!")
						bot.Send(msg)
					}

				} else if userDatabase[update.Message.From.ID].status == 4 {
					if update.Message.Text == "Yes" {
						msg := tgbotapi.NewMessage(userDatabase[update.Message.From.ID].id, "Here's the link to mint your token!")
						msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
						bot.Send(msg)
						delete(userDatabase, update.Message.From.ID)
					}

					if update.Message.Text == "No" { //TODO
					}

				}

			} else {
				userDatabase[update.Message.From.ID] = user{update.Message.Chat.ID, 0, "", "", 0, 0}
				msg := tgbotapi.NewMessage(userDatabase[update.Message.From.ID].id, "Heya, wanna mint your own ERC20, ERC20Snapshot or ERC20Votes? You've come to a right place! Let's begin. Tell me the name of your token!")
				msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
				bot.Send(msg)
			}
		}

	}
}
