package com.loadstrike.samplereference.features.runner_builder_and_context;

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

public final class ReusableContextExample {
  private ReusableContextExample() {}

  /** This feature area focuses on creating runners, reusing built contexts, and merging partial configuration. */
  public static Object build() {
    // Reusable Context is configured on a reusable context so you can compare the context flow separately from scenario authoring.
    var context = runner().buildContext();
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

    public static final class OrdersReportingSink implements LoadStrikeReportingSink {

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

    public static final class OrdersRuntimePolicy implements LoadStrikeRuntimePolicy {

      @Override

      public boolean shouldRunScenario(String scenarioName) {

        return scenarioName == null || !scenarioName.endsWith(".skip");

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

    public static final class OrdersWorkerPlugin implements LoadStrikeWorkerPlugin {

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

        var table = LoadStrikePluginDataTable.create("Captured Orders");

        table.headers = List.of("Order Number", "Tenant");

        var row = new LinkedHashMap<String, Object>();

        row.put("Order Number", OrdersWorkflowPlaceholders.EXAMPLE_ORDER_NUMBER);

        row.put("Tenant", OrdersWorkflowPlaceholders.EXAMPLE_TENANT);

        table.rows = List.of(row);



        var pluginData = LoadStrikePluginData.create(pluginName());

        pluginData.hints = List.of("Reference-only plugin payload for the sample repository.");

        pluginData.tables = List.of(table);

        return pluginData;

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

                OrdersWorkflowPlaceholders.scenarioName(suffix), ReusableContextExample::publishStep)

            .withInitAsync(context -> CompletableFuture.completedFuture(null))

            .withCleanAsync(context -> CompletableFuture.completedFuture(null))

            .withWeight(2)

            .withoutWarmUp()

            .withMaxFailCount(3)

            .withRestartIterationOnFail(true);

      }

      public static LoadStrikeRunner runner() {

        return LoadStrikeRunner.create()

            .addScenario(scenario())

            .withTestSuite("orders-reference")

            .withTestName("sample-reference")

            .withSessionId("orders-sample-session")

            .withReportFolder(OrdersWorkflowPlaceholders.REPORT_FOLDER)

            .withReportFileName("orders-reference")

            .withReportFormats(LoadStrikeReportFormat.Txt, LoadStrikeReportFormat.Html)

            .withReportingInterval(2.5d)

            .withoutReports()

            .configure(

                context ->

                    context

                        .withRunnerKey(OrdersWorkflowPlaceholders.RUNNER_KEY)

                        .withLoggerConfig(() -> new LoadStrikeLogger() {})

                        .withMinimumLogLevel("Warning")

                        .withRuntimePolicies(new OrdersRuntimePolicy())

                        .withRuntimePolicyErrorMode(LoadStrikeRuntimePolicyErrorMode.Continue)

                        .withWorkerPlugins(new OrdersWorkerPlugin())

                        .withReportingSinks(new OrdersReportingSink()));

      }

      public static Map<String, Object> mapOf(Object... values) {

        var result = new LinkedHashMap<String, Object>();

        for (int index = 0; index < values.length; index += 2) {

          result.put(String.valueOf(values[index]), values[index + 1]);

        }

        return result;

      }
}
