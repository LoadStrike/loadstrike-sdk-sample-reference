from __future__ import annotations

from typing import Any

from playwright.async_api import async_playwright
from selenium.webdriver import ChromeOptions

from loadstrike_sdk import (
    AzureEventHubsEndpointDefinition,
    CrossPlatformScenarioConfigurator,
    DatadogReportingSink,
    DelegateStreamEndpointDefinition,
    GrafanaLokiReportingSink,
    HttpAuthOptions,
    HttpEndpointDefinition,
    HttpOAuth2ClientCredentialsOptions,
    InfluxDbReportingSink,
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
    PushDiffusionEndpointDefinition,
    RabbitMqEndpointDefinition,
    RedisStreamsEndpointDefinition,
    SplunkReportingSink,
    TimescaleDbReportingSink,
    TrackingFieldSelector,
    TrackingPayloadBuilder,
)

FEATURE_OVERVIEW = "This feature area focuses on Kafka producer and consumer setup, consumer groups, security, SASL, and OAuth bearer branches."


def build():
    # OAuth Bearer Branches is the feature being explained. The sample returns that focused object alongside a minimal run result.
    endpoint = KafkaEndpointDefinition(name="orders-kafka-oauth", mode="Produce", tracking_field=TrackingFieldSelector.parse("header:x-id"), bootstrap_servers=placeholders.KAFKA_BOOTSTRAP_SERVERS, topic=placeholders.ORDER_TOPIC, security_protocol="SaslSsl", sasl=KafkaSaslOptions(mechanism="OAuthBearer", oauth_bearer_token_endpoint_url=f"{placeholders.ORDERS_API_BASE_URL}/oauth/token"))
    result = (LoadStrikeRunner.register_scenarios(_scenario("oauth-bearer-branches")).with_runner_key(placeholders.RUNNER_KEY).with_test_suite("orders-reference").with_test_name("oauth-bearer-branches").without_reports().run())
    return {"endpoint": endpoint, "Result": result}


class placeholders:
    RUNNER_KEY = "runner_dummy_orders_reference"
    ORDERS_API_BASE_URL = "https://orders.example.test"
    KAFKA_BOOTSTRAP_SERVERS = "localhost:9092"
    NATS_SERVER_URL = "nats://localhost:4222"
    REDIS_CONNECTION_STRING = "localhost:6379,abortConnect=false"
    RABBIT_HOST = "localhost"
    AZURE_EVENT_HUBS_CONNECTION_STRING = (
        "Endpoint=sb://orders.example.test/;SharedAccessKeyName=RootManageSharedAccessKey;SharedAccessKey=dummy"
    )
    REPORT_FOLDER = "./artifacts/reports"
    ORDER_TOPIC = "orders.created"
    EXAMPLE_ORDER_NUMBER = "ORD-10001"
    EXAMPLE_TENANT = "demo-tenant"

    @staticmethod
    def scenario_name(suffix: str) -> str:
        return f"orders.{suffix}"


def _step_reply(context: Any):
    return LoadStrikeStep.run(
        "publish-order",
        context,
        lambda: LoadStrikeResponse.ok(
            {
                "orderNumber": placeholders.EXAMPLE_ORDER_NUMBER,
                "tenant": placeholders.EXAMPLE_TENANT,
            },
            "201",
            128,
            "created",
            4.5,
        ),
    ).as_reply()


def _scenario(name_suffix: str = "publish"):
    return (
        LoadStrikeScenario.create(placeholders.scenario_name(name_suffix), _step_reply)
        .with_init(lambda _context: None)
        .with_clean(lambda _context: None)
        .with_weight(2)
        .without_warm_up()
        .with_max_fail_count(3)
        .with_restart_iteration_on_fail(True)
    )


def _kafka_pair():
    producer = KafkaEndpointDefinition(
        name="orders-kafka-producer",
        mode="Produce",
        tracking_field=TrackingFieldSelector.parse("header:x-id"),
        bootstrap_servers=placeholders.KAFKA_BOOTSTRAP_SERVERS,
        topic=placeholders.ORDER_TOPIC,
        security_protocol="SaslSsl",
        sasl=KafkaSaslOptions(
            mechanism="Plain",
            username="dummy-user",
            password="dummy-password",
        ),
    )
    consumer = KafkaEndpointDefinition(
        name="orders-kafka-consumer",
        mode="Consume",
        tracking_field=TrackingFieldSelector.parse("header:x-id"),
        bootstrap_servers=placeholders.KAFKA_BOOTSTRAP_SERVERS,
        topic=placeholders.ORDER_TOPIC,
        consumer_group_id="orders-sample-group",
        start_from_earliest=True,
    )
    return producer, consumer
