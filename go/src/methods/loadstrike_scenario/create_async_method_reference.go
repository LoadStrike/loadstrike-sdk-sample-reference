package loadstrike_scenario

        import (
        	"time"

        	loadstrike "loadstrike.com/sdk/go"
        )

        const createAsyncRunnerKey = "runner_dummy_orders_reference"

type CreateAsyncMethodReference struct{}

type createAsyncTempConfigPaths struct {
	ConfigPath string
	InfraPath  string
}

type createAsyncOrdersReportingSink struct{}

func newCreateAsyncOrdersReportingSink() createAsyncOrdersReportingSink {
	return createAsyncOrdersReportingSink{}
}
func (createAsyncOrdersReportingSink) SinkName() string {
	return "orders-sample-sink"
}
func (createAsyncOrdersReportingSink) Init(loadstrike.LoadStrikeBaseContext, loadstrike.IConfiguration) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (createAsyncOrdersReportingSink) Start(loadstrike.LoadStrikeSessionStartInfo) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (createAsyncOrdersReportingSink) SaveRealtimeStats([]loadstrike.LoadStrikeScenarioStats) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (createAsyncOrdersReportingSink) SaveRealtimeMetrics(loadstrike.LoadStrikeMetricStats) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (createAsyncOrdersReportingSink) SaveRunResult(loadstrike.LoadStrikeRunResult) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (createAsyncOrdersReportingSink) Stop() loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (createAsyncOrdersReportingSink) Dispose() {}

type createAsyncOrdersRuntimePolicy struct{}

func newCreateAsyncOrdersRuntimePolicy() createAsyncOrdersRuntimePolicy {
	return createAsyncOrdersRuntimePolicy{}
}
func (createAsyncOrdersRuntimePolicy) ShouldRunScenario(string) loadstrike.LoadStrikeBoolTask {
	return loadstrike.TaskFromBool(true)
}
func (createAsyncOrdersRuntimePolicy) BeforeScenario(string) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (createAsyncOrdersRuntimePolicy) AfterScenario(string, loadstrike.LoadStrikeScenarioRuntime) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (createAsyncOrdersRuntimePolicy) BeforeStep(string, string) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (createAsyncOrdersRuntimePolicy) AfterStep(string, string, loadstrike.LoadStrikeReply) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}

type createAsyncOrdersWorkerPlugin struct{}

func newCreateAsyncOrdersWorkerPlugin() createAsyncOrdersWorkerPlugin {
	return createAsyncOrdersWorkerPlugin{}
}
func (createAsyncOrdersWorkerPlugin) PluginName() string {
	return "orders-sample-plugin"
}
func (createAsyncOrdersWorkerPlugin) Init(loadstrike.LoadStrikeBaseContext, loadstrike.IConfiguration) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (createAsyncOrdersWorkerPlugin) Start(loadstrike.LoadStrikeSessionStartInfo) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (createAsyncOrdersWorkerPlugin) GetData(loadstrike.LoadStrikeRunResult) loadstrike.LoadStrikeValueTask[loadstrike.LoadStrikePluginData] {
	return loadstrike.TaskFromResult(loadstrike.CreatePluginData("orders-sample-plugin"))
}
func (createAsyncOrdersWorkerPlugin) Stop() loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (createAsyncOrdersWorkerPlugin) Dispose() loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}

func createAsyncPerformOrderGetReply() loadstrike.LoadStrikeReply {
	return loadstrike.LoadStrikeResponse.Ok("200", int64(128), "ok", loadstrike.TimeSpan(3*time.Millisecond))
}

func createAsyncExecuteOrderGet(context loadstrike.LoadStrikeScenarioContext) loadstrike.LoadStrikeReply {
	return loadstrike.LoadStrikeStep.Run("get-order", context, func(loadstrike.LoadStrikeScenarioContext) loadstrike.LoadStrikeReply {
		return createAsyncPerformOrderGetReply()
	})
}

func createAsyncExecuteOrderGetAsync(context loadstrike.LoadStrikeScenarioContext) loadstrike.LoadStrikeValueTask[loadstrike.LoadStrikeReply] {
	return loadstrike.TaskFromResult(createAsyncExecuteOrderGet(context))
}

func createAsyncBaselineScenario(names ...string) loadstrike.LoadStrikeScenario {
	name := "orders.get-by-id"
	if len(names) > 0 && names[0] != "" {
		name = names[0]
	}
	return loadstrike.CreateScenario(name, createAsyncExecuteOrderGet).
		WithLoadSimulations(loadstrike.LoadStrikeSimulation.IterationsForConstant(1, 1)).
		WithoutWarmUp()
}

func createAsyncBaseRunner() loadstrike.LoadStrikeRunner {
	return loadstrike.Create().
		AddScenario(createAsyncBaselineScenario()).
		WithRunnerKey(createAsyncRunnerKey).
		WithTestSuite("orders-reference").
		WithTestName("orders-get-by-id").
		WithoutReports()
}

func createAsyncBaseContext() loadstrike.LoadStrikeContext {
	return createAsyncBaseRunner().BuildContext()
}

func createAsyncHttpSource() *loadstrike.EndpointSpec {
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

func createAsyncHttpDestination() *loadstrike.EndpointSpec {
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

func createAsyncTrackingConfiguration() *loadstrike.TrackingConfigurationSpec {
	return &loadstrike.TrackingConfigurationSpec{
		Source:                      createAsyncHttpSource(),
		Destination:                 createAsyncHttpDestination(),
		RunMode:                     "GenerateAndCorrelate",
		CorrelationTimeoutSeconds:   30,
		TimeoutSweepIntervalSeconds: 1,
		TimeoutBatchSize:            200,
		TimeoutCountsAsFailure:      true,
		MetricPrefix:                "orders_tracking",
		ExecuteOriginalScenarioRun:  false,
	}
}

func createAsyncTrackedContext() loadstrike.LoadStrikeContext {
	return loadstrike.Create().
		AddScenario(createAsyncBaselineScenario("orders.tracked").WithCrossPlatformTracking(createAsyncTrackingConfiguration())).
		WithRunnerKey(createAsyncRunnerKey).
		WithoutReports().
		BuildContext()
}

func createAsyncBuildTrackingPayload() loadstrike.TrackingPayload {
	builder := loadstrike.TrackingPayloadBuilder{
		Headers:            map[string]string{"X-Correlation-Id": "ord-1"},
		ContentType:        "application/json",
		MessagePayloadType: "json",
	}
	builder.SetBody(map[string]any{"trackingId": "ord-1"})
	return builder.Build()
}

func createAsyncRunResult() loadstrike.LoadStrikeRunResult {
	var result loadstrike.LoadStrikeRunResult
	return result
}

func createAsyncScenarioStats() loadstrike.LoadStrikeScenarioStats {
	var stats loadstrike.LoadStrikeScenarioStats
	return stats
}

func createAsyncWriteTempConfigFiles() createAsyncTempConfigPaths {
	return createAsyncTempConfigPaths{
		ConfigPath: "method-reference.loadstrike.config.json",
		InfraPath:  "method-reference.loadstrike.infra.json",
	}
}

        // Create the GET scenario from an asynchronous callback.
func (reference CreateAsyncMethodReference) CreateAsyncScenarioExample() any {
    return loadstrike.CreateScenarioAsync("orders.get-by-id", createAsyncExecuteOrderGetAsync)
}

// Create the async scenario and then attach a simple simulation.
func (reference CreateAsyncMethodReference) CreateAsyncScenarioWithSimulationExample() any {
    return loadstrike.CreateScenarioAsync("orders.get-by-id", createAsyncExecuteOrderGetAsync).WithLoadSimulations(loadstrike.LoadStrikeSimulation.IterationsForConstant(1, 3))
}
