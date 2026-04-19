package loadstrike_scenario

        import (
        	"time"

        	loadstrike "loadstrike.com/sdk/go"
        )

        const withLoadSimulationsRunnerKey = "runner_dummy_orders_reference"

type WithLoadSimulationsMethodReference struct{}

type withLoadSimulationsTempConfigPaths struct {
	ConfigPath string
	InfraPath  string
}

type withLoadSimulationsOrdersReportingSink struct{}

func newWithLoadSimulationsOrdersReportingSink() withLoadSimulationsOrdersReportingSink {
	return withLoadSimulationsOrdersReportingSink{}
}
func (withLoadSimulationsOrdersReportingSink) SinkName() string {
	return "orders-sample-sink"
}
func (withLoadSimulationsOrdersReportingSink) Init(loadstrike.LoadStrikeBaseContext, loadstrike.IConfiguration) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withLoadSimulationsOrdersReportingSink) Start(loadstrike.LoadStrikeSessionStartInfo) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withLoadSimulationsOrdersReportingSink) SaveRealtimeStats([]loadstrike.LoadStrikeScenarioStats) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withLoadSimulationsOrdersReportingSink) SaveRealtimeMetrics(loadstrike.LoadStrikeMetricStats) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withLoadSimulationsOrdersReportingSink) SaveRunResult(loadstrike.LoadStrikeRunResult) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withLoadSimulationsOrdersReportingSink) Stop() loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withLoadSimulationsOrdersReportingSink) Dispose() {}

type withLoadSimulationsOrdersRuntimePolicy struct{}

func newWithLoadSimulationsOrdersRuntimePolicy() withLoadSimulationsOrdersRuntimePolicy {
	return withLoadSimulationsOrdersRuntimePolicy{}
}
func (withLoadSimulationsOrdersRuntimePolicy) ShouldRunScenario(string) loadstrike.LoadStrikeBoolTask {
	return loadstrike.TaskFromBool(true)
}
func (withLoadSimulationsOrdersRuntimePolicy) BeforeScenario(string) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withLoadSimulationsOrdersRuntimePolicy) AfterScenario(string, loadstrike.LoadStrikeScenarioRuntime) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withLoadSimulationsOrdersRuntimePolicy) BeforeStep(string, string) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withLoadSimulationsOrdersRuntimePolicy) AfterStep(string, string, loadstrike.LoadStrikeReply) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}

type withLoadSimulationsOrdersWorkerPlugin struct{}

func newWithLoadSimulationsOrdersWorkerPlugin() withLoadSimulationsOrdersWorkerPlugin {
	return withLoadSimulationsOrdersWorkerPlugin{}
}
func (withLoadSimulationsOrdersWorkerPlugin) PluginName() string {
	return "orders-sample-plugin"
}
func (withLoadSimulationsOrdersWorkerPlugin) Init(loadstrike.LoadStrikeBaseContext, loadstrike.IConfiguration) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withLoadSimulationsOrdersWorkerPlugin) Start(loadstrike.LoadStrikeSessionStartInfo) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withLoadSimulationsOrdersWorkerPlugin) GetData(loadstrike.LoadStrikeRunResult) loadstrike.LoadStrikeValueTask[loadstrike.LoadStrikePluginData] {
	return loadstrike.TaskFromResult(loadstrike.CreatePluginData("orders-sample-plugin"))
}
func (withLoadSimulationsOrdersWorkerPlugin) Stop() loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withLoadSimulationsOrdersWorkerPlugin) Dispose() loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}

func withLoadSimulationsPerformOrderGetReply() loadstrike.LoadStrikeReply {
	return loadstrike.LoadStrikeResponse.Ok("200", int64(128), "ok", loadstrike.TimeSpan(3*time.Millisecond))
}

func withLoadSimulationsExecuteOrderGet(context loadstrike.LoadStrikeScenarioContext) loadstrike.LoadStrikeReply {
	return loadstrike.LoadStrikeStep.Run("get-order", context, func(loadstrike.LoadStrikeScenarioContext) loadstrike.LoadStrikeReply {
		return withLoadSimulationsPerformOrderGetReply()
	})
}

func withLoadSimulationsExecuteOrderGetAsync(context loadstrike.LoadStrikeScenarioContext) loadstrike.LoadStrikeValueTask[loadstrike.LoadStrikeReply] {
	return loadstrike.TaskFromResult(withLoadSimulationsExecuteOrderGet(context))
}

func withLoadSimulationsBaselineScenario(names ...string) loadstrike.LoadStrikeScenario {
	name := "orders.get-by-id"
	if len(names) > 0 && names[0] != "" {
		name = names[0]
	}
	return loadstrike.CreateScenario(name, withLoadSimulationsExecuteOrderGet).
		WithLoadSimulations(loadstrike.LoadStrikeSimulation.IterationsForConstant(1, 1)).
		WithoutWarmUp()
}

func withLoadSimulationsBaseRunner() loadstrike.LoadStrikeRunner {
	return loadstrike.Create().
		AddScenario(withLoadSimulationsBaselineScenario()).
		WithRunnerKey(withLoadSimulationsRunnerKey).
		WithTestSuite("orders-reference").
		WithTestName("orders-get-by-id").
		WithoutReports()
}

func withLoadSimulationsBaseContext() loadstrike.LoadStrikeContext {
	return withLoadSimulationsBaseRunner().BuildContext()
}

func withLoadSimulationsHttpSource() *loadstrike.EndpointSpec {
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

func withLoadSimulationsHttpDestination() *loadstrike.EndpointSpec {
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

func withLoadSimulationsTrackingConfiguration() *loadstrike.TrackingConfigurationSpec {
	return &loadstrike.TrackingConfigurationSpec{
		Source:                      withLoadSimulationsHttpSource(),
		Destination:                 withLoadSimulationsHttpDestination(),
		RunMode:                     "GenerateAndCorrelate",
		CorrelationTimeoutSeconds:   30,
		TimeoutSweepIntervalSeconds: 1,
		TimeoutBatchSize:            200,
		TimeoutCountsAsFailure:      true,
		MetricPrefix:                "orders_tracking",
		ExecuteOriginalScenarioRun:  false,
	}
}

func withLoadSimulationsTrackedContext() loadstrike.LoadStrikeContext {
	return loadstrike.Create().
		AddScenario(withLoadSimulationsBaselineScenario("orders.tracked").WithCrossPlatformTracking(withLoadSimulationsTrackingConfiguration())).
		WithRunnerKey(withLoadSimulationsRunnerKey).
		WithoutReports().
		BuildContext()
}

func withLoadSimulationsBuildTrackingPayload() loadstrike.TrackingPayload {
	builder := loadstrike.TrackingPayloadBuilder{
		Headers:            map[string]string{"X-Correlation-Id": "ord-1"},
		ContentType:        "application/json",
		MessagePayloadType: "json",
	}
	builder.SetBody(map[string]any{"trackingId": "ord-1"})
	return builder.Build()
}

func withLoadSimulationsRunResult() loadstrike.LoadStrikeRunResult {
	var result loadstrike.LoadStrikeRunResult
	return result
}

func withLoadSimulationsScenarioStats() loadstrike.LoadStrikeScenarioStats {
	var stats loadstrike.LoadStrikeScenarioStats
	return stats
}

func withLoadSimulationsWriteTempConfigFiles() withLoadSimulationsTempConfigPaths {
	return withLoadSimulationsTempConfigPaths{
		ConfigPath: "method-reference.loadstrike.config.json",
		InfraPath:  "method-reference.loadstrike.infra.json",
	}
}

        // Attach one pacing rule to the GET scenario.
func (reference WithLoadSimulationsMethodReference) AttachSingleSimulationExample() any {
    return withLoadSimulationsBaselineScenario().WithLoadSimulations(loadstrike.LoadStrikeSimulation.IterationsForConstant(1, 3))
}

// Attach multiple pacing rules in one call.
func (reference WithLoadSimulationsMethodReference) AttachMultipleSimulationsExample() any {
    return withLoadSimulationsBaselineScenario().WithLoadSimulations(loadstrike.LoadStrikeSimulation.IterationsForConstant(1, 3), loadstrike.LoadStrikeSimulation.Pause(loadstrike.DurationFromSeconds(1)))
}
