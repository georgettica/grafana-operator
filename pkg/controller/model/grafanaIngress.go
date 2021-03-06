package model

import (
	"github.com/integr8ly/grafana-operator/v3/pkg/apis/integreatly/v1alpha1"
	"k8s.io/api/extensions/v1beta1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func GetIngressPathType(cr *v1alpha1.Grafana) *v1beta1.PathType {
	t := v1beta1.PathType(cr.Spec.Ingress.PathType)
	switch t {
	case v1beta1.PathTypeExact, v1beta1.PathTypePrefix:
		return &t
	}
	t = v1beta1.PathTypeImplementationSpecific
	return &t
}

func GetIngressClassName(cr *v1alpha1.Grafana) *string {
	if cr.Spec.Ingress.IngressClassName == "" {
		return nil
	}

	return &cr.Spec.Ingress.IngressClassName
}

func getIngressTLS(cr *v1alpha1.Grafana) []v1beta1.IngressTLS {
	if cr.Spec.Ingress == nil {
		return nil
	}

	if cr.Spec.Ingress.TLSEnabled {
		return []v1beta1.IngressTLS{
			{
				Hosts:      []string{cr.Spec.Ingress.Hostname},
				SecretName: cr.Spec.Ingress.TLSSecretName,
			},
		}
	}
	return nil
}

func getIngressSpec(cr *v1alpha1.Grafana) v1beta1.IngressSpec {
	serviceName := func(cr *v1alpha1.Grafana) string {
		if cr.Spec.Service != nil && cr.Spec.Service.Name != "" {
			return cr.Spec.Service.Name
		}
		return GrafanaServiceName
	}
	return v1beta1.IngressSpec{
		TLS:              getIngressTLS(cr),
		IngressClassName: GetIngressClassName(cr),
		Rules: []v1beta1.IngressRule{
			{
				Host: GetHost(cr),
				IngressRuleValue: v1beta1.IngressRuleValue{
					HTTP: &v1beta1.HTTPIngressRuleValue{
						Paths: []v1beta1.HTTPIngressPath{
							{
								Path:     GetPath(cr),
								PathType: GetIngressPathType(cr),
								Backend: v1beta1.IngressBackend{
									ServiceName: serviceName(cr),
									ServicePort: GetIngressTargetPort(cr),
								},
							},
						},
					},
				},
			},
		},
	}
}

func GrafanaIngress(cr *v1alpha1.Grafana) *v1beta1.Ingress {
	return &v1beta1.Ingress{
		ObjectMeta: v1.ObjectMeta{
			Name:        GrafanaIngressName,
			Namespace:   cr.Namespace,
			Labels:      GetIngressLabels(cr),
			Annotations: GetIngressAnnotations(cr, nil),
		},
		Spec: getIngressSpec(cr),
	}
}

func GrafanaIngressReconciled(cr *v1alpha1.Grafana, currentState *v1beta1.Ingress) *v1beta1.Ingress {
	reconciled := currentState.DeepCopy()
	reconciled.Labels = GetIngressLabels(cr)
	reconciled.Annotations = GetIngressAnnotations(cr, currentState.Annotations)
	reconciled.Spec = getIngressSpec(cr)
	return reconciled
}

func GrafanaIngressSelector(cr *v1alpha1.Grafana) client.ObjectKey {
	return client.ObjectKey{
		Namespace: cr.Namespace,
		Name:      GrafanaIngressName,
	}
}
