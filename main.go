package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"reflect"
	"sort"
	"strings"
)

func main() {
	b, err := ioutil.ReadFile("data.json")
	if err != nil {
		log.Fatal(err)
	}
	var data interface{}
	json.Unmarshal(b, &data)
	result := make([]string, 0)
	flatten(data, []string{}, &result)
	sort.Slice(result, func(i, j int) bool {
		a := result[i]
		b := result[j]
		return a < b
	})

	fmt.Print(strings.Join(result, ""))

}

/* Implements
function flatten(data, context) {
    if (Array.isArray(data)) {
        data.forEach(d => {
            if(d.name) {
                context = context.concat(d.name);
            }
            flatten(d, context);
            if(d.name) {
                context = context.slice(0, context.length - 1);
            }
        });
        return;
    }
    if (data && typeof data === 'object') {
        Object.keys(data).forEach(k => {
            context = context.concat(k);
            flatten(data[k], context);
            context = context.slice(0, context.length - 1);
        });
        return;
    }
    console.log(`${context.join(".")}=${data}`);
}
*/
func flatten(data interface{}, parents []string, result *[]string) {
	v := reflect.ValueOf(data)
	switch v.Kind() {
	case reflect.Slice:
		for i := 0; i < v.Len(); i++ {
			elem := reflect.ValueOf(v.Index(i).Interface())
			pushed := false
			switch elem.Kind() {
			case reflect.Map:
				m := elem.Interface().(map[string]interface{})
				if k, ok := m["name"]; ok {
					parents = append(parents, k.(string))
					pushed = true
				}
			}
			flatten(v.Index(i).Interface(), parents, result)
			if pushed {
				parents = parents[0 : len(parents)-1]
			}
		}
	case reflect.Map:
		for _, k := range v.MapKeys() {
			parents = append(parents, k.String())
			flatten(v.MapIndex(k).Interface(), parents, result)
			parents = parents[0 : len(parents)-1]
		}
	case reflect.Struct:
		for i := 0; i < v.NumField(); i++ {
			parents = append(parents, v.Type().Field(i).Name)
			flatten(v.Field(i).Interface(), parents, result)
			parents = parents[0 : len(parents)-1]
		}
	default:
		switch d := data.(type) {
		case string:
			d = strings.Join(strings.Split(data.(string), "\n"), "")
			*result = append(*result, fmt.Sprintf("%v=%s\n", strings.Join(parents, "."), d))
		default:
			*result = append(*result, fmt.Sprintf("%v=%s\n", strings.Join(parents, "."), data))
		}

	}
}
