package loadstrike_step

        import (
        	"time"

        	loadstrike "loadstrike.com/sdk/go"
        )

        const runAsyncRunnerKey = "runner_dummy_orders_reference"

type RunAsyncMethodReference struct{}

type runAsyncTempConfigPaths struct {
	ConfigPath string
	InfraPath  string
}

type runAsyncOrdersReportingSink struct{}

func newRunAsyncOrdersReportingSink() runAsyncOrdersReportingSink {
	return runAsyncOrdersReportingSink{}
}
func (runAsyncOrdersReportingSink) SinkName() string {
	return "orders-sample-sink"
}
func (runAsyncOrdersReportingSink) Init(loadstrike.LoadStrikeBaseContext, loadstrike.IConfiguration) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (runAsyncOrdersReportingSink) Start(loadstrike.LoadStrikeSessionStartInfo) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (runAsyncOrdersReportingSink) SaveRealtimeStats([]loadstrike.LoadStrikeScenarioStats) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (runAsyncOrdersReportingSink) SaveRealtimeMetrics(loadstrike.LoadStrikeMetricStats) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (runAsyncOrdersReportingSink) SaveRunResult(loadstrike.LoadStrikeRunResult) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (runAsyncOrdersReportingSink) Stop() loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (runAsyncOrdersReportingSink) Dispose() {}

type runAsyncOrdersRuntimePolicy struct{}

func newRunAsyncOrdersRuntimePolicy() runAsyncOrdersRuntimePolicy {
	return runAsyncOrdersRuntimePolicy{}
}
func (runAsyncOrdersRuntimePolicy) ShouldRunScenario(string) loadstrike.LoadStrikeBoolTask {
	return loadstrike.TaskFromBool(true)
}
func (runAsyncOrdersRuntimePolicy) BeforeScenario(string) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (runAsyncOrdersRuntimePolicy) AfterScenario(string, loadstrike.LoadStrikeScenarioRuntime) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (runAsyncOrdersRuntimePolicy) BeforeStep(string, string) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (runAsyncOrdersRuntimePolicy) AfterStep(string, string, loadstrike.LoadStrikeReply) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}

type runAsyncOrdersWorkerPlugin struct{}

func newRunAsyncOrdersWorkerPlugin() runAsyncOrdersWorkerPlugin {
	return runAsyncOrdersWorkerPlugin{}
}
func (runAsyncOrdersWorkerPlugin) PluginName() string {
	return "orders-sample-plugin"
}
func (runAsyncOrdersWorkerPlugin) Init(loadstrike.LoadStrikeBaseContext, loadstrike.IConfiguration) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (runAsyncOrdersWorkerPlugin) Start(loadstrike.LoadStrikeSessionStartInfo) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (runAsyncOrdersWorkerPlugin) GetData(loadstrike.LoadStrikeRunResult) loadstrike.LoadStrikeValueTask[loadstrike.LoadStrikePluginData] {
	return loadstrike.TaskFromResult(loadstrike.CreatePluginData("orders-sample-plugin"))
}
func (runAsyncOrdersWorkerPlugin) Stop() loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (runAsyncOrdersWorkerPlugin) Dispose() loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}

func runAsyncPerformOrderGetReply() loadstrike.LoadStrikeReply {
	return loadstrike.LoadStrikeResponse.Ok("200", int64(128), "ok", loadstrike.TimeSpan(3*time.Millisecond))
}

func runAsyncExecuteOrderGet(context loadstrike.LoadStrikeScenarioContext) loadstrike.LoadStrikeReply {
	return loadstrike.LoadStrikeStep.Run("get-order", context, func(loadstrike.LoadStrikeScenarioContext) loadstrike.LoadStrikeReply {
		return runAsyncPerformOrderGetReply()
	})
}

func runAsyncExecuteOrderGetAsync(context loadstrike.LoadStrikeScenarioContext) loadstrike.LoadStrikeValueTask[loadstrike.LoadStrikeReply] {
	return loadstrike.TaskFromResult(runAsyncExecuteOrderGet(context))
}

func runAsyncBaselineScenario(names ...string) loadstrike.LoadStrikeScenario {
	name := "orders.get-by-id"
	if len(names) > 0 && names[0] != "" {
		name = names[0]
	}
	return loadstrike.CreateScenario(name, runAsyncExecuteOrderGet).
		WithLoadSimulations(loadstrike.LoadStrikeSimulation.IterationsForConstant(1, 1)).
		WithoutWarmUp()
}

func runAsyncBaseRunner() loadstrike.LoadStrikeRunner {
	return loadstrike.Create().
		AddScenario(runAsyncBaselineScenario()).
		WithRunnerKey(runAsyncRunnerKey).
		WithTestSuite("orders-reference").
		WithTestName("orders-get-by-id").
		WithoutReports()
}

func runAsyncBaseContext() loadstrike.LoadStrikeContext {
	return runAsyncBaseRunner().BuildContext()
}

func runAsyncHttpSource() *loadstrike.EndpointSpec {
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

func runAsyncHttpDestination() *loadstrike.EndpointSpec {
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

func runAsyncTrackingConfiguration() *loadstrike.TrackingConfigurationSpec {
	return &loadstrike.TrackingConfigurationSpec{
		Source:                      runAsyncHttpSource(),
		Destination:                 runAsyncHttpDestination(),
		RunMode:                     "GenerateAndCorrelate",
		CorrelationTimeoutSeconds:   30,
		TimeoutSweepIntervalSeconds: 1,
		TimeoutBatchSize:            200,
		TimeoutCountsAsFailure:      true,
		MetricPrefix:                "orders_tracking",
		ExecuteOriginalScenarioRun:  false,
	}
}

func runAsyncTrackedContext() loadstrike.LoadStrikeContext {
	return loadstrike.Create().
		AddScenario(runAsyncBaselineScenario("orders.tracked").WithCrossPlatformTracking(runAsyncTrackingConfiguration())).
		WithRunnerKey(runAsyncRunnerKey).
		WithoutReports().
		BuildContext()
}

func runAsyncBuildTrackingPayload() loadstrike.TrackingPayload {
	builder := loadstrike.TrackingPayloadBuilder{
		Headers:            map[string]string{"X-Correlation-Id": "ord-1"},
		ContentType:        "application/json",
		MessagePayloadType: "json",
	}
	builder.SetBody(map[string]any{"trackingId": "ord-1"})
	return builder.Build()
}

func runAsyncRunResult() loadstrike.LoadStrikeRunResult {
	var result loadstrike.LoadStrikeRunResult
	return result
}

func runAsyncScenarioStats() loadstrike.LoadStrikeScenarioStats {
	var stats loadstrike.LoadStrikeScenarioStats
	return stats
}

func runAsyncWriteTempConfigFiles() runAsyncTempConfigPaths {
	return runAsyncTempConfigPaths{
		ConfigPath: "method-reference.loadstrike.config.json",
		InfraPath:  "method-reference.loadstrike.infra.json",
	}
}

        // Execute the GET request through the explicit async step helper.
func (reference RunAsyncMethodReference) RunAsyncStepExample() any {
    return loadstrike.LoadStrikeStep.RunAsync("get-order", createStepContext(), func(loadstrike.LoadStrikeScenarioContext) loadstrike.LoadStrikeValueTask[loadstrike.LoadStrikeReply] { return loadstrike.TaskFromResult(runAsyncPerformOrderGetReply()) })
}

// Project the async step reply to the untyped reply contract.
func (reference RunAsyncMethodReference) RunAsyncStepAndProjectReplyExample() any {
    return loadstrike.LoadStrikeStep.RunAsync("get-order", createStepContext(), func(loadstrike.LoadStrikeScenarioContext) loadstrike.LoadStrikeValueTask[loadstrike.LoadStrikeReply] { return loadstrike.TaskFromResult(runAsyncPerformOrderGetReply()) })
}
