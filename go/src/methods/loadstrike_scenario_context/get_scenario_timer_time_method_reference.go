package loadstrike_scenario_context

        import (
        	"time"

        	loadstrike "loadstrike.com/sdk/go"
        )

        const getScenarioTimerTimeRunnerKey = "runner_dummy_orders_reference"

type GetScenarioTimerTimeMethodReference struct{}

type getScenarioTimerTimeTempConfigPaths struct {
	ConfigPath string
	InfraPath  string
}

type getScenarioTimerTimeOrdersReportingSink struct{}

func newGetScenarioTimerTimeOrdersReportingSink() getScenarioTimerTimeOrdersReportingSink {
	return getScenarioTimerTimeOrdersReportingSink{}
}
func (getScenarioTimerTimeOrdersReportingSink) SinkName() string {
	return "orders-sample-sink"
}
func (getScenarioTimerTimeOrdersReportingSink) Init(loadstrike.LoadStrikeBaseContext, loadstrike.IConfiguration) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (getScenarioTimerTimeOrdersReportingSink) Start(loadstrike.LoadStrikeSessionStartInfo) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (getScenarioTimerTimeOrdersReportingSink) SaveRealtimeStats([]loadstrike.LoadStrikeScenarioStats) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (getScenarioTimerTimeOrdersReportingSink) SaveRealtimeMetrics(loadstrike.LoadStrikeMetricStats) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (getScenarioTimerTimeOrdersReportingSink) SaveRunResult(loadstrike.LoadStrikeRunResult) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (getScenarioTimerTimeOrdersReportingSink) Stop() loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (getScenarioTimerTimeOrdersReportingSink) Dispose() {}

type getScenarioTimerTimeOrdersRuntimePolicy struct{}

func newGetScenarioTimerTimeOrdersRuntimePolicy() getScenarioTimerTimeOrdersRuntimePolicy {
	return getScenarioTimerTimeOrdersRuntimePolicy{}
}
func (getScenarioTimerTimeOrdersRuntimePolicy) ShouldRunScenario(string) loadstrike.LoadStrikeBoolTask {
	return loadstrike.TaskFromBool(true)
}
func (getScenarioTimerTimeOrdersRuntimePolicy) BeforeScenario(string) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (getScenarioTimerTimeOrdersRuntimePolicy) AfterScenario(string, loadstrike.LoadStrikeScenarioRuntime) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (getScenarioTimerTimeOrdersRuntimePolicy) BeforeStep(string, string) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (getScenarioTimerTimeOrdersRuntimePolicy) AfterStep(string, string, loadstrike.LoadStrikeReply) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}

type getScenarioTimerTimeOrdersWorkerPlugin struct{}

func newGetScenarioTimerTimeOrdersWorkerPlugin() getScenarioTimerTimeOrdersWorkerPlugin {
	return getScenarioTimerTimeOrdersWorkerPlugin{}
}
func (getScenarioTimerTimeOrdersWorkerPlugin) PluginName() string {
	return "orders-sample-plugin"
}
func (getScenarioTimerTimeOrdersWorkerPlugin) Init(loadstrike.LoadStrikeBaseContext, loadstrike.IConfiguration) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (getScenarioTimerTimeOrdersWorkerPlugin) Start(loadstrike.LoadStrikeSessionStartInfo) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (getScenarioTimerTimeOrdersWorkerPlugin) GetData(loadstrike.LoadStrikeRunResult) loadstrike.LoadStrikeValueTask[loadstrike.LoadStrikePluginData] {
	return loadstrike.TaskFromResult(loadstrike.CreatePluginData("orders-sample-plugin"))
}
func (getScenarioTimerTimeOrdersWorkerPlugin) Stop() loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (getScenarioTimerTimeOrdersWorkerPlugin) Dispose() loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}

func getScenarioTimerTimePerformOrderGetReply() loadstrike.LoadStrikeReply {
	return loadstrike.LoadStrikeResponse.Ok("200", int64(128), "ok", loadstrike.TimeSpan(3*time.Millisecond))
}

func getScenarioTimerTimeExecuteOrderGet(context loadstrike.LoadStrikeScenarioContext) loadstrike.LoadStrikeReply {
	return loadstrike.LoadStrikeStep.Run("get-order", context, func(loadstrike.LoadStrikeScenarioContext) loadstrike.LoadStrikeReply {
		return getScenarioTimerTimePerformOrderGetReply()
	})
}

func getScenarioTimerTimeExecuteOrderGetAsync(context loadstrike.LoadStrikeScenarioContext) loadstrike.LoadStrikeValueTask[loadstrike.LoadStrikeReply] {
	return loadstrike.TaskFromResult(getScenarioTimerTimeExecuteOrderGet(context))
}

func getScenarioTimerTimeBaselineScenario(names ...string) loadstrike.LoadStrikeScenario {
	name := "orders.get-by-id"
	if len(names) > 0 && names[0] != "" {
		name = names[0]
	}
	return loadstrike.CreateScenario(name, getScenarioTimerTimeExecuteOrderGet).
		WithLoadSimulations(loadstrike.LoadStrikeSimulation.IterationsForConstant(1, 1)).
		WithoutWarmUp()
}

func getScenarioTimerTimeBaseRunner() loadstrike.LoadStrikeRunner {
	return loadstrike.Create().
		AddScenario(getScenarioTimerTimeBaselineScenario()).
		WithRunnerKey(getScenarioTimerTimeRunnerKey).
		WithTestSuite("orders-reference").
		WithTestName("orders-get-by-id").
		WithoutReports()
}

func getScenarioTimerTimeBaseContext() loadstrike.LoadStrikeContext {
	return getScenarioTimerTimeBaseRunner().BuildContext()
}

func getScenarioTimerTimeHttpSource() *loadstrike.EndpointSpec {
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

func getScenarioTimerTimeHttpDestination() *loadstrike.EndpointSpec {
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

func getScenarioTimerTimeTrackingConfiguration() *loadstrike.TrackingConfigurationSpec {
	return &loadstrike.TrackingConfigurationSpec{
		Source:                      getScenarioTimerTimeHttpSource(),
		Destination:                 getScenarioTimerTimeHttpDestination(),
		RunMode:                     "GenerateAndCorrelate",
		CorrelationTimeoutSeconds:   30,
		TimeoutSweepIntervalSeconds: 1,
		TimeoutBatchSize:            200,
		TimeoutCountsAsFailure:      true,
		MetricPrefix:                "orders_tracking",
		ExecuteOriginalScenarioRun:  false,
	}
}

func getScenarioTimerTimeTrackedContext() loadstrike.LoadStrikeContext {
	return loadstrike.Create().
		AddScenario(getScenarioTimerTimeBaselineScenario("orders.tracked").WithCrossPlatformTracking(getScenarioTimerTimeTrackingConfiguration())).
		WithRunnerKey(getScenarioTimerTimeRunnerKey).
		WithoutReports().
		BuildContext()
}

func getScenarioTimerTimeBuildTrackingPayload() loadstrike.TrackingPayload {
	builder := loadstrike.TrackingPayloadBuilder{
		Headers:            map[string]string{"X-Correlation-Id": "ord-1"},
		ContentType:        "application/json",
		MessagePayloadType: "json",
	}
	builder.SetBody(map[string]any{"trackingId": "ord-1"})
	return builder.Build()
}

func getScenarioTimerTimeRunResult() loadstrike.LoadStrikeRunResult {
	var result loadstrike.LoadStrikeRunResult
	return result
}

func getScenarioTimerTimeScenarioStats() loadstrike.LoadStrikeScenarioStats {
	var stats loadstrike.LoadStrikeScenarioStats
	return stats
}

func getScenarioTimerTimeWriteTempConfigFiles() getScenarioTimerTimeTempConfigPaths {
	return getScenarioTimerTimeTempConfigPaths{
		ConfigPath: "method-reference.loadstrike.config.json",
		InfraPath:  "method-reference.loadstrike.infra.json",
	}
}

        // Call the public context helper directly from the scenario context surface.
func (reference GetScenarioTimerTimeMethodReference) UseContextMethodExample() any {
    return loadstrike.CreateScenario("orders.timer", func(context loadstrike.LoadStrikeScenarioContext) loadstrike.LoadStrikeReply { _ = context.GetScenarioTimerTime(); return loadstrike.OK() })
}

// Show the same helper in the baseline GET-step flow.
func (reference GetScenarioTimerTimeMethodReference) UseContextMethodInStepExample() any {
    return loadstrike.CreateScenario("orders.timer", func(context loadstrike.LoadStrikeScenarioContext) loadstrike.LoadStrikeReply { _ = context.GetScenarioTimerTime(); return loadstrike.OK() })
}
