package loadstrike_simulation

        import (
        	"time"

        	loadstrike "loadstrike.com/sdk/go"
        )

        const iterationsForInjectRunnerKey = "runner_dummy_orders_reference"

type IterationsForInjectMethodReference struct{}

type iterationsForInjectTempConfigPaths struct {
	ConfigPath string
	InfraPath  string
}

type iterationsForInjectOrdersReportingSink struct{}

func newIterationsForInjectOrdersReportingSink() iterationsForInjectOrdersReportingSink {
	return iterationsForInjectOrdersReportingSink{}
}
func (iterationsForInjectOrdersReportingSink) SinkName() string {
	return "orders-sample-sink"
}
func (iterationsForInjectOrdersReportingSink) Init(loadstrike.LoadStrikeBaseContext, loadstrike.IConfiguration) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (iterationsForInjectOrdersReportingSink) Start(loadstrike.LoadStrikeSessionStartInfo) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (iterationsForInjectOrdersReportingSink) SaveRealtimeStats([]loadstrike.LoadStrikeScenarioStats) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (iterationsForInjectOrdersReportingSink) SaveRealtimeMetrics(loadstrike.LoadStrikeMetricStats) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (iterationsForInjectOrdersReportingSink) SaveRunResult(loadstrike.LoadStrikeRunResult) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (iterationsForInjectOrdersReportingSink) Stop() loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (iterationsForInjectOrdersReportingSink) Dispose() {}

type iterationsForInjectOrdersRuntimePolicy struct{}

func newIterationsForInjectOrdersRuntimePolicy() iterationsForInjectOrdersRuntimePolicy {
	return iterationsForInjectOrdersRuntimePolicy{}
}
func (iterationsForInjectOrdersRuntimePolicy) ShouldRunScenario(string) loadstrike.LoadStrikeBoolTask {
	return loadstrike.TaskFromBool(true)
}
func (iterationsForInjectOrdersRuntimePolicy) BeforeScenario(string) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (iterationsForInjectOrdersRuntimePolicy) AfterScenario(string, loadstrike.LoadStrikeScenarioRuntime) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (iterationsForInjectOrdersRuntimePolicy) BeforeStep(string, string) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (iterationsForInjectOrdersRuntimePolicy) AfterStep(string, string, loadstrike.LoadStrikeReply) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}

type iterationsForInjectOrdersWorkerPlugin struct{}

func newIterationsForInjectOrdersWorkerPlugin() iterationsForInjectOrdersWorkerPlugin {
	return iterationsForInjectOrdersWorkerPlugin{}
}
func (iterationsForInjectOrdersWorkerPlugin) PluginName() string {
	return "orders-sample-plugin"
}
func (iterationsForInjectOrdersWorkerPlugin) Init(loadstrike.LoadStrikeBaseContext, loadstrike.IConfiguration) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (iterationsForInjectOrdersWorkerPlugin) Start(loadstrike.LoadStrikeSessionStartInfo) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (iterationsForInjectOrdersWorkerPlugin) GetData(loadstrike.LoadStrikeRunResult) loadstrike.LoadStrikeValueTask[loadstrike.LoadStrikePluginData] {
	return loadstrike.TaskFromResult(loadstrike.CreatePluginData("orders-sample-plugin"))
}
func (iterationsForInjectOrdersWorkerPlugin) Stop() loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (iterationsForInjectOrdersWorkerPlugin) Dispose() loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}

func iterationsForInjectPerformOrderGetReply() loadstrike.LoadStrikeReply {
	return loadstrike.LoadStrikeResponse.Ok("200", int64(128), "ok", loadstrike.TimeSpan(3*time.Millisecond))
}

func iterationsForInjectExecuteOrderGet(context loadstrike.LoadStrikeScenarioContext) loadstrike.LoadStrikeReply {
	return loadstrike.LoadStrikeStep.Run("get-order", context, func(loadstrike.LoadStrikeScenarioContext) loadstrike.LoadStrikeReply {
		return iterationsForInjectPerformOrderGetReply()
	})
}

func iterationsForInjectExecuteOrderGetAsync(context loadstrike.LoadStrikeScenarioContext) loadstrike.LoadStrikeValueTask[loadstrike.LoadStrikeReply] {
	return loadstrike.TaskFromResult(iterationsForInjectExecuteOrderGet(context))
}

func iterationsForInjectBaselineScenario(names ...string) loadstrike.LoadStrikeScenario {
	name := "orders.get-by-id"
	if len(names) > 0 && names[0] != "" {
		name = names[0]
	}
	return loadstrike.CreateScenario(name, iterationsForInjectExecuteOrderGet).
		WithLoadSimulations(loadstrike.LoadStrikeSimulation.IterationsForConstant(1, 1)).
		WithoutWarmUp()
}

func iterationsForInjectBaseRunner() loadstrike.LoadStrikeRunner {
	return loadstrike.Create().
		AddScenario(iterationsForInjectBaselineScenario()).
		WithRunnerKey(iterationsForInjectRunnerKey).
		WithTestSuite("orders-reference").
		WithTestName("orders-get-by-id").
		WithoutReports()
}

func iterationsForInjectBaseContext() loadstrike.LoadStrikeContext {
	return iterationsForInjectBaseRunner().BuildContext()
}

func iterationsForInjectHttpSource() *loadstrike.EndpointSpec {
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

func iterationsForInjectHttpDestination() *loadstrike.EndpointSpec {
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

func iterationsForInjectTrackingConfiguration() *loadstrike.TrackingConfigurationSpec {
	return &loadstrike.TrackingConfigurationSpec{
		Source:                      iterationsForInjectHttpSource(),
		Destination:                 iterationsForInjectHttpDestination(),
		RunMode:                     "GenerateAndCorrelate",
		CorrelationTimeoutSeconds:   30,
		TimeoutSweepIntervalSeconds: 1,
		TimeoutBatchSize:            200,
		TimeoutCountsAsFailure:      true,
		MetricPrefix:                "orders_tracking",
		ExecuteOriginalScenarioRun:  false,
	}
}

func iterationsForInjectTrackedContext() loadstrike.LoadStrikeContext {
	return loadstrike.Create().
		AddScenario(iterationsForInjectBaselineScenario("orders.tracked").WithCrossPlatformTracking(iterationsForInjectTrackingConfiguration())).
		WithRunnerKey(iterationsForInjectRunnerKey).
		WithoutReports().
		BuildContext()
}

func iterationsForInjectBuildTrackingPayload() loadstrike.TrackingPayload {
	builder := loadstrike.TrackingPayloadBuilder{
		Headers:            map[string]string{"X-Correlation-Id": "ord-1"},
		ContentType:        "application/json",
		MessagePayloadType: "json",
	}
	builder.SetBody(map[string]any{"trackingId": "ord-1"})
	return builder.Build()
}

func iterationsForInjectRunResult() loadstrike.LoadStrikeRunResult {
	var result loadstrike.LoadStrikeRunResult
	return result
}

func iterationsForInjectScenarioStats() loadstrike.LoadStrikeScenarioStats {
	var stats loadstrike.LoadStrikeScenarioStats
	return stats
}

func iterationsForInjectWriteTempConfigFiles() iterationsForInjectTempConfigPaths {
	return iterationsForInjectTempConfigPaths{
		ConfigPath: "method-reference.loadstrike.config.json",
		InfraPath:  "method-reference.loadstrike.infra.json",
	}
}

        // Create the simulation shape named by this method reference.
func (reference IterationsForInjectMethodReference) CreatePrimarySimulationExample() any {
    return loadstrike.LoadStrikeSimulation.IterationsForInject(2, loadstrike.DurationFromSeconds(0.25), 5)
}

// Attach the simulation to the baseline GET scenario.
func (reference IterationsForInjectMethodReference) AttachSimulationToScenarioExample() any {
    return iterationsForInjectBaselineScenario().WithLoadSimulations(loadstrike.LoadStrikeSimulation.IterationsForInject(2, loadstrike.DurationFromSeconds(0.25), 5))
}
