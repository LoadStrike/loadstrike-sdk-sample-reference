package loadstrike_scenario_context

        import (
        	"time"

        	loadstrike "loadstrike.com/sdk/go"
        )

        const stopCurrentTestRunnerKey = "runner_dummy_orders_reference"

type StopCurrentTestMethodReference struct{}

type stopCurrentTestTempConfigPaths struct {
	ConfigPath string
	InfraPath  string
}

type stopCurrentTestOrdersReportingSink struct{}

func newStopCurrentTestOrdersReportingSink() stopCurrentTestOrdersReportingSink {
	return stopCurrentTestOrdersReportingSink{}
}
func (stopCurrentTestOrdersReportingSink) SinkName() string {
	return "orders-sample-sink"
}
func (stopCurrentTestOrdersReportingSink) Init(loadstrike.LoadStrikeBaseContext, loadstrike.IConfiguration) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (stopCurrentTestOrdersReportingSink) Start(loadstrike.LoadStrikeSessionStartInfo) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (stopCurrentTestOrdersReportingSink) SaveRealtimeStats([]loadstrike.LoadStrikeScenarioStats) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (stopCurrentTestOrdersReportingSink) SaveRealtimeMetrics(loadstrike.LoadStrikeMetricStats) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (stopCurrentTestOrdersReportingSink) SaveRunResult(loadstrike.LoadStrikeRunResult) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (stopCurrentTestOrdersReportingSink) Stop() loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (stopCurrentTestOrdersReportingSink) Dispose() {}

type stopCurrentTestOrdersRuntimePolicy struct{}

func newStopCurrentTestOrdersRuntimePolicy() stopCurrentTestOrdersRuntimePolicy {
	return stopCurrentTestOrdersRuntimePolicy{}
}
func (stopCurrentTestOrdersRuntimePolicy) ShouldRunScenario(string) loadstrike.LoadStrikeBoolTask {
	return loadstrike.TaskFromBool(true)
}
func (stopCurrentTestOrdersRuntimePolicy) BeforeScenario(string) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (stopCurrentTestOrdersRuntimePolicy) AfterScenario(string, loadstrike.LoadStrikeScenarioRuntime) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (stopCurrentTestOrdersRuntimePolicy) BeforeStep(string, string) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (stopCurrentTestOrdersRuntimePolicy) AfterStep(string, string, loadstrike.LoadStrikeReply) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}

type stopCurrentTestOrdersWorkerPlugin struct{}

func newStopCurrentTestOrdersWorkerPlugin() stopCurrentTestOrdersWorkerPlugin {
	return stopCurrentTestOrdersWorkerPlugin{}
}
func (stopCurrentTestOrdersWorkerPlugin) PluginName() string {
	return "orders-sample-plugin"
}
func (stopCurrentTestOrdersWorkerPlugin) Init(loadstrike.LoadStrikeBaseContext, loadstrike.IConfiguration) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (stopCurrentTestOrdersWorkerPlugin) Start(loadstrike.LoadStrikeSessionStartInfo) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (stopCurrentTestOrdersWorkerPlugin) GetData(loadstrike.LoadStrikeRunResult) loadstrike.LoadStrikeValueTask[loadstrike.LoadStrikePluginData] {
	return loadstrike.TaskFromResult(loadstrike.CreatePluginData("orders-sample-plugin"))
}
func (stopCurrentTestOrdersWorkerPlugin) Stop() loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (stopCurrentTestOrdersWorkerPlugin) Dispose() loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}

func stopCurrentTestPerformOrderGetReply() loadstrike.LoadStrikeReply {
	return loadstrike.LoadStrikeResponse.Ok("200", int64(128), "ok", loadstrike.TimeSpan(3*time.Millisecond))
}

func stopCurrentTestExecuteOrderGet(context loadstrike.LoadStrikeScenarioContext) loadstrike.LoadStrikeReply {
	return loadstrike.LoadStrikeStep.Run("get-order", context, func(loadstrike.LoadStrikeScenarioContext) loadstrike.LoadStrikeReply {
		return stopCurrentTestPerformOrderGetReply()
	})
}

func stopCurrentTestExecuteOrderGetAsync(context loadstrike.LoadStrikeScenarioContext) loadstrike.LoadStrikeValueTask[loadstrike.LoadStrikeReply] {
	return loadstrike.TaskFromResult(stopCurrentTestExecuteOrderGet(context))
}

func stopCurrentTestBaselineScenario(names ...string) loadstrike.LoadStrikeScenario {
	name := "orders.get-by-id"
	if len(names) > 0 && names[0] != "" {
		name = names[0]
	}
	return loadstrike.CreateScenario(name, stopCurrentTestExecuteOrderGet).
		WithLoadSimulations(loadstrike.LoadStrikeSimulation.IterationsForConstant(1, 1)).
		WithoutWarmUp()
}

func stopCurrentTestBaseRunner() loadstrike.LoadStrikeRunner {
	return loadstrike.Create().
		AddScenario(stopCurrentTestBaselineScenario()).
		WithRunnerKey(stopCurrentTestRunnerKey).
		WithTestSuite("orders-reference").
		WithTestName("orders-get-by-id").
		WithoutReports()
}

func stopCurrentTestBaseContext() loadstrike.LoadStrikeContext {
	return stopCurrentTestBaseRunner().BuildContext()
}

func stopCurrentTestHttpSource() *loadstrike.EndpointSpec {
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

func stopCurrentTestHttpDestination() *loadstrike.EndpointSpec {
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

func stopCurrentTestTrackingConfiguration() *loadstrike.TrackingConfigurationSpec {
	return &loadstrike.TrackingConfigurationSpec{
		Source:                      stopCurrentTestHttpSource(),
		Destination:                 stopCurrentTestHttpDestination(),
		RunMode:                     "GenerateAndCorrelate",
		CorrelationTimeoutSeconds:   30,
		TimeoutSweepIntervalSeconds: 1,
		TimeoutBatchSize:            200,
		TimeoutCountsAsFailure:      true,
		MetricPrefix:                "orders_tracking",
		ExecuteOriginalScenarioRun:  false,
	}
}

func stopCurrentTestTrackedContext() loadstrike.LoadStrikeContext {
	return loadstrike.Create().
		AddScenario(stopCurrentTestBaselineScenario("orders.tracked").WithCrossPlatformTracking(stopCurrentTestTrackingConfiguration())).
		WithRunnerKey(stopCurrentTestRunnerKey).
		WithoutReports().
		BuildContext()
}

func stopCurrentTestBuildTrackingPayload() loadstrike.TrackingPayload {
	builder := loadstrike.TrackingPayloadBuilder{
		Headers:            map[string]string{"X-Correlation-Id": "ord-1"},
		ContentType:        "application/json",
		MessagePayloadType: "json",
	}
	builder.SetBody(map[string]any{"trackingId": "ord-1"})
	return builder.Build()
}

func stopCurrentTestRunResult() loadstrike.LoadStrikeRunResult {
	var result loadstrike.LoadStrikeRunResult
	return result
}

func stopCurrentTestScenarioStats() loadstrike.LoadStrikeScenarioStats {
	var stats loadstrike.LoadStrikeScenarioStats
	return stats
}

func stopCurrentTestWriteTempConfigFiles() stopCurrentTestTempConfigPaths {
	return stopCurrentTestTempConfigPaths{
		ConfigPath: "method-reference.loadstrike.config.json",
		InfraPath:  "method-reference.loadstrike.infra.json",
	}
}

        // Call the public context helper directly from the scenario context surface.
func (reference StopCurrentTestMethodReference) UseContextMethodExample() any {
    return loadstrike.CreateScenario("orders.stop-test", func(context loadstrike.LoadStrikeScenarioContext) loadstrike.LoadStrikeReply { context.StopCurrentTest("stop-now"); return loadstrike.OK() })
}

// Show the same helper in the baseline GET-step flow.
func (reference StopCurrentTestMethodReference) UseContextMethodInStepExample() any {
    return loadstrike.CreateScenario("orders.stop-test", func(context loadstrike.LoadStrikeScenarioContext) loadstrike.LoadStrikeReply { context.StopCurrentTest("stop-now"); return loadstrike.OK() })
}
