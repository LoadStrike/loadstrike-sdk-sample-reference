package loadstrike_runner

        import (
        	"time"

        	loadstrike "loadstrike.com/sdk/go"
        )

        const withRuntimePolicyErrorModeRunnerKey = "runner_dummy_orders_reference"

type WithRuntimePolicyErrorModeMethodReference struct{}

type withRuntimePolicyErrorModeTempConfigPaths struct {
	ConfigPath string
	InfraPath  string
}

type withRuntimePolicyErrorModeOrdersReportingSink struct{}

func newWithRuntimePolicyErrorModeOrdersReportingSink() withRuntimePolicyErrorModeOrdersReportingSink {
	return withRuntimePolicyErrorModeOrdersReportingSink{}
}
func (withRuntimePolicyErrorModeOrdersReportingSink) SinkName() string {
	return "orders-sample-sink"
}
func (withRuntimePolicyErrorModeOrdersReportingSink) Init(loadstrike.LoadStrikeBaseContext, loadstrike.IConfiguration) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withRuntimePolicyErrorModeOrdersReportingSink) Start(loadstrike.LoadStrikeSessionStartInfo) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withRuntimePolicyErrorModeOrdersReportingSink) SaveRealtimeStats([]loadstrike.LoadStrikeScenarioStats) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withRuntimePolicyErrorModeOrdersReportingSink) SaveRealtimeMetrics(loadstrike.LoadStrikeMetricStats) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withRuntimePolicyErrorModeOrdersReportingSink) SaveRunResult(loadstrike.LoadStrikeRunResult) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withRuntimePolicyErrorModeOrdersReportingSink) Stop() loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withRuntimePolicyErrorModeOrdersReportingSink) Dispose() {}

type withRuntimePolicyErrorModeOrdersRuntimePolicy struct{}

func newWithRuntimePolicyErrorModeOrdersRuntimePolicy() withRuntimePolicyErrorModeOrdersRuntimePolicy {
	return withRuntimePolicyErrorModeOrdersRuntimePolicy{}
}
func (withRuntimePolicyErrorModeOrdersRuntimePolicy) ShouldRunScenario(string) loadstrike.LoadStrikeBoolTask {
	return loadstrike.TaskFromBool(true)
}
func (withRuntimePolicyErrorModeOrdersRuntimePolicy) BeforeScenario(string) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withRuntimePolicyErrorModeOrdersRuntimePolicy) AfterScenario(string, loadstrike.LoadStrikeScenarioRuntime) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withRuntimePolicyErrorModeOrdersRuntimePolicy) BeforeStep(string, string) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withRuntimePolicyErrorModeOrdersRuntimePolicy) AfterStep(string, string, loadstrike.LoadStrikeReply) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}

type withRuntimePolicyErrorModeOrdersWorkerPlugin struct{}

func newWithRuntimePolicyErrorModeOrdersWorkerPlugin() withRuntimePolicyErrorModeOrdersWorkerPlugin {
	return withRuntimePolicyErrorModeOrdersWorkerPlugin{}
}
func (withRuntimePolicyErrorModeOrdersWorkerPlugin) PluginName() string {
	return "orders-sample-plugin"
}
func (withRuntimePolicyErrorModeOrdersWorkerPlugin) Init(loadstrike.LoadStrikeBaseContext, loadstrike.IConfiguration) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withRuntimePolicyErrorModeOrdersWorkerPlugin) Start(loadstrike.LoadStrikeSessionStartInfo) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withRuntimePolicyErrorModeOrdersWorkerPlugin) GetData(loadstrike.LoadStrikeRunResult) loadstrike.LoadStrikeValueTask[loadstrike.LoadStrikePluginData] {
	return loadstrike.TaskFromResult(loadstrike.CreatePluginData("orders-sample-plugin"))
}
func (withRuntimePolicyErrorModeOrdersWorkerPlugin) Stop() loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withRuntimePolicyErrorModeOrdersWorkerPlugin) Dispose() loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}

func withRuntimePolicyErrorModePerformOrderGetReply() loadstrike.LoadStrikeReply {
	return loadstrike.LoadStrikeResponse.Ok("200", int64(128), "ok", loadstrike.TimeSpan(3*time.Millisecond))
}

func withRuntimePolicyErrorModeExecuteOrderGet(context loadstrike.LoadStrikeScenarioContext) loadstrike.LoadStrikeReply {
	return loadstrike.LoadStrikeStep.Run("get-order", context, func(loadstrike.LoadStrikeScenarioContext) loadstrike.LoadStrikeReply {
		return withRuntimePolicyErrorModePerformOrderGetReply()
	})
}

func withRuntimePolicyErrorModeExecuteOrderGetAsync(context loadstrike.LoadStrikeScenarioContext) loadstrike.LoadStrikeValueTask[loadstrike.LoadStrikeReply] {
	return loadstrike.TaskFromResult(withRuntimePolicyErrorModeExecuteOrderGet(context))
}

func withRuntimePolicyErrorModeBaselineScenario(names ...string) loadstrike.LoadStrikeScenario {
	name := "orders.get-by-id"
	if len(names) > 0 && names[0] != "" {
		name = names[0]
	}
	return loadstrike.CreateScenario(name, withRuntimePolicyErrorModeExecuteOrderGet).
		WithLoadSimulations(loadstrike.LoadStrikeSimulation.IterationsForConstant(1, 1)).
		WithoutWarmUp()
}

func withRuntimePolicyErrorModeBaseRunner() loadstrike.LoadStrikeRunner {
	return loadstrike.Create().
		AddScenario(withRuntimePolicyErrorModeBaselineScenario()).
		WithRunnerKey(withRuntimePolicyErrorModeRunnerKey).
		WithTestSuite("orders-reference").
		WithTestName("orders-get-by-id").
		WithoutReports()
}

func withRuntimePolicyErrorModeBaseContext() loadstrike.LoadStrikeContext {
	return withRuntimePolicyErrorModeBaseRunner().BuildContext()
}

func withRuntimePolicyErrorModeHttpSource() *loadstrike.EndpointSpec {
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

func withRuntimePolicyErrorModeHttpDestination() *loadstrike.EndpointSpec {
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

func withRuntimePolicyErrorModeTrackingConfiguration() *loadstrike.TrackingConfigurationSpec {
	return &loadstrike.TrackingConfigurationSpec{
		Source:                      withRuntimePolicyErrorModeHttpSource(),
		Destination:                 withRuntimePolicyErrorModeHttpDestination(),
		RunMode:                     "GenerateAndCorrelate",
		CorrelationTimeoutSeconds:   30,
		TimeoutSweepIntervalSeconds: 1,
		TimeoutBatchSize:            200,
		TimeoutCountsAsFailure:      true,
		MetricPrefix:                "orders_tracking",
		ExecuteOriginalScenarioRun:  false,
	}
}

func withRuntimePolicyErrorModeTrackedContext() loadstrike.LoadStrikeContext {
	return loadstrike.Create().
		AddScenario(withRuntimePolicyErrorModeBaselineScenario("orders.tracked").WithCrossPlatformTracking(withRuntimePolicyErrorModeTrackingConfiguration())).
		WithRunnerKey(withRuntimePolicyErrorModeRunnerKey).
		WithoutReports().
		BuildContext()
}

func withRuntimePolicyErrorModeBuildTrackingPayload() loadstrike.TrackingPayload {
	builder := loadstrike.TrackingPayloadBuilder{
		Headers:            map[string]string{"X-Correlation-Id": "ord-1"},
		ContentType:        "application/json",
		MessagePayloadType: "json",
	}
	builder.SetBody(map[string]any{"trackingId": "ord-1"})
	return builder.Build()
}

func withRuntimePolicyErrorModeRunResult() loadstrike.LoadStrikeRunResult {
	var result loadstrike.LoadStrikeRunResult
	return result
}

func withRuntimePolicyErrorModeScenarioStats() loadstrike.LoadStrikeScenarioStats {
	var stats loadstrike.LoadStrikeScenarioStats
	return stats
}

func withRuntimePolicyErrorModeWriteTempConfigFiles() withRuntimePolicyErrorModeTempConfigPaths {
	return withRuntimePolicyErrorModeTempConfigPaths{
		ConfigPath: "method-reference.loadstrike.config.json",
		InfraPath:  "method-reference.loadstrike.infra.json",
	}
}

        // Continue when policy hooks fail.
func (reference WithRuntimePolicyErrorModeMethodReference) ContinueOnPolicyErrorExample() any {
    return withRuntimePolicyErrorModeBaseRunner().WithRuntimePolicyErrorMode(loadstrike.RuntimePolicyErrorModeContinue)
}

// Fail the run when policy hooks fail.
func (reference WithRuntimePolicyErrorModeMethodReference) FailOnPolicyErrorExample() any {
    return withRuntimePolicyErrorModeBaseRunner().WithRuntimePolicyErrorMode(loadstrike.RuntimePolicyErrorModeFail)
}
