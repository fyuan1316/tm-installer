package cmd

import (
	"encoding/json"
	"fmt"
	v1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"testing"
)


func Test_loadMap(t *testing.T) {
	deploy1 := v1.Deployment{}
	deploy1.Namespace = "default"
	deploy1.Name = "name1"
	deploy1.Kind = "Deployment"


	bytes, err := json.Marshal(&deploy1)
	if err != nil {
		panic(err)
	}

	typemeta := &metav1.TypeMeta{}
	err = json.Unmarshal(bytes, typemeta)
	if err != nil {
		panic(err)
	}
	switch typemeta.Kind {
	case "Deployment":
		fmt.Println("got Deployment")
	default:
		fmt.Println("unknow kind")
	}

}
