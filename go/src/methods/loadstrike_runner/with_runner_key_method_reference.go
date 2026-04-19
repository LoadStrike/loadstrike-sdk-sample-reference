package loadstrike_runner

        import (
        	"time"

        	loadstrike "loadstrike.com/sdk/go"
        )

        const withRunnerKeyRunnerKey = "runner_dummy_orders_reference"

type WithRunnerKeyMethodReference struct{}

type withRunnerKeyTempConfigPaths struct {
	ConfigPath string
	InfraPath  string
}

type withRunnerKeyOrdersReportingSink struct{}

func newWithRunnerKeyOrdersReportingSink() withRunnerKeyOrdersReportingSink {
	return withRunnerKeyOrdersReportingSink{}
}
func (withRunnerKeyOrdersReportingSink) SinkName() string {
	return "orders-sample-sink"
}
func (withRunnerKeyOrdersReportingSink) Init(loadstrike.LoadStrikeBaseContext, loadstrike.IConfiguration) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withRunnerKeyOrdersReportingSink) Start(loadstrike.LoadStrikeSessionStartInfo) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withRunnerKeyOrdersReportingSink) SaveRealtimeStats([]loadstrike.LoadStrikeScenarioStats) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withRunnerKeyOrdersReportingSink) SaveRealtimeMetrics(loadstrike.LoadStrikeMetricStats) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withRunnerKeyOrdersReportingSink) SaveRunResult(loadstrike.LoadStrikeRunResult) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withRunnerKeyOrdersReportingSink) Stop() loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withRunnerKeyOrdersReportingSink) Dispose() {}

type withRunnerKeyOrdersRuntimePolicy struct{}

func newWithRunnerKeyOrdersRuntimePolicy() withRunnerKeyOrdersRuntimePolicy {
	return withRunnerKeyOrdersRuntimePolicy{}
}
func (withRunnerKeyOrdersRuntimePolicy) ShouldRunScenario(string) loadstrike.LoadStrikeBoolTask {
	return loadstrike.TaskFromBool(true)
}
func (withRunnerKeyOrdersRuntimePolicy) BeforeScenario(string) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withRunnerKeyOrdersRuntimePolicy) AfterScenario(string, loadstrike.LoadStrikeScenarioRuntime) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withRunnerKeyOrdersRuntimePolicy) BeforeStep(string, string) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withRunnerKeyOrdersRuntimePolicy) AfterStep(string, string, loadstrike.LoadStrikeReply) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}

type withRunnerKeyOrdersWorkerPlugin struct{}

func newWithRunnerKeyOrdersWorkerPlugin() withRunnerKeyOrdersWorkerPlugin {
	return withRunnerKeyOrdersWorkerPlugin{}
}
func (withRunnerKeyOrdersWorkerPlugin) PluginName() string {
	return "orders-sample-plugin"
}
func (withRunnerKeyOrdersWorkerPlugin) Init(loadstrike.LoadStrikeBaseContext, loadstrike.IConfiguration) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withRunnerKeyOrdersWorkerPlugin) Start(loadstrike.LoadStrikeSessionStartInfo) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withRunnerKeyOrdersWorkerPlugin) GetData(loadstrike.LoadStrikeRunResult) loadstrike.LoadStrikeValueTask[loadstrike.LoadStrikePluginData] {
	return loadstrike.TaskFromResult(loadstrike.CreatePluginData("orders-sample-plugin"))
}
func (withRunnerKeyOrdersWorkerPlugin) Stop() loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withRunnerKeyOrdersWorkerPlugin) Dispose() loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}

func withRunnerKeyPerformOrderGetReply() loadstrike.LoadStrikeReply {
	return loadstrike.LoadStrikeResponse.Ok("200", int64(128), "ok", loadstrike.TimeSpan(3*time.Millisecond))
}

func withRunnerKeyExecuteOrderGet(context loadstrike.LoadStrikeScenarioContext) loadstrike.LoadStrikeReply {
	return loadstrike.LoadStrikeStep.Run("get-order", context, func(loadstrike.LoadStrikeScenarioContext) loadstrike.LoadStrikeReply {
		return withRunnerKeyPerformOrderGetReply()
	})
}

func withRunnerKeyExecuteOrderGetAsync(context loadstrike.LoadStrikeScenarioContext) loadstrike.LoadStrikeValueTask[loadstrike.LoadStrikeReply] {
	return loadstrike.TaskFromResult(withRunnerKeyExecuteOrderGet(context))
}

func withRunnerKeyBaselineScenario(names ...string) loadstrike.LoadStrikeScenario {
	name := "orders.get-by-id"
	if len(names) > 0 && names[0] != "" {
		name = names[0]
	}
	return loadstrike.CreateScenario(name, withRunnerKeyExecuteOrderGet).
		WithLoadSimulations(loadstrike.LoadStrikeSimulation.IterationsForConstant(1, 1)).
		WithoutWarmUp()
}

func withRunnerKeyBaseRunner() loadstrike.LoadStrikeRunner {
	return loadstrike.Create().
		AddScenario(withRunnerKeyBaselineScenario()).
		WithRunnerKey(withRunnerKeyRunnerKey).
		WithTestSuite("orders-reference").
		WithTestName("orders-get-by-id").
		WithoutReports()
}

func withRunnerKeyBaseContext() loadstrike.LoadStrikeContext {
	return withRunnerKeyBaseRunner().BuildContext()
}

func withRunnerKeyHttpSource() *loadstrike.EndpointSpec {
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

func withRunnerKeyHttpDestination() *loadstrike.EndpointSpec {
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

func withRunnerKeyTrackingConfiguration() *loadstrike.TrackingConfigurationSpec {
	return &loadstrike.TrackingConfigurationSpec{
		Source:                      withRunnerKeyHttpSource(),
		Destination:                 withRunnerKeyHttpDestination(),
		RunMode:                     "GenerateAndCorrelate",
		CorrelationTimeoutSeconds:   30,
		TimeoutSweepIntervalSeconds: 1,
		TimeoutBatchSize:            200,
		TimeoutCountsAsFailure:      true,
		MetricPrefix:                "orders_tracking",
		ExecuteOriginalScenarioRun:  false,
	}
}

func withRunnerKeyTrackedContext() loadstrike.LoadStrikeContext {
	return loadstrike.Create().
		AddScenario(withRunnerKeyBaselineScenario("orders.tracked").WithCrossPlatformTracking(withRunnerKeyTrackingConfiguration())).
		WithRunnerKey(withRunnerKeyRunnerKey).
		WithoutReports().
		BuildContext()
}

func withRunnerKeyBuildTrackingPayload() loadstrike.TrackingPayload {
	builder := loadstrike.TrackingPayloadBuilder{
		Headers:            map[string]string{"X-Correlation-Id": "ord-1"},
		ContentType:        "application/json",
		MessagePayloadType: "json",
	}
	builder.SetBody(map[string]any{"trackingId": "ord-1"})
	return builder.Build()
}

func withRunnerKeyRunResult() loadstrike.LoadStrikeRunResult {
	var result loadstrike.LoadStrikeRunResult
	return result
}

func withRunnerKeyScenarioStats() loadstrike.LoadStrikeScenarioStats {
	var stats loadstrike.LoadStrikeScenarioStats
	return stats
}

func withRunnerKeyWriteTempConfigFiles() withRunnerKeyTempConfigPaths {
	return withRunnerKeyTempConfigPaths{
		ConfigPath: "method-reference.loadstrike.config.json",
		InfraPath:  "method-reference.loadstrike.infra.json",
	}
}

        // Apply the primary string value shown in the sample method reference.
func (reference WithRunnerKeyMethodReference) ApplyPrimaryValueExample() any {
    return withRunnerKeyBaseRunner().WithRunnerKey("runner_dummy_orders_reference")
}

// Show the same method with a second concrete value.
func (reference WithRunnerKeyMethodReference) ApplyAlternateValueExample() any {
    return withRunnerKeyBaseRunner().WithRunnerKey("runner_dummy_checkout_reference")
}
