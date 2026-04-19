package loadstrike_runner

        import (
        	"time"

        	loadstrike "loadstrike.com/sdk/go"
        )

        const withTargetScenariosRunnerKey = "runner_dummy_orders_reference"

type WithTargetScenariosMethodReference struct{}

type withTargetScenariosTempConfigPaths struct {
	ConfigPath string
	InfraPath  string
}

type withTargetScenariosOrdersReportingSink struct{}

func newWithTargetScenariosOrdersReportingSink() withTargetScenariosOrdersReportingSink {
	return withTargetScenariosOrdersReportingSink{}
}
func (withTargetScenariosOrdersReportingSink) SinkName() string {
	return "orders-sample-sink"
}
func (withTargetScenariosOrdersReportingSink) Init(loadstrike.LoadStrikeBaseContext, loadstrike.IConfiguration) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withTargetScenariosOrdersReportingSink) Start(loadstrike.LoadStrikeSessionStartInfo) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withTargetScenariosOrdersReportingSink) SaveRealtimeStats([]loadstrike.LoadStrikeScenarioStats) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withTargetScenariosOrdersReportingSink) SaveRealtimeMetrics(loadstrike.LoadStrikeMetricStats) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withTargetScenariosOrdersReportingSink) SaveRunResult(loadstrike.LoadStrikeRunResult) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withTargetScenariosOrdersReportingSink) Stop() loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withTargetScenariosOrdersReportingSink) Dispose() {}

type withTargetScenariosOrdersRuntimePolicy struct{}

func newWithTargetScenariosOrdersRuntimePolicy() withTargetScenariosOrdersRuntimePolicy {
	return withTargetScenariosOrdersRuntimePolicy{}
}
func (withTargetScenariosOrdersRuntimePolicy) ShouldRunScenario(string) loadstrike.LoadStrikeBoolTask {
	return loadstrike.TaskFromBool(true)
}
func (withTargetScenariosOrdersRuntimePolicy) BeforeScenario(string) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withTargetScenariosOrdersRuntimePolicy) AfterScenario(string, loadstrike.LoadStrikeScenarioRuntime) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withTargetScenariosOrdersRuntimePolicy) BeforeStep(string, string) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withTargetScenariosOrdersRuntimePolicy) AfterStep(string, string, loadstrike.LoadStrikeReply) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}

type withTargetScenariosOrdersWorkerPlugin struct{}

func newWithTargetScenariosOrdersWorkerPlugin() withTargetScenariosOrdersWorkerPlugin {
	return withTargetScenariosOrdersWorkerPlugin{}
}
func (withTargetScenariosOrdersWorkerPlugin) PluginName() string {
	return "orders-sample-plugin"
}
func (withTargetScenariosOrdersWorkerPlugin) Init(loadstrike.LoadStrikeBaseContext, loadstrike.IConfiguration) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withTargetScenariosOrdersWorkerPlugin) Start(loadstrike.LoadStrikeSessionStartInfo) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withTargetScenariosOrdersWorkerPlugin) GetData(loadstrike.LoadStrikeRunResult) loadstrike.LoadStrikeValueTask[loadstrike.LoadStrikePluginData] {
	return loadstrike.TaskFromResult(loadstrike.CreatePluginData("orders-sample-plugin"))
}
func (withTargetScenariosOrdersWorkerPlugin) Stop() loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withTargetScenariosOrdersWorkerPlugin) Dispose() loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}

func withTargetScenariosPerformOrderGetReply() loadstrike.LoadStrikeReply {
	return loadstrike.LoadStrikeResponse.Ok("200", int64(128), "ok", loadstrike.TimeSpan(3*time.Millisecond))
}

func withTargetScenariosExecuteOrderGet(context loadstrike.LoadStrikeScenarioContext) loadstrike.LoadStrikeReply {
	return loadstrike.LoadStrikeStep.Run("get-order", context, func(loadstrike.LoadStrikeScenarioContext) loadstrike.LoadStrikeReply {
		return withTargetScenariosPerformOrderGetReply()
	})
}

func withTargetScenariosExecuteOrderGetAsync(context loadstrike.LoadStrikeScenarioContext) loadstrike.LoadStrikeValueTask[loadstrike.LoadStrikeReply] {
	return loadstrike.TaskFromResult(withTargetScenariosExecuteOrderGet(context))
}

func withTargetScenariosBaselineScenario(names ...string) loadstrike.LoadStrikeScenario {
	name := "orders.get-by-id"
	if len(names) > 0 && names[0] != "" {
		name = names[0]
	}
	return loadstrike.CreateScenario(name, withTargetScenariosExecuteOrderGet).
		WithLoadSimulations(loadstrike.LoadStrikeSimulation.IterationsForConstant(1, 1)).
		WithoutWarmUp()
}

func withTargetScenariosBaseRunner() loadstrike.LoadStrikeRunner {
	return loadstrike.Create().
		AddScenario(withTargetScenariosBaselineScenario()).
		WithRunnerKey(withTargetScenariosRunnerKey).
		WithTestSuite("orders-reference").
		WithTestName("orders-get-by-id").
		WithoutReports()
}

func withTargetScenariosBaseContext() loadstrike.LoadStrikeContext {
	return withTargetScenariosBaseRunner().BuildContext()
}

func withTargetScenariosHttpSource() *loadstrike.EndpointSpec {
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

func withTargetScenariosHttpDestination() *loadstrike.EndpointSpec {
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

func withTargetScenariosTrackingConfiguration() *loadstrike.TrackingConfigurationSpec {
	return &loadstrike.TrackingConfigurationSpec{
		Source:                      withTargetScenariosHttpSource(),
		Destination:                 withTargetScenariosHttpDestination(),
		RunMode:                     "GenerateAndCorrelate",
		CorrelationTimeoutSeconds:   30,
		TimeoutSweepIntervalSeconds: 1,
		TimeoutBatchSize:            200,
		TimeoutCountsAsFailure:      true,
		MetricPrefix:                "orders_tracking",
		ExecuteOriginalScenarioRun:  false,
	}
}

func withTargetScenariosTrackedContext() loadstrike.LoadStrikeContext {
	return loadstrike.Create().
		AddScenario(withTargetScenariosBaselineScenario("orders.tracked").WithCrossPlatformTracking(withTargetScenariosTrackingConfiguration())).
		WithRunnerKey(withTargetScenariosRunnerKey).
		WithoutReports().
		BuildContext()
}

func withTargetScenariosBuildTrackingPayload() loadstrike.TrackingPayload {
	builder := loadstrike.TrackingPayloadBuilder{
		Headers:            map[string]string{"X-Correlation-Id": "ord-1"},
		ContentType:        "application/json",
		MessagePayloadType: "json",
	}
	builder.SetBody(map[string]any{"trackingId": "ord-1"})
	return builder.Build()
}

func withTargetScenariosRunResult() loadstrike.LoadStrikeRunResult {
	var result loadstrike.LoadStrikeRunResult
	return result
}

func withTargetScenariosScenarioStats() loadstrike.LoadStrikeScenarioStats {
	var stats loadstrike.LoadStrikeScenarioStats
	return stats
}

func withTargetScenariosWriteTempConfigFiles() withTargetScenariosTempConfigPaths {
	return withTargetScenariosTempConfigPaths{
		ConfigPath: "method-reference.loadstrike.config.json",
		InfraPath:  "method-reference.loadstrike.infra.json",
	}
}

        // Target a single scenario by name.
func (reference WithTargetScenariosMethodReference) SelectOneScenarioExample() any {
    return withTargetScenariosBaseRunner().WithTargetScenarios("orders.get-by-id")
}

// Target more than one named scenario in the same call.
func (reference WithTargetScenariosMethodReference) SelectMultipleScenariosExample() any {
    return withTargetScenariosBaseRunner().WithTargetScenarios("orders.get-by-id", "orders.audit")
}
