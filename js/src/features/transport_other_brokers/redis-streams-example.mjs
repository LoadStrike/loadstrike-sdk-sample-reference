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
 * This feature area focuses on the minimal endpoint options for NATS, RabbitMQ, Redis Streams, and Azure Event Hubs.
 */
export async function build() {
  // Redis Streams is the feature being explained. The sample returns that focused object alongside a minimal run result.
  const endpoints = otherBrokers().redisStreams;
  const result = await LoadStrikeRunner.registerScenarios(scenario("redis-streams")).withRunnerKey(RUNNER_KEY).withTestSuite("orders-reference").withTestName("redis-streams").withoutReports().run();
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

export function otherBrokers() {
  return {
    nats: new NatsEndpointDefinition({
      Name: "orders-nats",
      Mode: "Consume",
      TrackingField: TrackingFieldSelector.parse("header:x-id"),
      ServerUrl: NATS_SERVER_URL,
      Subject: "orders.events",
      QueueGroup: "orders-workers",
    }),
    rabbitmq: new RabbitMqEndpointDefinition({
      Name: "orders-rabbitmq",
      Mode: "Produce",
      TrackingField: TrackingFieldSelector.parse("header:x-id"),
      HostName: RABBIT_HOST,
      UserName: "guest",
      Password: "guest",
      QueueName: "orders.queue",
      RoutingKey: "orders.route",
    }),
    redisStreams: new RedisStreamsEndpointDefinition({
      Name: "orders-redis-streams",
      Mode: "Produce",
      TrackingField: TrackingFieldSelector.parse("header:x-id"),
      ConnectionString: REDIS_CONNECTION_STRING,
      StreamKey: "orders-stream",
      MaxLength: 1000,
    }),
    azureEventHubs: new AzureEventHubsEndpointDefinition({
      Name: "orders-event-hub",
      Mode: "Produce",
      TrackingField: TrackingFieldSelector.parse("header:x-id"),
      ConnectionString: AZURE_EVENT_HUBS_CONNECTION_STRING,
      EventHubName: "orders-hub",
      PartitionKey: EXAMPLE_TENANT,
      PartitionCount: 4,
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
