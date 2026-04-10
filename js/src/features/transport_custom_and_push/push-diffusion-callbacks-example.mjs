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
 * This feature area focuses on delegate and push callback contracts for custom transport integrations.
 */
export async function build() {
  // Push Diffusion Callbacks is the feature being explained. The sample returns that focused object alongside a minimal run result.
  const endpoints = customEndpoints().pushDiffusion;
  const result = await LoadStrikeRunner.registerScenarios(scenario("push-diffusion-callbacks")).withRunnerKey(RUNNER_KEY).withTestSuite("orders-reference").withTestName("push-diffusion-callbacks").withoutReports().run();
  return { endpoints, result };
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

export function trackingBuilder() {
  const builder = new TrackingPayloadBuilder();
  builder.withHeader("X-Correlation-Id", EXAMPLE_ORDER_NUMBER);
  builder.withJsonPath("$.tenantId", EXAMPLE_TENANT);
  return builder;
}

export function customEndpoints() {
  return {
    delegateStream: new DelegateStreamEndpointDefinition({
      Name: "orders-delegate",
      Mode: "Produce",
      TrackingField: TrackingFieldSelector.parse("header:x-id"),
      Produce: (payload) => payload,
      Consume: () => ({ headers: { "x-id": "delegate-track" } }),
    }),
    pushDiffusion: new PushDiffusionEndpointDefinition({
      Name: "orders-push",
      Mode: "Produce",
      TrackingField: TrackingFieldSelector.parse("header:x-id"),
      ServerUrl: "wss://push.example.test",
      TopicPath: "/orders/events",
      Principal: "dummy-principal",
      Password: "dummy-password",
      PublishAsync: async (request) => request,
      SubscribeAsync: async (callback) => callback,
    }),
  };
}

export function runner() {
  return LoadStrikeRunner.registerScenarios(scenario())
    .withRunnerKey(RUNNER_KEY)
    .withTestSuite("orders-reference")
    .withTestName("sample-reference")
    .withoutReports();
}
