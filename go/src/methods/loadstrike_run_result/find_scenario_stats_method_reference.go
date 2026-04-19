package loadstrike_run_result

        import (
        	"time"

        	loadstrike "loadstrike.com/sdk/go"
        )

        const findScenarioStatsRunnerKey = "runner_dummy_orders_reference"

type FindScenarioStatsMethodReference struct{}

type findScenarioStatsTempConfigPaths struct {
	ConfigPath string
	InfraPath  string
}

type findScenarioStatsOrdersReportingSink struct{}

func newFindScenarioStatsOrdersReportingSink() findScenarioStatsOrdersReportingSink {
	return findScenarioStatsOrdersReportingSink{}
}
func (findScenarioStatsOrdersReportingSink) SinkName() string {
	return "orders-sample-sink"
}
func (findScenarioStatsOrdersReportingSink) Init(loadstrike.LoadStrikeBaseContext, loadstrike.IConfiguration) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (findScenarioStatsOrdersReportingSink) Start(loadstrike.LoadStrikeSessionStartInfo) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (findScenarioStatsOrdersReportingSink) SaveRealtimeStats([]loadstrike.LoadStrikeScenarioStats) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (findScenarioStatsOrdersReportingSink) SaveRealtimeMetrics(loadstrike.LoadStrikeMetricStats) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (findScenarioStatsOrdersReportingSink) SaveRunResult(loadstrike.LoadStrikeRunResult) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (findScenarioStatsOrdersReportingSink) Stop() loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (findScenarioStatsOrdersReportingSink) Dispose() {}

type findScenarioStatsOrdersRuntimePolicy struct{}

func newFindScenarioStatsOrdersRuntimePolicy() findScenarioStatsOrdersRuntimePolicy {
	return findScenarioStatsOrdersRuntimePolicy{}
}
func (findScenarioStatsOrdersRuntimePolicy) ShouldRunScenario(string) loadstrike.LoadStrikeBoolTask {
	return loadstrike.TaskFromBool(true)
}
func (findScenarioStatsOrdersRuntimePolicy) BeforeScenario(string) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (findScenarioStatsOrdersRuntimePolicy) AfterScenario(string, loadstrike.LoadStrikeScenarioRuntime) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (findScenarioStatsOrdersRuntimePolicy) BeforeStep(string, string) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (findScenarioStatsOrdersRuntimePolicy) AfterStep(string, string, loadstrike.LoadStrikeReply) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}

type findScenarioStatsOrdersWorkerPlugin struct{}

func newFindScenarioStatsOrdersWorkerPlugin() findScenarioStatsOrdersWorkerPlugin {
	return findScenarioStatsOrdersWorkerPlugin{}
}
func (findScenarioStatsOrdersWorkerPlugin) PluginName() string {
	return "orders-sample-plugin"
}
func (findScenarioStatsOrdersWorkerPlugin) Init(loadstrike.LoadStrikeBaseContext, loadstrike.IConfiguration) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (findScenarioStatsOrdersWorkerPlugin) Start(loadstrike.LoadStrikeSessionStartInfo) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (findScenarioStatsOrdersWorkerPlugin) GetData(loadstrike.LoadStrikeRunResult) loadstrike.LoadStrikeValueTask[loadstrike.LoadStrikePluginData] {
	return loadstrike.TaskFromResult(loadstrike.CreatePluginData("orders-sample-plugin"))
}
func (findScenarioStatsOrdersWorkerPlugin) Stop() loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (findScenarioStatsOrdersWorkerPlugin) Dispose() loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}

func findScenarioStatsPerformOrderGetReply() loadstrike.LoadStrikeReply {
	return loadstrike.LoadStrikeResponse.Ok("200", int64(128), "ok", loadstrike.TimeSpan(3*time.Millisecond))
}

func findScenarioStatsExecuteOrderGet(context loadstrike.LoadStrikeScenarioContext) loadstrike.LoadStrikeReply {
	return loadstrike.LoadStrikeStep.Run("get-order", context, func(loadstrike.LoadStrikeScenarioContext) loadstrike.LoadStrikeReply {
		return findScenarioStatsPerformOrderGetReply()
	})
}

func findScenarioStatsExecuteOrderGetAsync(context loadstrike.LoadStrikeScenarioContext) loadstrike.LoadStrikeValueTask[loadstrike.LoadStrikeReply] {
	return loadstrike.TaskFromResult(findScenarioStatsExecuteOrderGet(context))
}

func findScenarioStatsBaselineScenario(names ...string) loadstrike.LoadStrikeScenario {
	name := "orders.get-by-id"
	if len(names) > 0 && names[0] != "" {
		name = names[0]
	}
	return loadstrike.CreateScenario(name, findScenarioStatsExecuteOrderGet).
		WithLoadSimulations(loadstrike.LoadStrikeSimulation.IterationsForConstant(1, 1)).
		WithoutWarmUp()
}

func findScenarioStatsBaseRunner() loadstrike.LoadStrikeRunner {
	return loadstrike.Create().
		AddScenario(findScenarioStatsBaselineScenario()).
		WithRunnerKey(findScenarioStatsRunnerKey).
		WithTestSuite("orders-reference").
		WithTestName("orders-get-by-id").
		WithoutReports()
}

func findScenarioStatsBaseContext() loadstrike.LoadStrikeContext {
	return findScenarioStatsBaseRunner().BuildContext()
}

func findScenarioStatsHttpSource() *loadstrike.EndpointSpec {
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

func findScenarioStatsHttpDestination() *loadstrike.EndpointSpec {
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

func findScenarioStatsTrackingConfiguration() *loadstrike.TrackingConfigurationSpec {
	return &loadstrike.TrackingConfigurationSpec{
		Source:                      findScenarioStatsHttpSource(),
		Destination:                 findScenarioStatsHttpDestination(),
		RunMode:                     "GenerateAndCorrelate",
		CorrelationTimeoutSeconds:   30,
		TimeoutSweepIntervalSeconds: 1,
		TimeoutBatchSize:            200,
		TimeoutCountsAsFailure:      true,
		MetricPrefix:                "orders_tracking",
		ExecuteOriginalScenarioRun:  false,
	}
}

func findScenarioStatsTrackedContext() loadstrike.LoadStrikeContext {
	return loadstrike.Create().
		AddScenario(findScenarioStatsBaselineScenario("orders.tracked").WithCrossPlatformTracking(findScenarioStatsTrackingConfiguration())).
		WithRunnerKey(findScenarioStatsRunnerKey).
		WithoutReports().
		BuildContext()
}

func findScenarioStatsBuildTrackingPayload() loadstrike.TrackingPayload {
	builder := loadstrike.TrackingPayloadBuilder{
		Headers:            map[string]string{"X-Correlation-Id": "ord-1"},
		ContentType:        "application/json",
		MessagePayloadType: "json",
	}
	builder.SetBody(map[string]any{"trackingId": "ord-1"})
	return builder.Build()
}

func findScenarioStatsRunResult() loadstrike.LoadStrikeRunResult {
	var result loadstrike.LoadStrikeRunResult
	return result
}

func findScenarioStatsScenarioStats() loadstrike.LoadStrikeScenarioStats {
	var stats loadstrike.LoadStrikeScenarioStats
	return stats
}

func findScenarioStatsWriteTempConfigFiles() findScenarioStatsTempConfigPaths {
	return findScenarioStatsTempConfigPaths{
		ConfigPath: "method-reference.loadstrike.config.json",
		InfraPath:  "method-reference.loadstrike.infra.json",
	}
}

        // Read data back from the run result surface with this helper.
func (reference FindScenarioStatsMethodReference) ReadRunResultsExample() any {
    return findScenarioStatsRunResult().FindScenarioStats("orders.get-by-id")
}

// Continue into the scenario or step stats surface after the first lookup.
func (reference FindScenarioStatsMethodReference) ReadDeeperResultSurfaceExample() any {
    return findScenarioStatsRunResult().GetScenarioStats("orders.get-by-id")
}
