package loadstrike_scenario

        import (
        	"time"

        	loadstrike "loadstrike.com/sdk/go"
        )

        const withRestartIterationOnFailRunnerKey = "runner_dummy_orders_reference"

type WithRestartIterationOnFailMethodReference struct{}

type withRestartIterationOnFailTempConfigPaths struct {
	ConfigPath string
	InfraPath  string
}

type withRestartIterationOnFailOrdersReportingSink struct{}

func newWithRestartIterationOnFailOrdersReportingSink() withRestartIterationOnFailOrdersReportingSink {
	return withRestartIterationOnFailOrdersReportingSink{}
}
func (withRestartIterationOnFailOrdersReportingSink) SinkName() string {
	return "orders-sample-sink"
}
func (withRestartIterationOnFailOrdersReportingSink) Init(loadstrike.LoadStrikeBaseContext, loadstrike.IConfiguration) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withRestartIterationOnFailOrdersReportingSink) Start(loadstrike.LoadStrikeSessionStartInfo) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withRestartIterationOnFailOrdersReportingSink) SaveRealtimeStats([]loadstrike.LoadStrikeScenarioStats) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withRestartIterationOnFailOrdersReportingSink) SaveRealtimeMetrics(loadstrike.LoadStrikeMetricStats) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withRestartIterationOnFailOrdersReportingSink) SaveRunResult(loadstrike.LoadStrikeRunResult) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withRestartIterationOnFailOrdersReportingSink) Stop() loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withRestartIterationOnFailOrdersReportingSink) Dispose() {}

type withRestartIterationOnFailOrdersRuntimePolicy struct{}

func newWithRestartIterationOnFailOrdersRuntimePolicy() withRestartIterationOnFailOrdersRuntimePolicy {
	return withRestartIterationOnFailOrdersRuntimePolicy{}
}
func (withRestartIterationOnFailOrdersRuntimePolicy) ShouldRunScenario(string) loadstrike.LoadStrikeBoolTask {
	return loadstrike.TaskFromBool(true)
}
func (withRestartIterationOnFailOrdersRuntimePolicy) BeforeScenario(string) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withRestartIterationOnFailOrdersRuntimePolicy) AfterScenario(string, loadstrike.LoadStrikeScenarioRuntime) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withRestartIterationOnFailOrdersRuntimePolicy) BeforeStep(string, string) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withRestartIterationOnFailOrdersRuntimePolicy) AfterStep(string, string, loadstrike.LoadStrikeReply) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}

type withRestartIterationOnFailOrdersWorkerPlugin struct{}

func newWithRestartIterationOnFailOrdersWorkerPlugin() withRestartIterationOnFailOrdersWorkerPlugin {
	return withRestartIterationOnFailOrdersWorkerPlugin{}
}
func (withRestartIterationOnFailOrdersWorkerPlugin) PluginName() string {
	return "orders-sample-plugin"
}
func (withRestartIterationOnFailOrdersWorkerPlugin) Init(loadstrike.LoadStrikeBaseContext, loadstrike.IConfiguration) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withRestartIterationOnFailOrdersWorkerPlugin) Start(loadstrike.LoadStrikeSessionStartInfo) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withRestartIterationOnFailOrdersWorkerPlugin) GetData(loadstrike.LoadStrikeRunResult) loadstrike.LoadStrikeValueTask[loadstrike.LoadStrikePluginData] {
	return loadstrike.TaskFromResult(loadstrike.CreatePluginData("orders-sample-plugin"))
}
func (withRestartIterationOnFailOrdersWorkerPlugin) Stop() loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withRestartIterationOnFailOrdersWorkerPlugin) Dispose() loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}

func withRestartIterationOnFailPerformOrderGetReply() loadstrike.LoadStrikeReply {
	return loadstrike.LoadStrikeResponse.Ok("200", int64(128), "ok", loadstrike.TimeSpan(3*time.Millisecond))
}

func withRestartIterationOnFailExecuteOrderGet(context loadstrike.LoadStrikeScenarioContext) loadstrike.LoadStrikeReply {
	return loadstrike.LoadStrikeStep.Run("get-order", context, func(loadstrike.LoadStrikeScenarioContext) loadstrike.LoadStrikeReply {
		return withRestartIterationOnFailPerformOrderGetReply()
	})
}

func withRestartIterationOnFailExecuteOrderGetAsync(context loadstrike.LoadStrikeScenarioContext) loadstrike.LoadStrikeValueTask[loadstrike.LoadStrikeReply] {
	return loadstrike.TaskFromResult(withRestartIterationOnFailExecuteOrderGet(context))
}

func withRestartIterationOnFailBaselineScenario(names ...string) loadstrike.LoadStrikeScenario {
	name := "orders.get-by-id"
	if len(names) > 0 && names[0] != "" {
		name = names[0]
	}
	return loadstrike.CreateScenario(name, withRestartIterationOnFailExecuteOrderGet).
		WithLoadSimulations(loadstrike.LoadStrikeSimulation.IterationsForConstant(1, 1)).
		WithoutWarmUp()
}

func withRestartIterationOnFailBaseRunner() loadstrike.LoadStrikeRunner {
	return loadstrike.Create().
		AddScenario(withRestartIterationOnFailBaselineScenario()).
		WithRunnerKey(withRestartIterationOnFailRunnerKey).
		WithTestSuite("orders-reference").
		WithTestName("orders-get-by-id").
		WithoutReports()
}

func withRestartIterationOnFailBaseContext() loadstrike.LoadStrikeContext {
	return withRestartIterationOnFailBaseRunner().BuildContext()
}

func withRestartIterationOnFailHttpSource() *loadstrike.EndpointSpec {
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

func withRestartIterationOnFailHttpDestination() *loadstrike.EndpointSpec {
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

func withRestartIterationOnFailTrackingConfiguration() *loadstrike.TrackingConfigurationSpec {
	return &loadstrike.TrackingConfigurationSpec{
		Source:                      withRestartIterationOnFailHttpSource(),
		Destination:                 withRestartIterationOnFailHttpDestination(),
		RunMode:                     "GenerateAndCorrelate",
		CorrelationTimeoutSeconds:   30,
		TimeoutSweepIntervalSeconds: 1,
		TimeoutBatchSize:            200,
		TimeoutCountsAsFailure:      true,
		MetricPrefix:                "orders_tracking",
		ExecuteOriginalScenarioRun:  false,
	}
}

func withRestartIterationOnFailTrackedContext() loadstrike.LoadStrikeContext {
	return loadstrike.Create().
		AddScenario(withRestartIterationOnFailBaselineScenario("orders.tracked").WithCrossPlatformTracking(withRestartIterationOnFailTrackingConfiguration())).
		WithRunnerKey(withRestartIterationOnFailRunnerKey).
		WithoutReports().
		BuildContext()
}

func withRestartIterationOnFailBuildTrackingPayload() loadstrike.TrackingPayload {
	builder := loadstrike.TrackingPayloadBuilder{
		Headers:            map[string]string{"X-Correlation-Id": "ord-1"},
		ContentType:        "application/json",
		MessagePayloadType: "json",
	}
	builder.SetBody(map[string]any{"trackingId": "ord-1"})
	return builder.Build()
}

func withRestartIterationOnFailRunResult() loadstrike.LoadStrikeRunResult {
	var result loadstrike.LoadStrikeRunResult
	return result
}

func withRestartIterationOnFailScenarioStats() loadstrike.LoadStrikeScenarioStats {
	var stats loadstrike.LoadStrikeScenarioStats
	return stats
}

func withRestartIterationOnFailWriteTempConfigFiles() withRestartIterationOnFailTempConfigPaths {
	return withRestartIterationOnFailTempConfigPaths{
		ConfigPath: "method-reference.loadstrike.config.json",
		InfraPath:  "method-reference.loadstrike.infra.json",
	}
}

        // Retry iterations when the GET step fails.
func (reference WithRestartIterationOnFailMethodReference) EnableRestartOnFailExample() any {
    return withRestartIterationOnFailBaselineScenario().WithRestartIterationOnFail(true)
}

// Leave iteration restart disabled.
func (reference WithRestartIterationOnFailMethodReference) DisableRestartOnFailExample() any {
    return withRestartIterationOnFailBaselineScenario().WithRestartIterationOnFail(false)
}
