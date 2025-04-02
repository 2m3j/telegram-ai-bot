package bot

import (
	"log/slog"
	"os"
	"time"

	"bot/config"
	cnvrepo "bot/internal/model/ai/conversation/repository"
	aiusecase "bot/internal/model/ai/conversation/usecase"
	botconvhdlr "bot/internal/model/bot/conversation/handler"
	botcommand "bot/internal/model/bot/conversation/handler/command"
	botconvstorage "bot/internal/model/bot/conversation/storage"
	botusecase "bot/internal/model/bot/conversation/usecase"
	userrepo "bot/internal/model/user/repository"
	userusecase "bot/internal/model/user/usecase"
	aiapi "bot/internal/pkg/ai"
	botapi "bot/internal/pkg/bot"
	"bot/internal/pkg/bot/telegram"
	"bot/internal/pkg/db/mysql"
	"bot/internal/pkg/decorator"
	"bot/internal/pkg/metrics"
	trmsql "github.com/avito-tech/go-transaction-manager/drivers/sql/v2"
	"github.com/avito-tech/go-transaction-manager/trm/v2"
	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"
)

type Dependencies struct {
	bot           botapi.Bot
	botCnvHandler *botconvhdlr.Handler
}

func NewDependencies(cfg *config.Config, mysqlClient *mysql.Client) *Dependencies {
	trManager := manager.Must(
		trmsql.NewDefaultFactory(mysqlClient.Pool),
	)
	aiConvRepository := cnvrepo.NewConversationMysqlRepository(mysqlClient, trmsql.DefaultCtxGetter)
	aiMsgRepository := cnvrepo.NewMessageMysqlRepository(mysqlClient, trmsql.DefaultCtxGetter)
	userRepository := userrepo.NewUserMysqlRepository(mysqlClient, trmsql.DefaultCtxGetter)
	logger := newLogger(cfg.App.Env)
	errorsHandler := func(err error) {
		logger.Error(err.Error())
	}
	bot := telegram.NewBot(cfg.Telegram.Token, errorsHandler, botconvstorage.MainMenu)
	metricsClient := metrics.NoOp{}
	return &Dependencies{
		bot:           bot,
		botCnvHandler: newBotConversationHandler(trManager, aiConvRepository, aiMsgRepository, userRepository, cfg.Bot, bot, newAIClient(cfg.AI, errorsHandler), errorsHandler, logger, metricsClient),
	}
}

func newAiSendConvHandler(
	convRepository cnvrepo.ConversationRepository,
	msgRepository cnvrepo.MessageRepository,
	bot botapi.Bot,
	ai *aiapi.Client,
	logger *slog.Logger,
	metrics decorator.MetricsClient,
) aiusecase.SendConversationMessageHandler {
	h := aiusecase.NewSendConversationMessageHandler(convRepository, msgRepository, bot, ai, logger, metrics)
	return h
}
func newAiCreateConvHandler(
	convRepository cnvrepo.ConversationRepository,
	msgRepository cnvrepo.MessageRepository,
	trm trm.Manager,
	bot botapi.Bot,
	ai *aiapi.Client,
	logger *slog.Logger,
	metrics decorator.MetricsClient,
) aiusecase.StartConversationHandler {
	h := aiusecase.NewStartConversationHandler(convRepository, msgRepository, trm, bot, ai, logger, metrics)
	return h
}
func newUserRegistrationHandler(
	userRepository userrepo.UserRepository,
	logger *slog.Logger,
	metrics decorator.MetricsClient,
) userusecase.RegistrationHandler {
	h := userusecase.NewRegistrationHandler(userRepository, logger, metrics)
	return h
}

func newLogger(env string) *slog.Logger {
	opts := &slog.HandlerOptions{}
	if env == "dev" {
		opts.Level = slog.LevelDebug
	} else {
		opts.Level = slog.LevelInfo
	}
	var handler slog.Handler = slog.NewJSONHandler(os.Stdout, opts)
	return slog.New(handler)
}

func newAIClient(cfg config.AI, errorHandler aiapi.ErrorsHandler) *aiapi.Client {
	openAIPlatform := aiapi.NewOpenAIProvider(cfg.OpenAI.Token)
	deepSeekPlatform := aiapi.NewDeepSeekProvider(cfg.DeepSeek.Token)
	opts := []aiapi.Option{
		aiapi.WithPlatform(aiapi.PlatformOpenAI, openAIPlatform),
		aiapi.WithPlatform(aiapi.PlatformDeepSeek, deepSeekPlatform),
		aiapi.WithErrorHandler(errorHandler),
	}
	return aiapi.NewClient(opts...)
}

func newBotConversationHandler(
	trm trm.Manager,
	convRepository cnvrepo.ConversationRepository,
	msgRepository cnvrepo.MessageRepository,
	userRepository userrepo.UserRepository,
	cfg config.Bot,
	bot botapi.Bot,
	aiClient *aiapi.Client,
	errorsHandler botconvhdlr.ErrorsHandler,
	logger *slog.Logger,
	metrics decorator.MetricsClient,
) *botconvhdlr.Handler {

	sendMessageHandler := botusecase.NewSendMessageHandler(bot, logger, metrics)
	changeAICommandsHandler := botcommand.NewChangeAIHandler(userusecase.NewChangeAIHandler(userRepository, logger, metrics), sendMessageHandler, botconvstorage.AICommands)

	commandHandlers := map[string]botconvhdlr.CommandHandler{
		"/start": botcommand.NewStartConversationHandler(botusecase.NewStartConvHandler(bot, logger, metrics)),
		"/clear": botcommand.NewClearConversationContextHandler(aiusecase.NewEndConversationsHandler(trm, bot, convRepository, logger, metrics)),
		"/ai":    botcommand.NewShowAIHandler(bot),
	}
	for changeAIc := range botconvstorage.AICommands {
		commandHandlers["/"+changeAIc] = changeAICommandsHandler
	}
	h := botconvhdlr.NewHandler(
		time.Duration(cfg.SkipMessageTimeout)*time.Second,
		botconvhdlr.NewMessageHandler(
			newAiCreateConvHandler(convRepository, msgRepository, trm, bot, aiClient, logger, metrics),
			newAiSendConvHandler(convRepository, msgRepository, bot, aiClient, logger, metrics),
			convRepository,
		),
		botconvhdlr.NewMessageCommandHandler(commandHandlers, botcommand.NewSendUnknownHandler(sendMessageHandler)),
		botconvhdlr.NewUserProcess(cfg.DefaultAIPlatform, cfg.DefaultAIModel, userRepository, newUserRegistrationHandler(userRepository, logger, metrics)),
		errorsHandler,
	)
	return h
}
