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
 * This feature area focuses on Kafka producer and consumer setup, consumer groups, security, SASL, and OAuth bearer branches.
 */
export async function build() {
  // OAuth Bearer Branches is the feature being explained. The sample returns that focused object alongside a minimal run result.
  const endpoint = new KafkaEndpointDefinition({ Name: "orders-kafka-oauth", Mode: "Produce", TrackingField: TrackingFieldSelector.parse("header:x-id"), BootstrapServers: KAFKA_BOOTSTRAP_SERVERS, Topic: ORDER_TOPIC, SecurityProtocol: "SaslSsl", Sasl: { Mechanism: "OAuthBearer", OAuthBearerTokenEndpointUrl: `${ORDERS_API_BASE_URL}/oauth/token` } });
  const result = await LoadStrikeRunner.registerScenarios(scenario("oauth-bearer-branches")).withRunnerKey(RUNNER_KEY).withTestSuite("orders-reference").withTestName("oauth-bearer-branches").withoutReports().run();
  return { endpoint, result };
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

export function kafkaPair() {
  return {
    producer: new KafkaEndpointDefinition({
      Name: "orders-kafka-producer",
      Mode: "Produce",
      TrackingField: TrackingFieldSelector.parse("header:x-id"),
      BootstrapServers: KAFKA_BOOTSTRAP_SERVERS,
      Topic: ORDER_TOPIC,
      SecurityProtocol: "SaslSsl",
      Sasl: new KafkaSaslOptions({
        Mechanism: "Plain",
        Username: "dummy-user",
        Password: "dummy-password",
      }),
    }),
    consumer: new KafkaEndpointDefinition({
      Name: "orders-kafka-consumer",
      Mode: "Consume",
      TrackingField: TrackingFieldSelector.parse("header:x-id"),
      BootstrapServers: KAFKA_BOOTSTRAP_SERVERS,
      Topic: ORDER_TOPIC,
      ConsumerGroupId: "orders-sample-group",
      StartFromEarliest: true,
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
