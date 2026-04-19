package tracking_payload

        import (
        	"time"

        	loadstrike "loadstrike.com/sdk/go"
        )

        const getBodyAsUtf8RunnerKey = "runner_dummy_orders_reference"

type GetBodyAsUtf8MethodReference struct{}

type getBodyAsUtf8TempConfigPaths struct {
	ConfigPath string
	InfraPath  string
}

type getBodyAsUtf8OrdersReportingSink struct{}

func newGetBodyAsUtf8OrdersReportingSink() getBodyAsUtf8OrdersReportingSink {
	return getBodyAsUtf8OrdersReportingSink{}
}
func (getBodyAsUtf8OrdersReportingSink) SinkName() string {
	return "orders-sample-sink"
}
func (getBodyAsUtf8OrdersReportingSink) Init(loadstrike.LoadStrikeBaseContext, loadstrike.IConfiguration) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (getBodyAsUtf8OrdersReportingSink) Start(loadstrike.LoadStrikeSessionStartInfo) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (getBodyAsUtf8OrdersReportingSink) SaveRealtimeStats([]loadstrike.LoadStrikeScenarioStats) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (getBodyAsUtf8OrdersReportingSink) SaveRealtimeMetrics(loadstrike.LoadStrikeMetricStats) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (getBodyAsUtf8OrdersReportingSink) SaveRunResult(loadstrike.LoadStrikeRunResult) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (getBodyAsUtf8OrdersReportingSink) Stop() loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (getBodyAsUtf8OrdersReportingSink) Dispose() {}

type getBodyAsUtf8OrdersRuntimePolicy struct{}

func newGetBodyAsUtf8OrdersRuntimePolicy() getBodyAsUtf8OrdersRuntimePolicy {
	return getBodyAsUtf8OrdersRuntimePolicy{}
}
func (getBodyAsUtf8OrdersRuntimePolicy) ShouldRunScenario(string) loadstrike.LoadStrikeBoolTask {
	return loadstrike.TaskFromBool(true)
}
func (getBodyAsUtf8OrdersRuntimePolicy) BeforeScenario(string) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (getBodyAsUtf8OrdersRuntimePolicy) AfterScenario(string, loadstrike.LoadStrikeScenarioRuntime) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (getBodyAsUtf8OrdersRuntimePolicy) BeforeStep(string, string) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (getBodyAsUtf8OrdersRuntimePolicy) AfterStep(string, string, loadstrike.LoadStrikeReply) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}

type getBodyAsUtf8OrdersWorkerPlugin struct{}

func newGetBodyAsUtf8OrdersWorkerPlugin() getBodyAsUtf8OrdersWorkerPlugin {
	return getBodyAsUtf8OrdersWorkerPlugin{}
}
func (getBodyAsUtf8OrdersWorkerPlugin) PluginName() string {
	return "orders-sample-plugin"
}
func (getBodyAsUtf8OrdersWorkerPlugin) Init(loadstrike.LoadStrikeBaseContext, loadstrike.IConfiguration) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (getBodyAsUtf8OrdersWorkerPlugin) Start(loadstrike.LoadStrikeSessionStartInfo) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (getBodyAsUtf8OrdersWorkerPlugin) GetData(loadstrike.LoadStrikeRunResult) loadstrike.LoadStrikeValueTask[loadstrike.LoadStrikePluginData] {
	return loadstrike.TaskFromResult(loadstrike.CreatePluginData("orders-sample-plugin"))
}
func (getBodyAsUtf8OrdersWorkerPlugin) Stop() loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (getBodyAsUtf8OrdersWorkerPlugin) Dispose() loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}

func getBodyAsUtf8PerformOrderGetReply() loadstrike.LoadStrikeReply {
	return loadstrike.LoadStrikeResponse.Ok("200", int64(128), "ok", loadstrike.TimeSpan(3*time.Millisecond))
}

func getBodyAsUtf8ExecuteOrderGet(context loadstrike.LoadStrikeScenarioContext) loadstrike.LoadStrikeReply {
	return loadstrike.LoadStrikeStep.Run("get-order", context, func(loadstrike.LoadStrikeScenarioContext) loadstrike.LoadStrikeReply {
		return getBodyAsUtf8PerformOrderGetReply()
	})
}

func getBodyAsUtf8ExecuteOrderGetAsync(context loadstrike.LoadStrikeScenarioContext) loadstrike.LoadStrikeValueTask[loadstrike.LoadStrikeReply] {
	return loadstrike.TaskFromResult(getBodyAsUtf8ExecuteOrderGet(context))
}

func getBodyAsUtf8BaselineScenario(names ...string) loadstrike.LoadStrikeScenario {
	name := "orders.get-by-id"
	if len(names) > 0 && names[0] != "" {
		name = names[0]
	}
	return loadstrike.CreateScenario(name, getBodyAsUtf8ExecuteOrderGet).
		WithLoadSimulations(loadstrike.LoadStrikeSimulation.IterationsForConstant(1, 1)).
		WithoutWarmUp()
}

func getBodyAsUtf8BaseRunner() loadstrike.LoadStrikeRunner {
	return loadstrike.Create().
		AddScenario(getBodyAsUtf8BaselineScenario()).
		WithRunnerKey(getBodyAsUtf8RunnerKey).
		WithTestSuite("orders-reference").
		WithTestName("orders-get-by-id").
		WithoutReports()
}

func getBodyAsUtf8BaseContext() loadstrike.LoadStrikeContext {
	return getBodyAsUtf8BaseRunner().BuildContext()
}

func getBodyAsUtf8HttpSource() *loadstrike.EndpointSpec {
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

func getBodyAsUtf8HttpDestination() *loadstrike.EndpointSpec {
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

func getBodyAsUtf8TrackingConfiguration() *loadstrike.TrackingConfigurationSpec {
	return &loadstrike.TrackingConfigurationSpec{
		Source:                      getBodyAsUtf8HttpSource(),
		Destination:                 getBodyAsUtf8HttpDestination(),
		RunMode:                     "GenerateAndCorrelate",
		CorrelationTimeoutSeconds:   30,
		TimeoutSweepIntervalSeconds: 1,
		TimeoutBatchSize:            200,
		TimeoutCountsAsFailure:      true,
		MetricPrefix:                "orders_tracking",
		ExecuteOriginalScenarioRun:  false,
	}
}

func getBodyAsUtf8TrackedContext() loadstrike.LoadStrikeContext {
	return loadstrike.Create().
		AddScenario(getBodyAsUtf8BaselineScenario("orders.tracked").WithCrossPlatformTracking(getBodyAsUtf8TrackingConfiguration())).
		WithRunnerKey(getBodyAsUtf8RunnerKey).
		WithoutReports().
		BuildContext()
}

func getBodyAsUtf8BuildTrackingPayload() loadstrike.TrackingPayload {
	builder := loadstrike.TrackingPayloadBuilder{
		Headers:            map[string]string{"X-Correlation-Id": "ord-1"},
		ContentType:        "application/json",
		MessagePayloadType: "json",
	}
	builder.SetBody(map[string]any{"trackingId": "ord-1"})
	return builder.Build()
}

func getBodyAsUtf8RunResult() loadstrike.LoadStrikeRunResult {
	var result loadstrike.LoadStrikeRunResult
	return result
}

func getBodyAsUtf8ScenarioStats() loadstrike.LoadStrikeScenarioStats {
	var stats loadstrike.LoadStrikeScenarioStats
	return stats
}

func getBodyAsUtf8WriteTempConfigFiles() getBodyAsUtf8TempConfigPaths {
	return getBodyAsUtf8TempConfigPaths{
		ConfigPath: "method-reference.loadstrike.config.json",
		InfraPath:  "method-reference.loadstrike.infra.json",
	}
}

        // Call the correlation helper directly with a concrete value.
func (reference GetBodyAsUtf8MethodReference) CreateCorrelationExample() any {
    return getBodyAsUtf8BuildTrackingPayload().GetBodyAsUtf8()
}

// Show how the same helper fits into the tracked source/destination example.
func (reference GetBodyAsUtf8MethodReference) UseCorrelationExampleInTrackedFlow() any {
    return getBodyAsUtf8BuildTrackingPayload()
}
