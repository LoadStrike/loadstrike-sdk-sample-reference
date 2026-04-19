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
  constructor() {
    this.SinkName = "orders-sample-sink";
  }
  init(_context, _infraConfig) {}
  start(_session) {}
  saveRealtimeStats(_stats) {}
  saveRealtimeMetrics(_metrics) {}
  saveRunResult(_result) {}
  stop() {}
  dispose() {}
}

class OrdersRuntimePolicy {
  shouldRunScenario(_scenarioName) {
    return true;
  }
  beforeScenario(_scenarioName) {}
  afterScenario(_scenarioName, _stats) {}
  beforeStep(_scenarioName, _stepName) {}
  afterStep(_scenarioName, _stepName, _reply) {}
}

class OrdersWorkerPlugin {
  constructor() {
    this.PluginName = "orders-sample-plugin";
  }
  init(_context, _infraConfig) {}
  start(_session) {}
  getData(_result) {
    return LoadStrikePluginData.create(this.PluginName);
  }
  stop() {}
  dispose() {}
}

function createOrderReply() {
  return LoadStrikeResponse.ok("200", 128, "ok", 3.2);
}

function executeOrderGet(context) {
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

function runResult() {
  return null;
}

function scenarioStats() {
  return null;
}

function buildLogger() {
  return { logger: "orders" };
}

function writeTempConfigFiles() {
  return {
    configPath: "method-reference.loadstrike.config.json",
    infraPath: "method-reference.loadstrike.infra.json",
  };
}

export class WithClusterIdMethodReference {
  /** Apply the primary string value shown in the sample method reference. */
  static async ApplyPrimaryValueExample() {
      return LoadStrikeRunner.WithClusterId(baseContext(), "orders-cluster");
  }

  /** Show the same method with a second concrete value. */
  static async ApplyAlternateValueExample() {
      return LoadStrikeRunner.WithClusterId(baseContext(), "checkout-cluster");
  }
}
