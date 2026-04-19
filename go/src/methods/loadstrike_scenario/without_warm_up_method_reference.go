package loadstrike_scenario

        import (
        	"time"

        	loadstrike "loadstrike.com/sdk/go"
        )

        const withoutWarmUpRunnerKey = "runner_dummy_orders_reference"

type WithoutWarmUpMethodReference struct{}

type withoutWarmUpTempConfigPaths struct {
	ConfigPath string
	InfraPath  string
}

type withoutWarmUpOrdersReportingSink struct{}

func newWithoutWarmUpOrdersReportingSink() withoutWarmUpOrdersReportingSink {
	return withoutWarmUpOrdersReportingSink{}
}
func (withoutWarmUpOrdersReportingSink) SinkName() string {
	return "orders-sample-sink"
}
func (withoutWarmUpOrdersReportingSink) Init(loadstrike.LoadStrikeBaseContext, loadstrike.IConfiguration) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withoutWarmUpOrdersReportingSink) Start(loadstrike.LoadStrikeSessionStartInfo) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withoutWarmUpOrdersReportingSink) SaveRealtimeStats([]loadstrike.LoadStrikeScenarioStats) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withoutWarmUpOrdersReportingSink) SaveRealtimeMetrics(loadstrike.LoadStrikeMetricStats) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withoutWarmUpOrdersReportingSink) SaveRunResult(loadstrike.LoadStrikeRunResult) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withoutWarmUpOrdersReportingSink) Stop() loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withoutWarmUpOrdersReportingSink) Dispose() {}

type withoutWarmUpOrdersRuntimePolicy struct{}

func newWithoutWarmUpOrdersRuntimePolicy() withoutWarmUpOrdersRuntimePolicy {
	return withoutWarmUpOrdersRuntimePolicy{}
}
func (withoutWarmUpOrdersRuntimePolicy) ShouldRunScenario(string) loadstrike.LoadStrikeBoolTask {
	return loadstrike.TaskFromBool(true)
}
func (withoutWarmUpOrdersRuntimePolicy) BeforeScenario(string) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withoutWarmUpOrdersRuntimePolicy) AfterScenario(string, loadstrike.LoadStrikeScenarioRuntime) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withoutWarmUpOrdersRuntimePolicy) BeforeStep(string, string) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withoutWarmUpOrdersRuntimePolicy) AfterStep(string, string, loadstrike.LoadStrikeReply) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}

type withoutWarmUpOrdersWorkerPlugin struct{}

func newWithoutWarmUpOrdersWorkerPlugin() withoutWarmUpOrdersWorkerPlugin {
	return withoutWarmUpOrdersWorkerPlugin{}
}
func (withoutWarmUpOrdersWorkerPlugin) PluginName() string {
	return "orders-sample-plugin"
}
func (withoutWarmUpOrdersWorkerPlugin) Init(loadstrike.LoadStrikeBaseContext, loadstrike.IConfiguration) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withoutWarmUpOrdersWorkerPlugin) Start(loadstrike.LoadStrikeSessionStartInfo) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withoutWarmUpOrdersWorkerPlugin) GetData(loadstrike.LoadStrikeRunResult) loadstrike.LoadStrikeValueTask[loadstrike.LoadStrikePluginData] {
	return loadstrike.TaskFromResult(loadstrike.CreatePluginData("orders-sample-plugin"))
}
func (withoutWarmUpOrdersWorkerPlugin) Stop() loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withoutWarmUpOrdersWorkerPlugin) Dispose() loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}

func withoutWarmUpPerformOrderGetReply() loadstrike.LoadStrikeReply {
	return loadstrike.LoadStrikeResponse.Ok("200", int64(128), "ok", loadstrike.TimeSpan(3*time.Millisecond))
}

func withoutWarmUpExecuteOrderGet(context loadstrike.LoadStrikeScenarioContext) loadstrike.LoadStrikeReply {
	return loadstrike.LoadStrikeStep.Run("get-order", context, func(loadstrike.LoadStrikeScenarioContext) loadstrike.LoadStrikeReply {
		return withoutWarmUpPerformOrderGetReply()
	})
}

func withoutWarmUpExecuteOrderGetAsync(context loadstrike.LoadStrikeScenarioContext) loadstrike.LoadStrikeValueTask[loadstrike.LoadStrikeReply] {
	return loadstrike.TaskFromResult(withoutWarmUpExecuteOrderGet(context))
}

func withoutWarmUpBaselineScenario(names ...string) loadstrike.LoadStrikeScenario {
	name := "orders.get-by-id"
	if len(names) > 0 && names[0] != "" {
		name = names[0]
	}
	return loadstrike.CreateScenario(name, withoutWarmUpExecuteOrderGet).
		WithLoadSimulations(loadstrike.LoadStrikeSimulation.IterationsForConstant(1, 1)).
		WithoutWarmUp()
}

func withoutWarmUpBaseRunner() loadstrike.LoadStrikeRunner {
	return loadstrike.Create().
		AddScenario(withoutWarmUpBaselineScenario()).
		WithRunnerKey(withoutWarmUpRunnerKey).
		WithTestSuite("orders-reference").
		WithTestName("orders-get-by-id").
		WithoutReports()
}

func withoutWarmUpBaseContext() loadstrike.LoadStrikeContext {
	return withoutWarmUpBaseRunner().BuildContext()
}

func withoutWarmUpHttpSource() *loadstrike.EndpointSpec {
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

func withoutWarmUpHttpDestination() *loadstrike.EndpointSpec {
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

func withoutWarmUpTrackingConfiguration() *loadstrike.TrackingConfigurationSpec {
	return &loadstrike.TrackingConfigurationSpec{
		Source:                      withoutWarmUpHttpSource(),
		Destination:                 withoutWarmUpHttpDestination(),
		RunMode:                     "GenerateAndCorrelate",
		CorrelationTimeoutSeconds:   30,
		TimeoutSweepIntervalSeconds: 1,
		TimeoutBatchSize:            200,
		TimeoutCountsAsFailure:      true,
		MetricPrefix:                "orders_tracking",
		ExecuteOriginalScenarioRun:  false,
	}
}

func withoutWarmUpTrackedContext() loadstrike.LoadStrikeContext {
	return loadstrike.Create().
		AddScenario(withoutWarmUpBaselineScenario("orders.tracked").WithCrossPlatformTracking(withoutWarmUpTrackingConfiguration())).
		WithRunnerKey(withoutWarmUpRunnerKey).
		WithoutReports().
		BuildContext()
}

func withoutWarmUpBuildTrackingPayload() loadstrike.TrackingPayload {
	builder := loadstrike.TrackingPayloadBuilder{
		Headers:            map[string]string{"X-Correlation-Id": "ord-1"},
		ContentType:        "application/json",
		MessagePayloadType: "json",
	}
	builder.SetBody(map[string]any{"trackingId": "ord-1"})
	return builder.Build()
}

func withoutWarmUpRunResult() loadstrike.LoadStrikeRunResult {
	var result loadstrike.LoadStrikeRunResult
	return result
}

func withoutWarmUpScenarioStats() loadstrike.LoadStrikeScenarioStats {
	var stats loadstrike.LoadStrikeScenarioStats
	return stats
}

func withoutWarmUpWriteTempConfigFiles() withoutWarmUpTempConfigPaths {
	return withoutWarmUpTempConfigPaths{
		ConfigPath: "method-reference.loadstrike.config.json",
		InfraPath:  "method-reference.loadstrike.infra.json",
	}
}

        // Disable warm-up for the GET scenario.
func (reference WithoutWarmUpMethodReference) DisableWarmUpExample() any {
    return withoutWarmUpBaselineScenario().WithoutWarmUp()
}

// Disable warm-up after attaching a simulation.
func (reference WithoutWarmUpMethodReference) DisableWarmUpAfterSimulationExample() any {
    return withoutWarmUpBaselineScenario().WithLoadSimulations(loadstrike.LoadStrikeSimulation.IterationsForConstant(1, 3)).WithoutWarmUp()
}
