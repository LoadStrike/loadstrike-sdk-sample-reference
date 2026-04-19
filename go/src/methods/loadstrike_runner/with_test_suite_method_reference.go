package loadstrike_runner

        import (
        	"time"

        	loadstrike "loadstrike.com/sdk/go"
        )

        const withTestSuiteRunnerKey = "runner_dummy_orders_reference"

type WithTestSuiteMethodReference struct{}

type withTestSuiteTempConfigPaths struct {
	ConfigPath string
	InfraPath  string
}

type withTestSuiteOrdersReportingSink struct{}

func newWithTestSuiteOrdersReportingSink() withTestSuiteOrdersReportingSink {
	return withTestSuiteOrdersReportingSink{}
}
func (withTestSuiteOrdersReportingSink) SinkName() string {
	return "orders-sample-sink"
}
func (withTestSuiteOrdersReportingSink) Init(loadstrike.LoadStrikeBaseContext, loadstrike.IConfiguration) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withTestSuiteOrdersReportingSink) Start(loadstrike.LoadStrikeSessionStartInfo) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withTestSuiteOrdersReportingSink) SaveRealtimeStats([]loadstrike.LoadStrikeScenarioStats) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withTestSuiteOrdersReportingSink) SaveRealtimeMetrics(loadstrike.LoadStrikeMetricStats) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withTestSuiteOrdersReportingSink) SaveRunResult(loadstrike.LoadStrikeRunResult) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withTestSuiteOrdersReportingSink) Stop() loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withTestSuiteOrdersReportingSink) Dispose() {}

type withTestSuiteOrdersRuntimePolicy struct{}

func newWithTestSuiteOrdersRuntimePolicy() withTestSuiteOrdersRuntimePolicy {
	return withTestSuiteOrdersRuntimePolicy{}
}
func (withTestSuiteOrdersRuntimePolicy) ShouldRunScenario(string) loadstrike.LoadStrikeBoolTask {
	return loadstrike.TaskFromBool(true)
}
func (withTestSuiteOrdersRuntimePolicy) BeforeScenario(string) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withTestSuiteOrdersRuntimePolicy) AfterScenario(string, loadstrike.LoadStrikeScenarioRuntime) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withTestSuiteOrdersRuntimePolicy) BeforeStep(string, string) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withTestSuiteOrdersRuntimePolicy) AfterStep(string, string, loadstrike.LoadStrikeReply) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}

type withTestSuiteOrdersWorkerPlugin struct{}

func newWithTestSuiteOrdersWorkerPlugin() withTestSuiteOrdersWorkerPlugin {
	return withTestSuiteOrdersWorkerPlugin{}
}
func (withTestSuiteOrdersWorkerPlugin) PluginName() string {
	return "orders-sample-plugin"
}
func (withTestSuiteOrdersWorkerPlugin) Init(loadstrike.LoadStrikeBaseContext, loadstrike.IConfiguration) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withTestSuiteOrdersWorkerPlugin) Start(loadstrike.LoadStrikeSessionStartInfo) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withTestSuiteOrdersWorkerPlugin) GetData(loadstrike.LoadStrikeRunResult) loadstrike.LoadStrikeValueTask[loadstrike.LoadStrikePluginData] {
	return loadstrike.TaskFromResult(loadstrike.CreatePluginData("orders-sample-plugin"))
}
func (withTestSuiteOrdersWorkerPlugin) Stop() loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withTestSuiteOrdersWorkerPlugin) Dispose() loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}

func withTestSuitePerformOrderGetReply() loadstrike.LoadStrikeReply {
	return loadstrike.LoadStrikeResponse.Ok("200", int64(128), "ok", loadstrike.TimeSpan(3*time.Millisecond))
}

func withTestSuiteExecuteOrderGet(context loadstrike.LoadStrikeScenarioContext) loadstrike.LoadStrikeReply {
	return loadstrike.LoadStrikeStep.Run("get-order", context, func(loadstrike.LoadStrikeScenarioContext) loadstrike.LoadStrikeReply {
		return withTestSuitePerformOrderGetReply()
	})
}

func withTestSuiteExecuteOrderGetAsync(context loadstrike.LoadStrikeScenarioContext) loadstrike.LoadStrikeValueTask[loadstrike.LoadStrikeReply] {
	return loadstrike.TaskFromResult(withTestSuiteExecuteOrderGet(context))
}

func withTestSuiteBaselineScenario(names ...string) loadstrike.LoadStrikeScenario {
	name := "orders.get-by-id"
	if len(names) > 0 && names[0] != "" {
		name = names[0]
	}
	return loadstrike.CreateScenario(name, withTestSuiteExecuteOrderGet).
		WithLoadSimulations(loadstrike.LoadStrikeSimulation.IterationsForConstant(1, 1)).
		WithoutWarmUp()
}

func withTestSuiteBaseRunner() loadstrike.LoadStrikeRunner {
	return loadstrike.Create().
		AddScenario(withTestSuiteBaselineScenario()).
		WithRunnerKey(withTestSuiteRunnerKey).
		WithTestSuite("orders-reference").
		WithTestName("orders-get-by-id").
		WithoutReports()
}

func withTestSuiteBaseContext() loadstrike.LoadStrikeContext {
	return withTestSuiteBaseRunner().BuildContext()
}

func withTestSuiteHttpSource() *loadstrike.EndpointSpec {
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

func withTestSuiteHttpDestination() *loadstrike.EndpointSpec {
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

func withTestSuiteTrackingConfiguration() *loadstrike.TrackingConfigurationSpec {
	return &loadstrike.TrackingConfigurationSpec{
		Source:                      withTestSuiteHttpSource(),
		Destination:                 withTestSuiteHttpDestination(),
		RunMode:                     "GenerateAndCorrelate",
		CorrelationTimeoutSeconds:   30,
		TimeoutSweepIntervalSeconds: 1,
		TimeoutBatchSize:            200,
		TimeoutCountsAsFailure:      true,
		MetricPrefix:                "orders_tracking",
		ExecuteOriginalScenarioRun:  false,
	}
}

func withTestSuiteTrackedContext() loadstrike.LoadStrikeContext {
	return loadstrike.Create().
		AddScenario(withTestSuiteBaselineScenario("orders.tracked").WithCrossPlatformTracking(withTestSuiteTrackingConfiguration())).
		WithRunnerKey(withTestSuiteRunnerKey).
		WithoutReports().
		BuildContext()
}

func withTestSuiteBuildTrackingPayload() loadstrike.TrackingPayload {
	builder := loadstrike.TrackingPayloadBuilder{
		Headers:            map[string]string{"X-Correlation-Id": "ord-1"},
		ContentType:        "application/json",
		MessagePayloadType: "json",
	}
	builder.SetBody(map[string]any{"trackingId": "ord-1"})
	return builder.Build()
}

func withTestSuiteRunResult() loadstrike.LoadStrikeRunResult {
	var result loadstrike.LoadStrikeRunResult
	return result
}

func withTestSuiteScenarioStats() loadstrike.LoadStrikeScenarioStats {
	var stats loadstrike.LoadStrikeScenarioStats
	return stats
}

func withTestSuiteWriteTempConfigFiles() withTestSuiteTempConfigPaths {
	return withTestSuiteTempConfigPaths{
		ConfigPath: "method-reference.loadstrike.config.json",
		InfraPath:  "method-reference.loadstrike.infra.json",
	}
}

        // Apply the primary string value shown in the sample method reference.
func (reference WithTestSuiteMethodReference) ApplyPrimaryValueExample() any {
    return withTestSuiteBaseRunner().WithTestSuite("orders-reference")
}

// Show the same method with a second concrete value.
func (reference WithTestSuiteMethodReference) ApplyAlternateValueExample() any {
    return withTestSuiteBaseRunner().WithTestSuite("checkout-reference")
}
