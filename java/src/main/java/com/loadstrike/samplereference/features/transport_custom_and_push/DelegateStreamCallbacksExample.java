package com.loadstrike.samplereference.features.transport_custom_and_push;

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

public final class DelegateStreamCallbacksExample {
  private DelegateStreamCallbacksExample() {}

  /** This feature area focuses on delegate and push callback contracts for custom transport integrations. */
  public static Object build() {
    // Delegate Stream Callbacks is the feature being explained. The sample returns that focused object alongside a minimal run result.
    var endpoints = delegateEndpoints();
    var result = LoadStrikeRunner.registerScenarios(scenario("delegate-stream-callbacks"))
        .withRunnerKey(OrdersWorkflowPlaceholders.RUNNER_KEY)
        .withTestSuite("orders-reference")
        .withTestName("delegate-stream-callbacks")
        .withoutReports()
        .run();
    return mapOf("endpoints", endpoints, "Result", result);
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

                OrdersWorkflowPlaceholders.scenarioName(suffix), DelegateStreamCallbacksExample::publishStep)

            .withInitAsync(context -> CompletableFuture.completedFuture(null))

            .withCleanAsync(context -> CompletableFuture.completedFuture(null))

            .withWeight(2)

            .withoutWarmUp()

            .withMaxFailCount(3)

            .withRestartIterationOnFail(true);

      }

      public static TrackingPayloadBuilder trackingBuilder() {

        var builder = new TrackingPayloadBuilder();

        builder.headers.put("X-Correlation-Id", OrdersWorkflowPlaceholders.EXAMPLE_ORDER_NUMBER);

        builder.headers.put("X-Tenant", OrdersWorkflowPlaceholders.EXAMPLE_TENANT);

        builder.contentType = "application/json";

        builder.messagePayloadType = Map.class;

        builder.setBody(

            "{\"trackingId\":\""

                + OrdersWorkflowPlaceholders.EXAMPLE_ORDER_NUMBER

                + "\",\"tenantId\":\""

                + OrdersWorkflowPlaceholders.EXAMPLE_TENANT

                + "\"}");

        return builder;

      }

      public static Object delegateEndpoints() {

        var producer = new DelegateStreamEndpointDefinition();

        producer.name = "orders-delegate-producer";

        producer.mode = LoadStrikeTransports.TrafficEndpointMode.Produce;

        producer.trackingField = TrackingFieldSelector.parse("header:x-id");

        producer.produceAsync =

            request -> {

              var result = new LoadStrikeTransports.ProducedMessageResult();

              result.isSuccess = true;

              result.timestampUtcMs = Instant.now().toEpochMilli();

              result.responsePayload = trackingBuilder().build();

              return result;

            };



        var consumer = new DelegateStreamEndpointDefinition();

        consumer.name = "orders-delegate-consumer";

        consumer.mode = LoadStrikeTransports.TrafficEndpointMode.Consume;

        consumer.trackingField = TrackingFieldSelector.parse("json:$.trackingId");

        consumer.consumeAsync =

            () -> {

              var message = new LoadStrikeTransports.ConsumedMessage();

              message.timestampUtcMs = Instant.now().toEpochMilli();

              message.payload = trackingBuilder().build();

              return message;

            };



        return mapOf("Producer", producer, "Consumer", consumer);

      }

      public static Object pushEndpoints() {

        var producer = new PushDiffusionEndpointDefinition();

        producer.name = "orders-push-producer";

        producer.mode = LoadStrikeTransports.TrafficEndpointMode.Produce;

        producer.trackingField = TrackingFieldSelector.parse("header:x-id");

        producer.serverUrl = "ws://push.loadstrike.test";

        producer.topicPath = "/orders";

        producer.publishAsync =

            request -> {

              var result = new LoadStrikeTransports.ProducedMessageResult();

              result.isSuccess = true;

              result.timestampUtcMs = Instant.now().toEpochMilli();

              result.responsePayload = trackingBuilder().build();

              return result;

            };



        var consumer = new PushDiffusionEndpointDefinition();

        consumer.name = "orders-push-consumer";

        consumer.mode = LoadStrikeTransports.TrafficEndpointMode.Consume;

        consumer.trackingField = TrackingFieldSelector.parse("json:$.trackingId");

        consumer.serverUrl = "ws://push.loadstrike.test";

        consumer.topicPath = "/orders";

        consumer.subscribeAsync =

            () -> {

              var message = new LoadStrikeTransports.ConsumedMessage();

              message.timestampUtcMs = Instant.now().toEpochMilli();

              message.payload = trackingBuilder().build();

              return message;

            };



        return mapOf("Producer", producer, "Consumer", consumer);

      }

      public static Map<String, Object> mapOf(Object... values) {

        var result = new LinkedHashMap<String, Object>();

        for (int index = 0; index < values.length; index += 2) {

          result.put(String.valueOf(values[index]), values[index + 1]);

        }

        return result;

      }
}
