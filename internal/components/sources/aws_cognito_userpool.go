package sources

import (
	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/zclconf/go-cty/cty"

	"bridgedl/internal/sdk/k8s"
	"bridgedl/internal/sdk/secrets"
	"bridgedl/translation"
)

type AWSCognitoUserPool struct{}

var (
	_ translation.Decodable    = (*AWSCognitoUserPool)(nil)
	_ translation.Translatable = (*AWSCognitoUserPool)(nil)
)

// Spec implements translation.Decodable.
func (*AWSCognitoUserPool) Spec() hcldec.Spec {
	return &hcldec.ObjectSpec{
		"arn": &hcldec.AttrSpec{
			Name:     "arn",
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
func (*AWSCognitoUserPool) Manifests(id string, config, eventDst cty.Value) []interface{} {
	var manifests []interface{}

	s := k8s.NewObject(k8s.APISources, "AWSCognitoUserPoolSource", k8s.RFC1123Name(id))

	arn := config.GetAttr("arn").AsString()
	s.SetNestedField(arn, "spec", "arn")

	credsSecretName := config.GetAttr("credentials").GetAttr("name").AsString()
	accKeySecretRef, secrKeySecretRef := secrets.SecretKeyRefsAWS(credsSecretName)
	s.SetNestedMap(accKeySecretRef, "spec", "credentials", "accessKeyID", "valueFromSecret")
	s.SetNestedMap(secrKeySecretRef, "spec", "credentials", "secretAccessKey", "valueFromSecret")

	sink := k8s.DecodeDestination(eventDst)
	s.SetNestedMap(sink, "spec", "sink", "ref")

	return append(manifests, s.Unstructured())
}
