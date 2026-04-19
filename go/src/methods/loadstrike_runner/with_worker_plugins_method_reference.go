package loadstrike_runner

        import (
        	"time"

        	loadstrike "loadstrike.com/sdk/go"
        )

        const withWorkerPluginsRunnerKey = "runner_dummy_orders_reference"

type WithWorkerPluginsMethodReference struct{}

type withWorkerPluginsTempConfigPaths struct {
	ConfigPath string
	InfraPath  string
}

type withWorkerPluginsOrdersReportingSink struct{}

func newWithWorkerPluginsOrdersReportingSink() withWorkerPluginsOrdersReportingSink {
	return withWorkerPluginsOrdersReportingSink{}
}
func (withWorkerPluginsOrdersReportingSink) SinkName() string {
	return "orders-sample-sink"
}
func (withWorkerPluginsOrdersReportingSink) Init(loadstrike.LoadStrikeBaseContext, loadstrike.IConfiguration) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withWorkerPluginsOrdersReportingSink) Start(loadstrike.LoadStrikeSessionStartInfo) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withWorkerPluginsOrdersReportingSink) SaveRealtimeStats([]loadstrike.LoadStrikeScenarioStats) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withWorkerPluginsOrdersReportingSink) SaveRealtimeMetrics(loadstrike.LoadStrikeMetricStats) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withWorkerPluginsOrdersReportingSink) SaveRunResult(loadstrike.LoadStrikeRunResult) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withWorkerPluginsOrdersReportingSink) Stop() loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withWorkerPluginsOrdersReportingSink) Dispose() {}

type withWorkerPluginsOrdersRuntimePolicy struct{}

func newWithWorkerPluginsOrdersRuntimePolicy() withWorkerPluginsOrdersRuntimePolicy {
	return withWorkerPluginsOrdersRuntimePolicy{}
}
func (withWorkerPluginsOrdersRuntimePolicy) ShouldRunScenario(string) loadstrike.LoadStrikeBoolTask {
	return loadstrike.TaskFromBool(true)
}
func (withWorkerPluginsOrdersRuntimePolicy) BeforeScenario(string) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withWorkerPluginsOrdersRuntimePolicy) AfterScenario(string, loadstrike.LoadStrikeScenarioRuntime) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withWorkerPluginsOrdersRuntimePolicy) BeforeStep(string, string) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withWorkerPluginsOrdersRuntimePolicy) AfterStep(string, string, loadstrike.LoadStrikeReply) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}

type withWorkerPluginsOrdersWorkerPlugin struct{}

func newWithWorkerPluginsOrdersWorkerPlugin() withWorkerPluginsOrdersWorkerPlugin {
	return withWorkerPluginsOrdersWorkerPlugin{}
}
func (withWorkerPluginsOrdersWorkerPlugin) PluginName() string {
	return "orders-sample-plugin"
}
func (withWorkerPluginsOrdersWorkerPlugin) Init(loadstrike.LoadStrikeBaseContext, loadstrike.IConfiguration) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withWorkerPluginsOrdersWorkerPlugin) Start(loadstrike.LoadStrikeSessionStartInfo) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withWorkerPluginsOrdersWorkerPlugin) GetData(loadstrike.LoadStrikeRunResult) loadstrike.LoadStrikeValueTask[loadstrike.LoadStrikePluginData] {
	return loadstrike.TaskFromResult(loadstrike.CreatePluginData("orders-sample-plugin"))
}
func (withWorkerPluginsOrdersWorkerPlugin) Stop() loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withWorkerPluginsOrdersWorkerPlugin) Dispose() loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}

func withWorkerPluginsPerformOrderGetReply() loadstrike.LoadStrikeReply {
	return loadstrike.LoadStrikeResponse.Ok("200", int64(128), "ok", loadstrike.TimeSpan(3*time.Millisecond))
}

func withWorkerPluginsExecuteOrderGet(context loadstrike.LoadStrikeScenarioContext) loadstrike.LoadStrikeReply {
	return loadstrike.LoadStrikeStep.Run("get-order", context, func(loadstrike.LoadStrikeScenarioContext) loadstrike.LoadStrikeReply {
		return withWorkerPluginsPerformOrderGetReply()
	})
}

func withWorkerPluginsExecuteOrderGetAsync(context loadstrike.LoadStrikeScenarioContext) loadstrike.LoadStrikeValueTask[loadstrike.LoadStrikeReply] {
	return loadstrike.TaskFromResult(withWorkerPluginsExecuteOrderGet(context))
}

func withWorkerPluginsBaselineScenario(names ...string) loadstrike.LoadStrikeScenario {
	name := "orders.get-by-id"
	if len(names) > 0 && names[0] != "" {
		name = names[0]
	}
	return loadstrike.CreateScenario(name, withWorkerPluginsExecuteOrderGet).
		WithLoadSimulations(loadstrike.LoadStrikeSimulation.IterationsForConstant(1, 1)).
		WithoutWarmUp()
}

func withWorkerPluginsBaseRunner() loadstrike.LoadStrikeRunner {
	return loadstrike.Create().
		AddScenario(withWorkerPluginsBaselineScenario()).
		WithRunnerKey(withWorkerPluginsRunnerKey).
		WithTestSuite("orders-reference").
		WithTestName("orders-get-by-id").
		WithoutReports()
}

func withWorkerPluginsBaseContext() loadstrike.LoadStrikeContext {
	return withWorkerPluginsBaseRunner().BuildContext()
}

func withWorkerPluginsHttpSource() *loadstrike.EndpointSpec {
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

func withWorkerPluginsHttpDestination() *loadstrike.EndpointSpec {
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

func withWorkerPluginsTrackingConfiguration() *loadstrike.TrackingConfigurationSpec {
	return &loadstrike.TrackingConfigurationSpec{
		Source:                      withWorkerPluginsHttpSource(),
		Destination:                 withWorkerPluginsHttpDestination(),
		RunMode:                     "GenerateAndCorrelate",
		CorrelationTimeoutSeconds:   30,
		TimeoutSweepIntervalSeconds: 1,
		TimeoutBatchSize:            200,
		TimeoutCountsAsFailure:      true,
		MetricPrefix:                "orders_tracking",
		ExecuteOriginalScenarioRun:  false,
	}
}

func withWorkerPluginsTrackedContext() loadstrike.LoadStrikeContext {
	return loadstrike.Create().
		AddScenario(withWorkerPluginsBaselineScenario("orders.tracked").WithCrossPlatformTracking(withWorkerPluginsTrackingConfiguration())).
		WithRunnerKey(withWorkerPluginsRunnerKey).
		WithoutReports().
		BuildContext()
}

func withWorkerPluginsBuildTrackingPayload() loadstrike.TrackingPayload {
	builder := loadstrike.TrackingPayloadBuilder{
		Headers:            map[string]string{"X-Correlation-Id": "ord-1"},
		ContentType:        "application/json",
		MessagePayloadType: "json",
	}
	builder.SetBody(map[string]any{"trackingId": "ord-1"})
	return builder.Build()
}

func withWorkerPluginsRunResult() loadstrike.LoadStrikeRunResult {
	var result loadstrike.LoadStrikeRunResult
	return result
}

func withWorkerPluginsScenarioStats() loadstrike.LoadStrikeScenarioStats {
	var stats loadstrike.LoadStrikeScenarioStats
	return stats
}

func withWorkerPluginsWriteTempConfigFiles() withWorkerPluginsTempConfigPaths {
	return withWorkerPluginsTempConfigPaths{
		ConfigPath: "method-reference.loadstrike.config.json",
		InfraPath:  "method-reference.loadstrike.infra.json",
	}
}

        // Attach one worker plugin instance.
func (reference WithWorkerPluginsMethodReference) AttachOneWorkerPluginExample() any {
    return withWorkerPluginsBaseRunner().WithWorkerPlugins(newWithWorkerPluginsOrdersWorkerPlugin())
}

// Attach multiple worker plugins in one call.
func (reference WithWorkerPluginsMethodReference) AttachTwoWorkerPluginsExample() any {
    return withWorkerPluginsBaseRunner().WithWorkerPlugins(newWithWorkerPluginsOrdersWorkerPlugin(), newWithWorkerPluginsOrdersWorkerPlugin())
}
