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
 * This feature area focuses on should-run filters, lifecycle callbacks, error modes, and the result surface for runtime policy failures.
 */
export async function build(): Promise<unknown> {
  // Error Mode is configured on a reusable context so you can compare the context flow separately from scenario authoring.
  const context = runner().buildContext().withRuntimePolicies(ordersRuntimePolicy).withRuntimePolicyErrorMode("continue");
  return await LoadStrikeRunner.Run(context);
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

export const ordersRuntimePolicy: LoadStrikeRuntimePolicy = {
  policyName: "orders-sample-policy",
  shouldRunScenario: (scenario) => !scenario.endsWith(".skip"),
  beforeScenario: (_scenario) => {},
  afterScenario: (_scenario, _stats: LoadStrikeScenarioRuntime) => {},
  beforeStep: (_scenario, _step) => {},
  afterStep: (_scenario, _step, _reply: LoadStrikeReply) => {},
};

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
