package services

type VmessConfigsFormatter struct{}

func (formatter *VmessConfigsFormatter) Format(rawConfig string) (string, error) {
	return " ", nil

}

type VlessConfigsFormatter struct{}

func (formatter *VlessConfigsFormatter) Format(rawConfig string) (string, error) {
	return " ", nil

}

type TrojanConfigsFormatter struct{}

func (formatter *TrojanConfigsFormatter) Format(rawConfig string) (string, error) {
	return " ", nil

}

type SsConfigsFormatter struct{}

func (formatter *SsConfigsFormatter) Format(rawConfig string) (string, error) {
	return " ", nil

}
