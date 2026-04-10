import { chromium } from "playwright";
import { Builder } from "selenium-webdriver";
import {
  AzureEventHubsEndpointDefinition,
  CrossPlatformScenarioConfigurator,
  DatadogReportingSink,
  DatadogReportingSinkOptions,
  DelegateStreamEndpointDefinition,
  GrafanaLokiReportingSink,
  GrafanaLokiReportingSinkOptions,
  HttpEndpointDefinition,
  HttpOAuth2ClientCredentialsOptions,
  InfluxDbReportingSink,
  InfluxDbReportingSinkOptions,
  KafkaEndpointDefinition,
  KafkaSaslOptions,
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
  NatsEndpointDefinition,
  OtelCollectorReportingSink,
  OtelCollectorReportingSinkOptions,
  PushDiffusionEndpointDefinition,
  RabbitMqEndpointDefinition,
  RedisStreamsEndpointDefinition,
  SplunkReportingSink,
  SplunkReportingSinkOptions,
  TimescaleDbReportingSink,
  TimescaleDbReportingSinkOptions,
  TrackingFieldSelector,
  TrackingPayloadBuilder,
} from "@loadstrike/loadstrike-sdk";

/**
 * This feature area focuses on local reports and sink registration so the reporting feature is the main thing being explained.
 */
export async function build() {
  // Sink Identity Hooks is the feature being explained. The sample returns that focused object alongside a minimal run result.
  const sink = new OrdersReportingSink();
  const result = await LoadStrikeRunner.registerScenarios(scenario("sink-identity-hooks")).withRunnerKey(RUNNER_KEY).withTestSuite("orders-reference").withTestName("sink-identity-hooks").withoutReports().run();
  return { sink, result };
}

export const RUNNER_KEY = "runner_dummy_orders_reference";
export const ORDERS_API_BASE_URL = "https://orders.example.test";
export const KAFKA_BOOTSTRAP_SERVERS = "localhost:9092";
export const NATS_SERVER_URL = "nats://localhost:4222";
export const REDIS_CONNECTION_STRING = "localhost:6379,abortConnect=false";
export const RABBIT_HOST = "localhost";
export const AZURE_EVENT_HUBS_CONNECTION_STRING =
  "Endpoint=sb://orders.example.test/;SharedAccessKeyName=RootManageSharedAccessKey;SharedAccessKey=dummy";
export const REPORT_FOLDER = "./artifacts/reports";
export const ORDER_TOPIC = "orders.created";
export const EXAMPLE_ORDER_NUMBER = "ORD-10001";
export const EXAMPLE_TENANT = "demo-tenant";
export const scenarioName = (suffix) => `orders.${suffix}`;

export class OrdersReportingSink {
  sinkName = "orders-sample-sink";
  SinkName = "orders-sample-sink";
  init() {}
  start() {}
  saveRealtimeStats() {}
  saveRealtimeMetrics() {}
  saveRunResult() {}
  stop() {}
}

export function publishStep(context) {
  return LoadStrikeStep.run("publish-order", context, async () =>
    LoadStrikeResponse.ok(
      {
        orderNumber: EXAMPLE_ORDER_NUMBER,
        tenant: EXAMPLE_TENANT,
      },
      "201",
      128,
      "created",
      4.5,
    ),
  );
}

export function scenario(nameSuffix = "publish") {
  return LoadStrikeScenario.create(scenarioName(nameSuffix), async (context) =>
    publishStep(context),
  );
}

export function runner() {
  return LoadStrikeRunner.registerScenarios(scenario())
    .withRunnerKey(RUNNER_KEY)
    .withTestSuite("orders-reference")
    .withTestName("sample-reference")
    .withoutReports();
}

export function builtInSinks() {
  return {
    influxdb: new InfluxDbReportingSink(
      new InfluxDbReportingSinkOptions({
        BaseUrl: "https://influxdb.example.test",
        Organization: "orders-demo",
        Bucket: "orders",
        Token: "dummy-token",
        StaticTags: { tenant: EXAMPLE_TENANT },
      }),
    ),
    grafanaLoki: new GrafanaLokiReportingSink(
      new GrafanaLokiReportingSinkOptions({
        BaseUrl: "https://loki.example.test",
        BearerToken: "dummy-token",
        TenantId: EXAMPLE_TENANT,
        StaticLabels: { service: "orders-api" },
      }),
    ),
    timescaleDb: new TimescaleDbReportingSink(
      new TimescaleDbReportingSinkOptions({
        ConnectionString:
          "Host=localhost;Database=orders;Username=dummy;Password=dummy",
        Schema: "public",
        TableName: "loadstrike_reporting_events",
      }),
    ),
    datadog: new DatadogReportingSink(
      new DatadogReportingSinkOptions({
        BaseUrl: "https://http-intake.logs.datadoghq.com",
        ApiKey: "dummy-api-key",
        ApplicationKey: "dummy-app-key",
        StaticTags: { team: "orders" },
      }),
    ),
    splunk: new SplunkReportingSink(
      new SplunkReportingSinkOptions({
        BaseUrl: "https://splunk.example.test",
        Token: "dummy-hec-token",
        Index: "orders",
        Source: "loadstrike-reference",
      }),
    ),
    otelCollector: new OtelCollectorReportingSink(
      new OtelCollectorReportingSinkOptions({
        BaseUrl: "https://otel.example.test",
        Headers: { authorization: "Bearer dummy" },
        StaticResourceAttributes: { "service.name": "orders-reference" },
      }),
    ),
  };
}
