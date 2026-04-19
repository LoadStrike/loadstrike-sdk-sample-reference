package loadstrike_runner

        import (
        	"time"

        	loadstrike "loadstrike.com/sdk/go"
        )

        const withReportFormatsRunnerKey = "runner_dummy_orders_reference"

type WithReportFormatsMethodReference struct{}

type withReportFormatsTempConfigPaths struct {
	ConfigPath string
	InfraPath  string
}

type withReportFormatsOrdersReportingSink struct{}

func newWithReportFormatsOrdersReportingSink() withReportFormatsOrdersReportingSink {
	return withReportFormatsOrdersReportingSink{}
}
func (withReportFormatsOrdersReportingSink) SinkName() string {
	return "orders-sample-sink"
}
func (withReportFormatsOrdersReportingSink) Init(loadstrike.LoadStrikeBaseContext, loadstrike.IConfiguration) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withReportFormatsOrdersReportingSink) Start(loadstrike.LoadStrikeSessionStartInfo) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withReportFormatsOrdersReportingSink) SaveRealtimeStats([]loadstrike.LoadStrikeScenarioStats) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withReportFormatsOrdersReportingSink) SaveRealtimeMetrics(loadstrike.LoadStrikeMetricStats) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withReportFormatsOrdersReportingSink) SaveRunResult(loadstrike.LoadStrikeRunResult) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withReportFormatsOrdersReportingSink) Stop() loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withReportFormatsOrdersReportingSink) Dispose() {}

type withReportFormatsOrdersRuntimePolicy struct{}

func newWithReportFormatsOrdersRuntimePolicy() withReportFormatsOrdersRuntimePolicy {
	return withReportFormatsOrdersRuntimePolicy{}
}
func (withReportFormatsOrdersRuntimePolicy) ShouldRunScenario(string) loadstrike.LoadStrikeBoolTask {
	return loadstrike.TaskFromBool(true)
}
func (withReportFormatsOrdersRuntimePolicy) BeforeScenario(string) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withReportFormatsOrdersRuntimePolicy) AfterScenario(string, loadstrike.LoadStrikeScenarioRuntime) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withReportFormatsOrdersRuntimePolicy) BeforeStep(string, string) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withReportFormatsOrdersRuntimePolicy) AfterStep(string, string, loadstrike.LoadStrikeReply) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}

type withReportFormatsOrdersWorkerPlugin struct{}

func newWithReportFormatsOrdersWorkerPlugin() withReportFormatsOrdersWorkerPlugin {
	return withReportFormatsOrdersWorkerPlugin{}
}
func (withReportFormatsOrdersWorkerPlugin) PluginName() string {
	return "orders-sample-plugin"
}
func (withReportFormatsOrdersWorkerPlugin) Init(loadstrike.LoadStrikeBaseContext, loadstrike.IConfiguration) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withReportFormatsOrdersWorkerPlugin) Start(loadstrike.LoadStrikeSessionStartInfo) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withReportFormatsOrdersWorkerPlugin) GetData(loadstrike.LoadStrikeRunResult) loadstrike.LoadStrikeValueTask[loadstrike.LoadStrikePluginData] {
	return loadstrike.TaskFromResult(loadstrike.CreatePluginData("orders-sample-plugin"))
}
func (withReportFormatsOrdersWorkerPlugin) Stop() loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withReportFormatsOrdersWorkerPlugin) Dispose() loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}

func withReportFormatsPerformOrderGetReply() loadstrike.LoadStrikeReply {
	return loadstrike.LoadStrikeResponse.Ok("200", int64(128), "ok", loadstrike.TimeSpan(3*time.Millisecond))
}

func withReportFormatsExecuteOrderGet(context loadstrike.LoadStrikeScenarioContext) loadstrike.LoadStrikeReply {
	return loadstrike.LoadStrikeStep.Run("get-order", context, func(loadstrike.LoadStrikeScenarioContext) loadstrike.LoadStrikeReply {
		return withReportFormatsPerformOrderGetReply()
	})
}

func withReportFormatsExecuteOrderGetAsync(context loadstrike.LoadStrikeScenarioContext) loadstrike.LoadStrikeValueTask[loadstrike.LoadStrikeReply] {
	return loadstrike.TaskFromResult(withReportFormatsExecuteOrderGet(context))
}

func withReportFormatsBaselineScenario(names ...string) loadstrike.LoadStrikeScenario {
	name := "orders.get-by-id"
	if len(names) > 0 && names[0] != "" {
		name = names[0]
	}
	return loadstrike.CreateScenario(name, withReportFormatsExecuteOrderGet).
		WithLoadSimulations(loadstrike.LoadStrikeSimulation.IterationsForConstant(1, 1)).
		WithoutWarmUp()
}

func withReportFormatsBaseRunner() loadstrike.LoadStrikeRunner {
	return loadstrike.Create().
		AddScenario(withReportFormatsBaselineScenario()).
		WithRunnerKey(withReportFormatsRunnerKey).
		WithTestSuite("orders-reference").
		WithTestName("orders-get-by-id").
		WithoutReports()
}

func withReportFormatsBaseContext() loadstrike.LoadStrikeContext {
	return withReportFormatsBaseRunner().BuildContext()
}

func withReportFormatsHttpSource() *loadstrike.EndpointSpec {
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

func withReportFormatsHttpDestination() *loadstrike.EndpointSpec {
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

func withReportFormatsTrackingConfiguration() *loadstrike.TrackingConfigurationSpec {
	return &loadstrike.TrackingConfigurationSpec{
		Source:                      withReportFormatsHttpSource(),
		Destination:                 withReportFormatsHttpDestination(),
		RunMode:                     "GenerateAndCorrelate",
		CorrelationTimeoutSeconds:   30,
		TimeoutSweepIntervalSeconds: 1,
		TimeoutBatchSize:            200,
		TimeoutCountsAsFailure:      true,
		MetricPrefix:                "orders_tracking",
		ExecuteOriginalScenarioRun:  false,
	}
}

func withReportFormatsTrackedContext() loadstrike.LoadStrikeContext {
	return loadstrike.Create().
		AddScenario(withReportFormatsBaselineScenario("orders.tracked").WithCrossPlatformTracking(withReportFormatsTrackingConfiguration())).
		WithRunnerKey(withReportFormatsRunnerKey).
		WithoutReports().
		BuildContext()
}

func withReportFormatsBuildTrackingPayload() loadstrike.TrackingPayload {
	builder := loadstrike.TrackingPayloadBuilder{
		Headers:            map[string]string{"X-Correlation-Id": "ord-1"},
		ContentType:        "application/json",
		MessagePayloadType: "json",
	}
	builder.SetBody(map[string]any{"trackingId": "ord-1"})
	return builder.Build()
}

func withReportFormatsRunResult() loadstrike.LoadStrikeRunResult {
	var result loadstrike.LoadStrikeRunResult
	return result
}

func withReportFormatsScenarioStats() loadstrike.LoadStrikeScenarioStats {
	var stats loadstrike.LoadStrikeScenarioStats
	return stats
}

func withReportFormatsWriteTempConfigFiles() withReportFormatsTempConfigPaths {
	return withReportFormatsTempConfigPaths{
		ConfigPath: "method-reference.loadstrike.config.json",
		InfraPath:  "method-reference.loadstrike.infra.json",
	}
}

        // Write both text and HTML report outputs.
func (reference WithReportFormatsMethodReference) WriteTextAndHtmlReportsExample() any {
    return withReportFormatsBaseRunner().WithReportFormats(loadstrike.ReportFormatTXT, loadstrike.ReportFormatHTML)
}

// Switch to a single markdown-style report format.
func (reference WithReportFormatsMethodReference) WriteMarkdownOnlyReportExample() any {
    return withReportFormatsBaseRunner().WithReportFormats(loadstrike.ReportFormatMD)
}
