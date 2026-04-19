package loadstrike_runner

        import (
        	"time"

        	loadstrike "loadstrike.com/sdk/go"
        )

        const withLoggerConfigRunnerKey = "runner_dummy_orders_reference"

type WithLoggerConfigMethodReference struct{}

type withLoggerConfigTempConfigPaths struct {
	ConfigPath string
	InfraPath  string
}

type withLoggerConfigOrdersReportingSink struct{}

func newWithLoggerConfigOrdersReportingSink() withLoggerConfigOrdersReportingSink {
	return withLoggerConfigOrdersReportingSink{}
}
func (withLoggerConfigOrdersReportingSink) SinkName() string {
	return "orders-sample-sink"
}
func (withLoggerConfigOrdersReportingSink) Init(loadstrike.LoadStrikeBaseContext, loadstrike.IConfiguration) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withLoggerConfigOrdersReportingSink) Start(loadstrike.LoadStrikeSessionStartInfo) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withLoggerConfigOrdersReportingSink) SaveRealtimeStats([]loadstrike.LoadStrikeScenarioStats) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withLoggerConfigOrdersReportingSink) SaveRealtimeMetrics(loadstrike.LoadStrikeMetricStats) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withLoggerConfigOrdersReportingSink) SaveRunResult(loadstrike.LoadStrikeRunResult) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withLoggerConfigOrdersReportingSink) Stop() loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withLoggerConfigOrdersReportingSink) Dispose() {}

type withLoggerConfigOrdersRuntimePolicy struct{}

func newWithLoggerConfigOrdersRuntimePolicy() withLoggerConfigOrdersRuntimePolicy {
	return withLoggerConfigOrdersRuntimePolicy{}
}
func (withLoggerConfigOrdersRuntimePolicy) ShouldRunScenario(string) loadstrike.LoadStrikeBoolTask {
	return loadstrike.TaskFromBool(true)
}
func (withLoggerConfigOrdersRuntimePolicy) BeforeScenario(string) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withLoggerConfigOrdersRuntimePolicy) AfterScenario(string, loadstrike.LoadStrikeScenarioRuntime) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withLoggerConfigOrdersRuntimePolicy) BeforeStep(string, string) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withLoggerConfigOrdersRuntimePolicy) AfterStep(string, string, loadstrike.LoadStrikeReply) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}

type withLoggerConfigOrdersWorkerPlugin struct{}

func newWithLoggerConfigOrdersWorkerPlugin() withLoggerConfigOrdersWorkerPlugin {
	return withLoggerConfigOrdersWorkerPlugin{}
}
func (withLoggerConfigOrdersWorkerPlugin) PluginName() string {
	return "orders-sample-plugin"
}
func (withLoggerConfigOrdersWorkerPlugin) Init(loadstrike.LoadStrikeBaseContext, loadstrike.IConfiguration) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withLoggerConfigOrdersWorkerPlugin) Start(loadstrike.LoadStrikeSessionStartInfo) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withLoggerConfigOrdersWorkerPlugin) GetData(loadstrike.LoadStrikeRunResult) loadstrike.LoadStrikeValueTask[loadstrike.LoadStrikePluginData] {
	return loadstrike.TaskFromResult(loadstrike.CreatePluginData("orders-sample-plugin"))
}
func (withLoggerConfigOrdersWorkerPlugin) Stop() loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withLoggerConfigOrdersWorkerPlugin) Dispose() loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}

func withLoggerConfigPerformOrderGetReply() loadstrike.LoadStrikeReply {
	return loadstrike.LoadStrikeResponse.Ok("200", int64(128), "ok", loadstrike.TimeSpan(3*time.Millisecond))
}

func withLoggerConfigExecuteOrderGet(context loadstrike.LoadStrikeScenarioContext) loadstrike.LoadStrikeReply {
	return loadstrike.LoadStrikeStep.Run("get-order", context, func(loadstrike.LoadStrikeScenarioContext) loadstrike.LoadStrikeReply {
		return withLoggerConfigPerformOrderGetReply()
	})
}

func withLoggerConfigExecuteOrderGetAsync(context loadstrike.LoadStrikeScenarioContext) loadstrike.LoadStrikeValueTask[loadstrike.LoadStrikeReply] {
	return loadstrike.TaskFromResult(withLoggerConfigExecuteOrderGet(context))
}

func withLoggerConfigBaselineScenario(names ...string) loadstrike.LoadStrikeScenario {
	name := "orders.get-by-id"
	if len(names) > 0 && names[0] != "" {
		name = names[0]
	}
	return loadstrike.CreateScenario(name, withLoggerConfigExecuteOrderGet).
		WithLoadSimulations(loadstrike.LoadStrikeSimulation.IterationsForConstant(1, 1)).
		WithoutWarmUp()
}

func withLoggerConfigBaseRunner() loadstrike.LoadStrikeRunner {
	return loadstrike.Create().
		AddScenario(withLoggerConfigBaselineScenario()).
		WithRunnerKey(withLoggerConfigRunnerKey).
		WithTestSuite("orders-reference").
		WithTestName("orders-get-by-id").
		WithoutReports()
}

func withLoggerConfigBaseContext() loadstrike.LoadStrikeContext {
	return withLoggerConfigBaseRunner().BuildContext()
}

func withLoggerConfigHttpSource() *loadstrike.EndpointSpec {
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

func withLoggerConfigHttpDestination() *loadstrike.EndpointSpec {
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

func withLoggerConfigTrackingConfiguration() *loadstrike.TrackingConfigurationSpec {
	return &loadstrike.TrackingConfigurationSpec{
		Source:                      withLoggerConfigHttpSource(),
		Destination:                 withLoggerConfigHttpDestination(),
		RunMode:                     "GenerateAndCorrelate",
		CorrelationTimeoutSeconds:   30,
		TimeoutSweepIntervalSeconds: 1,
		TimeoutBatchSize:            200,
		TimeoutCountsAsFailure:      true,
		MetricPrefix:                "orders_tracking",
		ExecuteOriginalScenarioRun:  false,
	}
}

func withLoggerConfigTrackedContext() loadstrike.LoadStrikeContext {
	return loadstrike.Create().
		AddScenario(withLoggerConfigBaselineScenario("orders.tracked").WithCrossPlatformTracking(withLoggerConfigTrackingConfiguration())).
		WithRunnerKey(withLoggerConfigRunnerKey).
		WithoutReports().
		BuildContext()
}

func withLoggerConfigBuildTrackingPayload() loadstrike.TrackingPayload {
	builder := loadstrike.TrackingPayloadBuilder{
		Headers:            map[string]string{"X-Correlation-Id": "ord-1"},
		ContentType:        "application/json",
		MessagePayloadType: "json",
	}
	builder.SetBody(map[string]any{"trackingId": "ord-1"})
	return builder.Build()
}

func withLoggerConfigRunResult() loadstrike.LoadStrikeRunResult {
	var result loadstrike.LoadStrikeRunResult
	return result
}

func withLoggerConfigScenarioStats() loadstrike.LoadStrikeScenarioStats {
	var stats loadstrike.LoadStrikeScenarioStats
	return stats
}

func withLoggerConfigWriteTempConfigFiles() withLoggerConfigTempConfigPaths {
	return withLoggerConfigTempConfigPaths{
		ConfigPath: "method-reference.loadstrike.config.json",
		InfraPath:  "method-reference.loadstrike.infra.json",
	}
}

        // Attach a logger factory directly on the runner.
func (reference WithLoggerConfigMethodReference) AttachInlineLoggerFactoryExample() any {
    return withLoggerConfigBaseRunner().WithLoggerConfig(func() loadstrike.LoggerConfiguration { return loadstrike.LoggerConfiguration{"logger": "orders"} })
}

// Attach the logger factory before building a reusable context.
func (reference WithLoggerConfigMethodReference) AttachLoggerBeforeContextBuildExample() any {
    return withLoggerConfigBaseRunner().WithLoggerConfig(func() loadstrike.LoggerConfiguration { return loadstrike.LoggerConfiguration{"logger": "orders"} }).BuildContext()
}
