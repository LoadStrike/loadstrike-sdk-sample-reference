package loadstrike_threshold

        import (
        	"time"

        	loadstrike "loadstrike.com/sdk/go"
        )

        const stepPredicateRunnerKey = "runner_dummy_orders_reference"

type StepPredicateMethodReference struct{}

type stepPredicateTempConfigPaths struct {
	ConfigPath string
	InfraPath  string
}

type stepPredicateOrdersReportingSink struct{}

func newStepPredicateOrdersReportingSink() stepPredicateOrdersReportingSink {
	return stepPredicateOrdersReportingSink{}
}
func (stepPredicateOrdersReportingSink) SinkName() string {
	return "orders-sample-sink"
}
func (stepPredicateOrdersReportingSink) Init(loadstrike.LoadStrikeBaseContext, loadstrike.IConfiguration) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (stepPredicateOrdersReportingSink) Start(loadstrike.LoadStrikeSessionStartInfo) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (stepPredicateOrdersReportingSink) SaveRealtimeStats([]loadstrike.LoadStrikeScenarioStats) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (stepPredicateOrdersReportingSink) SaveRealtimeMetrics(loadstrike.LoadStrikeMetricStats) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (stepPredicateOrdersReportingSink) SaveRunResult(loadstrike.LoadStrikeRunResult) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (stepPredicateOrdersReportingSink) Stop() loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (stepPredicateOrdersReportingSink) Dispose() {}

type stepPredicateOrdersRuntimePolicy struct{}

func newStepPredicateOrdersRuntimePolicy() stepPredicateOrdersRuntimePolicy {
	return stepPredicateOrdersRuntimePolicy{}
}
func (stepPredicateOrdersRuntimePolicy) ShouldRunScenario(string) loadstrike.LoadStrikeBoolTask {
	return loadstrike.TaskFromBool(true)
}
func (stepPredicateOrdersRuntimePolicy) BeforeScenario(string) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (stepPredicateOrdersRuntimePolicy) AfterScenario(string, loadstrike.LoadStrikeScenarioRuntime) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (stepPredicateOrdersRuntimePolicy) BeforeStep(string, string) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (stepPredicateOrdersRuntimePolicy) AfterStep(string, string, loadstrike.LoadStrikeReply) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}

type stepPredicateOrdersWorkerPlugin struct{}

func newStepPredicateOrdersWorkerPlugin() stepPredicateOrdersWorkerPlugin {
	return stepPredicateOrdersWorkerPlugin{}
}
func (stepPredicateOrdersWorkerPlugin) PluginName() string {
	return "orders-sample-plugin"
}
func (stepPredicateOrdersWorkerPlugin) Init(loadstrike.LoadStrikeBaseContext, loadstrike.IConfiguration) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (stepPredicateOrdersWorkerPlugin) Start(loadstrike.LoadStrikeSessionStartInfo) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (stepPredicateOrdersWorkerPlugin) GetData(loadstrike.LoadStrikeRunResult) loadstrike.LoadStrikeValueTask[loadstrike.LoadStrikePluginData] {
	return loadstrike.TaskFromResult(loadstrike.CreatePluginData("orders-sample-plugin"))
}
func (stepPredicateOrdersWorkerPlugin) Stop() loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (stepPredicateOrdersWorkerPlugin) Dispose() loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}

func stepPredicatePerformOrderGetReply() loadstrike.LoadStrikeReply {
	return loadstrike.LoadStrikeResponse.Ok("200", int64(128), "ok", loadstrike.TimeSpan(3*time.Millisecond))
}

func stepPredicateExecuteOrderGet(context loadstrike.LoadStrikeScenarioContext) loadstrike.LoadStrikeReply {
	return loadstrike.LoadStrikeStep.Run("get-order", context, func(loadstrike.LoadStrikeScenarioContext) loadstrike.LoadStrikeReply {
		return stepPredicatePerformOrderGetReply()
	})
}

func stepPredicateExecuteOrderGetAsync(context loadstrike.LoadStrikeScenarioContext) loadstrike.LoadStrikeValueTask[loadstrike.LoadStrikeReply] {
	return loadstrike.TaskFromResult(stepPredicateExecuteOrderGet(context))
}

func stepPredicateBaselineScenario(names ...string) loadstrike.LoadStrikeScenario {
	name := "orders.get-by-id"
	if len(names) > 0 && names[0] != "" {
		name = names[0]
	}
	return loadstrike.CreateScenario(name, stepPredicateExecuteOrderGet).
		WithLoadSimulations(loadstrike.LoadStrikeSimulation.IterationsForConstant(1, 1)).
		WithoutWarmUp()
}

func stepPredicateBaseRunner() loadstrike.LoadStrikeRunner {
	return loadstrike.Create().
		AddScenario(stepPredicateBaselineScenario()).
		WithRunnerKey(stepPredicateRunnerKey).
		WithTestSuite("orders-reference").
		WithTestName("orders-get-by-id").
		WithoutReports()
}

func stepPredicateBaseContext() loadstrike.LoadStrikeContext {
	return stepPredicateBaseRunner().BuildContext()
}

func stepPredicateHttpSource() *loadstrike.EndpointSpec {
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

func stepPredicateHttpDestination() *loadstrike.EndpointSpec {
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

func stepPredicateTrackingConfiguration() *loadstrike.TrackingConfigurationSpec {
	return &loadstrike.TrackingConfigurationSpec{
		Source:                      stepPredicateHttpSource(),
		Destination:                 stepPredicateHttpDestination(),
		RunMode:                     "GenerateAndCorrelate",
		CorrelationTimeoutSeconds:   30,
		TimeoutSweepIntervalSeconds: 1,
		TimeoutBatchSize:            200,
		TimeoutCountsAsFailure:      true,
		MetricPrefix:                "orders_tracking",
		ExecuteOriginalScenarioRun:  false,
	}
}

func stepPredicateTrackedContext() loadstrike.LoadStrikeContext {
	return loadstrike.Create().
		AddScenario(stepPredicateBaselineScenario("orders.tracked").WithCrossPlatformTracking(stepPredicateTrackingConfiguration())).
		WithRunnerKey(stepPredicateRunnerKey).
		WithoutReports().
		BuildContext()
}

func stepPredicateBuildTrackingPayload() loadstrike.TrackingPayload {
	builder := loadstrike.TrackingPayloadBuilder{
		Headers:            map[string]string{"X-Correlation-Id": "ord-1"},
		ContentType:        "application/json",
		MessagePayloadType: "json",
	}
	builder.SetBody(map[string]any{"trackingId": "ord-1"})
	return builder.Build()
}

func stepPredicateRunResult() loadstrike.LoadStrikeRunResult {
	var result loadstrike.LoadStrikeRunResult
	return result
}

func stepPredicateScenarioStats() loadstrike.LoadStrikeScenarioStats {
	var stats loadstrike.LoadStrikeScenarioStats
	return stats
}

func stepPredicateWriteTempConfigFiles() stepPredicateTempConfigPaths {
	return stepPredicateTempConfigPaths{
		ConfigPath: "method-reference.loadstrike.config.json",
		InfraPath:  "method-reference.loadstrike.infra.json",
	}
}

        // Create the threshold from the public threshold helper.
func (reference StepPredicateMethodReference) CreateThresholdExample() any {
    return loadstrike.LoadStrikeThreshold{}.StepPredicate("get-order", func(stats loadstrike.LoadStrikeDetailedStepStats) bool { return stats.OkCount() >= 1 })
}

// Attach the threshold to the baseline GET scenario.
func (reference StepPredicateMethodReference) AttachThresholdToScenarioExample() any {
    return stepPredicateBaselineScenario().WithThresholds(loadstrike.LoadStrikeThreshold{}.StepPredicate("get-order", func(stats loadstrike.LoadStrikeDetailedStepStats) bool { return stats.OkCount() >= 1 }))
}
