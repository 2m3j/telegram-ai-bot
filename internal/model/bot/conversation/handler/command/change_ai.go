package command

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"bot/internal/model/bot/conversation/storage"
	botconvuc "bot/internal/model/bot/conversation/usecase"
	useruc "bot/internal/model/user/usecase"
	botapi "bot/internal/pkg/bot"
)

type ChangeAIHandler struct {
	handler            useruc.ChangeAIHandler
	sendMessageHandler botconvuc.SendMessageHandler
	aiCommandMap       map[string]storage.AICommand
}

func NewChangeAIHandler(handler useruc.ChangeAIHandler, sendMessageHandler botconvuc.SendMessageHandler, aiCommandMap map[string]storage.AICommand) *ChangeAIHandler {
	return &ChangeAIHandler{
		handler:            handler,
		sendMessageHandler: sendMessageHandler,
		aiCommandMap:       aiCommandMap,
	}
}

func (c *ChangeAIHandler) Handle(ctx context.Context, im *botapi.IncomingMessage) error {
	aiCommand := c.aiCommandMap[strings.TrimLeft(im.Message, "/")]
	if aiCommand.Platform == "" || aiCommand.Model == "" {
		aiErr := errors.New("unknown platform or model")
		sendErr := c.sendMessageHandler.Handle(ctx, botconvuc.SendMessage{UserID: im.UserID, ChatID: im.ChatID, Message: "Выбрана неизвестная нейросеть."})
		if sendErr != nil {
			return fmt.Errorf("failed to send unknown AI message: %w", errors.Join(aiErr, sendErr))
		}
		return aiErr
	}
	changeErr := c.handler.Handle(ctx, useruc.ChangeAI{UserID: im.UserID, AIPlatform: aiCommand.Platform, AIModel: aiCommand.Model})
	if changeErr != nil {
		sendErr := c.sendMessageHandler.Handle(ctx, botconvuc.SendMessage{UserID: im.UserID, ChatID: im.ChatID, Message: "Возникла ошибка выбора нейросети, попробуйте повторить выбор снова."})
		if sendErr != nil {
			return fmt.Errorf("failed to send error message after AI change failure: %w", errors.Join(changeErr, sendErr))
		}
		return fmt.Errorf("failed to change AI: %w", changeErr)
	}
	sendErr := c.sendMessageHandler.Handle(ctx, botconvuc.SendMessage{UserID: im.UserID, ChatID: im.ChatID, Message: fmt.Sprintf("Нейросеть %s %s успешно задана.", aiCommand.Platform, aiCommand.Model)})
	if sendErr != nil {
		return fmt.Errorf("failed to send success message: %w", sendErr)
	}
	return nil
}
