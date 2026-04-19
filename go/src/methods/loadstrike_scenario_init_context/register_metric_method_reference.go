package loadstrike_scenario_init_context

        import (
        	"time"

        	loadstrike "loadstrike.com/sdk/go"
        )

        const registerMetricRunnerKey = "runner_dummy_orders_reference"

type RegisterMetricMethodReference struct{}

type registerMetricTempConfigPaths struct {
	ConfigPath string
	InfraPath  string
}

type registerMetricOrdersReportingSink struct{}

func newRegisterMetricOrdersReportingSink() registerMetricOrdersReportingSink {
	return registerMetricOrdersReportingSink{}
}
func (registerMetricOrdersReportingSink) SinkName() string {
	return "orders-sample-sink"
}
func (registerMetricOrdersReportingSink) Init(loadstrike.LoadStrikeBaseContext, loadstrike.IConfiguration) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (registerMetricOrdersReportingSink) Start(loadstrike.LoadStrikeSessionStartInfo) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (registerMetricOrdersReportingSink) SaveRealtimeStats([]loadstrike.LoadStrikeScenarioStats) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (registerMetricOrdersReportingSink) SaveRealtimeMetrics(loadstrike.LoadStrikeMetricStats) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (registerMetricOrdersReportingSink) SaveRunResult(loadstrike.LoadStrikeRunResult) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (registerMetricOrdersReportingSink) Stop() loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (registerMetricOrdersReportingSink) Dispose() {}

type registerMetricOrdersRuntimePolicy struct{}

func newRegisterMetricOrdersRuntimePolicy() registerMetricOrdersRuntimePolicy {
	return registerMetricOrdersRuntimePolicy{}
}
func (registerMetricOrdersRuntimePolicy) ShouldRunScenario(string) loadstrike.LoadStrikeBoolTask {
	return loadstrike.TaskFromBool(true)
}
func (registerMetricOrdersRuntimePolicy) BeforeScenario(string) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (registerMetricOrdersRuntimePolicy) AfterScenario(string, loadstrike.LoadStrikeScenarioRuntime) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (registerMetricOrdersRuntimePolicy) BeforeStep(string, string) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (registerMetricOrdersRuntimePolicy) AfterStep(string, string, loadstrike.LoadStrikeReply) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}

type registerMetricOrdersWorkerPlugin struct{}

func newRegisterMetricOrdersWorkerPlugin() registerMetricOrdersWorkerPlugin {
	return registerMetricOrdersWorkerPlugin{}
}
func (registerMetricOrdersWorkerPlugin) PluginName() string {
	return "orders-sample-plugin"
}
func (registerMetricOrdersWorkerPlugin) Init(loadstrike.LoadStrikeBaseContext, loadstrike.IConfiguration) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (registerMetricOrdersWorkerPlugin) Start(loadstrike.LoadStrikeSessionStartInfo) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (registerMetricOrdersWorkerPlugin) GetData(loadstrike.LoadStrikeRunResult) loadstrike.LoadStrikeValueTask[loadstrike.LoadStrikePluginData] {
	return loadstrike.TaskFromResult(loadstrike.CreatePluginData("orders-sample-plugin"))
}
func (registerMetricOrdersWorkerPlugin) Stop() loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (registerMetricOrdersWorkerPlugin) Dispose() loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}

func registerMetricPerformOrderGetReply() loadstrike.LoadStrikeReply {
	return loadstrike.LoadStrikeResponse.Ok("200", int64(128), "ok", loadstrike.TimeSpan(3*time.Millisecond))
}

func registerMetricExecuteOrderGet(context loadstrike.LoadStrikeScenarioContext) loadstrike.LoadStrikeReply {
	return loadstrike.LoadStrikeStep.Run("get-order", context, func(loadstrike.LoadStrikeScenarioContext) loadstrike.LoadStrikeReply {
		return registerMetricPerformOrderGetReply()
	})
}

func registerMetricExecuteOrderGetAsync(context loadstrike.LoadStrikeScenarioContext) loadstrike.LoadStrikeValueTask[loadstrike.LoadStrikeReply] {
	return loadstrike.TaskFromResult(registerMetricExecuteOrderGet(context))
}

func registerMetricBaselineScenario(names ...string) loadstrike.LoadStrikeScenario {
	name := "orders.get-by-id"
	if len(names) > 0 && names[0] != "" {
		name = names[0]
	}
	return loadstrike.CreateScenario(name, registerMetricExecuteOrderGet).
		WithLoadSimulations(loadstrike.LoadStrikeSimulation.IterationsForConstant(1, 1)).
		WithoutWarmUp()
}

func registerMetricBaseRunner() loadstrike.LoadStrikeRunner {
	return loadstrike.Create().
		AddScenario(registerMetricBaselineScenario()).
		WithRunnerKey(registerMetricRunnerKey).
		WithTestSuite("orders-reference").
		WithTestName("orders-get-by-id").
		WithoutReports()
}

func registerMetricBaseContext() loadstrike.LoadStrikeContext {
	return registerMetricBaseRunner().BuildContext()
}

func registerMetricHttpSource() *loadstrike.EndpointSpec {
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

func registerMetricHttpDestination() *loadstrike.EndpointSpec {
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

func registerMetricTrackingConfiguration() *loadstrike.TrackingConfigurationSpec {
	return &loadstrike.TrackingConfigurationSpec{
		Source:                      registerMetricHttpSource(),
		Destination:                 registerMetricHttpDestination(),
		RunMode:                     "GenerateAndCorrelate",
		CorrelationTimeoutSeconds:   30,
		TimeoutSweepIntervalSeconds: 1,
		TimeoutBatchSize:            200,
		TimeoutCountsAsFailure:      true,
		MetricPrefix:                "orders_tracking",
		ExecuteOriginalScenarioRun:  false,
	}
}

func registerMetricTrackedContext() loadstrike.LoadStrikeContext {
	return loadstrike.Create().
		AddScenario(registerMetricBaselineScenario("orders.tracked").WithCrossPlatformTracking(registerMetricTrackingConfiguration())).
		WithRunnerKey(registerMetricRunnerKey).
		WithoutReports().
		BuildContext()
}

func registerMetricBuildTrackingPayload() loadstrike.TrackingPayload {
	builder := loadstrike.TrackingPayloadBuilder{
		Headers:            map[string]string{"X-Correlation-Id": "ord-1"},
		ContentType:        "application/json",
		MessagePayloadType: "json",
	}
	builder.SetBody(map[string]any{"trackingId": "ord-1"})
	return builder.Build()
}

func registerMetricRunResult() loadstrike.LoadStrikeRunResult {
	var result loadstrike.LoadStrikeRunResult
	return result
}

func registerMetricScenarioStats() loadstrike.LoadStrikeScenarioStats {
	var stats loadstrike.LoadStrikeScenarioStats
	return stats
}

func registerMetricWriteTempConfigFiles() registerMetricTempConfigPaths {
	return registerMetricTempConfigPaths{
		ConfigPath: "method-reference.loadstrike.config.json",
		InfraPath:  "method-reference.loadstrike.infra.json",
	}
}

        // Call the public context helper directly from the scenario context surface.
func (reference RegisterMetricMethodReference) UseContextMethodExample() any {
    return loadstrike.Metric.CreateCounter("orders_total", "count")
}

// Show the same helper in the baseline GET-step flow.
        func (reference RegisterMetricMethodReference) UseContextMethodInStepExample() any {
            metric := loadstrike.Metric.CreateCounter("orders_total", "count")
scenario := loadstrike.CreateScenario("orders.metric", registerMetricExecuteOrderGet).WithInit(func(context loadstrike.LoadStrikeScenarioInitContext) error { context.RegisterMetric(metric); return nil })
return scenario
        }
