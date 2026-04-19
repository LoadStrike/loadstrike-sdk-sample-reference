package loadstrike_runner

        import (
        	"time"

        	loadstrike "loadstrike.com/sdk/go"
        )

        const withReportingSinksRunnerKey = "runner_dummy_orders_reference"

type WithReportingSinksMethodReference struct{}

type withReportingSinksTempConfigPaths struct {
	ConfigPath string
	InfraPath  string
}

type withReportingSinksOrdersReportingSink struct{}

func newWithReportingSinksOrdersReportingSink() withReportingSinksOrdersReportingSink {
	return withReportingSinksOrdersReportingSink{}
}
func (withReportingSinksOrdersReportingSink) SinkName() string {
	return "orders-sample-sink"
}
func (withReportingSinksOrdersReportingSink) Init(loadstrike.LoadStrikeBaseContext, loadstrike.IConfiguration) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withReportingSinksOrdersReportingSink) Start(loadstrike.LoadStrikeSessionStartInfo) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withReportingSinksOrdersReportingSink) SaveRealtimeStats([]loadstrike.LoadStrikeScenarioStats) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withReportingSinksOrdersReportingSink) SaveRealtimeMetrics(loadstrike.LoadStrikeMetricStats) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withReportingSinksOrdersReportingSink) SaveRunResult(loadstrike.LoadStrikeRunResult) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withReportingSinksOrdersReportingSink) Stop() loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withReportingSinksOrdersReportingSink) Dispose() {}

type withReportingSinksOrdersRuntimePolicy struct{}

func newWithReportingSinksOrdersRuntimePolicy() withReportingSinksOrdersRuntimePolicy {
	return withReportingSinksOrdersRuntimePolicy{}
}
func (withReportingSinksOrdersRuntimePolicy) ShouldRunScenario(string) loadstrike.LoadStrikeBoolTask {
	return loadstrike.TaskFromBool(true)
}
func (withReportingSinksOrdersRuntimePolicy) BeforeScenario(string) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withReportingSinksOrdersRuntimePolicy) AfterScenario(string, loadstrike.LoadStrikeScenarioRuntime) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withReportingSinksOrdersRuntimePolicy) BeforeStep(string, string) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withReportingSinksOrdersRuntimePolicy) AfterStep(string, string, loadstrike.LoadStrikeReply) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}

type withReportingSinksOrdersWorkerPlugin struct{}

func newWithReportingSinksOrdersWorkerPlugin() withReportingSinksOrdersWorkerPlugin {
	return withReportingSinksOrdersWorkerPlugin{}
}
func (withReportingSinksOrdersWorkerPlugin) PluginName() string {
	return "orders-sample-plugin"
}
func (withReportingSinksOrdersWorkerPlugin) Init(loadstrike.LoadStrikeBaseContext, loadstrike.IConfiguration) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withReportingSinksOrdersWorkerPlugin) Start(loadstrike.LoadStrikeSessionStartInfo) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withReportingSinksOrdersWorkerPlugin) GetData(loadstrike.LoadStrikeRunResult) loadstrike.LoadStrikeValueTask[loadstrike.LoadStrikePluginData] {
	return loadstrike.TaskFromResult(loadstrike.CreatePluginData("orders-sample-plugin"))
}
func (withReportingSinksOrdersWorkerPlugin) Stop() loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withReportingSinksOrdersWorkerPlugin) Dispose() loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}

func withReportingSinksPerformOrderGetReply() loadstrike.LoadStrikeReply {
	return loadstrike.LoadStrikeResponse.Ok("200", int64(128), "ok", loadstrike.TimeSpan(3*time.Millisecond))
}

func withReportingSinksExecuteOrderGet(context loadstrike.LoadStrikeScenarioContext) loadstrike.LoadStrikeReply {
	return loadstrike.LoadStrikeStep.Run("get-order", context, func(loadstrike.LoadStrikeScenarioContext) loadstrike.LoadStrikeReply {
		return withReportingSinksPerformOrderGetReply()
	})
}

func withReportingSinksExecuteOrderGetAsync(context loadstrike.LoadStrikeScenarioContext) loadstrike.LoadStrikeValueTask[loadstrike.LoadStrikeReply] {
	return loadstrike.TaskFromResult(withReportingSinksExecuteOrderGet(context))
}

func withReportingSinksBaselineScenario(names ...string) loadstrike.LoadStrikeScenario {
	name := "orders.get-by-id"
	if len(names) > 0 && names[0] != "" {
		name = names[0]
	}
	return loadstrike.CreateScenario(name, withReportingSinksExecuteOrderGet).
		WithLoadSimulations(loadstrike.LoadStrikeSimulation.IterationsForConstant(1, 1)).
		WithoutWarmUp()
}

func withReportingSinksBaseRunner() loadstrike.LoadStrikeRunner {
	return loadstrike.Create().
		AddScenario(withReportingSinksBaselineScenario()).
		WithRunnerKey(withReportingSinksRunnerKey).
		WithTestSuite("orders-reference").
		WithTestName("orders-get-by-id").
		WithoutReports()
}

func withReportingSinksBaseContext() loadstrike.LoadStrikeContext {
	return withReportingSinksBaseRunner().BuildContext()
}

func withReportingSinksHttpSource() *loadstrike.EndpointSpec {
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

func withReportingSinksHttpDestination() *loadstrike.EndpointSpec {
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

func withReportingSinksTrackingConfiguration() *loadstrike.TrackingConfigurationSpec {
	return &loadstrike.TrackingConfigurationSpec{
		Source:                      withReportingSinksHttpSource(),
		Destination:                 withReportingSinksHttpDestination(),
		RunMode:                     "GenerateAndCorrelate",
		CorrelationTimeoutSeconds:   30,
		TimeoutSweepIntervalSeconds: 1,
		TimeoutBatchSize:            200,
		TimeoutCountsAsFailure:      true,
		MetricPrefix:                "orders_tracking",
		ExecuteOriginalScenarioRun:  false,
	}
}

func withReportingSinksTrackedContext() loadstrike.LoadStrikeContext {
	return loadstrike.Create().
		AddScenario(withReportingSinksBaselineScenario("orders.tracked").WithCrossPlatformTracking(withReportingSinksTrackingConfiguration())).
		WithRunnerKey(withReportingSinksRunnerKey).
		WithoutReports().
		BuildContext()
}

func withReportingSinksBuildTrackingPayload() loadstrike.TrackingPayload {
	builder := loadstrike.TrackingPayloadBuilder{
		Headers:            map[string]string{"X-Correlation-Id": "ord-1"},
		ContentType:        "application/json",
		MessagePayloadType: "json",
	}
	builder.SetBody(map[string]any{"trackingId": "ord-1"})
	return builder.Build()
}

func withReportingSinksRunResult() loadstrike.LoadStrikeRunResult {
	var result loadstrike.LoadStrikeRunResult
	return result
}

func withReportingSinksScenarioStats() loadstrike.LoadStrikeScenarioStats {
	var stats loadstrike.LoadStrikeScenarioStats
	return stats
}

func withReportingSinksWriteTempConfigFiles() withReportingSinksTempConfigPaths {
	return withReportingSinksTempConfigPaths{
		ConfigPath: "method-reference.loadstrike.config.json",
		InfraPath:  "method-reference.loadstrike.infra.json",
	}
}

        // Attach one custom reporting sink instance.
func (reference WithReportingSinksMethodReference) AttachOneSinkExample() any {
    return withReportingSinksBaseRunner().WithReportingSinks(newWithReportingSinksOrdersReportingSink())
}

// Attach multiple sink instances in the same call.
func (reference WithReportingSinksMethodReference) AttachTwoSinksExample() any {
    return withReportingSinksBaseRunner().WithReportingSinks(newWithReportingSinksOrdersReportingSink(), newWithReportingSinksOrdersReportingSink())
}
