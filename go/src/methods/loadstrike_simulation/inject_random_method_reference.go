package loadstrike_simulation

        import (
        	"time"

        	loadstrike "loadstrike.com/sdk/go"
        )

        const injectRandomRunnerKey = "runner_dummy_orders_reference"

type InjectRandomMethodReference struct{}

type injectRandomTempConfigPaths struct {
	ConfigPath string
	InfraPath  string
}

type injectRandomOrdersReportingSink struct{}

func newInjectRandomOrdersReportingSink() injectRandomOrdersReportingSink {
	return injectRandomOrdersReportingSink{}
}
func (injectRandomOrdersReportingSink) SinkName() string {
	return "orders-sample-sink"
}
func (injectRandomOrdersReportingSink) Init(loadstrike.LoadStrikeBaseContext, loadstrike.IConfiguration) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (injectRandomOrdersReportingSink) Start(loadstrike.LoadStrikeSessionStartInfo) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (injectRandomOrdersReportingSink) SaveRealtimeStats([]loadstrike.LoadStrikeScenarioStats) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (injectRandomOrdersReportingSink) SaveRealtimeMetrics(loadstrike.LoadStrikeMetricStats) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (injectRandomOrdersReportingSink) SaveRunResult(loadstrike.LoadStrikeRunResult) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (injectRandomOrdersReportingSink) Stop() loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (injectRandomOrdersReportingSink) Dispose() {}

type injectRandomOrdersRuntimePolicy struct{}

func newInjectRandomOrdersRuntimePolicy() injectRandomOrdersRuntimePolicy {
	return injectRandomOrdersRuntimePolicy{}
}
func (injectRandomOrdersRuntimePolicy) ShouldRunScenario(string) loadstrike.LoadStrikeBoolTask {
	return loadstrike.TaskFromBool(true)
}
func (injectRandomOrdersRuntimePolicy) BeforeScenario(string) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (injectRandomOrdersRuntimePolicy) AfterScenario(string, loadstrike.LoadStrikeScenarioRuntime) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (injectRandomOrdersRuntimePolicy) BeforeStep(string, string) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (injectRandomOrdersRuntimePolicy) AfterStep(string, string, loadstrike.LoadStrikeReply) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}

type injectRandomOrdersWorkerPlugin struct{}

func newInjectRandomOrdersWorkerPlugin() injectRandomOrdersWorkerPlugin {
	return injectRandomOrdersWorkerPlugin{}
}
func (injectRandomOrdersWorkerPlugin) PluginName() string {
	return "orders-sample-plugin"
}
func (injectRandomOrdersWorkerPlugin) Init(loadstrike.LoadStrikeBaseContext, loadstrike.IConfiguration) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (injectRandomOrdersWorkerPlugin) Start(loadstrike.LoadStrikeSessionStartInfo) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (injectRandomOrdersWorkerPlugin) GetData(loadstrike.LoadStrikeRunResult) loadstrike.LoadStrikeValueTask[loadstrike.LoadStrikePluginData] {
	return loadstrike.TaskFromResult(loadstrike.CreatePluginData("orders-sample-plugin"))
}
func (injectRandomOrdersWorkerPlugin) Stop() loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (injectRandomOrdersWorkerPlugin) Dispose() loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}

func injectRandomPerformOrderGetReply() loadstrike.LoadStrikeReply {
	return loadstrike.LoadStrikeResponse.Ok("200", int64(128), "ok", loadstrike.TimeSpan(3*time.Millisecond))
}

func injectRandomExecuteOrderGet(context loadstrike.LoadStrikeScenarioContext) loadstrike.LoadStrikeReply {
	return loadstrike.LoadStrikeStep.Run("get-order", context, func(loadstrike.LoadStrikeScenarioContext) loadstrike.LoadStrikeReply {
		return injectRandomPerformOrderGetReply()
	})
}

func injectRandomExecuteOrderGetAsync(context loadstrike.LoadStrikeScenarioContext) loadstrike.LoadStrikeValueTask[loadstrike.LoadStrikeReply] {
	return loadstrike.TaskFromResult(injectRandomExecuteOrderGet(context))
}

func injectRandomBaselineScenario(names ...string) loadstrike.LoadStrikeScenario {
	name := "orders.get-by-id"
	if len(names) > 0 && names[0] != "" {
		name = names[0]
	}
	return loadstrike.CreateScenario(name, injectRandomExecuteOrderGet).
		WithLoadSimulations(loadstrike.LoadStrikeSimulation.IterationsForConstant(1, 1)).
		WithoutWarmUp()
}

func injectRandomBaseRunner() loadstrike.LoadStrikeRunner {
	return loadstrike.Create().
		AddScenario(injectRandomBaselineScenario()).
		WithRunnerKey(injectRandomRunnerKey).
		WithTestSuite("orders-reference").
		WithTestName("orders-get-by-id").
		WithoutReports()
}

func injectRandomBaseContext() loadstrike.LoadStrikeContext {
	return injectRandomBaseRunner().BuildContext()
}

func injectRandomHttpSource() *loadstrike.EndpointSpec {
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

func injectRandomHttpDestination() *loadstrike.EndpointSpec {
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

func injectRandomTrackingConfiguration() *loadstrike.TrackingConfigurationSpec {
	return &loadstrike.TrackingConfigurationSpec{
		Source:                      injectRandomHttpSource(),
		Destination:                 injectRandomHttpDestination(),
		RunMode:                     "GenerateAndCorrelate",
		CorrelationTimeoutSeconds:   30,
		TimeoutSweepIntervalSeconds: 1,
		TimeoutBatchSize:            200,
		TimeoutCountsAsFailure:      true,
		MetricPrefix:                "orders_tracking",
		ExecuteOriginalScenarioRun:  false,
	}
}

func injectRandomTrackedContext() loadstrike.LoadStrikeContext {
	return loadstrike.Create().
		AddScenario(injectRandomBaselineScenario("orders.tracked").WithCrossPlatformTracking(injectRandomTrackingConfiguration())).
		WithRunnerKey(injectRandomRunnerKey).
		WithoutReports().
		BuildContext()
}

func injectRandomBuildTrackingPayload() loadstrike.TrackingPayload {
	builder := loadstrike.TrackingPayloadBuilder{
		Headers:            map[string]string{"X-Correlation-Id": "ord-1"},
		ContentType:        "application/json",
		MessagePayloadType: "json",
	}
	builder.SetBody(map[string]any{"trackingId": "ord-1"})
	return builder.Build()
}

func injectRandomRunResult() loadstrike.LoadStrikeRunResult {
	var result loadstrike.LoadStrikeRunResult
	return result
}

func injectRandomScenarioStats() loadstrike.LoadStrikeScenarioStats {
	var stats loadstrike.LoadStrikeScenarioStats
	return stats
}

func injectRandomWriteTempConfigFiles() injectRandomTempConfigPaths {
	return injectRandomTempConfigPaths{
		ConfigPath: "method-reference.loadstrike.config.json",
		InfraPath:  "method-reference.loadstrike.infra.json",
	}
}

        // Create the simulation shape named by this method reference.
func (reference InjectRandomMethodReference) CreatePrimarySimulationExample() any {
    return loadstrike.LoadStrikeSimulation.InjectRandom(2, 4, loadstrike.DurationFromSeconds(0.25), loadstrike.DurationFromSeconds(20))
}

// Attach the simulation to the baseline GET scenario.
func (reference InjectRandomMethodReference) AttachSimulationToScenarioExample() any {
    return injectRandomBaselineScenario().WithLoadSimulations(loadstrike.LoadStrikeSimulation.InjectRandom(2, 4, loadstrike.DurationFromSeconds(0.25), loadstrike.DurationFromSeconds(20)))
}
