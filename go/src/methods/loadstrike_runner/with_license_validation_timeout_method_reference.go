package loadstrike_runner

        import (
        	"time"

        	loadstrike "loadstrike.com/sdk/go"
        )

        const withLicenseValidationTimeoutRunnerKey = "runner_dummy_orders_reference"

type WithLicenseValidationTimeoutMethodReference struct{}

type withLicenseValidationTimeoutTempConfigPaths struct {
	ConfigPath string
	InfraPath  string
}

type withLicenseValidationTimeoutOrdersReportingSink struct{}

func newWithLicenseValidationTimeoutOrdersReportingSink() withLicenseValidationTimeoutOrdersReportingSink {
	return withLicenseValidationTimeoutOrdersReportingSink{}
}
func (withLicenseValidationTimeoutOrdersReportingSink) SinkName() string {
	return "orders-sample-sink"
}
func (withLicenseValidationTimeoutOrdersReportingSink) Init(loadstrike.LoadStrikeBaseContext, loadstrike.IConfiguration) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withLicenseValidationTimeoutOrdersReportingSink) Start(loadstrike.LoadStrikeSessionStartInfo) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withLicenseValidationTimeoutOrdersReportingSink) SaveRealtimeStats([]loadstrike.LoadStrikeScenarioStats) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withLicenseValidationTimeoutOrdersReportingSink) SaveRealtimeMetrics(loadstrike.LoadStrikeMetricStats) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withLicenseValidationTimeoutOrdersReportingSink) SaveRunResult(loadstrike.LoadStrikeRunResult) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withLicenseValidationTimeoutOrdersReportingSink) Stop() loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withLicenseValidationTimeoutOrdersReportingSink) Dispose() {}

type withLicenseValidationTimeoutOrdersRuntimePolicy struct{}

func newWithLicenseValidationTimeoutOrdersRuntimePolicy() withLicenseValidationTimeoutOrdersRuntimePolicy {
	return withLicenseValidationTimeoutOrdersRuntimePolicy{}
}
func (withLicenseValidationTimeoutOrdersRuntimePolicy) ShouldRunScenario(string) loadstrike.LoadStrikeBoolTask {
	return loadstrike.TaskFromBool(true)
}
func (withLicenseValidationTimeoutOrdersRuntimePolicy) BeforeScenario(string) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withLicenseValidationTimeoutOrdersRuntimePolicy) AfterScenario(string, loadstrike.LoadStrikeScenarioRuntime) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withLicenseValidationTimeoutOrdersRuntimePolicy) BeforeStep(string, string) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withLicenseValidationTimeoutOrdersRuntimePolicy) AfterStep(string, string, loadstrike.LoadStrikeReply) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}

type withLicenseValidationTimeoutOrdersWorkerPlugin struct{}

func newWithLicenseValidationTimeoutOrdersWorkerPlugin() withLicenseValidationTimeoutOrdersWorkerPlugin {
	return withLicenseValidationTimeoutOrdersWorkerPlugin{}
}
func (withLicenseValidationTimeoutOrdersWorkerPlugin) PluginName() string {
	return "orders-sample-plugin"
}
func (withLicenseValidationTimeoutOrdersWorkerPlugin) Init(loadstrike.LoadStrikeBaseContext, loadstrike.IConfiguration) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withLicenseValidationTimeoutOrdersWorkerPlugin) Start(loadstrike.LoadStrikeSessionStartInfo) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withLicenseValidationTimeoutOrdersWorkerPlugin) GetData(loadstrike.LoadStrikeRunResult) loadstrike.LoadStrikeValueTask[loadstrike.LoadStrikePluginData] {
	return loadstrike.TaskFromResult(loadstrike.CreatePluginData("orders-sample-plugin"))
}
func (withLicenseValidationTimeoutOrdersWorkerPlugin) Stop() loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withLicenseValidationTimeoutOrdersWorkerPlugin) Dispose() loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}

func withLicenseValidationTimeoutPerformOrderGetReply() loadstrike.LoadStrikeReply {
	return loadstrike.LoadStrikeResponse.Ok("200", int64(128), "ok", loadstrike.TimeSpan(3*time.Millisecond))
}

func withLicenseValidationTimeoutExecuteOrderGet(context loadstrike.LoadStrikeScenarioContext) loadstrike.LoadStrikeReply {
	return loadstrike.LoadStrikeStep.Run("get-order", context, func(loadstrike.LoadStrikeScenarioContext) loadstrike.LoadStrikeReply {
		return withLicenseValidationTimeoutPerformOrderGetReply()
	})
}

func withLicenseValidationTimeoutExecuteOrderGetAsync(context loadstrike.LoadStrikeScenarioContext) loadstrike.LoadStrikeValueTask[loadstrike.LoadStrikeReply] {
	return loadstrike.TaskFromResult(withLicenseValidationTimeoutExecuteOrderGet(context))
}

func withLicenseValidationTimeoutBaselineScenario(names ...string) loadstrike.LoadStrikeScenario {
	name := "orders.get-by-id"
	if len(names) > 0 && names[0] != "" {
		name = names[0]
	}
	return loadstrike.CreateScenario(name, withLicenseValidationTimeoutExecuteOrderGet).
		WithLoadSimulations(loadstrike.LoadStrikeSimulation.IterationsForConstant(1, 1)).
		WithoutWarmUp()
}

func withLicenseValidationTimeoutBaseRunner() loadstrike.LoadStrikeRunner {
	return loadstrike.Create().
		AddScenario(withLicenseValidationTimeoutBaselineScenario()).
		WithRunnerKey(withLicenseValidationTimeoutRunnerKey).
		WithTestSuite("orders-reference").
		WithTestName("orders-get-by-id").
		WithoutReports()
}

func withLicenseValidationTimeoutBaseContext() loadstrike.LoadStrikeContext {
	return withLicenseValidationTimeoutBaseRunner().BuildContext()
}

func withLicenseValidationTimeoutHttpSource() *loadstrike.EndpointSpec {
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

func withLicenseValidationTimeoutHttpDestination() *loadstrike.EndpointSpec {
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

func withLicenseValidationTimeoutTrackingConfiguration() *loadstrike.TrackingConfigurationSpec {
	return &loadstrike.TrackingConfigurationSpec{
		Source:                      withLicenseValidationTimeoutHttpSource(),
		Destination:                 withLicenseValidationTimeoutHttpDestination(),
		RunMode:                     "GenerateAndCorrelate",
		CorrelationTimeoutSeconds:   30,
		TimeoutSweepIntervalSeconds: 1,
		TimeoutBatchSize:            200,
		TimeoutCountsAsFailure:      true,
		MetricPrefix:                "orders_tracking",
		ExecuteOriginalScenarioRun:  false,
	}
}

func withLicenseValidationTimeoutTrackedContext() loadstrike.LoadStrikeContext {
	return loadstrike.Create().
		AddScenario(withLicenseValidationTimeoutBaselineScenario("orders.tracked").WithCrossPlatformTracking(withLicenseValidationTimeoutTrackingConfiguration())).
		WithRunnerKey(withLicenseValidationTimeoutRunnerKey).
		WithoutReports().
		BuildContext()
}

func withLicenseValidationTimeoutBuildTrackingPayload() loadstrike.TrackingPayload {
	builder := loadstrike.TrackingPayloadBuilder{
		Headers:            map[string]string{"X-Correlation-Id": "ord-1"},
		ContentType:        "application/json",
		MessagePayloadType: "json",
	}
	builder.SetBody(map[string]any{"trackingId": "ord-1"})
	return builder.Build()
}

func withLicenseValidationTimeoutRunResult() loadstrike.LoadStrikeRunResult {
	var result loadstrike.LoadStrikeRunResult
	return result
}

func withLicenseValidationTimeoutScenarioStats() loadstrike.LoadStrikeScenarioStats {
	var stats loadstrike.LoadStrikeScenarioStats
	return stats
}

func withLicenseValidationTimeoutWriteTempConfigFiles() withLicenseValidationTimeoutTempConfigPaths {
	return withLicenseValidationTimeoutTempConfigPaths{
		ConfigPath: "method-reference.loadstrike.config.json",
		InfraPath:  "method-reference.loadstrike.infra.json",
	}
}

        // Apply a shorter timeout to the runner.
func (reference WithLicenseValidationTimeoutMethodReference) ApplyShortTimeoutExample() any {
    return withLicenseValidationTimeoutBaseRunner().WithLicenseValidationTimeout(loadstrike.DurationFromSeconds(5))
}

// Apply a longer timeout to the runner.
func (reference WithLicenseValidationTimeoutMethodReference) ApplyLongTimeoutExample() any {
    return withLicenseValidationTimeoutBaseRunner().WithLicenseValidationTimeout(loadstrike.DurationFromSeconds(15))
}
