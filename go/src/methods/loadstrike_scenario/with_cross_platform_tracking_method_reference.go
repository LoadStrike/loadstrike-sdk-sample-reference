package loadstrike_scenario

        import (
        	"time"

        	loadstrike "loadstrike.com/sdk/go"
        )

        const withCrossPlatformTrackingRunnerKey = "runner_dummy_orders_reference"

type WithCrossPlatformTrackingMethodReference struct{}

type withCrossPlatformTrackingTempConfigPaths struct {
	ConfigPath string
	InfraPath  string
}

type withCrossPlatformTrackingOrdersReportingSink struct{}

func newWithCrossPlatformTrackingOrdersReportingSink() withCrossPlatformTrackingOrdersReportingSink {
	return withCrossPlatformTrackingOrdersReportingSink{}
}
func (withCrossPlatformTrackingOrdersReportingSink) SinkName() string {
	return "orders-sample-sink"
}
func (withCrossPlatformTrackingOrdersReportingSink) Init(loadstrike.LoadStrikeBaseContext, loadstrike.IConfiguration) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withCrossPlatformTrackingOrdersReportingSink) Start(loadstrike.LoadStrikeSessionStartInfo) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withCrossPlatformTrackingOrdersReportingSink) SaveRealtimeStats([]loadstrike.LoadStrikeScenarioStats) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withCrossPlatformTrackingOrdersReportingSink) SaveRealtimeMetrics(loadstrike.LoadStrikeMetricStats) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withCrossPlatformTrackingOrdersReportingSink) SaveRunResult(loadstrike.LoadStrikeRunResult) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withCrossPlatformTrackingOrdersReportingSink) Stop() loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withCrossPlatformTrackingOrdersReportingSink) Dispose() {}

type withCrossPlatformTrackingOrdersRuntimePolicy struct{}

func newWithCrossPlatformTrackingOrdersRuntimePolicy() withCrossPlatformTrackingOrdersRuntimePolicy {
	return withCrossPlatformTrackingOrdersRuntimePolicy{}
}
func (withCrossPlatformTrackingOrdersRuntimePolicy) ShouldRunScenario(string) loadstrike.LoadStrikeBoolTask {
	return loadstrike.TaskFromBool(true)
}
func (withCrossPlatformTrackingOrdersRuntimePolicy) BeforeScenario(string) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withCrossPlatformTrackingOrdersRuntimePolicy) AfterScenario(string, loadstrike.LoadStrikeScenarioRuntime) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withCrossPlatformTrackingOrdersRuntimePolicy) BeforeStep(string, string) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withCrossPlatformTrackingOrdersRuntimePolicy) AfterStep(string, string, loadstrike.LoadStrikeReply) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}

type withCrossPlatformTrackingOrdersWorkerPlugin struct{}

func newWithCrossPlatformTrackingOrdersWorkerPlugin() withCrossPlatformTrackingOrdersWorkerPlugin {
	return withCrossPlatformTrackingOrdersWorkerPlugin{}
}
func (withCrossPlatformTrackingOrdersWorkerPlugin) PluginName() string {
	return "orders-sample-plugin"
}
func (withCrossPlatformTrackingOrdersWorkerPlugin) Init(loadstrike.LoadStrikeBaseContext, loadstrike.IConfiguration) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withCrossPlatformTrackingOrdersWorkerPlugin) Start(loadstrike.LoadStrikeSessionStartInfo) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withCrossPlatformTrackingOrdersWorkerPlugin) GetData(loadstrike.LoadStrikeRunResult) loadstrike.LoadStrikeValueTask[loadstrike.LoadStrikePluginData] {
	return loadstrike.TaskFromResult(loadstrike.CreatePluginData("orders-sample-plugin"))
}
func (withCrossPlatformTrackingOrdersWorkerPlugin) Stop() loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withCrossPlatformTrackingOrdersWorkerPlugin) Dispose() loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}

func withCrossPlatformTrackingPerformOrderGetReply() loadstrike.LoadStrikeReply {
	return loadstrike.LoadStrikeResponse.Ok("200", int64(128), "ok", loadstrike.TimeSpan(3*time.Millisecond))
}

func withCrossPlatformTrackingExecuteOrderGet(context loadstrike.LoadStrikeScenarioContext) loadstrike.LoadStrikeReply {
	return loadstrike.LoadStrikeStep.Run("get-order", context, func(loadstrike.LoadStrikeScenarioContext) loadstrike.LoadStrikeReply {
		return withCrossPlatformTrackingPerformOrderGetReply()
	})
}

func withCrossPlatformTrackingExecuteOrderGetAsync(context loadstrike.LoadStrikeScenarioContext) loadstrike.LoadStrikeValueTask[loadstrike.LoadStrikeReply] {
	return loadstrike.TaskFromResult(withCrossPlatformTrackingExecuteOrderGet(context))
}

func withCrossPlatformTrackingBaselineScenario(names ...string) loadstrike.LoadStrikeScenario {
	name := "orders.get-by-id"
	if len(names) > 0 && names[0] != "" {
		name = names[0]
	}
	return loadstrike.CreateScenario(name, withCrossPlatformTrackingExecuteOrderGet).
		WithLoadSimulations(loadstrike.LoadStrikeSimulation.IterationsForConstant(1, 1)).
		WithoutWarmUp()
}

func withCrossPlatformTrackingBaseRunner() loadstrike.LoadStrikeRunner {
	return loadstrike.Create().
		AddScenario(withCrossPlatformTrackingBaselineScenario()).
		WithRunnerKey(withCrossPlatformTrackingRunnerKey).
		WithTestSuite("orders-reference").
		WithTestName("orders-get-by-id").
		WithoutReports()
}

func withCrossPlatformTrackingBaseContext() loadstrike.LoadStrikeContext {
	return withCrossPlatformTrackingBaseRunner().BuildContext()
}

func withCrossPlatformTrackingHttpSource() *loadstrike.EndpointSpec {
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

func withCrossPlatformTrackingHttpDestination() *loadstrike.EndpointSpec {
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

func withCrossPlatformTrackingTrackingConfiguration() *loadstrike.TrackingConfigurationSpec {
	return &loadstrike.TrackingConfigurationSpec{
		Source:                      withCrossPlatformTrackingHttpSource(),
		Destination:                 withCrossPlatformTrackingHttpDestination(),
		RunMode:                     "GenerateAndCorrelate",
		CorrelationTimeoutSeconds:   30,
		TimeoutSweepIntervalSeconds: 1,
		TimeoutBatchSize:            200,
		TimeoutCountsAsFailure:      true,
		MetricPrefix:                "orders_tracking",
		ExecuteOriginalScenarioRun:  false,
	}
}

func withCrossPlatformTrackingTrackedContext() loadstrike.LoadStrikeContext {
	return loadstrike.Create().
		AddScenario(withCrossPlatformTrackingBaselineScenario("orders.tracked").WithCrossPlatformTracking(withCrossPlatformTrackingTrackingConfiguration())).
		WithRunnerKey(withCrossPlatformTrackingRunnerKey).
		WithoutReports().
		BuildContext()
}

func withCrossPlatformTrackingBuildTrackingPayload() loadstrike.TrackingPayload {
	builder := loadstrike.TrackingPayloadBuilder{
		Headers:            map[string]string{"X-Correlation-Id": "ord-1"},
		ContentType:        "application/json",
		MessagePayloadType: "json",
	}
	builder.SetBody(map[string]any{"trackingId": "ord-1"})
	return builder.Build()
}

func withCrossPlatformTrackingRunResult() loadstrike.LoadStrikeRunResult {
	var result loadstrike.LoadStrikeRunResult
	return result
}

func withCrossPlatformTrackingScenarioStats() loadstrike.LoadStrikeScenarioStats {
	var stats loadstrike.LoadStrikeScenarioStats
	return stats
}

func withCrossPlatformTrackingWriteTempConfigFiles() withCrossPlatformTrackingTempConfigPaths {
	return withCrossPlatformTrackingTempConfigPaths{
		ConfigPath: "method-reference.loadstrike.config.json",
		InfraPath:  "method-reference.loadstrike.infra.json",
	}
}

        // Attach the correlation configuration used for source and destination tracking.
func (reference WithCrossPlatformTrackingMethodReference) AttachDelegateTrackingExample() any {
    return withCrossPlatformTrackingBaselineScenario().WithCrossPlatformTracking(withCrossPlatformTrackingTrackingConfiguration())
}

// Start from an empty scenario and attach cross-platform tracking.
func (reference WithCrossPlatformTrackingMethodReference) AttachTrackingToEmptyScenarioExample() any {
    return loadstrike.Empty("orders.empty").WithCrossPlatformTracking(withCrossPlatformTrackingTrackingConfiguration())
}
