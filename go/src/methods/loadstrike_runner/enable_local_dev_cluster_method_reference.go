package loadstrike_runner

        import (
        	"time"

        	loadstrike "loadstrike.com/sdk/go"
        )

        const enableLocalDevClusterRunnerKey = "runner_dummy_orders_reference"

type EnableLocalDevClusterMethodReference struct{}

type enableLocalDevClusterTempConfigPaths struct {
	ConfigPath string
	InfraPath  string
}

type enableLocalDevClusterOrdersReportingSink struct{}

func newEnableLocalDevClusterOrdersReportingSink() enableLocalDevClusterOrdersReportingSink {
	return enableLocalDevClusterOrdersReportingSink{}
}
func (enableLocalDevClusterOrdersReportingSink) SinkName() string {
	return "orders-sample-sink"
}
func (enableLocalDevClusterOrdersReportingSink) Init(loadstrike.LoadStrikeBaseContext, loadstrike.IConfiguration) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (enableLocalDevClusterOrdersReportingSink) Start(loadstrike.LoadStrikeSessionStartInfo) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (enableLocalDevClusterOrdersReportingSink) SaveRealtimeStats([]loadstrike.LoadStrikeScenarioStats) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (enableLocalDevClusterOrdersReportingSink) SaveRealtimeMetrics(loadstrike.LoadStrikeMetricStats) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (enableLocalDevClusterOrdersReportingSink) SaveRunResult(loadstrike.LoadStrikeRunResult) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (enableLocalDevClusterOrdersReportingSink) Stop() loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (enableLocalDevClusterOrdersReportingSink) Dispose() {}

type enableLocalDevClusterOrdersRuntimePolicy struct{}

func newEnableLocalDevClusterOrdersRuntimePolicy() enableLocalDevClusterOrdersRuntimePolicy {
	return enableLocalDevClusterOrdersRuntimePolicy{}
}
func (enableLocalDevClusterOrdersRuntimePolicy) ShouldRunScenario(string) loadstrike.LoadStrikeBoolTask {
	return loadstrike.TaskFromBool(true)
}
func (enableLocalDevClusterOrdersRuntimePolicy) BeforeScenario(string) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (enableLocalDevClusterOrdersRuntimePolicy) AfterScenario(string, loadstrike.LoadStrikeScenarioRuntime) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (enableLocalDevClusterOrdersRuntimePolicy) BeforeStep(string, string) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (enableLocalDevClusterOrdersRuntimePolicy) AfterStep(string, string, loadstrike.LoadStrikeReply) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}

type enableLocalDevClusterOrdersWorkerPlugin struct{}

func newEnableLocalDevClusterOrdersWorkerPlugin() enableLocalDevClusterOrdersWorkerPlugin {
	return enableLocalDevClusterOrdersWorkerPlugin{}
}
func (enableLocalDevClusterOrdersWorkerPlugin) PluginName() string {
	return "orders-sample-plugin"
}
func (enableLocalDevClusterOrdersWorkerPlugin) Init(loadstrike.LoadStrikeBaseContext, loadstrike.IConfiguration) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (enableLocalDevClusterOrdersWorkerPlugin) Start(loadstrike.LoadStrikeSessionStartInfo) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (enableLocalDevClusterOrdersWorkerPlugin) GetData(loadstrike.LoadStrikeRunResult) loadstrike.LoadStrikeValueTask[loadstrike.LoadStrikePluginData] {
	return loadstrike.TaskFromResult(loadstrike.CreatePluginData("orders-sample-plugin"))
}
func (enableLocalDevClusterOrdersWorkerPlugin) Stop() loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (enableLocalDevClusterOrdersWorkerPlugin) Dispose() loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}

func enableLocalDevClusterPerformOrderGetReply() loadstrike.LoadStrikeReply {
	return loadstrike.LoadStrikeResponse.Ok("200", int64(128), "ok", loadstrike.TimeSpan(3*time.Millisecond))
}

func enableLocalDevClusterExecuteOrderGet(context loadstrike.LoadStrikeScenarioContext) loadstrike.LoadStrikeReply {
	return loadstrike.LoadStrikeStep.Run("get-order", context, func(loadstrike.LoadStrikeScenarioContext) loadstrike.LoadStrikeReply {
		return enableLocalDevClusterPerformOrderGetReply()
	})
}

func enableLocalDevClusterExecuteOrderGetAsync(context loadstrike.LoadStrikeScenarioContext) loadstrike.LoadStrikeValueTask[loadstrike.LoadStrikeReply] {
	return loadstrike.TaskFromResult(enableLocalDevClusterExecuteOrderGet(context))
}

func enableLocalDevClusterBaselineScenario(names ...string) loadstrike.LoadStrikeScenario {
	name := "orders.get-by-id"
	if len(names) > 0 && names[0] != "" {
		name = names[0]
	}
	return loadstrike.CreateScenario(name, enableLocalDevClusterExecuteOrderGet).
		WithLoadSimulations(loadstrike.LoadStrikeSimulation.IterationsForConstant(1, 1)).
		WithoutWarmUp()
}

func enableLocalDevClusterBaseRunner() loadstrike.LoadStrikeRunner {
	return loadstrike.Create().
		AddScenario(enableLocalDevClusterBaselineScenario()).
		WithRunnerKey(enableLocalDevClusterRunnerKey).
		WithTestSuite("orders-reference").
		WithTestName("orders-get-by-id").
		WithoutReports()
}

func enableLocalDevClusterBaseContext() loadstrike.LoadStrikeContext {
	return enableLocalDevClusterBaseRunner().BuildContext()
}

func enableLocalDevClusterHttpSource() *loadstrike.EndpointSpec {
	return &loadstrike.EndpointSpec{
		Kind:          "Http",
		Name:          "orders-http-source",
		Mode:          "Produce",
		TrackingField: "header:X-Correlation-Id",
		HTTP: &loadstrike.HTTPEndpointOptions{
			URL:                   "https://orders.example.test/api/orders",
			Method:                "GET",
			TrackingPayloadSource: "Request",
			ResponseSource:        "ResponseBody",
		},
	}
}

func enableLocalDevClusterHttpDestination() *loadstrike.EndpointSpec {
	return &loadstrike.EndpointSpec{
		Kind:          "Http",
		Name:          "orders-http-destination",
		Mode:          "Consume",
		TrackingField: "json:$.trackingId",
		GatherByField: "json:$.tenantId",
		HTTP: &loadstrike.HTTPEndpointOptions{
			URL:                      "https://orders.example.test/api/order-events",
			Method:                   "GET",
			ResponseSource:           "ResponseBody",
			ConsumeJSONArrayResponse: true,
			ConsumeArrayPath:         "$.items",
		},
	}
}

func enableLocalDevClusterTrackingConfiguration() *loadstrike.TrackingConfigurationSpec {
	return &loadstrike.TrackingConfigurationSpec{
		Source:                      enableLocalDevClusterHttpSource(),
		Destination:                 enableLocalDevClusterHttpDestination(),
		RunMode:                     "GenerateAndCorrelate",
		CorrelationTimeoutSeconds:   30,
		TimeoutSweepIntervalSeconds: 1,
		TimeoutBatchSize:            200,
		TimeoutCountsAsFailure:      true,
		MetricPrefix:                "orders_tracking",
		ExecuteOriginalScenarioRun:  false,
	}
}

func enableLocalDevClusterTrackedContext() loadstrike.LoadStrikeContext {
	return loadstrike.Create().
		AddScenario(enableLocalDevClusterBaselineScenario("orders.tracked").WithCrossPlatformTracking(enableLocalDevClusterTrackingConfiguration())).
		WithRunnerKey(enableLocalDevClusterRunnerKey).
		WithoutReports().
		BuildContext()
}

func enableLocalDevClusterBuildTrackingPayload() loadstrike.TrackingPayload {
	builder := loadstrike.TrackingPayloadBuilder{
		Headers:            map[string]string{"X-Correlation-Id": "ord-1"},
		ContentType:        "application/json",
		MessagePayloadType: "json",
	}
	builder.SetBody(map[string]any{"trackingId": "ord-1"})
	return builder.Build()
}

func enableLocalDevClusterRunResult() loadstrike.LoadStrikeRunResult {
	var result loadstrike.LoadStrikeRunResult
	return result
}

func enableLocalDevClusterScenarioStats() loadstrike.LoadStrikeScenarioStats {
	var stats loadstrike.LoadStrikeScenarioStats
	return stats
}

func enableLocalDevClusterWriteTempConfigFiles() enableLocalDevClusterTempConfigPaths {
	return enableLocalDevClusterTempConfigPaths{
		ConfigPath: "method-reference.loadstrike.config.json",
		InfraPath:  "method-reference.loadstrike.infra.json",
	}
}

        // Apply the boolean option in its enabled state.
func (reference EnableLocalDevClusterMethodReference) EnableOnExample() any {
    return enableLocalDevClusterBaseRunner().EnableLocalDevCluster(true)
}

// Apply the boolean option in its disabled state.
func (reference EnableLocalDevClusterMethodReference) EnableOffExample() any {
    return enableLocalDevClusterBaseRunner().EnableLocalDevCluster(false)
}
