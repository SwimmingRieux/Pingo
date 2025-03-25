package services_test

import (
	"github.com/stretchr/testify/assert"
	"pingo/internal/app/services"
	"strings"
	"testing"
)

type jsonFormatterTest struct {
	name           string
	rawConfig      string
	expectedConfig string
}

var jsonFormatterTests = []jsonFormatterTest{
	{
		name:           "should return valid json string for vless config",
		rawConfig:      "vless://8778b6fb-3547-47a9-8494-5b93e9fd5972@EXPRESSVPN_420.fastly80-3.hosting-ip.com:80?security=none&type=ws&headerType=&path=%2Ftelegram-EXPRESSVPN_420%2Ctelegram-EXPRESSVPN_420%2Ctelegram-EXPRESSVPN_420%2Ctelegram-EXPRESSVPN_420%2Ctelegram-EXPRESSVPN_420%2Ctelegram-EXPRESSVPN_420%2Ctelegram-EXPRESSVPN_420%2Ctelegram-EXPRESSVPN_420%3Fed%3D8080&host=Digikala.iranian.net#%40EXPRESSVPN_420%20%F0%9F%87%A9%F0%9F%87%AAAll2",
		expectedConfig: "{\n    \"dns\": {\n        \"disableFallback\": true,\n        \"servers\": [\n            {\n                \"address\": \"https://8.8.8.8/dns-query\",\n                \"domains\": [],\n                \"queryStrategy\": \"\"\n            },\n            {\n                \"address\": \"localhost\",\n                \"domains\": [\n                    \"full:expressvpn_420.fastly80-3.hosting-ip.com\"\n                ],\n                \"queryStrategy\": \"\"\n            }\n        ],\n        \"tag\": \"dns\"\n    },\n    \"inbounds\": [\n        {\n            \"listen\": \"127.0.0.1\",\n            \"port\": 12346,\n            \"protocol\": \"socks\",\n            \"settings\": {\n                \"udp\": true\n            },\n            \"sniffing\": {\n                \"destOverride\": [\n                    \"http\",\n                    \"tls\",\n                    \"quic\"\n                ],\n                \"enabled\": true,\n                \"metadataOnly\": false,\n                \"routeOnly\": true\n            },\n            \"tag\": \"socks-in\"\n        },\n        {\n            \"listen\": \"127.0.0.1\",\n            \"port\": 12346,\n            \"protocol\": \"http\",\n            \"sniffing\": {\n                \"destOverride\": [\n                    \"http\",\n                    \"tls\",\n                    \"quic\"\n                ],\n                \"enabled\": true,\n                \"metadataOnly\": false,\n                \"routeOnly\": true\n            },\n            \"tag\": \"http-in\"\n        }\n    ],\n    \"log\": {\n        \"loglevel\": \"warning\"\n    },\n    \"outbounds\": [\n        {\n            \"domainStrategy\": \"AsIs\",\n            \"flow\": null,\n            \"protocol\": \"vless\",\n            \"settings\": {\n                \"vnext\": [\n                    {\n                        \"address\": \"expressvpn_420.fastly80-3.hosting-ip.com\",\n                        \"port\": 80,\n                        \"users\": [\n                            {\n                                \"encryption\": \"none\",\n                                \"flow\": \"\",\n                                \"id\": \"8778b6fb-3547-47a9-8494-5b93e9fd5972\"\n                            }\n                        ]\n                    }\n                ]\n            },\n            \"streamSettings\": {\n                \"network\": \"ws\",\n                \"wsSettings\": {\n                    \"headers\": {\n                        \"Host\": \"Digikala.iranian.net\"\n                    },\n                    \"path\": \"/telegram-EXPRESSVPN_420,telegram-EXPRESSVPN_420,telegram-EXPRESSVPN_420,telegram-EXPRESSVPN_420,telegram-EXPRESSVPN_420,telegram-EXPRESSVPN_420,telegram-EXPRESSVPN_420,telegram-EXPRESSVPN_420?ed=8080\"\n                }\n            },\n            \"tag\": \"proxy\"\n        },\n        {\n            \"domainStrategy\": \"\",\n            \"protocol\": \"freedom\",\n            \"tag\": \"direct\"\n        },\n        {\n            \"domainStrategy\": \"\",\n            \"protocol\": \"freedom\",\n            \"tag\": \"bypass\"\n        },\n        {\n            \"protocol\": \"blackhole\",\n            \"tag\": \"block\"\n        },\n        {\n            \"protocol\": \"dns\",\n            \"proxySettings\": {\n                \"tag\": \"proxy\",\n                \"transportLayer\": true\n            },\n            \"settings\": {\n                \"address\": \"8.8.8.8\",\n                \"network\": \"tcp\",\n                \"port\": 53,\n                \"userLevel\": 1\n            },\n            \"tag\": \"dns-out\"\n        }\n    ],\n    \"policy\": {\n        \"levels\": {\n            \"1\": {\n                \"connIdle\": 30\n            }\n        },\n        \"system\": {\n            \"statsOutboundDownlink\": true,\n            \"statsOutboundUplink\": true\n        }\n    },\n    \"routing\": {\n        \"domainStrategy\": \"AsIs\",\n        \"rules\": [\n            {\n                \"inboundTag\": [\n                    \"socks-in\",\n                    \"http-in\"\n                ],\n                \"outboundTag\": \"dns-out\",\n                \"port\": \"53\",\n                \"type\": \"field\"\n            },\n            {\n                \"outboundTag\": \"proxy\",\n                \"port\": \"0-65535\",\n                \"type\": \"field\"\n            }\n        ]\n    },\n    \"stats\": {}\n}",
	},
}

func TestFormat(t *testing.T) {
	t.Parallel()
	for _, test := range jsonFormatterTests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			// Arrange
			factory := services.NewFormatterFactory(ConfigForTest)
			configType := strings.Split(test.rawConfig, "://")[0]
			formatter, err := factory.Fetch(configType)
			assert.NoError(t, err)

			// Act
			result, err := formatter.Format(test.rawConfig)

			// Assert
			assert.NoError(t, err)
			assert.Equal(t, test.expectedConfig, result)
		})
	}
}
