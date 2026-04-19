package loadstrike_simulation

        import (
        	"time"

        	loadstrike "loadstrike.com/sdk/go"
        )

        const rampingInjectRunnerKey = "runner_dummy_orders_reference"

type RampingInjectMethodReference struct{}

type rampingInjectTempConfigPaths struct {
	ConfigPath string
	InfraPath  string
}

type rampingInjectOrdersReportingSink struct{}

func newRampingInjectOrdersReportingSink() rampingInjectOrdersReportingSink {
	return rampingInjectOrdersReportingSink{}
}
func (rampingInjectOrdersReportingSink) SinkName() string {
	return "orders-sample-sink"
}
func (rampingInjectOrdersReportingSink) Init(loadstrike.LoadStrikeBaseContext, loadstrike.IConfiguration) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (rampingInjectOrdersReportingSink) Start(loadstrike.LoadStrikeSessionStartInfo) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (rampingInjectOrdersReportingSink) SaveRealtimeStats([]loadstrike.LoadStrikeScenarioStats) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (rampingInjectOrdersReportingSink) SaveRealtimeMetrics(loadstrike.LoadStrikeMetricStats) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (rampingInjectOrdersReportingSink) SaveRunResult(loadstrike.LoadStrikeRunResult) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (rampingInjectOrdersReportingSink) Stop() loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (rampingInjectOrdersReportingSink) Dispose() {}

type rampingInjectOrdersRuntimePolicy struct{}

func newRampingInjectOrdersRuntimePolicy() rampingInjectOrdersRuntimePolicy {
	return rampingInjectOrdersRuntimePolicy{}
}
func (rampingInjectOrdersRuntimePolicy) ShouldRunScenario(string) loadstrike.LoadStrikeBoolTask {
	return loadstrike.TaskFromBool(true)
}
func (rampingInjectOrdersRuntimePolicy) BeforeScenario(string) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (rampingInjectOrdersRuntimePolicy) AfterScenario(string, loadstrike.LoadStrikeScenarioRuntime) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (rampingInjectOrdersRuntimePolicy) BeforeStep(string, string) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (rampingInjectOrdersRuntimePolicy) AfterStep(string, string, loadstrike.LoadStrikeReply) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}

type rampingInjectOrdersWorkerPlugin struct{}

func newRampingInjectOrdersWorkerPlugin() rampingInjectOrdersWorkerPlugin {
	return rampingInjectOrdersWorkerPlugin{}
}
func (rampingInjectOrdersWorkerPlugin) PluginName() string {
	return "orders-sample-plugin"
}
func (rampingInjectOrdersWorkerPlugin) Init(loadstrike.LoadStrikeBaseContext, loadstrike.IConfiguration) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (rampingInjectOrdersWorkerPlugin) Start(loadstrike.LoadStrikeSessionStartInfo) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (rampingInjectOrdersWorkerPlugin) GetData(loadstrike.LoadStrikeRunResult) loadstrike.LoadStrikeValueTask[loadstrike.LoadStrikePluginData] {
	return loadstrike.TaskFromResult(loadstrike.CreatePluginData("orders-sample-plugin"))
}
func (rampingInjectOrdersWorkerPlugin) Stop() loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (rampingInjectOrdersWorkerPlugin) Dispose() loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}

func rampingInjectPerformOrderGetReply() loadstrike.LoadStrikeReply {
	return loadstrike.LoadStrikeResponse.Ok("200", int64(128), "ok", loadstrike.TimeSpan(3*time.Millisecond))
}

func rampingInjectExecuteOrderGet(context loadstrike.LoadStrikeScenarioContext) loadstrike.LoadStrikeReply {
	return loadstrike.LoadStrikeStep.Run("get-order", context, func(loadstrike.LoadStrikeScenarioContext) loadstrike.LoadStrikeReply {
		return rampingInjectPerformOrderGetReply()
	})
}

func rampingInjectExecuteOrderGetAsync(context loadstrike.LoadStrikeScenarioContext) loadstrike.LoadStrikeValueTask[loadstrike.LoadStrikeReply] {
	return loadstrike.TaskFromResult(rampingInjectExecuteOrderGet(context))
}

func rampingInjectBaselineScenario(names ...string) loadstrike.LoadStrikeScenario {
	name := "orders.get-by-id"
	if len(names) > 0 && names[0] != "" {
		name = names[0]
	}
	return loadstrike.CreateScenario(name, rampingInjectExecuteOrderGet).
		WithLoadSimulations(loadstrike.LoadStrikeSimulation.IterationsForConstant(1, 1)).
		WithoutWarmUp()
}

func rampingInjectBaseRunner() loadstrike.LoadStrikeRunner {
	return loadstrike.Create().
		AddScenario(rampingInjectBaselineScenario()).
		WithRunnerKey(rampingInjectRunnerKey).
		WithTestSuite("orders-reference").
		WithTestName("orders-get-by-id").
		WithoutReports()
}

func rampingInjectBaseContext() loadstrike.LoadStrikeContext {
	return rampingInjectBaseRunner().BuildContext()
}

func rampingInjectHttpSource() *loadstrike.EndpointSpec {
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

func rampingInjectHttpDestination() *loadstrike.EndpointSpec {
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

func rampingInjectTrackingConfiguration() *loadstrike.TrackingConfigurationSpec {
	return &loadstrike.TrackingConfigurationSpec{
		Source:                      rampingInjectHttpSource(),
		Destination:                 rampingInjectHttpDestination(),
		RunMode:                     "GenerateAndCorrelate",
		CorrelationTimeoutSeconds:   30,
		TimeoutSweepIntervalSeconds: 1,
		TimeoutBatchSize:            200,
		TimeoutCountsAsFailure:      true,
		MetricPrefix:                "orders_tracking",
		ExecuteOriginalScenarioRun:  false,
	}
}

func rampingInjectTrackedContext() loadstrike.LoadStrikeContext {
	return loadstrike.Create().
		AddScenario(rampingInjectBaselineScenario("orders.tracked").WithCrossPlatformTracking(rampingInjectTrackingConfiguration())).
		WithRunnerKey(rampingInjectRunnerKey).
		WithoutReports().
		BuildContext()
}

func rampingInjectBuildTrackingPayload() loadstrike.TrackingPayload {
	builder := loadstrike.TrackingPayloadBuilder{
		Headers:            map[string]string{"X-Correlation-Id": "ord-1"},
		ContentType:        "application/json",
		MessagePayloadType: "json",
	}
	builder.SetBody(map[string]any{"trackingId": "ord-1"})
	return builder.Build()
}

func rampingInjectRunResult() loadstrike.LoadStrikeRunResult {
	var result loadstrike.LoadStrikeRunResult
	return result
}

func rampingInjectScenarioStats() loadstrike.LoadStrikeScenarioStats {
	var stats loadstrike.LoadStrikeScenarioStats
	return stats
}

func rampingInjectWriteTempConfigFiles() rampingInjectTempConfigPaths {
	return rampingInjectTempConfigPaths{
		ConfigPath: "method-reference.loadstrike.config.json",
		InfraPath:  "method-reference.loadstrike.infra.json",
	}
}

        // Create the simulation shape named by this method reference.
func (reference RampingInjectMethodReference) CreatePrimarySimulationExample() any {
    return loadstrike.LoadStrikeSimulation.RampingInject(5, loadstrike.DurationFromSeconds(0.25), loadstrike.DurationFromSeconds(20))
}

// Attach the simulation to the baseline GET scenario.
func (reference RampingInjectMethodReference) AttachSimulationToScenarioExample() any {
    return rampingInjectBaselineScenario().WithLoadSimulations(loadstrike.LoadStrikeSimulation.RampingInject(5, loadstrike.DurationFromSeconds(0.25), loadstrike.DurationFromSeconds(20)))
}
