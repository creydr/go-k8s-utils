package testing

import (
	"encoding/base64"
	"fmt"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ktesting "k8s.io/client-go/testing"
)

func SecretDataReactor(action ktesting.Action) (bool, runtime.Object, error) {
	secret, ok := action.(ktesting.CreateAction).GetObject().(*v1.Secret)
	if !ok {
		return false, nil, fmt.Errorf("SecretDataReactor can only be applied on secrets")
	}

	if len(secret.StringData) > 0 {
		if secret.Data == nil {
			secret.Data = make(map[string][]byte)
		}

		for k, v := range secret.StringData {
			secret.Data[k] = make([]byte, base64.StdEncoding.EncodedLen(len(v)))
			base64.StdEncoding.Encode(secret.Data[k], []byte(v))
		}
	}

	return false, nil, nil
}
