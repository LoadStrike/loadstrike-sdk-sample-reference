package loadstrike_response

        import (
        	"time"

        	loadstrike "loadstrike.com/sdk/go"
        )

        const failWithPayloadRunnerKey = "runner_dummy_orders_reference"

type FailWithPayloadMethodReference struct{}

type failWithPayloadTempConfigPaths struct {
	ConfigPath string
	InfraPath  string
}

type failWithPayloadOrdersReportingSink struct{}

func newFailWithPayloadOrdersReportingSink() failWithPayloadOrdersReportingSink {
	return failWithPayloadOrdersReportingSink{}
}
func (failWithPayloadOrdersReportingSink) SinkName() string {
	return "orders-sample-sink"
}
func (failWithPayloadOrdersReportingSink) Init(loadstrike.LoadStrikeBaseContext, loadstrike.IConfiguration) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (failWithPayloadOrdersReportingSink) Start(loadstrike.LoadStrikeSessionStartInfo) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (failWithPayloadOrdersReportingSink) SaveRealtimeStats([]loadstrike.LoadStrikeScenarioStats) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (failWithPayloadOrdersReportingSink) SaveRealtimeMetrics(loadstrike.LoadStrikeMetricStats) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (failWithPayloadOrdersReportingSink) SaveRunResult(loadstrike.LoadStrikeRunResult) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (failWithPayloadOrdersReportingSink) Stop() loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (failWithPayloadOrdersReportingSink) Dispose() {}

type failWithPayloadOrdersRuntimePolicy struct{}

func newFailWithPayloadOrdersRuntimePolicy() failWithPayloadOrdersRuntimePolicy {
	return failWithPayloadOrdersRuntimePolicy{}
}
func (failWithPayloadOrdersRuntimePolicy) ShouldRunScenario(string) loadstrike.LoadStrikeBoolTask {
	return loadstrike.TaskFromBool(true)
}
func (failWithPayloadOrdersRuntimePolicy) BeforeScenario(string) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (failWithPayloadOrdersRuntimePolicy) AfterScenario(string, loadstrike.LoadStrikeScenarioRuntime) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (failWithPayloadOrdersRuntimePolicy) BeforeStep(string, string) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (failWithPayloadOrdersRuntimePolicy) AfterStep(string, string, loadstrike.LoadStrikeReply) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}

type failWithPayloadOrdersWorkerPlugin struct{}

func newFailWithPayloadOrdersWorkerPlugin() failWithPayloadOrdersWorkerPlugin {
	return failWithPayloadOrdersWorkerPlugin{}
}
func (failWithPayloadOrdersWorkerPlugin) PluginName() string {
	return "orders-sample-plugin"
}
func (failWithPayloadOrdersWorkerPlugin) Init(loadstrike.LoadStrikeBaseContext, loadstrike.IConfiguration) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (failWithPayloadOrdersWorkerPlugin) Start(loadstrike.LoadStrikeSessionStartInfo) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (failWithPayloadOrdersWorkerPlugin) GetData(loadstrike.LoadStrikeRunResult) loadstrike.LoadStrikeValueTask[loadstrike.LoadStrikePluginData] {
	return loadstrike.TaskFromResult(loadstrike.CreatePluginData("orders-sample-plugin"))
}
func (failWithPayloadOrdersWorkerPlugin) Stop() loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (failWithPayloadOrdersWorkerPlugin) Dispose() loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}

func failWithPayloadPerformOrderGetReply() loadstrike.LoadStrikeReply {
	return loadstrike.LoadStrikeResponse.Ok("200", int64(128), "ok", loadstrike.TimeSpan(3*time.Millisecond))
}

func failWithPayloadExecuteOrderGet(context loadstrike.LoadStrikeScenarioContext) loadstrike.LoadStrikeReply {
	return loadstrike.LoadStrikeStep.Run("get-order", context, func(loadstrike.LoadStrikeScenarioContext) loadstrike.LoadStrikeReply {
		return failWithPayloadPerformOrderGetReply()
	})
}

func failWithPayloadExecuteOrderGetAsync(context loadstrike.LoadStrikeScenarioContext) loadstrike.LoadStrikeValueTask[loadstrike.LoadStrikeReply] {
	return loadstrike.TaskFromResult(failWithPayloadExecuteOrderGet(context))
}

func failWithPayloadBaselineScenario(names ...string) loadstrike.LoadStrikeScenario {
	name := "orders.get-by-id"
	if len(names) > 0 && names[0] != "" {
		name = names[0]
	}
	return loadstrike.CreateScenario(name, failWithPayloadExecuteOrderGet).
		WithLoadSimulations(loadstrike.LoadStrikeSimulation.IterationsForConstant(1, 1)).
		WithoutWarmUp()
}

func failWithPayloadBaseRunner() loadstrike.LoadStrikeRunner {
	return loadstrike.Create().
		AddScenario(failWithPayloadBaselineScenario()).
		WithRunnerKey(failWithPayloadRunnerKey).
		WithTestSuite("orders-reference").
		WithTestName("orders-get-by-id").
		WithoutReports()
}

func failWithPayloadBaseContext() loadstrike.LoadStrikeContext {
	return failWithPayloadBaseRunner().BuildContext()
}

func failWithPayloadHttpSource() *loadstrike.EndpointSpec {
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

func failWithPayloadHttpDestination() *loadstrike.EndpointSpec {
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

func failWithPayloadTrackingConfiguration() *loadstrike.TrackingConfigurationSpec {
	return &loadstrike.TrackingConfigurationSpec{
		Source:                      failWithPayloadHttpSource(),
		Destination:                 failWithPayloadHttpDestination(),
		RunMode:                     "GenerateAndCorrelate",
		CorrelationTimeoutSeconds:   30,
		TimeoutSweepIntervalSeconds: 1,
		TimeoutBatchSize:            200,
		TimeoutCountsAsFailure:      true,
		MetricPrefix:                "orders_tracking",
		ExecuteOriginalScenarioRun:  false,
	}
}

func failWithPayloadTrackedContext() loadstrike.LoadStrikeContext {
	return loadstrike.Create().
		AddScenario(failWithPayloadBaselineScenario("orders.tracked").WithCrossPlatformTracking(failWithPayloadTrackingConfiguration())).
		WithRunnerKey(failWithPayloadRunnerKey).
		WithoutReports().
		BuildContext()
}

func failWithPayloadBuildTrackingPayload() loadstrike.TrackingPayload {
	builder := loadstrike.TrackingPayloadBuilder{
		Headers:            map[string]string{"X-Correlation-Id": "ord-1"},
		ContentType:        "application/json",
		MessagePayloadType: "json",
	}
	builder.SetBody(map[string]any{"trackingId": "ord-1"})
	return builder.Build()
}

func failWithPayloadRunResult() loadstrike.LoadStrikeRunResult {
	var result loadstrike.LoadStrikeRunResult
	return result
}

func failWithPayloadScenarioStats() loadstrike.LoadStrikeScenarioStats {
	var stats loadstrike.LoadStrikeScenarioStats
	return stats
}

func failWithPayloadWriteTempConfigFiles() failWithPayloadTempConfigPaths {
	return failWithPayloadTempConfigPaths{
		ConfigPath: "method-reference.loadstrike.config.json",
		InfraPath:  "method-reference.loadstrike.infra.json",
	}
}

        // Create the primary reply shape for this helper.
func (reference FailWithPayloadMethodReference) CreatePrimaryReplyExample() any {
    return loadstrike.LoadStrikeResponse.FailWith(map[string]any{"error": "not-found"}, "409", "conflict", int64(64), loadstrike.TimeSpan(3 * time.Millisecond))
}

// Create an alternate reply shape that changes status or message metadata.
func (reference FailWithPayloadMethodReference) CreateAlternateReplyExample() any {
    return loadstrike.LoadStrikeResponse.FailWith(map[string]any{"error": "timeout"}, "504", "retry", int64(8), loadstrike.TimeSpan(4 * time.Millisecond))
}
