package loadstrike_runner

        import (
        	"time"

        	loadstrike "loadstrike.com/sdk/go"
        )

        const withAgentTargetScenariosRunnerKey = "runner_dummy_orders_reference"

type WithAgentTargetScenariosMethodReference struct{}

type withAgentTargetScenariosTempConfigPaths struct {
	ConfigPath string
	InfraPath  string
}

type withAgentTargetScenariosOrdersReportingSink struct{}

func newWithAgentTargetScenariosOrdersReportingSink() withAgentTargetScenariosOrdersReportingSink {
	return withAgentTargetScenariosOrdersReportingSink{}
}
func (withAgentTargetScenariosOrdersReportingSink) SinkName() string {
	return "orders-sample-sink"
}
func (withAgentTargetScenariosOrdersReportingSink) Init(loadstrike.LoadStrikeBaseContext, loadstrike.IConfiguration) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withAgentTargetScenariosOrdersReportingSink) Start(loadstrike.LoadStrikeSessionStartInfo) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withAgentTargetScenariosOrdersReportingSink) SaveRealtimeStats([]loadstrike.LoadStrikeScenarioStats) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withAgentTargetScenariosOrdersReportingSink) SaveRealtimeMetrics(loadstrike.LoadStrikeMetricStats) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withAgentTargetScenariosOrdersReportingSink) SaveRunResult(loadstrike.LoadStrikeRunResult) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withAgentTargetScenariosOrdersReportingSink) Stop() loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withAgentTargetScenariosOrdersReportingSink) Dispose() {}

type withAgentTargetScenariosOrdersRuntimePolicy struct{}

func newWithAgentTargetScenariosOrdersRuntimePolicy() withAgentTargetScenariosOrdersRuntimePolicy {
	return withAgentTargetScenariosOrdersRuntimePolicy{}
}
func (withAgentTargetScenariosOrdersRuntimePolicy) ShouldRunScenario(string) loadstrike.LoadStrikeBoolTask {
	return loadstrike.TaskFromBool(true)
}
func (withAgentTargetScenariosOrdersRuntimePolicy) BeforeScenario(string) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withAgentTargetScenariosOrdersRuntimePolicy) AfterScenario(string, loadstrike.LoadStrikeScenarioRuntime) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withAgentTargetScenariosOrdersRuntimePolicy) BeforeStep(string, string) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withAgentTargetScenariosOrdersRuntimePolicy) AfterStep(string, string, loadstrike.LoadStrikeReply) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}

type withAgentTargetScenariosOrdersWorkerPlugin struct{}

func newWithAgentTargetScenariosOrdersWorkerPlugin() withAgentTargetScenariosOrdersWorkerPlugin {
	return withAgentTargetScenariosOrdersWorkerPlugin{}
}
func (withAgentTargetScenariosOrdersWorkerPlugin) PluginName() string {
	return "orders-sample-plugin"
}
func (withAgentTargetScenariosOrdersWorkerPlugin) Init(loadstrike.LoadStrikeBaseContext, loadstrike.IConfiguration) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withAgentTargetScenariosOrdersWorkerPlugin) Start(loadstrike.LoadStrikeSessionStartInfo) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withAgentTargetScenariosOrdersWorkerPlugin) GetData(loadstrike.LoadStrikeRunResult) loadstrike.LoadStrikeValueTask[loadstrike.LoadStrikePluginData] {
	return loadstrike.TaskFromResult(loadstrike.CreatePluginData("orders-sample-plugin"))
}
func (withAgentTargetScenariosOrdersWorkerPlugin) Stop() loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withAgentTargetScenariosOrdersWorkerPlugin) Dispose() loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}

func withAgentTargetScenariosPerformOrderGetReply() loadstrike.LoadStrikeReply {
	return loadstrike.LoadStrikeResponse.Ok("200", int64(128), "ok", loadstrike.TimeSpan(3*time.Millisecond))
}

func withAgentTargetScenariosExecuteOrderGet(context loadstrike.LoadStrikeScenarioContext) loadstrike.LoadStrikeReply {
	return loadstrike.LoadStrikeStep.Run("get-order", context, func(loadstrike.LoadStrikeScenarioContext) loadstrike.LoadStrikeReply {
		return withAgentTargetScenariosPerformOrderGetReply()
	})
}

func withAgentTargetScenariosExecuteOrderGetAsync(context loadstrike.LoadStrikeScenarioContext) loadstrike.LoadStrikeValueTask[loadstrike.LoadStrikeReply] {
	return loadstrike.TaskFromResult(withAgentTargetScenariosExecuteOrderGet(context))
}

func withAgentTargetScenariosBaselineScenario(names ...string) loadstrike.LoadStrikeScenario {
	name := "orders.get-by-id"
	if len(names) > 0 && names[0] != "" {
		name = names[0]
	}
	return loadstrike.CreateScenario(name, withAgentTargetScenariosExecuteOrderGet).
		WithLoadSimulations(loadstrike.LoadStrikeSimulation.IterationsForConstant(1, 1)).
		WithoutWarmUp()
}

func withAgentTargetScenariosBaseRunner() loadstrike.LoadStrikeRunner {
	return loadstrike.Create().
		AddScenario(withAgentTargetScenariosBaselineScenario()).
		WithRunnerKey(withAgentTargetScenariosRunnerKey).
		WithTestSuite("orders-reference").
		WithTestName("orders-get-by-id").
		WithoutReports()
}

func withAgentTargetScenariosBaseContext() loadstrike.LoadStrikeContext {
	return withAgentTargetScenariosBaseRunner().BuildContext()
}

func withAgentTargetScenariosHttpSource() *loadstrike.EndpointSpec {
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

func withAgentTargetScenariosHttpDestination() *loadstrike.EndpointSpec {
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

func withAgentTargetScenariosTrackingConfiguration() *loadstrike.TrackingConfigurationSpec {
	return &loadstrike.TrackingConfigurationSpec{
		Source:                      withAgentTargetScenariosHttpSource(),
		Destination:                 withAgentTargetScenariosHttpDestination(),
		RunMode:                     "GenerateAndCorrelate",
		CorrelationTimeoutSeconds:   30,
		TimeoutSweepIntervalSeconds: 1,
		TimeoutBatchSize:            200,
		TimeoutCountsAsFailure:      true,
		MetricPrefix:                "orders_tracking",
		ExecuteOriginalScenarioRun:  false,
	}
}

func withAgentTargetScenariosTrackedContext() loadstrike.LoadStrikeContext {
	return loadstrike.Create().
		AddScenario(withAgentTargetScenariosBaselineScenario("orders.tracked").WithCrossPlatformTracking(withAgentTargetScenariosTrackingConfiguration())).
		WithRunnerKey(withAgentTargetScenariosRunnerKey).
		WithoutReports().
		BuildContext()
}

func withAgentTargetScenariosBuildTrackingPayload() loadstrike.TrackingPayload {
	builder := loadstrike.TrackingPayloadBuilder{
		Headers:            map[string]string{"X-Correlation-Id": "ord-1"},
		ContentType:        "application/json",
		MessagePayloadType: "json",
	}
	builder.SetBody(map[string]any{"trackingId": "ord-1"})
	return builder.Build()
}

func withAgentTargetScenariosRunResult() loadstrike.LoadStrikeRunResult {
	var result loadstrike.LoadStrikeRunResult
	return result
}

func withAgentTargetScenariosScenarioStats() loadstrike.LoadStrikeScenarioStats {
	var stats loadstrike.LoadStrikeScenarioStats
	return stats
}

func withAgentTargetScenariosWriteTempConfigFiles() withAgentTargetScenariosTempConfigPaths {
	return withAgentTargetScenariosTempConfigPaths{
		ConfigPath: "method-reference.loadstrike.config.json",
		InfraPath:  "method-reference.loadstrike.infra.json",
	}
}

        // Target a single scenario by name.
func (reference WithAgentTargetScenariosMethodReference) SelectOneScenarioExample() any {
    return withAgentTargetScenariosBaseRunner().WithAgentTargetScenarios("orders.get-by-id")
}

// Target more than one named scenario in the same call.
func (reference WithAgentTargetScenariosMethodReference) SelectMultipleScenariosExample() any {
    return withAgentTargetScenariosBaseRunner().WithAgentTargetScenarios("orders.get-by-id", "orders.audit")
}
