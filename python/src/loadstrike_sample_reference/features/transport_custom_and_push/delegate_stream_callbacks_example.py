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

FEATURE_OVERVIEW = "This feature area focuses on delegate and push callback contracts for custom transport integrations."


def build():
    # Delegate Stream Callbacks is the feature being explained. The sample returns that focused object alongside a minimal run result.
    endpoints = _custom_endpoints()["delegate_stream"]
    result = (LoadStrikeRunner.register_scenarios(_scenario("delegate-stream-callbacks")).with_runner_key(placeholders.RUNNER_KEY).with_test_suite("orders-reference").with_test_name("delegate-stream-callbacks").without_reports().run())
    return {"endpoints": endpoints, "Result": result}


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


def _tracking_builder():
    builder = TrackingPayloadBuilder()
    builder.with_header("X-Correlation-Id", placeholders.EXAMPLE_ORDER_NUMBER)
    builder.with_json_path("$.tenantId", placeholders.EXAMPLE_TENANT)
    return builder


def _custom_endpoints():
    return {
        "delegate_stream": DelegateStreamEndpointDefinition(
            name="orders-delegate",
            mode="Produce",
            tracking_field=TrackingFieldSelector.parse("header:x-id"),
            produce=lambda payload: payload,
            consume=lambda: {"headers": {"x-id": "delegate-track"}},
        ),
        "push_diffusion": PushDiffusionEndpointDefinition(
            name="orders-push",
            mode="Produce",
            tracking_field=TrackingFieldSelector.parse("header:x-id"),
            server_url="wss://push.example.test",
            topic_path="/orders/events",
            principal="dummy-principal",
            password="dummy-password",
            publish_async=lambda request: request,
            subscribe_async=lambda callback: callback,
        ),
    }
