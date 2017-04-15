package confWriters

import "gopkg.in/yaml.v2"

func ConfToYaml(conf map[string]string) (string, error) {
	structConf, err := mapToStructMap(conf, "")
	if err != nil {
		return "", nil
	}
	out, err := yaml.Marshal(structConf)

	return string(out), err
}
