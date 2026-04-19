import * as loadstrike from "@loadstrike/loadstrike-sdk";

const {
  CorrelationStoreConfiguration,
  CrossPlatformScenarioConfigurator,
  HttpEndpointDefinition,
  LoadStrikeMetric,
  LoadStrikePluginData,
  LoadStrikePluginDataTable,
  LoadStrikeReportFormat,
  LoadStrikeResponse,
  LoadStrikeRunner,
  LoadStrikeScenario,
  LoadStrikeSimulation,
  LoadStrikeStep,
  LoadStrikeThreshold,
  TrackingFieldSelector,
  TrackingPayloadBuilder,
} = loadstrike;

const RUNNER_KEY = "runner_dummy_orders_reference";

class OrdersReportingSink {
  readonly SinkName = "orders-sample-sink";
  init(_context: any, _infraConfig: any) {}
  start(_session: any) {}
  saveRealtimeStats(_stats: any) {}
  saveRealtimeMetrics(_metrics: any) {}
  saveRunResult(_result: any) {}
  stop() {}
  dispose() {}
}

class OrdersRuntimePolicy {
  shouldRunScenario(_scenarioName: string) {
    return true;
  }
  beforeScenario(_scenarioName: string) {}
  afterScenario(_scenarioName: string, _stats: any) {}
  beforeStep(_scenarioName: string, _stepName: string) {}
  afterStep(_scenarioName: string, _stepName: string, _reply: any) {}
}

class OrdersWorkerPlugin {
  readonly PluginName = "orders-sample-plugin";
  init(_context?: any, _infraConfig?: any) {}
  start(_session?: any) {}
  getData(_result: any) {
    return LoadStrikePluginData.create(this.PluginName);
  }
  stop() {}
  dispose() {}
}

function createOrderReply() {
  return LoadStrikeResponse.ok("200", 128, "ok", 3.2);
}

function executeOrderGet(context: any) {
  return LoadStrikeStep.run("get-order", context, async () => createOrderReply());
}

function baselineScenario(name = "orders.get-by-id") {
  return LoadStrikeScenario.create(name, async context => executeOrderGet(context))
    .withLoadSimulations(LoadStrikeSimulation.iterationsForConstant(1, 1))
    .withoutWarmUp();
}

function baseRunner() {
  return LoadStrikeRunner.create()
    .addScenario(baselineScenario())
    .withTestSuite("orders-reference")
    .withTestName("orders-get-by-id")
    .withoutReports();
}

function baseContext() {
  return LoadStrikeRunner.WithRunnerKey(baseRunner().buildContext(), RUNNER_KEY);
}

function httpSource() {
  return new HttpEndpointDefinition({
    Name: "orders-http-source",
    Mode: "Produce",
    TrackingField: TrackingFieldSelector.parse("header:X-Correlation-Id"),
    Url: "https://orders.example.test/api/orders",
    Method: "GET",
    ResponseSource: "Body",
  });
}

function httpDestination() {
  return new HttpEndpointDefinition({
    Name: "orders-http-destination",
    Mode: "Consume",
    TrackingField: TrackingFieldSelector.parse("json:$.trackingId"),
    GatherByField: TrackingFieldSelector.parse("json:$.tenantId"),
    Url: "https://orders.example.test/api/order-events",
    Method: "GET",
    ResponseSource: "Body",
    ConsumeJsonArrayResponse: true,
    ConsumeArrayPath: "$.items",
  });
}

function trackingConfiguration() {
  return {
    Source: httpSource(),
    Destination: httpDestination(),
    RunMode: "GenerateAndCorrelate",
    CorrelationTimeoutSeconds: 30,
    TimeoutSweepIntervalSeconds: 1,
    TimeoutBatchSize: 200,
    TimeoutCountsAsFailure: true,
    MetricPrefix: "orders_tracking",
    ExecuteOriginalScenarioRun: false,
    CorrelationStore: CorrelationStoreConfiguration.inMemory(),
  };
}

function trackedContext() {
  return LoadStrikeRunner.WithRunnerKey(
    LoadStrikeRunner.create()
    .addScenario(
      CrossPlatformScenarioConfigurator.configure(
        baselineScenario("orders.tracked"),
        trackingConfiguration(),
      ),
    )
    .withoutReports()
    .buildContext(),
    RUNNER_KEY,
  );
}

function buildTrackingPayload() {
  const builder = new TrackingPayloadBuilder();
  builder.setBody('{"trackingId":"ord-1"}');
  return builder.build();
}

function runResult(): any {
  return null;
}

function scenarioStats(): any {
  return null;
}

function buildLogger(): any {
  return { logger: "orders" };
}

function writeTempConfigFiles() {
  return {
    configPath: "method-reference.loadstrike.config.json",
    infraPath: "method-reference.loadstrike.infra.json",
  };
}

export class RegisterMetricMethodReference {
  /** Call the public context helper directly from the scenario context surface. */
  static async UseContextMethodExample(): Promise<unknown> {
      return LoadStrikeMetric.counter("orders_total", "count");
  }

  /** Show the same helper in the baseline GET-step flow. */
          static async UseContextMethodInStepExample(): Promise<unknown> {
              const metric = LoadStrikeMetric.counter("orders_total", "count");
  const scenario = LoadStrikeScenario.create("orders.metric", executeOrderGet).withInit(context => context.registerMetric(metric));
  return scenario;
          }
}
