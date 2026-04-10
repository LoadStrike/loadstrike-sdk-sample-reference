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
import type {
  KafkaSaslOptions,
  LoadStrikeBaseContext,
  LoadStrikeMetricStats,
  LoadStrikeReply,
  LoadStrikeRunResult,
  LoadStrikeRuntimePolicy,
  LoadStrikeScenarioRuntime,
  LoadStrikeScenarioStats,
  LoadStrikeSessionStartInfo,
} from "@loadstrike/loadstrike-sdk";

/**
 * This feature area focuses on building runner contexts once and reusing payload or argument helpers without mixing in unrelated transport setup.
 */
export async function build(): Promise<unknown> {
  // Run Args Reuse is easiest to understand when the runner stays fixed and only the external args change.
  const options = { context: runner().buildContext(), args: ["--sessionId=orders-sample-session"] };
  return await LoadStrikeRunner.Run(options.context, ...options.args);
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
export const scenarioName = (suffix: string) => `orders.${suffix}`;

export function publishStep(context: any) {
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

export function trackingBuilder() {
  const builder = new TrackingPayloadBuilder();
  builder.Headers["X-Correlation-Id"] = EXAMPLE_ORDER_NUMBER;
  builder.Headers["X-Tenant"] = EXAMPLE_TENANT;
  builder.ContentType = "application/json";
  builder.MessagePayloadType = "OrdersPayload";
  builder.setBody(
    JSON.stringify({
      trackingId: EXAMPLE_ORDER_NUMBER,
      tenantId: EXAMPLE_TENANT,
    }),
  );
  return builder;
}
