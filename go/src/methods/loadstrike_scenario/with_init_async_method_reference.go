package loadstrike_scenario

        import (
        	"time"

        	loadstrike "loadstrike.com/sdk/go"
        )

        const withInitAsyncRunnerKey = "runner_dummy_orders_reference"

type WithInitAsyncMethodReference struct{}

type withInitAsyncTempConfigPaths struct {
	ConfigPath string
	InfraPath  string
}

type withInitAsyncOrdersReportingSink struct{}

func newWithInitAsyncOrdersReportingSink() withInitAsyncOrdersReportingSink {
	return withInitAsyncOrdersReportingSink{}
}
func (withInitAsyncOrdersReportingSink) SinkName() string {
	return "orders-sample-sink"
}
func (withInitAsyncOrdersReportingSink) Init(loadstrike.LoadStrikeBaseContext, loadstrike.IConfiguration) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withInitAsyncOrdersReportingSink) Start(loadstrike.LoadStrikeSessionStartInfo) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withInitAsyncOrdersReportingSink) SaveRealtimeStats([]loadstrike.LoadStrikeScenarioStats) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withInitAsyncOrdersReportingSink) SaveRealtimeMetrics(loadstrike.LoadStrikeMetricStats) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withInitAsyncOrdersReportingSink) SaveRunResult(loadstrike.LoadStrikeRunResult) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withInitAsyncOrdersReportingSink) Stop() loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withInitAsyncOrdersReportingSink) Dispose() {}

type withInitAsyncOrdersRuntimePolicy struct{}

func newWithInitAsyncOrdersRuntimePolicy() withInitAsyncOrdersRuntimePolicy {
	return withInitAsyncOrdersRuntimePolicy{}
}
func (withInitAsyncOrdersRuntimePolicy) ShouldRunScenario(string) loadstrike.LoadStrikeBoolTask {
	return loadstrike.TaskFromBool(true)
}
func (withInitAsyncOrdersRuntimePolicy) BeforeScenario(string) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withInitAsyncOrdersRuntimePolicy) AfterScenario(string, loadstrike.LoadStrikeScenarioRuntime) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withInitAsyncOrdersRuntimePolicy) BeforeStep(string, string) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withInitAsyncOrdersRuntimePolicy) AfterStep(string, string, loadstrike.LoadStrikeReply) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}

type withInitAsyncOrdersWorkerPlugin struct{}

func newWithInitAsyncOrdersWorkerPlugin() withInitAsyncOrdersWorkerPlugin {
	return withInitAsyncOrdersWorkerPlugin{}
}
func (withInitAsyncOrdersWorkerPlugin) PluginName() string {
	return "orders-sample-plugin"
}
func (withInitAsyncOrdersWorkerPlugin) Init(loadstrike.LoadStrikeBaseContext, loadstrike.IConfiguration) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withInitAsyncOrdersWorkerPlugin) Start(loadstrike.LoadStrikeSessionStartInfo) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withInitAsyncOrdersWorkerPlugin) GetData(loadstrike.LoadStrikeRunResult) loadstrike.LoadStrikeValueTask[loadstrike.LoadStrikePluginData] {
	return loadstrike.TaskFromResult(loadstrike.CreatePluginData("orders-sample-plugin"))
}
func (withInitAsyncOrdersWorkerPlugin) Stop() loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withInitAsyncOrdersWorkerPlugin) Dispose() loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}

func withInitAsyncPerformOrderGetReply() loadstrike.LoadStrikeReply {
	return loadstrike.LoadStrikeResponse.Ok("200", int64(128), "ok", loadstrike.TimeSpan(3*time.Millisecond))
}

func withInitAsyncExecuteOrderGet(context loadstrike.LoadStrikeScenarioContext) loadstrike.LoadStrikeReply {
	return loadstrike.LoadStrikeStep.Run("get-order", context, func(loadstrike.LoadStrikeScenarioContext) loadstrike.LoadStrikeReply {
		return withInitAsyncPerformOrderGetReply()
	})
}

func withInitAsyncExecuteOrderGetAsync(context loadstrike.LoadStrikeScenarioContext) loadstrike.LoadStrikeValueTask[loadstrike.LoadStrikeReply] {
	return loadstrike.TaskFromResult(withInitAsyncExecuteOrderGet(context))
}

func withInitAsyncBaselineScenario(names ...string) loadstrike.LoadStrikeScenario {
	name := "orders.get-by-id"
	if len(names) > 0 && names[0] != "" {
		name = names[0]
	}
	return loadstrike.CreateScenario(name, withInitAsyncExecuteOrderGet).
		WithLoadSimulations(loadstrike.LoadStrikeSimulation.IterationsForConstant(1, 1)).
		WithoutWarmUp()
}

func withInitAsyncBaseRunner() loadstrike.LoadStrikeRunner {
	return loadstrike.Create().
		AddScenario(withInitAsyncBaselineScenario()).
		WithRunnerKey(withInitAsyncRunnerKey).
		WithTestSuite("orders-reference").
		WithTestName("orders-get-by-id").
		WithoutReports()
}

func withInitAsyncBaseContext() loadstrike.LoadStrikeContext {
	return withInitAsyncBaseRunner().BuildContext()
}

func withInitAsyncHttpSource() *loadstrike.EndpointSpec {
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

func withInitAsyncHttpDestination() *loadstrike.EndpointSpec {
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

func withInitAsyncTrackingConfiguration() *loadstrike.TrackingConfigurationSpec {
	return &loadstrike.TrackingConfigurationSpec{
		Source:                      withInitAsyncHttpSource(),
		Destination:                 withInitAsyncHttpDestination(),
		RunMode:                     "GenerateAndCorrelate",
		CorrelationTimeoutSeconds:   30,
		TimeoutSweepIntervalSeconds: 1,
		TimeoutBatchSize:            200,
		TimeoutCountsAsFailure:      true,
		MetricPrefix:                "orders_tracking",
		ExecuteOriginalScenarioRun:  false,
	}
}

func withInitAsyncTrackedContext() loadstrike.LoadStrikeContext {
	return loadstrike.Create().
		AddScenario(withInitAsyncBaselineScenario("orders.tracked").WithCrossPlatformTracking(withInitAsyncTrackingConfiguration())).
		WithRunnerKey(withInitAsyncRunnerKey).
		WithoutReports().
		BuildContext()
}

func withInitAsyncBuildTrackingPayload() loadstrike.TrackingPayload {
	builder := loadstrike.TrackingPayloadBuilder{
		Headers:            map[string]string{"X-Correlation-Id": "ord-1"},
		ContentType:        "application/json",
		MessagePayloadType: "json",
	}
	builder.SetBody(map[string]any{"trackingId": "ord-1"})
	return builder.Build()
}

func withInitAsyncRunResult() loadstrike.LoadStrikeRunResult {
	var result loadstrike.LoadStrikeRunResult
	return result
}

func withInitAsyncScenarioStats() loadstrike.LoadStrikeScenarioStats {
	var stats loadstrike.LoadStrikeScenarioStats
	return stats
}

func withInitAsyncWriteTempConfigFiles() withInitAsyncTempConfigPaths {
	return withInitAsyncTempConfigPaths{
		ConfigPath: "method-reference.loadstrike.config.json",
		InfraPath:  "method-reference.loadstrike.infra.json",
	}
}

        // Attach an asynchronous lifecycle hook with no extra side effects.
func (reference WithInitAsyncMethodReference) AttachAsyncNoOpHookExample() any {
    return withInitAsyncBaselineScenario().WithInitAsync(func(context loadstrike.LoadStrikeScenarioInitContext) loadstrike.LoadStrikeTask { return loadstrike.CompletedTask() })
}

// Attach an asynchronous lifecycle hook that still registers a metric.
func (reference WithInitAsyncMethodReference) AttachAsyncMetricHookExample() any {
    return withInitAsyncBaselineScenario().WithInitAsync(func(context loadstrike.LoadStrikeScenarioInitContext) loadstrike.LoadStrikeTask { context.RegisterMetric(loadstrike.Metric.CreateCounter("orders_seen", "count")); return loadstrike.CompletedTask() })
}
