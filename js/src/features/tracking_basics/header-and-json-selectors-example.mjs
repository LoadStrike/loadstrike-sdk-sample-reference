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
 * This feature area focuses on selectors, run modes, and cross-platform tracking configuration needed to correlate generated traffic.
 */
export async function build() {
  // Header And Json Selectors is the feature being explained. The sample returns that focused object alongside a minimal run result.
  const selectors = { sourceSelector: TrackingFieldSelector.parse("header:X-Correlation-Id"), destinationSelector: TrackingFieldSelector.parse("json:$.trackingId"), source: httpSource(), destination: httpDestination() };
  const result = await LoadStrikeRunner.registerScenarios(scenario("header-and-json-selectors")).withRunnerKey(RUNNER_KEY).withTestSuite("orders-reference").withTestName("header-and-json-selectors").withoutReports().run();
  return { selectors, result };
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

export function httpSource() {
  return new HttpEndpointDefinition({
    Name: "orders-http-source",
    Mode: "Produce",
    TrackingField: TrackingFieldSelector.parse("header:X-Correlation-Id"),
    Url: `${ORDERS_API_BASE_URL}/api/orders`,
    Method: "POST",
    MessageHeaders: { "X-Tenant": EXAMPLE_TENANT },
    MessagePayload: {
      orderNumber: EXAMPLE_ORDER_NUMBER,
      tenantId: EXAMPLE_TENANT,
    },
    Auth: {
      Type: "OAuth2ClientCredentials",
      OAuth2ClientCredentials: new HttpOAuth2ClientCredentialsOptions({
        TokenUrl: `${ORDERS_API_BASE_URL}/oauth/token`,
        ClientId: "dummy-client-id",
        ClientSecret: "dummy-client-secret",
        Scope: "orders.publish",
      }),
    },
    TrackingPayloadSource: "Request",
    ResponseSource: "Body",
  });
}

export function httpDestination() {
  return new HttpEndpointDefinition({
    Name: "orders-http-destination",
    Mode: "Consume",
    TrackingField: TrackingFieldSelector.parse("json:$.trackingId"),
    GatherByField: TrackingFieldSelector.parse("json:$.tenantId"),
    Url: `${ORDERS_API_BASE_URL}/api/order-events`,
    Method: "GET",
    ResponseSource: "Body",
    ConsumeJsonArrayResponse: true,
    ConsumeArrayPath: "$.items",
  });
}

export function trackingConfiguration() {
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
  };
}

export function trackedScenario() {
  return CrossPlatformScenarioConfigurator.configure(
    scenario("tracking"),
    trackingConfiguration(),
  );
}

export function runner() {
  return LoadStrikeRunner.registerScenarios(scenario())
    .withRunnerKey(RUNNER_KEY)
    .withTestSuite("orders-reference")
    .withTestName("sample-reference")
    .withoutReports();
}
