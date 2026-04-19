package loadstrike_reply

        import (
        	"time"

        	loadstrike "loadstrike.com/sdk/go"
        )

        const asReplyRunnerKey = "runner_dummy_orders_reference"

type AsReplyMethodReference struct{}

type asReplyTempConfigPaths struct {
	ConfigPath string
	InfraPath  string
}

type asReplyOrdersReportingSink struct{}

func newAsReplyOrdersReportingSink() asReplyOrdersReportingSink {
	return asReplyOrdersReportingSink{}
}
func (asReplyOrdersReportingSink) SinkName() string {
	return "orders-sample-sink"
}
func (asReplyOrdersReportingSink) Init(loadstrike.LoadStrikeBaseContext, loadstrike.IConfiguration) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (asReplyOrdersReportingSink) Start(loadstrike.LoadStrikeSessionStartInfo) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (asReplyOrdersReportingSink) SaveRealtimeStats([]loadstrike.LoadStrikeScenarioStats) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (asReplyOrdersReportingSink) SaveRealtimeMetrics(loadstrike.LoadStrikeMetricStats) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (asReplyOrdersReportingSink) SaveRunResult(loadstrike.LoadStrikeRunResult) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (asReplyOrdersReportingSink) Stop() loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (asReplyOrdersReportingSink) Dispose() {}

type asReplyOrdersRuntimePolicy struct{}

func newAsReplyOrdersRuntimePolicy() asReplyOrdersRuntimePolicy {
	return asReplyOrdersRuntimePolicy{}
}
func (asReplyOrdersRuntimePolicy) ShouldRunScenario(string) loadstrike.LoadStrikeBoolTask {
	return loadstrike.TaskFromBool(true)
}
func (asReplyOrdersRuntimePolicy) BeforeScenario(string) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (asReplyOrdersRuntimePolicy) AfterScenario(string, loadstrike.LoadStrikeScenarioRuntime) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (asReplyOrdersRuntimePolicy) BeforeStep(string, string) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (asReplyOrdersRuntimePolicy) AfterStep(string, string, loadstrike.LoadStrikeReply) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}

type asReplyOrdersWorkerPlugin struct{}

func newAsReplyOrdersWorkerPlugin() asReplyOrdersWorkerPlugin {
	return asReplyOrdersWorkerPlugin{}
}
func (asReplyOrdersWorkerPlugin) PluginName() string {
	return "orders-sample-plugin"
}
func (asReplyOrdersWorkerPlugin) Init(loadstrike.LoadStrikeBaseContext, loadstrike.IConfiguration) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (asReplyOrdersWorkerPlugin) Start(loadstrike.LoadStrikeSessionStartInfo) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (asReplyOrdersWorkerPlugin) GetData(loadstrike.LoadStrikeRunResult) loadstrike.LoadStrikeValueTask[loadstrike.LoadStrikePluginData] {
	return loadstrike.TaskFromResult(loadstrike.CreatePluginData("orders-sample-plugin"))
}
func (asReplyOrdersWorkerPlugin) Stop() loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (asReplyOrdersWorkerPlugin) Dispose() loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}

func asReplyPerformOrderGetReply() loadstrike.LoadStrikeReply {
	return loadstrike.LoadStrikeResponse.Ok("200", int64(128), "ok", loadstrike.TimeSpan(3*time.Millisecond))
}

func asReplyExecuteOrderGet(context loadstrike.LoadStrikeScenarioContext) loadstrike.LoadStrikeReply {
	return loadstrike.LoadStrikeStep.Run("get-order", context, func(loadstrike.LoadStrikeScenarioContext) loadstrike.LoadStrikeReply {
		return asReplyPerformOrderGetReply()
	})
}

func asReplyExecuteOrderGetAsync(context loadstrike.LoadStrikeScenarioContext) loadstrike.LoadStrikeValueTask[loadstrike.LoadStrikeReply] {
	return loadstrike.TaskFromResult(asReplyExecuteOrderGet(context))
}

func asReplyBaselineScenario(names ...string) loadstrike.LoadStrikeScenario {
	name := "orders.get-by-id"
	if len(names) > 0 && names[0] != "" {
		name = names[0]
	}
	return loadstrike.CreateScenario(name, asReplyExecuteOrderGet).
		WithLoadSimulations(loadstrike.LoadStrikeSimulation.IterationsForConstant(1, 1)).
		WithoutWarmUp()
}

func asReplyBaseRunner() loadstrike.LoadStrikeRunner {
	return loadstrike.Create().
		AddScenario(asReplyBaselineScenario()).
		WithRunnerKey(asReplyRunnerKey).
		WithTestSuite("orders-reference").
		WithTestName("orders-get-by-id").
		WithoutReports()
}

func asReplyBaseContext() loadstrike.LoadStrikeContext {
	return asReplyBaseRunner().BuildContext()
}

func asReplyHttpSource() *loadstrike.EndpointSpec {
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

func asReplyHttpDestination() *loadstrike.EndpointSpec {
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

func asReplyTrackingConfiguration() *loadstrike.TrackingConfigurationSpec {
	return &loadstrike.TrackingConfigurationSpec{
		Source:                      asReplyHttpSource(),
		Destination:                 asReplyHttpDestination(),
		RunMode:                     "GenerateAndCorrelate",
		CorrelationTimeoutSeconds:   30,
		TimeoutSweepIntervalSeconds: 1,
		TimeoutBatchSize:            200,
		TimeoutCountsAsFailure:      true,
		MetricPrefix:                "orders_tracking",
		ExecuteOriginalScenarioRun:  false,
	}
}

func asReplyTrackedContext() loadstrike.LoadStrikeContext {
	return loadstrike.Create().
		AddScenario(asReplyBaselineScenario("orders.tracked").WithCrossPlatformTracking(asReplyTrackingConfiguration())).
		WithRunnerKey(asReplyRunnerKey).
		WithoutReports().
		BuildContext()
}

func asReplyBuildTrackingPayload() loadstrike.TrackingPayload {
	builder := loadstrike.TrackingPayloadBuilder{
		Headers:            map[string]string{"X-Correlation-Id": "ord-1"},
		ContentType:        "application/json",
		MessagePayloadType: "json",
	}
	builder.SetBody(map[string]any{"trackingId": "ord-1"})
	return builder.Build()
}

func asReplyRunResult() loadstrike.LoadStrikeRunResult {
	var result loadstrike.LoadStrikeRunResult
	return result
}

func asReplyScenarioStats() loadstrike.LoadStrikeScenarioStats {
	var stats loadstrike.LoadStrikeScenarioStats
	return stats
}

func asReplyWriteTempConfigFiles() asReplyTempConfigPaths {
	return asReplyTempConfigPaths{
		ConfigPath: "method-reference.loadstrike.config.json",
		InfraPath:  "method-reference.loadstrike.infra.json",
	}
}

        // Project a typed reply into the untyped reply contract.
func (reference AsReplyMethodReference) ProjectTypedReplyExample() any {
    return loadstrike.LoadStrikeResponse.OkWith(map[string]any{"orderId": 42}, "200", int64(128), "ok", loadstrike.TimeSpan(3*time.Millisecond)).AsReply()
}

// Project a failed typed reply into the untyped reply contract.
func (reference AsReplyMethodReference) ProjectFailedTypedReplyExample() any {
    return loadstrike.LoadStrikeResponse.FailWith(map[string]any{"error": "not-found"}, "404", "missing", int64(64), loadstrike.TimeSpan(3*time.Millisecond)).AsReply()
}
