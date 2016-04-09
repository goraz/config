package config

import (
	// "os"
	"reflect"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestConfig(t *testing.T) {
	Convey("Config basic functionality", t, func() {
		mainConfig := New()
		slice := make([]interface{}, 3)
		slice[0] = "1"
		slice[1] = 2
		slice[2] = true

		slice2 := make([]interface{}, 4)
		slice2[0] = "foo"
		slice2[1] = 2
		slice2[2] = true
		slice2[3] = map[string]interface{}{
			"val": 1,
		}

		configData := map[string]interface{}{
			"foo":         1,
			"slice":       slice,
			"slice2":      slice2,
			"stringSlice": []string{"foo"},
			"intSlice":    []int64{1},
			"floatSlice":  []float64{1},
			"mapString": map[string]string{
				"key1": "val1",
				"key2": "val2",
			},
			"floatData":             1.1,
			"bar":                   "baz",
			"bool":                  "true",
			"duration":              time.Duration(123),
			"durationString":        "1h2m3s",
			"durationInt":           int(123),
			"durationIntString":     "123",
			"durationInvalidString": "baz",
			"db": map[string]interface{}{
				"username": "user",
			},
			"db2": map[interface{}]interface{}{
				"1": "test",
			},
			"deep": map[string]interface{}{
				"notvalid": map[int]interface{}{
					1: "bar",
				},
			},
			"deep2": map[interface{}]interface{}{
				"notvalid": map[int]interface{}{
					1: "bar",
				},
			},
		}
		config2Data := map[string]interface{}{
			"db": map[string]interface{}{
				"password": "pass",
			},
		}
		mainConfig.AddMap(configData)
		mainConfig.AddMap(config2Data)
		Convey("Sub Config", func() {
			dbConfig := mainConfig.GetConfig("db")
			notConfig := mainConfig.GetConfig("bar")
			notFoundConfig := mainConfig.GetConfig("notfound")
			So(dbConfig, ShouldNotBeNil)
			So(notConfig, ShouldBeNil)
			So(notFoundConfig, ShouldBeNil)
			So(dbConfig.GetString("username"), ShouldEqual, "user")

		})

		Convey("Get variable", func() {
			So(mainConfig.GetString("bar"), ShouldEqual, "baz")
			So(mainConfig.GetInt("foo"), ShouldEqual, 1)
			So(mainConfig.GetFloat32("floatData"), ShouldEqual, 1.1)
			So(mainConfig.GetFloat64("floatData"), ShouldEqual, 1.1)
			So(mainConfig.GetBool("bool"), ShouldEqual, true)
			duration, _ := time.ParseDuration("1h2m3s")
			So(mainConfig.GetDuration("duration"), ShouldEqual, 123)
			So(mainConfig.GetDuration("durationString"), ShouldEqual, duration)
			So(mainConfig.GetDuration("durationIntString"), ShouldEqual, 123)
			So(mainConfig.GetDuration("durationInt"), ShouldEqual, 123)
			So(reflect.DeepEqual(mainConfig.GetStringSlice("slice"), []string{"1", "2", "true"}), ShouldBeTrue)
			So(reflect.DeepEqual(mainConfig.GetIntSlice("slice"), []int64{1, 2, 1}), ShouldBeTrue)
			So(reflect.DeepEqual(mainConfig.GetFloatSlice("slice"), []float64{1, 2, 1}), ShouldBeTrue)
			So(reflect.DeepEqual(mainConfig.GetStringSlice("stringSlice"), []string{"foo"}), ShouldBeTrue)
			So(reflect.DeepEqual(mainConfig.GetIntSlice("intSlice"), []int64{1}), ShouldBeTrue)
			So(reflect.DeepEqual(mainConfig.GetFloatSlice("floatSlice"), []float64{1}), ShouldBeTrue)
			So(reflect.DeepEqual(mainConfig.GetStringSlice("intSlice"), []string{"1"}), ShouldBeTrue)
		})

		Convey("Get not found variable", func() {
			So(mainConfig.GetString(""), ShouldEqual, "")
			So(mainConfig.GetString("notfound"), ShouldEqual, "")
			So(mainConfig.GetInt("notfound"), ShouldEqual, 0)
			So(mainConfig.GetInt64("notfound"), ShouldEqual, 0)
			So(mainConfig.GetFloat32("notfound"), ShouldEqual, 0)
			So(mainConfig.GetFloat64("notfound"), ShouldEqual, 0)
			So(mainConfig.GetBool("notfound"), ShouldEqual, false)
			So(mainConfig.GetStringSlice("notfound"), ShouldBeNil)
			So(mainConfig.GetIntSlice("notfound"), ShouldBeNil)
			So(mainConfig.GetFloatSlice("notfound"), ShouldBeNil)
			So(mainConfig.GetDuration("notfound"), ShouldEqual, 0)
			So(mainConfig.GetDuration("db2.bar"), ShouldEqual, 0)
		})

		Convey("Get not valid variable", func() {
			So(mainConfig.GetString("slice"), ShouldEqual, "")
			So(mainConfig.GetInt("slice"), ShouldEqual, 0)
			So(mainConfig.GetInt64("slice"), ShouldEqual, 0)
			So(mainConfig.GetFloat32("slice"), ShouldEqual, 0)
			So(mainConfig.GetFloat64("slice"), ShouldEqual, 0)
			So(mainConfig.GetBool("slice"), ShouldEqual, false)
			So(mainConfig.GetDuration("slice"), ShouldEqual, 0)
			So(mainConfig.GetDuration("durationInvalidString"), ShouldEqual, 0)
			So(mainConfig.GetStringSlice("slice2"), ShouldBeNil)

			So(mainConfig.GetIntSlice("stringSlice"), ShouldBeNil)
			So(mainConfig.GetIntSlice("slice2"), ShouldBeNil)
			So(mainConfig.GetFloatSlice("stringSlice"), ShouldBeNil)
			So(mainConfig.GetFloatSlice("slice2"), ShouldBeNil)
		})

		Convey("Test Delimeter", func() {
			mainConfig.SetDelimiter("")
			So(mainConfig.GetDelimiter(), ShouldEqual, DefaultDelimiter)
			mainConfig.SetDelimiter(",")
			So(mainConfig.GetString("db,username"), ShouldEqual, "user")
			So(mainConfig.GetString("db2,1"), ShouldEqual, "test")
		})

		Convey("Test Not valid key type", func() {
			So(mainConfig.GetString("deep.notvalid.1"), ShouldEqual, "")
			So(mainConfig.GetString("deep2.notvalid.1"), ShouldEqual, "")
		})

	})
}
