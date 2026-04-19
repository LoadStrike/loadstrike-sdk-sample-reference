package loadstrike_scenario

        import (
        	"time"

        	loadstrike "loadstrike.com/sdk/go"
        )

        const withMaxFailCountRunnerKey = "runner_dummy_orders_reference"

type WithMaxFailCountMethodReference struct{}

type withMaxFailCountTempConfigPaths struct {
	ConfigPath string
	InfraPath  string
}

type withMaxFailCountOrdersReportingSink struct{}

func newWithMaxFailCountOrdersReportingSink() withMaxFailCountOrdersReportingSink {
	return withMaxFailCountOrdersReportingSink{}
}
func (withMaxFailCountOrdersReportingSink) SinkName() string {
	return "orders-sample-sink"
}
func (withMaxFailCountOrdersReportingSink) Init(loadstrike.LoadStrikeBaseContext, loadstrike.IConfiguration) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withMaxFailCountOrdersReportingSink) Start(loadstrike.LoadStrikeSessionStartInfo) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withMaxFailCountOrdersReportingSink) SaveRealtimeStats([]loadstrike.LoadStrikeScenarioStats) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withMaxFailCountOrdersReportingSink) SaveRealtimeMetrics(loadstrike.LoadStrikeMetricStats) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withMaxFailCountOrdersReportingSink) SaveRunResult(loadstrike.LoadStrikeRunResult) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withMaxFailCountOrdersReportingSink) Stop() loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withMaxFailCountOrdersReportingSink) Dispose() {}

type withMaxFailCountOrdersRuntimePolicy struct{}

func newWithMaxFailCountOrdersRuntimePolicy() withMaxFailCountOrdersRuntimePolicy {
	return withMaxFailCountOrdersRuntimePolicy{}
}
func (withMaxFailCountOrdersRuntimePolicy) ShouldRunScenario(string) loadstrike.LoadStrikeBoolTask {
	return loadstrike.TaskFromBool(true)
}
func (withMaxFailCountOrdersRuntimePolicy) BeforeScenario(string) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withMaxFailCountOrdersRuntimePolicy) AfterScenario(string, loadstrike.LoadStrikeScenarioRuntime) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withMaxFailCountOrdersRuntimePolicy) BeforeStep(string, string) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withMaxFailCountOrdersRuntimePolicy) AfterStep(string, string, loadstrike.LoadStrikeReply) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}

type withMaxFailCountOrdersWorkerPlugin struct{}

func newWithMaxFailCountOrdersWorkerPlugin() withMaxFailCountOrdersWorkerPlugin {
	return withMaxFailCountOrdersWorkerPlugin{}
}
func (withMaxFailCountOrdersWorkerPlugin) PluginName() string {
	return "orders-sample-plugin"
}
func (withMaxFailCountOrdersWorkerPlugin) Init(loadstrike.LoadStrikeBaseContext, loadstrike.IConfiguration) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withMaxFailCountOrdersWorkerPlugin) Start(loadstrike.LoadStrikeSessionStartInfo) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withMaxFailCountOrdersWorkerPlugin) GetData(loadstrike.LoadStrikeRunResult) loadstrike.LoadStrikeValueTask[loadstrike.LoadStrikePluginData] {
	return loadstrike.TaskFromResult(loadstrike.CreatePluginData("orders-sample-plugin"))
}
func (withMaxFailCountOrdersWorkerPlugin) Stop() loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withMaxFailCountOrdersWorkerPlugin) Dispose() loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}

func withMaxFailCountPerformOrderGetReply() loadstrike.LoadStrikeReply {
	return loadstrike.LoadStrikeResponse.Ok("200", int64(128), "ok", loadstrike.TimeSpan(3*time.Millisecond))
}

func withMaxFailCountExecuteOrderGet(context loadstrike.LoadStrikeScenarioContext) loadstrike.LoadStrikeReply {
	return loadstrike.LoadStrikeStep.Run("get-order", context, func(loadstrike.LoadStrikeScenarioContext) loadstrike.LoadStrikeReply {
		return withMaxFailCountPerformOrderGetReply()
	})
}

func withMaxFailCountExecuteOrderGetAsync(context loadstrike.LoadStrikeScenarioContext) loadstrike.LoadStrikeValueTask[loadstrike.LoadStrikeReply] {
	return loadstrike.TaskFromResult(withMaxFailCountExecuteOrderGet(context))
}

func withMaxFailCountBaselineScenario(names ...string) loadstrike.LoadStrikeScenario {
	name := "orders.get-by-id"
	if len(names) > 0 && names[0] != "" {
		name = names[0]
	}
	return loadstrike.CreateScenario(name, withMaxFailCountExecuteOrderGet).
		WithLoadSimulations(loadstrike.LoadStrikeSimulation.IterationsForConstant(1, 1)).
		WithoutWarmUp()
}

func withMaxFailCountBaseRunner() loadstrike.LoadStrikeRunner {
	return loadstrike.Create().
		AddScenario(withMaxFailCountBaselineScenario()).
		WithRunnerKey(withMaxFailCountRunnerKey).
		WithTestSuite("orders-reference").
		WithTestName("orders-get-by-id").
		WithoutReports()
}

func withMaxFailCountBaseContext() loadstrike.LoadStrikeContext {
	return withMaxFailCountBaseRunner().BuildContext()
}

func withMaxFailCountHttpSource() *loadstrike.EndpointSpec {
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

func withMaxFailCountHttpDestination() *loadstrike.EndpointSpec {
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

func withMaxFailCountTrackingConfiguration() *loadstrike.TrackingConfigurationSpec {
	return &loadstrike.TrackingConfigurationSpec{
		Source:                      withMaxFailCountHttpSource(),
		Destination:                 withMaxFailCountHttpDestination(),
		RunMode:                     "GenerateAndCorrelate",
		CorrelationTimeoutSeconds:   30,
		TimeoutSweepIntervalSeconds: 1,
		TimeoutBatchSize:            200,
		TimeoutCountsAsFailure:      true,
		MetricPrefix:                "orders_tracking",
		ExecuteOriginalScenarioRun:  false,
	}
}

func withMaxFailCountTrackedContext() loadstrike.LoadStrikeContext {
	return loadstrike.Create().
		AddScenario(withMaxFailCountBaselineScenario("orders.tracked").WithCrossPlatformTracking(withMaxFailCountTrackingConfiguration())).
		WithRunnerKey(withMaxFailCountRunnerKey).
		WithoutReports().
		BuildContext()
}

func withMaxFailCountBuildTrackingPayload() loadstrike.TrackingPayload {
	builder := loadstrike.TrackingPayloadBuilder{
		Headers:            map[string]string{"X-Correlation-Id": "ord-1"},
		ContentType:        "application/json",
		MessagePayloadType: "json",
	}
	builder.SetBody(map[string]any{"trackingId": "ord-1"})
	return builder.Build()
}

func withMaxFailCountRunResult() loadstrike.LoadStrikeRunResult {
	var result loadstrike.LoadStrikeRunResult
	return result
}

func withMaxFailCountScenarioStats() loadstrike.LoadStrikeScenarioStats {
	var stats loadstrike.LoadStrikeScenarioStats
	return stats
}

func withMaxFailCountWriteTempConfigFiles() withMaxFailCountTempConfigPaths {
	return withMaxFailCountTempConfigPaths{
		ConfigPath: "method-reference.loadstrike.config.json",
		InfraPath:  "method-reference.loadstrike.infra.json",
	}
}

        // Stop the scenario after the first failure.
func (reference WithMaxFailCountMethodReference) StopAfterOneFailureExample() any {
    return withMaxFailCountBaselineScenario().WithMaxFailCount(1)
}

// Allow a few failures before stopping the scenario.
func (reference WithMaxFailCountMethodReference) StopAfterThreeFailuresExample() any {
    return withMaxFailCountBaselineScenario().WithMaxFailCount(3)
}
