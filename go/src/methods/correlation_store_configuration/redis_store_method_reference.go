package correlation_store_configuration

        import (
        	"time"

        	loadstrike "loadstrike.com/sdk/go"
        )

        const redisStoreRunnerKey = "runner_dummy_orders_reference"

type RedisStoreMethodReference struct{}

type redisStoreTempConfigPaths struct {
	ConfigPath string
	InfraPath  string
}

type redisStoreOrdersReportingSink struct{}

func newRedisStoreOrdersReportingSink() redisStoreOrdersReportingSink {
	return redisStoreOrdersReportingSink{}
}
func (redisStoreOrdersReportingSink) SinkName() string {
	return "orders-sample-sink"
}
func (redisStoreOrdersReportingSink) Init(loadstrike.LoadStrikeBaseContext, loadstrike.IConfiguration) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (redisStoreOrdersReportingSink) Start(loadstrike.LoadStrikeSessionStartInfo) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (redisStoreOrdersReportingSink) SaveRealtimeStats([]loadstrike.LoadStrikeScenarioStats) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (redisStoreOrdersReportingSink) SaveRealtimeMetrics(loadstrike.LoadStrikeMetricStats) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (redisStoreOrdersReportingSink) SaveRunResult(loadstrike.LoadStrikeRunResult) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (redisStoreOrdersReportingSink) Stop() loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (redisStoreOrdersReportingSink) Dispose() {}

type redisStoreOrdersRuntimePolicy struct{}

func newRedisStoreOrdersRuntimePolicy() redisStoreOrdersRuntimePolicy {
	return redisStoreOrdersRuntimePolicy{}
}
func (redisStoreOrdersRuntimePolicy) ShouldRunScenario(string) loadstrike.LoadStrikeBoolTask {
	return loadstrike.TaskFromBool(true)
}
func (redisStoreOrdersRuntimePolicy) BeforeScenario(string) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (redisStoreOrdersRuntimePolicy) AfterScenario(string, loadstrike.LoadStrikeScenarioRuntime) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (redisStoreOrdersRuntimePolicy) BeforeStep(string, string) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (redisStoreOrdersRuntimePolicy) AfterStep(string, string, loadstrike.LoadStrikeReply) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}

type redisStoreOrdersWorkerPlugin struct{}

func newRedisStoreOrdersWorkerPlugin() redisStoreOrdersWorkerPlugin {
	return redisStoreOrdersWorkerPlugin{}
}
func (redisStoreOrdersWorkerPlugin) PluginName() string {
	return "orders-sample-plugin"
}
func (redisStoreOrdersWorkerPlugin) Init(loadstrike.LoadStrikeBaseContext, loadstrike.IConfiguration) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (redisStoreOrdersWorkerPlugin) Start(loadstrike.LoadStrikeSessionStartInfo) loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (redisStoreOrdersWorkerPlugin) GetData(loadstrike.LoadStrikeRunResult) loadstrike.LoadStrikeValueTask[loadstrike.LoadStrikePluginData] {
	return loadstrike.TaskFromResult(loadstrike.CreatePluginData("orders-sample-plugin"))
}
func (redisStoreOrdersWorkerPlugin) Stop() loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}
func (redisStoreOrdersWorkerPlugin) Dispose() loadstrike.LoadStrikeTask {
	return loadstrike.CompletedTask()
}

func redisStorePerformOrderGetReply() loadstrike.LoadStrikeReply {
	return loadstrike.LoadStrikeResponse.Ok("200", int64(128), "ok", loadstrike.TimeSpan(3*time.Millisecond))
}

func redisStoreExecuteOrderGet(context loadstrike.LoadStrikeScenarioContext) loadstrike.LoadStrikeReply {
	return loadstrike.LoadStrikeStep.Run("get-order", context, func(loadstrike.LoadStrikeScenarioContext) loadstrike.LoadStrikeReply {
		return redisStorePerformOrderGetReply()
	})
}

func redisStoreExecuteOrderGetAsync(context loadstrike.LoadStrikeScenarioContext) loadstrike.LoadStrikeValueTask[loadstrike.LoadStrikeReply] {
	return loadstrike.TaskFromResult(redisStoreExecuteOrderGet(context))
}

func redisStoreBaselineScenario(names ...string) loadstrike.LoadStrikeScenario {
	name := "orders.get-by-id"
	if len(names) > 0 && names[0] != "" {
		name = names[0]
	}
	return loadstrike.CreateScenario(name, redisStoreExecuteOrderGet).
		WithLoadSimulations(loadstrike.LoadStrikeSimulation.IterationsForConstant(1, 1)).
		WithoutWarmUp()
}

func redisStoreBaseRunner() loadstrike.LoadStrikeRunner {
	return loadstrike.Create().
		AddScenario(redisStoreBaselineScenario()).
		WithRunnerKey(redisStoreRunnerKey).
		WithTestSuite("orders-reference").
		WithTestName("orders-get-by-id").
		WithoutReports()
}

func redisStoreBaseContext() loadstrike.LoadStrikeContext {
	return redisStoreBaseRunner().BuildContext()
}

func redisStoreHttpSource() *loadstrike.EndpointSpec {
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

func redisStoreHttpDestination() *loadstrike.EndpointSpec {
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

func redisStoreTrackingConfiguration() *loadstrike.TrackingConfigurationSpec {
	return &loadstrike.TrackingConfigurationSpec{
		Source:                      redisStoreHttpSource(),
		Destination:                 redisStoreHttpDestination(),
		RunMode:                     "GenerateAndCorrelate",
		CorrelationTimeoutSeconds:   30,
		TimeoutSweepIntervalSeconds: 1,
		TimeoutBatchSize:            200,
		TimeoutCountsAsFailure:      true,
		MetricPrefix:                "orders_tracking",
		ExecuteOriginalScenarioRun:  false,
	}
}

func redisStoreTrackedContext() loadstrike.LoadStrikeContext {
	return loadstrike.Create().
		AddScenario(redisStoreBaselineScenario("orders.tracked").WithCrossPlatformTracking(redisStoreTrackingConfiguration())).
		WithRunnerKey(redisStoreRunnerKey).
		WithoutReports().
		BuildContext()
}

func redisStoreBuildTrackingPayload() loadstrike.TrackingPayload {
	builder := loadstrike.TrackingPayloadBuilder{
		Headers:            map[string]string{"X-Correlation-Id": "ord-1"},
		ContentType:        "application/json",
		MessagePayloadType: "json",
	}
	builder.SetBody(map[string]any{"trackingId": "ord-1"})
	return builder.Build()
}

func redisStoreRunResult() loadstrike.LoadStrikeRunResult {
	var result loadstrike.LoadStrikeRunResult
	return result
}

func redisStoreScenarioStats() loadstrike.LoadStrikeScenarioStats {
	var stats loadstrike.LoadStrikeScenarioStats
	return stats
}

func redisStoreWriteTempConfigFiles() redisStoreTempConfigPaths {
	return redisStoreTempConfigPaths{
		ConfigPath: "method-reference.loadstrike.config.json",
		InfraPath:  "method-reference.loadstrike.infra.json",
	}
}

        // Call the correlation helper directly with a concrete value.
        func (reference RedisStoreMethodReference) CreateCorrelationExample() any {
            options := loadstrike.RedisCorrelationStoreOptions{
	ConnectionString: "redis://localhost:6379/1",
	Database:         1,
	KeyPrefix:        "orders",
	EntryTTLSeconds:  300,
}
return loadstrike.CorrelationStoreConfiguration{}.RedisStore(options)
        }

// Show how the same helper fits into the tracked source/destination example.
func (reference RedisStoreMethodReference) UseCorrelationExampleInTrackedFlow() any {
    return redisStoreTrackingConfiguration()
}
