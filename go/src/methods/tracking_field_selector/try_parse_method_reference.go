package tracking_field_selector

        import (
        	"time"

        	loadstrike "loadstrike.com/sdk/go"
        )

        const tryParseRunnerKey = "runner_dummy_orders_reference"

type TryParseMethodReference struct{}

type tryParseTempConfigPaths struct {
	ConfigPath string
	InfraPath  string
}

type tryParseOrdersReportingSink struct{}

func newTryParseOrdersReportingSink() tryParseOrdersReportingSink {
	return tryParseOrdersReportingSink{}
}
func (tryParseOrdersReportingSink) SinkName() string {
	return "orders-sample-sink"
}
func (tryParseOrdersReportingSink) Init(loadstrike.LoadStrikeBaseContext, loadstrike.IConfiguration) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (tryParseOrdersReportingSink) Start(loadstrike.LoadStrikeSessionStartInfo) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (tryParseOrdersReportingSink) SaveRealtimeStats([]loadstrike.LoadStrikeScenarioStats) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (tryParseOrdersReportingSink) SaveRealtimeMetrics(loadstrike.LoadStrikeMetricStats) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (tryParseOrdersReportingSink) SaveRunResult(loadstrike.LoadStrikeRunResult) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (tryParseOrdersReportingSink) Stop() loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (tryParseOrdersReportingSink) Dispose() {}

type tryParseOrdersRuntimePolicy struct{}

func newTryParseOrdersRuntimePolicy() tryParseOrdersRuntimePolicy {
	return tryParseOrdersRuntimePolicy{}
}
func (tryParseOrdersRuntimePolicy) ShouldRunScenario(string) loadstrike.LoadStrikeBoolTask {
	return loadstrike.TaskFromBool(true)
}
func (tryParseOrdersRuntimePolicy) BeforeScenario(string) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (tryParseOrdersRuntimePolicy) AfterScenario(string, loadstrike.LoadStrikeScenarioRuntime) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (tryParseOrdersRuntimePolicy) BeforeStep(string, string) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (tryParseOrdersRuntimePolicy) AfterStep(string, string, loadstrike.LoadStrikeReply) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}

type tryParseOrdersWorkerPlugin struct{}

func newTryParseOrdersWorkerPlugin() tryParseOrdersWorkerPlugin {
	return tryParseOrdersWorkerPlugin{}
}
func (tryParseOrdersWorkerPlugin) PluginName() string {
	return "orders-sample-plugin"
}
func (tryParseOrdersWorkerPlugin) Init(loadstrike.LoadStrikeBaseContext, loadstrike.IConfiguration) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (tryParseOrdersWorkerPlugin) Start(loadstrike.LoadStrikeSessionStartInfo) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (tryParseOrdersWorkerPlugin) GetData(loadstrike.LoadStrikeRunResult) loadstrike.LoadStrikeValueTask[loadstrike.LoadStrikePluginData] {
	return loadstrike.TaskFromResult(loadstrike.CreatePluginData("orders-sample-plugin"))
}
func (tryParseOrdersWorkerPlugin) Stop() loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (tryParseOrdersWorkerPlugin) Dispose() loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}

func tryParsePerformOrderGetReply() loadstrike.LoadStrikeReply {
	return loadstrike.LoadStrikeResponse.Ok("200", int64(128), "ok", loadstrike.TimeSpan(3*time.Millisecond))
}

func tryParseExecuteOrderGet(context loadstrike.LoadStrikeScenarioContext) loadstrike.LoadStrikeReply {
	return loadstrike.LoadStrikeStep.Run("get-order", context, func(loadstrike.LoadStrikeScenarioContext) loadstrike.LoadStrikeReply {
		return tryParsePerformOrderGetReply()
	})
}

func tryParseExecuteOrderGetAsync(context loadstrike.LoadStrikeScenarioContext) loadstrike.LoadStrikeValueTask[loadstrike.LoadStrikeReply] {
	return loadstrike.TaskFromResult(tryParseExecuteOrderGet(context))
}

func tryParseBaselineScenario(names ...string) loadstrike.LoadStrikeScenario {
	name := "orders.get-by-id"
	if len(names) > 0 && names[0] != "" {
		name = names[0]
	}
	return loadstrike.CreateScenario(name, tryParseExecuteOrderGet).
		WithLoadSimulations(loadstrike.LoadStrikeSimulation.IterationsForConstant(1, 1)).
		WithoutWarmUp()
}

func tryParseBaseRunner() loadstrike.LoadStrikeRunner {
	return loadstrike.Create().
		AddScenario(tryParseBaselineScenario()).
		WithRunnerKey(tryParseRunnerKey).
		WithTestSuite("orders-reference").
		WithTestName("orders-get-by-id").
		WithoutReports()
}

func tryParseBaseContext() loadstrike.LoadStrikeContext {
	return tryParseBaseRunner().BuildContext()
}

func tryParseHttpSource() *loadstrike.EndpointSpec {
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

func tryParseHttpDestination() *loadstrike.EndpointSpec {
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

func tryParseTrackingConfiguration() *loadstrike.TrackingConfigurationSpec {
	return &loadstrike.TrackingConfigurationSpec{
		Source:                      tryParseHttpSource(),
		Destination:                 tryParseHttpDestination(),
		RunMode:                     "GenerateAndCorrelate",
		CorrelationTimeoutSeconds:   30,
		TimeoutSweepIntervalSeconds: 1,
		TimeoutBatchSize:            200,
		TimeoutCountsAsFailure:      true,
		MetricPrefix:                "orders_tracking",
		ExecuteOriginalScenarioRun:  false,
	}
}

func tryParseTrackedContext() loadstrike.LoadStrikeContext {
	return loadstrike.Create().
		AddScenario(tryParseBaselineScenario("orders.tracked").WithCrossPlatformTracking(tryParseTrackingConfiguration())).
		WithRunnerKey(tryParseRunnerKey).
		WithoutReports().
		BuildContext()
}

func tryParseBuildTrackingPayload() loadstrike.TrackingPayload {
	builder := loadstrike.TrackingPayloadBuilder{
		Headers:            map[string]string{"X-Correlation-Id": "ord-1"},
		ContentType:        "application/json",
		MessagePayloadType: "json",
	}
	builder.SetBody(map[string]any{"trackingId": "ord-1"})
	return builder.Build()
}

func tryParseRunResult() loadstrike.LoadStrikeRunResult {
	var result loadstrike.LoadStrikeRunResult
	return result
}

func tryParseScenarioStats() loadstrike.LoadStrikeScenarioStats {
	var stats loadstrike.LoadStrikeScenarioStats
	return stats
}

func tryParseWriteTempConfigFiles() tryParseTempConfigPaths {
	return tryParseTempConfigPaths{
		ConfigPath: "method-reference.loadstrike.config.json",
		InfraPath:  "method-reference.loadstrike.infra.json",
	}
}

        // Call the correlation helper directly with a concrete value.
        func (reference TryParseMethodReference) CreateCorrelationExample() any {
            parsed, selector := loadstrike.TrackingFieldSelector{}.TryParse("header:x-id")
return map[string]any{"parsed": parsed, "selector": selector}
        }

// Show how the same helper fits into the tracked source/destination example.
        func (reference TryParseMethodReference) UseCorrelationExampleInTrackedFlow() any {
            parsed, selector := loadstrike.TrackingFieldSelector{}.TryParse("header:x-id")
return map[string]any{"parsed": parsed, "selector": selector}
        }
