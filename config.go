package config

import (
	"github.com/goraz/cast"
	"strings"
	"sync"
	"time"
)

// DefaultDelimiter is the default delimiter for the config scope
const DefaultDelimiter = "."

// Layer is an interface to handle the load phase.
type Layer interface {
	// Load a layer into the Config. the call is only done in the
	// registration
	Load() (map[string]interface{}, error)
}

// Config is a layer base configuration system
type Config struct {
	delimiter string
	data      map[string]interface{}
	lock      sync.RWMutex
}

// AddLayer add a new layer to the end of config layers. last layer is loaded after all other
// layer
func (o *Config) AddLayer(l Layer) error {
	o.lock.Lock()
	defer o.lock.Unlock()

	data, err := l.Load()
	if err != nil {
		return err
	}

	lowerStringMap(o.data, data)

	return nil
}

// AddMap add map to config; usually use for default config
func (o *Config) AddMap(data map[string]interface{}) {
	lowerStringMap(o.data, data)
}

//GetConfig Get New Config for a key
func (o *Config) GetConfig(key string) *Config {
	val, found := o.Get(key)
	if found == false {
		return nil
	}
	switch data := val.(type) {
	case map[string]interface{}:
		return &Config{
			delimiter: o.delimiter,
			data:      data,
		}
	default:
		return nil
	}

}

// GetDelimiter return the delimiter for nested key
func (o *Config) GetDelimiter() string {
	if o.delimiter == "" {
		o.delimiter = DefaultDelimiter
	}

	return o.delimiter
}

// SetDelimiter set the current delimiter
func (o *Config) SetDelimiter(d string) {
	o.delimiter = d
}

// Get try to get the key from config layers
func (o *Config) Get(key string) (interface{}, bool) {
	o.lock.RLock()
	defer o.lock.RUnlock()
	key = strings.Trim(key, " ")
	if len(key) == 0 {
		return nil, false
	}
	path := strings.Split(strings.ToLower(key), o.GetDelimiter())

	return searchStringMap(path, o.data)
}

// The following two function are identical. but converting between map[string] and
// map[interface{}] is not easy, and there is no _Generic_ way to do it, so I decide to create
// two almost identical function instead of writing a converter each time.
//
// Some of the loaders like yaml, load inner keys in map[interface{}]interface{}
// some other like json do it in map[string]interface{} so we should support both
func searchStringMap(path []string, m map[string]interface{}) (interface{}, bool) {
	v, ok := m[path[0]]

	if !ok {
		return nil, false
	}

	if len(path) == 1 {
		return v, true
	}

	switch m := v.(type) {
	case map[string]interface{}:
		return searchStringMap(path[1:], m)
	case map[interface{}]interface{}:
		return searchInterfaceMap(path[1:], m)
	}
	return nil, false
}

func searchInterfaceMap(path []string, m map[interface{}]interface{}) (interface{}, bool) {
	v, ok := m[path[0]]
	if !ok {
		return nil, false
	}

	if len(path) == 1 {
		return v, true
	}

	switch m := v.(type) {
	case map[string]interface{}:
		return searchStringMap(path[1:], m)
	case map[interface{}]interface{}:
		return searchInterfaceMap(path[1:], m)
	}
	return nil, false
}

func lowerStringMap(dst map[string]interface{}, src map[string]interface{}) {
	for k := range src {
		lowerK := strings.ToLower(k)
		switch nm := src[k].(type) {
		case map[string]interface{}:
			if dst[lowerK] == nil {
				dst[lowerK] = make((map[string]interface{}))
			}
			lowerStringMap(dst[lowerK].(map[string]interface{}), nm)
		case map[interface{}]interface{}:
			if dst[lowerK] == nil {
				dst[lowerK] = make((map[interface{}]interface{}))
			}
			lowerInterfaceMap(dst[lowerK].(map[interface{}]interface{}), nm)
		default:
			dst[lowerK] = src[k]
		}
	}

}

func lowerInterfaceMap(dst map[interface{}]interface{}, src map[interface{}]interface{}) {
	for k := range src {
		switch k.(type) {
		case string:
			lowerK := strings.ToLower(k.(string))
			switch nm := src[k].(type) {
			case map[string]interface{}:
				if dst[lowerK] == nil {
					dst[lowerK] = make((map[string]interface{}))
				}
				lowerStringMap(dst[lowerK].(map[string]interface{}), nm)
			case map[interface{}]interface{}:
				if dst[lowerK] == nil {
					dst[lowerK] = make((map[interface{}]interface{}))
				}
				lowerInterfaceMap(dst[lowerK].(map[interface{}]interface{}), nm)
			default:
				dst[lowerK] = src[k]
			}
			//todo
			// default:
			// dst[k] = src[k]
		}
	}

}

// GetIntDefault return an int value from Config, if the value is not exists or its not an
// integer , default is returned
func (o *Config) GetIntDefault(key string, def int) int {
	return int(o.GetInt64Default(key, int64(def)))
}

// GetInt return an int value, if the value is not there, then it return zero value
func (o *Config) GetInt(key string) int {
	return o.GetIntDefault(key, 0)
}

// GetInt64Default return an int64 value from Config, if the value is not exists or if the value is not
// int64 then return the default
func (o *Config) GetInt64Default(key string, def int64) int64 {
	v, ok := o.Get(key)
	if !ok {
		return def
	}
	output, err := cast.Int(v)
	if err != nil {
		return def
	}
	return output
}

// GetInt64 return the int64 value from config, if its not there, return zero
func (o *Config) GetInt64(key string) int64 {
	return o.GetInt64Default(key, 0)
}

// GetFloat32Default return an float32 value from Config, if the value is not exists or its not a
// float32, default is returned
func (o *Config) GetFloat32Default(key string, def float32) float32 {
	return float32(o.GetFloat64Default(key, float64(def)))
}

// GetFloat32 return an float32 value, if the value is not there, then it returns zero value
func (o *Config) GetFloat32(key string) float32 {
	return o.GetFloat32Default(key, 0)
}

// GetFloat64Default return an float64 value from Config, if the value is not exists or if the value is not
// float64 then return the default
func (o *Config) GetFloat64Default(key string, def float64) float64 {
	v, ok := o.Get(key)
	if !ok {
		return def
	}

	output, err := cast.Float(v)
	if err != nil {
		return def
	}
	return output

}

// GetFloat64 return the float64 value from config, if its not there, return zero
func (o *Config) GetFloat64(key string) float64 {
	return o.GetFloat64Default(key, 0)
}

// GetStringDefault get a string from Config. if the value is not exists or if tha value is not
// string, return the default
func (o *Config) GetStringDefault(key string, def string) string {
	v, ok := o.Get(key)
	if !ok {
		return def
	}

	output, err := cast.String(v)
	if err != nil {
		return def
	}
	return output
}

// GetString is for getting an string from conig. if the key is not
func (o *Config) GetString(key string) string {
	return o.GetStringDefault(key, "")
}

// GetBoolDefault return bool value from Config. if the value is not exists or if tha value is not
// boolean, return the default
func (o *Config) GetBoolDefault(key string, def bool) bool {
	v, ok := o.Get(key)
	if !ok {
		return def
	}

	output, err := cast.Bool(v)
	if err != nil {
		return def
	}
	return output
}

// GetBool is used to get a boolean value fro config, with false as default
func (o *Config) GetBool(key string) bool {
	return o.GetBoolDefault(key, false)
}

// GetDurationDefault is a function to get duration from config. it support both
// string duration (like 1h3m2s) and integer duration
func (o *Config) GetDurationDefault(key string, def time.Duration) time.Duration {
	v, ok := o.Get(key)
	if !ok {
		return def
	}

	switch nv := v.(type) {
	case string:
		d, err := time.ParseDuration(nv)
		if err != nil {
			output, err := cast.Int(v)
			if err != nil {
				return def
			}
			return time.Duration(output)
		}
		return d
	case time.Duration:
		return nv
	default:
		output, err := cast.Int(v)
		if err == nil {
			return time.Duration(output)
		}
		return def
	}
}

// GetDuration is for getting duration from config, it cast both int and string
// to duration
func (o *Config) GetDuration(key string) time.Duration {
	return o.GetDurationDefault(key, 0)
}

// GetStringSlice try to get a string slice from the config
func (o *Config) GetStringSlice(key string) []string {
	var ok bool
	v, ok := o.Get(key)
	if !ok {
		return nil
	}
	output, err := cast.StringSlice(v)
	if err != nil {
		return nil
	}
	return output
}

// GetIntSlice try to get a int64 slice from the config
func (o *Config) GetIntSlice(key string) []int64 {
	var ok bool
	v, ok := o.Get(key)
	if !ok {
		return nil
	}
	output, err := cast.IntSlice(v)
	if err != nil {
		return nil
	}
	return output
}

// GetFloatSlice try to get a float64 slice from the config
func (o *Config) GetFloatSlice(key string) []float64 {
	var ok bool
	v, ok := o.Get(key)
	if !ok {
		return nil
	}
	output, err := cast.FloatSlice(v)
	if err != nil {
		return nil
	}
	return output
}

// New return a new Config
func New() *Config {
	return &Config{
		delimiter: DefaultDelimiter,
		data:      make(map[string]interface{}),
	}
}
