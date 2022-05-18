package main

import (
	"log"
	"os"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

//those are variables which we get from User to pass them into a smart-contract
var exportTokenName string
var exportTokenSymbol string
var exportTokenSupply uint64
var exportTokenType uint64

//variable for asking questions to correct data
var ercType string

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

func main() {
	bot, err := tgbotapi.NewBotAPI(string(tgApiKey))
	if err != nil {
		log.Panic(err)
	}

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	//greetings & ask for tokenname
	for update := range updates {
		if update.Message != nil { // If we got a message
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Heya, wanna mint your own ERC20, ERC20Snapshot or ERC20Votes? You've come to a right place! Let's begin. Tell me the name of your token!")
			bot.Send(msg)
		}

		//tokenname export acquired here & symbol asked
		for update := range updates {
			if update.Message != nil { // If we got a message
				tokenname := update.Message.Text
				message2 := tokenname + "? That's a cool name! Now tell me the symbol of your token? Usually it's like Bitcoin - BTC, you get the idea"
				msg1 := tgbotapi.NewMessage(update.Message.Chat.ID, message2)
				bot.Send(msg1)
				exportTokenName = tokenname
				break
			}
		}

		//tokensymbol export acquired here & supply asked
		for update := range updates {
			if update.Message != nil { // If we got a message
				tokensymbol := update.Message.Text
				message2 := tokensymbol + ", alright. Now tell me, what's your desired supply of the tokens?"
				msg1 := tgbotapi.NewMessage(update.Message.Chat.ID, message2)
				bot.Send(msg1)
				exportTokenSymbol = tokensymbol
				break
			}
		}

		//tokensupply acquired here, keyboard for type provided
		for update := range updates {
			if update.Message != nil { // If we got a message
				TokenSupply := update.Message.Text
				var err2 error
				exportTokenSupply, err2 = strconv.ParseUint(TokenSupply, 10, 64)
				if err2 == nil {
					message3 := TokenSupply + " tokens may exist at max, great. Now let's decide about what type of token you want to use - ERC20, ERC20Snapshot or ERC20Votes?"
					msg3 := tgbotapi.NewMessage(update.Message.Chat.ID, message3)
					msg3.ReplyMarkup = numericKeyboard
					bot.Send(msg3)
					break
				} else {
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Please, enter a number of tokens you want to exist!")
					bot.Send(msg)
				}
			}
		}

		//type acquired here, final check asked
		for update := range updates {
			if update.Message != nil {

				if update.Message.Text == "ERC20Snapshot" || update.Message.Text == "ERC20" || update.Message.Text == "ERC20Votes" {
					if update.Message.Text == "ERC20" {
						exportTokenType = 0
						ercType = "ERC20"
					} else if update.Message.Text == "ERC20Snapshot" {
						exportTokenType = 1
						ercType = "ERC20Snapshot"
					} else if update.Message.Text == "ERC20Votes" {
						exportTokenType = 2
						ercType = "ERC20Votes"
					}
					supplyToStr := strconv.FormatUint(exportTokenSupply, 10)
					message4 := "Okay, let's check it.\n \n" +
						"Token name: " + exportTokenName + "\n" +
						"Token symbol: " + exportTokenSymbol + "\n" +
						"Total supply: " + supplyToStr + "\n" +
						"Token type: " + ercType + "\n \n" +
						"Is this right?"

					msg4 := tgbotapi.NewMessage(update.Message.Chat.ID, message4)
					msg4.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
					msg4.ReplyMarkup = yesNoKeyboard
					bot.Send(msg4)
					break
				} else {
					msg4 := tgbotapi.NewMessage(update.Message.Chat.ID, "That's not the type!")
					bot.Send(msg4)
				}
			}

		}

		//final check happens, form to correct appears
		for update := range updates {
			if update.Message != nil {
				var message5 string

				//после вопроса все ли ок -- ответ yes приводит к завершению программы
				if update.Message.Text == "Yes" {
					message5 = "Cool! Here's the link to confirm and mint your token! (хехе, а это я еще не дописал)" //TODO
					msg5 := tgbotapi.NewMessage(update.Message.Chat.ID, message5)
					msg5.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
					bot.Send(msg5)
					supplyToStr := strconv.FormatUint(exportTokenSupply, 10)
					log.Println("\n СОБРАННЫЕ ДАННЫЕ:" + " " + exportTokenName + " " + exportTokenSymbol + " " + supplyToStr + " " + ercType)
					break

					// любой другой приводит к вопросу "что необходимо поменять?"
				} else {

					message5 = "Alright, let's see what's wrong. What do you want to correct?"
					msg5 := tgbotapi.NewMessage(update.Message.Chat.ID, message5)
					msg5.ReplyMarkup = correctKeyboard
					bot.Send(msg5)

					for update := range updates {
						if update.Message != nil && update.Message.Text == "It's all correct" {
							message5 = "Cool! Here's the link to confirm and mint your token! (хехе, а это я еще не дописал)" //TODO
							msg5 := tgbotapi.NewMessage(update.Message.Chat.ID, message5)
							msg5.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
							bot.Send(msg5)
							supplyToStr := strconv.FormatUint(exportTokenSupply, 10)
							typeToStr := strconv.FormatUint(exportTokenType, 10)
							log.Println("\n СОБРАННЫЕ ДАННЫЕ:" + " " + exportTokenName + " " + exportTokenSymbol + " " + supplyToStr + " " + ercType + " " + typeToStr)
							break

						} else {

							switch update.Message.Text {

							case "Name":
								enterName := "What's the correct name?"
								msg := tgbotapi.NewMessage(update.Message.Chat.ID, enterName)
								msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
								bot.Send(msg)
								for update := range updates {
									if update.Message != nil {
										exportTokenName = update.Message.Text
										supplyToStr := strconv.FormatUint(exportTokenSupply, 10)
										checkMsg :=
											"Okay, let's check it.\n \n" +
												"Token name: " + exportTokenName + "\n" +
												"Token symbol: " + exportTokenSymbol + "\n" +
												"Total supply: " + supplyToStr + "\n" +
												"Token type: " + ercType + "\n \n" +
												"Is it all correct or something needs to be changed?"
										msg := tgbotapi.NewMessage(update.Message.Chat.ID, checkMsg)
										msg.ReplyMarkup = correctKeyboard
										bot.Send(msg)
										break
									}
								}

							case "Supply":
								enterSupply := "What's the correct supply?"
								msg := tgbotapi.NewMessage(update.Message.Chat.ID, enterSupply)
								msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
								bot.Send(msg)
								for update := range updates {
									if update.Message != nil {
										TokenSupply := update.Message.Text
										var err2 error
										exportTokenSupply, err2 = strconv.ParseUint(TokenSupply, 10, 64)
										if err2 == nil {
											supplyToStr := strconv.FormatUint(exportTokenSupply, 10)
											checkMsg :=
												"Okay, let's check it.\n \n" +
													"Token name: " + exportTokenName + "\n" +
													"Token symbol: " + exportTokenSymbol + "\n" +
													"Total supply: " + supplyToStr + "\n" +
													"Token type: " + ercType + "\n \n" +
													"Is it all correct or something needs to be changed?"
											msg := tgbotapi.NewMessage(update.Message.Chat.ID, checkMsg)
											msg.ReplyMarkup = correctKeyboard
											bot.Send(msg)
											break
										} else {
											msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Please, enter a number of tokens you want to exist!")
											bot.Send(msg)
										}
									}
								}

							case "Symbol":
								enterSymbol := "What's the correct symbol?"
								msg := tgbotapi.NewMessage(update.Message.Chat.ID, enterSymbol)
								msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
								bot.Send(msg)
								for update := range updates {
									if update.Message != nil {
										exportTokenSymbol = update.Message.Text
										supplyToStr := strconv.FormatUint(exportTokenSupply, 10)
										checkMsg :=
											"Okay, let's check it.\n \n" +
												"Token name: " + exportTokenName + "\n" +
												"Token symbol: " + exportTokenSymbol + "\n" +
												"Total supply: " + supplyToStr + "\n" +
												"Token type: " + ercType + "\n \n" +
												"Is it all correct or something needs to be changed?"
										msg := tgbotapi.NewMessage(update.Message.Chat.ID, checkMsg)
										msg.ReplyMarkup = correctKeyboard
										bot.Send(msg)
										break
									}
								}

							case "Type":
								enterType := "What's the correct type?"
								msg := tgbotapi.NewMessage(update.Message.Chat.ID, enterType)
								msg.ReplyMarkup = numericKeyboard
								bot.Send(msg)
								for update := range updates {
									if update.Message != nil {
										if update.Message.Text == "ERC20Snapshot" || update.Message.Text == "ERC20" || update.Message.Text == "ERC20Votes" {

											if update.Message.Text == "ERC20" {
												exportTokenType = 0
												ercType = "ERC20"
											} else if update.Message.Text == "ERC20Snapshot" {
												exportTokenType = 1
												ercType = "ERC20Snapshot"
											} else if update.Message.Text == "ERC20Votes" {
												exportTokenType = 2
												ercType = "ERC20Votes"
											}

											supplyToStr := strconv.FormatUint(exportTokenSupply, 10)

											checkMsg :=
												"Okay, let's check it.\n \n" +
													"Token name: " + exportTokenName + "\n" +
													"Token symbol: " + exportTokenSymbol + "\n" +
													"Total supply: " + supplyToStr + "\n" +
													"Token type: " + ercType + "\n \n" +
													"Is it all correct or something needs to be changed?"
											msg := tgbotapi.NewMessage(update.Message.Chat.ID, checkMsg)
											msg.ReplyMarkup = correctKeyboard
											bot.Send(msg)
											break
										} else {
											msg4 := tgbotapi.NewMessage(update.Message.Chat.ID, "That's not the type!")
											bot.Send(msg4)
										}
									}
								}

							case "Путин":
								msg := tgbotapi.NewMessage(update.Message.Chat.ID, "хуйло, конечно же :3")
								bot.Send(msg)

							default:
								msg := tgbotapi.NewMessage(update.Message.Chat.ID, "That's not a valid command!")
								bot.Send(msg)
							}

						}

					}

				}

				break
			}
		}
	}
}
