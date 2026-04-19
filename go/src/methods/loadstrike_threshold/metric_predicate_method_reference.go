package loadstrike_threshold

        import (
        	"time"

        	loadstrike "loadstrike.com/sdk/go"
        )

        const metricPredicateRunnerKey = "runner_dummy_orders_reference"

type MetricPredicateMethodReference struct{}

type metricPredicateTempConfigPaths struct {
	ConfigPath string
	InfraPath  string
}

type metricPredicateOrdersReportingSink struct{}

func newMetricPredicateOrdersReportingSink() metricPredicateOrdersReportingSink {
	return metricPredicateOrdersReportingSink{}
}
func (metricPredicateOrdersReportingSink) SinkName() string {
	return "orders-sample-sink"
}
func (metricPredicateOrdersReportingSink) Init(loadstrike.LoadStrikeBaseContext, loadstrike.IConfiguration) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (metricPredicateOrdersReportingSink) Start(loadstrike.LoadStrikeSessionStartInfo) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (metricPredicateOrdersReportingSink) SaveRealtimeStats([]loadstrike.LoadStrikeScenarioStats) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (metricPredicateOrdersReportingSink) SaveRealtimeMetrics(loadstrike.LoadStrikeMetricStats) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (metricPredicateOrdersReportingSink) SaveRunResult(loadstrike.LoadStrikeRunResult) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (metricPredicateOrdersReportingSink) Stop() loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (metricPredicateOrdersReportingSink) Dispose() {}

type metricPredicateOrdersRuntimePolicy struct{}

func newMetricPredicateOrdersRuntimePolicy() metricPredicateOrdersRuntimePolicy {
	return metricPredicateOrdersRuntimePolicy{}
}
func (metricPredicateOrdersRuntimePolicy) ShouldRunScenario(string) loadstrike.LoadStrikeBoolTask {
	return loadstrike.TaskFromBool(true)
}
func (metricPredicateOrdersRuntimePolicy) BeforeScenario(string) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (metricPredicateOrdersRuntimePolicy) AfterScenario(string, loadstrike.LoadStrikeScenarioRuntime) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (metricPredicateOrdersRuntimePolicy) BeforeStep(string, string) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (metricPredicateOrdersRuntimePolicy) AfterStep(string, string, loadstrike.LoadStrikeReply) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}

type metricPredicateOrdersWorkerPlugin struct{}

func newMetricPredicateOrdersWorkerPlugin() metricPredicateOrdersWorkerPlugin {
	return metricPredicateOrdersWorkerPlugin{}
}
func (metricPredicateOrdersWorkerPlugin) PluginName() string {
	return "orders-sample-plugin"
}
func (metricPredicateOrdersWorkerPlugin) Init(loadstrike.LoadStrikeBaseContext, loadstrike.IConfiguration) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (metricPredicateOrdersWorkerPlugin) Start(loadstrike.LoadStrikeSessionStartInfo) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (metricPredicateOrdersWorkerPlugin) GetData(loadstrike.LoadStrikeRunResult) loadstrike.LoadStrikeValueTask[loadstrike.LoadStrikePluginData] {
	return loadstrike.TaskFromResult(loadstrike.CreatePluginData("orders-sample-plugin"))
}
func (metricPredicateOrdersWorkerPlugin) Stop() loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (metricPredicateOrdersWorkerPlugin) Dispose() loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}

func metricPredicatePerformOrderGetReply() loadstrike.LoadStrikeReply {
	return loadstrike.LoadStrikeResponse.Ok("200", int64(128), "ok", loadstrike.TimeSpan(3*time.Millisecond))
}

func metricPredicateExecuteOrderGet(context loadstrike.LoadStrikeScenarioContext) loadstrike.LoadStrikeReply {
	return loadstrike.LoadStrikeStep.Run("get-order", context, func(loadstrike.LoadStrikeScenarioContext) loadstrike.LoadStrikeReply {
		return metricPredicatePerformOrderGetReply()
	})
}

func metricPredicateExecuteOrderGetAsync(context loadstrike.LoadStrikeScenarioContext) loadstrike.LoadStrikeValueTask[loadstrike.LoadStrikeReply] {
	return loadstrike.TaskFromResult(metricPredicateExecuteOrderGet(context))
}

func metricPredicateBaselineScenario(names ...string) loadstrike.LoadStrikeScenario {
	name := "orders.get-by-id"
	if len(names) > 0 && names[0] != "" {
		name = names[0]
	}
	return loadstrike.CreateScenario(name, metricPredicateExecuteOrderGet).
		WithLoadSimulations(loadstrike.LoadStrikeSimulation.IterationsForConstant(1, 1)).
		WithoutWarmUp()
}

func metricPredicateBaseRunner() loadstrike.LoadStrikeRunner {
	return loadstrike.Create().
		AddScenario(metricPredicateBaselineScenario()).
		WithRunnerKey(metricPredicateRunnerKey).
		WithTestSuite("orders-reference").
		WithTestName("orders-get-by-id").
		WithoutReports()
}

func metricPredicateBaseContext() loadstrike.LoadStrikeContext {
	return metricPredicateBaseRunner().BuildContext()
}

func metricPredicateHttpSource() *loadstrike.EndpointSpec {
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

func metricPredicateHttpDestination() *loadstrike.EndpointSpec {
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

func metricPredicateTrackingConfiguration() *loadstrike.TrackingConfigurationSpec {
	return &loadstrike.TrackingConfigurationSpec{
		Source:                      metricPredicateHttpSource(),
		Destination:                 metricPredicateHttpDestination(),
		RunMode:                     "GenerateAndCorrelate",
		CorrelationTimeoutSeconds:   30,
		TimeoutSweepIntervalSeconds: 1,
		TimeoutBatchSize:            200,
		TimeoutCountsAsFailure:      true,
		MetricPrefix:                "orders_tracking",
		ExecuteOriginalScenarioRun:  false,
	}
}

func metricPredicateTrackedContext() loadstrike.LoadStrikeContext {
	return loadstrike.Create().
		AddScenario(metricPredicateBaselineScenario("orders.tracked").WithCrossPlatformTracking(metricPredicateTrackingConfiguration())).
		WithRunnerKey(metricPredicateRunnerKey).
		WithoutReports().
		BuildContext()
}

func metricPredicateBuildTrackingPayload() loadstrike.TrackingPayload {
	builder := loadstrike.TrackingPayloadBuilder{
		Headers:            map[string]string{"X-Correlation-Id": "ord-1"},
		ContentType:        "application/json",
		MessagePayloadType: "json",
	}
	builder.SetBody(map[string]any{"trackingId": "ord-1"})
	return builder.Build()
}

func metricPredicateRunResult() loadstrike.LoadStrikeRunResult {
	var result loadstrike.LoadStrikeRunResult
	return result
}

func metricPredicateScenarioStats() loadstrike.LoadStrikeScenarioStats {
	var stats loadstrike.LoadStrikeScenarioStats
	return stats
}

func metricPredicateWriteTempConfigFiles() metricPredicateTempConfigPaths {
	return metricPredicateTempConfigPaths{
		ConfigPath: "method-reference.loadstrike.config.json",
		InfraPath:  "method-reference.loadstrike.infra.json",
	}
}

        // Create the threshold from the public threshold helper.
func (reference MetricPredicateMethodReference) CreateThresholdExample() any {
    return loadstrike.LoadStrikeThreshold{}.MetricPredicate(func(stats loadstrike.LoadStrikeMetricStats) bool { return len(stats.Counters()) >= 1 })
}

// Attach the threshold to the baseline GET scenario.
func (reference MetricPredicateMethodReference) AttachThresholdToScenarioExample() any {
    return metricPredicateBaselineScenario().WithThresholds(loadstrike.LoadStrikeThreshold{}.MetricPredicate(func(stats loadstrike.LoadStrikeMetricStats) bool { return len(stats.Counters()) >= 1 }))
}
