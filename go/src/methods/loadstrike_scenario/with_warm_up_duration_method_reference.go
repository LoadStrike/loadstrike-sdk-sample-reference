package loadstrike_scenario

        import (
        	"time"

        	loadstrike "loadstrike.com/sdk/go"
        )

        const withWarmUpDurationRunnerKey = "runner_dummy_orders_reference"

type WithWarmUpDurationMethodReference struct{}

type withWarmUpDurationTempConfigPaths struct {
	ConfigPath string
	InfraPath  string
}

type withWarmUpDurationOrdersReportingSink struct{}

func newWithWarmUpDurationOrdersReportingSink() withWarmUpDurationOrdersReportingSink {
	return withWarmUpDurationOrdersReportingSink{}
}
func (withWarmUpDurationOrdersReportingSink) SinkName() string {
	return "orders-sample-sink"
}
func (withWarmUpDurationOrdersReportingSink) Init(loadstrike.LoadStrikeBaseContext, loadstrike.IConfiguration) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withWarmUpDurationOrdersReportingSink) Start(loadstrike.LoadStrikeSessionStartInfo) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withWarmUpDurationOrdersReportingSink) SaveRealtimeStats([]loadstrike.LoadStrikeScenarioStats) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withWarmUpDurationOrdersReportingSink) SaveRealtimeMetrics(loadstrike.LoadStrikeMetricStats) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withWarmUpDurationOrdersReportingSink) SaveRunResult(loadstrike.LoadStrikeRunResult) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withWarmUpDurationOrdersReportingSink) Stop() loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withWarmUpDurationOrdersReportingSink) Dispose() {}

type withWarmUpDurationOrdersRuntimePolicy struct{}

func newWithWarmUpDurationOrdersRuntimePolicy() withWarmUpDurationOrdersRuntimePolicy {
	return withWarmUpDurationOrdersRuntimePolicy{}
}
func (withWarmUpDurationOrdersRuntimePolicy) ShouldRunScenario(string) loadstrike.LoadStrikeBoolTask {
	return loadstrike.TaskFromBool(true)
}
func (withWarmUpDurationOrdersRuntimePolicy) BeforeScenario(string) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withWarmUpDurationOrdersRuntimePolicy) AfterScenario(string, loadstrike.LoadStrikeScenarioRuntime) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withWarmUpDurationOrdersRuntimePolicy) BeforeStep(string, string) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withWarmUpDurationOrdersRuntimePolicy) AfterStep(string, string, loadstrike.LoadStrikeReply) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}

type withWarmUpDurationOrdersWorkerPlugin struct{}

func newWithWarmUpDurationOrdersWorkerPlugin() withWarmUpDurationOrdersWorkerPlugin {
	return withWarmUpDurationOrdersWorkerPlugin{}
}
func (withWarmUpDurationOrdersWorkerPlugin) PluginName() string {
	return "orders-sample-plugin"
}
func (withWarmUpDurationOrdersWorkerPlugin) Init(loadstrike.LoadStrikeBaseContext, loadstrike.IConfiguration) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withWarmUpDurationOrdersWorkerPlugin) Start(loadstrike.LoadStrikeSessionStartInfo) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withWarmUpDurationOrdersWorkerPlugin) GetData(loadstrike.LoadStrikeRunResult) loadstrike.LoadStrikeValueTask[loadstrike.LoadStrikePluginData] {
	return loadstrike.TaskFromResult(loadstrike.CreatePluginData("orders-sample-plugin"))
}
func (withWarmUpDurationOrdersWorkerPlugin) Stop() loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withWarmUpDurationOrdersWorkerPlugin) Dispose() loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}

func withWarmUpDurationPerformOrderGetReply() loadstrike.LoadStrikeReply {
	return loadstrike.LoadStrikeResponse.Ok("200", int64(128), "ok", loadstrike.TimeSpan(3*time.Millisecond))
}

func withWarmUpDurationExecuteOrderGet(context loadstrike.LoadStrikeScenarioContext) loadstrike.LoadStrikeReply {
	return loadstrike.LoadStrikeStep.Run("get-order", context, func(loadstrike.LoadStrikeScenarioContext) loadstrike.LoadStrikeReply {
		return withWarmUpDurationPerformOrderGetReply()
	})
}

func withWarmUpDurationExecuteOrderGetAsync(context loadstrike.LoadStrikeScenarioContext) loadstrike.LoadStrikeValueTask[loadstrike.LoadStrikeReply] {
	return loadstrike.TaskFromResult(withWarmUpDurationExecuteOrderGet(context))
}

func withWarmUpDurationBaselineScenario(names ...string) loadstrike.LoadStrikeScenario {
	name := "orders.get-by-id"
	if len(names) > 0 && names[0] != "" {
		name = names[0]
	}
	return loadstrike.CreateScenario(name, withWarmUpDurationExecuteOrderGet).
		WithLoadSimulations(loadstrike.LoadStrikeSimulation.IterationsForConstant(1, 1)).
		WithoutWarmUp()
}

func withWarmUpDurationBaseRunner() loadstrike.LoadStrikeRunner {
	return loadstrike.Create().
		AddScenario(withWarmUpDurationBaselineScenario()).
		WithRunnerKey(withWarmUpDurationRunnerKey).
		WithTestSuite("orders-reference").
		WithTestName("orders-get-by-id").
		WithoutReports()
}

func withWarmUpDurationBaseContext() loadstrike.LoadStrikeContext {
	return withWarmUpDurationBaseRunner().BuildContext()
}

func withWarmUpDurationHttpSource() *loadstrike.EndpointSpec {
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

func withWarmUpDurationHttpDestination() *loadstrike.EndpointSpec {
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

func withWarmUpDurationTrackingConfiguration() *loadstrike.TrackingConfigurationSpec {
	return &loadstrike.TrackingConfigurationSpec{
		Source:                      withWarmUpDurationHttpSource(),
		Destination:                 withWarmUpDurationHttpDestination(),
		RunMode:                     "GenerateAndCorrelate",
		CorrelationTimeoutSeconds:   30,
		TimeoutSweepIntervalSeconds: 1,
		TimeoutBatchSize:            200,
		TimeoutCountsAsFailure:      true,
		MetricPrefix:                "orders_tracking",
		ExecuteOriginalScenarioRun:  false,
	}
}

func withWarmUpDurationTrackedContext() loadstrike.LoadStrikeContext {
	return loadstrike.Create().
		AddScenario(withWarmUpDurationBaselineScenario("orders.tracked").WithCrossPlatformTracking(withWarmUpDurationTrackingConfiguration())).
		WithRunnerKey(withWarmUpDurationRunnerKey).
		WithoutReports().
		BuildContext()
}

func withWarmUpDurationBuildTrackingPayload() loadstrike.TrackingPayload {
	builder := loadstrike.TrackingPayloadBuilder{
		Headers:            map[string]string{"X-Correlation-Id": "ord-1"},
		ContentType:        "application/json",
		MessagePayloadType: "json",
	}
	builder.SetBody(map[string]any{"trackingId": "ord-1"})
	return builder.Build()
}

func withWarmUpDurationRunResult() loadstrike.LoadStrikeRunResult {
	var result loadstrike.LoadStrikeRunResult
	return result
}

func withWarmUpDurationScenarioStats() loadstrike.LoadStrikeScenarioStats {
	var stats loadstrike.LoadStrikeScenarioStats
	return stats
}

func withWarmUpDurationWriteTempConfigFiles() withWarmUpDurationTempConfigPaths {
	return withWarmUpDurationTempConfigPaths{
		ConfigPath: "method-reference.loadstrike.config.json",
		InfraPath:  "method-reference.loadstrike.infra.json",
	}
}

        // Set a short warm-up duration.
func (reference WithWarmUpDurationMethodReference) UseShortWarmUpExample() any {
    return withWarmUpDurationBaselineScenario().WithWarmUpDuration(2)
}

// Set a longer warm-up duration.
func (reference WithWarmUpDurationMethodReference) UseLongWarmUpExample() any {
    return withWarmUpDurationBaselineScenario().WithWarmUpDuration(10)
}
