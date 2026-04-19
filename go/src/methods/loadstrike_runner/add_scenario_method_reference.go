package loadstrike_runner

        import (
        	"time"

        	loadstrike "loadstrike.com/sdk/go"
        )

        const addScenarioRunnerKey = "runner_dummy_orders_reference"

type AddScenarioMethodReference struct{}

type addScenarioTempConfigPaths struct {
	ConfigPath string
	InfraPath  string
}

type addScenarioOrdersReportingSink struct{}

func newAddScenarioOrdersReportingSink() addScenarioOrdersReportingSink {
	return addScenarioOrdersReportingSink{}
}
func (addScenarioOrdersReportingSink) SinkName() string {
	return "orders-sample-sink"
}
func (addScenarioOrdersReportingSink) Init(loadstrike.LoadStrikeBaseContext, loadstrike.IConfiguration) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (addScenarioOrdersReportingSink) Start(loadstrike.LoadStrikeSessionStartInfo) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (addScenarioOrdersReportingSink) SaveRealtimeStats([]loadstrike.LoadStrikeScenarioStats) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (addScenarioOrdersReportingSink) SaveRealtimeMetrics(loadstrike.LoadStrikeMetricStats) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (addScenarioOrdersReportingSink) SaveRunResult(loadstrike.LoadStrikeRunResult) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (addScenarioOrdersReportingSink) Stop() loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (addScenarioOrdersReportingSink) Dispose() {}

type addScenarioOrdersRuntimePolicy struct{}

func newAddScenarioOrdersRuntimePolicy() addScenarioOrdersRuntimePolicy {
	return addScenarioOrdersRuntimePolicy{}
}
func (addScenarioOrdersRuntimePolicy) ShouldRunScenario(string) loadstrike.LoadStrikeBoolTask {
	return loadstrike.TaskFromBool(true)
}
func (addScenarioOrdersRuntimePolicy) BeforeScenario(string) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (addScenarioOrdersRuntimePolicy) AfterScenario(string, loadstrike.LoadStrikeScenarioRuntime) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (addScenarioOrdersRuntimePolicy) BeforeStep(string, string) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (addScenarioOrdersRuntimePolicy) AfterStep(string, string, loadstrike.LoadStrikeReply) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}

type addScenarioOrdersWorkerPlugin struct{}

func newAddScenarioOrdersWorkerPlugin() addScenarioOrdersWorkerPlugin {
	return addScenarioOrdersWorkerPlugin{}
}
func (addScenarioOrdersWorkerPlugin) PluginName() string {
	return "orders-sample-plugin"
}
func (addScenarioOrdersWorkerPlugin) Init(loadstrike.LoadStrikeBaseContext, loadstrike.IConfiguration) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (addScenarioOrdersWorkerPlugin) Start(loadstrike.LoadStrikeSessionStartInfo) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (addScenarioOrdersWorkerPlugin) GetData(loadstrike.LoadStrikeRunResult) loadstrike.LoadStrikeValueTask[loadstrike.LoadStrikePluginData] {
	return loadstrike.TaskFromResult(loadstrike.CreatePluginData("orders-sample-plugin"))
}
func (addScenarioOrdersWorkerPlugin) Stop() loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (addScenarioOrdersWorkerPlugin) Dispose() loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}

func addScenarioPerformOrderGetReply() loadstrike.LoadStrikeReply {
	return loadstrike.LoadStrikeResponse.Ok("200", int64(128), "ok", loadstrike.TimeSpan(3*time.Millisecond))
}

func addScenarioExecuteOrderGet(context loadstrike.LoadStrikeScenarioContext) loadstrike.LoadStrikeReply {
	return loadstrike.LoadStrikeStep.Run("get-order", context, func(loadstrike.LoadStrikeScenarioContext) loadstrike.LoadStrikeReply {
		return addScenarioPerformOrderGetReply()
	})
}

func addScenarioExecuteOrderGetAsync(context loadstrike.LoadStrikeScenarioContext) loadstrike.LoadStrikeValueTask[loadstrike.LoadStrikeReply] {
	return loadstrike.TaskFromResult(addScenarioExecuteOrderGet(context))
}

func addScenarioBaselineScenario(names ...string) loadstrike.LoadStrikeScenario {
	name := "orders.get-by-id"
	if len(names) > 0 && names[0] != "" {
		name = names[0]
	}
	return loadstrike.CreateScenario(name, addScenarioExecuteOrderGet).
		WithLoadSimulations(loadstrike.LoadStrikeSimulation.IterationsForConstant(1, 1)).
		WithoutWarmUp()
}

func addScenarioBaseRunner() loadstrike.LoadStrikeRunner {
	return loadstrike.Create().
		AddScenario(addScenarioBaselineScenario()).
		WithRunnerKey(addScenarioRunnerKey).
		WithTestSuite("orders-reference").
		WithTestName("orders-get-by-id").
		WithoutReports()
}

func addScenarioBaseContext() loadstrike.LoadStrikeContext {
	return addScenarioBaseRunner().BuildContext()
}

func addScenarioHttpSource() *loadstrike.EndpointSpec {
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

func addScenarioHttpDestination() *loadstrike.EndpointSpec {
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

func addScenarioTrackingConfiguration() *loadstrike.TrackingConfigurationSpec {
	return &loadstrike.TrackingConfigurationSpec{
		Source:                      addScenarioHttpSource(),
		Destination:                 addScenarioHttpDestination(),
		RunMode:                     "GenerateAndCorrelate",
		CorrelationTimeoutSeconds:   30,
		TimeoutSweepIntervalSeconds: 1,
		TimeoutBatchSize:            200,
		TimeoutCountsAsFailure:      true,
		MetricPrefix:                "orders_tracking",
		ExecuteOriginalScenarioRun:  false,
	}
}

func addScenarioTrackedContext() loadstrike.LoadStrikeContext {
	return loadstrike.Create().
		AddScenario(addScenarioBaselineScenario("orders.tracked").WithCrossPlatformTracking(addScenarioTrackingConfiguration())).
		WithRunnerKey(addScenarioRunnerKey).
		WithoutReports().
		BuildContext()
}

func addScenarioBuildTrackingPayload() loadstrike.TrackingPayload {
	builder := loadstrike.TrackingPayloadBuilder{
		Headers:            map[string]string{"X-Correlation-Id": "ord-1"},
		ContentType:        "application/json",
		MessagePayloadType: "json",
	}
	builder.SetBody(map[string]any{"trackingId": "ord-1"})
	return builder.Build()
}

func addScenarioRunResult() loadstrike.LoadStrikeRunResult {
	var result loadstrike.LoadStrikeRunResult
	return result
}

func addScenarioScenarioStats() loadstrike.LoadStrikeScenarioStats {
	var stats loadstrike.LoadStrikeScenarioStats
	return stats
}

func addScenarioWriteTempConfigFiles() addScenarioTempConfigPaths {
	return addScenarioTempConfigPaths{
		ConfigPath: "method-reference.loadstrike.config.json",
		InfraPath:  "method-reference.loadstrike.infra.json",
	}
}

        // Attach one scenario to a reusable runner.
func (reference AddScenarioMethodReference) AddOneScenarioExample() any {
    return addScenarioBaseRunner().AddScenario(addScenarioBaselineScenario())
}

// Attach a second scenario after the first GET scenario is already registered.
func (reference AddScenarioMethodReference) AddSecondScenarioExample() any {
    return addScenarioBaseRunner().AddScenario(addScenarioBaselineScenario()).AddScenario(addScenarioBaselineScenario("orders.audit"))
}
