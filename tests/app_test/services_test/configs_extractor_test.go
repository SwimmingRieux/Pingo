package services_test

import (
	"github.com/stretchr/testify/assert"
	"pingo/internal/app/services"
	"testing"
	"time"
)

type configsExtractorTest struct {
	name               string
	text               string
	expectedName       string
	expectedRawConfigs []string
}

var configsExtractorTests = []configsExtractorTest{
	{
		name:         "should return valid configs and current timestamp when input contains multiple valid configs",
		text:         "ws\nLocation: 🇩🇪\n\nتمامی نت ها\n\nvless://8778b6fb-3547-47a9-8494-5b93e9fd5972@EXPRESSVPN_420.fastly80-3.hosting-ip.com:80?security=none&type=ws&headerType=&path=%2Ftelegram-EXPRESSVPN_420%2Ctelegram-EXPRESSVPN_420%2Ctelegram-EXPRESSVPN_420%2Ctelegram-EXPRESSVPN_420%2Ctelegram-EXPRESSVPN_420%2Ctelegram-EXPRESSVPN_420%2Ctelegram-EXPRESSVPN_420%2Ctelegram-EXPRESSVPN_420%3Fed%3D8080&host=Digikala.iranian.net#%40EXPRESSVPN_420%20%F0%9F%87%A9%F0%9F%87%AAAll2\n\nتمامی نت ها\n\n\nvless://8778b6fb-3547-47a9-8494-5b93e9fd5972@EXPRESSVPN_420.fast.hosting-ip.com:80?security=none&type=ws&headerType=&path=%2Ftelegram-EXPRESSVPN_420%2Ctelegram-EXPRESSVPN_420%2Ctelegram-EXPRESSVPN_420%2Ctelegram-EXPRESSVPN_420%2Ctelegram-EXPRESSVPN_420%2Ctelegram-EXPRESSVPN_420%2Ctelegram-EXPRESSVPN_420%2Ctelegram-EXPRESSVPN_420%3Fed%3D8080&host=Digikala.iranian.net#%40EXPRESSVPN_420%20-All\n\nتمامی نت ها\n\nvless://8778b6fb-3547-47a9-8494-5b93e9fd5972@speedtest.net:80?security=none&type=ws&headerType=&path=%2Ftelegram-EXPRESSVPN_420%2Ctelegram-EXPRESSVPN_420%2Ctelegram-EXPRESSVPN_420%2Ctelegram-EXPRESSVPN_420%2Ctelegram-EXPRESSVPN_420%2Ctelegram-EXPRESSVPN_420%2Ctelegram-EXPRESSVPN_420%2Ctelegram-EXPRESSVPN_420%3Fed%3D8080&host=Digikala.iranian.net#%40EXPRESSVPN_420%20%F0%9F%87%A9%F0%9F%87%AAall\n\nتیک Mux و Fragment خاموش باشه!",
		expectedName: time.Now().Format("Jan 02 2006 15:04"),
		expectedRawConfigs: []string{
			"vless://8778b6fb-3547-47a9-8494-5b93e9fd5972@EXPRESSVPN_420.fastly80-3.hosting-ip.com:80?security=none&type=ws&headerType=&path=%2Ftelegram-EXPRESSVPN_420%2Ctelegram-EXPRESSVPN_420%2Ctelegram-EXPRESSVPN_420%2Ctelegram-EXPRESSVPN_420%2Ctelegram-EXPRESSVPN_420%2Ctelegram-EXPRESSVPN_420%2Ctelegram-EXPRESSVPN_420%2Ctelegram-EXPRESSVPN_420%3Fed%3D8080&host=Digikala.iranian.net#%40EXPRESSVPN_420%20%F0%9F%87%A9%F0%9F%87%AAAll2",
			"vless://8778b6fb-3547-47a9-8494-5b93e9fd5972@EXPRESSVPN_420.fast.hosting-ip.com:80?security=none&type=ws&headerType=&path=%2Ftelegram-EXPRESSVPN_420%2Ctelegram-EXPRESSVPN_420%2Ctelegram-EXPRESSVPN_420%2Ctelegram-EXPRESSVPN_420%2Ctelegram-EXPRESSVPN_420%2Ctelegram-EXPRESSVPN_420%2Ctelegram-EXPRESSVPN_420%2Ctelegram-EXPRESSVPN_420%3Fed%3D8080&host=Digikala.iranian.net#%40EXPRESSVPN_420%20-All",
			"vless://8778b6fb-3547-47a9-8494-5b93e9fd5972@speedtest.net:80?security=none&type=ws&headerType=&path=%2Ftelegram-EXPRESSVPN_420%2Ctelegram-EXPRESSVPN_420%2Ctelegram-EXPRESSVPN_420%2Ctelegram-EXPRESSVPN_420%2Ctelegram-EXPRESSVPN_420%2Ctelegram-EXPRESSVPN_420%2Ctelegram-EXPRESSVPN_420%2Ctelegram-EXPRESSVPN_420%3Fed%3D8080&host=Digikala.iranian.net#%40EXPRESSVPN_420%20%F0%9F%87%A9%F0%9F%87%AAall",
		},
	},
	{
		name:         "should return valid configs and current timestamp when input contains mixed valid and invalid configs",
		text:         "همه اپراتورها📶📶📶🛜\n\n🗺لوکشین: آلمان🇩🇪\n\nvless://d342d11e-d424-4583-b36e-524ab1f0afa4@zula.ir:80?path=tmevpnAndroid2%2F%3Fed%3D2560&security=none&encryption=none&host=a.xn--i-sx6a60i.us.kg.&type=ws#%E2%87%9D%20%40v2mystery%20%E2%80%A2%20%F0%9F%8C%BB%D9%82%D8%B7%D8%B9%20%D8%B4%D8%AF%DB%8C%20%D8%A8%DB%8C%D8%A7\n\nvmess://eyJhZGQiOiIxMDQuMTkuMTUwLjEwIiwiYWlkIjoiMCIsImFscG4iOiIiLCJmcCI6IiIsImhvc3QiOiJtY2kubXRuLmlyLmNvbS5vcmcubmV0Lm9tLm11c2ljZmExMjMuaXIuIiwiaWQiOiI5NGYzMzJiMC1jNWQzLTQ1MzEtYTFkNi02ZTYzNThjYzZjNzIiLCJuZXQiOiJodHRwdXBncmFkZSIsInBhdGgiOiIvZDNkM0xtbHlZVzVvYjNOMExtTnZiUVx1MDAzZFx1MDAzZD9lZFx1MDAzZDI1NjAiLCJwb3J0IjoiMjA5NSIsInBzIjoi4oedIEB2Mm15c3Rlcnkg4oCiIPCfjLvZgti32Lkg2LTYr9uMINio24zYpyIsInNjeSI6ImF1dG8iLCJzbmkiOiIiLCJ0bHMiOiIiLCJ0eXBlIjoiLS0tIiwidiI6IjIifQ==\nvless://1ca553ab-29e6-480e-b226-ad72493ba0e1@EXPRESSVPN_420.fastly80-3.hosting-ip.com:80?path=%2Ftelegram-EXPRESSVPN_420%2Ctelegram-EXPRESSVPN_420%2Ctelegram-EXPRESSVPN_420%2Ctelegram-EXPRESSVPN_420%2Ctelegram-EXPRESSVPN_420%2Ctelegram-EXPRESSVPN_420%2Ctelegram-EXPRESSVPN_420%2Ctelegram-EXPRESSVPN_420%3Fed%3D8080&security=none&encryption=none&host=Digikala.iranian.net&type=ws#%E2%87%9D%20%40v2mystery%20%E2%80%A2%20%F0%9F%8C%BB%D9%82%D8%B7%D8%B9%20%D8%B4%D8%AF%DB%8C%20%D8%A8%DB%8C%D8%A7\n\nvless://1ca553ab-29e6-480e-b226-ad72493ba0e1@speedtest.net:80?path=%2Ftelegram-EXPRESSVPN_420%2Ctelegram-EXPRESSVPN_420%2Ctelegram-EXPRESSVPN_420%2Ctelegram-EXPRESSVPN_420%2Ctelegram-EXPRESSVPN_420%2Ctelegram-EXPRESSVPN_420%2Ctelegram-EXPRESSVPN_420%2Ctelegram-EXPRESSVPN_420%3Fed%3D8080&security=none&encryption=none&host=Digikala.iranian.net&type=ws#%E2%87%9D%20%40v2mystery%20%E2%80%A2%20%F0%9F%8C%BB%D9%82%D8%B7%D8%B9%20%D8%B4%D8%AF%DB%8C%20%D8%A8%DB%8C%D8%A7\n\nوضعیت: 🔜فعال تا زمان نامشخص⏳🔙",
		expectedName: time.Now().Format("Jan 02 2006 15:04"),
		expectedRawConfigs: []string{
			"vless://d342d11e-d424-4583-b36e-524ab1f0afa4@zula.ir:80?path=tmevpnAndroid2%2F%3Fed%3D2560&security=none&encryption=none&host=a.xn--i-sx6a60i.us.kg.&type=ws#%E2%87%9D%20%40v2mystery%20%E2%80%A2%20%F0%9F%8C%BB%D9%82%D8%B7%D8%B9%20%D8%B4%D8%AF%DB%8C%20%D8%A8%DB%8C%D8%A7",
			"vmess://eyJhZGQiOiIxMDQuMTkuMTUwLjEwIiwiYWlkIjoiMCIsImFscG4iOiIiLCJmcCI6IiIsImhvc3QiOiJtY2kubXRuLmlyLmNvbS5vcmcubmV0Lm9tLm11c2ljZmExMjMuaXIuIiwiaWQiOiI5NGYzMzJiMC1jNWQzLTQ1MzEtYTFkNi02ZTYzNThjYzZjNzIiLCJuZXQiOiJodHRwdXBncmFkZSIsInBhdGgiOiIvZDNkM0xtbHlZVzVvYjNOMExtTnZiUVx1MDAzZFx1MDAzZD9lZFx1MDAzZDI1NjAiLCJwb3J0IjoiMjA5NSIsInBzIjoi4oedIEB2Mm15c3Rlcnkg4oCiIPCfjLvZgti32Lkg2LTYr9uMINio24zYpyIsInNjeSI6ImF1dG8iLCJzbmkiOiIiLCJ0bHMiOiIiLCJ0eXBlIjoiLS0tIiwidiI6IjIifQ==",
			"vless://1ca553ab-29e6-480e-b226-ad72493ba0e1@EXPRESSVPN_420.fastly80-3.hosting-ip.com:80?path=%2Ftelegram-EXPRESSVPN_420%2Ctelegram-EXPRESSVPN_420%2Ctelegram-EXPRESSVPN_420%2Ctelegram-EXPRESSVPN_420%2Ctelegram-EXPRESSVPN_420%2Ctelegram-EXPRESSVPN_420%2Ctelegram-EXPRESSVPN_420%2Ctelegram-EXPRESSVPN_420%3Fed%3D8080&security=none&encryption=none&host=Digikala.iranian.net&type=ws#%E2%87%9D%20%40v2mystery%20%E2%80%A2%20%F0%9F%8C%BB%D9%82%D8%B7%D8%B9%20%D8%B4%D8%AF%DB%8C%20%D8%A8%DB%8C%D8%A7",
			"vless://1ca553ab-29e6-480e-b226-ad72493ba0e1@speedtest.net:80?path=%2Ftelegram-EXPRESSVPN_420%2Ctelegram-EXPRESSVPN_420%2Ctelegram-EXPRESSVPN_420%2Ctelegram-EXPRESSVPN_420%2Ctelegram-EXPRESSVPN_420%2Ctelegram-EXPRESSVPN_420%2Ctelegram-EXPRESSVPN_420%2Ctelegram-EXPRESSVPN_420%3Fed%3D8080&security=none&encryption=none&host=Digikala.iranian.net&type=ws#%E2%87%9D%20%40v2mystery%20%E2%80%A2%20%F0%9F%8C%BB%D9%82%D8%B7%D8%B9%20%D8%B4%D8%AF%DB%8C%20%D8%A8%DB%8C%D8%A7",
		},
	},
	{
		name:               "should return empty configs and current timestamp when input contains no valid configs",
		text:               "This is a test string with no valid configs.",
		expectedName:       time.Now().Format("Jan 02 2006 15:04"),
		expectedRawConfigs: []string{},
	},
	{
		name:               "should return empty configs and current timestamp when input is empty",
		text:               "",
		expectedName:       time.Now().Format("Jan 02 2006 15:04"),
		expectedRawConfigs: []string{},
	},
	{
		name:         "should return valid configs and current timestamp when input contains only one valid config",
		text:         "vless://8778b6fb-3547-47a9-8494-5b93e9fd5972@EXPRESSVPN_420.fastly80-3.hosting-ip.com:80?security=none&type=ws&headerType=&path=%2Ftelegram-EXPRESSVPN_420%2Ctelegram-EXPRESSVPN_420%2Ctelegram-EXPRESSVPN_420%2Ctelegram-EXPRESSVPN_420%2Ctelegram-EXPRESSVPN_420%2Ctelegram-EXPRESSVPN_420%2Ctelegram-EXPRESSVPN_420%2Ctelegram-EXPRESSVPN_420%3Fed%3D8080&host=Digikala.iranian.net#%40EXPRESSVPN_420%20%F0%9F%87%A9%F0%9F%87%AAAll2",
		expectedName: time.Now().Format("Jan 02 2006 15:04"),
		expectedRawConfigs: []string{
			"vless://8778b6fb-3547-47a9-8494-5b93e9fd5972@EXPRESSVPN_420.fastly80-3.hosting-ip.com:80?security=none&type=ws&headerType=&path=%2Ftelegram-EXPRESSVPN_420%2Ctelegram-EXPRESSVPN_420%2Ctelegram-EXPRESSVPN_420%2Ctelegram-EXPRESSVPN_420%2Ctelegram-EXPRESSVPN_420%2Ctelegram-EXPRESSVPN_420%2Ctelegram-EXPRESSVPN_420%2Ctelegram-EXPRESSVPN_420%3Fed%3D8080&host=Digikala.iranian.net#%40EXPRESSVPN_420%20%F0%9F%87%A9%F0%9F%87%AAAll2",
		},
	},
	{
		name:         "should return valid configs and current timestamp when input contains base64 encoded configs",
		text:         "dmxlc3M6Ly9mZTMwODBiYS1iOWY2LTRkMzAtOWQ0My1jYTYyMWJmZjE4OWZAYzIudmFjdC5pcjoyMDgzP21vZGU9Z3VuJnNlY3VyaXR5PXJlYWxpdHkmZW5jcnlwdGlvbj1ub25lJnBiaz1nbWZaeVF0b1hYcmNrNDVra1htX3ZPM1RsQjBDY0ROdmtMa2QzLUhBZkZBJmZwPWNocm9tZSZzcHg9JTJGJnR5cGU9Z3JwYyZzZXJ2aWNlTmFtZT1UZWxlZ3JhbS1pcFYyUmF5LVRlbGVncmFtLWlwVjJSYXktVGVsZWdyYW0taXBWMlJheS1UZWxlZ3JhbS1pcFYyUmF5LVRlbGVncmFtLWlwVjJSYXkmc25pPXd3dy5jZG43Ny5jb20mc2lkPWI1NDJhYyMlNDB2MnJheTMxMyUyMCVEQSVBOSVEOCVBNyVEOSU4NiVEOCVBNyVEOSU4NCVEOSU4NSVEOSU4OCVEOSU4NiUyMCVGMCU5RiU5OCU4RA==",
		expectedName: time.Now().Format("Jan 02 2006 15:04"),
		expectedRawConfigs: []string{
			"vless://fe3080ba-b9f6-4d30-9d43-ca621bff189f@c2.vact.ir:2083?mode=gun&security=reality&encryption=none&pbk=gmfZyQtoXXrck45kkXm_vO3TlB0CcDNvkLkd3-HAfFA&fp=chrome&spx=%2F&type=grpc&serviceName=Telegram-ipV2Ray-Telegram-ipV2Ray-Telegram-ipV2Ray-Telegram-ipV2Ray-Telegram-ipV2Ray&sni=www.cdn77.com&sid=b542ac#%40v2ray313%20%DA%A9%D8%A7%D9%86%D8%A7%D9%84%D9%85%D9%88%D9%86%20%F0%9F%98%8D",
		},
	},
	{
		name:               "should return empty configs and current timestamp when input contains invalid base64 encoded configs",
		text:               "dGhpcyBpcyBhIHRlc3Qgc3RyaW5nIHdpdGggaW52YWxpZCBiYXNlNjQgZW5jb2RlZCBjb25maWdzLg==",
		expectedName:       time.Now().Format("Jan 02 2006 15:04"),
		expectedRawConfigs: []string{},
	},
	{
		name:               "should return empty configs and current timestamp when input contains only invalid configs",
		text:               "invalid://config\nanotherinvalid://config",
		expectedName:       time.Now().Format("Jan 02 2006 15:04"),
		expectedRawConfigs: []string{},
	},
	{
		name:         "should return valid configs and current timestamp when input contains valid configs with special characters",
		text:         "vless://8778b6fb-3547-47a9-8494-5b93e9fd5972@EXPRESSVPN_420.fastly80-3.hosting-ip.com:80?security=none&type=ws&headerType=&path=%2Ftelegram-EXPRESSVPN_420%2Ctelegram-EXPRESSVPN_420%2Ctelegram-EXPRESSVPN_420%2Ctelegram-EXPRESSVPN_420%2Ctelegram-EXPRESSVPN_420%2Ctelegram-EXPRESSVPN_420%2Ctelegram-EXPRESSVPN_420%2Ctelegram-EXPRESSVPN_420%3Fed%3D8080&host=Digikala.iranian.net#%40EXPRESSVPN_420%20%F0%9F%87%A9%F0%9F%87%AAAll2\n\nvless://8778b6fb-3547-47a9-8494-5b93e9fd5972@EXPRESSVPN_420.fast.hosting-ip.com:80?security=none&type=ws&headerType=&path=%2Ftelegram-EXPRESSVPN_420%2Ctelegram-EXPRESSVPN_420%2Ctelegram-EXPRESSVPN_420%2Ctelegram-EXPRESSVPN_420%2Ctelegram-EXPRESSVPN_420%2Ctelegram-EXPRESSVPN_420%2Ctelegram-EXPRESSVPN_420%2Ctelegram-EXPRESSVPN_420%3Fed%3D8080&host=Digikala.iranian.net#%40EXPRESSVPN_420%20-All\n\nvless://8778b6fb-3547-47a9-8494-5b93e9fd5972@speedtest.net:80?security=none&type=ws&headerType=&path=%2Ftelegram-EXPRESSVPN_420%2Ctelegram-EXPRESSVPN_420%2Ctelegram-EXPRESSVPN_420%2Ctelegram-EXPRESSVPN_420%2Ctelegram-EXPRESSVPN_420%2Ctelegram-EXPRESSVPN_420%2Ctelegram-EXPRESSVPN_420%2Ctelegram-EXPRESSVPN_420%3Fed%3D8080&host=Digikala.iranian.net#%40EXPRESSVPN_420%20%F0%9F%87%A9%F0%9F%87%AAall\n\nتیک Mux و Fragment خاموش باشه!",
		expectedName: time.Now().Format("Jan 02 2006 15:04"),
		expectedRawConfigs: []string{
			"vless://8778b6fb-3547-47a9-8494-5b93e9fd5972@EXPRESSVPN_420.fastly80-3.hosting-ip.com:80?security=none&type=ws&headerType=&path=%2Ftelegram-EXPRESSVPN_420%2Ctelegram-EXPRESSVPN_420%2Ctelegram-EXPRESSVPN_420%2Ctelegram-EXPRESSVPN_420%2Ctelegram-EXPRESSVPN_420%2Ctelegram-EXPRESSVPN_420%2Ctelegram-EXPRESSVPN_420%2Ctelegram-EXPRESSVPN_420%3Fed%3D8080&host=Digikala.iranian.net#%40EXPRESSVPN_420%20%F0%9F%87%A9%F0%9F%87%AAAll2",
			"vless://8778b6fb-3547-47a9-8494-5b93e9fd5972@EXPRESSVPN_420.fast.hosting-ip.com:80?security=none&type=ws&headerType=&path=%2Ftelegram-EXPRESSVPN_420%2Ctelegram-EXPRESSVPN_420%2Ctelegram-EXPRESSVPN_420%2Ctelegram-EXPRESSVPN_420%2Ctelegram-EXPRESSVPN_420%2Ctelegram-EXPRESSVPN_420%2Ctelegram-EXPRESSVPN_420%2Ctelegram-EXPRESSVPN_420%3Fed%3D8080&host=Digikala.iranian.net#%40EXPRESSVPN_420%20-All",
			"vless://8778b6fb-3547-47a9-8494-5b93e9fd5972@speedtest.net:80?security=none&type=ws&headerType=&path=%2Ftelegram-EXPRESSVPN_420%2Ctelegram-EXPRESSVPN_420%2Ctelegram-EXPRESSVPN_420%2Ctelegram-EXPRESSVPN_420%2Ctelegram-EXPRESSVPN_420%2Ctelegram-EXPRESSVPN_420%2Ctelegram-EXPRESSVPN_420%2Ctelegram-EXPRESSVPN_420%3Fed%3D8080&host=Digikala.iranian.net#%40EXPRESSVPN_420%20%F0%9F%87%A9%F0%9F%87%AAall",
		},
	},
	{
		name:               "should return empty configs when input is only spaces and new lines",
		text:               "   \n  \n    ",
		expectedName:       time.Now().Format("Jan 02 2006 15:04"),
		expectedRawConfigs: []string{},
	},
}

func TestExtract(t *testing.T) {

	t.Parallel()
	for _, testCase := range configsExtractorTests {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			// Arrange
			service := services.ConfigsExtractor{}
			expectedNameTime, _ := time.Parse("Jan 02 2006 15:04", testCase.expectedName)
			// Act
			name, rawConfigs := service.Extract(testCase.text)
			nameTime, nameErr := time.Parse("Jan 02 2006 15:04", name)
			// Assert
			assert.NoError(t, nameErr)
			assert.Equal(t, testCase.expectedRawConfigs, rawConfigs)
			assert.WithinDuration(t, expectedNameTime, nameTime, 5*time.Minute)
		})
	}
}
