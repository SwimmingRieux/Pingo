package services_test

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"net"
	"os"
	"pingo/internal/app/services"
	"testing"
)

type portSetterFailCase struct {
	name       string
	configType string
	tempPath   string
	jsonText   string
}

var actualConfigPath = "test_config.json"

var portSetterFailCases = []portSetterFailCase{
	{
		name:       "should return error when path is invalid",
		configType: "trojan",
		tempPath:   "somepaththatdoesntexist",
		jsonText:   "",
	},
	{
		name:       "should return error when config doesn't have port entry",
		configType: "vless",
		tempPath:   actualConfigPath,
		jsonText:   "{\"inbounds\": [{}]}",
	},
	{
		name:       "should return error when config doesn't have inbounds entry",
		configType: "vmess",
		tempPath:   actualConfigPath,
		jsonText:   "{}",
	},
}

func TestSetPort(t *testing.T) {
	t.Parallel()
	for _, testCase := range portSetterFailCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			os.WriteFile(actualConfigPath, []byte(testCase.jsonText), 0644)
			defer os.Remove(actualConfigPath)
			listener, err := net.Listen("tcp", ":0")
			assert.NoError(t, err)
			defer listener.Close()

			factory := services.NewPortSetterFactory(ConfigForTest)
			setter, _ := factory.Fetch(testCase.configType)

			err = setter.SetPort(listener, testCase.tempPath)
			assert.Error(t, err)
		})
	}
	t.Run("should successfully update port in config", func(t *testing.T) {
		initialConfig := "{\"inbounds\": [{\"port\": 12345}]}"
		os.WriteFile(actualConfigPath, []byte(initialConfig), 0644)
		defer os.Remove(actualConfigPath)

		listener, err := net.Listen("tcp", ":0")
		assert.NoError(t, err)
		defer listener.Close()

		newPort := listener.Addr().(*net.TCPAddr).Port
		factory := services.NewPortSetterFactory(ConfigForTest)
		setter, _ := factory.Fetch("ss")
		err = setter.SetPort(listener, actualConfigPath)
		assert.NoError(t, err)

		updatedData, err := os.ReadFile(actualConfigPath)
		assert.NoError(t, err)

		var updatedConfig map[string]interface{}
		err = json.Unmarshal(updatedData, &updatedConfig)
		assert.NoError(t, err)

		inbounds, ok := updatedConfig["inbounds"].([]interface{})
		assert.True(t, ok)
		assert.Len(t, inbounds, 1)

		entry, ok := inbounds[0].(map[string]interface{})
		assert.True(t, ok)
		assert.Equal(t, float64(newPort), entry["port"])

	})
}
