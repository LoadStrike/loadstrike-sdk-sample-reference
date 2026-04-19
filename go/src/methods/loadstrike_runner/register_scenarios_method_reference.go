package loadstrike_runner

        import (
        	"time"

        	loadstrike "loadstrike.com/sdk/go"
        )

        const registerScenariosRunnerKey = "runner_dummy_orders_reference"

type RegisterScenariosMethodReference struct{}

type registerScenariosTempConfigPaths struct {
	ConfigPath string
	InfraPath  string
}

type registerScenariosOrdersReportingSink struct{}

func newRegisterScenariosOrdersReportingSink() registerScenariosOrdersReportingSink {
	return registerScenariosOrdersReportingSink{}
}
func (registerScenariosOrdersReportingSink) SinkName() string {
	return "orders-sample-sink"
}
func (registerScenariosOrdersReportingSink) Init(loadstrike.LoadStrikeBaseContext, loadstrike.IConfiguration) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (registerScenariosOrdersReportingSink) Start(loadstrike.LoadStrikeSessionStartInfo) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (registerScenariosOrdersReportingSink) SaveRealtimeStats([]loadstrike.LoadStrikeScenarioStats) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (registerScenariosOrdersReportingSink) SaveRealtimeMetrics(loadstrike.LoadStrikeMetricStats) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (registerScenariosOrdersReportingSink) SaveRunResult(loadstrike.LoadStrikeRunResult) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (registerScenariosOrdersReportingSink) Stop() loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (registerScenariosOrdersReportingSink) Dispose() {}

type registerScenariosOrdersRuntimePolicy struct{}

func newRegisterScenariosOrdersRuntimePolicy() registerScenariosOrdersRuntimePolicy {
	return registerScenariosOrdersRuntimePolicy{}
}
func (registerScenariosOrdersRuntimePolicy) ShouldRunScenario(string) loadstrike.LoadStrikeBoolTask {
	return loadstrike.TaskFromBool(true)
}
func (registerScenariosOrdersRuntimePolicy) BeforeScenario(string) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (registerScenariosOrdersRuntimePolicy) AfterScenario(string, loadstrike.LoadStrikeScenarioRuntime) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (registerScenariosOrdersRuntimePolicy) BeforeStep(string, string) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (registerScenariosOrdersRuntimePolicy) AfterStep(string, string, loadstrike.LoadStrikeReply) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}

type registerScenariosOrdersWorkerPlugin struct{}

func newRegisterScenariosOrdersWorkerPlugin() registerScenariosOrdersWorkerPlugin {
	return registerScenariosOrdersWorkerPlugin{}
}
func (registerScenariosOrdersWorkerPlugin) PluginName() string {
	return "orders-sample-plugin"
}
func (registerScenariosOrdersWorkerPlugin) Init(loadstrike.LoadStrikeBaseContext, loadstrike.IConfiguration) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (registerScenariosOrdersWorkerPlugin) Start(loadstrike.LoadStrikeSessionStartInfo) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (registerScenariosOrdersWorkerPlugin) GetData(loadstrike.LoadStrikeRunResult) loadstrike.LoadStrikeValueTask[loadstrike.LoadStrikePluginData] {
	return loadstrike.TaskFromResult(loadstrike.CreatePluginData("orders-sample-plugin"))
}
func (registerScenariosOrdersWorkerPlugin) Stop() loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (registerScenariosOrdersWorkerPlugin) Dispose() loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}

func registerScenariosPerformOrderGetReply() loadstrike.LoadStrikeReply {
	return loadstrike.LoadStrikeResponse.Ok("200", int64(128), "ok", loadstrike.TimeSpan(3*time.Millisecond))
}

func registerScenariosExecuteOrderGet(context loadstrike.LoadStrikeScenarioContext) loadstrike.LoadStrikeReply {
	return loadstrike.LoadStrikeStep.Run("get-order", context, func(loadstrike.LoadStrikeScenarioContext) loadstrike.LoadStrikeReply {
		return registerScenariosPerformOrderGetReply()
	})
}

func registerScenariosExecuteOrderGetAsync(context loadstrike.LoadStrikeScenarioContext) loadstrike.LoadStrikeValueTask[loadstrike.LoadStrikeReply] {
	return loadstrike.TaskFromResult(registerScenariosExecuteOrderGet(context))
}

func registerScenariosBaselineScenario(names ...string) loadstrike.LoadStrikeScenario {
	name := "orders.get-by-id"
	if len(names) > 0 && names[0] != "" {
		name = names[0]
	}
	return loadstrike.CreateScenario(name, registerScenariosExecuteOrderGet).
		WithLoadSimulations(loadstrike.LoadStrikeSimulation.IterationsForConstant(1, 1)).
		WithoutWarmUp()
}

func registerScenariosBaseRunner() loadstrike.LoadStrikeRunner {
	return loadstrike.Create().
		AddScenario(registerScenariosBaselineScenario()).
		WithRunnerKey(registerScenariosRunnerKey).
		WithTestSuite("orders-reference").
		WithTestName("orders-get-by-id").
		WithoutReports()
}

func registerScenariosBaseContext() loadstrike.LoadStrikeContext {
	return registerScenariosBaseRunner().BuildContext()
}

func registerScenariosHttpSource() *loadstrike.EndpointSpec {
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

func registerScenariosHttpDestination() *loadstrike.EndpointSpec {
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

func registerScenariosTrackingConfiguration() *loadstrike.TrackingConfigurationSpec {
	return &loadstrike.TrackingConfigurationSpec{
		Source:                      registerScenariosHttpSource(),
		Destination:                 registerScenariosHttpDestination(),
		RunMode:                     "GenerateAndCorrelate",
		CorrelationTimeoutSeconds:   30,
		TimeoutSweepIntervalSeconds: 1,
		TimeoutBatchSize:            200,
		TimeoutCountsAsFailure:      true,
		MetricPrefix:                "orders_tracking",
		ExecuteOriginalScenarioRun:  false,
	}
}

func registerScenariosTrackedContext() loadstrike.LoadStrikeContext {
	return loadstrike.Create().
		AddScenario(registerScenariosBaselineScenario("orders.tracked").WithCrossPlatformTracking(registerScenariosTrackingConfiguration())).
		WithRunnerKey(registerScenariosRunnerKey).
		WithoutReports().
		BuildContext()
}

func registerScenariosBuildTrackingPayload() loadstrike.TrackingPayload {
	builder := loadstrike.TrackingPayloadBuilder{
		Headers:            map[string]string{"X-Correlation-Id": "ord-1"},
		ContentType:        "application/json",
		MessagePayloadType: "json",
	}
	builder.SetBody(map[string]any{"trackingId": "ord-1"})
	return builder.Build()
}

func registerScenariosRunResult() loadstrike.LoadStrikeRunResult {
	var result loadstrike.LoadStrikeRunResult
	return result
}

func registerScenariosScenarioStats() loadstrike.LoadStrikeScenarioStats {
	var stats loadstrike.LoadStrikeScenarioStats
	return stats
}

func registerScenariosWriteTempConfigFiles() registerScenariosTempConfigPaths {
	return registerScenariosTempConfigPaths{
		ConfigPath: "method-reference.loadstrike.config.json",
		InfraPath:  "method-reference.loadstrike.infra.json",
	}
}

        // Register one baseline GET scenario in a single call.
func (reference RegisterScenariosMethodReference) RegisterSingleScenarioExample() any {
    return loadstrike.RegisterScenarios(registerScenariosBaselineScenario())
}

// Register multiple scenarios without switching back to the builder flow.
func (reference RegisterScenariosMethodReference) RegisterMultipleScenariosExample() any {
    return loadstrike.RegisterScenarios(registerScenariosBaselineScenario(), registerScenariosBaselineScenario("orders.audit"))
}
