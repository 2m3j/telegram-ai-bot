package ai

type Option func(h *Client)

func WithPlatform(platform PlatformName, client PlatformProvider) Option {
	return func(h *Client) {
		h.platforms[platform] = client
	}
}
func WithErrorHandler(errorHandler ErrorsHandler) Option {
	return func(h *Client) {
		h.errorHandler = errorHandler
	}
}
