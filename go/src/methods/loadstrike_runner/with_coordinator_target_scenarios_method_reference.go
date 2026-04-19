package loadstrike_runner

        import (
        	"time"

        	loadstrike "loadstrike.com/sdk/go"
        )

        const withCoordinatorTargetScenariosRunnerKey = "runner_dummy_orders_reference"

type WithCoordinatorTargetScenariosMethodReference struct{}

type withCoordinatorTargetScenariosTempConfigPaths struct {
	ConfigPath string
	InfraPath  string
}

type withCoordinatorTargetScenariosOrdersReportingSink struct{}

func newWithCoordinatorTargetScenariosOrdersReportingSink() withCoordinatorTargetScenariosOrdersReportingSink {
	return withCoordinatorTargetScenariosOrdersReportingSink{}
}
func (withCoordinatorTargetScenariosOrdersReportingSink) SinkName() string {
	return "orders-sample-sink"
}
func (withCoordinatorTargetScenariosOrdersReportingSink) Init(loadstrike.LoadStrikeBaseContext, loadstrike.IConfiguration) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withCoordinatorTargetScenariosOrdersReportingSink) Start(loadstrike.LoadStrikeSessionStartInfo) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withCoordinatorTargetScenariosOrdersReportingSink) SaveRealtimeStats([]loadstrike.LoadStrikeScenarioStats) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withCoordinatorTargetScenariosOrdersReportingSink) SaveRealtimeMetrics(loadstrike.LoadStrikeMetricStats) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withCoordinatorTargetScenariosOrdersReportingSink) SaveRunResult(loadstrike.LoadStrikeRunResult) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withCoordinatorTargetScenariosOrdersReportingSink) Stop() loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withCoordinatorTargetScenariosOrdersReportingSink) Dispose() {}

type withCoordinatorTargetScenariosOrdersRuntimePolicy struct{}

func newWithCoordinatorTargetScenariosOrdersRuntimePolicy() withCoordinatorTargetScenariosOrdersRuntimePolicy {
	return withCoordinatorTargetScenariosOrdersRuntimePolicy{}
}
func (withCoordinatorTargetScenariosOrdersRuntimePolicy) ShouldRunScenario(string) loadstrike.LoadStrikeBoolTask {
	return loadstrike.TaskFromBool(true)
}
func (withCoordinatorTargetScenariosOrdersRuntimePolicy) BeforeScenario(string) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withCoordinatorTargetScenariosOrdersRuntimePolicy) AfterScenario(string, loadstrike.LoadStrikeScenarioRuntime) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withCoordinatorTargetScenariosOrdersRuntimePolicy) BeforeStep(string, string) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withCoordinatorTargetScenariosOrdersRuntimePolicy) AfterStep(string, string, loadstrike.LoadStrikeReply) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}

type withCoordinatorTargetScenariosOrdersWorkerPlugin struct{}

func newWithCoordinatorTargetScenariosOrdersWorkerPlugin() withCoordinatorTargetScenariosOrdersWorkerPlugin {
	return withCoordinatorTargetScenariosOrdersWorkerPlugin{}
}
func (withCoordinatorTargetScenariosOrdersWorkerPlugin) PluginName() string {
	return "orders-sample-plugin"
}
func (withCoordinatorTargetScenariosOrdersWorkerPlugin) Init(loadstrike.LoadStrikeBaseContext, loadstrike.IConfiguration) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withCoordinatorTargetScenariosOrdersWorkerPlugin) Start(loadstrike.LoadStrikeSessionStartInfo) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withCoordinatorTargetScenariosOrdersWorkerPlugin) GetData(loadstrike.LoadStrikeRunResult) loadstrike.LoadStrikeValueTask[loadstrike.LoadStrikePluginData] {
	return loadstrike.TaskFromResult(loadstrike.CreatePluginData("orders-sample-plugin"))
}
func (withCoordinatorTargetScenariosOrdersWorkerPlugin) Stop() loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withCoordinatorTargetScenariosOrdersWorkerPlugin) Dispose() loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}

func withCoordinatorTargetScenariosPerformOrderGetReply() loadstrike.LoadStrikeReply {
	return loadstrike.LoadStrikeResponse.Ok("200", int64(128), "ok", loadstrike.TimeSpan(3*time.Millisecond))
}

func withCoordinatorTargetScenariosExecuteOrderGet(context loadstrike.LoadStrikeScenarioContext) loadstrike.LoadStrikeReply {
	return loadstrike.LoadStrikeStep.Run("get-order", context, func(loadstrike.LoadStrikeScenarioContext) loadstrike.LoadStrikeReply {
		return withCoordinatorTargetScenariosPerformOrderGetReply()
	})
}

func withCoordinatorTargetScenariosExecuteOrderGetAsync(context loadstrike.LoadStrikeScenarioContext) loadstrike.LoadStrikeValueTask[loadstrike.LoadStrikeReply] {
	return loadstrike.TaskFromResult(withCoordinatorTargetScenariosExecuteOrderGet(context))
}

func withCoordinatorTargetScenariosBaselineScenario(names ...string) loadstrike.LoadStrikeScenario {
	name := "orders.get-by-id"
	if len(names) > 0 && names[0] != "" {
		name = names[0]
	}
	return loadstrike.CreateScenario(name, withCoordinatorTargetScenariosExecuteOrderGet).
		WithLoadSimulations(loadstrike.LoadStrikeSimulation.IterationsForConstant(1, 1)).
		WithoutWarmUp()
}

func withCoordinatorTargetScenariosBaseRunner() loadstrike.LoadStrikeRunner {
	return loadstrike.Create().
		AddScenario(withCoordinatorTargetScenariosBaselineScenario()).
		WithRunnerKey(withCoordinatorTargetScenariosRunnerKey).
		WithTestSuite("orders-reference").
		WithTestName("orders-get-by-id").
		WithoutReports()
}

func withCoordinatorTargetScenariosBaseContext() loadstrike.LoadStrikeContext {
	return withCoordinatorTargetScenariosBaseRunner().BuildContext()
}

func withCoordinatorTargetScenariosHttpSource() *loadstrike.EndpointSpec {
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

func withCoordinatorTargetScenariosHttpDestination() *loadstrike.EndpointSpec {
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

func withCoordinatorTargetScenariosTrackingConfiguration() *loadstrike.TrackingConfigurationSpec {
	return &loadstrike.TrackingConfigurationSpec{
		Source:                      withCoordinatorTargetScenariosHttpSource(),
		Destination:                 withCoordinatorTargetScenariosHttpDestination(),
		RunMode:                     "GenerateAndCorrelate",
		CorrelationTimeoutSeconds:   30,
		TimeoutSweepIntervalSeconds: 1,
		TimeoutBatchSize:            200,
		TimeoutCountsAsFailure:      true,
		MetricPrefix:                "orders_tracking",
		ExecuteOriginalScenarioRun:  false,
	}
}

func withCoordinatorTargetScenariosTrackedContext() loadstrike.LoadStrikeContext {
	return loadstrike.Create().
		AddScenario(withCoordinatorTargetScenariosBaselineScenario("orders.tracked").WithCrossPlatformTracking(withCoordinatorTargetScenariosTrackingConfiguration())).
		WithRunnerKey(withCoordinatorTargetScenariosRunnerKey).
		WithoutReports().
		BuildContext()
}

func withCoordinatorTargetScenariosBuildTrackingPayload() loadstrike.TrackingPayload {
	builder := loadstrike.TrackingPayloadBuilder{
		Headers:            map[string]string{"X-Correlation-Id": "ord-1"},
		ContentType:        "application/json",
		MessagePayloadType: "json",
	}
	builder.SetBody(map[string]any{"trackingId": "ord-1"})
	return builder.Build()
}

func withCoordinatorTargetScenariosRunResult() loadstrike.LoadStrikeRunResult {
	var result loadstrike.LoadStrikeRunResult
	return result
}

func withCoordinatorTargetScenariosScenarioStats() loadstrike.LoadStrikeScenarioStats {
	var stats loadstrike.LoadStrikeScenarioStats
	return stats
}

func withCoordinatorTargetScenariosWriteTempConfigFiles() withCoordinatorTargetScenariosTempConfigPaths {
	return withCoordinatorTargetScenariosTempConfigPaths{
		ConfigPath: "method-reference.loadstrike.config.json",
		InfraPath:  "method-reference.loadstrike.infra.json",
	}
}

        // Target a single scenario by name.
func (reference WithCoordinatorTargetScenariosMethodReference) SelectOneScenarioExample() any {
    return withCoordinatorTargetScenariosBaseRunner().WithCoordinatorTargetScenarios("orders.get-by-id")
}

// Target more than one named scenario in the same call.
func (reference WithCoordinatorTargetScenariosMethodReference) SelectMultipleScenariosExample() any {
    return withCoordinatorTargetScenariosBaseRunner().WithCoordinatorTargetScenarios("orders.get-by-id", "orders.audit")
}
