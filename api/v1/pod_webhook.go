package v1

import (
	"context"
	admissionv1 "k8s.io/api/admission/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/json"
	"net/http"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

type PodAnnotator struct {
	Client  client.Client
	decoder *admission.Decoder
}

var (
	log = ctrl.Log.WithName("webhook")
)

//+kubebuilder:webhook:path=/mutate-v1-pod,mutating=true,sideEffects=NoneOnDryRun,admissionReviewVersions=v1,failurePolicy=fail,groups="",resources=pods,verbs=create;update;delete;connect,versions=v1,name=xy.meteor.io

func (a *PodAnnotator) Handle(ctx context.Context, req admission.Request) admission.Response {
	if req.Kind.Kind != "Pod" {
		return admission.Allowed("not a pod,skip")
	}
	var pod, oldPod *corev1.Pod
	var err error
	pod, err = decodePod(req.Object, a.decoder)
	if err != nil {
		return admission.Errored(http.StatusBadRequest, err)
	}
	oldPod, err = decodePod(req.OldObject, a.decoder)
	if err != nil {
		return admission.Errored(http.StatusBadRequest, err)
	}
	switch req.Operation {
	case admissionv1.Create:
		handleCreate(pod)
	case admissionv1.Update:
		handleUpdate(pod, oldPod)
	case admissionv1.Delete:
		handleDelete(oldPod)
		return admission.Allowed("skip")
	default:
		return admission.Allowed("skip")
	}
	//在 pod 中修改字段
	marshaledPod, err := json.Marshal(pod)
	if err != nil {
		return admission.Errored(http.StatusInternalServerError, err)
	}
	return admission.PatchResponseFromRaw(req.Object.Raw, marshaledPod)
}
func decodePod(raw runtime.RawExtension, decoder *admission.Decoder) (*corev1.Pod, error) {
	if len(raw.Raw) == 0 {
		return nil, nil
	}

	pod := &corev1.Pod{}
	err := decoder.DecodeRaw(raw, pod)
	if err != nil {
		return nil, err
	}
	return pod, nil
}

func handleDelete(pod *corev1.Pod) {

}

func handleUpdate(pod *corev1.Pod, oldPod *corev1.Pod) {

}

func handleCreate(pod *corev1.Pod) {

}

func (a *PodAnnotator) InjectDecoder(d *admission.Decoder) error {
	a.decoder = d
	return nil
}
