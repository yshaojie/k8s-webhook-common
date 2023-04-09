package pods

import (
	api "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

type ContainerVisitorWithPath func(container *api.Container, path *field.Path) bool

func VisitContainersWithPath(podSpec *api.PodSpec, specPath *field.Path, visitor ContainerVisitorWithPath) bool {
	fldPath := specPath.Child("initContainers")
	for i := range podSpec.InitContainers {
		if !visitor(&podSpec.InitContainers[i], fldPath.Index(i)) {
			return false
		}
	}
	fldPath = specPath.Child("containers")
	for i := range podSpec.Containers {
		if !visitor(&podSpec.Containers[i], fldPath.Index(i)) {
			return false
		}
	}
	fldPath = specPath.Child("ephemeralContainers")
	for i := range podSpec.EphemeralContainers {
		if !visitor((*api.Container)(&podSpec.EphemeralContainers[i].EphemeralContainerCommon), fldPath.Index(i)) {
			return false
		}
	}
	return true
}
