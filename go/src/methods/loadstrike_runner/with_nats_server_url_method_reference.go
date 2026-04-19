package loadstrike_runner

        import (
        	"time"

        	loadstrike "loadstrike.com/sdk/go"
        )

        const withNatsServerUrlRunnerKey = "runner_dummy_orders_reference"

type WithNatsServerUrlMethodReference struct{}

type withNatsServerUrlTempConfigPaths struct {
	ConfigPath string
	InfraPath  string
}

type withNatsServerUrlOrdersReportingSink struct{}

func newWithNatsServerUrlOrdersReportingSink() withNatsServerUrlOrdersReportingSink {
	return withNatsServerUrlOrdersReportingSink{}
}
func (withNatsServerUrlOrdersReportingSink) SinkName() string {
	return "orders-sample-sink"
}
func (withNatsServerUrlOrdersReportingSink) Init(loadstrike.LoadStrikeBaseContext, loadstrike.IConfiguration) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withNatsServerUrlOrdersReportingSink) Start(loadstrike.LoadStrikeSessionStartInfo) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withNatsServerUrlOrdersReportingSink) SaveRealtimeStats([]loadstrike.LoadStrikeScenarioStats) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withNatsServerUrlOrdersReportingSink) SaveRealtimeMetrics(loadstrike.LoadStrikeMetricStats) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withNatsServerUrlOrdersReportingSink) SaveRunResult(loadstrike.LoadStrikeRunResult) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withNatsServerUrlOrdersReportingSink) Stop() loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withNatsServerUrlOrdersReportingSink) Dispose() {}

type withNatsServerUrlOrdersRuntimePolicy struct{}

func newWithNatsServerUrlOrdersRuntimePolicy() withNatsServerUrlOrdersRuntimePolicy {
	return withNatsServerUrlOrdersRuntimePolicy{}
}
func (withNatsServerUrlOrdersRuntimePolicy) ShouldRunScenario(string) loadstrike.LoadStrikeBoolTask {
	return loadstrike.TaskFromBool(true)
}
func (withNatsServerUrlOrdersRuntimePolicy) BeforeScenario(string) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withNatsServerUrlOrdersRuntimePolicy) AfterScenario(string, loadstrike.LoadStrikeScenarioRuntime) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withNatsServerUrlOrdersRuntimePolicy) BeforeStep(string, string) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withNatsServerUrlOrdersRuntimePolicy) AfterStep(string, string, loadstrike.LoadStrikeReply) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}

type withNatsServerUrlOrdersWorkerPlugin struct{}

func newWithNatsServerUrlOrdersWorkerPlugin() withNatsServerUrlOrdersWorkerPlugin {
	return withNatsServerUrlOrdersWorkerPlugin{}
}
func (withNatsServerUrlOrdersWorkerPlugin) PluginName() string {
	return "orders-sample-plugin"
}
func (withNatsServerUrlOrdersWorkerPlugin) Init(loadstrike.LoadStrikeBaseContext, loadstrike.IConfiguration) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withNatsServerUrlOrdersWorkerPlugin) Start(loadstrike.LoadStrikeSessionStartInfo) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withNatsServerUrlOrdersWorkerPlugin) GetData(loadstrike.LoadStrikeRunResult) loadstrike.LoadStrikeValueTask[loadstrike.LoadStrikePluginData] {
	return loadstrike.TaskFromResult(loadstrike.CreatePluginData("orders-sample-plugin"))
}
func (withNatsServerUrlOrdersWorkerPlugin) Stop() loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withNatsServerUrlOrdersWorkerPlugin) Dispose() loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}

func withNatsServerUrlPerformOrderGetReply() loadstrike.LoadStrikeReply {
	return loadstrike.LoadStrikeResponse.Ok("200", int64(128), "ok", loadstrike.TimeSpan(3*time.Millisecond))
}

func withNatsServerUrlExecuteOrderGet(context loadstrike.LoadStrikeScenarioContext) loadstrike.LoadStrikeReply {
	return loadstrike.LoadStrikeStep.Run("get-order", context, func(loadstrike.LoadStrikeScenarioContext) loadstrike.LoadStrikeReply {
		return withNatsServerUrlPerformOrderGetReply()
	})
}

func withNatsServerUrlExecuteOrderGetAsync(context loadstrike.LoadStrikeScenarioContext) loadstrike.LoadStrikeValueTask[loadstrike.LoadStrikeReply] {
	return loadstrike.TaskFromResult(withNatsServerUrlExecuteOrderGet(context))
}

func withNatsServerUrlBaselineScenario(names ...string) loadstrike.LoadStrikeScenario {
	name := "orders.get-by-id"
	if len(names) > 0 && names[0] != "" {
		name = names[0]
	}
	return loadstrike.CreateScenario(name, withNatsServerUrlExecuteOrderGet).
		WithLoadSimulations(loadstrike.LoadStrikeSimulation.IterationsForConstant(1, 1)).
		WithoutWarmUp()
}

func withNatsServerUrlBaseRunner() loadstrike.LoadStrikeRunner {
	return loadstrike.Create().
		AddScenario(withNatsServerUrlBaselineScenario()).
		WithRunnerKey(withNatsServerUrlRunnerKey).
		WithTestSuite("orders-reference").
		WithTestName("orders-get-by-id").
		WithoutReports()
}

func withNatsServerUrlBaseContext() loadstrike.LoadStrikeContext {
	return withNatsServerUrlBaseRunner().BuildContext()
}

func withNatsServerUrlHttpSource() *loadstrike.EndpointSpec {
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

func withNatsServerUrlHttpDestination() *loadstrike.EndpointSpec {
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

func withNatsServerUrlTrackingConfiguration() *loadstrike.TrackingConfigurationSpec {
	return &loadstrike.TrackingConfigurationSpec{
		Source:                      withNatsServerUrlHttpSource(),
		Destination:                 withNatsServerUrlHttpDestination(),
		RunMode:                     "GenerateAndCorrelate",
		CorrelationTimeoutSeconds:   30,
		TimeoutSweepIntervalSeconds: 1,
		TimeoutBatchSize:            200,
		TimeoutCountsAsFailure:      true,
		MetricPrefix:                "orders_tracking",
		ExecuteOriginalScenarioRun:  false,
	}
}

func withNatsServerUrlTrackedContext() loadstrike.LoadStrikeContext {
	return loadstrike.Create().
		AddScenario(withNatsServerUrlBaselineScenario("orders.tracked").WithCrossPlatformTracking(withNatsServerUrlTrackingConfiguration())).
		WithRunnerKey(withNatsServerUrlRunnerKey).
		WithoutReports().
		BuildContext()
}

func withNatsServerUrlBuildTrackingPayload() loadstrike.TrackingPayload {
	builder := loadstrike.TrackingPayloadBuilder{
		Headers:            map[string]string{"X-Correlation-Id": "ord-1"},
		ContentType:        "application/json",
		MessagePayloadType: "json",
	}
	builder.SetBody(map[string]any{"trackingId": "ord-1"})
	return builder.Build()
}

func withNatsServerUrlRunResult() loadstrike.LoadStrikeRunResult {
	var result loadstrike.LoadStrikeRunResult
	return result
}

func withNatsServerUrlScenarioStats() loadstrike.LoadStrikeScenarioStats {
	var stats loadstrike.LoadStrikeScenarioStats
	return stats
}

func withNatsServerUrlWriteTempConfigFiles() withNatsServerUrlTempConfigPaths {
	return withNatsServerUrlTempConfigPaths{
		ConfigPath: "method-reference.loadstrike.config.json",
		InfraPath:  "method-reference.loadstrike.infra.json",
	}
}

        // Apply the primary string value shown in the sample method reference.
func (reference WithNatsServerUrlMethodReference) ApplyPrimaryValueExample() any {
    return withNatsServerUrlBaseRunner().WithNatsServerURL("nats://localhost:4222")
}

// Show the same method with a second concrete value.
func (reference WithNatsServerUrlMethodReference) ApplyAlternateValueExample() any {
    return withNatsServerUrlBaseRunner().WithNatsServerURL("nats://demo-cluster:4222")
}
