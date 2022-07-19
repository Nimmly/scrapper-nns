package util

import (
	"io/ioutil"
	"log"

	conf "github.com/Nimmly/scrapper-nns/models/config"
	"gopkg.in/yaml.v2"
)

func ReadPostgresConfig(config_file string) (*conf.PostgresConfig, error) {

	config := &conf.PostgresConfig{}
	source, err := ioutil.ReadFile(config_file)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal([]byte(source), &config)
	if err != nil {
		log.Println("Reader can't unmarshal Postgres Config from YAML file!")
		return nil, err
	}

	//log.Printf("Postgres Config is loaded: %v\n",config)
	return config, nil
}
