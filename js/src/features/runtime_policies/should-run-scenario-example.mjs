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
 * This feature area focuses on should-run filters, lifecycle callbacks, error modes, and the result surface for runtime policy failures.
 */
export async function build() {
  // Should Run Scenario is the feature being explained. The sample returns that focused object alongside a minimal run result.
  const policy = ordersRuntimePolicy;
  const result = await LoadStrikeRunner.registerScenarios(scenario("should-run-scenario")).withRunnerKey(RUNNER_KEY).withTestSuite("orders-reference").withTestName("should-run-scenario").withoutReports().run();
  return { policy, result };
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

export const ordersRuntimePolicy = {
  policyName: "orders-sample-policy",
  shouldRunScenario: (scenario) => !scenario.endsWith(".skip"),
  beforeScenario() {},
  afterScenario() {},
  beforeStep() {},
  afterStep() {},
};

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
