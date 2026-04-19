package loadstrike_scenario_context

        import (
        	"time"

        	loadstrike "loadstrike.com/sdk/go"
        )

        const stopScenarioRunnerKey = "runner_dummy_orders_reference"

type StopScenarioMethodReference struct{}

type stopScenarioTempConfigPaths struct {
	ConfigPath string
	InfraPath  string
}

type stopScenarioOrdersReportingSink struct{}

func newStopScenarioOrdersReportingSink() stopScenarioOrdersReportingSink {
	return stopScenarioOrdersReportingSink{}
}
func (stopScenarioOrdersReportingSink) SinkName() string {
	return "orders-sample-sink"
}
func (stopScenarioOrdersReportingSink) Init(loadstrike.LoadStrikeBaseContext, loadstrike.IConfiguration) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (stopScenarioOrdersReportingSink) Start(loadstrike.LoadStrikeSessionStartInfo) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (stopScenarioOrdersReportingSink) SaveRealtimeStats([]loadstrike.LoadStrikeScenarioStats) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (stopScenarioOrdersReportingSink) SaveRealtimeMetrics(loadstrike.LoadStrikeMetricStats) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (stopScenarioOrdersReportingSink) SaveRunResult(loadstrike.LoadStrikeRunResult) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (stopScenarioOrdersReportingSink) Stop() loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (stopScenarioOrdersReportingSink) Dispose() {}

type stopScenarioOrdersRuntimePolicy struct{}

func newStopScenarioOrdersRuntimePolicy() stopScenarioOrdersRuntimePolicy {
	return stopScenarioOrdersRuntimePolicy{}
}
func (stopScenarioOrdersRuntimePolicy) ShouldRunScenario(string) loadstrike.LoadStrikeBoolTask {
	return loadstrike.TaskFromBool(true)
}
func (stopScenarioOrdersRuntimePolicy) BeforeScenario(string) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (stopScenarioOrdersRuntimePolicy) AfterScenario(string, loadstrike.LoadStrikeScenarioRuntime) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (stopScenarioOrdersRuntimePolicy) BeforeStep(string, string) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (stopScenarioOrdersRuntimePolicy) AfterStep(string, string, loadstrike.LoadStrikeReply) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}

type stopScenarioOrdersWorkerPlugin struct{}

func newStopScenarioOrdersWorkerPlugin() stopScenarioOrdersWorkerPlugin {
	return stopScenarioOrdersWorkerPlugin{}
}
func (stopScenarioOrdersWorkerPlugin) PluginName() string {
	return "orders-sample-plugin"
}
func (stopScenarioOrdersWorkerPlugin) Init(loadstrike.LoadStrikeBaseContext, loadstrike.IConfiguration) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (stopScenarioOrdersWorkerPlugin) Start(loadstrike.LoadStrikeSessionStartInfo) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (stopScenarioOrdersWorkerPlugin) GetData(loadstrike.LoadStrikeRunResult) loadstrike.LoadStrikeValueTask[loadstrike.LoadStrikePluginData] {
	return loadstrike.TaskFromResult(loadstrike.CreatePluginData("orders-sample-plugin"))
}
func (stopScenarioOrdersWorkerPlugin) Stop() loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (stopScenarioOrdersWorkerPlugin) Dispose() loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}

func stopScenarioPerformOrderGetReply() loadstrike.LoadStrikeReply {
	return loadstrike.LoadStrikeResponse.Ok("200", int64(128), "ok", loadstrike.TimeSpan(3*time.Millisecond))
}

func stopScenarioExecuteOrderGet(context loadstrike.LoadStrikeScenarioContext) loadstrike.LoadStrikeReply {
	return loadstrike.LoadStrikeStep.Run("get-order", context, func(loadstrike.LoadStrikeScenarioContext) loadstrike.LoadStrikeReply {
		return stopScenarioPerformOrderGetReply()
	})
}

func stopScenarioExecuteOrderGetAsync(context loadstrike.LoadStrikeScenarioContext) loadstrike.LoadStrikeValueTask[loadstrike.LoadStrikeReply] {
	return loadstrike.TaskFromResult(stopScenarioExecuteOrderGet(context))
}

func stopScenarioBaselineScenario(names ...string) loadstrike.LoadStrikeScenario {
	name := "orders.get-by-id"
	if len(names) > 0 && names[0] != "" {
		name = names[0]
	}
	return loadstrike.CreateScenario(name, stopScenarioExecuteOrderGet).
		WithLoadSimulations(loadstrike.LoadStrikeSimulation.IterationsForConstant(1, 1)).
		WithoutWarmUp()
}

func stopScenarioBaseRunner() loadstrike.LoadStrikeRunner {
	return loadstrike.Create().
		AddScenario(stopScenarioBaselineScenario()).
		WithRunnerKey(stopScenarioRunnerKey).
		WithTestSuite("orders-reference").
		WithTestName("orders-get-by-id").
		WithoutReports()
}

func stopScenarioBaseContext() loadstrike.LoadStrikeContext {
	return stopScenarioBaseRunner().BuildContext()
}

func stopScenarioHttpSource() *loadstrike.EndpointSpec {
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

func stopScenarioHttpDestination() *loadstrike.EndpointSpec {
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

func stopScenarioTrackingConfiguration() *loadstrike.TrackingConfigurationSpec {
	return &loadstrike.TrackingConfigurationSpec{
		Source:                      stopScenarioHttpSource(),
		Destination:                 stopScenarioHttpDestination(),
		RunMode:                     "GenerateAndCorrelate",
		CorrelationTimeoutSeconds:   30,
		TimeoutSweepIntervalSeconds: 1,
		TimeoutBatchSize:            200,
		TimeoutCountsAsFailure:      true,
		MetricPrefix:                "orders_tracking",
		ExecuteOriginalScenarioRun:  false,
	}
}

func stopScenarioTrackedContext() loadstrike.LoadStrikeContext {
	return loadstrike.Create().
		AddScenario(stopScenarioBaselineScenario("orders.tracked").WithCrossPlatformTracking(stopScenarioTrackingConfiguration())).
		WithRunnerKey(stopScenarioRunnerKey).
		WithoutReports().
		BuildContext()
}

func stopScenarioBuildTrackingPayload() loadstrike.TrackingPayload {
	builder := loadstrike.TrackingPayloadBuilder{
		Headers:            map[string]string{"X-Correlation-Id": "ord-1"},
		ContentType:        "application/json",
		MessagePayloadType: "json",
	}
	builder.SetBody(map[string]any{"trackingId": "ord-1"})
	return builder.Build()
}

func stopScenarioRunResult() loadstrike.LoadStrikeRunResult {
	var result loadstrike.LoadStrikeRunResult
	return result
}

func stopScenarioScenarioStats() loadstrike.LoadStrikeScenarioStats {
	var stats loadstrike.LoadStrikeScenarioStats
	return stats
}

func stopScenarioWriteTempConfigFiles() stopScenarioTempConfigPaths {
	return stopScenarioTempConfigPaths{
		ConfigPath: "method-reference.loadstrike.config.json",
		InfraPath:  "method-reference.loadstrike.infra.json",
	}
}

        // Call the public context helper directly from the scenario context surface.
func (reference StopScenarioMethodReference) UseContextMethodExample() any {
    return loadstrike.CreateScenario("orders.stop", func(context loadstrike.LoadStrikeScenarioContext) loadstrike.LoadStrikeReply { context.StopScenario("orders.stop", "stop-now"); return loadstrike.OK() })
}

// Show the same helper in the baseline GET-step flow.
func (reference StopScenarioMethodReference) UseContextMethodInStepExample() any {
    return loadstrike.CreateScenario("orders.stop", func(context loadstrike.LoadStrikeScenarioContext) loadstrike.LoadStrikeReply { context.StopScenario("orders.stop", "stop-now"); return loadstrike.OK() })
}
