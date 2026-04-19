package loadstrike_runner

        import (
        	"time"

        	loadstrike "loadstrike.com/sdk/go"
        )

        const withRuntimePoliciesRunnerKey = "runner_dummy_orders_reference"

type WithRuntimePoliciesMethodReference struct{}

type withRuntimePoliciesTempConfigPaths struct {
	ConfigPath string
	InfraPath  string
}

type withRuntimePoliciesOrdersReportingSink struct{}

func newWithRuntimePoliciesOrdersReportingSink() withRuntimePoliciesOrdersReportingSink {
	return withRuntimePoliciesOrdersReportingSink{}
}
func (withRuntimePoliciesOrdersReportingSink) SinkName() string {
	return "orders-sample-sink"
}
func (withRuntimePoliciesOrdersReportingSink) Init(loadstrike.LoadStrikeBaseContext, loadstrike.IConfiguration) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withRuntimePoliciesOrdersReportingSink) Start(loadstrike.LoadStrikeSessionStartInfo) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withRuntimePoliciesOrdersReportingSink) SaveRealtimeStats([]loadstrike.LoadStrikeScenarioStats) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withRuntimePoliciesOrdersReportingSink) SaveRealtimeMetrics(loadstrike.LoadStrikeMetricStats) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withRuntimePoliciesOrdersReportingSink) SaveRunResult(loadstrike.LoadStrikeRunResult) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withRuntimePoliciesOrdersReportingSink) Stop() loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withRuntimePoliciesOrdersReportingSink) Dispose() {}

type withRuntimePoliciesOrdersRuntimePolicy struct{}

func newWithRuntimePoliciesOrdersRuntimePolicy() withRuntimePoliciesOrdersRuntimePolicy {
	return withRuntimePoliciesOrdersRuntimePolicy{}
}
func (withRuntimePoliciesOrdersRuntimePolicy) ShouldRunScenario(string) loadstrike.LoadStrikeBoolTask {
	return loadstrike.TaskFromBool(true)
}
func (withRuntimePoliciesOrdersRuntimePolicy) BeforeScenario(string) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withRuntimePoliciesOrdersRuntimePolicy) AfterScenario(string, loadstrike.LoadStrikeScenarioRuntime) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withRuntimePoliciesOrdersRuntimePolicy) BeforeStep(string, string) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withRuntimePoliciesOrdersRuntimePolicy) AfterStep(string, string, loadstrike.LoadStrikeReply) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}

type withRuntimePoliciesOrdersWorkerPlugin struct{}

func newWithRuntimePoliciesOrdersWorkerPlugin() withRuntimePoliciesOrdersWorkerPlugin {
	return withRuntimePoliciesOrdersWorkerPlugin{}
}
func (withRuntimePoliciesOrdersWorkerPlugin) PluginName() string {
	return "orders-sample-plugin"
}
func (withRuntimePoliciesOrdersWorkerPlugin) Init(loadstrike.LoadStrikeBaseContext, loadstrike.IConfiguration) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withRuntimePoliciesOrdersWorkerPlugin) Start(loadstrike.LoadStrikeSessionStartInfo) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withRuntimePoliciesOrdersWorkerPlugin) GetData(loadstrike.LoadStrikeRunResult) loadstrike.LoadStrikeValueTask[loadstrike.LoadStrikePluginData] {
	return loadstrike.TaskFromResult(loadstrike.CreatePluginData("orders-sample-plugin"))
}
func (withRuntimePoliciesOrdersWorkerPlugin) Stop() loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (withRuntimePoliciesOrdersWorkerPlugin) Dispose() loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}

func withRuntimePoliciesPerformOrderGetReply() loadstrike.LoadStrikeReply {
	return loadstrike.LoadStrikeResponse.Ok("200", int64(128), "ok", loadstrike.TimeSpan(3*time.Millisecond))
}

func withRuntimePoliciesExecuteOrderGet(context loadstrike.LoadStrikeScenarioContext) loadstrike.LoadStrikeReply {
	return loadstrike.LoadStrikeStep.Run("get-order", context, func(loadstrike.LoadStrikeScenarioContext) loadstrike.LoadStrikeReply {
		return withRuntimePoliciesPerformOrderGetReply()
	})
}

func withRuntimePoliciesExecuteOrderGetAsync(context loadstrike.LoadStrikeScenarioContext) loadstrike.LoadStrikeValueTask[loadstrike.LoadStrikeReply] {
	return loadstrike.TaskFromResult(withRuntimePoliciesExecuteOrderGet(context))
}

func withRuntimePoliciesBaselineScenario(names ...string) loadstrike.LoadStrikeScenario {
	name := "orders.get-by-id"
	if len(names) > 0 && names[0] != "" {
		name = names[0]
	}
	return loadstrike.CreateScenario(name, withRuntimePoliciesExecuteOrderGet).
		WithLoadSimulations(loadstrike.LoadStrikeSimulation.IterationsForConstant(1, 1)).
		WithoutWarmUp()
}

func withRuntimePoliciesBaseRunner() loadstrike.LoadStrikeRunner {
	return loadstrike.Create().
		AddScenario(withRuntimePoliciesBaselineScenario()).
		WithRunnerKey(withRuntimePoliciesRunnerKey).
		WithTestSuite("orders-reference").
		WithTestName("orders-get-by-id").
		WithoutReports()
}

func withRuntimePoliciesBaseContext() loadstrike.LoadStrikeContext {
	return withRuntimePoliciesBaseRunner().BuildContext()
}

func withRuntimePoliciesHttpSource() *loadstrike.EndpointSpec {
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

func withRuntimePoliciesHttpDestination() *loadstrike.EndpointSpec {
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

func withRuntimePoliciesTrackingConfiguration() *loadstrike.TrackingConfigurationSpec {
	return &loadstrike.TrackingConfigurationSpec{
		Source:                      withRuntimePoliciesHttpSource(),
		Destination:                 withRuntimePoliciesHttpDestination(),
		RunMode:                     "GenerateAndCorrelate",
		CorrelationTimeoutSeconds:   30,
		TimeoutSweepIntervalSeconds: 1,
		TimeoutBatchSize:            200,
		TimeoutCountsAsFailure:      true,
		MetricPrefix:                "orders_tracking",
		ExecuteOriginalScenarioRun:  false,
	}
}

func withRuntimePoliciesTrackedContext() loadstrike.LoadStrikeContext {
	return loadstrike.Create().
		AddScenario(withRuntimePoliciesBaselineScenario("orders.tracked").WithCrossPlatformTracking(withRuntimePoliciesTrackingConfiguration())).
		WithRunnerKey(withRuntimePoliciesRunnerKey).
		WithoutReports().
		BuildContext()
}

func withRuntimePoliciesBuildTrackingPayload() loadstrike.TrackingPayload {
	builder := loadstrike.TrackingPayloadBuilder{
		Headers:            map[string]string{"X-Correlation-Id": "ord-1"},
		ContentType:        "application/json",
		MessagePayloadType: "json",
	}
	builder.SetBody(map[string]any{"trackingId": "ord-1"})
	return builder.Build()
}

func withRuntimePoliciesRunResult() loadstrike.LoadStrikeRunResult {
	var result loadstrike.LoadStrikeRunResult
	return result
}

func withRuntimePoliciesScenarioStats() loadstrike.LoadStrikeScenarioStats {
	var stats loadstrike.LoadStrikeScenarioStats
	return stats
}

func withRuntimePoliciesWriteTempConfigFiles() withRuntimePoliciesTempConfigPaths {
	return withRuntimePoliciesTempConfigPaths{
		ConfigPath: "method-reference.loadstrike.config.json",
		InfraPath:  "method-reference.loadstrike.infra.json",
	}
}

        // Attach one runtime policy implementation.
func (reference WithRuntimePoliciesMethodReference) AttachOnePolicyExample() any {
    return withRuntimePoliciesBaseRunner().WithRuntimePolicies(newWithRuntimePoliciesOrdersRuntimePolicy())
}

// Attach multiple runtime policy instances.
func (reference WithRuntimePoliciesMethodReference) AttachTwoPoliciesExample() any {
    return withRuntimePoliciesBaseRunner().WithRuntimePolicies(newWithRuntimePoliciesOrdersRuntimePolicy(), newWithRuntimePoliciesOrdersRuntimePolicy())
}
