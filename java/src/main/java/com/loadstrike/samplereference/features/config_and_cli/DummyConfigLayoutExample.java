package com.loadstrike.samplereference.features.config_and_cli;

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

public final class DummyConfigLayoutExample {
  private DummyConfigLayoutExample() {}

  /** This feature area focuses on loading config, overriding config with run arguments, and documenting the expected config layout. */
  public static Object build() {
    // Dummy Config Layout keeps the file-based setup visible while the runtime path stays minimal.
    var configLayout = mapOf("ConfigFile", "loadstrike.config.json", "InfraFile", "loadstrike.infra.json", "RunnerKey", OrdersWorkflowPlaceholders.RUNNER_KEY);
    var context = LoadStrikeRunner.registerScenarios(scenario("config-layout"))
        .loadConfig("features/config_and_cli/" + configLayout.get("ConfigFile"))
        .loadInfraConfig("features/config_and_cli/" + configLayout.get("InfraFile"))
        .withRunnerKey(String.valueOf(configLayout.get("RunnerKey")))
        .buildContext();
    return LoadStrikeRunner.Run(context);
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

                OrdersWorkflowPlaceholders.scenarioName(suffix), DummyConfigLayoutExample::publishStep)

            .withInitAsync(context -> CompletableFuture.completedFuture(null))

            .withCleanAsync(context -> CompletableFuture.completedFuture(null))

            .withWeight(2)

            .withoutWarmUp()

            .withMaxFailCount(3)

            .withRestartIterationOnFail(true);

      }

      public static LoadStrikeContext configuredContext() {

        return LoadStrikeRunner.registerScenarios(scenario())

            .loadConfig("features/config_and_cli/loadstrike.config.json")

            .loadInfraConfig("features/config_and_cli/loadstrike.infra.json")

            .withRunnerKey(OrdersWorkflowPlaceholders.RUNNER_KEY)

            .withSessionId("orders-config-sample")

            .buildContext();

      }

      public static Map<String, Object> mapOf(Object... values) {

        var result = new LinkedHashMap<String, Object>();

        for (int index = 0; index < values.length; index += 2) {

          result.put(String.valueOf(values[index]), values[index + 1]);

        }

        return result;

      }
}
