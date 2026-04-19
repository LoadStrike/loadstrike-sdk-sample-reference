package loadstrike_simulation

        import (
        	"time"

        	loadstrike "loadstrike.com/sdk/go"
        )

        const rampingConstantRunnerKey = "runner_dummy_orders_reference"

type RampingConstantMethodReference struct{}

type rampingConstantTempConfigPaths struct {
	ConfigPath string
	InfraPath  string
}

type rampingConstantOrdersReportingSink struct{}

func newRampingConstantOrdersReportingSink() rampingConstantOrdersReportingSink {
	return rampingConstantOrdersReportingSink{}
}
func (rampingConstantOrdersReportingSink) SinkName() string {
	return "orders-sample-sink"
}
func (rampingConstantOrdersReportingSink) Init(loadstrike.LoadStrikeBaseContext, loadstrike.IConfiguration) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (rampingConstantOrdersReportingSink) Start(loadstrike.LoadStrikeSessionStartInfo) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (rampingConstantOrdersReportingSink) SaveRealtimeStats([]loadstrike.LoadStrikeScenarioStats) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (rampingConstantOrdersReportingSink) SaveRealtimeMetrics(loadstrike.LoadStrikeMetricStats) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (rampingConstantOrdersReportingSink) SaveRunResult(loadstrike.LoadStrikeRunResult) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (rampingConstantOrdersReportingSink) Stop() loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (rampingConstantOrdersReportingSink) Dispose() {}

type rampingConstantOrdersRuntimePolicy struct{}

func newRampingConstantOrdersRuntimePolicy() rampingConstantOrdersRuntimePolicy {
	return rampingConstantOrdersRuntimePolicy{}
}
func (rampingConstantOrdersRuntimePolicy) ShouldRunScenario(string) loadstrike.LoadStrikeBoolTask {
	return loadstrike.TaskFromBool(true)
}
func (rampingConstantOrdersRuntimePolicy) BeforeScenario(string) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (rampingConstantOrdersRuntimePolicy) AfterScenario(string, loadstrike.LoadStrikeScenarioRuntime) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (rampingConstantOrdersRuntimePolicy) BeforeStep(string, string) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (rampingConstantOrdersRuntimePolicy) AfterStep(string, string, loadstrike.LoadStrikeReply) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}

type rampingConstantOrdersWorkerPlugin struct{}

func newRampingConstantOrdersWorkerPlugin() rampingConstantOrdersWorkerPlugin {
	return rampingConstantOrdersWorkerPlugin{}
}
func (rampingConstantOrdersWorkerPlugin) PluginName() string {
	return "orders-sample-plugin"
}
func (rampingConstantOrdersWorkerPlugin) Init(loadstrike.LoadStrikeBaseContext, loadstrike.IConfiguration) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (rampingConstantOrdersWorkerPlugin) Start(loadstrike.LoadStrikeSessionStartInfo) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (rampingConstantOrdersWorkerPlugin) GetData(loadstrike.LoadStrikeRunResult) loadstrike.LoadStrikeValueTask[loadstrike.LoadStrikePluginData] {
	return loadstrike.TaskFromResult(loadstrike.CreatePluginData("orders-sample-plugin"))
}
func (rampingConstantOrdersWorkerPlugin) Stop() loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (rampingConstantOrdersWorkerPlugin) Dispose() loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}

func rampingConstantPerformOrderGetReply() loadstrike.LoadStrikeReply {
	return loadstrike.LoadStrikeResponse.Ok("200", int64(128), "ok", loadstrike.TimeSpan(3*time.Millisecond))
}

func rampingConstantExecuteOrderGet(context loadstrike.LoadStrikeScenarioContext) loadstrike.LoadStrikeReply {
	return loadstrike.LoadStrikeStep.Run("get-order", context, func(loadstrike.LoadStrikeScenarioContext) loadstrike.LoadStrikeReply {
		return rampingConstantPerformOrderGetReply()
	})
}

func rampingConstantExecuteOrderGetAsync(context loadstrike.LoadStrikeScenarioContext) loadstrike.LoadStrikeValueTask[loadstrike.LoadStrikeReply] {
	return loadstrike.TaskFromResult(rampingConstantExecuteOrderGet(context))
}

func rampingConstantBaselineScenario(names ...string) loadstrike.LoadStrikeScenario {
	name := "orders.get-by-id"
	if len(names) > 0 && names[0] != "" {
		name = names[0]
	}
	return loadstrike.CreateScenario(name, rampingConstantExecuteOrderGet).
		WithLoadSimulations(loadstrike.LoadStrikeSimulation.IterationsForConstant(1, 1)).
		WithoutWarmUp()
}

func rampingConstantBaseRunner() loadstrike.LoadStrikeRunner {
	return loadstrike.Create().
		AddScenario(rampingConstantBaselineScenario()).
		WithRunnerKey(rampingConstantRunnerKey).
		WithTestSuite("orders-reference").
		WithTestName("orders-get-by-id").
		WithoutReports()
}

func rampingConstantBaseContext() loadstrike.LoadStrikeContext {
	return rampingConstantBaseRunner().BuildContext()
}

func rampingConstantHttpSource() *loadstrike.EndpointSpec {
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

func rampingConstantHttpDestination() *loadstrike.EndpointSpec {
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

func rampingConstantTrackingConfiguration() *loadstrike.TrackingConfigurationSpec {
	return &loadstrike.TrackingConfigurationSpec{
		Source:                      rampingConstantHttpSource(),
		Destination:                 rampingConstantHttpDestination(),
		RunMode:                     "GenerateAndCorrelate",
		CorrelationTimeoutSeconds:   30,
		TimeoutSweepIntervalSeconds: 1,
		TimeoutBatchSize:            200,
		TimeoutCountsAsFailure:      true,
		MetricPrefix:                "orders_tracking",
		ExecuteOriginalScenarioRun:  false,
	}
}

func rampingConstantTrackedContext() loadstrike.LoadStrikeContext {
	return loadstrike.Create().
		AddScenario(rampingConstantBaselineScenario("orders.tracked").WithCrossPlatformTracking(rampingConstantTrackingConfiguration())).
		WithRunnerKey(rampingConstantRunnerKey).
		WithoutReports().
		BuildContext()
}

func rampingConstantBuildTrackingPayload() loadstrike.TrackingPayload {
	builder := loadstrike.TrackingPayloadBuilder{
		Headers:            map[string]string{"X-Correlation-Id": "ord-1"},
		ContentType:        "application/json",
		MessagePayloadType: "json",
	}
	builder.SetBody(map[string]any{"trackingId": "ord-1"})
	return builder.Build()
}

func rampingConstantRunResult() loadstrike.LoadStrikeRunResult {
	var result loadstrike.LoadStrikeRunResult
	return result
}

func rampingConstantScenarioStats() loadstrike.LoadStrikeScenarioStats {
	var stats loadstrike.LoadStrikeScenarioStats
	return stats
}

func rampingConstantWriteTempConfigFiles() rampingConstantTempConfigPaths {
	return rampingConstantTempConfigPaths{
		ConfigPath: "method-reference.loadstrike.config.json",
		InfraPath:  "method-reference.loadstrike.infra.json",
	}
}

        // Create the simulation shape named by this method reference.
func (reference RampingConstantMethodReference) CreatePrimarySimulationExample() any {
    return loadstrike.LoadStrikeSimulation.RampingConstant(3, loadstrike.DurationFromSeconds(20))
}

// Attach the simulation to the baseline GET scenario.
func (reference RampingConstantMethodReference) AttachSimulationToScenarioExample() any {
    return rampingConstantBaselineScenario().WithLoadSimulations(loadstrike.LoadStrikeSimulation.RampingConstant(3, loadstrike.DurationFromSeconds(20)))
}
