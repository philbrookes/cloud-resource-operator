package azure

import (
	"context"
	"fmt"
	flexibleservers "github.com/Azure/azure-sdk-for-go/services/preview/postgresql/mgmt/2020-02-14-preview/postgresqlflexibleservers"
	"github.com/Azure/go-autorest/autorest/to"
	errorUtil "github.com/pkg/errors"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	croType "github.com/integr8ly/cloud-resource-operator/apis/integreatly/v1alpha1/types"
	"github.com/integr8ly/cloud-resource-operator/pkg/resources"
	"time"

	"github.com/integr8ly/cloud-resource-operator/apis/integreatly/v1alpha1"
	"github.com/integr8ly/cloud-resource-operator/pkg/providers"
	"github.com/sirupsen/logrus"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	postgresProviderName                 = "azure-postgres"
	defaultCredSecSuffix                 = "-azure-rds-credentials"
	defaultPostgresUserKey               = "user"
	defaultPostgresPasswordKey           = "password"
	defaultAzurePostgresUser             = "postgres"
	DefaultAzurePostgresSKU				 = "Standard_D4s_v3"
	DefaultAzurePostgresTier			 = "GeneralPurpose"
	DefaultAzurePostgresStorageMB		 = 524288
	DefaultAzurePostgresBackupDays		 = 14
)

var (
	defaultSupportedEngineVersions = []string{"10.16", "10.15", "10.13", "10.6", "9.6", "9.5"}
	healthyAzureDBInstanceStatuses   = []string{
		"available",
	}
)

var _ providers.PostgresProvider = (*PostgresProvider)(nil)

type PostgresProvider struct {
	Client            client.Client
	Logger            *logrus.Entry
	CredentialManager CredentialManager
	ConfigManager     ConfigManager
}

func (p *PostgresProvider) DeletePostgres(ctx context.Context, ps *v1alpha1.Postgres) (croType.StatusMessage, error) {
	panic("implement me")
}

func NewAzurePostgresProvider(client client.Client, logger *logrus.Entry) *PostgresProvider {
	return &PostgresProvider{
		Client:            client,
		Logger:            logger.WithFields(logrus.Fields{"provider": postgresProviderName}),
		CredentialManager: NewCredentialMinterCredentialManager(client),
		ConfigManager:     NewDefaultConfigMapConfigManager(client),
	}
}

func (p *PostgresProvider) GetName() string {
	return postgresProviderName
}

func (p *PostgresProvider) SupportsStrategy(d string) bool {
	return d == providers.AzureDeploymentStrategy
}

func (p *PostgresProvider) GetReconcileTime(pg *v1alpha1.Postgres) time.Duration {
	if pg.Status.Phase != croType.PhaseComplete {
		return time.Second * 60
	}
	return resources.GetForcedReconcileTimeOrDefault(defaultReconcileTime)
}

// ReconcilePostgres creates an RDS Instance from strategy config
func (p *PostgresProvider) ReconcilePostgres(ctx context.Context, pg *v1alpha1.Postgres) (*providers.PostgresInstance, croType.StatusMessage, error) {
	logger := p.Logger.WithField("action", "ReconcilePostgres")
	logger.Infof("reconciling postgres %s", pg.Name)

	// handle provider-specific finalizer
	if err := resources.CreateFinalizer(ctx, p.Client, pg, DefaultFinalizer); err != nil {
		return nil, "failed to set finalizer", err
	}

	// create the credentials to be used by the aws resource providers, not to be used by end-user
	providerCreds, err := p.CredentialManager.ReconcileProviderCredentials(ctx, pg.Namespace)
	if err != nil {
		msg := "failed to reconcile rds credentials"
		return nil, croType.StatusMessage(msg), errorUtil.Wrap(err, msg)
	}

	environment, err := resources.GetAzureEnvironment(ctx, p.Client)
	if err != nil {
		msg := "failed to retrieve cluster environment"
		return nil, croType.StatusMessage(msg), errorUtil.Wrap(err, msg)
	}

	// create credentials secret
	sec := buildDefaultRDSSecret(pg)
	or, err := controllerutil.CreateOrUpdate(ctx, p.Client, sec, func() error {
		return nil
	})
	if err != nil {
		errMsg := fmt.Sprintf("failed to create or update secret %s, action was %s", sec.Name, or)
		return nil, croType.StatusMessage(errMsg), errorUtil.Wrapf(err, errMsg)
	}

	// setup aws RDS instance sdk session
	serversClient, err := CreateServersClient(providerCreds, environment)
	if err != nil {
		errMsg := "failed to create Azure servers client to create rds db instance"
		return nil, croType.StatusMessage(errMsg), errorUtil.Wrap(err, errMsg)
	}


	future, err := serversClient.Create(
		ctx,
		providerCreds.ResourceGroup,
		pg.Name,
		flexibleservers.Server{
			Sku: &flexibleservers.Sku{
				Name: to.StringPtr(DefaultAzurePostgresSKU),
				Tier: DefaultAzurePostgresTier,
			},
			ServerProperties: &flexibleservers.ServerProperties{
				AdministratorLogin:         to.StringPtr(sec.StringData["user"]),
				AdministratorLoginPassword: to.StringPtr(sec.StringData["password"]),
				Version:                    flexibleservers.OneTwo,
				StorageProfile:             &flexibleservers.StorageProfile{
					StorageMB: to.Int32Ptr(DefaultAzurePostgresStorageMB),
					BackupRetentionDays: to.Int32Ptr(DefaultAzurePostgresBackupDays),
				},
			},
			Location: to.StringPtr(providerCreds.Region),
		})

	if err != nil {
		return &providers.PostgresInstance{}, croType.StatusMessage("cannot update pg server"), err
	}

	if err := future.WaitForCompletionRef(ctx, serversClient.Client); err != nil {
		return &providers.PostgresInstance{}, croType.StatusMessage("cannot get the pg server update future response"), err
	}

	server, err := future.Result(serversClient)
	if err != nil {
		return &providers.PostgresInstance{}, croType.StatusMessage("cannot update pg server"), err
	}

	logger.Infof("providerCreds: %+v, server: %+v", providerCreds, server)
	return &providers.PostgresInstance{}, croType.StatusMessage(""), nil

}

func buildDefaultRDSSecret(ps *v1alpha1.Postgres) *v1.Secret {
	password, err := resources.GeneratePassword()
	if err != nil {
		return nil
	}
	return &v1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      ps.Name + defaultCredSecSuffix,
			Namespace: ps.Namespace,
		},
		StringData: map[string]string{
			defaultPostgresUserKey:     defaultAzurePostgresUser,
			defaultPostgresPasswordKey: password,
		},
		Type: v1.SecretTypeOpaque,
	}
}