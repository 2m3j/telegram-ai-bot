package storage

import (
	"bot/internal/pkg/ai"
	botapi "bot/internal/pkg/bot"
)

const menuRowLength = 2

var MainMenu = map[string]string{
	"/clear": "Очистить контекст",
	"/ai":    "Выбрать нейросеть",
}

var AIMenu []botapi.DialogNode

func init() {
	AIMenu = aIMenu(menuRowLength)
}
func aIMenu(rowLength int) []botapi.DialogNode {
	var deepSeekKeyboard, openAIKeyboard [][]botapi.DialogButton
	deepSeekKeyboardRow := []botapi.DialogButton{{Text: "назад", NodeID: "start"}}
	openAIKeyboardRow := []botapi.DialogButton{{Text: "назад", NodeID: "start"}}
	menu := []botapi.DialogNode{
		{ID: "start", Text: "Выберите нейросеть", Keyboard: [][]botapi.DialogButton{{{Text: ai.PlatformOpenAI, NodeID: "openai"}, {Text: ai.PlatformDeepSeek, NodeID: "deepseek"}}}},
	}
	for name, command := range AICommands {
		if command.Platform == ai.PlatformOpenAI {
			openAIKeyboardRow = append(openAIKeyboardRow, botapi.DialogButton{Text: command.Model, NodeID: name})
			if len(openAIKeyboardRow) == rowLength {
				openAIKeyboard = append(openAIKeyboard, openAIKeyboardRow)
				openAIKeyboardRow = nil
			}
		} else {
			deepSeekKeyboardRow = append(deepSeekKeyboardRow, botapi.DialogButton{Text: command.Model, NodeID: name})
			if len(deepSeekKeyboardRow) == rowLength {
				deepSeekKeyboard = append(deepSeekKeyboard, deepSeekKeyboardRow)
				deepSeekKeyboardRow = nil
			}
		}
		menu = append(menu, botapi.DialogNode{ID: name, Command: "/" + name})
	}
	if openAIKeyboardRow != nil {
		openAIKeyboard = append(openAIKeyboard, openAIKeyboardRow)
	}
	if deepSeekKeyboardRow != nil {
		deepSeekKeyboard = append(deepSeekKeyboard, deepSeekKeyboardRow)
	}
	menu = append(menu, botapi.DialogNode{ID: "openai", Text: "Выберите модель OpenAI", Keyboard: openAIKeyboard})
	menu = append(menu, botapi.DialogNode{ID: "deepseek", Text: "Выберите модель DeepSeek", Keyboard: deepSeekKeyboard})
	return menu
}
