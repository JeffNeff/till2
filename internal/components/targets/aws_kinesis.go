package targets

import (
	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/zclconf/go-cty/cty"

	"bridgedl/internal/sdk/k8s"
	"bridgedl/internal/sdk/secrets"
	"bridgedl/translation"
)

type AWSKinesis struct{}

var (
	_ translation.Decodable    = (*AWSKinesis)(nil)
	_ translation.Translatable = (*AWSKinesis)(nil)
	_ translation.Addressable  = (*AWSKinesis)(nil)
)

// Spec implements translation.Decodable.
func (*AWSKinesis) Spec() hcldec.Spec {
	return &hcldec.ObjectSpec{
		"arn": &hcldec.AttrSpec{
			Name:     "arn",
			Type:     cty.String,
			Required: true,
		},
		"partition": &hcldec.AttrSpec{
			Name:     "partition",
			Type:     cty.String,
			Required: true,
		},
		"credentials": &hcldec.AttrSpec{
			Name:     "credentials",
			Type:     k8s.ObjectReferenceCty,
			Required: true,
		},
	}
}

// Manifests implements translation.Translatable.
func (*AWSKinesis) Manifests(id string, config, eventDst cty.Value) []interface{} {
	var manifests []interface{}

	name := k8s.RFC1123Name(id)

	t := k8s.NewObject(k8s.APITargets, "AWSKinesisTarget", name)

	arn := config.GetAttr("arn").AsString()
	t.SetNestedField(arn, "spec", "arn")

	partition := config.GetAttr("partition").AsString()
	t.SetNestedField(partition, "spec", "partition")

	credsSecretName := config.GetAttr("credentials").GetAttr("name").AsString()
	accKeySecretRef, secrKeySecretRef := secrets.SecretKeyRefsAWS(credsSecretName)
	t.SetNestedMap(accKeySecretRef, "spec", "awsApiKey", "secretKeyRef")
	t.SetNestedMap(secrKeySecretRef, "spec", "awsApiSecret", "secretKeyRef")

	manifests = append(manifests, t.Unstructured())

	if !eventDst.IsNull() {
		ch := k8s.NewChannel(name)
		subs := k8s.NewSubscription(name, name, k8s.NewDestination(k8s.APITargets, "AWSKinesisTarget", name), eventDst)
		manifests = append(manifests, ch, subs)
	}

	return manifests
}

// Address implements translation.Addressable.
func (*AWSKinesis) Address(id string, _, eventDst cty.Value) cty.Value {
	name := k8s.RFC1123Name(id)

	if eventDst.IsNull() {
		return k8s.NewDestination(k8s.APITargets, "AWSKinesisTarget", name)
	}
	return k8s.NewDestination(k8s.APIMessaging, "Channel", name)
}
