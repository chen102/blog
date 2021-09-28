package tool

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Mysql struct {
		Name            string `yml:"name"`
		Dns             string `yml:"dns"`
		Log             bool   `yml:"log"`
		Logfile         string `yml:"logfile"`
		MaxIdleConns    int    `yml:"maxidleConns"`
		MaxOpenConns    int    `yml:"maxopenConns"`
		ConnMaxLifetime int    `yml:"connmaxlifetime"`
	} `yml:mysql`
	Redis struct {
		MasterIP   string `yml:"masterip"`
		MasterPort string `yml:"masterport"`
		MasterAuth string `yml:"masterauth"`
		MasterDB   int    `yml:"masterdb"`
		SlaveOfIp  string `yml:"slaveofip"`
		SlaveIP    string `yml:"salveip"`
		SlavePort  string `yml:"slaveport"`
		SlaveAuth  string `yml:"slaveauth"`
		SlaveDB    int    `yml:"salvedb"`
	} `yml:redis`
}

func NewConfig() *Config {
	return &Config{}
}
func (config *Config) ReadConfig() {
	file, err := ioutil.ReadFile("config/config.yml")
	if err != nil {
		panic("读取文件config.yml发生错误")
	}
	if yaml.Unmarshal(file, config) != nil {
		panic("解析文件config.yml发生错误")
	}
}
func (config *Config) TestReadConfig() {
	file, err := ioutil.ReadFile("../config/config.yml")
	if err != nil {
		panic("读取文件config.yml发生错误")
	}
	if yaml.Unmarshal(file, config) != nil {
		panic("解析文件config.yml发生错误")
	}
}
