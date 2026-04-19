package loadstrike_runner

        import (
        	"time"

        	loadstrike "loadstrike.com/sdk/go"
        )

        const withMinimumLogLevelRunnerKey = "runner_dummy_orders_reference"

type WithMinimumLogLevelMethodReference struct{}

type withMinimumLogLevelTempConfigPaths struct {
	ConfigPath string
	InfraPath  string
}

type withMinimumLogLevelOrdersReportingSink struct{}

func newWithMinimumLogLevelOrdersReportingSink() withMinimumLogLevelOrdersReportingSink {
	return withMinimumLogLevelOrdersReportingSink{}
}
func (withMinimumLogLevelOrdersReportingSink) SinkName() string {
	return "orders-sample-sink"
}
func (withMinimumLogLevelOrdersReportingSink) Init(loadstrike.LoadStrikeBaseContext, loadstrike.IConfiguration) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withMinimumLogLevelOrdersReportingSink) Start(loadstrike.LoadStrikeSessionStartInfo) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withMinimumLogLevelOrdersReportingSink) SaveRealtimeStats([]loadstrike.LoadStrikeScenarioStats) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withMinimumLogLevelOrdersReportingSink) SaveRealtimeMetrics(loadstrike.LoadStrikeMetricStats) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withMinimumLogLevelOrdersReportingSink) SaveRunResult(loadstrike.LoadStrikeRunResult) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withMinimumLogLevelOrdersReportingSink) Stop() loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withMinimumLogLevelOrdersReportingSink) Dispose() {}

type withMinimumLogLevelOrdersRuntimePolicy struct{}

func newWithMinimumLogLevelOrdersRuntimePolicy() withMinimumLogLevelOrdersRuntimePolicy {
	return withMinimumLogLevelOrdersRuntimePolicy{}
}
func (withMinimumLogLevelOrdersRuntimePolicy) ShouldRunScenario(string) loadstrike.LoadStrikeBoolTask {
	return loadstrike.TaskFromBool(true)
}
func (withMinimumLogLevelOrdersRuntimePolicy) BeforeScenario(string) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withMinimumLogLevelOrdersRuntimePolicy) AfterScenario(string, loadstrike.LoadStrikeScenarioRuntime) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withMinimumLogLevelOrdersRuntimePolicy) BeforeStep(string, string) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withMinimumLogLevelOrdersRuntimePolicy) AfterStep(string, string, loadstrike.LoadStrikeReply) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}

type withMinimumLogLevelOrdersWorkerPlugin struct{}

func newWithMinimumLogLevelOrdersWorkerPlugin() withMinimumLogLevelOrdersWorkerPlugin {
	return withMinimumLogLevelOrdersWorkerPlugin{}
}
func (withMinimumLogLevelOrdersWorkerPlugin) PluginName() string {
	return "orders-sample-plugin"
}
func (withMinimumLogLevelOrdersWorkerPlugin) Init(loadstrike.LoadStrikeBaseContext, loadstrike.IConfiguration) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withMinimumLogLevelOrdersWorkerPlugin) Start(loadstrike.LoadStrikeSessionStartInfo) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withMinimumLogLevelOrdersWorkerPlugin) GetData(loadstrike.LoadStrikeRunResult) loadstrike.LoadStrikeValueTask[loadstrike.LoadStrikePluginData] {
	return loadstrike.TaskFromResult(loadstrike.CreatePluginData("orders-sample-plugin"))
}
func (withMinimumLogLevelOrdersWorkerPlugin) Stop() loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withMinimumLogLevelOrdersWorkerPlugin) Dispose() loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}

func withMinimumLogLevelPerformOrderGetReply() loadstrike.LoadStrikeReply {
	return loadstrike.LoadStrikeResponse.Ok("200", int64(128), "ok", loadstrike.TimeSpan(3*time.Millisecond))
}

func withMinimumLogLevelExecuteOrderGet(context loadstrike.LoadStrikeScenarioContext) loadstrike.LoadStrikeReply {
	return loadstrike.LoadStrikeStep.Run("get-order", context, func(loadstrike.LoadStrikeScenarioContext) loadstrike.LoadStrikeReply {
		return withMinimumLogLevelPerformOrderGetReply()
	})
}

func withMinimumLogLevelExecuteOrderGetAsync(context loadstrike.LoadStrikeScenarioContext) loadstrike.LoadStrikeValueTask[loadstrike.LoadStrikeReply] {
	return loadstrike.TaskFromResult(withMinimumLogLevelExecuteOrderGet(context))
}

func withMinimumLogLevelBaselineScenario(names ...string) loadstrike.LoadStrikeScenario {
	name := "orders.get-by-id"
	if len(names) > 0 && names[0] != "" {
		name = names[0]
	}
	return loadstrike.CreateScenario(name, withMinimumLogLevelExecuteOrderGet).
		WithLoadSimulations(loadstrike.LoadStrikeSimulation.IterationsForConstant(1, 1)).
		WithoutWarmUp()
}

func withMinimumLogLevelBaseRunner() loadstrike.LoadStrikeRunner {
	return loadstrike.Create().
		AddScenario(withMinimumLogLevelBaselineScenario()).
		WithRunnerKey(withMinimumLogLevelRunnerKey).
		WithTestSuite("orders-reference").
		WithTestName("orders-get-by-id").
		WithoutReports()
}

func withMinimumLogLevelBaseContext() loadstrike.LoadStrikeContext {
	return withMinimumLogLevelBaseRunner().BuildContext()
}

func withMinimumLogLevelHttpSource() *loadstrike.EndpointSpec {
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

func withMinimumLogLevelHttpDestination() *loadstrike.EndpointSpec {
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

func withMinimumLogLevelTrackingConfiguration() *loadstrike.TrackingConfigurationSpec {
	return &loadstrike.TrackingConfigurationSpec{
		Source:                      withMinimumLogLevelHttpSource(),
		Destination:                 withMinimumLogLevelHttpDestination(),
		RunMode:                     "GenerateAndCorrelate",
		CorrelationTimeoutSeconds:   30,
		TimeoutSweepIntervalSeconds: 1,
		TimeoutBatchSize:            200,
		TimeoutCountsAsFailure:      true,
		MetricPrefix:                "orders_tracking",
		ExecuteOriginalScenarioRun:  false,
	}
}

func withMinimumLogLevelTrackedContext() loadstrike.LoadStrikeContext {
	return loadstrike.Create().
		AddScenario(withMinimumLogLevelBaselineScenario("orders.tracked").WithCrossPlatformTracking(withMinimumLogLevelTrackingConfiguration())).
		WithRunnerKey(withMinimumLogLevelRunnerKey).
		WithoutReports().
		BuildContext()
}

func withMinimumLogLevelBuildTrackingPayload() loadstrike.TrackingPayload {
	builder := loadstrike.TrackingPayloadBuilder{
		Headers:            map[string]string{"X-Correlation-Id": "ord-1"},
		ContentType:        "application/json",
		MessagePayloadType: "json",
	}
	builder.SetBody(map[string]any{"trackingId": "ord-1"})
	return builder.Build()
}

func withMinimumLogLevelRunResult() loadstrike.LoadStrikeRunResult {
	var result loadstrike.LoadStrikeRunResult
	return result
}

func withMinimumLogLevelScenarioStats() loadstrike.LoadStrikeScenarioStats {
	var stats loadstrike.LoadStrikeScenarioStats
	return stats
}

func withMinimumLogLevelWriteTempConfigFiles() withMinimumLogLevelTempConfigPaths {
	return withMinimumLogLevelTempConfigPaths{
		ConfigPath: "method-reference.loadstrike.config.json",
		InfraPath:  "method-reference.loadstrike.infra.json",
	}
}

        // Lower the logger to warning-level output.
func (reference WithMinimumLogLevelMethodReference) SetWarningLevelExample() any {
    return withMinimumLogLevelBaseRunner().WithMinimumLogLevel("Warning")
}

// Raise the logger detail for troubleshooting a sample run.
func (reference WithMinimumLogLevelMethodReference) SetDebugLevelExample() any {
    return withMinimumLogLevelBaseRunner().WithMinimumLogLevel("Debug")
}
