//go:build exec

package uci

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Sergii-Kirichok/DTekSpeachParser/pkg/helpers/testtools/exec"
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const (
	testdata = "testdata"
	trigger  = "trigger"
)

var (
	deviceConfig = filepath.Join(testdata, "device-config.json")

	deviceConfigStub = map[string]string{
		"network.wan.dns": ToList([]string{
			"8.8.8.8",
			"8.8.4.4",
		}),
		"network.wlan.ipaddr":  "109.207.201.210",
		"network.wlan.netmask": "'255.255.255.0'",
		"system.ntp.server": ToList([]string{
			"0.lede.pool.ntp.org",
			"1.lede.pool.ntp.org",
			"2.lede.pool.ntp.org",
		}),
		"wireless.radio0.country": "UA",
		"trigger.target":          "nothing more than a trigger target",
	}
)

func TestGet(t *testing.T) {
	rollback := exec.SetCommandMock("TestUCIMock")
	clean := createDeviceConfig(t)

	defer rollback()
	defer clean()

	type ret struct {
		value  string
		hasErr bool
	}

	testCases := []struct {
		name   string
		input  string
		output ret
	}{
		{
			input:  "unknown.target",
			output: ret{hasErr: true},
		},
		{
			input:  "wireless.radio0.country",
			output: ret{value: deviceConfigStub["wireless.radio0.country"]},
		},
		{
			input:  "system.ntp.server",
			output: ret{value: deviceConfigStub["system.ntp.server"]},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			value, err := get(tc.input)
			assert.Equal(t, tc.output.value, value)
			if tc.output.hasErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestSet(t *testing.T) {
	rollback := exec.SetCommandMock("TestAllSuccessfulCalls")
	clean := createDeviceConfig(t)

	defer rollback()
	defer clean()

	assert.NoError(t, set("network.wan.dns", "'4.4.4.4' 8.8.8.8"))
}

func getConfigHelper(t *testing.T) map[string]string {
	t.Helper()

	file, err := os.OpenFile(deviceConfig, os.O_CREATE|os.O_RDONLY, 0644)
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	var config map[string]string
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		t.Fatal(err)
	}

	return config
}

func createDeviceConfig(t *testing.T) func() {
	t.Helper()

	if err := os.MkdirAll(testdata, 0744); err != nil {
		t.Fatal(err)
	}
	setConfigHelper(t, deviceConfigStub)

	// clean function
	return func() {
		t.Helper()

		if err := os.RemoveAll(testdata); err != nil {
			t.Fatal(err)
		}
	}
}

func setConfigHelper(t *testing.T, config map[string]string) {
	t.Helper()

	file, err := os.OpenFile(deviceConfig, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "\t")
	if err := encoder.Encode(config); err != nil {
		t.Fatal(err)
	}
}

func TestAllSuccessfulCalls(t *testing.T) {
	exec.TestProcessWrapper(
		func(cmd string, args []string) (success bool) {
			return true
		},
	)
}

/*
You can test this UCI mock using these following commands:
1) export GO_TEST_PROCESS=2                                   | we use it for enable debug mode
2) go test -test.run=TestUCIMock -- uci get system.ntp.server | it just one example of use
3) unset GO_TEST_PROCESS                                      | we use it for disable debug mode
*/
func TestUCIMock(t *testing.T) {
	if _, err := os.Stat(deviceConfig); errors.Is(err, os.ErrNotExist) {
		clean := createDeviceConfig(t)
		defer clean()
	}

	exec.TestProcessWrapper(
		func(cmd string, args []string) (success bool) {
			switch {
			case cmd == "uci" && len(args) == 2:
				switch args[0] {
				case "get":
					return uciGetMock(t, args[1])
				case "set":
					return uciSetMock(t, args[1])
				case "add_list":
					return uciAddListMock(t, args[1])
				case "delete":
					return uciDeleteMock(t, args[1])
				case "commit":
					return uciCommitMock(t, args[1])
				}
			case cmd == "reload_config" && len(args) == 0:
				return uciReloadConfigMock()
			}

			return false
		},
	)
}

func uciGetMock(t *testing.T, aim string) bool {
	config := getConfigHelper(t)

	if value, ok := config[aim]; ok || strings.Contains(aim, trigger) {
		fmt.Printf("%s\n", value)
		return true
	}

	return false
}

func uciSetMock(t *testing.T, binding string) bool {
	config := getConfigHelper(t)

	if parts := strings.Split(binding, "="); len(parts) == 2 {
		aim, value := parts[0], parts[1]
		if _, ok := config[aim]; ok && !strings.Contains(value, trigger) {
			config[aim] = value
			setConfigHelper(t, config)
			return true
		}
	}

	return false
}

func uciAddListMock(t *testing.T, binding string) bool {
	config := getConfigHelper(t)

	if parts := strings.Split(binding, "="); len(parts) == 2 {
		aim, value := parts[0], parts[1]
		if !strings.Contains(value, trigger) {
			var prefix string
			if _, ok := config[aim]; ok {
				prefix = " "
			}

			config[aim] += prefix + value
			setConfigHelper(t, config)
			return true
		}
	}

	return false
}

func uciDeleteMock(t *testing.T, aim string) bool {
	config := getConfigHelper(t)

	if _, ok := config[aim]; ok {
		delete(config, aim)
		setConfigHelper(t, config)
		return true
	}

	return false
}

func uciCommitMock(t *testing.T, aim string) bool {
	config := getConfigHelper(t)

	for tg := range config {
		if aim == strings.Split(tg, ".")[0] && !strings.Contains(tg, trigger) {
			return true
		}
	}

	return false
}

func uciReloadConfigMock() bool { return false }
