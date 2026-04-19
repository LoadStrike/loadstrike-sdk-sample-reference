package loadstrike_run_result

        import (
        	"time"

        	loadstrike "loadstrike.com/sdk/go"
        )

        const getScenarioStatsRunnerKey = "runner_dummy_orders_reference"

type GetScenarioStatsMethodReference struct{}

type getScenarioStatsTempConfigPaths struct {
	ConfigPath string
	InfraPath  string
}

type getScenarioStatsOrdersReportingSink struct{}

func newGetScenarioStatsOrdersReportingSink() getScenarioStatsOrdersReportingSink {
	return getScenarioStatsOrdersReportingSink{}
}
func (getScenarioStatsOrdersReportingSink) SinkName() string {
	return "orders-sample-sink"
}
func (getScenarioStatsOrdersReportingSink) Init(loadstrike.LoadStrikeBaseContext, loadstrike.IConfiguration) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (getScenarioStatsOrdersReportingSink) Start(loadstrike.LoadStrikeSessionStartInfo) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (getScenarioStatsOrdersReportingSink) SaveRealtimeStats([]loadstrike.LoadStrikeScenarioStats) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (getScenarioStatsOrdersReportingSink) SaveRealtimeMetrics(loadstrike.LoadStrikeMetricStats) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (getScenarioStatsOrdersReportingSink) SaveRunResult(loadstrike.LoadStrikeRunResult) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (getScenarioStatsOrdersReportingSink) Stop() loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (getScenarioStatsOrdersReportingSink) Dispose() {}

type getScenarioStatsOrdersRuntimePolicy struct{}

func newGetScenarioStatsOrdersRuntimePolicy() getScenarioStatsOrdersRuntimePolicy {
	return getScenarioStatsOrdersRuntimePolicy{}
}
func (getScenarioStatsOrdersRuntimePolicy) ShouldRunScenario(string) loadstrike.LoadStrikeBoolTask {
	return loadstrike.TaskFromBool(true)
}
func (getScenarioStatsOrdersRuntimePolicy) BeforeScenario(string) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (getScenarioStatsOrdersRuntimePolicy) AfterScenario(string, loadstrike.LoadStrikeScenarioRuntime) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (getScenarioStatsOrdersRuntimePolicy) BeforeStep(string, string) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (getScenarioStatsOrdersRuntimePolicy) AfterStep(string, string, loadstrike.LoadStrikeReply) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}

type getScenarioStatsOrdersWorkerPlugin struct{}

func newGetScenarioStatsOrdersWorkerPlugin() getScenarioStatsOrdersWorkerPlugin {
	return getScenarioStatsOrdersWorkerPlugin{}
}
func (getScenarioStatsOrdersWorkerPlugin) PluginName() string {
	return "orders-sample-plugin"
}
func (getScenarioStatsOrdersWorkerPlugin) Init(loadstrike.LoadStrikeBaseContext, loadstrike.IConfiguration) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (getScenarioStatsOrdersWorkerPlugin) Start(loadstrike.LoadStrikeSessionStartInfo) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (getScenarioStatsOrdersWorkerPlugin) GetData(loadstrike.LoadStrikeRunResult) loadstrike.LoadStrikeValueTask[loadstrike.LoadStrikePluginData] {
	return loadstrike.TaskFromResult(loadstrike.CreatePluginData("orders-sample-plugin"))
}
func (getScenarioStatsOrdersWorkerPlugin) Stop() loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (getScenarioStatsOrdersWorkerPlugin) Dispose() loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}

func getScenarioStatsPerformOrderGetReply() loadstrike.LoadStrikeReply {
	return loadstrike.LoadStrikeResponse.Ok("200", int64(128), "ok", loadstrike.TimeSpan(3*time.Millisecond))
}

func getScenarioStatsExecuteOrderGet(context loadstrike.LoadStrikeScenarioContext) loadstrike.LoadStrikeReply {
	return loadstrike.LoadStrikeStep.Run("get-order", context, func(loadstrike.LoadStrikeScenarioContext) loadstrike.LoadStrikeReply {
		return getScenarioStatsPerformOrderGetReply()
	})
}

func getScenarioStatsExecuteOrderGetAsync(context loadstrike.LoadStrikeScenarioContext) loadstrike.LoadStrikeValueTask[loadstrike.LoadStrikeReply] {
	return loadstrike.TaskFromResult(getScenarioStatsExecuteOrderGet(context))
}

func getScenarioStatsBaselineScenario(names ...string) loadstrike.LoadStrikeScenario {
	name := "orders.get-by-id"
	if len(names) > 0 && names[0] != "" {
		name = names[0]
	}
	return loadstrike.CreateScenario(name, getScenarioStatsExecuteOrderGet).
		WithLoadSimulations(loadstrike.LoadStrikeSimulation.IterationsForConstant(1, 1)).
		WithoutWarmUp()
}

func getScenarioStatsBaseRunner() loadstrike.LoadStrikeRunner {
	return loadstrike.Create().
		AddScenario(getScenarioStatsBaselineScenario()).
		WithRunnerKey(getScenarioStatsRunnerKey).
		WithTestSuite("orders-reference").
		WithTestName("orders-get-by-id").
		WithoutReports()
}

func getScenarioStatsBaseContext() loadstrike.LoadStrikeContext {
	return getScenarioStatsBaseRunner().BuildContext()
}

func getScenarioStatsHttpSource() *loadstrike.EndpointSpec {
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

func getScenarioStatsHttpDestination() *loadstrike.EndpointSpec {
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

func getScenarioStatsTrackingConfiguration() *loadstrike.TrackingConfigurationSpec {
	return &loadstrike.TrackingConfigurationSpec{
		Source:                      getScenarioStatsHttpSource(),
		Destination:                 getScenarioStatsHttpDestination(),
		RunMode:                     "GenerateAndCorrelate",
		CorrelationTimeoutSeconds:   30,
		TimeoutSweepIntervalSeconds: 1,
		TimeoutBatchSize:            200,
		TimeoutCountsAsFailure:      true,
		MetricPrefix:                "orders_tracking",
		ExecuteOriginalScenarioRun:  false,
	}
}

func getScenarioStatsTrackedContext() loadstrike.LoadStrikeContext {
	return loadstrike.Create().
		AddScenario(getScenarioStatsBaselineScenario("orders.tracked").WithCrossPlatformTracking(getScenarioStatsTrackingConfiguration())).
		WithRunnerKey(getScenarioStatsRunnerKey).
		WithoutReports().
		BuildContext()
}

func getScenarioStatsBuildTrackingPayload() loadstrike.TrackingPayload {
	builder := loadstrike.TrackingPayloadBuilder{
		Headers:            map[string]string{"X-Correlation-Id": "ord-1"},
		ContentType:        "application/json",
		MessagePayloadType: "json",
	}
	builder.SetBody(map[string]any{"trackingId": "ord-1"})
	return builder.Build()
}

func getScenarioStatsRunResult() loadstrike.LoadStrikeRunResult {
	var result loadstrike.LoadStrikeRunResult
	return result
}

func getScenarioStatsScenarioStats() loadstrike.LoadStrikeScenarioStats {
	var stats loadstrike.LoadStrikeScenarioStats
	return stats
}

func getScenarioStatsWriteTempConfigFiles() getScenarioStatsTempConfigPaths {
	return getScenarioStatsTempConfigPaths{
		ConfigPath: "method-reference.loadstrike.config.json",
		InfraPath:  "method-reference.loadstrike.infra.json",
	}
}

        // Read data back from the run result surface with this helper.
func (reference GetScenarioStatsMethodReference) ReadRunResultsExample() any {
    return getScenarioStatsRunResult().GetScenarioStats("orders.get-by-id")
}

// Continue into the scenario or step stats surface after the first lookup.
func (reference GetScenarioStatsMethodReference) ReadDeeperResultSurfaceExample() any {
    return getScenarioStatsScenarioStats().FindStepStats("get-order")
}
