package com.loadstrike.samplereference.features.tracking_basics;

import com.loadstrike.runtime.AzureEventHubsEndpointDefinition;
import com.loadstrike.runtime.CorrelationStoreConfiguration;
import com.loadstrike.runtime.CrossPlatformScenarioConfigurator;
import com.loadstrike.runtime.CrossPlatformTrackingConfiguration;
import com.loadstrike.runtime.DelegateStreamEndpointDefinition;
import com.loadstrike.runtime.HttpAuthOptions;
import com.loadstrike.runtime.HttpEndpointDefinition;
import com.loadstrike.runtime.HttpOAuth2ClientCredentialsOptions;
import com.loadstrike.runtime.KafkaEndpointDefinition;
import com.loadstrike.runtime.KafkaSaslOptions;
import com.loadstrike.runtime.LoadStrikeCorrelation.TrackingFieldSelector;
import com.loadstrike.runtime.LoadStrikeCorrelation.TrackingPayloadBuilder;
import com.loadstrike.runtime.LoadStrikeRuntime.LoadStrikeBaseContext;
import com.loadstrike.runtime.LoadStrikeRuntime.LoadStrikeContext;
import com.loadstrike.runtime.LoadStrikeRuntime.LoadStrikeCounter;
import com.loadstrike.runtime.LoadStrikeRuntime.LoadStrikeGauge;
import com.loadstrike.runtime.LoadStrikeRuntime.LoadStrikeLogger;
import com.loadstrike.runtime.LoadStrikeRuntime.LoadStrikeMetric;
import com.loadstrike.runtime.LoadStrikeRuntime.LoadStrikeMetricStats;
import com.loadstrike.runtime.LoadStrikeRuntime.LoadStrikeNodeType;
import com.loadstrike.runtime.LoadStrikeRuntime.LoadStrikePluginData;
import com.loadstrike.runtime.LoadStrikeRuntime.LoadStrikePluginDataTable;
import com.loadstrike.runtime.LoadStrikeRuntime.LoadStrikeReply;
import com.loadstrike.runtime.LoadStrikeRuntime.LoadStrikeReportFormat;
import com.loadstrike.runtime.LoadStrikeRuntime.LoadStrikeReportingSink;
import com.loadstrike.runtime.LoadStrikeRuntime.LoadStrikeResponse;
import com.loadstrike.runtime.LoadStrikeRuntime.LoadStrikeRunResult;
import com.loadstrike.runtime.LoadStrikeRuntime.LoadStrikeRunner;
import com.loadstrike.runtime.LoadStrikeRuntime.LoadStrikeRuntimePolicy;
import com.loadstrike.runtime.LoadStrikeRuntime.LoadStrikeRuntimePolicyError;
import com.loadstrike.runtime.LoadStrikeRuntime.LoadStrikeRuntimePolicyErrorMode;
import com.loadstrike.runtime.LoadStrikeRuntime.LoadStrikeScenario;
import com.loadstrike.runtime.LoadStrikeRuntime.LoadStrikeScenarioContext;
import com.loadstrike.runtime.LoadStrikeRuntime.LoadStrikeScenarioRuntime;
import com.loadstrike.runtime.LoadStrikeRuntime.LoadStrikeScenarioStats;
import com.loadstrike.runtime.LoadStrikeRuntime.LoadStrikeSessionStartInfo;
import com.loadstrike.runtime.LoadStrikeRuntime.LoadStrikeSimulation;
import com.loadstrike.runtime.LoadStrikeRuntime.LoadStrikeSinkError;
import com.loadstrike.runtime.LoadStrikeRuntime.LoadStrikeStep;
import com.loadstrike.runtime.LoadStrikeRuntime.LoadStrikeThreshold;
import com.loadstrike.runtime.LoadStrikeRuntime.LoadStrikeWorkerPlugin;
import com.loadstrike.runtime.LoadStrikeSinks.DatadogReportingSink;
import com.loadstrike.runtime.LoadStrikeSinks.DatadogSinkOptions;
import com.loadstrike.runtime.LoadStrikeSinks.GrafanaLokiReportingSink;
import com.loadstrike.runtime.LoadStrikeSinks.GrafanaLokiSinkOptions;
import com.loadstrike.runtime.LoadStrikeSinks.InfluxDbReportingSink;
import com.loadstrike.runtime.LoadStrikeSinks.InfluxDbSinkOptions;
import com.loadstrike.runtime.LoadStrikeSinks.OtelCollectorReportingSink;
import com.loadstrike.runtime.LoadStrikeSinks.OtelCollectorSinkOptions;
import com.loadstrike.runtime.LoadStrikeSinks.SplunkReportingSink;
import com.loadstrike.runtime.LoadStrikeSinks.SplunkSinkOptions;
import com.loadstrike.runtime.LoadStrikeSinks.TimescaleDbReportingSink;
import com.loadstrike.runtime.LoadStrikeSinks.TimescaleDbSinkOptions;
import com.loadstrike.runtime.LoadStrikeTransports;
import com.loadstrike.runtime.NatsEndpointDefinition;
import com.loadstrike.runtime.PushDiffusionEndpointDefinition;
import com.loadstrike.runtime.RabbitMqEndpointDefinition;
import com.loadstrike.runtime.RedisStreamsEndpointDefinition;
import com.microsoft.playwright.BrowserType;
import com.microsoft.playwright.Playwright;
import java.time.Duration;
import java.time.Instant;
import java.util.LinkedHashMap;
import java.util.List;
import java.util.Map;
import java.util.concurrent.CompletableFuture;
import java.util.concurrent.CompletionStage;
import org.openqa.selenium.chrome.ChromeOptions;

public final class ExactCaseSelectorsExample {
  private ExactCaseSelectorsExample() {}

  /** This feature area focuses on selectors, run modes, and cross-platform tracking configuration needed to correlate generated traffic. */
  public static Object build() {
    // Exact Case Selectors is the feature being explained. The sample returns that focused object alongside a minimal run result.
    var trackingConfig = exactCaseTrackingConfiguration();
    var result = LoadStrikeRunner.registerScenarios(scenario("exact-case-selectors"))
        .withRunnerKey(OrdersWorkflowPlaceholders.RUNNER_KEY)
        .withTestSuite("orders-reference")
        .withTestName("exact-case-selectors")
        .withoutReports()
        .run();
    return mapOf("trackingConfig", trackingConfig, "Result", result);
  }

    public static final class OrdersWorkflowPlaceholders {

      private OrdersWorkflowPlaceholders() {}



      public static final String RUNNER_KEY = "runner_dummy_orders_reference";

      public static final String ORDERS_API_BASE_URL = "https://orders.example.test";

      public static final String KAFKA_BOOTSTRAP_SERVERS = "localhost:9092";

      public static final String NATS_SERVER_URL = "nats://localhost:4222";

      public static final String REDIS_CONNECTION_STRING = "localhost:6379,abortConnect=false";

      public static final String RABBIT_HOST = "localhost";

      public static final String AZURE_EVENT_HUBS_CONNECTION_STRING =

          "Endpoint=sb://orders.example.test/;SharedAccessKeyName=RootManageSharedAccessKey;SharedAccessKey=dummy";

      public static final String REPORT_FOLDER = "./artifacts/reports";

      public static final String ORDER_TOPIC = "orders.created";

      public static final String EXAMPLE_ORDER_NUMBER = "ORD-10001";

      public static final String EXAMPLE_TENANT = "demo-tenant";



      public static String scenarioName(String suffix) {

        return "orders." + suffix;

      }

    }

      public static CompletionStage<LoadStrikeReply<?>> publishStep(

          LoadStrikeScenarioContext context) {

        var payload = new LinkedHashMap<String, Object>();

        payload.put("orderNumber", OrdersWorkflowPlaceholders.EXAMPLE_ORDER_NUMBER);

        payload.put("tenant", OrdersWorkflowPlaceholders.EXAMPLE_TENANT);



        return LoadStrikeStep.runAsync(

                "publish-order",

                context,

                () ->

                    CompletableFuture.completedFuture(

                        LoadStrikeResponse.ok(payload, "201", 128, "created", 4.5d)))

            .thenApply(step -> step.AsReply());

      }

      public static LoadStrikeScenario scenario() {

        return scenario("publish");

      }

      public static LoadStrikeScenario scenario(String suffix) {

        return LoadStrikeScenario.createAsync(

                OrdersWorkflowPlaceholders.scenarioName(suffix), ExactCaseSelectorsExample::publishStep)

            .withInitAsync(context -> CompletableFuture.completedFuture(null))

            .withCleanAsync(context -> CompletableFuture.completedFuture(null))

            .withWeight(2)

            .withoutWarmUp()

            .withMaxFailCount(3)

            .withRestartIterationOnFail(true);

      }

      public static HttpEndpointDefinition httpSource() {

        var endpoint = new HttpEndpointDefinition();

        endpoint.name = "orders-http-source";

        endpoint.mode = LoadStrikeTransports.TrafficEndpointMode.Produce;

        endpoint.trackingField = TrackingFieldSelector.parse("header:X-Correlation-Id");

        endpoint.url = OrdersWorkflowPlaceholders.ORDERS_API_BASE_URL + "/api/orders";

        endpoint.method = "POST";

        endpoint.messageHeaders.put(

            "X-Correlation-Id", OrdersWorkflowPlaceholders.EXAMPLE_ORDER_NUMBER);

        endpoint.messageHeaders.put("X-Tenant", OrdersWorkflowPlaceholders.EXAMPLE_TENANT);

        endpoint.messagePayload =

            mapOf(

                "orderNumber", OrdersWorkflowPlaceholders.EXAMPLE_ORDER_NUMBER,

                "tenantId", OrdersWorkflowPlaceholders.EXAMPLE_TENANT);

        endpoint.trackingPayloadSource = LoadStrikeTransports.HttpTrackingPayloadSource.Request;

        endpoint.responseSource = "ResponseBody";



        var auth = new HttpAuthOptions();

        auth.type = LoadStrikeTransports.HttpAuthType.OAuth2ClientCredentials;

        var oauth = new HttpOAuth2ClientCredentialsOptions();

        oauth.tokenEndpoint = OrdersWorkflowPlaceholders.ORDERS_API_BASE_URL + "/oauth/token";

        oauth.clientId = "dummy-client-id";

        oauth.clientSecret = "dummy-client-secret";

        oauth.scopes = List.of("orders.publish");

        auth.oauth2ClientCredentials = oauth;

        endpoint.auth = auth;

        return endpoint;

      }

      public static HttpEndpointDefinition httpDestination() {

        var endpoint = new HttpEndpointDefinition();

        endpoint.name = "orders-http-destination";

        endpoint.mode = LoadStrikeTransports.TrafficEndpointMode.Consume;

        endpoint.trackingField = TrackingFieldSelector.parse("json:$.trackingId");

        endpoint.gatherByField = TrackingFieldSelector.parse("json:$.tenantId");

        endpoint.url = OrdersWorkflowPlaceholders.ORDERS_API_BASE_URL + "/api/order-events";

        endpoint.method = "GET";

        endpoint.responseSource = "ResponseBody";

        endpoint.consumeJsonArrayResponse = true;

        endpoint.consumeArrayPath = "$.items";

        return endpoint;

      }

      public static CrossPlatformTrackingConfiguration exactCaseTrackingConfiguration() {

        var tracking = new CrossPlatformTrackingConfiguration();

        tracking.source = httpSource();

        tracking.destination = httpDestination();

        tracking.runMode = LoadStrikeTransports.TrackingRunMode.GenerateAndCorrelate;

        tracking.trackingFieldValueCaseSensitive = true;

        tracking.gatherByFieldValueCaseSensitive = true;

        tracking.correlationStore = CorrelationStoreConfiguration.inMemory();

        return tracking;

      }

      public static LoadStrikeScenario trackedScenario() {

        var tracking = exactCaseTrackingConfiguration();

        tracking.correlationTimeout = Duration.ofSeconds(45);

        tracking.timeoutSweepInterval = Duration.ofMillis(500);

        tracking.timeoutBatchSize = 100;

        tracking.timeoutCountsAsFailure = true;

        tracking.executeOriginalScenarioRun = true;

        tracking.metricPrefix = "orders_tracking";



        return CrossPlatformScenarioConfigurator.Configure(

            scenario("tracked").withLoadSimulations(LoadStrikeSimulation.iterationsForConstant(1, 1)),

            tracking);

      }

  public static Map<String, Object> mapOf(Object... values) {
        var result = new LinkedHashMap<String, Object>();
        for (int index = 0; index < values.length; index += 2) {
          result.put(String.valueOf(values[index]), values[index + 1]);
        }
        return result;
      }
}
