package loadstrike_response

        import (
        	"time"

        	loadstrike "loadstrike.com/sdk/go"
        )

        const okWithPayloadRunnerKey = "runner_dummy_orders_reference"

type OkWithPayloadMethodReference struct{}

type okWithPayloadTempConfigPaths struct {
	ConfigPath string
	InfraPath  string
}

type okWithPayloadOrdersReportingSink struct{}

func newOkWithPayloadOrdersReportingSink() okWithPayloadOrdersReportingSink {
	return okWithPayloadOrdersReportingSink{}
}
func (okWithPayloadOrdersReportingSink) SinkName() string {
	return "orders-sample-sink"
}
func (okWithPayloadOrdersReportingSink) Init(loadstrike.LoadStrikeBaseContext, loadstrike.IConfiguration) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (okWithPayloadOrdersReportingSink) Start(loadstrike.LoadStrikeSessionStartInfo) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (okWithPayloadOrdersReportingSink) SaveRealtimeStats([]loadstrike.LoadStrikeScenarioStats) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (okWithPayloadOrdersReportingSink) SaveRealtimeMetrics(loadstrike.LoadStrikeMetricStats) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (okWithPayloadOrdersReportingSink) SaveRunResult(loadstrike.LoadStrikeRunResult) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (okWithPayloadOrdersReportingSink) Stop() loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (okWithPayloadOrdersReportingSink) Dispose() {}

type okWithPayloadOrdersRuntimePolicy struct{}

func newOkWithPayloadOrdersRuntimePolicy() okWithPayloadOrdersRuntimePolicy {
	return okWithPayloadOrdersRuntimePolicy{}
}
func (okWithPayloadOrdersRuntimePolicy) ShouldRunScenario(string) loadstrike.LoadStrikeBoolTask {
	return loadstrike.TaskFromBool(true)
}
func (okWithPayloadOrdersRuntimePolicy) BeforeScenario(string) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (okWithPayloadOrdersRuntimePolicy) AfterScenario(string, loadstrike.LoadStrikeScenarioRuntime) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (okWithPayloadOrdersRuntimePolicy) BeforeStep(string, string) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (okWithPayloadOrdersRuntimePolicy) AfterStep(string, string, loadstrike.LoadStrikeReply) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}

type okWithPayloadOrdersWorkerPlugin struct{}

func newOkWithPayloadOrdersWorkerPlugin() okWithPayloadOrdersWorkerPlugin {
	return okWithPayloadOrdersWorkerPlugin{}
}
func (okWithPayloadOrdersWorkerPlugin) PluginName() string {
	return "orders-sample-plugin"
}
func (okWithPayloadOrdersWorkerPlugin) Init(loadstrike.LoadStrikeBaseContext, loadstrike.IConfiguration) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (okWithPayloadOrdersWorkerPlugin) Start(loadstrike.LoadStrikeSessionStartInfo) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (okWithPayloadOrdersWorkerPlugin) GetData(loadstrike.LoadStrikeRunResult) loadstrike.LoadStrikeValueTask[loadstrike.LoadStrikePluginData] {
	return loadstrike.TaskFromResult(loadstrike.CreatePluginData("orders-sample-plugin"))
}
func (okWithPayloadOrdersWorkerPlugin) Stop() loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (okWithPayloadOrdersWorkerPlugin) Dispose() loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}

func okWithPayloadPerformOrderGetReply() loadstrike.LoadStrikeReply {
	return loadstrike.LoadStrikeResponse.Ok("200", int64(128), "ok", loadstrike.TimeSpan(3*time.Millisecond))
}

func okWithPayloadExecuteOrderGet(context loadstrike.LoadStrikeScenarioContext) loadstrike.LoadStrikeReply {
	return loadstrike.LoadStrikeStep.Run("get-order", context, func(loadstrike.LoadStrikeScenarioContext) loadstrike.LoadStrikeReply {
		return okWithPayloadPerformOrderGetReply()
	})
}

func okWithPayloadExecuteOrderGetAsync(context loadstrike.LoadStrikeScenarioContext) loadstrike.LoadStrikeValueTask[loadstrike.LoadStrikeReply] {
	return loadstrike.TaskFromResult(okWithPayloadExecuteOrderGet(context))
}

func okWithPayloadBaselineScenario(names ...string) loadstrike.LoadStrikeScenario {
	name := "orders.get-by-id"
	if len(names) > 0 && names[0] != "" {
		name = names[0]
	}
	return loadstrike.CreateScenario(name, okWithPayloadExecuteOrderGet).
		WithLoadSimulations(loadstrike.LoadStrikeSimulation.IterationsForConstant(1, 1)).
		WithoutWarmUp()
}

func okWithPayloadBaseRunner() loadstrike.LoadStrikeRunner {
	return loadstrike.Create().
		AddScenario(okWithPayloadBaselineScenario()).
		WithRunnerKey(okWithPayloadRunnerKey).
		WithTestSuite("orders-reference").
		WithTestName("orders-get-by-id").
		WithoutReports()
}

func okWithPayloadBaseContext() loadstrike.LoadStrikeContext {
	return okWithPayloadBaseRunner().BuildContext()
}

func okWithPayloadHttpSource() *loadstrike.EndpointSpec {
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

func okWithPayloadHttpDestination() *loadstrike.EndpointSpec {
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

func okWithPayloadTrackingConfiguration() *loadstrike.TrackingConfigurationSpec {
	return &loadstrike.TrackingConfigurationSpec{
		Source:                      okWithPayloadHttpSource(),
		Destination:                 okWithPayloadHttpDestination(),
		RunMode:                     "GenerateAndCorrelate",
		CorrelationTimeoutSeconds:   30,
		TimeoutSweepIntervalSeconds: 1,
		TimeoutBatchSize:            200,
		TimeoutCountsAsFailure:      true,
		MetricPrefix:                "orders_tracking",
		ExecuteOriginalScenarioRun:  false,
	}
}

func okWithPayloadTrackedContext() loadstrike.LoadStrikeContext {
	return loadstrike.Create().
		AddScenario(okWithPayloadBaselineScenario("orders.tracked").WithCrossPlatformTracking(okWithPayloadTrackingConfiguration())).
		WithRunnerKey(okWithPayloadRunnerKey).
		WithoutReports().
		BuildContext()
}

func okWithPayloadBuildTrackingPayload() loadstrike.TrackingPayload {
	builder := loadstrike.TrackingPayloadBuilder{
		Headers:            map[string]string{"X-Correlation-Id": "ord-1"},
		ContentType:        "application/json",
		MessagePayloadType: "json",
	}
	builder.SetBody(map[string]any{"trackingId": "ord-1"})
	return builder.Build()
}

func okWithPayloadRunResult() loadstrike.LoadStrikeRunResult {
	var result loadstrike.LoadStrikeRunResult
	return result
}

func okWithPayloadScenarioStats() loadstrike.LoadStrikeScenarioStats {
	var stats loadstrike.LoadStrikeScenarioStats
	return stats
}

func okWithPayloadWriteTempConfigFiles() okWithPayloadTempConfigPaths {
	return okWithPayloadTempConfigPaths{
		ConfigPath: "method-reference.loadstrike.config.json",
		InfraPath:  "method-reference.loadstrike.infra.json",
	}
}

        // Create the primary reply shape for this helper.
func (reference OkWithPayloadMethodReference) CreatePrimaryReplyExample() any {
    return loadstrike.LoadStrikeResponse.OkWith(map[string]any{"orderId": "ORD-10001"}, "202", int64(128), "accepted", loadstrike.TimeSpan(3 * time.Millisecond))
}

// Create an alternate reply shape that changes status or message metadata.
func (reference OkWithPayloadMethodReference) CreateAlternateReplyExample() any {
    return loadstrike.LoadStrikeResponse.OkWith(map[string]any{"orderId": "ORD-10002"}, "201", int64(64), "created", loadstrike.TimeSpan(4 * time.Millisecond))
}
