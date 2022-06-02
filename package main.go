package main

import (
	"log"
	"os"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

//Those are keyboards that appear for some of the questions
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

//to operate the bot, put a text file containing key for your bot acquired from telegram "botfather" to the same directory with this file
var tgApiKey, err = os.ReadFile(".secret")
var bot, error1 = tgbotapi.NewBotAPI(string(tgApiKey))

//type containing all the info about user input
type user struct {
	id                int64
	status            int64
	exportTokenName   string
	exportTokenSymbol string
	exportTokenSupply uint64
	exportTokenType   uint64
	tokenTypeString   string
}

//main database, key (int64) is telegram user id
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

	//whenever bot gets a new message, check for user id in the database happens, if it's a new user, the entry in the database is created.
	for update := range updates {

		if update.Message != nil {
			if _, ok := userDatabase[update.Message.From.ID]; !ok {

				userDatabase[update.Message.From.ID] = user{update.Message.Chat.ID, 0, "", "", 0, 0, ""}
				msg := tgbotapi.NewMessage(userDatabase[update.Message.From.ID].id, "Heya, wanna mint your own ERC20, ERC20Snapshot or ERC20Votes? You've come to the right place! Let's begin. Tell me the name of your token!")
				msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
				bot.Send(msg)
			} else {

				//first check for user status, (for a new user status 0 is set automatically), then user reply for the first bot message is logged to a database as name AND user status is updated
				if userDatabase[update.Message.From.ID].status == 0 {
					if updateDb, ok := userDatabase[update.Message.From.ID]; ok {
						updateDb.exportTokenName = update.Message.Text
						updateDb.status = 1
						userDatabase[update.Message.From.ID] = updateDb
					}
					msg := tgbotapi.NewMessage(userDatabase[update.Message.From.ID].id, userDatabase[update.Message.From.ID].exportTokenName+"? That's a cool name! Now tell me the symbol of your token? Usually it's like Bitcoin - BTC, you get the idea")
					bot.Send(msg)

					//logic is that 1 incoming message fro the user equals one status check in database, so each status check ends with the message asking the next question
				} else if userDatabase[update.Message.From.ID].status == 1 {
					if updateDb, ok := userDatabase[update.Message.From.ID]; ok {
						updateDb.exportTokenSymbol = update.Message.Text
						updateDb.status = 2
						userDatabase[update.Message.From.ID] = updateDb
					}
					msg := tgbotapi.NewMessage(userDatabase[update.Message.From.ID].id, userDatabase[update.Message.From.ID].exportTokenSymbol+", alright. Now tell me, what's your desired supply of the tokens?")
					bot.Send(msg)

					//decimals asked, check if user input is uint, token type asked, keyboard is provided
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

					//desired tokentype asked here, it is collected both as string and uint numbers. string is used inside this program, uint is exported. check message asked
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
							updateDb.tokenTypeString = tokenTypeString
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

					//after check message is sent, keyboard is provided. if user answers yes, link to a front-end (WIP) is provided, his entry in the database is deleted, so
					//next time the same user contacts the bot, the process will begin all over again
					//any other answer than "yes" brings the options to correct the info
				} else if userDatabase[update.Message.From.ID].status == 4 {
					if update.Message.Text == "Yes" || update.Message.Text == "It's all correct" {
						msg := tgbotapi.NewMessage(userDatabase[update.Message.From.ID].id, "Here's the link to mint your token!")
						msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
						bot.Send(msg)
						delete(userDatabase, update.Message.From.ID)

					} else {
						if updateDb, ok := userDatabase[update.Message.From.ID]; ok {
							updateDb.status = 5
							userDatabase[update.Message.From.ID] = updateDb
						}
						msg := tgbotapi.NewMessage(userDatabase[update.Message.From.ID].id, "What needs to be corrected?")
						msg.ReplyMarkup = correctKeyboard
						bot.Send(msg)
					}

					//status 5-9 are used for data correction
				} else if userDatabase[update.Message.From.ID].status == 5 {
					if update.Message.Text == "Name" {
						if updateDb, ok := userDatabase[update.Message.From.ID]; ok {
							updateDb.status = 6
							userDatabase[update.Message.From.ID] = updateDb
						}
						msg := tgbotapi.NewMessage(userDatabase[update.Message.From.ID].id, "What's the correct name?")
						msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
						bot.Send(msg)

					} else if update.Message.Text == "Symbol" {
						if updateDb, ok := userDatabase[update.Message.From.ID]; ok {
							updateDb.status = 7
							userDatabase[update.Message.From.ID] = updateDb
						}
						msg := tgbotapi.NewMessage(userDatabase[update.Message.From.ID].id, "What's the correct symbol?")
						msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
						bot.Send(msg)

					} else if update.Message.Text == "Supply" {
						if updateDb, ok := userDatabase[update.Message.From.ID]; ok {
							updateDb.status = 8
							userDatabase[update.Message.From.ID] = updateDb
						}
						msg := tgbotapi.NewMessage(userDatabase[update.Message.From.ID].id, "What's the correct supply?")
						msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
						bot.Send(msg)

					} else if update.Message.Text == "Type" {
						if updateDb, ok := userDatabase[update.Message.From.ID]; ok {
							updateDb.status = 9
							userDatabase[update.Message.From.ID] = updateDb
						}
						msg := tgbotapi.NewMessage(userDatabase[update.Message.From.ID].id, "What's the correct type?")
						msg.ReplyMarkup = numericKeyboard
						bot.Send(msg)

						//keyboard is provided, so whenever user input this one, the link is provided and user entry is deleted from the database
					} else if update.Message.Text == "It's all correct" {
						msg := tgbotapi.NewMessage(userDatabase[update.Message.From.ID].id, "Here's the link to mint your token!")
						msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
						bot.Send(msg)
						delete(userDatabase, update.Message.From.ID)

					} else {
						msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Please, select what needs to be edited!")
						bot.Send(msg)
					}

					//those are used to correct the data. after each correction status is set to 5 AND keyboard to select what needs to be edited is provided,
					//so if something else needs to be corrected, it may be done infinitely, process terminates when the input is "it's all correct"

					//name edit
				} else if userDatabase[update.Message.From.ID].status == 6 {
					if updateDb, ok := userDatabase[update.Message.From.ID]; ok {
						updateDb.exportTokenName = update.Message.Text
						updateDb.status = 5
						userDatabase[update.Message.From.ID] = updateDb
					}

					supplyString := strconv.FormatUint(userDatabase[update.Message.From.ID].exportTokenSupply, 10)

					checkMsgFinal := "Okay, let's check it.\n \n" +
						"Token name: " + userDatabase[update.Message.From.ID].exportTokenName + "\n" +
						"Token symbol: " + userDatabase[update.Message.From.ID].exportTokenSymbol + "\n" +
						"Total supply: " + supplyString + "\n" +
						"Token type: " + userDatabase[update.Message.From.ID].tokenTypeString + "\n \n" +
						"What else needs to be corrected?"
					msg := tgbotapi.NewMessage(userDatabase[update.Message.From.ID].id, checkMsgFinal)
					msg.ReplyMarkup = correctKeyboard
					bot.Send(msg)

					//symbol edit
				} else if userDatabase[update.Message.From.ID].status == 7 {
					if updateDb, ok := userDatabase[update.Message.From.ID]; ok {
						updateDb.exportTokenSymbol = update.Message.Text
						updateDb.status = 5
						userDatabase[update.Message.From.ID] = updateDb
					}

					supplyString := strconv.FormatUint(userDatabase[update.Message.From.ID].exportTokenSupply, 10)

					checkMsgFinal := "Okay, let's check it.\n \n" +
						"Token name: " + userDatabase[update.Message.From.ID].exportTokenName + "\n" +
						"Token symbol: " + userDatabase[update.Message.From.ID].exportTokenSymbol + "\n" +
						"Total supply: " + supplyString + "\n" +
						"Token type: " + userDatabase[update.Message.From.ID].tokenTypeString + "\n \n" +
						"What else needs to be corrected?"

					msg := tgbotapi.NewMessage(userDatabase[update.Message.From.ID].id, checkMsgFinal)
					msg.ReplyMarkup = correctKeyboard
					bot.Send(msg)

					//decimals edit
				} else if userDatabase[update.Message.From.ID].status == 8 {
					TokenSupplyString := update.Message.Text
					tokenSupply, err2 := strconv.ParseUint(TokenSupplyString, 10, 64)
					if err2 == nil {
						if updateDb, ok := userDatabase[update.Message.From.ID]; ok {
							updateDb.exportTokenSupply = tokenSupply
							updateDb.status = 5
							userDatabase[update.Message.From.ID] = updateDb
						}

						supplyString := strconv.FormatUint(userDatabase[update.Message.From.ID].exportTokenSupply, 10)
						checkMsgFinal := "Okay, let's check it.\n \n" +
							"Token name: " + userDatabase[update.Message.From.ID].exportTokenName + "\n" +
							"Token symbol: " + userDatabase[update.Message.From.ID].exportTokenSymbol + "\n" +
							"Total supply: " + supplyString + "\n" +
							"Token type: " + userDatabase[update.Message.From.ID].tokenTypeString + "\n \n" +
							"What else needs to be corrected?"

						msg := tgbotapi.NewMessage(userDatabase[update.Message.From.ID].id, checkMsgFinal)
						msg.ReplyMarkup = correctKeyboard
						bot.Send(msg)
					} else {
						msg := tgbotapi.NewMessage(userDatabase[update.Message.From.ID].id, "Please, enter a number of tokens you want to exist!")
						bot.Send(msg)
					}

					//type edit
				} else if userDatabase[update.Message.From.ID].status == 9 {

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
							updateDb.tokenTypeString = tokenTypeString
							updateDb.status = 5
							userDatabase[update.Message.From.ID] = updateDb
						}

						supplyString := strconv.FormatUint(userDatabase[update.Message.From.ID].exportTokenSupply, 10)
						checkMsgFinal := "Okay, let's check it.\n \n" +
							"Token name: " + userDatabase[update.Message.From.ID].exportTokenName + "\n" +
							"Token symbol: " + userDatabase[update.Message.From.ID].exportTokenSymbol + "\n" +
							"Total supply: " + supplyString + "\n" +
							"Token type: " + userDatabase[update.Message.From.ID].tokenTypeString + "\n \n" +
							"What else needs to be corrected?"

						msg := tgbotapi.NewMessage(userDatabase[update.Message.From.ID].id, checkMsgFinal)
						msg.ReplyMarkup = correctKeyboard
						bot.Send(msg)
					} else {
						msg := tgbotapi.NewMessage(update.Message.Chat.ID, "That's not the type!")
						bot.Send(msg)
					}
				}
			}
		}

	}
}
