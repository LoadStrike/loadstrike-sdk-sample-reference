package loadstrike_metric

        import (
        	"time"

        	loadstrike "loadstrike.com/sdk/go"
        )

        const createGaugeRunnerKey = "runner_dummy_orders_reference"

type CreateGaugeMethodReference struct{}

type createGaugeTempConfigPaths struct {
	ConfigPath string
	InfraPath  string
}

type createGaugeOrdersReportingSink struct{}

func newCreateGaugeOrdersReportingSink() createGaugeOrdersReportingSink {
	return createGaugeOrdersReportingSink{}
}
func (createGaugeOrdersReportingSink) SinkName() string {
	return "orders-sample-sink"
}
func (createGaugeOrdersReportingSink) Init(loadstrike.LoadStrikeBaseContext, loadstrike.IConfiguration) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (createGaugeOrdersReportingSink) Start(loadstrike.LoadStrikeSessionStartInfo) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (createGaugeOrdersReportingSink) SaveRealtimeStats([]loadstrike.LoadStrikeScenarioStats) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (createGaugeOrdersReportingSink) SaveRealtimeMetrics(loadstrike.LoadStrikeMetricStats) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (createGaugeOrdersReportingSink) SaveRunResult(loadstrike.LoadStrikeRunResult) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (createGaugeOrdersReportingSink) Stop() loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (createGaugeOrdersReportingSink) Dispose() {}

type createGaugeOrdersRuntimePolicy struct{}

func newCreateGaugeOrdersRuntimePolicy() createGaugeOrdersRuntimePolicy {
	return createGaugeOrdersRuntimePolicy{}
}
func (createGaugeOrdersRuntimePolicy) ShouldRunScenario(string) loadstrike.LoadStrikeBoolTask {
	return loadstrike.TaskFromBool(true)
}
func (createGaugeOrdersRuntimePolicy) BeforeScenario(string) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (createGaugeOrdersRuntimePolicy) AfterScenario(string, loadstrike.LoadStrikeScenarioRuntime) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (createGaugeOrdersRuntimePolicy) BeforeStep(string, string) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (createGaugeOrdersRuntimePolicy) AfterStep(string, string, loadstrike.LoadStrikeReply) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}

type createGaugeOrdersWorkerPlugin struct{}

func newCreateGaugeOrdersWorkerPlugin() createGaugeOrdersWorkerPlugin {
	return createGaugeOrdersWorkerPlugin{}
}
func (createGaugeOrdersWorkerPlugin) PluginName() string {
	return "orders-sample-plugin"
}
func (createGaugeOrdersWorkerPlugin) Init(loadstrike.LoadStrikeBaseContext, loadstrike.IConfiguration) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (createGaugeOrdersWorkerPlugin) Start(loadstrike.LoadStrikeSessionStartInfo) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (createGaugeOrdersWorkerPlugin) GetData(loadstrike.LoadStrikeRunResult) loadstrike.LoadStrikeValueTask[loadstrike.LoadStrikePluginData] {
	return loadstrike.TaskFromResult(loadstrike.CreatePluginData("orders-sample-plugin"))
}
func (createGaugeOrdersWorkerPlugin) Stop() loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (createGaugeOrdersWorkerPlugin) Dispose() loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}

func createGaugePerformOrderGetReply() loadstrike.LoadStrikeReply {
	return loadstrike.LoadStrikeResponse.Ok("200", int64(128), "ok", loadstrike.TimeSpan(3*time.Millisecond))
}

func createGaugeExecuteOrderGet(context loadstrike.LoadStrikeScenarioContext) loadstrike.LoadStrikeReply {
	return loadstrike.LoadStrikeStep.Run("get-order", context, func(loadstrike.LoadStrikeScenarioContext) loadstrike.LoadStrikeReply {
		return createGaugePerformOrderGetReply()
	})
}

func createGaugeExecuteOrderGetAsync(context loadstrike.LoadStrikeScenarioContext) loadstrike.LoadStrikeValueTask[loadstrike.LoadStrikeReply] {
	return loadstrike.TaskFromResult(createGaugeExecuteOrderGet(context))
}

func createGaugeBaselineScenario(names ...string) loadstrike.LoadStrikeScenario {
	name := "orders.get-by-id"
	if len(names) > 0 && names[0] != "" {
		name = names[0]
	}
	return loadstrike.CreateScenario(name, createGaugeExecuteOrderGet).
		WithLoadSimulations(loadstrike.LoadStrikeSimulation.IterationsForConstant(1, 1)).
		WithoutWarmUp()
}

func createGaugeBaseRunner() loadstrike.LoadStrikeRunner {
	return loadstrike.Create().
		AddScenario(createGaugeBaselineScenario()).
		WithRunnerKey(createGaugeRunnerKey).
		WithTestSuite("orders-reference").
		WithTestName("orders-get-by-id").
		WithoutReports()
}

func createGaugeBaseContext() loadstrike.LoadStrikeContext {
	return createGaugeBaseRunner().BuildContext()
}

func createGaugeHttpSource() *loadstrike.EndpointSpec {
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

func createGaugeHttpDestination() *loadstrike.EndpointSpec {
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

func createGaugeTrackingConfiguration() *loadstrike.TrackingConfigurationSpec {
	return &loadstrike.TrackingConfigurationSpec{
		Source:                      createGaugeHttpSource(),
		Destination:                 createGaugeHttpDestination(),
		RunMode:                     "GenerateAndCorrelate",
		CorrelationTimeoutSeconds:   30,
		TimeoutSweepIntervalSeconds: 1,
		TimeoutBatchSize:            200,
		TimeoutCountsAsFailure:      true,
		MetricPrefix:                "orders_tracking",
		ExecuteOriginalScenarioRun:  false,
	}
}

func createGaugeTrackedContext() loadstrike.LoadStrikeContext {
	return loadstrike.Create().
		AddScenario(createGaugeBaselineScenario("orders.tracked").WithCrossPlatformTracking(createGaugeTrackingConfiguration())).
		WithRunnerKey(createGaugeRunnerKey).
		WithoutReports().
		BuildContext()
}

func createGaugeBuildTrackingPayload() loadstrike.TrackingPayload {
	builder := loadstrike.TrackingPayloadBuilder{
		Headers:            map[string]string{"X-Correlation-Id": "ord-1"},
		ContentType:        "application/json",
		MessagePayloadType: "json",
	}
	builder.SetBody(map[string]any{"trackingId": "ord-1"})
	return builder.Build()
}

func createGaugeRunResult() loadstrike.LoadStrikeRunResult {
	var result loadstrike.LoadStrikeRunResult
	return result
}

func createGaugeScenarioStats() loadstrike.LoadStrikeScenarioStats {
	var stats loadstrike.LoadStrikeScenarioStats
	return stats
}

func createGaugeWriteTempConfigFiles() createGaugeTempConfigPaths {
	return createGaugeTempConfigPaths{
		ConfigPath: "method-reference.loadstrike.config.json",
		InfraPath:  "method-reference.loadstrike.infra.json",
	}
}

        // Create the metric directly from the public metric helper.
func (reference CreateGaugeMethodReference) CreateMetricExample() any {
    return loadstrike.Metric.CreateGauge("orders_latency", "ms")
}

// Register the metric during scenario init so the GET scenario can report it.
        func (reference CreateGaugeMethodReference) RegisterMetricDuringInitExample() any {
            metric := loadstrike.Metric.CreateGauge("orders_latency", "ms")
scenario := loadstrike.CreateScenario("orders.metric", createGaugeExecuteOrderGet).WithInit(func(context loadstrike.LoadStrikeScenarioInitContext) error { context.RegisterMetric(metric); return nil })
return scenario
        }
