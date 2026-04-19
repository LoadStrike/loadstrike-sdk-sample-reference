package loadstrike_runner

        import (
        	"time"

        	loadstrike "loadstrike.com/sdk/go"
        )

        const loadConfigRunnerKey = "runner_dummy_orders_reference"

type LoadConfigMethodReference struct{}

type loadConfigTempConfigPaths struct {
	ConfigPath string
	InfraPath  string
}

type loadConfigOrdersReportingSink struct{}

func newLoadConfigOrdersReportingSink() loadConfigOrdersReportingSink {
	return loadConfigOrdersReportingSink{}
}
func (loadConfigOrdersReportingSink) SinkName() string {
	return "orders-sample-sink"
}
func (loadConfigOrdersReportingSink) Init(loadstrike.LoadStrikeBaseContext, loadstrike.IConfiguration) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (loadConfigOrdersReportingSink) Start(loadstrike.LoadStrikeSessionStartInfo) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (loadConfigOrdersReportingSink) SaveRealtimeStats([]loadstrike.LoadStrikeScenarioStats) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (loadConfigOrdersReportingSink) SaveRealtimeMetrics(loadstrike.LoadStrikeMetricStats) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (loadConfigOrdersReportingSink) SaveRunResult(loadstrike.LoadStrikeRunResult) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (loadConfigOrdersReportingSink) Stop() loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (loadConfigOrdersReportingSink) Dispose() {}

type loadConfigOrdersRuntimePolicy struct{}

func newLoadConfigOrdersRuntimePolicy() loadConfigOrdersRuntimePolicy {
	return loadConfigOrdersRuntimePolicy{}
}
func (loadConfigOrdersRuntimePolicy) ShouldRunScenario(string) loadstrike.LoadStrikeBoolTask {
	return loadstrike.TaskFromBool(true)
}
func (loadConfigOrdersRuntimePolicy) BeforeScenario(string) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (loadConfigOrdersRuntimePolicy) AfterScenario(string, loadstrike.LoadStrikeScenarioRuntime) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (loadConfigOrdersRuntimePolicy) BeforeStep(string, string) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (loadConfigOrdersRuntimePolicy) AfterStep(string, string, loadstrike.LoadStrikeReply) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}

type loadConfigOrdersWorkerPlugin struct{}

func newLoadConfigOrdersWorkerPlugin() loadConfigOrdersWorkerPlugin {
	return loadConfigOrdersWorkerPlugin{}
}
func (loadConfigOrdersWorkerPlugin) PluginName() string {
	return "orders-sample-plugin"
}
func (loadConfigOrdersWorkerPlugin) Init(loadstrike.LoadStrikeBaseContext, loadstrike.IConfiguration) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (loadConfigOrdersWorkerPlugin) Start(loadstrike.LoadStrikeSessionStartInfo) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (loadConfigOrdersWorkerPlugin) GetData(loadstrike.LoadStrikeRunResult) loadstrike.LoadStrikeValueTask[loadstrike.LoadStrikePluginData] {
	return loadstrike.TaskFromResult(loadstrike.CreatePluginData("orders-sample-plugin"))
}
func (loadConfigOrdersWorkerPlugin) Stop() loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (loadConfigOrdersWorkerPlugin) Dispose() loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}

func loadConfigPerformOrderGetReply() loadstrike.LoadStrikeReply {
	return loadstrike.LoadStrikeResponse.Ok("200", int64(128), "ok", loadstrike.TimeSpan(3*time.Millisecond))
}

func loadConfigExecuteOrderGet(context loadstrike.LoadStrikeScenarioContext) loadstrike.LoadStrikeReply {
	return loadstrike.LoadStrikeStep.Run("get-order", context, func(loadstrike.LoadStrikeScenarioContext) loadstrike.LoadStrikeReply {
		return loadConfigPerformOrderGetReply()
	})
}

func loadConfigExecuteOrderGetAsync(context loadstrike.LoadStrikeScenarioContext) loadstrike.LoadStrikeValueTask[loadstrike.LoadStrikeReply] {
	return loadstrike.TaskFromResult(loadConfigExecuteOrderGet(context))
}

func loadConfigBaselineScenario(names ...string) loadstrike.LoadStrikeScenario {
	name := "orders.get-by-id"
	if len(names) > 0 && names[0] != "" {
		name = names[0]
	}
	return loadstrike.CreateScenario(name, loadConfigExecuteOrderGet).
		WithLoadSimulations(loadstrike.LoadStrikeSimulation.IterationsForConstant(1, 1)).
		WithoutWarmUp()
}

func loadConfigBaseRunner() loadstrike.LoadStrikeRunner {
	return loadstrike.Create().
		AddScenario(loadConfigBaselineScenario()).
		WithRunnerKey(loadConfigRunnerKey).
		WithTestSuite("orders-reference").
		WithTestName("orders-get-by-id").
		WithoutReports()
}

func loadConfigBaseContext() loadstrike.LoadStrikeContext {
	return loadConfigBaseRunner().BuildContext()
}

func loadConfigHttpSource() *loadstrike.EndpointSpec {
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

func loadConfigHttpDestination() *loadstrike.EndpointSpec {
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

func loadConfigTrackingConfiguration() *loadstrike.TrackingConfigurationSpec {
	return &loadstrike.TrackingConfigurationSpec{
		Source:                      loadConfigHttpSource(),
		Destination:                 loadConfigHttpDestination(),
		RunMode:                     "GenerateAndCorrelate",
		CorrelationTimeoutSeconds:   30,
		TimeoutSweepIntervalSeconds: 1,
		TimeoutBatchSize:            200,
		TimeoutCountsAsFailure:      true,
		MetricPrefix:                "orders_tracking",
		ExecuteOriginalScenarioRun:  false,
	}
}

func loadConfigTrackedContext() loadstrike.LoadStrikeContext {
	return loadstrike.Create().
		AddScenario(loadConfigBaselineScenario("orders.tracked").WithCrossPlatformTracking(loadConfigTrackingConfiguration())).
		WithRunnerKey(loadConfigRunnerKey).
		WithoutReports().
		BuildContext()
}

func loadConfigBuildTrackingPayload() loadstrike.TrackingPayload {
	builder := loadstrike.TrackingPayloadBuilder{
		Headers:            map[string]string{"X-Correlation-Id": "ord-1"},
		ContentType:        "application/json",
		MessagePayloadType: "json",
	}
	builder.SetBody(map[string]any{"trackingId": "ord-1"})
	return builder.Build()
}

func loadConfigRunResult() loadstrike.LoadStrikeRunResult {
	var result loadstrike.LoadStrikeRunResult
	return result
}

func loadConfigScenarioStats() loadstrike.LoadStrikeScenarioStats {
	var stats loadstrike.LoadStrikeScenarioStats
	return stats
}

func loadConfigWriteTempConfigFiles() loadConfigTempConfigPaths {
	return loadConfigTempConfigPaths{
		ConfigPath: "method-reference.loadstrike.config.json",
		InfraPath:  "method-reference.loadstrike.infra.json",
	}
}

        // Point the runner at a temporary JSON file generated inside the sample.
        func (reference LoadConfigMethodReference) UseGeneratedConfigFileExample() any {
            paths := loadConfigWriteTempConfigFiles()
return loadConfigBaseRunner().LoadConfig(paths.ConfigPath)
        }

// Load the file and continue into a built context.
        func (reference LoadConfigMethodReference) BuildContextFromGeneratedConfigExample() any {
            paths := loadConfigWriteTempConfigFiles()
return loadConfigBaseRunner().LoadConfig(paths.ConfigPath).BuildContext()
        }
