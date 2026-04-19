package loadstrike_runner

        import (
        	"time"

        	loadstrike "loadstrike.com/sdk/go"
        )

        const displayConsoleMetricsRunnerKey = "runner_dummy_orders_reference"

type DisplayConsoleMetricsMethodReference struct{}

type displayConsoleMetricsTempConfigPaths struct {
	ConfigPath string
	InfraPath  string
}

type displayConsoleMetricsOrdersReportingSink struct{}

func newDisplayConsoleMetricsOrdersReportingSink() displayConsoleMetricsOrdersReportingSink {
	return displayConsoleMetricsOrdersReportingSink{}
}
func (displayConsoleMetricsOrdersReportingSink) SinkName() string {
	return "orders-sample-sink"
}
func (displayConsoleMetricsOrdersReportingSink) Init(loadstrike.LoadStrikeBaseContext, loadstrike.IConfiguration) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (displayConsoleMetricsOrdersReportingSink) Start(loadstrike.LoadStrikeSessionStartInfo) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (displayConsoleMetricsOrdersReportingSink) SaveRealtimeStats([]loadstrike.LoadStrikeScenarioStats) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (displayConsoleMetricsOrdersReportingSink) SaveRealtimeMetrics(loadstrike.LoadStrikeMetricStats) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (displayConsoleMetricsOrdersReportingSink) SaveRunResult(loadstrike.LoadStrikeRunResult) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (displayConsoleMetricsOrdersReportingSink) Stop() loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (displayConsoleMetricsOrdersReportingSink) Dispose() {}

type displayConsoleMetricsOrdersRuntimePolicy struct{}

func newDisplayConsoleMetricsOrdersRuntimePolicy() displayConsoleMetricsOrdersRuntimePolicy {
	return displayConsoleMetricsOrdersRuntimePolicy{}
}
func (displayConsoleMetricsOrdersRuntimePolicy) ShouldRunScenario(string) loadstrike.LoadStrikeBoolTask {
	return loadstrike.TaskFromBool(true)
}
func (displayConsoleMetricsOrdersRuntimePolicy) BeforeScenario(string) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (displayConsoleMetricsOrdersRuntimePolicy) AfterScenario(string, loadstrike.LoadStrikeScenarioRuntime) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (displayConsoleMetricsOrdersRuntimePolicy) BeforeStep(string, string) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (displayConsoleMetricsOrdersRuntimePolicy) AfterStep(string, string, loadstrike.LoadStrikeReply) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}

type displayConsoleMetricsOrdersWorkerPlugin struct{}

func newDisplayConsoleMetricsOrdersWorkerPlugin() displayConsoleMetricsOrdersWorkerPlugin {
	return displayConsoleMetricsOrdersWorkerPlugin{}
}
func (displayConsoleMetricsOrdersWorkerPlugin) PluginName() string {
	return "orders-sample-plugin"
}
func (displayConsoleMetricsOrdersWorkerPlugin) Init(loadstrike.LoadStrikeBaseContext, loadstrike.IConfiguration) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (displayConsoleMetricsOrdersWorkerPlugin) Start(loadstrike.LoadStrikeSessionStartInfo) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (displayConsoleMetricsOrdersWorkerPlugin) GetData(loadstrike.LoadStrikeRunResult) loadstrike.LoadStrikeValueTask[loadstrike.LoadStrikePluginData] {
	return loadstrike.TaskFromResult(loadstrike.CreatePluginData("orders-sample-plugin"))
}
func (displayConsoleMetricsOrdersWorkerPlugin) Stop() loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (displayConsoleMetricsOrdersWorkerPlugin) Dispose() loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}

func displayConsoleMetricsPerformOrderGetReply() loadstrike.LoadStrikeReply {
	return loadstrike.LoadStrikeResponse.Ok("200", int64(128), "ok", loadstrike.TimeSpan(3*time.Millisecond))
}

func displayConsoleMetricsExecuteOrderGet(context loadstrike.LoadStrikeScenarioContext) loadstrike.LoadStrikeReply {
	return loadstrike.LoadStrikeStep.Run("get-order", context, func(loadstrike.LoadStrikeScenarioContext) loadstrike.LoadStrikeReply {
		return displayConsoleMetricsPerformOrderGetReply()
	})
}

func displayConsoleMetricsExecuteOrderGetAsync(context loadstrike.LoadStrikeScenarioContext) loadstrike.LoadStrikeValueTask[loadstrike.LoadStrikeReply] {
	return loadstrike.TaskFromResult(displayConsoleMetricsExecuteOrderGet(context))
}

func displayConsoleMetricsBaselineScenario(names ...string) loadstrike.LoadStrikeScenario {
	name := "orders.get-by-id"
	if len(names) > 0 && names[0] != "" {
		name = names[0]
	}
	return loadstrike.CreateScenario(name, displayConsoleMetricsExecuteOrderGet).
		WithLoadSimulations(loadstrike.LoadStrikeSimulation.IterationsForConstant(1, 1)).
		WithoutWarmUp()
}

func displayConsoleMetricsBaseRunner() loadstrike.LoadStrikeRunner {
	return loadstrike.Create().
		AddScenario(displayConsoleMetricsBaselineScenario()).
		WithRunnerKey(displayConsoleMetricsRunnerKey).
		WithTestSuite("orders-reference").
		WithTestName("orders-get-by-id").
		WithoutReports()
}

func displayConsoleMetricsBaseContext() loadstrike.LoadStrikeContext {
	return displayConsoleMetricsBaseRunner().BuildContext()
}

func displayConsoleMetricsHttpSource() *loadstrike.EndpointSpec {
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

func displayConsoleMetricsHttpDestination() *loadstrike.EndpointSpec {
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

func displayConsoleMetricsTrackingConfiguration() *loadstrike.TrackingConfigurationSpec {
	return &loadstrike.TrackingConfigurationSpec{
		Source:                      displayConsoleMetricsHttpSource(),
		Destination:                 displayConsoleMetricsHttpDestination(),
		RunMode:                     "GenerateAndCorrelate",
		CorrelationTimeoutSeconds:   30,
		TimeoutSweepIntervalSeconds: 1,
		TimeoutBatchSize:            200,
		TimeoutCountsAsFailure:      true,
		MetricPrefix:                "orders_tracking",
		ExecuteOriginalScenarioRun:  false,
	}
}

func displayConsoleMetricsTrackedContext() loadstrike.LoadStrikeContext {
	return loadstrike.Create().
		AddScenario(displayConsoleMetricsBaselineScenario("orders.tracked").WithCrossPlatformTracking(displayConsoleMetricsTrackingConfiguration())).
		WithRunnerKey(displayConsoleMetricsRunnerKey).
		WithoutReports().
		BuildContext()
}

func displayConsoleMetricsBuildTrackingPayload() loadstrike.TrackingPayload {
	builder := loadstrike.TrackingPayloadBuilder{
		Headers:            map[string]string{"X-Correlation-Id": "ord-1"},
		ContentType:        "application/json",
		MessagePayloadType: "json",
	}
	builder.SetBody(map[string]any{"trackingId": "ord-1"})
	return builder.Build()
}

func displayConsoleMetricsRunResult() loadstrike.LoadStrikeRunResult {
	var result loadstrike.LoadStrikeRunResult
	return result
}

func displayConsoleMetricsScenarioStats() loadstrike.LoadStrikeScenarioStats {
	var stats loadstrike.LoadStrikeScenarioStats
	return stats
}

func displayConsoleMetricsWriteTempConfigFiles() displayConsoleMetricsTempConfigPaths {
	return displayConsoleMetricsTempConfigPaths{
		ConfigPath: "method-reference.loadstrike.config.json",
		InfraPath:  "method-reference.loadstrike.infra.json",
	}
}

        // Apply the boolean option in its enabled state.
func (reference DisplayConsoleMetricsMethodReference) ToggleOnExample() any {
    return displayConsoleMetricsBaseRunner().WithDisplayConsoleMetrics(true)
}

// Apply the boolean option in its disabled state.
func (reference DisplayConsoleMetricsMethodReference) ToggleOffExample() any {
    return displayConsoleMetricsBaseRunner().WithDisplayConsoleMetrics(false)
}
