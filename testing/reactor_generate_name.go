package testing

import (
	"fmt"

	"k8s.io/apimachinery/pkg/runtime"
	krand "k8s.io/apimachinery/pkg/util/rand"
	ktesting "k8s.io/client-go/testing"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var stringRandFunc func(n int) string = krand.String

// GenerateNameReactor sets the metav1.Name of an object, if metav1.GenerateName was used.
// It returns "handled" == false, so the test client can continue to the next ReactionFunc.
func GenerateNameReactor(action ktesting.Action) (bool, runtime.Object, error) {
	obj := action.(ktesting.CreateAction).GetObject().(client.Object)
	if obj.GetName() == "" && obj.GetGenerateName() != "" {
		obj.SetName(fmt.Sprintf("%s%s", obj.GetGenerateName(), stringRandFunc(8)))
	}

	return false, nil, nil
}
