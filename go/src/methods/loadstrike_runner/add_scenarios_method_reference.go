package loadstrike_runner

        import (
        	"time"

        	loadstrike "loadstrike.com/sdk/go"
        )

        const addScenariosRunnerKey = "runner_dummy_orders_reference"

type AddScenariosMethodReference struct{}

type addScenariosTempConfigPaths struct {
	ConfigPath string
	InfraPath  string
}

type addScenariosOrdersReportingSink struct{}

func newAddScenariosOrdersReportingSink() addScenariosOrdersReportingSink {
	return addScenariosOrdersReportingSink{}
}
func (addScenariosOrdersReportingSink) SinkName() string {
	return "orders-sample-sink"
}
func (addScenariosOrdersReportingSink) Init(loadstrike.LoadStrikeBaseContext, loadstrike.IConfiguration) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (addScenariosOrdersReportingSink) Start(loadstrike.LoadStrikeSessionStartInfo) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (addScenariosOrdersReportingSink) SaveRealtimeStats([]loadstrike.LoadStrikeScenarioStats) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (addScenariosOrdersReportingSink) SaveRealtimeMetrics(loadstrike.LoadStrikeMetricStats) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (addScenariosOrdersReportingSink) SaveRunResult(loadstrike.LoadStrikeRunResult) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (addScenariosOrdersReportingSink) Stop() loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (addScenariosOrdersReportingSink) Dispose() {}

type addScenariosOrdersRuntimePolicy struct{}

func newAddScenariosOrdersRuntimePolicy() addScenariosOrdersRuntimePolicy {
	return addScenariosOrdersRuntimePolicy{}
}
func (addScenariosOrdersRuntimePolicy) ShouldRunScenario(string) loadstrike.LoadStrikeBoolTask {
	return loadstrike.TaskFromBool(true)
}
func (addScenariosOrdersRuntimePolicy) BeforeScenario(string) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (addScenariosOrdersRuntimePolicy) AfterScenario(string, loadstrike.LoadStrikeScenarioRuntime) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (addScenariosOrdersRuntimePolicy) BeforeStep(string, string) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (addScenariosOrdersRuntimePolicy) AfterStep(string, string, loadstrike.LoadStrikeReply) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}

type addScenariosOrdersWorkerPlugin struct{}

func newAddScenariosOrdersWorkerPlugin() addScenariosOrdersWorkerPlugin {
	return addScenariosOrdersWorkerPlugin{}
}
func (addScenariosOrdersWorkerPlugin) PluginName() string {
	return "orders-sample-plugin"
}
func (addScenariosOrdersWorkerPlugin) Init(loadstrike.LoadStrikeBaseContext, loadstrike.IConfiguration) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (addScenariosOrdersWorkerPlugin) Start(loadstrike.LoadStrikeSessionStartInfo) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (addScenariosOrdersWorkerPlugin) GetData(loadstrike.LoadStrikeRunResult) loadstrike.LoadStrikeValueTask[loadstrike.LoadStrikePluginData] {
	return loadstrike.TaskFromResult(loadstrike.CreatePluginData("orders-sample-plugin"))
}
func (addScenariosOrdersWorkerPlugin) Stop() loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (addScenariosOrdersWorkerPlugin) Dispose() loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}

func addScenariosPerformOrderGetReply() loadstrike.LoadStrikeReply {
	return loadstrike.LoadStrikeResponse.Ok("200", int64(128), "ok", loadstrike.TimeSpan(3*time.Millisecond))
}

func addScenariosExecuteOrderGet(context loadstrike.LoadStrikeScenarioContext) loadstrike.LoadStrikeReply {
	return loadstrike.LoadStrikeStep.Run("get-order", context, func(loadstrike.LoadStrikeScenarioContext) loadstrike.LoadStrikeReply {
		return addScenariosPerformOrderGetReply()
	})
}

func addScenariosExecuteOrderGetAsync(context loadstrike.LoadStrikeScenarioContext) loadstrike.LoadStrikeValueTask[loadstrike.LoadStrikeReply] {
	return loadstrike.TaskFromResult(addScenariosExecuteOrderGet(context))
}

func addScenariosBaselineScenario(names ...string) loadstrike.LoadStrikeScenario {
	name := "orders.get-by-id"
	if len(names) > 0 && names[0] != "" {
		name = names[0]
	}
	return loadstrike.CreateScenario(name, addScenariosExecuteOrderGet).
		WithLoadSimulations(loadstrike.LoadStrikeSimulation.IterationsForConstant(1, 1)).
		WithoutWarmUp()
}

func addScenariosBaseRunner() loadstrike.LoadStrikeRunner {
	return loadstrike.Create().
		AddScenario(addScenariosBaselineScenario()).
		WithRunnerKey(addScenariosRunnerKey).
		WithTestSuite("orders-reference").
		WithTestName("orders-get-by-id").
		WithoutReports()
}

func addScenariosBaseContext() loadstrike.LoadStrikeContext {
	return addScenariosBaseRunner().BuildContext()
}

func addScenariosHttpSource() *loadstrike.EndpointSpec {
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

func addScenariosHttpDestination() *loadstrike.EndpointSpec {
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

func addScenariosTrackingConfiguration() *loadstrike.TrackingConfigurationSpec {
	return &loadstrike.TrackingConfigurationSpec{
		Source:                      addScenariosHttpSource(),
		Destination:                 addScenariosHttpDestination(),
		RunMode:                     "GenerateAndCorrelate",
		CorrelationTimeoutSeconds:   30,
		TimeoutSweepIntervalSeconds: 1,
		TimeoutBatchSize:            200,
		TimeoutCountsAsFailure:      true,
		MetricPrefix:                "orders_tracking",
		ExecuteOriginalScenarioRun:  false,
	}
}

func addScenariosTrackedContext() loadstrike.LoadStrikeContext {
	return loadstrike.Create().
		AddScenario(addScenariosBaselineScenario("orders.tracked").WithCrossPlatformTracking(addScenariosTrackingConfiguration())).
		WithRunnerKey(addScenariosRunnerKey).
		WithoutReports().
		BuildContext()
}

func addScenariosBuildTrackingPayload() loadstrike.TrackingPayload {
	builder := loadstrike.TrackingPayloadBuilder{
		Headers:            map[string]string{"X-Correlation-Id": "ord-1"},
		ContentType:        "application/json",
		MessagePayloadType: "json",
	}
	builder.SetBody(map[string]any{"trackingId": "ord-1"})
	return builder.Build()
}

func addScenariosRunResult() loadstrike.LoadStrikeRunResult {
	var result loadstrike.LoadStrikeRunResult
	return result
}

func addScenariosScenarioStats() loadstrike.LoadStrikeScenarioStats {
	var stats loadstrike.LoadStrikeScenarioStats
	return stats
}

func addScenariosWriteTempConfigFiles() addScenariosTempConfigPaths {
	return addScenariosTempConfigPaths{
		ConfigPath: "method-reference.loadstrike.config.json",
		InfraPath:  "method-reference.loadstrike.infra.json",
	}
}

        // Attach two scenarios in a single batch call.
func (reference AddScenariosMethodReference) AddTwoScenariosExample() any {
    return addScenariosBaseRunner().AddScenarios(addScenariosBaselineScenario(), addScenariosBaselineScenario("orders.audit"))
}

// Keep chaining runner options after adding a batch of scenarios.
func (reference AddScenariosMethodReference) AddBatchAndKeepBuildingExample() any {
    return addScenariosBaseRunner().AddScenarios(addScenariosBaselineScenario(), addScenariosBaselineScenario("orders.audit")).WithTestName("batched-add")
}
