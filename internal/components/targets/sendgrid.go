package targets

import (
	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/zclconf/go-cty/cty"

	"bridgedl/internal/sdk/k8s"
	"bridgedl/internal/sdk/secrets"
	"bridgedl/translation"
)

type Sendgrid struct{}

var (
	_ translation.Decodable    = (*Sendgrid)(nil)
	_ translation.Translatable = (*Sendgrid)(nil)
	_ translation.Addressable  = (*Sendgrid)(nil)
)

// Spec implements translation.Decodable.
func (*Sendgrid) Spec() hcldec.Spec {
	return &hcldec.ObjectSpec{
		"default_from_email": &hcldec.AttrSpec{
			Name:     "default_from_email",
			Type:     cty.String,
			Required: false,
		},
		"default_from_name": &hcldec.AttrSpec{
			Name:     "default_from_name",
			Type:     cty.String,
			Required: false,
		},
		"default_subject": &hcldec.AttrSpec{
			Name:     "default_subject",
			Type:     cty.String,
			Required: false,
		},
		"default_to_email": &hcldec.AttrSpec{
			Name:     "default_to_email",
			Type:     cty.String,
			Required: false,
		},
		"default_to_name": &hcldec.AttrSpec{
			Name:     "default_to_name",
			Type:     cty.String,
			Required: false,
		},
		"auth": &hcldec.AttrSpec{
			Name:     "auth",
			Type:     k8s.ObjectReferenceCty,
			Required: true,
		},
	}
}

// Manifests implements translation.Translatable.
func (*Sendgrid) Manifests(id string, config, eventDst cty.Value) []interface{} {
	var manifests []interface{}

	name := k8s.RFC1123Name(id)

	t := k8s.NewObject(k8s.APITargets, "SendgridTarget", name)

	if v := config.GetAttr("default_from_email"); !v.IsNull() {
		fromEmail := v.AsString()
		t.SetNestedField(fromEmail, "spec", "defaultFromEmail")
	}

	if v := config.GetAttr("default_from_name"); !v.IsNull() {
		fromName := v.AsString()
		t.SetNestedField(fromName, "spec", "defaultFromName")
	}

	if v := config.GetAttr("default_subject"); !v.IsNull() {
		subject := v.AsString()
		t.SetNestedField(subject, "spec", "defaultSubject")
	}

	if v := config.GetAttr("default_to_email"); !v.IsNull() {
		toEmail := v.AsString()
		t.SetNestedField(toEmail, "spec", "defaultToEmail")
	}

	if v := config.GetAttr("default_to_name"); !v.IsNull() {
		toName := v.AsString()
		t.SetNestedField(toName, "spec", "defaultToName")
	}

	authSecretName := config.GetAttr("auth").GetAttr("name").AsString()
	apiKeySecretRef := secrets.SecretKeyRefsSendgrid(authSecretName)
	t.SetNestedMap(apiKeySecretRef, "spec", "apiKey", "secretKeyRef")

	manifests = append(manifests, t.Unstructured())

	if !eventDst.IsNull() {
		ch := k8s.NewChannel(name)
		subs := k8s.NewSubscription(name, name, k8s.NewDestination(k8s.APITargets, "SendgridTarget", name), eventDst)
		manifests = append(manifests, ch, subs)
	}

	return manifests
}

// Address implements translation.Addressable.
func (*Sendgrid) Address(id string, _, eventDst cty.Value) cty.Value {
	name := k8s.RFC1123Name(id)

	if eventDst.IsNull() {
		return k8s.NewDestination(k8s.APITargets, "SendgridTarget", name)
	}
	return k8s.NewDestination(k8s.APIMessaging, "Channel", name)
}