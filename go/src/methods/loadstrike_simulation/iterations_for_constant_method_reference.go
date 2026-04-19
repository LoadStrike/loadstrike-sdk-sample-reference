package loadstrike_simulation

        import (
        	"time"

        	loadstrike "loadstrike.com/sdk/go"
        )

        const iterationsForConstantRunnerKey = "runner_dummy_orders_reference"

type IterationsForConstantMethodReference struct{}

type iterationsForConstantTempConfigPaths struct {
	ConfigPath string
	InfraPath  string
}

type iterationsForConstantOrdersReportingSink struct{}

func newIterationsForConstantOrdersReportingSink() iterationsForConstantOrdersReportingSink {
	return iterationsForConstantOrdersReportingSink{}
}
func (iterationsForConstantOrdersReportingSink) SinkName() string {
	return "orders-sample-sink"
}
func (iterationsForConstantOrdersReportingSink) Init(loadstrike.LoadStrikeBaseContext, loadstrike.IConfiguration) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (iterationsForConstantOrdersReportingSink) Start(loadstrike.LoadStrikeSessionStartInfo) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (iterationsForConstantOrdersReportingSink) SaveRealtimeStats([]loadstrike.LoadStrikeScenarioStats) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (iterationsForConstantOrdersReportingSink) SaveRealtimeMetrics(loadstrike.LoadStrikeMetricStats) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (iterationsForConstantOrdersReportingSink) SaveRunResult(loadstrike.LoadStrikeRunResult) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (iterationsForConstantOrdersReportingSink) Stop() loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (iterationsForConstantOrdersReportingSink) Dispose() {}

type iterationsForConstantOrdersRuntimePolicy struct{}

func newIterationsForConstantOrdersRuntimePolicy() iterationsForConstantOrdersRuntimePolicy {
	return iterationsForConstantOrdersRuntimePolicy{}
}
func (iterationsForConstantOrdersRuntimePolicy) ShouldRunScenario(string) loadstrike.LoadStrikeBoolTask {
	return loadstrike.TaskFromBool(true)
}
func (iterationsForConstantOrdersRuntimePolicy) BeforeScenario(string) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (iterationsForConstantOrdersRuntimePolicy) AfterScenario(string, loadstrike.LoadStrikeScenarioRuntime) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (iterationsForConstantOrdersRuntimePolicy) BeforeStep(string, string) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (iterationsForConstantOrdersRuntimePolicy) AfterStep(string, string, loadstrike.LoadStrikeReply) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}

type iterationsForConstantOrdersWorkerPlugin struct{}

func newIterationsForConstantOrdersWorkerPlugin() iterationsForConstantOrdersWorkerPlugin {
	return iterationsForConstantOrdersWorkerPlugin{}
}
func (iterationsForConstantOrdersWorkerPlugin) PluginName() string {
	return "orders-sample-plugin"
}
func (iterationsForConstantOrdersWorkerPlugin) Init(loadstrike.LoadStrikeBaseContext, loadstrike.IConfiguration) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (iterationsForConstantOrdersWorkerPlugin) Start(loadstrike.LoadStrikeSessionStartInfo) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (iterationsForConstantOrdersWorkerPlugin) GetData(loadstrike.LoadStrikeRunResult) loadstrike.LoadStrikeValueTask[loadstrike.LoadStrikePluginData] {
	return loadstrike.TaskFromResult(loadstrike.CreatePluginData("orders-sample-plugin"))
}
func (iterationsForConstantOrdersWorkerPlugin) Stop() loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (iterationsForConstantOrdersWorkerPlugin) Dispose() loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}

func iterationsForConstantPerformOrderGetReply() loadstrike.LoadStrikeReply {
	return loadstrike.LoadStrikeResponse.Ok("200", int64(128), "ok", loadstrike.TimeSpan(3*time.Millisecond))
}

func iterationsForConstantExecuteOrderGet(context loadstrike.LoadStrikeScenarioContext) loadstrike.LoadStrikeReply {
	return loadstrike.LoadStrikeStep.Run("get-order", context, func(loadstrike.LoadStrikeScenarioContext) loadstrike.LoadStrikeReply {
		return iterationsForConstantPerformOrderGetReply()
	})
}

func iterationsForConstantExecuteOrderGetAsync(context loadstrike.LoadStrikeScenarioContext) loadstrike.LoadStrikeValueTask[loadstrike.LoadStrikeReply] {
	return loadstrike.TaskFromResult(iterationsForConstantExecuteOrderGet(context))
}

func iterationsForConstantBaselineScenario(names ...string) loadstrike.LoadStrikeScenario {
	name := "orders.get-by-id"
	if len(names) > 0 && names[0] != "" {
		name = names[0]
	}
	return loadstrike.CreateScenario(name, iterationsForConstantExecuteOrderGet).
		WithLoadSimulations(loadstrike.LoadStrikeSimulation.IterationsForConstant(1, 1)).
		WithoutWarmUp()
}

func iterationsForConstantBaseRunner() loadstrike.LoadStrikeRunner {
	return loadstrike.Create().
		AddScenario(iterationsForConstantBaselineScenario()).
		WithRunnerKey(iterationsForConstantRunnerKey).
		WithTestSuite("orders-reference").
		WithTestName("orders-get-by-id").
		WithoutReports()
}

func iterationsForConstantBaseContext() loadstrike.LoadStrikeContext {
	return iterationsForConstantBaseRunner().BuildContext()
}

func iterationsForConstantHttpSource() *loadstrike.EndpointSpec {
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

func iterationsForConstantHttpDestination() *loadstrike.EndpointSpec {
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

func iterationsForConstantTrackingConfiguration() *loadstrike.TrackingConfigurationSpec {
	return &loadstrike.TrackingConfigurationSpec{
		Source:                      iterationsForConstantHttpSource(),
		Destination:                 iterationsForConstantHttpDestination(),
		RunMode:                     "GenerateAndCorrelate",
		CorrelationTimeoutSeconds:   30,
		TimeoutSweepIntervalSeconds: 1,
		TimeoutBatchSize:            200,
		TimeoutCountsAsFailure:      true,
		MetricPrefix:                "orders_tracking",
		ExecuteOriginalScenarioRun:  false,
	}
}

func iterationsForConstantTrackedContext() loadstrike.LoadStrikeContext {
	return loadstrike.Create().
		AddScenario(iterationsForConstantBaselineScenario("orders.tracked").WithCrossPlatformTracking(iterationsForConstantTrackingConfiguration())).
		WithRunnerKey(iterationsForConstantRunnerKey).
		WithoutReports().
		BuildContext()
}

func iterationsForConstantBuildTrackingPayload() loadstrike.TrackingPayload {
	builder := loadstrike.TrackingPayloadBuilder{
		Headers:            map[string]string{"X-Correlation-Id": "ord-1"},
		ContentType:        "application/json",
		MessagePayloadType: "json",
	}
	builder.SetBody(map[string]any{"trackingId": "ord-1"})
	return builder.Build()
}

func iterationsForConstantRunResult() loadstrike.LoadStrikeRunResult {
	var result loadstrike.LoadStrikeRunResult
	return result
}

func iterationsForConstantScenarioStats() loadstrike.LoadStrikeScenarioStats {
	var stats loadstrike.LoadStrikeScenarioStats
	return stats
}

func iterationsForConstantWriteTempConfigFiles() iterationsForConstantTempConfigPaths {
	return iterationsForConstantTempConfigPaths{
		ConfigPath: "method-reference.loadstrike.config.json",
		InfraPath:  "method-reference.loadstrike.infra.json",
	}
}

        // Create the simulation shape named by this method reference.
func (reference IterationsForConstantMethodReference) CreatePrimarySimulationExample() any {
    return loadstrike.LoadStrikeSimulation.IterationsForConstant(1, 3)
}

// Attach the simulation to the baseline GET scenario.
func (reference IterationsForConstantMethodReference) AttachSimulationToScenarioExample() any {
    return iterationsForConstantBaselineScenario().WithLoadSimulations(loadstrike.LoadStrikeSimulation.IterationsForConstant(1, 3))
}
