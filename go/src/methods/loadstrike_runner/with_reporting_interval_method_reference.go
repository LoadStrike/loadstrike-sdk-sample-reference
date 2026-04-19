package loadstrike_runner

        import (
        	"time"

        	loadstrike "loadstrike.com/sdk/go"
        )

        const withReportingIntervalRunnerKey = "runner_dummy_orders_reference"

type WithReportingIntervalMethodReference struct{}

type withReportingIntervalTempConfigPaths struct {
	ConfigPath string
	InfraPath  string
}

type withReportingIntervalOrdersReportingSink struct{}

func newWithReportingIntervalOrdersReportingSink() withReportingIntervalOrdersReportingSink {
	return withReportingIntervalOrdersReportingSink{}
}
func (withReportingIntervalOrdersReportingSink) SinkName() string {
	return "orders-sample-sink"
}
func (withReportingIntervalOrdersReportingSink) Init(loadstrike.LoadStrikeBaseContext, loadstrike.IConfiguration) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withReportingIntervalOrdersReportingSink) Start(loadstrike.LoadStrikeSessionStartInfo) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withReportingIntervalOrdersReportingSink) SaveRealtimeStats([]loadstrike.LoadStrikeScenarioStats) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withReportingIntervalOrdersReportingSink) SaveRealtimeMetrics(loadstrike.LoadStrikeMetricStats) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withReportingIntervalOrdersReportingSink) SaveRunResult(loadstrike.LoadStrikeRunResult) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withReportingIntervalOrdersReportingSink) Stop() loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withReportingIntervalOrdersReportingSink) Dispose() {}

type withReportingIntervalOrdersRuntimePolicy struct{}

func newWithReportingIntervalOrdersRuntimePolicy() withReportingIntervalOrdersRuntimePolicy {
	return withReportingIntervalOrdersRuntimePolicy{}
}
func (withReportingIntervalOrdersRuntimePolicy) ShouldRunScenario(string) loadstrike.LoadStrikeBoolTask {
	return loadstrike.TaskFromBool(true)
}
func (withReportingIntervalOrdersRuntimePolicy) BeforeScenario(string) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withReportingIntervalOrdersRuntimePolicy) AfterScenario(string, loadstrike.LoadStrikeScenarioRuntime) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withReportingIntervalOrdersRuntimePolicy) BeforeStep(string, string) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withReportingIntervalOrdersRuntimePolicy) AfterStep(string, string, loadstrike.LoadStrikeReply) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}

type withReportingIntervalOrdersWorkerPlugin struct{}

func newWithReportingIntervalOrdersWorkerPlugin() withReportingIntervalOrdersWorkerPlugin {
	return withReportingIntervalOrdersWorkerPlugin{}
}
func (withReportingIntervalOrdersWorkerPlugin) PluginName() string {
	return "orders-sample-plugin"
}
func (withReportingIntervalOrdersWorkerPlugin) Init(loadstrike.LoadStrikeBaseContext, loadstrike.IConfiguration) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withReportingIntervalOrdersWorkerPlugin) Start(loadstrike.LoadStrikeSessionStartInfo) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withReportingIntervalOrdersWorkerPlugin) GetData(loadstrike.LoadStrikeRunResult) loadstrike.LoadStrikeValueTask[loadstrike.LoadStrikePluginData] {
	return loadstrike.TaskFromResult(loadstrike.CreatePluginData("orders-sample-plugin"))
}
func (withReportingIntervalOrdersWorkerPlugin) Stop() loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withReportingIntervalOrdersWorkerPlugin) Dispose() loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}

func withReportingIntervalPerformOrderGetReply() loadstrike.LoadStrikeReply {
	return loadstrike.LoadStrikeResponse.Ok("200", int64(128), "ok", loadstrike.TimeSpan(3*time.Millisecond))
}

func withReportingIntervalExecuteOrderGet(context loadstrike.LoadStrikeScenarioContext) loadstrike.LoadStrikeReply {
	return loadstrike.LoadStrikeStep.Run("get-order", context, func(loadstrike.LoadStrikeScenarioContext) loadstrike.LoadStrikeReply {
		return withReportingIntervalPerformOrderGetReply()
	})
}

func withReportingIntervalExecuteOrderGetAsync(context loadstrike.LoadStrikeScenarioContext) loadstrike.LoadStrikeValueTask[loadstrike.LoadStrikeReply] {
	return loadstrike.TaskFromResult(withReportingIntervalExecuteOrderGet(context))
}

func withReportingIntervalBaselineScenario(names ...string) loadstrike.LoadStrikeScenario {
	name := "orders.get-by-id"
	if len(names) > 0 && names[0] != "" {
		name = names[0]
	}
	return loadstrike.CreateScenario(name, withReportingIntervalExecuteOrderGet).
		WithLoadSimulations(loadstrike.LoadStrikeSimulation.IterationsForConstant(1, 1)).
		WithoutWarmUp()
}

func withReportingIntervalBaseRunner() loadstrike.LoadStrikeRunner {
	return loadstrike.Create().
		AddScenario(withReportingIntervalBaselineScenario()).
		WithRunnerKey(withReportingIntervalRunnerKey).
		WithTestSuite("orders-reference").
		WithTestName("orders-get-by-id").
		WithoutReports()
}

func withReportingIntervalBaseContext() loadstrike.LoadStrikeContext {
	return withReportingIntervalBaseRunner().BuildContext()
}

func withReportingIntervalHttpSource() *loadstrike.EndpointSpec {
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

func withReportingIntervalHttpDestination() *loadstrike.EndpointSpec {
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

func withReportingIntervalTrackingConfiguration() *loadstrike.TrackingConfigurationSpec {
	return &loadstrike.TrackingConfigurationSpec{
		Source:                      withReportingIntervalHttpSource(),
		Destination:                 withReportingIntervalHttpDestination(),
		RunMode:                     "GenerateAndCorrelate",
		CorrelationTimeoutSeconds:   30,
		TimeoutSweepIntervalSeconds: 1,
		TimeoutBatchSize:            200,
		TimeoutCountsAsFailure:      true,
		MetricPrefix:                "orders_tracking",
		ExecuteOriginalScenarioRun:  false,
	}
}

func withReportingIntervalTrackedContext() loadstrike.LoadStrikeContext {
	return loadstrike.Create().
		AddScenario(withReportingIntervalBaselineScenario("orders.tracked").WithCrossPlatformTracking(withReportingIntervalTrackingConfiguration())).
		WithRunnerKey(withReportingIntervalRunnerKey).
		WithoutReports().
		BuildContext()
}

func withReportingIntervalBuildTrackingPayload() loadstrike.TrackingPayload {
	builder := loadstrike.TrackingPayloadBuilder{
		Headers:            map[string]string{"X-Correlation-Id": "ord-1"},
		ContentType:        "application/json",
		MessagePayloadType: "json",
	}
	builder.SetBody(map[string]any{"trackingId": "ord-1"})
	return builder.Build()
}

func withReportingIntervalRunResult() loadstrike.LoadStrikeRunResult {
	var result loadstrike.LoadStrikeRunResult
	return result
}

func withReportingIntervalScenarioStats() loadstrike.LoadStrikeScenarioStats {
	var stats loadstrike.LoadStrikeScenarioStats
	return stats
}

func withReportingIntervalWriteTempConfigFiles() withReportingIntervalTempConfigPaths {
	return withReportingIntervalTempConfigPaths{
		ConfigPath: "method-reference.loadstrike.config.json",
		InfraPath:  "method-reference.loadstrike.infra.json",
	}
}

        // Set a faster realtime reporting cadence.
func (reference WithReportingIntervalMethodReference) UseFastIntervalExample() any {
    return withReportingIntervalBaseRunner().WithReportingInterval(loadstrike.DurationFromSeconds(1))
}

// Set a slower reporting cadence for noisier runs.
func (reference WithReportingIntervalMethodReference) UseSlowerIntervalExample() any {
    return withReportingIntervalBaseRunner().WithReportingInterval(loadstrike.DurationFromSeconds(5))
}
