// Package tomlloader is used to handle toml file in Config file/folder layer.
// for using this package, just import it
//
// 		import (
//			_ "github.com/fzerorubigd/config/tomlloader"
//		)
//
// There is no need to do anything else, if you load a file with toml
// extension, the toml loader is doing his job.
package tomlloader

import (
	"io"
	"io/ioutil"

	"github.com/BurntSushi/toml"

	"github.com/goraz/config"
)

type tomlLoader struct {
}

func (tl tomlLoader) SupportedEXT() []string {
	return []string{".toml"}
}

func (tl tomlLoader) Convert(r io.Reader) (map[string]interface{}, error) {
	data, _ := ioutil.ReadAll(r)
	ret := make(map[string]interface{})
	err := toml.Unmarshal(data, &ret)
	if err != nil {
		return nil, err
	}

	return ret, nil
}

func init() {
	config.RegisterLoader(&tomlLoader{})
}
