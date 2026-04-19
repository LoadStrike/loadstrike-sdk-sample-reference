package com.loadstrike.samplereference.methods.loadstrike_runner;

      import com.loadstrike.runtime.CorrelationStoreConfiguration;
      import com.loadstrike.runtime.CrossPlatformScenarioConfigurator;
      import com.loadstrike.runtime.CrossPlatformTrackingConfiguration;
      import com.loadstrike.runtime.HttpEndpointDefinition;
      import com.loadstrike.runtime.LoadStrikeCorrelation.TrackingFieldSelector;
      import com.loadstrike.runtime.LoadStrikeCorrelation.TrackingPayload;
      import com.loadstrike.runtime.LoadStrikeCorrelation.TrackingPayloadBuilder;
      import com.loadstrike.runtime.LoadStrikeRuntime.*;
      import com.loadstrike.runtime.LoadStrikeTransports;
      import java.time.Duration;
      import java.util.Map;
      import java.util.concurrent.CompletableFuture;
      import java.util.concurrent.CompletionStage;

      public final class CreateMethodReference {
        private CreateMethodReference() {}

        private static final String RUNNER_KEY = "runner_dummy_orders_reference";

private static LoadStrikeScenario baselineScenario() {
  return baselineScenario("orders.get-by-id");
}

private static LoadStrikeScenario baselineScenario(String name) {
  return LoadStrikeScenario.create(name, context -> executeOrderGet(context))
      .withLoadSimulations(LoadStrikeSimulation.iterationsForConstant(1, 1))
      .withoutWarmUp();
}

private static LoadStrikeReply<?> executeOrderGet(LoadStrikeScenarioContext context) {
  return LoadStrikeStep.run("get-order", context, () -> createOrderReply()).asReply();
}

private static CompletionStage<LoadStrikeReply<?>> executeOrderGetAsyncStage(
    LoadStrikeScenarioContext context) {
  return CompletableFuture.completedFuture(executeOrderGet(context));
}

private static LoadStrikeReply<Map<String, Object>> createOrderReply() {
  return LoadStrikeResponse.ok(Map.of("orderId", "ORD-10001"), "200", 128, "ok", 3.2d);
}

private static LoadStrikeRunner baseRunner() {
  return LoadStrikeRunner.create()
      .addScenario(baselineScenario())
      .withRunnerKey(RUNNER_KEY)
      .withTestSuite("orders-reference")
      .withTestName("orders-get-by-id")
      .withoutReports();
}

private static LoadStrikeContext baseContext() {
  return baseRunner().buildContext();
}

private static LoadStrikeContext trackedContext() {
  return LoadStrikeRunner.create()
      .addScenario(
          CrossPlatformScenarioConfigurator.Configure(
              baselineScenario("orders.tracked"), trackingConfiguration()))
      .withRunnerKey(RUNNER_KEY)
      .withoutReports()
      .buildContext();
}

private static HttpEndpointDefinition httpSource() {
  var endpoint = new HttpEndpointDefinition();
  endpoint.name = "orders-http-source";
  endpoint.mode = LoadStrikeTransports.TrafficEndpointMode.Produce;
  endpoint.trackingField = TrackingFieldSelector.parse("header:X-Correlation-Id");
  endpoint.url = "https://orders.example.test/api/orders";
  endpoint.method = "GET";
  endpoint.responseSource = "ResponseBody";
  return endpoint;
}

private static HttpEndpointDefinition httpDestination() {
  var endpoint = new HttpEndpointDefinition();
  endpoint.name = "orders-http-destination";
  endpoint.mode = LoadStrikeTransports.TrafficEndpointMode.Consume;
  endpoint.trackingField = TrackingFieldSelector.parse("json:$.trackingId");
  endpoint.gatherByField = TrackingFieldSelector.parse("json:$.tenantId");
  endpoint.url = "https://orders.example.test/api/order-events";
  endpoint.method = "GET";
  endpoint.responseSource = "ResponseBody";
  endpoint.consumeJsonArrayResponse = true;
  endpoint.consumeArrayPath = "$.items";
  return endpoint;
}

private static CrossPlatformTrackingConfiguration trackingConfiguration() {
  var tracking = new CrossPlatformTrackingConfiguration();
  tracking.source = httpSource();
  tracking.destination = httpDestination();
  tracking.runMode = LoadStrikeTransports.TrackingRunMode.GenerateAndCorrelate;
  tracking.correlationTimeout = Duration.ofSeconds(30);
  tracking.timeoutSweepInterval = Duration.ofSeconds(1);
  tracking.timeoutBatchSize = 200;
  tracking.timeoutCountsAsFailure = true;
  tracking.metricPrefix = "orders_tracking";
  tracking.executeOriginalScenarioRun = false;
  tracking.correlationStore = CorrelationStoreConfiguration.inMemory();
  return tracking;
}

private static TrackingPayload buildTrackingPayload() {
  var builder = new TrackingPayloadBuilder();
  builder.setBody("{\"trackingId\":\"ord-1\"}");
  return builder.build();
}

private static LoadStrikeScenarioContext createStepContext() {
  return null;
}

private static LoadStrikeRunResult runResult() {
  return null;
}

private static LoadStrikeScenarioStats scenarioStats() {
  return null;
}

private static TempConfigPaths writeTempConfigFiles() {
  return new TempConfigPaths("method-reference.loadstrike.config.json", "method-reference.loadstrike.infra.json");
}

private record TempConfigPaths(String ConfigPath, String InfraPath) {}

private static final class OrdersReportingSink implements LoadStrikeReportingSink {
  @Override
  public String sinkName() {
    return "orders-sample-sink";
  }

  @Override
  public void init(LoadStrikeBaseContext context, Map<String, Object> infraConfig) {}

  @Override
  public void start(LoadStrikeSessionStartInfo sessionInfo) {}

  @Override
  public void saveRealtimeStats(LoadStrikeScenarioStats[] scenarioStats) {}

  @Override
  public void saveRealtimeMetrics(LoadStrikeMetricStats metrics) {}

  @Override
  public void saveRunResult(LoadStrikeRunResult result) {}

  @Override
  public void stop() {}

  @Override
  public void dispose() {}
}

private static final class OrdersRuntimePolicy implements LoadStrikeRuntimePolicy {
  @Override
  public boolean shouldRunScenario(String scenarioName) {
    return true;
  }

  @Override
  public void beforeScenario(String scenarioName) {}

  @Override
  public void afterScenario(String scenarioName, LoadStrikeScenarioRuntime stats) {}

  @Override
  public void beforeStep(String scenarioName, String stepName) {}

  @Override
  public void afterStep(String scenarioName, String stepName, LoadStrikeReply<?> reply) {}
}

private static final class OrdersWorkerPlugin implements LoadStrikeWorkerPlugin {
  @Override
  public String pluginName() {
    return "orders-sample-plugin";
  }

  @Override
  public void init() {}

  @Override
  public void start(Map<String, Object> session) {}

  @Override
  public LoadStrikePluginData getData(LoadStrikeRunResult result) {
    return LoadStrikePluginData.create(pluginName());
  }
}

        /** Create an empty runner before scenarios are added. */
public static Object CreateEmptyRunnerExample() {
    return LoadStrikeRunner.create();
}

/** Create a runner and immediately attach the baseline GET scenario. */
public static Object CreateRunnerAndAddScenarioExample() {
    return LoadStrikeRunner.create().addScenario(baselineScenario());
}
      }
