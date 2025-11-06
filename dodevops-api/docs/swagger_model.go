package docs

import "k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

// K8sUnstructured is an alias for unstructured.Unstructured
// swagger:model
type K8sUnstructured unstructured.Unstructured

// K8sUnstructuredList is an alias for unstructured.UnstructuredList
// swagger:model
type K8sUnstructuredList unstructured.UnstructuredList
