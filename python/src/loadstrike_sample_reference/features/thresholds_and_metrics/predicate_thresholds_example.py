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

FEATURE_OVERVIEW = "This feature area focuses on where thresholds attach, how custom metrics are registered, and how threshold results surface in the run output."


def build():
    # Predicate Thresholds keeps the scenario setup focused on one behavior. Change the option below first when comparing suboptions.
    scenario = _scenario_with_threshold("predicate-threshold", LoadStrikeThreshold.create_metric(lambda stats: True, abort_when_error_count=1))
    return (LoadStrikeRunner.register_scenarios(scenario).with_runner_key(placeholders.RUNNER_KEY).with_test_suite("orders-reference").with_test_name("predicate-thresholds").without_reports().run())


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


def _scenario_with_threshold(name_suffix: str, threshold: dict[str, Any]):
    return _scenario(name_suffix).with_thresholds(threshold)
