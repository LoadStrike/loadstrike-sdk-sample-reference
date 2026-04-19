package loadstrike_simulation

        import (
        	"time"

        	loadstrike "loadstrike.com/sdk/go"
        )

        const keepConstantRunnerKey = "runner_dummy_orders_reference"

type KeepConstantMethodReference struct{}

type keepConstantTempConfigPaths struct {
	ConfigPath string
	InfraPath  string
}

type keepConstantOrdersReportingSink struct{}

func newKeepConstantOrdersReportingSink() keepConstantOrdersReportingSink {
	return keepConstantOrdersReportingSink{}
}
func (keepConstantOrdersReportingSink) SinkName() string {
	return "orders-sample-sink"
}
func (keepConstantOrdersReportingSink) Init(loadstrike.LoadStrikeBaseContext, loadstrike.IConfiguration) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (keepConstantOrdersReportingSink) Start(loadstrike.LoadStrikeSessionStartInfo) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (keepConstantOrdersReportingSink) SaveRealtimeStats([]loadstrike.LoadStrikeScenarioStats) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (keepConstantOrdersReportingSink) SaveRealtimeMetrics(loadstrike.LoadStrikeMetricStats) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (keepConstantOrdersReportingSink) SaveRunResult(loadstrike.LoadStrikeRunResult) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (keepConstantOrdersReportingSink) Stop() loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (keepConstantOrdersReportingSink) Dispose() {}

type keepConstantOrdersRuntimePolicy struct{}

func newKeepConstantOrdersRuntimePolicy() keepConstantOrdersRuntimePolicy {
	return keepConstantOrdersRuntimePolicy{}
}
func (keepConstantOrdersRuntimePolicy) ShouldRunScenario(string) loadstrike.LoadStrikeBoolTask {
	return loadstrike.TaskFromBool(true)
}
func (keepConstantOrdersRuntimePolicy) BeforeScenario(string) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (keepConstantOrdersRuntimePolicy) AfterScenario(string, loadstrike.LoadStrikeScenarioRuntime) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (keepConstantOrdersRuntimePolicy) BeforeStep(string, string) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (keepConstantOrdersRuntimePolicy) AfterStep(string, string, loadstrike.LoadStrikeReply) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}

type keepConstantOrdersWorkerPlugin struct{}

func newKeepConstantOrdersWorkerPlugin() keepConstantOrdersWorkerPlugin {
	return keepConstantOrdersWorkerPlugin{}
}
func (keepConstantOrdersWorkerPlugin) PluginName() string {
	return "orders-sample-plugin"
}
func (keepConstantOrdersWorkerPlugin) Init(loadstrike.LoadStrikeBaseContext, loadstrike.IConfiguration) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (keepConstantOrdersWorkerPlugin) Start(loadstrike.LoadStrikeSessionStartInfo) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (keepConstantOrdersWorkerPlugin) GetData(loadstrike.LoadStrikeRunResult) loadstrike.LoadStrikeValueTask[loadstrike.LoadStrikePluginData] {
	return loadstrike.TaskFromResult(loadstrike.CreatePluginData("orders-sample-plugin"))
}
func (keepConstantOrdersWorkerPlugin) Stop() loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (keepConstantOrdersWorkerPlugin) Dispose() loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}

func keepConstantPerformOrderGetReply() loadstrike.LoadStrikeReply {
	return loadstrike.LoadStrikeResponse.Ok("200", int64(128), "ok", loadstrike.TimeSpan(3*time.Millisecond))
}

func keepConstantExecuteOrderGet(context loadstrike.LoadStrikeScenarioContext) loadstrike.LoadStrikeReply {
	return loadstrike.LoadStrikeStep.Run("get-order", context, func(loadstrike.LoadStrikeScenarioContext) loadstrike.LoadStrikeReply {
		return keepConstantPerformOrderGetReply()
	})
}

func keepConstantExecuteOrderGetAsync(context loadstrike.LoadStrikeScenarioContext) loadstrike.LoadStrikeValueTask[loadstrike.LoadStrikeReply] {
	return loadstrike.TaskFromResult(keepConstantExecuteOrderGet(context))
}

func keepConstantBaselineScenario(names ...string) loadstrike.LoadStrikeScenario {
	name := "orders.get-by-id"
	if len(names) > 0 && names[0] != "" {
		name = names[0]
	}
	return loadstrike.CreateScenario(name, keepConstantExecuteOrderGet).
		WithLoadSimulations(loadstrike.LoadStrikeSimulation.IterationsForConstant(1, 1)).
		WithoutWarmUp()
}

func keepConstantBaseRunner() loadstrike.LoadStrikeRunner {
	return loadstrike.Create().
		AddScenario(keepConstantBaselineScenario()).
		WithRunnerKey(keepConstantRunnerKey).
		WithTestSuite("orders-reference").
		WithTestName("orders-get-by-id").
		WithoutReports()
}

func keepConstantBaseContext() loadstrike.LoadStrikeContext {
	return keepConstantBaseRunner().BuildContext()
}

func keepConstantHttpSource() *loadstrike.EndpointSpec {
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

func keepConstantHttpDestination() *loadstrike.EndpointSpec {
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

func keepConstantTrackingConfiguration() *loadstrike.TrackingConfigurationSpec {
	return &loadstrike.TrackingConfigurationSpec{
		Source:                      keepConstantHttpSource(),
		Destination:                 keepConstantHttpDestination(),
		RunMode:                     "GenerateAndCorrelate",
		CorrelationTimeoutSeconds:   30,
		TimeoutSweepIntervalSeconds: 1,
		TimeoutBatchSize:            200,
		TimeoutCountsAsFailure:      true,
		MetricPrefix:                "orders_tracking",
		ExecuteOriginalScenarioRun:  false,
	}
}

func keepConstantTrackedContext() loadstrike.LoadStrikeContext {
	return loadstrike.Create().
		AddScenario(keepConstantBaselineScenario("orders.tracked").WithCrossPlatformTracking(keepConstantTrackingConfiguration())).
		WithRunnerKey(keepConstantRunnerKey).
		WithoutReports().
		BuildContext()
}

func keepConstantBuildTrackingPayload() loadstrike.TrackingPayload {
	builder := loadstrike.TrackingPayloadBuilder{
		Headers:            map[string]string{"X-Correlation-Id": "ord-1"},
		ContentType:        "application/json",
		MessagePayloadType: "json",
	}
	builder.SetBody(map[string]any{"trackingId": "ord-1"})
	return builder.Build()
}

func keepConstantRunResult() loadstrike.LoadStrikeRunResult {
	var result loadstrike.LoadStrikeRunResult
	return result
}

func keepConstantScenarioStats() loadstrike.LoadStrikeScenarioStats {
	var stats loadstrike.LoadStrikeScenarioStats
	return stats
}

func keepConstantWriteTempConfigFiles() keepConstantTempConfigPaths {
	return keepConstantTempConfigPaths{
		ConfigPath: "method-reference.loadstrike.config.json",
		InfraPath:  "method-reference.loadstrike.infra.json",
	}
}

        // Create the simulation shape named by this method reference.
func (reference KeepConstantMethodReference) CreatePrimarySimulationExample() any {
    return loadstrike.LoadStrikeSimulation.KeepConstant(2, loadstrike.DurationFromSeconds(20))
}

// Attach the simulation to the baseline GET scenario.
func (reference KeepConstantMethodReference) AttachSimulationToScenarioExample() any {
    return keepConstantBaselineScenario().WithLoadSimulations(loadstrike.LoadStrikeSimulation.KeepConstant(2, loadstrike.DurationFromSeconds(20)))
}
