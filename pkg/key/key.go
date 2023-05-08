package key

import (
	"fmt"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	capi "sigs.k8s.io/cluster-api/api/v1beta1"
)

func GetClusterIDFromLabels(t v1.ObjectMeta) string {
	return t.GetLabels()[capi.ClusterLabelName]
}

func KubeconfigSecretName(clusterName string) string {
	return fmt.Sprintf("%s-kubeconfig", clusterName)
}
