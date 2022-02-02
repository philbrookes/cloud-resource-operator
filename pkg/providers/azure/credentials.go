package azure

import (
	"context"
	"fmt"
	"time"

	"k8s.io/apimachinery/pkg/api/errors"

	"k8s.io/apimachinery/pkg/util/wait"

	"k8s.io/apimachinery/pkg/types"

	"sigs.k8s.io/controller-runtime/pkg/client"

	v1 "github.com/openshift/cloud-credential-operator/pkg/apis/cloudcredential/v1"
	errorUtil "github.com/pkg/errors"
	v12 "k8s.io/api/core/v1"
	controllerruntime "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

const (
	defaultProviderCredentialName = "cloud-resources-azure-credentials"
	defaultCredentialsKeyIDName = "azure_client_id"
	defaultCredentialsSecretKeyName = "azure_client_secret"
	DefaultAzureRole = "phils-awesome-role"
)

type Credentials struct {
	ClientID        string
	ClientSecret      string
	Region     string
	ResourcePrefix string
	ResourceGroup string
	SubscriptionId string
	TenantId string
}

//go:generate moq -out credentials_moq.go . CredentialManager
type CredentialManager interface {
	ReconcileProviderCredentials(ctx context.Context, ns string) (*Credentials, error)
	ReconcileSESCredentials(ctx context.Context, name, ns string) (*Credentials, error)
	ReconcileBucketOwnerCredentials(ctx context.Context, name, ns string) (*Credentials, *v1.CredentialsRequest, error)
	ReconcileCredentials(ctx context.Context, name, ns string) (*v1.CredentialsRequest, *Credentials, error)
}

var _ CredentialManager = (*CredentialMinterCredentialManager)(nil)

// CredentialMinterCredentialManager Implementation of CredentialManager using the openshift cloud credential minter
type CredentialMinterCredentialManager struct {
	ProviderCredentialName string
	Client                 client.Client
}

func NewCredentialMinterCredentialManager(client client.Client) *CredentialMinterCredentialManager {
	return &CredentialMinterCredentialManager{
		ProviderCredentialName: defaultProviderCredentialName,
		Client:                 client,
	}
}

//ReconcileProviderCredentials Ensure the credentials the AWS provider requires are available
func (m *CredentialMinterCredentialManager) ReconcileProviderCredentials(ctx context.Context, ns string) (*Credentials, error) {
	_, creds, err := m.ReconcileCredentials(ctx, m.ProviderCredentialName, ns)
	if err != nil {
		return nil, err
	}
	return creds, nil
}

func (m *CredentialMinterCredentialManager) ReconcileSESCredentials(ctx context.Context, name, ns string) (*Credentials, error) {
	_, creds, err := m.ReconcileCredentials(ctx, name, ns)
	if err != nil {
		return nil, err
	}
	return creds, nil
}

func (m *CredentialMinterCredentialManager) ReconcileBucketOwnerCredentials(ctx context.Context, name, ns string) (*Credentials, *v1.CredentialsRequest, error) {
	cr, creds, err := m.ReconcileCredentials(ctx, name, ns)
	if err != nil {
		return nil, nil, err
	}
	return creds, cr, nil
}

func (m *CredentialMinterCredentialManager) ReconcileCredentials(ctx context.Context, name string, ns string) (*v1.CredentialsRequest, *Credentials, error) {
	cr, err := m.reconcileCredentialRequest(ctx, name, ns)
	if err != nil {
		return nil, nil, errorUtil.Wrapf(err, "failed to reconcile aws credential request %s", name)
	}
	err = wait.PollImmediate(time.Second*5, time.Minute*5, func() (done bool, err error) {
		if err = m.Client.Get(ctx, types.NamespacedName{Name: cr.Name, Namespace: cr.Namespace}, cr); err != nil {
			if errors.IsNotFound(err) {
				return false, nil
			}
			return false, err
		}
		return cr.Status.Provisioned, nil
	})
	if err != nil {
		return nil, nil, errorUtil.Wrap(err, "timed out waiting for credential request to become provisioned")
	}

	credentials := &v12.Secret{}
	err = m.Client.Get(ctx, types.NamespacedName{Namespace: cr.Spec.SecretRef.Namespace, Name: cr.Spec.SecretRef.Name}, credentials)
	if err != nil {
		return nil, nil, errorUtil.Wrap(err, "unable to retrieve azure credentials secret")
	}

	return cr, &Credentials{
		ClientID:       string(credentials.Data["azure_client_id"]),
		ClientSecret:   string(credentials.Data["azure_client_secret"]),
		Region:         string(credentials.Data["azure_region"]),
		ResourcePrefix: string(credentials.Data["azure_resource_prefix"]),
		ResourceGroup:  string(credentials.Data["azure_resourcegroup"]),
		SubscriptionId: string(credentials.Data["azure_subscription_id"]),
		TenantId:       string(credentials.Data["azure_tenant_id"]),
	}, nil
}

func (m *CredentialMinterCredentialManager) reconcileCredentialRequest(ctx context.Context, name string, ns string) (*v1.CredentialsRequest, error) {
	codec, err := v1.NewCodec()
	if err != nil {
		return nil, errorUtil.Wrap(err, "failed to create provider codec")
	}
	providerSpec, err := codec.EncodeProviderSpec(&v1.AzureProviderSpec{
		TypeMeta: controllerruntime.TypeMeta{
			Kind: "AzureProviderSpec",
		},
		RoleBindings: []v1.RoleBinding{
			{Role: DefaultAzureRole},
		},

	})
	if err != nil {
		return nil, errorUtil.Wrap(err, "failed to encode provider spec")
	}
	cr := &v1.CredentialsRequest{
		ObjectMeta: controllerruntime.ObjectMeta{
			Name:      name,
			Namespace: ns,
		},
	}
	_, err = controllerutil.CreateOrUpdate(ctx, m.Client, cr, func() error {
		cr.Spec.ProviderSpec = providerSpec
		cr.Spec.SecretRef = v12.ObjectReference{
			Name:      name,
			Namespace: ns,
		}
		return nil
	})
	if err != nil {
		return nil, errorUtil.Wrapf(err, "failed to reconcile credential request %s in namespace %s", cr.Name, cr.Namespace)
	}
	return cr, nil
}

func (m *CredentialMinterCredentialManager) reconcileAzureCredentials(ctx context.Context, cr *v1.CredentialsRequest) (string, string, error) {
	sec := &v12.Secret{}
	err := m.Client.Get(ctx, types.NamespacedName{Name: cr.Spec.SecretRef.Name, Namespace: cr.Spec.SecretRef.Namespace}, sec)
	if err != nil {
		return "", "", errorUtil.Wrapf(err, "failed to get azure credentials secret %s", cr.Spec.SecretRef.Name)
	}
	azureAccessKeyID := string(sec.Data[defaultCredentialsKeyIDName])
	azureSecretAccessKey := string(sec.Data[defaultCredentialsSecretKeyName])
	if azureAccessKeyID == "" {
		return "", "", errorUtil.New(fmt.Sprintf("aws access key id is undefined in secret %s", sec.Name))
	}
	if azureSecretAccessKey == "" {
		return "", "", errorUtil.New(fmt.Sprintf("aws secret access key is undefined in secret %s", sec.Name))
	}
	return azureAccessKeyID, azureSecretAccessKey, nil
}
