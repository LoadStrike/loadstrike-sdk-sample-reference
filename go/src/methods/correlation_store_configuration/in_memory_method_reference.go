package correlation_store_configuration

        import (
        	"time"

        	loadstrike "loadstrike.com/sdk/go"
        )

        const inMemoryRunnerKey = "runner_dummy_orders_reference"

type InMemoryMethodReference struct{}

type inMemoryTempConfigPaths struct {
	ConfigPath string
	InfraPath  string
}

type inMemoryOrdersReportingSink struct{}

func newInMemoryOrdersReportingSink() inMemoryOrdersReportingSink {
	return inMemoryOrdersReportingSink{}
}
func (inMemoryOrdersReportingSink) SinkName() string {
	return "orders-sample-sink"
}
func (inMemoryOrdersReportingSink) Init(loadstrike.LoadStrikeBaseContext, loadstrike.IConfiguration) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (inMemoryOrdersReportingSink) Start(loadstrike.LoadStrikeSessionStartInfo) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (inMemoryOrdersReportingSink) SaveRealtimeStats([]loadstrike.LoadStrikeScenarioStats) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (inMemoryOrdersReportingSink) SaveRealtimeMetrics(loadstrike.LoadStrikeMetricStats) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (inMemoryOrdersReportingSink) SaveRunResult(loadstrike.LoadStrikeRunResult) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (inMemoryOrdersReportingSink) Stop() loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (inMemoryOrdersReportingSink) Dispose() {}

type inMemoryOrdersRuntimePolicy struct{}

func newInMemoryOrdersRuntimePolicy() inMemoryOrdersRuntimePolicy {
	return inMemoryOrdersRuntimePolicy{}
}
func (inMemoryOrdersRuntimePolicy) ShouldRunScenario(string) loadstrike.LoadStrikeBoolTask {
	return loadstrike.TaskFromBool(true)
}
func (inMemoryOrdersRuntimePolicy) BeforeScenario(string) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (inMemoryOrdersRuntimePolicy) AfterScenario(string, loadstrike.LoadStrikeScenarioRuntime) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (inMemoryOrdersRuntimePolicy) BeforeStep(string, string) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (inMemoryOrdersRuntimePolicy) AfterStep(string, string, loadstrike.LoadStrikeReply) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}

type inMemoryOrdersWorkerPlugin struct{}

func newInMemoryOrdersWorkerPlugin() inMemoryOrdersWorkerPlugin {
	return inMemoryOrdersWorkerPlugin{}
}
func (inMemoryOrdersWorkerPlugin) PluginName() string {
	return "orders-sample-plugin"
}
func (inMemoryOrdersWorkerPlugin) Init(loadstrike.LoadStrikeBaseContext, loadstrike.IConfiguration) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (inMemoryOrdersWorkerPlugin) Start(loadstrike.LoadStrikeSessionStartInfo) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (inMemoryOrdersWorkerPlugin) GetData(loadstrike.LoadStrikeRunResult) loadstrike.LoadStrikeValueTask[loadstrike.LoadStrikePluginData] {
	return loadstrike.TaskFromResult(loadstrike.CreatePluginData("orders-sample-plugin"))
}
func (inMemoryOrdersWorkerPlugin) Stop() loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (inMemoryOrdersWorkerPlugin) Dispose() loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}

func inMemoryPerformOrderGetReply() loadstrike.LoadStrikeReply {
	return loadstrike.LoadStrikeResponse.Ok("200", int64(128), "ok", loadstrike.TimeSpan(3*time.Millisecond))
}

func inMemoryExecuteOrderGet(context loadstrike.LoadStrikeScenarioContext) loadstrike.LoadStrikeReply {
	return loadstrike.LoadStrikeStep.Run("get-order", context, func(loadstrike.LoadStrikeScenarioContext) loadstrike.LoadStrikeReply {
		return inMemoryPerformOrderGetReply()
	})
}

func inMemoryExecuteOrderGetAsync(context loadstrike.LoadStrikeScenarioContext) loadstrike.LoadStrikeValueTask[loadstrike.LoadStrikeReply] {
	return loadstrike.TaskFromResult(inMemoryExecuteOrderGet(context))
}

func inMemoryBaselineScenario(names ...string) loadstrike.LoadStrikeScenario {
	name := "orders.get-by-id"
	if len(names) > 0 && names[0] != "" {
		name = names[0]
	}
	return loadstrike.CreateScenario(name, inMemoryExecuteOrderGet).
		WithLoadSimulations(loadstrike.LoadStrikeSimulation.IterationsForConstant(1, 1)).
		WithoutWarmUp()
}

func inMemoryBaseRunner() loadstrike.LoadStrikeRunner {
	return loadstrike.Create().
		AddScenario(inMemoryBaselineScenario()).
		WithRunnerKey(inMemoryRunnerKey).
		WithTestSuite("orders-reference").
		WithTestName("orders-get-by-id").
		WithoutReports()
}

func inMemoryBaseContext() loadstrike.LoadStrikeContext {
	return inMemoryBaseRunner().BuildContext()
}

func inMemoryHttpSource() *loadstrike.EndpointSpec {
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

func inMemoryHttpDestination() *loadstrike.EndpointSpec {
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

func inMemoryTrackingConfiguration() *loadstrike.TrackingConfigurationSpec {
	return &loadstrike.TrackingConfigurationSpec{
		Source:                      inMemoryHttpSource(),
		Destination:                 inMemoryHttpDestination(),
		RunMode:                     "GenerateAndCorrelate",
		CorrelationTimeoutSeconds:   30,
		TimeoutSweepIntervalSeconds: 1,
		TimeoutBatchSize:            200,
		TimeoutCountsAsFailure:      true,
		MetricPrefix:                "orders_tracking",
		ExecuteOriginalScenarioRun:  false,
	}
}

func inMemoryTrackedContext() loadstrike.LoadStrikeContext {
	return loadstrike.Create().
		AddScenario(inMemoryBaselineScenario("orders.tracked").WithCrossPlatformTracking(inMemoryTrackingConfiguration())).
		WithRunnerKey(inMemoryRunnerKey).
		WithoutReports().
		BuildContext()
}

func inMemoryBuildTrackingPayload() loadstrike.TrackingPayload {
	builder := loadstrike.TrackingPayloadBuilder{
		Headers:            map[string]string{"X-Correlation-Id": "ord-1"},
		ContentType:        "application/json",
		MessagePayloadType: "json",
	}
	builder.SetBody(map[string]any{"trackingId": "ord-1"})
	return builder.Build()
}

func inMemoryRunResult() loadstrike.LoadStrikeRunResult {
	var result loadstrike.LoadStrikeRunResult
	return result
}

func inMemoryScenarioStats() loadstrike.LoadStrikeScenarioStats {
	var stats loadstrike.LoadStrikeScenarioStats
	return stats
}

func inMemoryWriteTempConfigFiles() inMemoryTempConfigPaths {
	return inMemoryTempConfigPaths{
		ConfigPath: "method-reference.loadstrike.config.json",
		InfraPath:  "method-reference.loadstrike.infra.json",
	}
}

        // Call the correlation helper directly with a concrete value.
func (reference InMemoryMethodReference) CreateCorrelationExample() any {
    return loadstrike.CorrelationStoreConfiguration{}.InMemory()
}

// Show how the same helper fits into the tracked source/destination example.
func (reference InMemoryMethodReference) UseCorrelationExampleInTrackedFlow() any {
    return inMemoryTrackingConfiguration()
}
