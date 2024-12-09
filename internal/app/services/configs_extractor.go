package services

type ConfigsExtractor struct{}

func (this *ConfigsExtractor) Extract(string) (string, []string) {
	// todo: decode first, if it's base64
	// returns a name for group and list of configs in the given string in raw format
}
