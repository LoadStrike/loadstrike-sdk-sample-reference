package loadstrike_threshold

        import (
        	"time"

        	loadstrike "loadstrike.com/sdk/go"
        )

        const createMetricRunnerKey = "runner_dummy_orders_reference"

type CreateMetricMethodReference struct{}

type createMetricTempConfigPaths struct {
	ConfigPath string
	InfraPath  string
}

type createMetricOrdersReportingSink struct{}

func newCreateMetricOrdersReportingSink() createMetricOrdersReportingSink {
	return createMetricOrdersReportingSink{}
}
func (createMetricOrdersReportingSink) SinkName() string {
	return "orders-sample-sink"
}
func (createMetricOrdersReportingSink) Init(loadstrike.LoadStrikeBaseContext, loadstrike.IConfiguration) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (createMetricOrdersReportingSink) Start(loadstrike.LoadStrikeSessionStartInfo) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (createMetricOrdersReportingSink) SaveRealtimeStats([]loadstrike.LoadStrikeScenarioStats) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (createMetricOrdersReportingSink) SaveRealtimeMetrics(loadstrike.LoadStrikeMetricStats) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (createMetricOrdersReportingSink) SaveRunResult(loadstrike.LoadStrikeRunResult) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (createMetricOrdersReportingSink) Stop() loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (createMetricOrdersReportingSink) Dispose() {}

type createMetricOrdersRuntimePolicy struct{}

func newCreateMetricOrdersRuntimePolicy() createMetricOrdersRuntimePolicy {
	return createMetricOrdersRuntimePolicy{}
}
func (createMetricOrdersRuntimePolicy) ShouldRunScenario(string) loadstrike.LoadStrikeBoolTask {
	return loadstrike.TaskFromBool(true)
}
func (createMetricOrdersRuntimePolicy) BeforeScenario(string) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (createMetricOrdersRuntimePolicy) AfterScenario(string, loadstrike.LoadStrikeScenarioRuntime) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (createMetricOrdersRuntimePolicy) BeforeStep(string, string) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (createMetricOrdersRuntimePolicy) AfterStep(string, string, loadstrike.LoadStrikeReply) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}

type createMetricOrdersWorkerPlugin struct{}

func newCreateMetricOrdersWorkerPlugin() createMetricOrdersWorkerPlugin {
	return createMetricOrdersWorkerPlugin{}
}
func (createMetricOrdersWorkerPlugin) PluginName() string {
	return "orders-sample-plugin"
}
func (createMetricOrdersWorkerPlugin) Init(loadstrike.LoadStrikeBaseContext, loadstrike.IConfiguration) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (createMetricOrdersWorkerPlugin) Start(loadstrike.LoadStrikeSessionStartInfo) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (createMetricOrdersWorkerPlugin) GetData(loadstrike.LoadStrikeRunResult) loadstrike.LoadStrikeValueTask[loadstrike.LoadStrikePluginData] {
	return loadstrike.TaskFromResult(loadstrike.CreatePluginData("orders-sample-plugin"))
}
func (createMetricOrdersWorkerPlugin) Stop() loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (createMetricOrdersWorkerPlugin) Dispose() loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}

func createMetricPerformOrderGetReply() loadstrike.LoadStrikeReply {
	return loadstrike.LoadStrikeResponse.Ok("200", int64(128), "ok", loadstrike.TimeSpan(3*time.Millisecond))
}

func createMetricExecuteOrderGet(context loadstrike.LoadStrikeScenarioContext) loadstrike.LoadStrikeReply {
	return loadstrike.LoadStrikeStep.Run("get-order", context, func(loadstrike.LoadStrikeScenarioContext) loadstrike.LoadStrikeReply {
		return createMetricPerformOrderGetReply()
	})
}

func createMetricExecuteOrderGetAsync(context loadstrike.LoadStrikeScenarioContext) loadstrike.LoadStrikeValueTask[loadstrike.LoadStrikeReply] {
	return loadstrike.TaskFromResult(createMetricExecuteOrderGet(context))
}

func createMetricBaselineScenario(names ...string) loadstrike.LoadStrikeScenario {
	name := "orders.get-by-id"
	if len(names) > 0 && names[0] != "" {
		name = names[0]
	}
	return loadstrike.CreateScenario(name, createMetricExecuteOrderGet).
		WithLoadSimulations(loadstrike.LoadStrikeSimulation.IterationsForConstant(1, 1)).
		WithoutWarmUp()
}

func createMetricBaseRunner() loadstrike.LoadStrikeRunner {
	return loadstrike.Create().
		AddScenario(createMetricBaselineScenario()).
		WithRunnerKey(createMetricRunnerKey).
		WithTestSuite("orders-reference").
		WithTestName("orders-get-by-id").
		WithoutReports()
}

func createMetricBaseContext() loadstrike.LoadStrikeContext {
	return createMetricBaseRunner().BuildContext()
}

func createMetricHttpSource() *loadstrike.EndpointSpec {
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

func createMetricHttpDestination() *loadstrike.EndpointSpec {
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

func createMetricTrackingConfiguration() *loadstrike.TrackingConfigurationSpec {
	return &loadstrike.TrackingConfigurationSpec{
		Source:                      createMetricHttpSource(),
		Destination:                 createMetricHttpDestination(),
		RunMode:                     "GenerateAndCorrelate",
		CorrelationTimeoutSeconds:   30,
		TimeoutSweepIntervalSeconds: 1,
		TimeoutBatchSize:            200,
		TimeoutCountsAsFailure:      true,
		MetricPrefix:                "orders_tracking",
		ExecuteOriginalScenarioRun:  false,
	}
}

func createMetricTrackedContext() loadstrike.LoadStrikeContext {
	return loadstrike.Create().
		AddScenario(createMetricBaselineScenario("orders.tracked").WithCrossPlatformTracking(createMetricTrackingConfiguration())).
		WithRunnerKey(createMetricRunnerKey).
		WithoutReports().
		BuildContext()
}

func createMetricBuildTrackingPayload() loadstrike.TrackingPayload {
	builder := loadstrike.TrackingPayloadBuilder{
		Headers:            map[string]string{"X-Correlation-Id": "ord-1"},
		ContentType:        "application/json",
		MessagePayloadType: "json",
	}
	builder.SetBody(map[string]any{"trackingId": "ord-1"})
	return builder.Build()
}

func createMetricRunResult() loadstrike.LoadStrikeRunResult {
	var result loadstrike.LoadStrikeRunResult
	return result
}

func createMetricScenarioStats() loadstrike.LoadStrikeScenarioStats {
	var stats loadstrike.LoadStrikeScenarioStats
	return stats
}

func createMetricWriteTempConfigFiles() createMetricTempConfigPaths {
	return createMetricTempConfigPaths{
		ConfigPath: "method-reference.loadstrike.config.json",
		InfraPath:  "method-reference.loadstrike.infra.json",
	}
}

        // Create the threshold from the public threshold helper.
func (reference CreateMetricMethodReference) CreateThresholdExample() any {
    return loadstrike.LoadStrikeThreshold{}.CreateMetric("orders_total", "gte", 1)
}

// Attach the threshold to the baseline GET scenario.
func (reference CreateMetricMethodReference) AttachThresholdToScenarioExample() any {
    return createMetricBaselineScenario().WithThresholds(loadstrike.LoadStrikeThreshold{}.CreateMetric("orders_total", "gte", 1))
}
