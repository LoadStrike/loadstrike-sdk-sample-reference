package loadstrike_runner

        import (
        	"time"

        	loadstrike "loadstrike.com/sdk/go"
        )

        const withClusterCommandTimeoutRunnerKey = "runner_dummy_orders_reference"

type WithClusterCommandTimeoutMethodReference struct{}

type withClusterCommandTimeoutTempConfigPaths struct {
	ConfigPath string
	InfraPath  string
}

type withClusterCommandTimeoutOrdersReportingSink struct{}

func newWithClusterCommandTimeoutOrdersReportingSink() withClusterCommandTimeoutOrdersReportingSink {
	return withClusterCommandTimeoutOrdersReportingSink{}
}
func (withClusterCommandTimeoutOrdersReportingSink) SinkName() string {
	return "orders-sample-sink"
}
func (withClusterCommandTimeoutOrdersReportingSink) Init(loadstrike.LoadStrikeBaseContext, loadstrike.IConfiguration) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withClusterCommandTimeoutOrdersReportingSink) Start(loadstrike.LoadStrikeSessionStartInfo) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withClusterCommandTimeoutOrdersReportingSink) SaveRealtimeStats([]loadstrike.LoadStrikeScenarioStats) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withClusterCommandTimeoutOrdersReportingSink) SaveRealtimeMetrics(loadstrike.LoadStrikeMetricStats) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withClusterCommandTimeoutOrdersReportingSink) SaveRunResult(loadstrike.LoadStrikeRunResult) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withClusterCommandTimeoutOrdersReportingSink) Stop() loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withClusterCommandTimeoutOrdersReportingSink) Dispose() {}

type withClusterCommandTimeoutOrdersRuntimePolicy struct{}

func newWithClusterCommandTimeoutOrdersRuntimePolicy() withClusterCommandTimeoutOrdersRuntimePolicy {
	return withClusterCommandTimeoutOrdersRuntimePolicy{}
}
func (withClusterCommandTimeoutOrdersRuntimePolicy) ShouldRunScenario(string) loadstrike.LoadStrikeBoolTask {
	return loadstrike.TaskFromBool(true)
}
func (withClusterCommandTimeoutOrdersRuntimePolicy) BeforeScenario(string) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withClusterCommandTimeoutOrdersRuntimePolicy) AfterScenario(string, loadstrike.LoadStrikeScenarioRuntime) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withClusterCommandTimeoutOrdersRuntimePolicy) BeforeStep(string, string) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withClusterCommandTimeoutOrdersRuntimePolicy) AfterStep(string, string, loadstrike.LoadStrikeReply) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}

type withClusterCommandTimeoutOrdersWorkerPlugin struct{}

func newWithClusterCommandTimeoutOrdersWorkerPlugin() withClusterCommandTimeoutOrdersWorkerPlugin {
	return withClusterCommandTimeoutOrdersWorkerPlugin{}
}
func (withClusterCommandTimeoutOrdersWorkerPlugin) PluginName() string {
	return "orders-sample-plugin"
}
func (withClusterCommandTimeoutOrdersWorkerPlugin) Init(loadstrike.LoadStrikeBaseContext, loadstrike.IConfiguration) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withClusterCommandTimeoutOrdersWorkerPlugin) Start(loadstrike.LoadStrikeSessionStartInfo) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withClusterCommandTimeoutOrdersWorkerPlugin) GetData(loadstrike.LoadStrikeRunResult) loadstrike.LoadStrikeValueTask[loadstrike.LoadStrikePluginData] {
	return loadstrike.TaskFromResult(loadstrike.CreatePluginData("orders-sample-plugin"))
}
func (withClusterCommandTimeoutOrdersWorkerPlugin) Stop() loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withClusterCommandTimeoutOrdersWorkerPlugin) Dispose() loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}

func withClusterCommandTimeoutPerformOrderGetReply() loadstrike.LoadStrikeReply {
	return loadstrike.LoadStrikeResponse.Ok("200", int64(128), "ok", loadstrike.TimeSpan(3*time.Millisecond))
}

func withClusterCommandTimeoutExecuteOrderGet(context loadstrike.LoadStrikeScenarioContext) loadstrike.LoadStrikeReply {
	return loadstrike.LoadStrikeStep.Run("get-order", context, func(loadstrike.LoadStrikeScenarioContext) loadstrike.LoadStrikeReply {
		return withClusterCommandTimeoutPerformOrderGetReply()
	})
}

func withClusterCommandTimeoutExecuteOrderGetAsync(context loadstrike.LoadStrikeScenarioContext) loadstrike.LoadStrikeValueTask[loadstrike.LoadStrikeReply] {
	return loadstrike.TaskFromResult(withClusterCommandTimeoutExecuteOrderGet(context))
}

func withClusterCommandTimeoutBaselineScenario(names ...string) loadstrike.LoadStrikeScenario {
	name := "orders.get-by-id"
	if len(names) > 0 && names[0] != "" {
		name = names[0]
	}
	return loadstrike.CreateScenario(name, withClusterCommandTimeoutExecuteOrderGet).
		WithLoadSimulations(loadstrike.LoadStrikeSimulation.IterationsForConstant(1, 1)).
		WithoutWarmUp()
}

func withClusterCommandTimeoutBaseRunner() loadstrike.LoadStrikeRunner {
	return loadstrike.Create().
		AddScenario(withClusterCommandTimeoutBaselineScenario()).
		WithRunnerKey(withClusterCommandTimeoutRunnerKey).
		WithTestSuite("orders-reference").
		WithTestName("orders-get-by-id").
		WithoutReports()
}

func withClusterCommandTimeoutBaseContext() loadstrike.LoadStrikeContext {
	return withClusterCommandTimeoutBaseRunner().BuildContext()
}

func withClusterCommandTimeoutHttpSource() *loadstrike.EndpointSpec {
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

func withClusterCommandTimeoutHttpDestination() *loadstrike.EndpointSpec {
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

func withClusterCommandTimeoutTrackingConfiguration() *loadstrike.TrackingConfigurationSpec {
	return &loadstrike.TrackingConfigurationSpec{
		Source:                      withClusterCommandTimeoutHttpSource(),
		Destination:                 withClusterCommandTimeoutHttpDestination(),
		RunMode:                     "GenerateAndCorrelate",
		CorrelationTimeoutSeconds:   30,
		TimeoutSweepIntervalSeconds: 1,
		TimeoutBatchSize:            200,
		TimeoutCountsAsFailure:      true,
		MetricPrefix:                "orders_tracking",
		ExecuteOriginalScenarioRun:  false,
	}
}

func withClusterCommandTimeoutTrackedContext() loadstrike.LoadStrikeContext {
	return loadstrike.Create().
		AddScenario(withClusterCommandTimeoutBaselineScenario("orders.tracked").WithCrossPlatformTracking(withClusterCommandTimeoutTrackingConfiguration())).
		WithRunnerKey(withClusterCommandTimeoutRunnerKey).
		WithoutReports().
		BuildContext()
}

func withClusterCommandTimeoutBuildTrackingPayload() loadstrike.TrackingPayload {
	builder := loadstrike.TrackingPayloadBuilder{
		Headers:            map[string]string{"X-Correlation-Id": "ord-1"},
		ContentType:        "application/json",
		MessagePayloadType: "json",
	}
	builder.SetBody(map[string]any{"trackingId": "ord-1"})
	return builder.Build()
}

func withClusterCommandTimeoutRunResult() loadstrike.LoadStrikeRunResult {
	var result loadstrike.LoadStrikeRunResult
	return result
}

func withClusterCommandTimeoutScenarioStats() loadstrike.LoadStrikeScenarioStats {
	var stats loadstrike.LoadStrikeScenarioStats
	return stats
}

func withClusterCommandTimeoutWriteTempConfigFiles() withClusterCommandTimeoutTempConfigPaths {
	return withClusterCommandTimeoutTempConfigPaths{
		ConfigPath: "method-reference.loadstrike.config.json",
		InfraPath:  "method-reference.loadstrike.infra.json",
	}
}

        // Apply a shorter timeout to the runner.
func (reference WithClusterCommandTimeoutMethodReference) ApplyShortTimeoutExample() any {
    return withClusterCommandTimeoutBaseRunner().WithClusterCommandTimeout(loadstrike.DurationFromSeconds(5))
}

// Apply a longer timeout to the runner.
func (reference WithClusterCommandTimeoutMethodReference) ApplyLongTimeoutExample() any {
    return withClusterCommandTimeoutBaseRunner().WithClusterCommandTimeout(loadstrike.DurationFromSeconds(15))
}
