package loadstrike_scenario

        import (
        	"time"

        	loadstrike "loadstrike.com/sdk/go"
        )

        const withThresholdsRunnerKey = "runner_dummy_orders_reference"

type WithThresholdsMethodReference struct{}

type withThresholdsTempConfigPaths struct {
	ConfigPath string
	InfraPath  string
}

type withThresholdsOrdersReportingSink struct{}

func newWithThresholdsOrdersReportingSink() withThresholdsOrdersReportingSink {
	return withThresholdsOrdersReportingSink{}
}
func (withThresholdsOrdersReportingSink) SinkName() string {
	return "orders-sample-sink"
}
func (withThresholdsOrdersReportingSink) Init(loadstrike.LoadStrikeBaseContext, loadstrike.IConfiguration) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withThresholdsOrdersReportingSink) Start(loadstrike.LoadStrikeSessionStartInfo) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withThresholdsOrdersReportingSink) SaveRealtimeStats([]loadstrike.LoadStrikeScenarioStats) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withThresholdsOrdersReportingSink) SaveRealtimeMetrics(loadstrike.LoadStrikeMetricStats) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withThresholdsOrdersReportingSink) SaveRunResult(loadstrike.LoadStrikeRunResult) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withThresholdsOrdersReportingSink) Stop() loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withThresholdsOrdersReportingSink) Dispose() {}

type withThresholdsOrdersRuntimePolicy struct{}

func newWithThresholdsOrdersRuntimePolicy() withThresholdsOrdersRuntimePolicy {
	return withThresholdsOrdersRuntimePolicy{}
}
func (withThresholdsOrdersRuntimePolicy) ShouldRunScenario(string) loadstrike.LoadStrikeBoolTask {
	return loadstrike.TaskFromBool(true)
}
func (withThresholdsOrdersRuntimePolicy) BeforeScenario(string) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withThresholdsOrdersRuntimePolicy) AfterScenario(string, loadstrike.LoadStrikeScenarioRuntime) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withThresholdsOrdersRuntimePolicy) BeforeStep(string, string) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withThresholdsOrdersRuntimePolicy) AfterStep(string, string, loadstrike.LoadStrikeReply) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}

type withThresholdsOrdersWorkerPlugin struct{}

func newWithThresholdsOrdersWorkerPlugin() withThresholdsOrdersWorkerPlugin {
	return withThresholdsOrdersWorkerPlugin{}
}
func (withThresholdsOrdersWorkerPlugin) PluginName() string {
	return "orders-sample-plugin"
}
func (withThresholdsOrdersWorkerPlugin) Init(loadstrike.LoadStrikeBaseContext, loadstrike.IConfiguration) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withThresholdsOrdersWorkerPlugin) Start(loadstrike.LoadStrikeSessionStartInfo) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withThresholdsOrdersWorkerPlugin) GetData(loadstrike.LoadStrikeRunResult) loadstrike.LoadStrikeValueTask[loadstrike.LoadStrikePluginData] {
	return loadstrike.TaskFromResult(loadstrike.CreatePluginData("orders-sample-plugin"))
}
func (withThresholdsOrdersWorkerPlugin) Stop() loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withThresholdsOrdersWorkerPlugin) Dispose() loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}

func withThresholdsPerformOrderGetReply() loadstrike.LoadStrikeReply {
	return loadstrike.LoadStrikeResponse.Ok("200", int64(128), "ok", loadstrike.TimeSpan(3*time.Millisecond))
}

func withThresholdsExecuteOrderGet(context loadstrike.LoadStrikeScenarioContext) loadstrike.LoadStrikeReply {
	return loadstrike.LoadStrikeStep.Run("get-order", context, func(loadstrike.LoadStrikeScenarioContext) loadstrike.LoadStrikeReply {
		return withThresholdsPerformOrderGetReply()
	})
}

func withThresholdsExecuteOrderGetAsync(context loadstrike.LoadStrikeScenarioContext) loadstrike.LoadStrikeValueTask[loadstrike.LoadStrikeReply] {
	return loadstrike.TaskFromResult(withThresholdsExecuteOrderGet(context))
}

func withThresholdsBaselineScenario(names ...string) loadstrike.LoadStrikeScenario {
	name := "orders.get-by-id"
	if len(names) > 0 && names[0] != "" {
		name = names[0]
	}
	return loadstrike.CreateScenario(name, withThresholdsExecuteOrderGet).
		WithLoadSimulations(loadstrike.LoadStrikeSimulation.IterationsForConstant(1, 1)).
		WithoutWarmUp()
}

func withThresholdsBaseRunner() loadstrike.LoadStrikeRunner {
	return loadstrike.Create().
		AddScenario(withThresholdsBaselineScenario()).
		WithRunnerKey(withThresholdsRunnerKey).
		WithTestSuite("orders-reference").
		WithTestName("orders-get-by-id").
		WithoutReports()
}

func withThresholdsBaseContext() loadstrike.LoadStrikeContext {
	return withThresholdsBaseRunner().BuildContext()
}

func withThresholdsHttpSource() *loadstrike.EndpointSpec {
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

func withThresholdsHttpDestination() *loadstrike.EndpointSpec {
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

func withThresholdsTrackingConfiguration() *loadstrike.TrackingConfigurationSpec {
	return &loadstrike.TrackingConfigurationSpec{
		Source:                      withThresholdsHttpSource(),
		Destination:                 withThresholdsHttpDestination(),
		RunMode:                     "GenerateAndCorrelate",
		CorrelationTimeoutSeconds:   30,
		TimeoutSweepIntervalSeconds: 1,
		TimeoutBatchSize:            200,
		TimeoutCountsAsFailure:      true,
		MetricPrefix:                "orders_tracking",
		ExecuteOriginalScenarioRun:  false,
	}
}

func withThresholdsTrackedContext() loadstrike.LoadStrikeContext {
	return loadstrike.Create().
		AddScenario(withThresholdsBaselineScenario("orders.tracked").WithCrossPlatformTracking(withThresholdsTrackingConfiguration())).
		WithRunnerKey(withThresholdsRunnerKey).
		WithoutReports().
		BuildContext()
}

func withThresholdsBuildTrackingPayload() loadstrike.TrackingPayload {
	builder := loadstrike.TrackingPayloadBuilder{
		Headers:            map[string]string{"X-Correlation-Id": "ord-1"},
		ContentType:        "application/json",
		MessagePayloadType: "json",
	}
	builder.SetBody(map[string]any{"trackingId": "ord-1"})
	return builder.Build()
}

func withThresholdsRunResult() loadstrike.LoadStrikeRunResult {
	var result loadstrike.LoadStrikeRunResult
	return result
}

func withThresholdsScenarioStats() loadstrike.LoadStrikeScenarioStats {
	var stats loadstrike.LoadStrikeScenarioStats
	return stats
}

func withThresholdsWriteTempConfigFiles() withThresholdsTempConfigPaths {
	return withThresholdsTempConfigPaths{
		ConfigPath: "method-reference.loadstrike.config.json",
		InfraPath:  "method-reference.loadstrike.infra.json",
	}
}

        // Attach one scenario-level threshold.
func (reference WithThresholdsMethodReference) AttachScenarioThresholdExample() any {
    return withThresholdsBaselineScenario().WithThresholds(loadstrike.LoadStrikeThreshold{}.CreateScenario("allokcount", "gte", 1))
}

// Attach more than one threshold in the same call.
func (reference WithThresholdsMethodReference) AttachScenarioAndStepThresholdExample() any {
    return withThresholdsBaselineScenario().WithThresholds(loadstrike.LoadStrikeThreshold{}.CreateScenario("allokcount", "gte", 1), loadstrike.LoadStrikeThreshold{}.CreateStep("get-order", "allokcount", "gte", 1))
}
