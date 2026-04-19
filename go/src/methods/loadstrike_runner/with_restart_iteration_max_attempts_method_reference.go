package loadstrike_runner

        import (
        	"time"

        	loadstrike "loadstrike.com/sdk/go"
        )

        const withRestartIterationMaxAttemptsRunnerKey = "runner_dummy_orders_reference"

type WithRestartIterationMaxAttemptsMethodReference struct{}

type withRestartIterationMaxAttemptsTempConfigPaths struct {
	ConfigPath string
	InfraPath  string
}

type withRestartIterationMaxAttemptsOrdersReportingSink struct{}

func newWithRestartIterationMaxAttemptsOrdersReportingSink() withRestartIterationMaxAttemptsOrdersReportingSink {
	return withRestartIterationMaxAttemptsOrdersReportingSink{}
}
func (withRestartIterationMaxAttemptsOrdersReportingSink) SinkName() string {
	return "orders-sample-sink"
}
func (withRestartIterationMaxAttemptsOrdersReportingSink) Init(loadstrike.LoadStrikeBaseContext, loadstrike.IConfiguration) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withRestartIterationMaxAttemptsOrdersReportingSink) Start(loadstrike.LoadStrikeSessionStartInfo) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withRestartIterationMaxAttemptsOrdersReportingSink) SaveRealtimeStats([]loadstrike.LoadStrikeScenarioStats) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withRestartIterationMaxAttemptsOrdersReportingSink) SaveRealtimeMetrics(loadstrike.LoadStrikeMetricStats) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withRestartIterationMaxAttemptsOrdersReportingSink) SaveRunResult(loadstrike.LoadStrikeRunResult) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withRestartIterationMaxAttemptsOrdersReportingSink) Stop() loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withRestartIterationMaxAttemptsOrdersReportingSink) Dispose() {}

type withRestartIterationMaxAttemptsOrdersRuntimePolicy struct{}

func newWithRestartIterationMaxAttemptsOrdersRuntimePolicy() withRestartIterationMaxAttemptsOrdersRuntimePolicy {
	return withRestartIterationMaxAttemptsOrdersRuntimePolicy{}
}
func (withRestartIterationMaxAttemptsOrdersRuntimePolicy) ShouldRunScenario(string) loadstrike.LoadStrikeBoolTask {
	return loadstrike.TaskFromBool(true)
}
func (withRestartIterationMaxAttemptsOrdersRuntimePolicy) BeforeScenario(string) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withRestartIterationMaxAttemptsOrdersRuntimePolicy) AfterScenario(string, loadstrike.LoadStrikeScenarioRuntime) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withRestartIterationMaxAttemptsOrdersRuntimePolicy) BeforeStep(string, string) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withRestartIterationMaxAttemptsOrdersRuntimePolicy) AfterStep(string, string, loadstrike.LoadStrikeReply) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}

type withRestartIterationMaxAttemptsOrdersWorkerPlugin struct{}

func newWithRestartIterationMaxAttemptsOrdersWorkerPlugin() withRestartIterationMaxAttemptsOrdersWorkerPlugin {
	return withRestartIterationMaxAttemptsOrdersWorkerPlugin{}
}
func (withRestartIterationMaxAttemptsOrdersWorkerPlugin) PluginName() string {
	return "orders-sample-plugin"
}
func (withRestartIterationMaxAttemptsOrdersWorkerPlugin) Init(loadstrike.LoadStrikeBaseContext, loadstrike.IConfiguration) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withRestartIterationMaxAttemptsOrdersWorkerPlugin) Start(loadstrike.LoadStrikeSessionStartInfo) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withRestartIterationMaxAttemptsOrdersWorkerPlugin) GetData(loadstrike.LoadStrikeRunResult) loadstrike.LoadStrikeValueTask[loadstrike.LoadStrikePluginData] {
	return loadstrike.TaskFromResult(loadstrike.CreatePluginData("orders-sample-plugin"))
}
func (withRestartIterationMaxAttemptsOrdersWorkerPlugin) Stop() loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withRestartIterationMaxAttemptsOrdersWorkerPlugin) Dispose() loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}

func withRestartIterationMaxAttemptsPerformOrderGetReply() loadstrike.LoadStrikeReply {
	return loadstrike.LoadStrikeResponse.Ok("200", int64(128), "ok", loadstrike.TimeSpan(3*time.Millisecond))
}

func withRestartIterationMaxAttemptsExecuteOrderGet(context loadstrike.LoadStrikeScenarioContext) loadstrike.LoadStrikeReply {
	return loadstrike.LoadStrikeStep.Run("get-order", context, func(loadstrike.LoadStrikeScenarioContext) loadstrike.LoadStrikeReply {
		return withRestartIterationMaxAttemptsPerformOrderGetReply()
	})
}

func withRestartIterationMaxAttemptsExecuteOrderGetAsync(context loadstrike.LoadStrikeScenarioContext) loadstrike.LoadStrikeValueTask[loadstrike.LoadStrikeReply] {
	return loadstrike.TaskFromResult(withRestartIterationMaxAttemptsExecuteOrderGet(context))
}

func withRestartIterationMaxAttemptsBaselineScenario(names ...string) loadstrike.LoadStrikeScenario {
	name := "orders.get-by-id"
	if len(names) > 0 && names[0] != "" {
		name = names[0]
	}
	return loadstrike.CreateScenario(name, withRestartIterationMaxAttemptsExecuteOrderGet).
		WithLoadSimulations(loadstrike.LoadStrikeSimulation.IterationsForConstant(1, 1)).
		WithoutWarmUp()
}

func withRestartIterationMaxAttemptsBaseRunner() loadstrike.LoadStrikeRunner {
	return loadstrike.Create().
		AddScenario(withRestartIterationMaxAttemptsBaselineScenario()).
		WithRunnerKey(withRestartIterationMaxAttemptsRunnerKey).
		WithTestSuite("orders-reference").
		WithTestName("orders-get-by-id").
		WithoutReports()
}

func withRestartIterationMaxAttemptsBaseContext() loadstrike.LoadStrikeContext {
	return withRestartIterationMaxAttemptsBaseRunner().BuildContext()
}

func withRestartIterationMaxAttemptsHttpSource() *loadstrike.EndpointSpec {
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

func withRestartIterationMaxAttemptsHttpDestination() *loadstrike.EndpointSpec {
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

func withRestartIterationMaxAttemptsTrackingConfiguration() *loadstrike.TrackingConfigurationSpec {
	return &loadstrike.TrackingConfigurationSpec{
		Source:                      withRestartIterationMaxAttemptsHttpSource(),
		Destination:                 withRestartIterationMaxAttemptsHttpDestination(),
		RunMode:                     "GenerateAndCorrelate",
		CorrelationTimeoutSeconds:   30,
		TimeoutSweepIntervalSeconds: 1,
		TimeoutBatchSize:            200,
		TimeoutCountsAsFailure:      true,
		MetricPrefix:                "orders_tracking",
		ExecuteOriginalScenarioRun:  false,
	}
}

func withRestartIterationMaxAttemptsTrackedContext() loadstrike.LoadStrikeContext {
	return loadstrike.Create().
		AddScenario(withRestartIterationMaxAttemptsBaselineScenario("orders.tracked").WithCrossPlatformTracking(withRestartIterationMaxAttemptsTrackingConfiguration())).
		WithRunnerKey(withRestartIterationMaxAttemptsRunnerKey).
		WithoutReports().
		BuildContext()
}

func withRestartIterationMaxAttemptsBuildTrackingPayload() loadstrike.TrackingPayload {
	builder := loadstrike.TrackingPayloadBuilder{
		Headers:            map[string]string{"X-Correlation-Id": "ord-1"},
		ContentType:        "application/json",
		MessagePayloadType: "json",
	}
	builder.SetBody(map[string]any{"trackingId": "ord-1"})
	return builder.Build()
}

func withRestartIterationMaxAttemptsRunResult() loadstrike.LoadStrikeRunResult {
	var result loadstrike.LoadStrikeRunResult
	return result
}

func withRestartIterationMaxAttemptsScenarioStats() loadstrike.LoadStrikeScenarioStats {
	var stats loadstrike.LoadStrikeScenarioStats
	return stats
}

func withRestartIterationMaxAttemptsWriteTempConfigFiles() withRestartIterationMaxAttemptsTempConfigPaths {
	return withRestartIterationMaxAttemptsTempConfigPaths{
		ConfigPath: "method-reference.loadstrike.config.json",
		InfraPath:  "method-reference.loadstrike.infra.json",
	}
}

        // Apply the option with a smaller numeric value.
func (reference WithRestartIterationMaxAttemptsMethodReference) ApplySmallerCountExample() any {
    return withRestartIterationMaxAttemptsBaseRunner().WithRestartIterationMaxAttempts(2)
}

// Apply the option with a larger numeric value.
func (reference WithRestartIterationMaxAttemptsMethodReference) ApplyLargerCountExample() any {
    return withRestartIterationMaxAttemptsBaseRunner().WithRestartIterationMaxAttempts(4)
}
