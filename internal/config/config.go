package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/viper"
)

const (
	EnvPrefix      = "WEFLOW"
	DefaultHost    = "http://127.0.0.1:5031"
	DefaultTimeout = 30 * time.Second
)

type OutputConfig struct {
	Color bool `mapstructure:"color" json:"color"`
}

type Config struct {
	Host    string        `mapstructure:"host"`
	Token   string        `mapstructure:"token"`
	Timeout time.Duration `mapstructure:"timeout"`
	Output  OutputConfig  `mapstructure:"output"`

	SourcePath string `mapstructure:"-"`
}

// MarshalJSON 让 `config show --json` 输出人类可读：字段小写、timeout 用字符串。
func (c Config) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Host       string       `json:"host"`
		Token      string       `json:"token"`
		Timeout    string       `json:"timeout"`
		Output     OutputConfig `json:"output"`
		SourcePath string       `json:"source_path,omitempty"`
	}{
		Host:       c.Host,
		Token:      c.Token,
		Timeout:    c.Timeout.String(),
		Output:     c.Output,
		SourcePath: c.SourcePath,
	})
}

const Template = `# WeFlow CLI 配置
host: http://127.0.0.1:5031
token: ""             # 在 WeFlow 应用设置 → API 服务里复制 Access Token
timeout: 30s
output:
  color: true
`

func Load(explicitPath string) (*Config, error) {
	v := viper.New()
	v.SetDefault("host", DefaultHost)
	v.SetDefault("token", "")
	v.SetDefault("timeout", DefaultTimeout)
	v.SetDefault("output.color", true)

	v.SetEnvPrefix(EnvPrefix)
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	path := explicitPath
	if path == "" {
		path = findConfigFile()
	}

	if path != "" {
		v.SetConfigFile(path)
		if err := v.ReadInConfig(); err != nil {
			var notFound viper.ConfigFileNotFoundError
			if !errors.As(err, &notFound) {
				return nil, fmt.Errorf("读取配置 %s: %w", path, err)
			}
		}
	}

	var c Config
	if err := v.Unmarshal(&c); err != nil {
		return nil, fmt.Errorf("解析配置: %w", err)
	}
	c.SourcePath = v.ConfigFileUsed()
	return &c, nil
}

func findConfigFile() string {
	for _, p := range CandidatePaths() {
		if st, err := os.Stat(p); err == nil && !st.IsDir() {
			return p
		}
	}
	return ""
}

func CandidatePaths() []string {
	paths := []string{"weflow.yaml", "weflow.yml"}
	if xdg := xdgConfigDir(); xdg != "" {
		paths = append(paths,
			filepath.Join(xdg, "weflow", "config.yaml"),
			filepath.Join(xdg, "weflow", "config.yml"),
		)
	}
	if home, err := os.UserHomeDir(); err == nil {
		paths = append(paths,
			filepath.Join(home, ".weflow.yaml"),
			filepath.Join(home, ".weflow.yml"),
		)
	}
	return paths
}

func DefaultPath() string {
	if xdg := xdgConfigDir(); xdg != "" {
		return filepath.Join(xdg, "weflow", "config.yaml")
	}
	if home, err := os.UserHomeDir(); err == nil {
		return filepath.Join(home, ".config", "weflow", "config.yaml")
	}
	return "weflow.yaml"
}

func WriteTemplate(path string) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("创建目录: %w", err)
	}
	return os.WriteFile(path, []byte(Template), 0o600)
}

func xdgConfigDir() string {
	if v := os.Getenv("XDG_CONFIG_HOME"); v != "" {
		return v
	}
	if home, err := os.UserHomeDir(); err == nil {
		return filepath.Join(home, ".config")
	}
	return ""
}
