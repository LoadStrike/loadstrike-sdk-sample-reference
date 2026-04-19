package loadstrike_runner

        import (
        	"time"

        	loadstrike "loadstrike.com/sdk/go"
        )

        const loadInfraConfigRunnerKey = "runner_dummy_orders_reference"

type LoadInfraConfigMethodReference struct{}

type loadInfraConfigTempConfigPaths struct {
	ConfigPath string
	InfraPath  string
}

type loadInfraConfigOrdersReportingSink struct{}

func newLoadInfraConfigOrdersReportingSink() loadInfraConfigOrdersReportingSink {
	return loadInfraConfigOrdersReportingSink{}
}
func (loadInfraConfigOrdersReportingSink) SinkName() string {
	return "orders-sample-sink"
}
func (loadInfraConfigOrdersReportingSink) Init(loadstrike.LoadStrikeBaseContext, loadstrike.IConfiguration) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (loadInfraConfigOrdersReportingSink) Start(loadstrike.LoadStrikeSessionStartInfo) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (loadInfraConfigOrdersReportingSink) SaveRealtimeStats([]loadstrike.LoadStrikeScenarioStats) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (loadInfraConfigOrdersReportingSink) SaveRealtimeMetrics(loadstrike.LoadStrikeMetricStats) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (loadInfraConfigOrdersReportingSink) SaveRunResult(loadstrike.LoadStrikeRunResult) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (loadInfraConfigOrdersReportingSink) Stop() loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (loadInfraConfigOrdersReportingSink) Dispose() {}

type loadInfraConfigOrdersRuntimePolicy struct{}

func newLoadInfraConfigOrdersRuntimePolicy() loadInfraConfigOrdersRuntimePolicy {
	return loadInfraConfigOrdersRuntimePolicy{}
}
func (loadInfraConfigOrdersRuntimePolicy) ShouldRunScenario(string) loadstrike.LoadStrikeBoolTask {
	return loadstrike.TaskFromBool(true)
}
func (loadInfraConfigOrdersRuntimePolicy) BeforeScenario(string) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (loadInfraConfigOrdersRuntimePolicy) AfterScenario(string, loadstrike.LoadStrikeScenarioRuntime) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (loadInfraConfigOrdersRuntimePolicy) BeforeStep(string, string) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (loadInfraConfigOrdersRuntimePolicy) AfterStep(string, string, loadstrike.LoadStrikeReply) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}

type loadInfraConfigOrdersWorkerPlugin struct{}

func newLoadInfraConfigOrdersWorkerPlugin() loadInfraConfigOrdersWorkerPlugin {
	return loadInfraConfigOrdersWorkerPlugin{}
}
func (loadInfraConfigOrdersWorkerPlugin) PluginName() string {
	return "orders-sample-plugin"
}
func (loadInfraConfigOrdersWorkerPlugin) Init(loadstrike.LoadStrikeBaseContext, loadstrike.IConfiguration) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (loadInfraConfigOrdersWorkerPlugin) Start(loadstrike.LoadStrikeSessionStartInfo) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (loadInfraConfigOrdersWorkerPlugin) GetData(loadstrike.LoadStrikeRunResult) loadstrike.LoadStrikeValueTask[loadstrike.LoadStrikePluginData] {
	return loadstrike.TaskFromResult(loadstrike.CreatePluginData("orders-sample-plugin"))
}
func (loadInfraConfigOrdersWorkerPlugin) Stop() loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (loadInfraConfigOrdersWorkerPlugin) Dispose() loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}

func loadInfraConfigPerformOrderGetReply() loadstrike.LoadStrikeReply {
	return loadstrike.LoadStrikeResponse.Ok("200", int64(128), "ok", loadstrike.TimeSpan(3*time.Millisecond))
}

func loadInfraConfigExecuteOrderGet(context loadstrike.LoadStrikeScenarioContext) loadstrike.LoadStrikeReply {
	return loadstrike.LoadStrikeStep.Run("get-order", context, func(loadstrike.LoadStrikeScenarioContext) loadstrike.LoadStrikeReply {
		return loadInfraConfigPerformOrderGetReply()
	})
}

func loadInfraConfigExecuteOrderGetAsync(context loadstrike.LoadStrikeScenarioContext) loadstrike.LoadStrikeValueTask[loadstrike.LoadStrikeReply] {
	return loadstrike.TaskFromResult(loadInfraConfigExecuteOrderGet(context))
}

func loadInfraConfigBaselineScenario(names ...string) loadstrike.LoadStrikeScenario {
	name := "orders.get-by-id"
	if len(names) > 0 && names[0] != "" {
		name = names[0]
	}
	return loadstrike.CreateScenario(name, loadInfraConfigExecuteOrderGet).
		WithLoadSimulations(loadstrike.LoadStrikeSimulation.IterationsForConstant(1, 1)).
		WithoutWarmUp()
}

func loadInfraConfigBaseRunner() loadstrike.LoadStrikeRunner {
	return loadstrike.Create().
		AddScenario(loadInfraConfigBaselineScenario()).
		WithRunnerKey(loadInfraConfigRunnerKey).
		WithTestSuite("orders-reference").
		WithTestName("orders-get-by-id").
		WithoutReports()
}

func loadInfraConfigBaseContext() loadstrike.LoadStrikeContext {
	return loadInfraConfigBaseRunner().BuildContext()
}

func loadInfraConfigHttpSource() *loadstrike.EndpointSpec {
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

func loadInfraConfigHttpDestination() *loadstrike.EndpointSpec {
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

func loadInfraConfigTrackingConfiguration() *loadstrike.TrackingConfigurationSpec {
	return &loadstrike.TrackingConfigurationSpec{
		Source:                      loadInfraConfigHttpSource(),
		Destination:                 loadInfraConfigHttpDestination(),
		RunMode:                     "GenerateAndCorrelate",
		CorrelationTimeoutSeconds:   30,
		TimeoutSweepIntervalSeconds: 1,
		TimeoutBatchSize:            200,
		TimeoutCountsAsFailure:      true,
		MetricPrefix:                "orders_tracking",
		ExecuteOriginalScenarioRun:  false,
	}
}

func loadInfraConfigTrackedContext() loadstrike.LoadStrikeContext {
	return loadstrike.Create().
		AddScenario(loadInfraConfigBaselineScenario("orders.tracked").WithCrossPlatformTracking(loadInfraConfigTrackingConfiguration())).
		WithRunnerKey(loadInfraConfigRunnerKey).
		WithoutReports().
		BuildContext()
}

func loadInfraConfigBuildTrackingPayload() loadstrike.TrackingPayload {
	builder := loadstrike.TrackingPayloadBuilder{
		Headers:            map[string]string{"X-Correlation-Id": "ord-1"},
		ContentType:        "application/json",
		MessagePayloadType: "json",
	}
	builder.SetBody(map[string]any{"trackingId": "ord-1"})
	return builder.Build()
}

func loadInfraConfigRunResult() loadstrike.LoadStrikeRunResult {
	var result loadstrike.LoadStrikeRunResult
	return result
}

func loadInfraConfigScenarioStats() loadstrike.LoadStrikeScenarioStats {
	var stats loadstrike.LoadStrikeScenarioStats
	return stats
}

func loadInfraConfigWriteTempConfigFiles() loadInfraConfigTempConfigPaths {
	return loadInfraConfigTempConfigPaths{
		ConfigPath: "method-reference.loadstrike.config.json",
		InfraPath:  "method-reference.loadstrike.infra.json",
	}
}

        // Point the runner at a temporary JSON file generated inside the sample.
        func (reference LoadInfraConfigMethodReference) UseGeneratedConfigFileExample() any {
            paths := loadInfraConfigWriteTempConfigFiles()
return loadInfraConfigBaseRunner().LoadInfraConfig(paths.InfraPath)
        }

// Load the file and continue into a built context.
        func (reference LoadInfraConfigMethodReference) BuildContextFromGeneratedConfigExample() any {
            paths := loadInfraConfigWriteTempConfigFiles()
return loadInfraConfigBaseRunner().LoadInfraConfig(paths.InfraPath).BuildContext()
        }
