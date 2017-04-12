package confWriters

import (
	"strings"
	"github.com/Sirupsen/logrus"
	"reflect"
)


func mapToStructMap(conf map[string]string, prefix string) (map[string]interface{}, error) {
	res := make(map[string]interface{})
	var err error
	for k, v := range conf {
		if len(prefix) == 0 || strings.HasPrefix(k, prefix) {
			key := strings.Replace(k, prefix, "", 1)
			logrus.Debug(key)
			if !strings.Contains(key, ".") {
				if _, alreadyPresent := res[key]; alreadyPresent {
					return nil, InvalidConfigurationStructure.Details("Duplicated key:" + key)
				}
				res[key] = v
			} else {
				keys := strings.Split(key, ".")
				if _, alreadyPresent := res[keys[0]]; !alreadyPresent {
					if len(prefix) == 0 {
						res[keys[0]], err = mapToStructMap(conf, keys[0] + ".")
					} else {
						res[keys[0]], err = mapToStructMap(conf, prefix + keys[0] + ".")
					}
					if err != nil {
						return nil, err
					}
				}else{
					if reflect.TypeOf(res[keys[0]]).Kind()!=reflect.Map{
						return nil, InvalidConfigurationStructure.Details("Duplicated key:" + keys[0])
					}
				}
			}
		}
	}
	return res, nil
}
