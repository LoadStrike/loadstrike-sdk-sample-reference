package loadstrike_threshold

        import (
        	"time"

        	loadstrike "loadstrike.com/sdk/go"
        )

        const createScenarioRunnerKey = "runner_dummy_orders_reference"

type CreateScenarioMethodReference struct{}

type createScenarioTempConfigPaths struct {
	ConfigPath string
	InfraPath  string
}

type createScenarioOrdersReportingSink struct{}

func newCreateScenarioOrdersReportingSink() createScenarioOrdersReportingSink {
	return createScenarioOrdersReportingSink{}
}
func (createScenarioOrdersReportingSink) SinkName() string {
	return "orders-sample-sink"
}
func (createScenarioOrdersReportingSink) Init(loadstrike.LoadStrikeBaseContext, loadstrike.IConfiguration) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (createScenarioOrdersReportingSink) Start(loadstrike.LoadStrikeSessionStartInfo) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (createScenarioOrdersReportingSink) SaveRealtimeStats([]loadstrike.LoadStrikeScenarioStats) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (createScenarioOrdersReportingSink) SaveRealtimeMetrics(loadstrike.LoadStrikeMetricStats) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (createScenarioOrdersReportingSink) SaveRunResult(loadstrike.LoadStrikeRunResult) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (createScenarioOrdersReportingSink) Stop() loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (createScenarioOrdersReportingSink) Dispose() {}

type createScenarioOrdersRuntimePolicy struct{}

func newCreateScenarioOrdersRuntimePolicy() createScenarioOrdersRuntimePolicy {
	return createScenarioOrdersRuntimePolicy{}
}
func (createScenarioOrdersRuntimePolicy) ShouldRunScenario(string) loadstrike.LoadStrikeBoolTask {
	return loadstrike.TaskFromBool(true)
}
func (createScenarioOrdersRuntimePolicy) BeforeScenario(string) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (createScenarioOrdersRuntimePolicy) AfterScenario(string, loadstrike.LoadStrikeScenarioRuntime) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (createScenarioOrdersRuntimePolicy) BeforeStep(string, string) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (createScenarioOrdersRuntimePolicy) AfterStep(string, string, loadstrike.LoadStrikeReply) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}

type createScenarioOrdersWorkerPlugin struct{}

func newCreateScenarioOrdersWorkerPlugin() createScenarioOrdersWorkerPlugin {
	return createScenarioOrdersWorkerPlugin{}
}
func (createScenarioOrdersWorkerPlugin) PluginName() string {
	return "orders-sample-plugin"
}
func (createScenarioOrdersWorkerPlugin) Init(loadstrike.LoadStrikeBaseContext, loadstrike.IConfiguration) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (createScenarioOrdersWorkerPlugin) Start(loadstrike.LoadStrikeSessionStartInfo) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (createScenarioOrdersWorkerPlugin) GetData(loadstrike.LoadStrikeRunResult) loadstrike.LoadStrikeValueTask[loadstrike.LoadStrikePluginData] {
	return loadstrike.TaskFromResult(loadstrike.CreatePluginData("orders-sample-plugin"))
}
func (createScenarioOrdersWorkerPlugin) Stop() loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (createScenarioOrdersWorkerPlugin) Dispose() loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}

func createScenarioPerformOrderGetReply() loadstrike.LoadStrikeReply {
	return loadstrike.LoadStrikeResponse.Ok("200", int64(128), "ok", loadstrike.TimeSpan(3*time.Millisecond))
}

func createScenarioExecuteOrderGet(context loadstrike.LoadStrikeScenarioContext) loadstrike.LoadStrikeReply {
	return loadstrike.LoadStrikeStep.Run("get-order", context, func(loadstrike.LoadStrikeScenarioContext) loadstrike.LoadStrikeReply {
		return createScenarioPerformOrderGetReply()
	})
}

func createScenarioExecuteOrderGetAsync(context loadstrike.LoadStrikeScenarioContext) loadstrike.LoadStrikeValueTask[loadstrike.LoadStrikeReply] {
	return loadstrike.TaskFromResult(createScenarioExecuteOrderGet(context))
}

func createScenarioBaselineScenario(names ...string) loadstrike.LoadStrikeScenario {
	name := "orders.get-by-id"
	if len(names) > 0 && names[0] != "" {
		name = names[0]
	}
	return loadstrike.CreateScenario(name, createScenarioExecuteOrderGet).
		WithLoadSimulations(loadstrike.LoadStrikeSimulation.IterationsForConstant(1, 1)).
		WithoutWarmUp()
}

func createScenarioBaseRunner() loadstrike.LoadStrikeRunner {
	return loadstrike.Create().
		AddScenario(createScenarioBaselineScenario()).
		WithRunnerKey(createScenarioRunnerKey).
		WithTestSuite("orders-reference").
		WithTestName("orders-get-by-id").
		WithoutReports()
}

func createScenarioBaseContext() loadstrike.LoadStrikeContext {
	return createScenarioBaseRunner().BuildContext()
}

func createScenarioHttpSource() *loadstrike.EndpointSpec {
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

func createScenarioHttpDestination() *loadstrike.EndpointSpec {
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

func createScenarioTrackingConfiguration() *loadstrike.TrackingConfigurationSpec {
	return &loadstrike.TrackingConfigurationSpec{
		Source:                      createScenarioHttpSource(),
		Destination:                 createScenarioHttpDestination(),
		RunMode:                     "GenerateAndCorrelate",
		CorrelationTimeoutSeconds:   30,
		TimeoutSweepIntervalSeconds: 1,
		TimeoutBatchSize:            200,
		TimeoutCountsAsFailure:      true,
		MetricPrefix:                "orders_tracking",
		ExecuteOriginalScenarioRun:  false,
	}
}

func createScenarioTrackedContext() loadstrike.LoadStrikeContext {
	return loadstrike.Create().
		AddScenario(createScenarioBaselineScenario("orders.tracked").WithCrossPlatformTracking(createScenarioTrackingConfiguration())).
		WithRunnerKey(createScenarioRunnerKey).
		WithoutReports().
		BuildContext()
}

func createScenarioBuildTrackingPayload() loadstrike.TrackingPayload {
	builder := loadstrike.TrackingPayloadBuilder{
		Headers:            map[string]string{"X-Correlation-Id": "ord-1"},
		ContentType:        "application/json",
		MessagePayloadType: "json",
	}
	builder.SetBody(map[string]any{"trackingId": "ord-1"})
	return builder.Build()
}

func createScenarioRunResult() loadstrike.LoadStrikeRunResult {
	var result loadstrike.LoadStrikeRunResult
	return result
}

func createScenarioScenarioStats() loadstrike.LoadStrikeScenarioStats {
	var stats loadstrike.LoadStrikeScenarioStats
	return stats
}

func createScenarioWriteTempConfigFiles() createScenarioTempConfigPaths {
	return createScenarioTempConfigPaths{
		ConfigPath: "method-reference.loadstrike.config.json",
		InfraPath:  "method-reference.loadstrike.infra.json",
	}
}

        // Create the threshold from the public threshold helper.
func (reference CreateScenarioMethodReference) CreateThresholdExample() any {
    return loadstrike.LoadStrikeThreshold{}.CreateScenario("allokcount", "gte", 1)
}

// Attach the threshold to the baseline GET scenario.
func (reference CreateScenarioMethodReference) AttachThresholdToScenarioExample() any {
    return createScenarioBaselineScenario().WithThresholds(loadstrike.LoadStrikeThreshold{}.CreateScenario("allokcount", "gte", 1))
}
