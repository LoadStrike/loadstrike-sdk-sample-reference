package loadstrike_threshold

        import (
        	"time"

        	loadstrike "loadstrike.com/sdk/go"
        )

        const scenarioPredicateRunnerKey = "runner_dummy_orders_reference"

type ScenarioPredicateMethodReference struct{}

type scenarioPredicateTempConfigPaths struct {
	ConfigPath string
	InfraPath  string
}

type scenarioPredicateOrdersReportingSink struct{}

func newScenarioPredicateOrdersReportingSink() scenarioPredicateOrdersReportingSink {
	return scenarioPredicateOrdersReportingSink{}
}
func (scenarioPredicateOrdersReportingSink) SinkName() string {
	return "orders-sample-sink"
}
func (scenarioPredicateOrdersReportingSink) Init(loadstrike.LoadStrikeBaseContext, loadstrike.IConfiguration) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (scenarioPredicateOrdersReportingSink) Start(loadstrike.LoadStrikeSessionStartInfo) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (scenarioPredicateOrdersReportingSink) SaveRealtimeStats([]loadstrike.LoadStrikeScenarioStats) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (scenarioPredicateOrdersReportingSink) SaveRealtimeMetrics(loadstrike.LoadStrikeMetricStats) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (scenarioPredicateOrdersReportingSink) SaveRunResult(loadstrike.LoadStrikeRunResult) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (scenarioPredicateOrdersReportingSink) Stop() loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (scenarioPredicateOrdersReportingSink) Dispose() {}

type scenarioPredicateOrdersRuntimePolicy struct{}

func newScenarioPredicateOrdersRuntimePolicy() scenarioPredicateOrdersRuntimePolicy {
	return scenarioPredicateOrdersRuntimePolicy{}
}
func (scenarioPredicateOrdersRuntimePolicy) ShouldRunScenario(string) loadstrike.LoadStrikeBoolTask {
	return loadstrike.TaskFromBool(true)
}
func (scenarioPredicateOrdersRuntimePolicy) BeforeScenario(string) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (scenarioPredicateOrdersRuntimePolicy) AfterScenario(string, loadstrike.LoadStrikeScenarioRuntime) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (scenarioPredicateOrdersRuntimePolicy) BeforeStep(string, string) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (scenarioPredicateOrdersRuntimePolicy) AfterStep(string, string, loadstrike.LoadStrikeReply) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}

type scenarioPredicateOrdersWorkerPlugin struct{}

func newScenarioPredicateOrdersWorkerPlugin() scenarioPredicateOrdersWorkerPlugin {
	return scenarioPredicateOrdersWorkerPlugin{}
}
func (scenarioPredicateOrdersWorkerPlugin) PluginName() string {
	return "orders-sample-plugin"
}
func (scenarioPredicateOrdersWorkerPlugin) Init(loadstrike.LoadStrikeBaseContext, loadstrike.IConfiguration) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (scenarioPredicateOrdersWorkerPlugin) Start(loadstrike.LoadStrikeSessionStartInfo) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (scenarioPredicateOrdersWorkerPlugin) GetData(loadstrike.LoadStrikeRunResult) loadstrike.LoadStrikeValueTask[loadstrike.LoadStrikePluginData] {
	return loadstrike.TaskFromResult(loadstrike.CreatePluginData("orders-sample-plugin"))
}
func (scenarioPredicateOrdersWorkerPlugin) Stop() loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (scenarioPredicateOrdersWorkerPlugin) Dispose() loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}

func scenarioPredicatePerformOrderGetReply() loadstrike.LoadStrikeReply {
	return loadstrike.LoadStrikeResponse.Ok("200", int64(128), "ok", loadstrike.TimeSpan(3*time.Millisecond))
}

func scenarioPredicateExecuteOrderGet(context loadstrike.LoadStrikeScenarioContext) loadstrike.LoadStrikeReply {
	return loadstrike.LoadStrikeStep.Run("get-order", context, func(loadstrike.LoadStrikeScenarioContext) loadstrike.LoadStrikeReply {
		return scenarioPredicatePerformOrderGetReply()
	})
}

func scenarioPredicateExecuteOrderGetAsync(context loadstrike.LoadStrikeScenarioContext) loadstrike.LoadStrikeValueTask[loadstrike.LoadStrikeReply] {
	return loadstrike.TaskFromResult(scenarioPredicateExecuteOrderGet(context))
}

func scenarioPredicateBaselineScenario(names ...string) loadstrike.LoadStrikeScenario {
	name := "orders.get-by-id"
	if len(names) > 0 && names[0] != "" {
		name = names[0]
	}
	return loadstrike.CreateScenario(name, scenarioPredicateExecuteOrderGet).
		WithLoadSimulations(loadstrike.LoadStrikeSimulation.IterationsForConstant(1, 1)).
		WithoutWarmUp()
}

func scenarioPredicateBaseRunner() loadstrike.LoadStrikeRunner {
	return loadstrike.Create().
		AddScenario(scenarioPredicateBaselineScenario()).
		WithRunnerKey(scenarioPredicateRunnerKey).
		WithTestSuite("orders-reference").
		WithTestName("orders-get-by-id").
		WithoutReports()
}

func scenarioPredicateBaseContext() loadstrike.LoadStrikeContext {
	return scenarioPredicateBaseRunner().BuildContext()
}

func scenarioPredicateHttpSource() *loadstrike.EndpointSpec {
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

func scenarioPredicateHttpDestination() *loadstrike.EndpointSpec {
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

func scenarioPredicateTrackingConfiguration() *loadstrike.TrackingConfigurationSpec {
	return &loadstrike.TrackingConfigurationSpec{
		Source:                      scenarioPredicateHttpSource(),
		Destination:                 scenarioPredicateHttpDestination(),
		RunMode:                     "GenerateAndCorrelate",
		CorrelationTimeoutSeconds:   30,
		TimeoutSweepIntervalSeconds: 1,
		TimeoutBatchSize:            200,
		TimeoutCountsAsFailure:      true,
		MetricPrefix:                "orders_tracking",
		ExecuteOriginalScenarioRun:  false,
	}
}

func scenarioPredicateTrackedContext() loadstrike.LoadStrikeContext {
	return loadstrike.Create().
		AddScenario(scenarioPredicateBaselineScenario("orders.tracked").WithCrossPlatformTracking(scenarioPredicateTrackingConfiguration())).
		WithRunnerKey(scenarioPredicateRunnerKey).
		WithoutReports().
		BuildContext()
}

func scenarioPredicateBuildTrackingPayload() loadstrike.TrackingPayload {
	builder := loadstrike.TrackingPayloadBuilder{
		Headers:            map[string]string{"X-Correlation-Id": "ord-1"},
		ContentType:        "application/json",
		MessagePayloadType: "json",
	}
	builder.SetBody(map[string]any{"trackingId": "ord-1"})
	return builder.Build()
}

func scenarioPredicateRunResult() loadstrike.LoadStrikeRunResult {
	var result loadstrike.LoadStrikeRunResult
	return result
}

func scenarioPredicateScenarioStats() loadstrike.LoadStrikeScenarioStats {
	var stats loadstrike.LoadStrikeScenarioStats
	return stats
}

func scenarioPredicateWriteTempConfigFiles() scenarioPredicateTempConfigPaths {
	return scenarioPredicateTempConfigPaths{
		ConfigPath: "method-reference.loadstrike.config.json",
		InfraPath:  "method-reference.loadstrike.infra.json",
	}
}

        // Create the threshold from the public threshold helper.
func (reference ScenarioPredicateMethodReference) CreateThresholdExample() any {
    return loadstrike.LoadStrikeThreshold{}.ScenarioPredicate(func(stats loadstrike.LoadStrikeScenarioStats) bool { return stats.AllRequestCount >= 1 })
}

// Attach the threshold to the baseline GET scenario.
func (reference ScenarioPredicateMethodReference) AttachThresholdToScenarioExample() any {
    return scenarioPredicateBaselineScenario().WithThresholds(loadstrike.LoadStrikeThreshold{}.ScenarioPredicate(func(stats loadstrike.LoadStrikeScenarioStats) bool { return stats.AllRequestCount >= 1 }))
}
