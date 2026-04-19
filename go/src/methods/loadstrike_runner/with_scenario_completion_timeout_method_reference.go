package loadstrike_runner

        import (
        	"time"

        	loadstrike "loadstrike.com/sdk/go"
        )

        const withScenarioCompletionTimeoutRunnerKey = "runner_dummy_orders_reference"

type WithScenarioCompletionTimeoutMethodReference struct{}

type withScenarioCompletionTimeoutTempConfigPaths struct {
	ConfigPath string
	InfraPath  string
}

type withScenarioCompletionTimeoutOrdersReportingSink struct{}

func newWithScenarioCompletionTimeoutOrdersReportingSink() withScenarioCompletionTimeoutOrdersReportingSink {
	return withScenarioCompletionTimeoutOrdersReportingSink{}
}
func (withScenarioCompletionTimeoutOrdersReportingSink) SinkName() string {
	return "orders-sample-sink"
}
func (withScenarioCompletionTimeoutOrdersReportingSink) Init(loadstrike.LoadStrikeBaseContext, loadstrike.IConfiguration) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withScenarioCompletionTimeoutOrdersReportingSink) Start(loadstrike.LoadStrikeSessionStartInfo) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withScenarioCompletionTimeoutOrdersReportingSink) SaveRealtimeStats([]loadstrike.LoadStrikeScenarioStats) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withScenarioCompletionTimeoutOrdersReportingSink) SaveRealtimeMetrics(loadstrike.LoadStrikeMetricStats) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withScenarioCompletionTimeoutOrdersReportingSink) SaveRunResult(loadstrike.LoadStrikeRunResult) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withScenarioCompletionTimeoutOrdersReportingSink) Stop() loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withScenarioCompletionTimeoutOrdersReportingSink) Dispose() {}

type withScenarioCompletionTimeoutOrdersRuntimePolicy struct{}

func newWithScenarioCompletionTimeoutOrdersRuntimePolicy() withScenarioCompletionTimeoutOrdersRuntimePolicy {
	return withScenarioCompletionTimeoutOrdersRuntimePolicy{}
}
func (withScenarioCompletionTimeoutOrdersRuntimePolicy) ShouldRunScenario(string) loadstrike.LoadStrikeBoolTask {
	return loadstrike.TaskFromBool(true)
}
func (withScenarioCompletionTimeoutOrdersRuntimePolicy) BeforeScenario(string) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withScenarioCompletionTimeoutOrdersRuntimePolicy) AfterScenario(string, loadstrike.LoadStrikeScenarioRuntime) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withScenarioCompletionTimeoutOrdersRuntimePolicy) BeforeStep(string, string) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withScenarioCompletionTimeoutOrdersRuntimePolicy) AfterStep(string, string, loadstrike.LoadStrikeReply) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}

type withScenarioCompletionTimeoutOrdersWorkerPlugin struct{}

func newWithScenarioCompletionTimeoutOrdersWorkerPlugin() withScenarioCompletionTimeoutOrdersWorkerPlugin {
	return withScenarioCompletionTimeoutOrdersWorkerPlugin{}
}
func (withScenarioCompletionTimeoutOrdersWorkerPlugin) PluginName() string {
	return "orders-sample-plugin"
}
func (withScenarioCompletionTimeoutOrdersWorkerPlugin) Init(loadstrike.LoadStrikeBaseContext, loadstrike.IConfiguration) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withScenarioCompletionTimeoutOrdersWorkerPlugin) Start(loadstrike.LoadStrikeSessionStartInfo) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withScenarioCompletionTimeoutOrdersWorkerPlugin) GetData(loadstrike.LoadStrikeRunResult) loadstrike.LoadStrikeValueTask[loadstrike.LoadStrikePluginData] {
	return loadstrike.TaskFromResult(loadstrike.CreatePluginData("orders-sample-plugin"))
}
func (withScenarioCompletionTimeoutOrdersWorkerPlugin) Stop() loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withScenarioCompletionTimeoutOrdersWorkerPlugin) Dispose() loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}

func withScenarioCompletionTimeoutPerformOrderGetReply() loadstrike.LoadStrikeReply {
	return loadstrike.LoadStrikeResponse.Ok("200", int64(128), "ok", loadstrike.TimeSpan(3*time.Millisecond))
}

func withScenarioCompletionTimeoutExecuteOrderGet(context loadstrike.LoadStrikeScenarioContext) loadstrike.LoadStrikeReply {
	return loadstrike.LoadStrikeStep.Run("get-order", context, func(loadstrike.LoadStrikeScenarioContext) loadstrike.LoadStrikeReply {
		return withScenarioCompletionTimeoutPerformOrderGetReply()
	})
}

func withScenarioCompletionTimeoutExecuteOrderGetAsync(context loadstrike.LoadStrikeScenarioContext) loadstrike.LoadStrikeValueTask[loadstrike.LoadStrikeReply] {
	return loadstrike.TaskFromResult(withScenarioCompletionTimeoutExecuteOrderGet(context))
}

func withScenarioCompletionTimeoutBaselineScenario(names ...string) loadstrike.LoadStrikeScenario {
	name := "orders.get-by-id"
	if len(names) > 0 && names[0] != "" {
		name = names[0]
	}
	return loadstrike.CreateScenario(name, withScenarioCompletionTimeoutExecuteOrderGet).
		WithLoadSimulations(loadstrike.LoadStrikeSimulation.IterationsForConstant(1, 1)).
		WithoutWarmUp()
}

func withScenarioCompletionTimeoutBaseRunner() loadstrike.LoadStrikeRunner {
	return loadstrike.Create().
		AddScenario(withScenarioCompletionTimeoutBaselineScenario()).
		WithRunnerKey(withScenarioCompletionTimeoutRunnerKey).
		WithTestSuite("orders-reference").
		WithTestName("orders-get-by-id").
		WithoutReports()
}

func withScenarioCompletionTimeoutBaseContext() loadstrike.LoadStrikeContext {
	return withScenarioCompletionTimeoutBaseRunner().BuildContext()
}

func withScenarioCompletionTimeoutHttpSource() *loadstrike.EndpointSpec {
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

func withScenarioCompletionTimeoutHttpDestination() *loadstrike.EndpointSpec {
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

func withScenarioCompletionTimeoutTrackingConfiguration() *loadstrike.TrackingConfigurationSpec {
	return &loadstrike.TrackingConfigurationSpec{
		Source:                      withScenarioCompletionTimeoutHttpSource(),
		Destination:                 withScenarioCompletionTimeoutHttpDestination(),
		RunMode:                     "GenerateAndCorrelate",
		CorrelationTimeoutSeconds:   30,
		TimeoutSweepIntervalSeconds: 1,
		TimeoutBatchSize:            200,
		TimeoutCountsAsFailure:      true,
		MetricPrefix:                "orders_tracking",
		ExecuteOriginalScenarioRun:  false,
	}
}

func withScenarioCompletionTimeoutTrackedContext() loadstrike.LoadStrikeContext {
	return loadstrike.Create().
		AddScenario(withScenarioCompletionTimeoutBaselineScenario("orders.tracked").WithCrossPlatformTracking(withScenarioCompletionTimeoutTrackingConfiguration())).
		WithRunnerKey(withScenarioCompletionTimeoutRunnerKey).
		WithoutReports().
		BuildContext()
}

func withScenarioCompletionTimeoutBuildTrackingPayload() loadstrike.TrackingPayload {
	builder := loadstrike.TrackingPayloadBuilder{
		Headers:            map[string]string{"X-Correlation-Id": "ord-1"},
		ContentType:        "application/json",
		MessagePayloadType: "json",
	}
	builder.SetBody(map[string]any{"trackingId": "ord-1"})
	return builder.Build()
}

func withScenarioCompletionTimeoutRunResult() loadstrike.LoadStrikeRunResult {
	var result loadstrike.LoadStrikeRunResult
	return result
}

func withScenarioCompletionTimeoutScenarioStats() loadstrike.LoadStrikeScenarioStats {
	var stats loadstrike.LoadStrikeScenarioStats
	return stats
}

func withScenarioCompletionTimeoutWriteTempConfigFiles() withScenarioCompletionTimeoutTempConfigPaths {
	return withScenarioCompletionTimeoutTempConfigPaths{
		ConfigPath: "method-reference.loadstrike.config.json",
		InfraPath:  "method-reference.loadstrike.infra.json",
	}
}

        // Apply a shorter timeout to the runner.
func (reference WithScenarioCompletionTimeoutMethodReference) ApplyShortTimeoutExample() any {
    return withScenarioCompletionTimeoutBaseRunner().WithScenarioCompletionTimeout(loadstrike.DurationFromSeconds(5))
}

// Apply a longer timeout to the runner.
func (reference WithScenarioCompletionTimeoutMethodReference) ApplyLongTimeoutExample() any {
    return withScenarioCompletionTimeoutBaseRunner().WithScenarioCompletionTimeout(loadstrike.DurationFromSeconds(15))
}
