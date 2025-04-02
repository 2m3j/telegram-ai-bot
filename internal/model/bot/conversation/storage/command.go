package storage

import "bot/internal/pkg/ai"

type AICommand struct {
	Platform string
	Model    string
}

var AICommands = map[string]AICommand{
	"aigpt4omini": {
		Platform: ai.PlatformOpenAI,
		Model:    ai.ModelOpenAIGPT4oMini,
	},
	"aio1mini": {
		Platform: ai.PlatformOpenAI,
		Model:    ai.ModelOpenAIO1Mini,
	},
	"aigpt4o": {
		Platform: ai.PlatformOpenAI,
		Model:    ai.ModelOpenAIGPT4o,
	},
	"aideepseekchat": {
		Platform: ai.PlatformDeepSeek,
		Model:    ai.ModelDeepSeekChat,
	},
	/*"aideepseekcoder": {
		Platform: ai.PlatformDeepSeek,
		Model:    ai.ModelDeepSeekCoder,
	},
	"aideepseekreasoner": {
		Platform: ai.PlatformDeepSeek,
		Model:    ai.ModelDeepSeekReasoner,
	},*/
}
