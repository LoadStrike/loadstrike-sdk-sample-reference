package loadstrike_scenario

        import (
        	"time"

        	loadstrike "loadstrike.com/sdk/go"
        )

        const withCleanAsyncRunnerKey = "runner_dummy_orders_reference"

type WithCleanAsyncMethodReference struct{}

type withCleanAsyncTempConfigPaths struct {
	ConfigPath string
	InfraPath  string
}

type withCleanAsyncOrdersReportingSink struct{}

func newWithCleanAsyncOrdersReportingSink() withCleanAsyncOrdersReportingSink {
	return withCleanAsyncOrdersReportingSink{}
}
func (withCleanAsyncOrdersReportingSink) SinkName() string {
	return "orders-sample-sink"
}
func (withCleanAsyncOrdersReportingSink) Init(loadstrike.LoadStrikeBaseContext, loadstrike.IConfiguration) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withCleanAsyncOrdersReportingSink) Start(loadstrike.LoadStrikeSessionStartInfo) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withCleanAsyncOrdersReportingSink) SaveRealtimeStats([]loadstrike.LoadStrikeScenarioStats) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withCleanAsyncOrdersReportingSink) SaveRealtimeMetrics(loadstrike.LoadStrikeMetricStats) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withCleanAsyncOrdersReportingSink) SaveRunResult(loadstrike.LoadStrikeRunResult) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withCleanAsyncOrdersReportingSink) Stop() loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withCleanAsyncOrdersReportingSink) Dispose() {}

type withCleanAsyncOrdersRuntimePolicy struct{}

func newWithCleanAsyncOrdersRuntimePolicy() withCleanAsyncOrdersRuntimePolicy {
	return withCleanAsyncOrdersRuntimePolicy{}
}
func (withCleanAsyncOrdersRuntimePolicy) ShouldRunScenario(string) loadstrike.LoadStrikeBoolTask {
	return loadstrike.TaskFromBool(true)
}
func (withCleanAsyncOrdersRuntimePolicy) BeforeScenario(string) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withCleanAsyncOrdersRuntimePolicy) AfterScenario(string, loadstrike.LoadStrikeScenarioRuntime) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withCleanAsyncOrdersRuntimePolicy) BeforeStep(string, string) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withCleanAsyncOrdersRuntimePolicy) AfterStep(string, string, loadstrike.LoadStrikeReply) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}

type withCleanAsyncOrdersWorkerPlugin struct{}

func newWithCleanAsyncOrdersWorkerPlugin() withCleanAsyncOrdersWorkerPlugin {
	return withCleanAsyncOrdersWorkerPlugin{}
}
func (withCleanAsyncOrdersWorkerPlugin) PluginName() string {
	return "orders-sample-plugin"
}
func (withCleanAsyncOrdersWorkerPlugin) Init(loadstrike.LoadStrikeBaseContext, loadstrike.IConfiguration) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withCleanAsyncOrdersWorkerPlugin) Start(loadstrike.LoadStrikeSessionStartInfo) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withCleanAsyncOrdersWorkerPlugin) GetData(loadstrike.LoadStrikeRunResult) loadstrike.LoadStrikeValueTask[loadstrike.LoadStrikePluginData] {
	return loadstrike.TaskFromResult(loadstrike.CreatePluginData("orders-sample-plugin"))
}
func (withCleanAsyncOrdersWorkerPlugin) Stop() loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withCleanAsyncOrdersWorkerPlugin) Dispose() loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}

func withCleanAsyncPerformOrderGetReply() loadstrike.LoadStrikeReply {
	return loadstrike.LoadStrikeResponse.Ok("200", int64(128), "ok", loadstrike.TimeSpan(3*time.Millisecond))
}

func withCleanAsyncExecuteOrderGet(context loadstrike.LoadStrikeScenarioContext) loadstrike.LoadStrikeReply {
	return loadstrike.LoadStrikeStep.Run("get-order", context, func(loadstrike.LoadStrikeScenarioContext) loadstrike.LoadStrikeReply {
		return withCleanAsyncPerformOrderGetReply()
	})
}

func withCleanAsyncExecuteOrderGetAsync(context loadstrike.LoadStrikeScenarioContext) loadstrike.LoadStrikeValueTask[loadstrike.LoadStrikeReply] {
	return loadstrike.TaskFromResult(withCleanAsyncExecuteOrderGet(context))
}

func withCleanAsyncBaselineScenario(names ...string) loadstrike.LoadStrikeScenario {
	name := "orders.get-by-id"
	if len(names) > 0 && names[0] != "" {
		name = names[0]
	}
	return loadstrike.CreateScenario(name, withCleanAsyncExecuteOrderGet).
		WithLoadSimulations(loadstrike.LoadStrikeSimulation.IterationsForConstant(1, 1)).
		WithoutWarmUp()
}

func withCleanAsyncBaseRunner() loadstrike.LoadStrikeRunner {
	return loadstrike.Create().
		AddScenario(withCleanAsyncBaselineScenario()).
		WithRunnerKey(withCleanAsyncRunnerKey).
		WithTestSuite("orders-reference").
		WithTestName("orders-get-by-id").
		WithoutReports()
}

func withCleanAsyncBaseContext() loadstrike.LoadStrikeContext {
	return withCleanAsyncBaseRunner().BuildContext()
}

func withCleanAsyncHttpSource() *loadstrike.EndpointSpec {
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

func withCleanAsyncHttpDestination() *loadstrike.EndpointSpec {
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

func withCleanAsyncTrackingConfiguration() *loadstrike.TrackingConfigurationSpec {
	return &loadstrike.TrackingConfigurationSpec{
		Source:                      withCleanAsyncHttpSource(),
		Destination:                 withCleanAsyncHttpDestination(),
		RunMode:                     "GenerateAndCorrelate",
		CorrelationTimeoutSeconds:   30,
		TimeoutSweepIntervalSeconds: 1,
		TimeoutBatchSize:            200,
		TimeoutCountsAsFailure:      true,
		MetricPrefix:                "orders_tracking",
		ExecuteOriginalScenarioRun:  false,
	}
}

func withCleanAsyncTrackedContext() loadstrike.LoadStrikeContext {
	return loadstrike.Create().
		AddScenario(withCleanAsyncBaselineScenario("orders.tracked").WithCrossPlatformTracking(withCleanAsyncTrackingConfiguration())).
		WithRunnerKey(withCleanAsyncRunnerKey).
		WithoutReports().
		BuildContext()
}

func withCleanAsyncBuildTrackingPayload() loadstrike.TrackingPayload {
	builder := loadstrike.TrackingPayloadBuilder{
		Headers:            map[string]string{"X-Correlation-Id": "ord-1"},
		ContentType:        "application/json",
		MessagePayloadType: "json",
	}
	builder.SetBody(map[string]any{"trackingId": "ord-1"})
	return builder.Build()
}

func withCleanAsyncRunResult() loadstrike.LoadStrikeRunResult {
	var result loadstrike.LoadStrikeRunResult
	return result
}

func withCleanAsyncScenarioStats() loadstrike.LoadStrikeScenarioStats {
	var stats loadstrike.LoadStrikeScenarioStats
	return stats
}

func withCleanAsyncWriteTempConfigFiles() withCleanAsyncTempConfigPaths {
	return withCleanAsyncTempConfigPaths{
		ConfigPath: "method-reference.loadstrike.config.json",
		InfraPath:  "method-reference.loadstrike.infra.json",
	}
}

        // Attach an asynchronous lifecycle hook with no extra side effects.
func (reference WithCleanAsyncMethodReference) AttachAsyncNoOpHookExample() any {
    return withCleanAsyncBaselineScenario().WithCleanAsync(func(context loadstrike.LoadStrikeScenarioInitContext) loadstrike.LoadStrikeTask { return loadstrike.CompletedTask() })
}

// Attach an asynchronous lifecycle hook that still registers a metric.
func (reference WithCleanAsyncMethodReference) AttachAsyncMetricHookExample() any {
    return withCleanAsyncBaselineScenario().WithCleanAsync(func(context loadstrike.LoadStrikeScenarioInitContext) loadstrike.LoadStrikeTask { context.RegisterMetric(loadstrike.Metric.CreateCounter("orders_seen", "count")); return loadstrike.CompletedTask() })
}
