package cmd

import (
	"context"
	"errors"
	"fmt"
	"math/bits"
	"os"
	"os/signal"
	sync2 "sync"
	"syscall"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/robfig/cron/v3"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/raito-io/cli/internal/clitrigger"
	"github.com/raito-io/cli/internal/constants"
	"github.com/raito-io/cli/internal/health_check"
	"github.com/raito-io/cli/internal/logging"
	"github.com/raito-io/cli/internal/target"
	"github.com/raito-io/cli/internal/target/types"
	"github.com/raito-io/cli/internal/target_sync"
	"github.com/raito-io/cli/internal/version"
	"github.com/raito-io/cli/internal/version_management"
)

func initRunCommand(rootCmd *cobra.Command) {
	var cmd = &cobra.Command{
		Hidden: false,
		Use:    "run",
		Short:  "Run all the configured synchronizations",
		Long:   `Run all the configured synchronizations`,
		Run:    executeRun,
	}

	cmd.PersistentFlags().String(constants.IdentityStoreIdFlag, "", "The ID of the identity store in Raito to import the user and group information to. This is only applicable if specifying the (single) target information in the commandline.")
	cmd.PersistentFlags().String(constants.DataSourceIdFlag, "", "The ID of the data source in Raito to import the meta data information to or get the access permissions from. This is only applicable if specifying the (single) target information in the commandline.")
	cmd.PersistentFlags().StringP(constants.DomainFlag, "d", "", "The subdomain to your Raito instance (https://<subdomain>.raito.io). This parameter can be overridden in the target configs if needed.")
	cmd.PersistentFlags().StringP(constants.ApiUserFlag, "u", "", "The username of the API user to authenticate against Raito. This parameter can be overridden in the target configs if needed.")
	cmd.PersistentFlags().StringP(constants.ApiSecretFlag, "s", "", "The API key secret to authenticate against Raito. This parameter can be overridden in the target configs if needed.")
	cmd.PersistentFlags().String(constants.URLOverrideFlag, "", "")
	cmd.PersistentFlags().Bool(constants.SkipAuthentication, false, "")
	cmd.PersistentFlags().Bool(constants.SkipFileUpload, false, "")
	cmd.PersistentFlags().StringP(constants.OnlyTargetsFlag, "t", "", "Can be used to only execute a subset of the defined targets in the configuration file. To specify multiple, use a comma-separated list.")
	cmd.PersistentFlags().String(constants.ConnectorNameFlag, "", "The name of the connector to use. If not set, the CLI will use a configuration file to define the targets.")
	cmd.PersistentFlags().String(constants.ConnectorVersionFlag, "", "The version of the connector to use. This is only relevant if the 'connector' flag is set as well. If not set (but the 'connector' flag is), then 'latest' is used.")
	cmd.PersistentFlags().StringP(constants.NameFlag, "n", "", "The name for the target. This is only relevant if the 'connector' flag is set as well. If not set, the name of the connector will be used.")
	cmd.PersistentFlags().String(constants.ContainerLivenessFile, "", "If set, we will create/remove a health-check file based on the webhook state. This is only relevant if you are running the CLI in long running mode.")

	cmd.PersistentFlags().StringP(constants.CronFlag, "c", "", "If set, the cron expression will define when a sync should run. When not set (and no frequency is defined), the sync will run once and quit after. (e.g. '0 0/2 * * *' initiates a sync evey 2 hours)")
	cmd.PersistentFlags().Bool(constants.SyncAtStartupFlag, false, "If set, a sync will be run at startup independent of the cron expression. Only applicable if cron expression is defined.")
	cmd.PersistentFlags().IntP(constants.FrequencyFlag, "f", 0, "The frequency used to do the sync (in minutes). When not set (and no cron expression is defined), the default value '0' is used, which means the sync will run once and quit after.")
	cmd.PersistentFlags().Bool(constants.SkipDataSourceSyncFlag, false, "If set, the data source meta data synchronization step to Raito will be skipped for each of the targets.")
	cmd.PersistentFlags().Bool(constants.SkipIdentityStoreSyncFlag, false, "If set, the identity store synchronization step to Raito will be skipped for each of the targets.")
	cmd.PersistentFlags().Bool(constants.SkipDataAccessSyncFlag, false, "If set, the data access information from Raito will not be synced to the data sources in the target list.")
	cmd.PersistentFlags().Bool(constants.SkipDataUsageSyncFlag, false, "If set, the data usage information synchronization step to Raito will be skipped for each of the targets.")
	cmd.PersistentFlags().Bool(constants.SkipResourceProviderFlag, false, "If set, the resource provider synchronization step to Raito will be skipped for each of the targets.")
	cmd.PersistentFlags().Bool(constants.SkipTagFlag, false, "If set, the tags synchronization step to Raito will be skipped for each of the targets")

	cmd.PersistentFlags().Bool(constants.LockAllWhoFlag, false, "If set, the 'who' (users and groups) of all access providers imported into Raito Cloud will be locked. Note that this only makes sense for access providers that represent a named entity (like a Snowflake Role or AWS Policy).")
	cmd.PersistentFlags().String(constants.LockWhoByNameFlag, "", "Allows you to specify a comma-separated list of access provider names for which the 'who' (users and groups) should be locked when imported into Raito Cloud. The names in the list are interpreted as regular expressions which allows for partial matches (e.g. '.+-prod,.+-dev' will match all access providers ending with '-prod' or '-dev'). Note that this only makes sense for access providers that represent a named entity (like a Snowflake Role or AWS Policy).")
	cmd.PersistentFlags().String(constants.LockWhoByTagFlag, "", "Allows you to specify a comma-separated list of access provider tags for which the 'who' (users and groups) should be locked when imported into Raito Cloud. The tags in the list are interpreted as regular expressions which allows for partial matches. The format for an item should be in the form 'key:value' (e.g. 'key1:value1,key2:.+' will match all access providers that have the tag 'key1:value1' or have the tag with key 'key2'). Note that this only makes sense for access providers that represent a named entity (like a Snowflake Role or AWS Policy).")
	cmd.PersistentFlags().Bool(constants.LockWhoWhenIncompleteFlag, false, "If set, the 'who' (users and groups) of all access providers imported into Raito Cloud will be locked if the access provider is incomplete (not all elements could be understood by Raito). This can be used to protect accidental removal of permissions unknown to Raito by blocking the possibility to edit it. Note that this only makes sense for access providers that represent a named entity (like a Snowflake Role or AWS Policy).")

	cmd.PersistentFlags().Bool(constants.LockAllInheritanceFlag, false, fmt.Sprintf("Same as %q, but for the 'inheritance' of the access providers.", constants.LockAllWhoFlag))
	cmd.PersistentFlags().String(constants.LockInheritanceByNameFlag, "", fmt.Sprintf("Same as %q, but for the 'inheritance' of the access providers.", constants.LockWhoByNameFlag))
	cmd.PersistentFlags().String(constants.LockInheritanceByTagFlag, "", fmt.Sprintf("Same as %q, but for the 'inheritance' of the access providers.", constants.LockWhoByTagFlag))
	cmd.PersistentFlags().Bool(constants.LockInheritanceWhenIncompleteFlag, false, fmt.Sprintf("Same as %q, but for the 'inheritance' of the access providers.", constants.LockWhoWhenIncompleteFlag))

	cmd.PersistentFlags().Bool(constants.LockAllWhatFlag, false, fmt.Sprintf("Same as %q, but for the 'what' of the access providers.", constants.LockAllWhoFlag))
	cmd.PersistentFlags().String(constants.LockWhatByNameFlag, "", fmt.Sprintf("Same as %q, but for the 'what' of the access providers.", constants.LockWhoByNameFlag))
	cmd.PersistentFlags().String(constants.LockWhatByTagFlag, "", fmt.Sprintf("Same as %q, but for the 'what' of the access providers.", constants.LockWhoByTagFlag))
	cmd.PersistentFlags().Bool(constants.LockWhatWhenIncompleteFlag, false, fmt.Sprintf("Same as %q, but for the 'what' of the access providers.", constants.LockWhoWhenIncompleteFlag))

	cmd.PersistentFlags().Bool(constants.LockAllNamesFlag, false, fmt.Sprintf("Same as %q, but for the name of the access providers.", constants.LockAllWhoFlag))
	cmd.PersistentFlags().String(constants.LockNamesByNameFlag, "", fmt.Sprintf("Same as %q, but for the name of the access providers.", constants.LockWhoByNameFlag))
	cmd.PersistentFlags().String(constants.LockNamesByTagFlag, "", fmt.Sprintf("Same as %q, but for the name of the access providers.", constants.LockWhoByTagFlag))
	cmd.PersistentFlags().Bool(constants.LockNamesWhenIncompleteFlag, false, fmt.Sprintf("Same as %q, but for the name of the access providers.", constants.LockWhoWhenIncompleteFlag))

	cmd.PersistentFlags().Bool(constants.LockAllDeleteFlag, false, fmt.Sprintf("Same as %q, but for deleting the access providers. Note: setting the delete lock on an access control also means that when it is internalized in Raito by a user, it will be switched back to external once a successful sync is done.", constants.LockAllWhoFlag))
	cmd.PersistentFlags().String(constants.LockDeleteByNameFlag, "", fmt.Sprintf("Same as %q, but for deleting the access providers. Note: setting the delete lock on an access control also means that when it is internalized in Raito by a user, it will be switched back to external once a successful sync is done.", constants.LockWhoByNameFlag))
	cmd.PersistentFlags().String(constants.LockDeleteByTagFlag, "", fmt.Sprintf("Same as %q, but for deleting the access providers. Note: setting the delete lock on an access control also means that when it is internalized in Raito by a user, it will be switched back to external once a successful sync is done.", constants.LockWhoByTagFlag))
	cmd.PersistentFlags().Bool(constants.LockDeleteWhenIncompleteFlag, false, fmt.Sprintf("Same as %q, but for deleting the access providers. Note: setting the delete lock on an access control also means that when it is internalized in Raito by a user, it will be switched back to external once a successful sync is done.", constants.LockWhoWhenIncompleteFlag))

	cmd.PersistentFlags().Bool(constants.LockAllOwnersFlag, false, "If set, the owners of all access providers imported into Raito Cloud will be locked.")

	cmd.PersistentFlags().String(constants.MakeNotInternalizableFlag, "", "Allows you to specify a comma-separated list of access provider names that should be made not-internalizable when imported into Raito Cloud. This means that these access providers will not be editable in Raito Cloud. The names in the list are interpreted as regular expressions so allow for partial matches. (e.g. '.+-prod,.+-dev' will match all access providers ending with '-prod' or '-dev')")
	cmd.PersistentFlags().Lookup(constants.MakeNotInternalizableFlag).Deprecated = fmt.Sprintf("use %q instead", constants.FullyLockByNameFlag)
	cmd.PersistentFlags().Bool(constants.FullyLockAllFlag, false, fmt.Sprintf("Same as %q, but will fully lock the access providers from being edited in Raito Cloud.", constants.LockAllWhoFlag))
	cmd.PersistentFlags().String(constants.FullyLockByNameFlag, "", fmt.Sprintf("Same as %q, but will fully lock the access providers from being edited in Raito Cloud.", constants.LockWhoByNameFlag))
	cmd.PersistentFlags().String(constants.FullyLockByTagFlag, "", fmt.Sprintf("Same as %q, but will fully lock the access providers from being edited in Raito Cloud.", constants.LockWhoByTagFlag))
	cmd.PersistentFlags().Bool(constants.FullyLockWhenIncompleteFlag, false, fmt.Sprintf("Same as %q, but will fully lock the access providers from being edited in Raito Cloud.", constants.LockWhoWhenIncompleteFlag))

	cmd.PersistentFlags().Bool(constants.DisableWebsocketFlag, false, "If set, raito will not setup a websocket to trigger new syncs. This flag has only effect if frequency is set.")
	cmd.PersistentFlags().Bool(constants.DisableLogForwarding, false, "If set, sync logs will not be forwarded to Raito Cloud.")
	cmd.PersistentFlags().Bool(constants.DisableLogForwardingDataSourceSync, false, "If set, data source sync logs will not be forwarded to Raito Cloud.")
	cmd.PersistentFlags().Bool(constants.DisableLogForwardingDataAccessSync, false, "If set, data access sync logs will not be forwarded to Raito Cloud.")
	cmd.PersistentFlags().Bool(constants.DisableLogForwardingIdentityStoreSync, false, "If set, identity store sync logs will not be forwarded to Raito Cloud.")
	cmd.PersistentFlags().Bool(constants.DisableLogForwardingDataUsageSync, false, "If set, data usage sync logs will not be forwarded to Raito Cloud.")
	cmd.PersistentFlags().Bool(constants.DisableLogForwardingResourceProviderSync, false, "If set, resource provider synchronization logs will not be forwarded to Raito Cloud.")
	cmd.PersistentFlags().Bool(constants.DisableLogForwardingTagSync, false, "If set, tag synchronization logs will not be forwarded to Raito Cloud.")

	cmd.PersistentFlags().String(constants.TagOverwriteKeyForAccessProviderName, "", "If set, will determine the tag-key used for overwriting the display-name of the Access Control when imported in to Raito Cloud.")
	cmd.PersistentFlags().String(constants.TagOverwriteKeyForAccessProviderOwners, "", "If set, will determine the tag-key used for assigning owners of the Access Control when imported in to Raito Cloud.")
	cmd.PersistentFlags().String(constants.TagOverwriteKeyForDataObjectOwners, "", "If set, will determine the tag-key used for assigning owners of the Data Objects when imported in to Raito Cloud.")

	cmd.PersistentFlags().String(constants.TagKeyAndValueForUserIsMachine, "", "If set, we automatically flag a user as machine user when the combination of tag key:value (split by `:`) is matched during the import to Raito Cloud.")

	cmd.PersistentFlags().String(constants.FileBackupLocationFlag, "", "If set, this filepath is used to store backups of the files that are used during synchronization jobs. A sub-folder is created per target, using the target name + the type of run (full, manual or webhook) as name for the folder. Underneath that, another sub-folder is created per run, using a timestamp as the folder name. The backed up files are then stored in that folder. This parameter can be overridden in the target configs if needed.")
	cmd.PersistentFlags().Int(constants.MaximumBackupsPerTargetFlag, 0, fmt.Sprintf("When %q is defined, this parameter can be used to control how many backups should be kept per target+type. When this number is exceeded, older backups will be removed automatically. By default, this is 0, which means there is no maximum. This parameter can be overridden in the target configs if needed.", constants.FileBackupLocationFlag))

	BindFlag(constants.IdentityStoreIdFlag, cmd)
	BindFlag(constants.DataSourceIdFlag, cmd)
	BindFlag(constants.OnlyTargetsFlag, cmd)
	BindFlag(constants.ConnectorNameFlag, cmd)
	BindFlag(constants.ConnectorVersionFlag, cmd)
	BindFlag(constants.NameFlag, cmd)
	BindFlag(constants.DomainFlag, cmd)
	BindFlag(constants.ApiUserFlag, cmd)
	BindFlag(constants.ApiSecretFlag, cmd)
	BindFlag(constants.URLOverrideFlag, cmd)
	BindFlag(constants.SkipAuthentication, cmd)
	BindFlag(constants.SkipFileUpload, cmd)
	BindFlag(constants.ContainerLivenessFile, cmd)

	BindFlag(constants.CronFlag, cmd)
	BindFlag(constants.SyncAtStartupFlag, cmd)
	BindFlag(constants.FrequencyFlag, cmd)
	BindFlag(constants.SkipDataSourceSyncFlag, cmd)
	BindFlag(constants.SkipIdentityStoreSyncFlag, cmd)
	BindFlag(constants.SkipDataAccessSyncFlag, cmd)
	BindFlag(constants.SkipDataUsageSyncFlag, cmd)
	BindFlag(constants.SkipResourceProviderFlag, cmd)
	BindFlag(constants.SkipTagFlag, cmd)
	BindFlag(constants.LockAllWhoFlag, cmd)
	BindFlag(constants.LockWhoByNameFlag, cmd)
	BindFlag(constants.LockWhoByTagFlag, cmd)
	BindFlag(constants.LockWhoWhenIncompleteFlag, cmd)
	BindFlag(constants.LockAllInheritanceFlag, cmd)
	BindFlag(constants.LockInheritanceByNameFlag, cmd)
	BindFlag(constants.LockInheritanceByTagFlag, cmd)
	BindFlag(constants.LockInheritanceWhenIncompleteFlag, cmd)
	BindFlag(constants.LockAllWhatFlag, cmd)
	BindFlag(constants.LockWhatByNameFlag, cmd)
	BindFlag(constants.LockWhatByTagFlag, cmd)
	BindFlag(constants.LockWhatWhenIncompleteFlag, cmd)
	BindFlag(constants.LockAllNamesFlag, cmd)
	BindFlag(constants.LockNamesByNameFlag, cmd)
	BindFlag(constants.LockNamesByTagFlag, cmd)
	BindFlag(constants.LockNamesWhenIncompleteFlag, cmd)
	BindFlag(constants.LockAllDeleteFlag, cmd)
	BindFlag(constants.LockDeleteByNameFlag, cmd)
	BindFlag(constants.LockDeleteByTagFlag, cmd)
	BindFlag(constants.LockDeleteWhenIncompleteFlag, cmd)
	BindFlag(constants.LockAllOwnersFlag, cmd)
	BindFlag(constants.MakeNotInternalizableFlag, cmd)
	BindFlag(constants.FullyLockAllFlag, cmd)
	BindFlag(constants.FullyLockByNameFlag, cmd)
	BindFlag(constants.FullyLockByTagFlag, cmd)
	BindFlag(constants.FullyLockWhenIncompleteFlag, cmd)
	BindFlag(constants.DisableWebsocketFlag, cmd)
	BindFlag(constants.DisableLogForwarding, cmd)
	BindFlag(constants.DisableLogForwardingDataSourceSync, cmd)
	BindFlag(constants.DisableLogForwardingDataAccessSync, cmd)
	BindFlag(constants.DisableLogForwardingIdentityStoreSync, cmd)
	BindFlag(constants.DisableLogForwardingDataUsageSync, cmd)
	BindFlag(constants.DisableLogForwardingResourceProviderSync, cmd)
	BindFlag(constants.DisableLogForwardingTagSync, cmd)

	BindFlag(constants.TagOverwriteKeyForAccessProviderName, cmd)
	BindFlag(constants.TagOverwriteKeyForAccessProviderOwners, cmd)
	BindFlag(constants.TagOverwriteKeyForDataObjectOwners, cmd)

	BindFlag(constants.TagKeyAndValueForUserIsMachine, cmd)

	BindFlag(constants.FileBackupLocationFlag, cmd)
	BindFlag(constants.MaximumBackupsPerTargetFlag, cmd)

	hideConfigOptions(cmd, constants.URLOverrideFlag, constants.SkipAuthentication, constants.SkipFileUpload, constants.ContainerLivenessFile)

	cmd.FParseErrWhitelist.UnknownFlags = true

	rootCmd.AddCommand(cmd)
}

func executeRun(cmd *cobra.Command, args []string) {
	logging.SetupLogging(false)

	otherArgs := cmd.Flags().Args()

	ctx := context.Background()

	baseLogger := hclog.L()
	healthChecker := createHealthChecker(baseLogger)

	baseConfig, err := target.BuildBaseConfigFromFlags(baseLogger, healthChecker, otherArgs)
	if err != nil {
		hclog.L().Error(err.Error())
		os.Exit(1)
	}

	executeSyncAtStartup, scheduler, err := createSyncScheduler(baseConfig)
	if err != nil {
		hclog.L().Error(err.Error())
		os.Exit(1)
	}

	if scheduler == nil {
		hclog.L().Info("Running synchronization just once.")

		baseConfig.BaseLogger = baseConfig.BaseLogger.With("iteration", 0)
		err = executeSingleRun(ctx, baseConfig)

		if err != nil {
			os.Exit(1)
		} else {
			os.Exit(0)
		}
	} else {
		executeContinuousRun(ctx, executeSyncAtStartup, scheduler, baseConfig)
	}
}

func executeContinuousRun(ctx context.Context, executeSyncAtStartup bool, scheduler cron.Schedule, baseConfig *types.BaseConfig) {
	hclog.L().Info("Starting continuous synchronization.")
	hclog.L().Info("Press 'ctrl+c' to stop the program.")

	cancelCtx, cancelFn := context.WithCancel(ctx)

	waitGroup := sync2.WaitGroup{}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGUSR1, os.Interrupt)

	returnSignal := 0

	waitGroup.Add(1)

	go func() {
		defer waitGroup.Done()

		cliTriggerCtx, cliTriggerCancel := context.WithCancel(cancelCtx)
		cliTrigger, apUpdateTrigger, syncTrigger := startListingToCliTriggers(cliTriggerCtx, baseConfig)

		defer cliTriggerCancel()

		if syncTrigger == nil {
			cancelFn()
			return
		}

		defer func() {
			cliTrigger.Wait()
			syncTrigger.Close()
			apUpdateTrigger.Close()
		}()

		it := 1

		if executeSyncAtStartup {
			baseConfig.BaseLogger = baseConfig.BaseLogger.With("iteration", it)
			if runErr := executeSingleRun(cancelCtx, baseConfig); runErr != nil {
				baseConfig.BaseLogger.Error(fmt.Sprintf("Run failed: %s", runErr.Error()))
			}

			it++
		}

		timer := cronTimer(baseConfig.BaseLogger, nil, scheduler)
		defer timer.Stop()

		for {
			select {
			case <-timer.C:
				cliTrigger.Reset()

				baseConfig.BaseLogger = baseConfig.BaseLogger.With("iteration", it)
				if runErr := executeSingleRun(cancelCtx, baseConfig); runErr != nil {
					baseConfig.BaseLogger.Error(fmt.Sprintf("Run failed: %s", runErr.Error()))
				}

				cronTimer(baseConfig.BaseLogger, timer, scheduler)

				it++
			case <-apUpdateTrigger.TriggerChannel():
				apUpdate := apUpdateTrigger.Pop()
				if apUpdate == nil {
					continue
				}

				baseConfig.BaseLogger = baseConfig.BaseLogger.With("iteration", it)
				err := handleApUpdateTrigger(cancelCtx, baseConfig, apUpdate)

				if err != nil {
					baseConfig.BaseLogger.Warn(fmt.Sprintf("ClI ApUpdate Trigger failed: %s", err.Error()))
				}

				it++
			case <-syncTrigger.TriggerChannel():
				syncRequest := syncTrigger.Pop()
				if syncRequest == nil {
					continue
				}

				baseConfig.BaseLogger = baseConfig.BaseLogger.With("iteration", it)
				err := handleSyncTrigger(cancelCtx, baseConfig, syncRequest)

				if err != nil {
					baseConfig.BaseLogger.Warn(fmt.Sprintf("ClI Sync Trigger failed: %s", err.Error()))
				}

				it++
			case <-cancelCtx.Done():
				baseConfig.BaseLogger.Debug("Context done: closing syncing routine.")
				return
			}

			hclog.L().Info("Press 'ctrl+c' to stop the program.")
		}
	}()

	waitGroup.Add(1)

	go func() {
		defer waitGroup.Done()
		defer cancelFn()
		defer hclog.L().Info("Waiting for the current synchronization run to end ...")

		for {
			select {
			case <-cancelCtx.Done():
				hclog.L().Debug("Context done: Will stop all running routines...")
				return
			case s := <-sigs:
				hclog.L().Debug(fmt.Sprintf("Received signal: %s. Will stop all running routines...", s.String()))

				if sysc, ok := s.(syscall.Signal); ok {
					returnSignal = int(sysc)
				}

				return
			}
		}
	}()

	waitGroup.Wait()
	hclog.L().Info("All routines finished. Bye!")

	if returnSignal != 0 {
		hclog.L().Debug(fmt.Sprintf("Exit with code: %d", returnSignal))
		syscall.Exit(returnSignal)
	}
}

func createHealthChecker(baseLogger hclog.Logger) health_check.HealthChecker {
	livenessFilePath := viper.GetString(constants.ContainerLivenessFile)

	return health_check.NewHealthChecker(baseLogger, livenessFilePath)
}

func createSyncScheduler(baseConfig *types.BaseConfig) (bool, cron.Schedule, error) {
	executeSyncAtStartup := viper.GetBool(constants.SyncAtStartupFlag)
	var scheduler cron.Schedule

	cronExpression := viper.GetString(constants.CronFlag)
	freq := viper.GetInt(constants.FrequencyFlag)

	if cronExpression == "" {
		if freq > 0 {
			if freq < 60 {
				return false, nil, fmt.Errorf("the 'frequency' flag must be at least 60 seconds. The value is: %d", freq)
			}

			executeSyncAtStartup = true
			scheduler = cron.Every(time.Minute * time.Duration(freq))
		}
	} else {
		if freq > 0 {
			baseConfig.BaseLogger.Warn("The 'frequency' flag is ignored when the 'cron' flag is set.")
		}

		var cronParserErr error
		scheduler, cronParserErr = cron.ParseStandard(cronExpression)

		specSchedule, ok := scheduler.(*cron.SpecSchedule)
		if ok {
			if moreThanOneExecutionWithinAnHour(specSchedule) {
				return false, nil, errors.New("cron expression will trigger sync multiple times within an hour")
			}
		}

		if cronParserErr != nil {
			return false, nil, cronParserErr
		}
	}

	return executeSyncAtStartup, scheduler, nil
}

func moreThanOneExecutionWithinAnHour(cronSchedule *cron.SpecSchedule) bool {
	prefix := uint64(1) << 63

	if bits.OnesCount64(cronSchedule.Second&^prefix) > 1 || bits.OnesCount64(cronSchedule.Minute&^prefix) > 1 {
		return true
	}

	return false
}

func cronTimer(logger hclog.Logger, timer *time.Timer, scheduler cron.Schedule) *time.Timer {
	next := scheduler.Next(time.Now())

	logger.Info(fmt.Sprintf("Next execution at %s", next.Format(time.RFC822)))

	waitTime := time.Until(next)

	if timer == nil {
		return time.NewTimer(waitTime)
	} else {
		timer.Reset(waitTime)
		return timer
	}
}

func executeSingleRun(ctx context.Context, baseconfig *types.BaseConfig) error {
	start := time.Now()

	err := runSync(ctx, baseconfig)

	sec := time.Since(start).Round(time.Millisecond)
	baseconfig.BaseLogger.Info(fmt.Sprintf("Finished execution of all targets in %s", sec))

	return err
}

func runSync(ctx context.Context, baseconfig *types.BaseConfig) error {
	compatibilityInformation, err := version_management.IsCompatibleWithRaitoCloud(baseconfig)
	if err != nil {
		baseconfig.BaseLogger.Error(fmt.Sprintf("Failed to check compatibility with Raito Cloud: %s", err.Error()))

		return fmt.Errorf("compatibility check failed: %s", err.Error())
	}

	switch compatibilityInformation.Compatibility {
	case version_management.NotSupported:
		baseconfig.BaseLogger.Error(fmt.Sprintf("CLI version is not compatible with Raito Cloud. Please upgrade to a supported version (%s).", compatibilityInformation.SupportedVersions))

		return errors.New("unsupported CLI version")
	case version_management.Deprecated:
		warning := " "
		if compatibilityInformation.DeprecatedWarningMsg != nil {
			warning += *compatibilityInformation.DeprecatedWarningMsg
		}

		baseconfig.BaseLogger.Warn(fmt.Sprintf("CLI version %s is deprecated.%s Please upgrade to supported version (%s) soon.", version.GetCliVersion().String(), warning, compatibilityInformation.SupportedVersions))

		fallthrough
	case version_management.Supported:
		return target.RunTargets(ctx, baseconfig, target_sync.NewSyncJob("full"))
	case version_management.CompatibilityUnknown:
	}

	return errors.New("unknown CLI version")
}

func handleApUpdateTrigger(ctx context.Context, config *types.BaseConfig, apUpdate *clitrigger.ApUpdate) error {
	return target.RunTargets(ctx, config, target_sync.NewSyncJob("webhook"), target.WithDataSourceIds(apUpdate.DataSourceNames...), target.WithConfigOption(func(targetConfig *types.BaseTargetConfig) {
		targetConfig.SkipIdentityStoreSync = true
		targetConfig.SkipDataSourceSync = true
		targetConfig.SkipDataUsageSync = true
		targetConfig.SkipResourceProvider = true

		targetConfig.SkipDataAccessImport = true
		targetConfig.OnlyOutOfSyncData = true
	}))
}

func handleSyncTrigger(ctx context.Context, config *types.BaseConfig, syncTrigger *clitrigger.SyncTrigger) error {
	opts := []func(*target.Options){
		target.WithConfigOption(func(targetConfig *types.BaseTargetConfig) {
			targetConfig.SkipIdentityStoreSync = !syncTrigger.IdentityStoreSync
			targetConfig.SkipDataSourceSync = !syncTrigger.DataSourceSync
			targetConfig.SkipDataUsageSync = !syncTrigger.DataUsageSync
			targetConfig.SkipDataAccessSync = !syncTrigger.DataAccessSync
			targetConfig.DataObjectParent = syncTrigger.DataObjectParent
			targetConfig.DataObjectExcludes = syncTrigger.DataObjectExcludes
			targetConfig.SkipResourceProvider = !syncTrigger.ResourceProviderSync
		}),
	}

	if syncTrigger.IdentityStore != nil {
		opts = append(opts, target.WithIdentityStoreIds(*syncTrigger.IdentityStore))
	}

	if syncTrigger.DataSource != nil {
		opts = append(opts, target.WithDataSourceIds(*syncTrigger.DataSource))
	}

	return target.RunTargets(ctx, config, target_sync.NewSyncJob("manual"), opts...)
}

func startListingToCliTriggers(ctx context.Context, baseConfig *types.BaseConfig) (clitrigger.CliTrigger, *clitrigger.ApUpdateTriggerHandler, *clitrigger.SyncTriggerHandler) {
	cliTrigger, err := clitrigger.CreateCliTrigger(baseConfig)
	if err != nil {
		baseConfig.BaseLogger.Error(fmt.Sprintf("Unable to start asynchronous access provider sync: %s", err.Error()))
		return cliTrigger, nil, nil
	}

	apUpdateTriggerHandler := clitrigger.NewApUpdateTrigger(cliTrigger)
	syncTriggerHandler := clitrigger.NewSyncTrigger(cliTrigger)

	cliTrigger.Start(ctx)

	return cliTrigger, apUpdateTriggerHandler, syncTriggerHandler
}
